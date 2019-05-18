package db

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Microkubes/backends"
	"github.com/Microkubes/microservice-apps-management/app"
	"github.com/Microkubes/microservice-tools/config"
	"github.com/asaskevich/govalidator"
	"github.com/keitaroinc/goa"
)

// AppsManagementStore defaines the interface for accessing the application data.
type AppsManagementStore interface {
	// GetApp looks up a applications by the app ID.
	GetApp(appID string) (*app.Apps, error)
	GetMyApps(userID string) ([]byte, error)
	GetUserApps(userID string) ([]byte, error)
	RegisterApp(payload *app.AppPayload, userID string) (*app.RegApps, error)
	DeleteApp(appID string) error
	UpdateApp(payload *app.AppPayload, appID string) (*app.Apps, error)
	RegenerateSecret(appID string) ([]byte, error)
	FindApp(id, secret string) (*ClientApp, error)
}

// ClientApp holds the data for a registered application (client).
type ClientApp struct {
	ID           string `json:"id" bson:"_id"`
	Name         string `json:"name" bson:"name"`
	Description  string `json:"description,omitempty" bson:"description"`
	Domain       string `json:"domain,omitempty" bson:"domain"`
	Owner        string `json:"owner" bson:"owner"`
	RegisteredAt int64  `json:"registeredAt" bson:"registeredAt"`
	Secret       string `json:"secret" bson:"secret"`
}

// BackendAppsManagementStore holds a repository for a certain backend.
// Implements the AppsManagementStore interface.
type BackendAppsManagementStore struct {
	repository backends.Repository
}

// GetApp retrieves an application by id
func (c *BackendAppsManagementStore) GetApp(appID string) (*app.Apps, error) {
	res, err := c.repository.GetOne(backends.NewFilter().Match("id", appID), &ClientApp{})
	if err != nil {
		return nil, err
	}

	clientApp := res.(*ClientApp)

	return &app.Apps{
		Description:  clientApp.Description,
		Domain:       clientApp.Domain,
		ID:           appID,
		Name:         clientApp.Name,
		Owner:        clientApp.Owner,
		RegisteredAt: int(clientApp.RegisteredAt),
	}, nil
}

// GetMyApps retrieves applications for current user
func (c *BackendAppsManagementStore) GetMyApps(userID string) ([]byte, error) {
	var typeHint map[string]interface{}
	apps, err := c.repository.GetAll(backends.NewFilter().Match("owner", userID), typeHint, "registeredat", "asc", 0, 0)
	if err != nil {
		return nil, err
	}

	appsSerialized := apps.(*[]*map[string]interface{})
	appsValue := *appsSerialized

	if len(appsValue) == 0 {
		return nil, goa.ErrNotFound("no apps found")
	}

	for _, client := range appsValue {
		clientValue := *client
		delete(clientValue, "secret")
	}

	res, err := json.Marshal(apps)

	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return res, nil
}

// GetUserApps retrieves applications for a user
func (c *BackendAppsManagementStore) GetUserApps(userID string) ([]byte, error) {
	var typeHint map[string]interface{}
	apps, err := c.repository.GetAll(backends.NewFilter().Match("owner", userID), typeHint, "registeredat", "asc", 0, 0)
	if err != nil {
		return nil, err
	}

	appsSerialized := apps.(*[]*map[string]interface{})
	appsValue := *appsSerialized

	if len(appsValue) == 0 {
		return nil, goa.ErrNotFound("no apps found")
	}

	res, err := json.Marshal(apps)

	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return res, nil
}

// RegisterApp creates a new application for a user
func (c *BackendAppsManagementStore) RegisterApp(payload *app.AppPayload, userID string) (*app.RegApps, error) {
	existing, err := c.repository.GetOne(backends.NewFilter().Match("name", payload.Name), &ClientApp{})
	if err != nil && err.Error() != "not found" {
		return nil, goa.ErrInternal(err)
	}
	if existing != nil {
		return nil, goa.ErrBadRequest("that application already exists")
	}

	registeredAt := int64(time.Now().Unix())
	secret, err := GenerateRandomString(42)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}
	if err := validateDomain(*payload.Domain); err != nil {
		return nil, err
	}

	clientApp := &ClientApp{
		Name:         payload.Name,
		Description:  *payload.Description,
		Domain:       *payload.Domain,
		Owner:        userID,
		Secret:       secret,
		RegisteredAt: registeredAt,
	}

	res, err := c.repository.Save(clientApp, nil)

	if err != nil {
		return nil, err
	}

	ca := res.(*ClientApp)
	regApp := &app.RegApps{
		ID:     ca.ID,
		Secret: secret,
	}

	return regApp, nil
}

// DeleteApp deletes an application by id
func (c *BackendAppsManagementStore) DeleteApp(appID string) error {
	err := c.repository.DeleteOne(backends.NewFilter().Match("id", appID))
	if err != nil {
		if err.Error() == "not found" {
			return goa.ErrNotFound("no app found!")
		}
		return goa.ErrInternal(err)
	}

	return nil
}

// UpdateApp updates an application by id
func (c *BackendAppsManagementStore) UpdateApp(payload *app.AppPayload, appID string) (*app.Apps, error) {
	res, err := c.repository.GetOne(backends.NewFilter().Match("id", appID), &ClientApp{})
	if err != nil {
		return nil, err
	}
	existing := res.(*ClientApp)

	existing.Name = payload.Name

	if payload.Description != nil {
		existing.Description = *payload.Description
	}

	if payload.Domain != nil {
		existing.Domain = *payload.Domain
	}

	res, err = c.repository.Save(existing, backends.NewFilter().Match("id", appID))
	if err != nil {
		if err.Error() == "not found" {
			return nil, goa.ErrNotFound("application not found.")
		}
		return nil, goa.ErrInternal(err)
	}

	clientApp := res.(*ClientApp)

	return &app.Apps{
		Description:  clientApp.Description,
		Domain:       clientApp.Domain,
		ID:           appID,
		Name:         clientApp.Name,
		Owner:        clientApp.Owner,
		RegisteredAt: int(clientApp.RegisteredAt),
	}, nil
}

// RegenerateSecret creates a new secret for an application by id
func (c *BackendAppsManagementStore) RegenerateSecret(appID string) ([]byte, error) {
	secret, err := GenerateRandomString(42)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	res, err := c.repository.GetOne(backends.NewFilter().Match("id", appID), &ClientApp{})
	if err != nil {
		return nil, err
	}
	existing := res.(*ClientApp)

	existing.Secret = secret

	client, err := c.repository.Save(existing, backends.NewFilter().Match("id", appID))
	if err != nil {
		if err.Error() == "not found" {
			return nil, goa.ErrNotFound("application not found.")
		}
		return nil, goa.ErrInternal(err)
	}

	resp, err := json.Marshal(client)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return resp, nil
}

// FindApp tries to find an application (client) by its ID and secret.
// Returns nil if no such app is found.
func (c *BackendAppsManagementStore) FindApp(ID, secret string) (*ClientApp, error) {
	clientApp, err := c.repository.GetOne(backends.NewFilter().Match("id", ID), &ClientApp{})
	if err != nil {
		return nil, err
	}

	ca := clientApp.(*ClientApp)

	if ca.Secret == secret {
		return ca, nil
	}

	return nil, nil
}

// NewAppsManagementStore creates new AppsManagementStore implementation that supports multiple backend types.
func NewAppsManagementStore(cfg *config.DBConfig) (store AppsManagementStore, cleanup func(), err error) {
	manager := backends.NewBackendSupport(map[string]*config.DBInfo{
		cfg.DBName: &cfg.DBInfo,
	})

	noop := func() {}
	backend, err := manager.GetBackend(cfg.DBName)
	if err != nil {
		return nil, noop, err
	}
	cleanup = func() {
		backend.Shutdown()
	}

	repo, err := backend.DefineRepository("apps-management", backends.RepositoryDefinitionMap{
		"name": "apps-management",
		"indexes": []backends.Index{
			backends.NewUniqueIndex("id"),
			backends.NewUniqueIndex("name"),
			backends.NewNonUniqueIndex("registeredAt"),
		},
		"hashKey":       "id",
		"rangeKey":      "name",
		"readCapacity":  10,
		"writeCapacity": 10,
		"GSI": map[string]interface{}{
			"name": map[string]interface{}{
				"readCapacity":  1,
				"writeCapacity": 1,
			},
		},
	})
	if err != nil {
		return nil, noop, err
	}

	store = &BackendAppsManagementStore{
		repository: repo,
	}

	return store, cleanup, err
}

// Convert hex representation of object id to bson object id
func hexToObjectID(hexID string) (bson.ObjectId, error) {
	// Return whether userID is a valid hex representation of object id.
	if bson.IsObjectIdHex(hexID) != true {
		return "", goa.ErrBadRequest("invalid ID")
	}

	// Return an ObjectId from the provided hex representation.
	objectID := bson.ObjectIdHex(hexID)

	// Return true if objectID is valid. A valid objectID must contain exactly 12 bytes.
	if objectID.Valid() != true {
		return "", goa.ErrInternal("invalid object ID")
	}

	return objectID, nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded securely generated random string.
// It will return an error if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Validate domain
func validateDomain(domain string) error {
	if ok := govalidator.IsURL(domain); !ok {
		return goa.ErrBadRequest("invalid domain name")
	}

	return nil
}
