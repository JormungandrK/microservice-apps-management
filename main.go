//go:generate goagen bootstrap -d github.com/Microkubes/microservice-apps-management/design

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Microkubes/microservice-apps-management/app"
	pkgdb "github.com/Microkubes/microservice-apps-management/db"
	"github.com/Microkubes/microservice-security/chain"
	"github.com/Microkubes/microservice-security/flow"
	toolscfg "github.com/Microkubes/microservice-tools/config"
	"github.com/Microkubes/microservice-tools/gateway"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/JormungandrK/backends"

)

func main() {
	// Create service
	service := goa.New("apps-management")

	// Load configuration
	gatewayAdminURL, configFile := loadGatewaySettings()

	cfg, err := toolscfg.LoadConfig(configFile)
	if err != nil {
		service.LogError("config", "err", err)
		return
	}

	// Gateway self-registration
	unregisterService := registerMicroservice(gatewayAdminURL, cfg)
	defer unregisterService() // defer the unregister for after main exits

	// Setup apps-management service
	appsService, err := setupAppsService(cfg)
	if err != nil {
		service.LogError("config", err)
	}

	// create security chain
	securityChain, cleanup, err := flow.NewSecurityFromConfig(cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	service.Use(chain.AsGoaMiddleware(securityChain))

	// Mount "apps" controller
	c := NewAppsController(service, appsService)
	app.MountAppsController(service, c)

	// Mount "swagger" controller
	// c2 := NewSwaggerController(service)
	// app.MountSwaggerController(service, c2)

	// Start service
	if err := service.ListenAndServe(fmt.Sprintf(":%d", cfg.Service.MicroservicePort)); err != nil {
		service.LogError("startup", "err", err)
	}

}

func setupRepository(backend backends.Backend) (backends.Repository, error) {
	return backend.DefineRepository("apps-management", backends.RepositoryDefinitionMap{
		"name": "apps-management",
		"indexes": []backends.Index{
			backends.NewUniqueIndex("id"),
			backends.NewUniqueIndex("domain"),
			backends.NewUniqueIndex("name"),
		},
		"hashKey":       "id",
		"readCapacity":  int64(5),
		"writeCapacity": int64(5),
		"GSI": map[string]interface{}{
			"name": map[string]interface{}{
				"readCapacity":  1,
				"writeCapacity": 1,
			},
		},
	})
}

func setupBackend(dbConfig toolscfg.DBConfig) (backends.Backend, backends.BackendManager, error) {
	dbinfoMap := map[string]*toolscfg.DBInfo{}
	dbinfoMap[dbConfig.DBName] = &dbConfig.DBInfo
	backendManager := backends.NewBackendSupport(dbinfoMap)
	backend, err := backendManager.GetBackend(dbConfig.DBName)
	return backend, backendManager, err
}

func setupAppsService(serviceConfig *toolscfg.ServiceConfig) (pkgdb.AppRepository, error) {
	backend, _, err := setupBackend(serviceConfig.DBConfig)
	if err != nil {
		return nil, err
	}

	appsRepo, err := setupRepository(backend)
	if err != nil {
		return nil, err
	}

	return pkgdb.NewAppsService(appsRepo), err
}

func loadGatewaySettings() (string, string) {
	gatewayURL := os.Getenv("API_GATEWAY_URL")
	serviceConfigFile := os.Getenv("SERVICE_CONFIG_FILE")

	if gatewayURL == "" {
		gatewayURL = "http://localhost:8001"
	}
	if serviceConfigFile == "" {
		serviceConfigFile = "./config.json"
	}

	return gatewayURL, serviceConfigFile
}

func registerMicroservice(gatewayAdminURL string, conf *toolscfg.ServiceConfig) func() {
	registration := gateway.NewKongGateway(gatewayAdminURL, &http.Client{}, conf.Service)

	err := registration.SelfRegister()
	if err != nil {
		panic(err)
	}

	return func() {
		registration.Unregister()
	}
}
