package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/handlers"
	"github.com/statping-ng/statping-ng/notifiers"
	"github.com/statping-ng/statping-ng/source"
	"github.com/statping-ng/statping-ng/types/configs"
	"github.com/statping-ng/statping-ng/types/core"
	"github.com/statping-ng/statping-ng/types/metrics"
	"github.com/statping-ng/statping-ng/types/null"
	"github.com/statping-ng/statping-ng/types/services"
	"github.com/statping-ng/statping-ng/utils"
	"os"
	"os/signal"
	"syscall"
	"strconv"
)

var (
	// VERSION stores the current version of Statping
	VERSION string = "dev"
	// COMMIT stores the git commit hash for this version of Statping
	COMMIT  string
	log     = utils.Log.WithField("type", "cmd")
	confgs  *configs.DbConfig
	stopped chan bool
)

func init() {
	stopped = make(chan bool, 1)
	core.New(VERSION, COMMIT)
	utils.InitEnvs()
	utils.Params.Set("VERSION", VERSION)
	utils.Params.Set("COMMIT", COMMIT)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(assetsCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(sassCmd)
	rootCmd.AddCommand(onceCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(systemctlCmd)
	rootCmd.AddCommand(resetCmd)

	parseFlags(rootCmd)
}

// exit will return an error and return an exit code 1 due to this error
func exit(err error) {
	utils.SentryErr(err)
	log.Fatalln(err)
	os.Exit(1)
}

// Close will gracefully stop the database connection, and log file
func Close() {
	utils.CloseLogs()
	confgs.Close()
	fmt.Println("Shutting down Statping")
}

// main will run the Statping application
func main() {
	go Execute()
	<-stopped
	Close()
}

// main will run the Statping application
func start() {
	go sigterm()
	var err error
	if err := source.Assets(); err != nil {
		exit(err)
	}

	utils.VerboseMode = verboseMode

	if err := utils.InitLogs(); err != nil {
		log.Errorf("Statping Log Error: %v\n", err)
	}

	log.Info(fmt.Sprintf("Starting Statping %s", VERSION))

	utils.Params.Set("SERVER_IP", ipAddress)
	utils.Params.Set("SERVER_PORT", port)

	confgs, err = configs.LoadConfigs(configFile)
	if err != nil {
		log.Infoln("Starting in Setup Mode")
		if err = handlers.RunHTTPServer(); err != nil {
			exit(err)
		}
	}

	if err = configs.ConnectConfigs(confgs, true); err != nil {
		exit(err)
	}

	if err = confgs.ResetCore(); err != nil {
		exit(err)
	}

	if err = confgs.DatabaseChanges(); err != nil {
		exit(err)
	}

	if err := confgs.MigrateDatabase(); err != nil {
		exit(err)
	}

	InitKeycloakConfig()

	if err := mainProcess(); err != nil {
		exit(err)
	}
}

// sigterm will attempt to close the database connections gracefully
func sigterm() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	stopped <- true
}

// mainProcess will initialize the Statping application and run the HTTP server
func mainProcess() error {
	if err := InitApp(); err != nil {
		return err
	}

	services.LoadServicesYaml()

	if err := handlers.RunHTTPServer(); err != nil {
		log.Fatalln(err)
		return errors.Wrap(err, "http server")
	}
	return nil
}

// InitApp will start the Statping instance with a valid database connection
// This function will gather all services in database, add/init Notifiers,
// and start the database cleanup routine
func InitApp() error {
	// fetch Core row information about this instance.
	if _, err := core.Select(); err != nil {
		return err
	}
	// init Sentry error monitoring (its useful)
	utils.SentryInit(core.App.AllowReports.Bool)
	// init prometheus metrics
	metrics.InitMetrics()
	// connect each notifier, added them into database if needed
	notifiers.InitNotifiers()
	// select all services in database and store services in a mapping of Service pointers
	if _, err := services.SelectAllServices(true); err != nil {
		return err
	}
	// start routines for each service checking process
	services.CheckServices()
	// start routine to delete old records (failures, hits)
	go database.Maintenance()
	core.App.Setup = true
	core.App.Started = utils.Now()
	return nil
}

// Initialize Keycloak configuration from environment variables.
// Ensure that the following variables are set:
// KEYCLOAK_CLIENT_ID, KEYCLOAK_CLIENT_SECRET, KEYCLOAK_ENDPOINT_AUTH,
// KEYCLOAK_ENDPOINT_TOKEN, KEYCLOAK_ENDPOINT_USERINFO, KEYCLOAK_SCOPES, and KEYCLOAK_IS_OPEN_ID.
//
// Note: Before deploying, configure your Keycloak client with the required mappers:
// - GroupToRoleMapper (Token mapper, Group Membership, priority 0)
// - User Realm Role (roles-mapper, priority 40)
// Also, create and map the `statping-admin` role for admin groups in Keycloak.
// This setup will include a 'roles' array with 'statping-admin' in the userinfo token.

func InitKeycloakConfig() {
	keycloakClientID := utils.Params.GetString("KEYCLOAK_CLIENT_ID")
	keycloakClientSecret := utils.Params.GetString("KEYCLOAK_CLIENT_SECRET")
	keycloakEndpointAuth := utils.Params.GetString("KEYCLOAK_ENDPOINT_AUTH")
	keycloakEndpointToken := utils.Params.GetString("KEYCLOAK_ENDPOINT_TOKEN")
	keycloakEndpointUserinfo := utils.Params.GetString("KEYCLOAK_ENDPOINT_USERINFO")
	keycloakScopes := utils.Params.GetString("KEYCLOAK_SCOPES")
	keycloakIsOpenID := utils.Params.GetString("KEYCLOAK_IS_OPEN_ID")
	domain := utils.Params.GetString("DOMAIN")

	if keycloakClientID != "" && keycloakClientSecret != "" && keycloakEndpointAuth != "" && keycloakEndpointToken != "" && keycloakEndpointUserinfo != "" && keycloakScopes != "" {
		core.App.OAuth.KeycloakClientID = keycloakClientID
		core.App.OAuth.KeycloakClientSecret = keycloakClientSecret
		core.App.OAuth.KeycloakEndpointAuth = keycloakEndpointAuth
		core.App.OAuth.KeycloakEndpointToken = keycloakEndpointToken
		core.App.OAuth.KeycloakEndpointUserinfo = keycloakEndpointUserinfo
		core.App.OAuth.KeycloakScopes = keycloakScopes
		core.App.Domain = domain
		
		// Convert the string value of KEYCLOAK_IS_OPEN_ID to a boolean
		var isOpenID null.NullBool
		if keycloakIsOpenID != "" {
			parsedIsOpenID, err := strconv.ParseBool(keycloakIsOpenID)
			if err != nil {
				log.Errorf("Invalid value for KEYCLOAK_IS_OPEN_ID: %v", err)
				isOpenID = null.NewNullBool(false)
			} else {
				isOpenID = null.NewNullBool(parsedIsOpenID)
			}
		} else {
			isOpenID = null.NewNullBool(false)
		}

		core.App.OAuth.KeycloakIsOpenID = isOpenID
		
		coreInstance := &core.Core{
			OAuth: core.App.OAuth,
			Domain: core.App.Domain,
		}

		updates := map[string]interface{}{
			"domain": coreInstance.Domain,
			"keycloak_client_id": coreInstance.OAuth.KeycloakClientID,
			"keycloak_client_secret": coreInstance.OAuth.KeycloakClientSecret,
			"keycloak_endpoint_auth": coreInstance.OAuth.KeycloakEndpointAuth,
			"keycloak_endpoint_token": coreInstance.OAuth.KeycloakEndpointToken,
			"keycloak_endpoint_userinfo": coreInstance.OAuth.KeycloakEndpointUserinfo,
			"keycloak_is_open_id": coreInstance.OAuth.KeycloakIsOpenID,
			"keycloak_scopes": coreInstance.OAuth.KeycloakScopes,
		}
		
		result := confgs.Db.Table("core").Model(&core.Core{}).Updates(updates)
		log.Infof("Saving Keycloak data to the database: %v", result)		
	} else {
		log.Warn("Missing Keycloak environment variables.")
	}
}