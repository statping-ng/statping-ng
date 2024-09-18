package services

import (
	"bytes"
	"context"
        "crypto/x509"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/statping-ng/statping-ng/types/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/emersion/go-imap/client"
	"github.com/statping-ng/statping-ng/types/failures"
	"github.com/statping-ng/statping-ng/types/hits"
	"github.com/statping-ng/statping-ng/utils"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
        "github.com/Ullaakut/nmap/v3"
)

// checkServices will start the checking go routine for each service
func CheckServices() {
	log.Infoln(fmt.Sprintf("Starting monitoring process for %v Services", len(allServices)))
	for _, s := range allServices {
		time.Sleep(50 * time.Millisecond)
		go ServiceCheckQueue(s, true)
	}
}

// CheckQueue is the main go routine for checking a service
func ServiceCheckQueue(s *Service, record bool) {
	s.Start()
	s.Checkpoint = utils.Now()
	s.SleepDuration = (time.Duration(s.Id) * 100) * time.Millisecond

CheckLoop:
	for {
		select {
		case <-s.Running:
			log.Infoln(fmt.Sprintf("Stopping service: %v", s.Name))
			break CheckLoop
		case <-time.After(s.SleepDuration):
			s.CheckService(record)
			s.UpdateStats()
			s.Checkpoint = s.Checkpoint.Add(s.Duration())
			if !s.Online {
				s.SleepDuration = s.Duration()
			} else {
				s.SleepDuration = s.Checkpoint.Sub(time.Now())
			}
		}
	}
}

func parseHost(s *Service) string {
	if s.Type == "tcp" || s.Type == "udp" || s.Type == "grpc" || s.Type == "smtp" || s.Type == "imap" {
		return s.Domain
	} else {
		u, err := url.Parse(s.Domain)
		if err != nil {
			return s.Domain
		}
		return strings.Split(u.Host, ":")[0]
	}
}

// dnsCheck will check the domain name and return a float64 for the amount of time the DNS check took
func dnsCheck(s *Service) (int64, error) {
	var err error
	t1 := utils.Now()
	host := parseHost(s)
	if s.Type == "tcp" || s.Type == "udp" || s.Type == "grpc" || s.Type == "smtp" {
		_, err = net.LookupHost(host)
	} else {
		_, err = net.LookupIP(host)
	}
	if err != nil {
		return 0, err
	}
	return utils.Now().Sub(t1).Microseconds(), err
}

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

// checkIcmp will send a ICMP ping packet to the service
func CheckIcmp(s *Service, record bool) (*Service, error) {
	defer s.updateLastCheck()
	timer := prometheus.NewTimer(metrics.ServiceTimer(s.Name))
	defer timer.ObserveDuration()

	dur, err := utils.Ping(s.Domain, s.Timeout)
	if err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("Could not send ICMP to service %v, %v", s.Domain, err), "lookup")
		}
		return s, err
	}

	s.PingTime = dur
	s.Latency = dur
	s.LastResponse = ""
	s.Online = true
	if record {
		RecordSuccess(s)
	}
	return s, nil
}

// CheckGrpc will check a gRPC service
func CheckGrpc(s *Service, record bool) (*Service, error) {
	defer s.updateLastCheck()
	timer := prometheus.NewTimer(metrics.ServiceTimer(s.Name))
	defer timer.ObserveDuration()

	// Strip URL scheme if present. Eg: https:// , http://
	if strings.Contains(s.Domain, "://") {
		u, err := url.Parse(s.Domain)
		if err != nil {
			// Unable to parse.
			log.Warnln(fmt.Sprintf("GRPC Service: '%s', Unable to parse URL: '%v'", s.Name, s.Domain))
			if record {
				RecordFailure(s, fmt.Sprintf("Unable to parse GRPC domain %v, %v", s.Domain, err), "parse_domain")
			}
		}

		// Set domain as hostname without port number.
		s.Domain = u.Hostname()
	}

	// Calculate DNS check time
	dnsLookup, err := dnsCheck(s)
	if err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("Could not get IP address for GRPC service %v, %v", s.Domain, err), "lookup")
		}
		return s, err
	}

	// Connect to grpc service without TLS certs.
	grpcOption := grpc.WithInsecure()

	// Check if TLS is enabled
	// Upgrade GRPC connection if using TLS
	// Force to connect on HTTP2 with TLS. Needed when using a reverse proxy such as nginx.
	if s.VerifySSL.Bool {
		h2creds := credentials.NewTLS(&tls.Config{NextProtos: []string{"h2"}})
		grpcOption = grpc.WithTransportCredentials(h2creds)
	}

	s.PingTime = dnsLookup
	t1 := utils.Now()
	domain := fmt.Sprintf("%v", s.Domain)
	if s.Port != 0 {
		domain = fmt.Sprintf("%v:%v", s.Domain, s.Port)
		if isIPv6(s.Domain) {
			domain = fmt.Sprintf("[%v]:%v", s.Domain, s.Port)
		}
	}

	// Context will cancel the request when timeout is exceeded.
	// Cancel the context when request is served within the timeout limit.
	timeout := time.Duration(s.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, domain, grpcOption, grpc.WithBlock())
	if err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("Dial Error %v", err), "connection")
		}
		return s, err
	}

	if s.GrpcHealthCheck.Bool {
		// Create a new health check client
		c := healthpb.NewHealthClient(conn)
		in := &healthpb.HealthCheckRequest{}
		res, err := c.Check(ctx, in)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("GRPC Error %v", err), "healthcheck")
			}
			return s, nil
		}

		// Record responses
		s.LastResponse = strings.TrimSpace(res.String())
		s.LastStatusCode = int(res.GetStatus())
	}

	if err := conn.Close(); err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("%v Socket Close Error %v", strings.ToUpper(s.Type), err), "close")
		}
		return s, err
	}

	// Record latency
	s.Latency = utils.Now().Sub(t1).Microseconds()
	s.Online = true

	if s.GrpcHealthCheck.Bool {
		if s.ExpectedStatus != s.LastStatusCode {
			if record {
				RecordFailure(s, fmt.Sprintf("GRPC Service: '%s', Status Code: expected '%v', got '%v'", s.Name, s.ExpectedStatus, s.LastStatusCode), "response_code")
			}
			return s, nil
		}

		if s.Expected.String != s.LastResponse {
			log.Warnln(fmt.Sprintf("GRPC Service: '%s', Response: expected '%v', got '%v'", s.Name, s.Expected.String, s.LastResponse))
			if record {
				RecordFailure(s, fmt.Sprintf("GRPC Response Body '%v' did not match '%v'", s.LastResponse, s.Expected.String), "response_body")
			}
			return s, nil
		}
	}

	if record {
		RecordSuccess(s)
	}

	return s, nil
}

// checkUdp will check a UDP service using nmap
func CheckUdp(s *Service, record bool) (*Service, error) {
        defer s.updateLastCheck()
        timer := prometheus.NewTimer(metrics.ServiceTimer(s.Name))
        defer timer.ObserveDuration()

        dnsLookup, err := dnsCheck(s)
        if err != nil {
                if record {
                        RecordFailure(s, fmt.Sprintf("Could not get IP address for UDP service %v, %v", s.Domain, err), "lookup")
                }
                return s, err
        }
        s.PingTime = dnsLookup
        t1 := utils.Now()

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        scanner, err := nmap.NewScanner(
                ctx,
                nmap.WithTargets(s.Domain),
                nmap.WithPorts(fmt.Sprintf("%d",s.Port)),
                nmap.WithCustomArguments("-Pn"),
                nmap.WithUDPScan(),
        )
        result, _, err := scanner.Run()

        if err != nil {
				return s, err
        }

	active := false
        for _, host := range result.Hosts {
                if len(host.Ports) == 0 || len(host.Addresses) == 0 {
                        continue
                }

                for _, port := range host.Ports {
                        //fmt.Printf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
                        dport := fmt.Sprintf("%d",port.ID)
                        status := fmt.Sprintf("%s",port.Status())
                        if (dport == fmt.Sprintf("%d",s.Port) && strings.Contains(status,"open")) {
		            active = true
                        }
                }
        }
		if !active {
		if record {
		       RecordFailure(s, fmt.Sprintf("port is closed"), "tls")
			   err = errors.New("port is closed")
		}
		       return s, err
		}

        s.Latency = utils.Now().Sub(t1).Microseconds()
        s.LastResponse = ""
        s.Online = true
        if record {
                RecordSuccess(s)
        }
        return s, nil
}

// checkTcp will check a TCP service
func CheckTcp(s *Service, record bool) (*Service, error) {
	defer s.updateLastCheck()
	timer := prometheus.NewTimer(metrics.ServiceTimer(s.Name))
	defer timer.ObserveDuration()

	dnsLookup, err := dnsCheck(s)
	if err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("Could not get IP address for TCP service %v, %v", s.Domain, err), "lookup")
		}
		return s, err
	}
	s.PingTime = dnsLookup
	t1 := utils.Now()
	domain := fmt.Sprintf("%v", s.Domain)
	if s.Port != 0 {
		domain = fmt.Sprintf("%v:%v", s.Domain, s.Port)
		if isIPv6(s.Domain) {
			domain = fmt.Sprintf("[%v]:%v", s.Domain, s.Port)
		}
	}

	tlsConfig, err := s.LoadTLSCert()
	if err != nil {
		log.Errorln(err)
	}

	// test TCP connection if there is no TLS Certificate set
	if s.TLSCert.String == "" {
		conn, err := net.DialTimeout(s.Type, domain, time.Duration(s.Timeout)*time.Second)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("Dial Error: %v", err), "tls")
			}
			return s, err
		}
		defer conn.Close()
	} else {
		// test TCP connection if TLS Certificate was set
		dialer := &net.Dialer{
			KeepAlive: time.Duration(s.Timeout) * time.Second,
			Timeout:   time.Duration(s.Timeout) * time.Second,
		}
		conn, err := tls.DialWithDialer(dialer, s.Type, domain, tlsConfig)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("Dial Error: %v", err), "tls")
			}
			return s, err
		}
		defer conn.Close()
	}

	s.Latency = utils.Now().Sub(t1).Microseconds()
	s.LastResponse = ""
	s.Online = true
	if record {
		RecordSuccess(s)
	}
	return s, nil
}

// checkSmtp will check an SMTP service
func CheckSmtp(s *Service, record bool) (*Service, error) {
	defer s.updateLastCheck()
	timer := prometheus.NewTimer(metrics.ServiceTimer(s.Name))
	defer timer.ObserveDuration()

	dnsLookup, err := dnsCheck(s)
	if err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("Could not get IP address for %s service %v, %v", strings.ToUpper(s.Type), s.Domain, err), "lookup")
		}
		return s, err
	}
	s.PingTime = dnsLookup
	t1 := utils.Now()
	domain := fmt.Sprintf("%v", s.Domain)
	if s.Port != 0 {
		domain = fmt.Sprintf("%v:%v", s.Domain, s.Port)
		if isIPv6(s.Domain) {
			domain = fmt.Sprintf("[%v]:%v", s.Domain, s.Port)
		}
	}

	tlsConfig, err := s.LoadTLSCert()
	if err != nil {
		log.Errorln(err)
	}

	var c *smtp.Client
	var headers []string
	var username, password string
	if s.Headers.Valid {
		headers = strings.Split(s.Headers.String, ",")
	} else {
		headers = nil
	}

	// check if 'Content-Type' header was defined
	for _, header := range headers {
		if len(strings.Split(header, "=")) < 2 {
			continue
		}
		switch strings.ToLower(strings.Split(header, "=")[0]) {
		case "username":
			username = strings.Split(header, "=")[1]
		case "password":
			password = strings.Split(header, "=")[1]
		}
	}

	if s.requiresTLS() || s.TLSCert.String != "" {
		// test TCP connection if TLS Certificate was set
		dialer := &net.Dialer{
			KeepAlive: time.Duration(s.Timeout) * time.Second,
			Timeout:   time.Duration(s.Timeout) * time.Second,
		}
		conn, err := tls.DialWithDialer(dialer, "tcp", domain, tlsConfig)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("Dial Error: %v", err), "tls")
			}
			return s, err
		}
		defer conn.Close()
		c, err = smtp.NewClient(conn, s.Domain)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("%s Connection Error: %v", strings.ToUpper(s.Type), err), s.Type)
			}
			return s, err
		}
	} else {
		// test TCP connection if there is no TLS Certificate set
		conn, err := net.DialTimeout("tcp", domain, time.Duration(s.Timeout)*time.Second)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("Dial Error: %v", err), "tls")
			}
			return s, err
		}
		defer conn.Close()
		c, err = smtp.NewClient(conn, s.Domain)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("%s Connection Error: %v", strings.ToUpper(s.Type), err), s.Type)
			}
			return s, err
		}
	}

	// Auth
	if s.Port != 25 {
		if username == "" || password == "" {
			err = errors.New("no credentials configured")
			if record {
				RecordFailure(s, fmt.Sprintf("%s Authentication Error: %v", strings.ToUpper(s.Type), err), s.Type)
			}
			return s, err
		}

		if err = c.Auth(smtp.PlainAuth("", username, password, s.Domain)); err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("%s Authentication Error: %v", strings.ToUpper(s.Type), err), s.Type)
			}
			return s, err
		}
	}

	s.Latency = utils.Now().Sub(t1).Microseconds()
	s.LastResponse = ""
	s.Online = true
	if record {
		RecordSuccess(s)
	}
	return s, nil
}

func CheckImap(s *Service, record bool) (*Service, error) {
	defer s.updateLastCheck()
	timer := prometheus.NewTimer(metrics.ServiceTimer(s.Name))
	defer timer.ObserveDuration()

	dnsLookup, err := dnsCheck(s)
	if err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("Could not get IP address for %s service %v, %v", strings.ToUpper(s.Type), s.Domain, err), "lookup")
		}
		return s, err
	}
	s.PingTime = dnsLookup
	t1 := utils.Now()
	domain := fmt.Sprintf("%v", s.Domain)
	if s.Port != 0 {
		domain = fmt.Sprintf("%v:%v", s.Domain, s.Port)
		if isIPv6(s.Domain) {
			domain = fmt.Sprintf("[%v]:%v", s.Domain, s.Port)
		}
	}

	tlsConfig, err := s.LoadTLSCert()
	if err != nil {
		log.Errorln(err)
	}

	var headers []string
	var username, password string
	if s.Headers.Valid {
		headers = strings.Split(s.Headers.String, ",")
	} else {
		headers = nil
	}

	// check if 'Content-Type' header was defined
	for _, header := range headers {
		if len(strings.Split(header, "=")) < 2 {
			continue
		}
		switch strings.ToLower(strings.Split(header, "=")[0]) {
		case "username":
			username = strings.Split(header, "=")[1]
		case "password":
			password = strings.Split(header, "=")[1]
		}
	}

	var conn *client.Client
	if s.requiresTLS() || s.TLSCert.String != "" {
		// test TCP connection if TLS Certificate was set
		dialer := &net.Dialer{
			KeepAlive: time.Duration(s.Timeout) * time.Second,
			Timeout:   time.Duration(s.Timeout) * time.Second,
		}
		conn, err = client.DialWithDialerTLS(dialer, domain, tlsConfig)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("Dial Error: %v", err), "tls")
			}
			return s, err
		}
	} else {
		// test TCP connection if there is no TLS Certificate set
		dialer := &net.Dialer{
			KeepAlive: time.Duration(s.Timeout) * time.Second,
			Timeout:   time.Duration(s.Timeout) * time.Second,
		}
		conn, err = client.DialWithDialer(dialer, domain)
		if err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("Dial Error: %v", err), "tls")
			}
			return s, err
		}
	}
	defer conn.Logout()

	// Auth
	if s.Port != 143 {
		if username == "" || password == "" {
			err = errors.New("no credentials configured")
			if record {
				RecordFailure(s, fmt.Sprintf("%s Authentication Error: %v", strings.ToUpper(s.Type), err), s.Type)
			}
			return s, err
		}

		if err = conn.Login(username, password); err != nil {
			if record {
				RecordFailure(s, fmt.Sprintf("%s Authentication Error: %v", strings.ToUpper(s.Type), err), s.Type)
			}
			return s, err
		}
	}

	s.Latency = utils.Now().Sub(t1).Microseconds()
	s.LastResponse = ""
	s.Online = true
	if record {
		RecordSuccess(s)
	}
	return s, nil
}

func (s *Service) updateLastCheck() {
	s.LastCheck = time.Now()
}

// checkHttp will check a HTTP service
func CheckHttp(s *Service, record bool) (*Service, error) {
	defer s.updateLastCheck()
	timer := prometheus.NewTimer(metrics.ServiceTimer(s.Name))
	defer timer.ObserveDuration()

	dnsLookup, err := dnsCheck(s)
	if err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("Could not get IP address for domain %v, %v", s.Domain, err), "lookup")
		}
		return s, err
	}
	s.PingTime = dnsLookup
	t1 := utils.Now()

	timeout := time.Duration(s.Timeout) * time.Second
	var content []byte
	var res *http.Response
	var data *bytes.Buffer
	var headers []string
	contentType := "application/json" // default Content-Type

	if s.Headers.Valid {
		headers = strings.Split(s.Headers.String, ",")
	} else {
		headers = nil
	}

	// check if 'Content-Type' header was defined
	for _, header := range headers {
		if len(strings.Split(header, "=")) < 2 {
			continue
		}
		if strings.Split(header, "=")[0] == "Content-Type" {
			contentType = strings.Split(header, "=")[1]
			break
		}
	}

	if s.Redirect.Bool {
		headers = append(headers, "Redirect=true")
	}

	if s.PostData.String != "" {
		data = bytes.NewBuffer([]byte(s.PostData.String))
	} else {
		data = bytes.NewBuffer(nil)
	}

	// force set Content-Type to 'application/json' if requests are made
	// with POST method
	if s.Method == "POST" && contentType != "application/json" {
		contentType = "application/json"
	}

	customTLS, err := s.LoadTLSCert()
	if err != nil {
		log.Errorln(err)
	}

        if s.ShowSSL.Bool  {
           checkssl, ssldays := TestSSL(s.Domain,30)
           if err != nil {
                         RecordFailure(s, fmt.Sprintf("%s",checkssl), "ssl")
           }
           s.SSLDays = ssldays
        }


	content, res, err = utils.HttpRequest(s.Domain, s.Method, contentType, headers, data, timeout, s.VerifySSL.Bool, customTLS)
	if err != nil {
		if record {
			RecordFailure(s, fmt.Sprintf("HTTP Error %v", err), "request")
		}
		return s, err
	}
	s.Latency = utils.Now().Sub(t1).Microseconds()
	s.LastResponse = string(content)
	s.LastStatusCode = res.StatusCode

	metrics.Gauge("status_code", float64(res.StatusCode), s.Name)

	if s.Expected.String != "" {
		match, err := regexp.MatchString(s.Expected.String, string(content))
		if err != nil {
			log.Warnln(fmt.Sprintf("Service %v expected: %v to match %v", s.Name, string(content), s.Expected.String))
		}
		if !match {
			if record {
				RecordFailure(s, fmt.Sprintf("HTTP Response Body did not match '%v'", s.Expected), "regex")
			}
			return s, err
		}
	}
	if s.ExpectedStatus != res.StatusCode {
		if record {
			RecordFailure(s, fmt.Sprintf("HTTP Status Code %v did not match %v", res.StatusCode, s.ExpectedStatus), "status_code")
		}
		return s, err
	}
	if record {
		RecordSuccess(s)
	}
	s.Online = true
	return s, err
}

// RecordSuccess will create a new 'hit' record in the database for a successful/online service
func RecordSuccess(s *Service) {
	s.LastOnline = utils.Now()
	s.Online = true
	hit := &hits.Hit{
		Service:   s.Id,
		Latency:   s.Latency,
		PingTime:  s.PingTime,
		CreatedAt: utils.Now(),
	}
	if err := hit.Create(); err != nil {
		log.Error(err)
	}
	log.WithFields(utils.ToFields(hit, s)).Infoln(
		fmt.Sprintf("Service #%d '%v' Successful Response: %s | Lookup in: %s | Online: %v | Interval: %d seconds", s.Id, s.Name, humanMicro(hit.Latency), humanMicro(hit.PingTime), s.Online, s.Interval))
	s.LastLookupTime = hit.PingTime
	s.LastLatency = hit.Latency
	metrics.Gauge("online", 1., s.Name, s.Type)
	metrics.Inc("success", s.Name)
	sendSuccess(s)
}

// RecordFailure will create a new 'Failure' record in the database for a offline service
func RecordFailure(s *Service, issue, reason string) {
	s.LastOffline = utils.Now()

	fail := &failures.Failure{
		Service:   s.Id,
		Issue:     issue,
		PingTime:  s.PingTime,
		CreatedAt: utils.Now(),
		ErrorCode: s.LastStatusCode,
		Reason:    reason,
	}
	log.WithFields(utils.ToFields(fail, s)).
		Warnln(fmt.Sprintf("Service %v Failing: %v | Lookup in: %v", s.Name, issue, humanMicro(fail.PingTime)))

	if err := fail.Create(); err != nil {
		log.Error(err)
	}
	s.Online = false
	s.DownText = s.DowntimeText()

	limitOffset := len(s.Failures)
	if len(s.Failures) >= limitFailures {
		limitOffset = limitFailures - 1
	}

	s.Failures = append([]*failures.Failure{fail}, s.Failures[:limitOffset]...)

	metrics.Gauge("online", 0., s.Name, s.Type)
	metrics.Inc("failure", s.Name)
	sendFailure(s, fail)
}

// Check will run checkHttp for HTTP services and checkTcp for TCP services
// if record param is set to true, it will add a record into the database.
func (s *Service) CheckService(record bool) {
	switch s.Type {
	case "http":
		CheckHttp(s, record)
	case "tcp":
		CheckTcp(s, record)
        case "udp":
                CheckUdp(s, record)
	case "grpc":
		CheckGrpc(s, record)
	case "icmp":
		CheckIcmp(s, record)
	case "smtp":
		CheckSmtp(s, record)
	case "imap":
		CheckImap(s, record)
	}
}

func findMin(arr []int) int {
    if len(arr) == 0 {
        return 0 // or handle the empty array case as needed
    }
    min := arr[0]
    for _, value := range arr {
        if value < min {
            min = value
        }
    }
    return min
}

func TestSSL(host string, expiry int) (error, int) {
    errSunsetAlg       := "%s: '%s' (S/N %X) expires after the sunset date for its signature algorithm '%s'."
    type sigAlgSunset struct {
        name      string    // Human readable name of signature algorithm
        sunsetsAt time.Time // Time the algorithm will be sunset
    }
    var sunsetSigAlgs = map[x509.SignatureAlgorithm]sigAlgSunset{
        x509.MD2WithRSA: sigAlgSunset{
                name:      "MD2 with RSA",
                sunsetsAt: time.Now(),
        },
        x509.MD5WithRSA: sigAlgSunset{
                name:      "MD5 with RSA",
                sunsetsAt: time.Now(),
        },
        x509.SHA1WithRSA: sigAlgSunset{
                name:      "SHA1 with RSA",
                sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
        },
        x509.DSAWithSHA1: sigAlgSunset{
                name:      "DSA with SHA1",
                sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
        },
        x509.ECDSAWithSHA1: sigAlgSunset{
                name:      "ECDSA with SHA1",
                sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
        },
        }



        checkSigAlg := true // Verify that non-root certificates are using a good signature algorithm.

        if strings.Contains(host,"http") {
               parsedURL, err := url.Parse(host)
               if err != nil {
                  log.Errorln(fmt.Sprintf("error in ssl check for %s err: %s", host, err))
               }
               host = parsedURL.Hostname()
        }

        if !strings.Contains(host,":") {
            host = host + ":443"
        }
        conn, err := tls.Dial("tcp", host, nil)
        if err != nil {
                log.Errorln(fmt.Sprintf("error in ssl check for %s err: %s", host, err))
                return err,0
        }
        defer conn.Close()
        timeNow := time.Now()
        checkedCerts := make(map[string]struct{})
        expirations := []int{}
        for _, chain := range conn.ConnectionState().VerifiedChains {
                for certNum, cert := range chain {
                        if _, checked := checkedCerts[string(cert.Signature)]; checked {
                                continue
                        }
                        checkedCerts[string(cert.Signature)] = struct{}{}
                        cErrs := []error{}
                        expirations = append(expirations,int(cert.NotAfter.Sub(timeNow).Hours()/24))
                        // Check the expiration.
                        if timeNow.AddDate(0, 0, expiry).After(cert.NotAfter) {
                                expiresIn := int64(cert.NotAfter.Sub(timeNow).Hours())
                                expiresInDays := int(cert.NotAfter.Sub(timeNow).Hours()/24)
                                if expiry >= expiresInDays {
                                  return errors.New(fmt.Sprintf("SSL Expires in: %d days",expiresInDays)),findMin(expirations)
                                }
                                if expiresIn <= 48 {
                                  return errors.New("SSL Expires in 48 hours"),findMin(expirations)
                                }
                        }

                        // Check the signature algorithm, ignoring the root certificate.
                        if alg, exists := sunsetSigAlgs[cert.SignatureAlgorithm]; checkSigAlg && exists && certNum != len(chain)-1 {
                                if cert.NotAfter.Equal(alg.sunsetsAt) || cert.NotAfter.After(alg.sunsetsAt) {
                                        cErrs = append(cErrs, fmt.Errorf(errSunsetAlg, host, cert.Subject.CommonName, cert.SerialNumber, alg.name))
                                }
                        }

                }
        }
        log.Infoln(fmt.Sprintf("Expiration for %s in days: %d", host, findMin(expirations)))
        return nil,findMin(expirations)
}
