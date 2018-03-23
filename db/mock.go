package db

import (
	"encoding/json"
	"sync"

	"github.com/Microkubes/microservice-apps-management/app"
	"github.com/goadesign/goa"
)

// DB emulates a database driver using in-memory data structures.
type DB struct {
	sync.Mutex
	apps map[string]*app.AppPayload
}

// New initializes a new "DB" with dummy data.
func New() *DB {
	name := "app-name"
	desc := "Some description"
	domain := "example.com"
	client := &app.AppPayload{
		Name:        name,
		Description: &desc,
		Domain:      &domain,
	}
	return &DB{apps: map[string]*app.AppPayload{"5975c461f9f8eb02aae053f3": client}}
}

// Mock GetApp method
func (db *DB) GetApp(appID string) (*app.Apps, error) {
	if appID == "internal-error" {
		return nil, goa.ErrInternal("inertnal-server-error")
	}
	if appID == "bad-request-error" {
		return nil, goa.ErrBadRequest("invalid metadata ID")
	}

	client, ok := db.apps[appID]
	if !ok {
		return nil, goa.ErrNotFound("app not found")
	}

	res := &app.Apps{
		ID:           "5975c461f9f8eb02aae053f3",
		Name:         client.Name,
		Description:  *client.Description,
		Domain:       *client.Domain,
		Owner:        "ada5c461f9f8eb02aae05zzz",
		RegisteredAt: 1505746311,
	}

	return res, nil
}

// Mock GetUserApps method
func (db *DB) GetMyApps(userID string) ([]byte, error) {
	if userID == "internal-error" {
		return nil, goa.ErrInternal("inertnal-server-error")
	}

	client, ok := db.apps[userID]
	if !ok {
		return nil, goa.ErrNotFound("app not found")
	}

	var apps []map[string]interface{}
	app1 := map[string]interface{}{
		"ID":           "5975c461f9f8eb02aae053f3",
		"Name":         client.Name,
		"Description":  *client.Description,
		"Domain":       *client.Domain,
		"Owner":        "ada5c461f9f8eb02aae05zzz",
		"RegisteredAt": 1505746311,
	}

	apps = append(apps, app1)

	res, err := json.Marshal(apps)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return res, nil
}

// Mock GetUserApps method
func (db *DB) GetUserApps(userID string) ([]byte, error) {
	if userID == "internal-error" {
		return nil, goa.ErrInternal("inertnal-server-error")
	}

	client, ok := db.apps[userID]
	if !ok {
		return nil, goa.ErrNotFound("app not found")
	}

	var apps []map[string]interface{}
	app1 := map[string]interface{}{
		"ID":           "5975c461f9f8eb02aae053f3",
		"Name":         client.Name,
		"Description":  *client.Description,
		"Domain":       *client.Domain,
		"Owner":        "ada5c461f9f8eb02aae05zzz",
		"Secret":       "GO70Gpt-y8tEYq8HFPDrtva7HhbEb10pEdVu4qjLCLoAXpT9q5DnvM7D",
		"RegisteredAt": 1505746311,
	}

	apps = append(apps, app1)

	res, err := json.Marshal(apps)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return res, nil
}

// Mock RegisterApp method
func (db *DB) RegisterApp(payload *app.AppPayload, userID string) (*app.RegApps, error) {
	if userID == "internal-error" {
		return nil, goa.ErrInternal("inertnal-server-error")
	}
	if userID == "bad-request-error" {
		return nil, goa.ErrBadRequest("invalid user ID")
	}

	db.apps["qwe5c461f9f8ebrtaae05zzz"] = payload

	client := &app.RegApps{
		ID:     "5975c461f9f8eb02aae053f3",
		Secret: "some-secret",
	}

	return client, nil

}

// Mock DeleteApp method
func (db *DB) DeleteApp(appID string) error {
	if appID == "internal-error" {
		return goa.ErrInternal("inertnal-server-error")
	}
	if appID == "bad-request-error" {
		return goa.ErrBadRequest("inalid ID")
	}

	if _, ok := db.apps[appID]; ok {
		delete(db.apps, appID)
	} else {
		return goa.ErrNotFound("app not found!")
	}

	return nil
}

// Mock UpdateApp method
func (db *DB) UpdateApp(payload *app.AppPayload, appID string) (*app.Apps, error) {
	if appID == "internal-error" {
		return nil, goa.ErrInternal("inertnal-server-error")
	}
	if appID == "bad-request-error" {
		return nil, goa.ErrBadRequest("invalid user ID")
	}

	if _, ok := db.apps[appID]; ok {
		db.apps[appID] = payload
	} else {
		return nil, goa.ErrNotFound("app not found!")
	}

	res := &app.Apps{
		ID:           appID,
		Name:         payload.Name,
		Description:  *payload.Description,
		Domain:       *payload.Domain,
		Owner:        "ada5c461f9f8eb02aae05zzz",
		RegisteredAt: 1505746311,
	}

	return res, nil
}

// Mock RegenerateSecret method
func (db *DB) RegenerateSecret(appID string) ([]byte, error) {
	if appID == "internal-error" {
		return nil, goa.ErrInternal("inertnal-server-error")
	}
	if appID == "bad-request-error" {
		return nil, goa.ErrBadRequest("invalid user ID")
	}

	if _, ok := db.apps[appID]; !ok {
		return nil, goa.ErrNotFound("app not found!")
	}

	name := "app-name"
	desc := "Some description"
	domain := "example.com"
	client := &app.Apps{
		ID:           appID,
		Name:         name,
		Description:  desc,
		Domain:       domain,
		Owner:        "ada5c461f9f8eb02aae05zzz",
		RegisteredAt: 1505746311,
	}

	res, err := json.Marshal(client)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return res, nil
}

// FindApp tries to find an app with the supplied app ID and secret.
func (db *DB) FindApp(ID, secret string) (*ClientApp, error) {
	return nil, nil
}
