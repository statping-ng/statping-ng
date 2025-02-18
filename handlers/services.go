package handlers

import (
	"github.com/gorilla/mux"
	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/types/errors"
	"github.com/statping-ng/statping-ng/types/failures"
	"github.com/statping-ng/statping-ng/types/hits"
	"github.com/statping-ng/statping-ng/types/services"
	"github.com/statping-ng/statping-ng/utils"
	"net/http"
	"encoding/json"
    "time"    
)

type serviceOrder struct {
	Id    int64 `json:"service"`
	Order int   `json:"order"`
}

func findService(r *http.Request) (*services.Service, error) {
	vars := mux.Vars(r)
	id := utils.ToInt(vars["id"])
	servicer, err := services.Find(id)
	if err != nil {
		return nil, err
	}
	if !IsReadAuthenticated(r) {
		return nil, errors.NotAuthenticated
	}
	return servicer, nil
}

func reorderServiceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var newOrder []*serviceOrder
	if err := DecodeJSON(r, &newOrder); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	for _, s := range newOrder {
		service, err := services.Find(s.Id)
		if err != nil {
			sendErrorJson(err, w, r)
			return
		}
		service.Order = s.Order
		service.Update()
	}
	returnJson(newOrder, w, r)
}

func apiServiceHandler(r *http.Request) interface{} {
	srv, err := findService(r)
	if err != nil {
		return err
	}
	srv = srv.UpdateStats()
	return *srv
}

func apiCreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	var service *services.Service
	if err := DecodeJSON(r, &service); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	if err := service.Create(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	go services.ServiceCheckQueue(service, true)

	sendJsonAction(service, "create", w, r)
}

type servicePatchReq struct {
	Online  bool   `json:"online"`
	Issue   string `json:"issue,omitempty"`
	Latency int64  `json:"latency,omitempty"`
}

func apiServicePatchHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	var req servicePatchReq
	if err := DecodeJSON(r, &req); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	service.Online = req.Online
	service.Latency = req.Latency

	issueDefault := "Service was triggered to be offline"
	if req.Issue != "" {
		issueDefault = req.Issue
	}

	if !req.Online {
		services.RecordFailure(service, issueDefault, "trigger")
	} else {
		services.RecordSuccess(service)
	}

	if err := service.Update(); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	sendJsonAction(service, "update", w, r)
}

func apiServiceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := DecodeJSON(r, &service); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := service.Update(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	go service.CheckService(true)
	sendJsonAction(service, "update", w, r)
}

func apiServiceDataHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupQuery, err := database.ParseQueries(r, service.AllHits())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	objs, err := groupQuery.GraphData(database.ByAverage("latency", 1000))
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	returnJson(objs, w, r)
}

func apiServiceFailureDataHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupQuery, err := database.ParseQueries(r, service.AllFailures())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	objs, err := groupQuery.GraphData(database.ByCount)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	returnJson(objs, w, r)
}

func apiServicePingDataHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupQuery, err := database.ParseQueries(r, service.AllHits())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	objs, err := groupQuery.GraphData(database.ByAverage("ping_time", 1000))
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	returnJson(objs, w, r)
}

func apiServiceTimeDataHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupHits, err := database.ParseQueries(r, service.AllHits())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupFailures, err := database.ParseQueries(r, service.AllFailures())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	var allFailures []*failures.Failure
	var allHits []*hits.Hit

	if err := groupHits.Find(&allHits); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	if err := groupFailures.Find(&allFailures); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	uptimeData, err := service.UptimeData(allHits, allFailures)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	returnJson(uptimeData, w, r)
}

func apiServiceHitsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := service.AllHits().DeleteAll(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(service, "delete", w, r)
}

func apiServiceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	err = service.Delete()
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	sendJsonAction(service, "delete", w, r)
}

func apiAllServicesHandler(r *http.Request) interface{} {
	var srvs []services.Service
	for _, v := range services.AllInOrder() {
		srvs = append(srvs, v)
	}
	return srvs
}

func servicesDeleteFailuresHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := service.AllFailures().DeleteAll(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(service, "delete_failures", w, r)
}

func apiServiceFailuresHandler(r *http.Request) interface{} {
	service, err := findService(r)
	if err != nil {
		return err
	}
	var fails []*failures.Failure
	query, err := database.ParseQueries(r, service.AllFailures())
	if err != nil {
		return err
	}
	query.Find(&fails)
	return fails
}

func apiServiceHitsHandler(r *http.Request) interface{} {
	service, err := findService(r)
	if err != nil {
		return err
	}
	var hts []*hits.Hit
	query, err := database.ParseQueries(r, service.AllHits())
	if err != nil {
		return err
	}
	query.Find(&hts)
	return hts
}

// apiServiceOutageHandler forces a service into outage mode by setting
// IsOutageEnabled to true and defining the OutageType (e.g., "Minor", "Major", "Critical").
// It also records the outage event in the failures history.
func apiServiceOutageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the service ID from the URL parameters.
	vars := mux.Vars(r)
	id := utils.ToInt(vars["id"])

	// Find the service using its ID.
	service, err := services.Find(int64(id))
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	// Verify that the request is authenticated.
	if !IsReadAuthenticated(r) {
		sendErrorJson(errors.NotAuthenticated, w, r)
		return
	}

	// Decode the JSON payload.
	var payload struct {
		OutageType string `json:"outage_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	// Update the service to force an outage.
	service.IsOutageEnabled = true
	service.OutageType = payload.OutageType

	// Save the updated service record.
	if err := service.Update(); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	// Record the outage event in the failures history.
	newFailure := &failures.Failure{
		Issue:      "Manual outage set (" + payload.OutageType + ")",
		Method:     "manual", // Indicates this is a manual trigger.
		ErrorCode:  0,
		Service:    service.Id,
		PingTime:   0,
		OutageType: payload.OutageType,
		CreatedAt:  time.Now(),
	}
	// Use the Create method defined in types/failures/database.go
	if err := newFailure.Create(); err != nil {
		utils.Log.Error("Unable to record outage in failures: ", err)
		// Continue even if recording the failure fails.
	}

	// Return the updated service as a JSON response.
	sendJsonAction(service, "update", w, r)
}