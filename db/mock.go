package db

import (
	"encoding/json"
	"sync"

	"github.com/JormungandrK/microservice-apps-management/app"
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
func (db *DB) GetUserApps(userID string) ([]byte, error) {
	if userID == "internal-error" {
		return nil, goa.ErrInternal("inertnal-server-error")
	}
	if userID == "bad-request-error" {
		return nil, goa.ErrBadRequest("invalid metadata ID")
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
