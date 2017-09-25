//go:generate goagen bootstrap -d github.com/JormungandrK/microservice-apps-management/design

package main

import (
	"net/http"
	"os"

	"github.com/JormungandrK/microservice-apps-management/app"
	"github.com/JormungandrK/microservice-apps-management/db"
	"github.com/JormungandrK/microservice-tools/gateway"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	// Gateway self-registration
	unregisterService := registerMicroservice()
	defer unregisterService() // defer the unregister for after main exits

	// Create service
	service := goa.New("apps-management")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Load MongoDB ENV variables
	host, username, password, database := loadMongnoSettings()
	// Create new session to MongoDB
	session := db.NewSession(host, username, password, database)

	// At the end close session
	defer session.Close()

	// Create apps collection and indexes
	index1 := []string{"domain"}
	index2 := []string{"name"}
	indexes := [][]string{index1, index2}
	collectionName := "apps"
	collection := db.PrepareDB(session, database, collectionName, indexes)

	// Mount "apps" controller
	c := NewAppsController(service, &db.MongoCollection{collection})
	app.MountAppsController(service, c)
	// Mount "swagger" controller
	c2 := NewSwaggerController(service)
	app.MountSwaggerController(service, c2)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}

func loadMongnoSettings() (string, string, string, string) {
	host := os.Getenv("MONGO_URL")
	username := os.Getenv("MS_USERNAME")
	password := os.Getenv("MS_PASSWORD")
	database := os.Getenv("MS_DBNAME")

	if host == "" {
		host = "127.0.0.1:27017"
	}
	if username == "" {
		username = "restapi"
	}
	if password == "" {
		password = "restapi"
	}
	if database == "" {
		database = "apps-management"
	}

	return host, username, password, database
}

func loadGatewaySettings() (string, string) {
	gatewayURL := os.Getenv("API_GATEWAY_URL")
	serviceConfigFile := os.Getenv("SERVICE_CONFIG_FILE")

	if gatewayURL == "" {
		gatewayURL = "http://localhost:8001"
	}
	if serviceConfigFile == "" {
		serviceConfigFile = "config.json"
	}

	return gatewayURL, serviceConfigFile
}

func registerMicroservice() func() {
	gatewayURL, configFile := loadGatewaySettings()
	registration, err := gateway.NewKongGatewayFromConfigFile(gatewayURL, &http.Client{}, configFile)
	if err != nil {
		panic(err)
	}
	err = registration.SelfRegister()
	if err != nil {
		panic(err)
	}

	return func() {
		registration.Unregister()
	}
}
