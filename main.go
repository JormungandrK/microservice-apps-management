//go:generate goagen bootstrap -d github.com/Microkubes/microservice-apps-management/design

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Microkubes/microservice-apps-management/app"
	"github.com/Microkubes/microservice-apps-management/db"
	"github.com/Microkubes/microservice-security/chain"
	"github.com/Microkubes/microservice-security/flow"
	"github.com/Microkubes/microservice-tools/config"
	"github.com/Microkubes/microservice-tools/utils/healthcheck"
	"github.com/Microkubes/microservice-tools/utils/version"
	"github.com/keitaroinc/goa"
	"github.com/keitaroinc/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("apps-management")

	// Load configuration
	configFile := loadConfigSettings()

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		service.LogError("config", "err", err)
		return
	}

	// create security chain
	securityChain, cleanup, err := flow.NewSecurityFromConfig(conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	store, cleanup, err := db.NewAppsManagementStore(&conf.DBConfig)
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	defer cleanup()

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	service.Use(chain.AsGoaMiddleware(securityChain))

	service.Use(healthcheck.NewCheckMiddleware("/healthcheck"))

	service.Use(version.NewVersionMiddleware(conf.Version, "/version"))

	// Mount "apps" controller
	c := NewAppsController(service, store)
	app.MountAppsController(service, c)
	// Mount "swagger" controller
	c2 := NewSwaggerController(service)
	app.MountSwaggerController(service, c2)

	// Start service
	if err := service.ListenAndServe(fmt.Sprintf(":%d", conf.Service.MicroservicePort)); err != nil {
		service.LogError("startup", "err", err)
	}

}

func loadConfigSettings() string {
	serviceConfigFile := os.Getenv("SERVICE_CONFIG_FILE")

	if serviceConfigFile == "" {
		serviceConfigFile = "/run/secrets/microservice_apps_management_config.json"
	}

	return serviceConfigFile
}
