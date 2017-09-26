package main

import (
	"context"
	"testing"

	"github.com/JormungandrK/microservice-apps-management/app"
	"github.com/JormungandrK/microservice-apps-management/app/test"
	"github.com/JormungandrK/microservice-apps-management/db"
	"github.com/JormungandrK/microservice-security/auth"
	"github.com/goadesign/goa"
)

var (
	service       = goa.New("apps-test")
	database      = db.New()
	ctrl          = NewAppsController(service, database)
	ID            = "5975c461f9f8eb02aae053f3"
	notFoundID    = "rrr5c461f9f8eb02aae05zzz"
	badReqID      = "bad-request-error"
	errInternalID = "internal-error"
	name          = "app-name"
	desc          = "Some description"
	domain        = "example.com"
)

var client = &app.AppPayload{
	Name:        name,
	Description: &desc,
	Domain:      &domain,
}

var ctx = context.Background()

// Call generated test helper, this checks that the returned media type is of the
// correct type (i.e. uses view "default") and validates the media type.
// Also, it ckecks the returned status code
func TestGetAppsOK(t *testing.T) {
	_, clientApp := test.GetAppsOK(t, ctx, service, ctrl, ID)

	if clientApp == nil {
		t.Fatal("Nil client app")
	}

	if clientApp.Name != name {
		t.Errorf("Invalid app name, expected %s, got %s", name, clientApp.Name)
	}
}

func TestGetAppsNotFound(t *testing.T) {
	test.GetAppsNotFound(t, ctx, service, ctrl, notFoundID)
}

func TestGetAppsInternalServerError(t *testing.T) {
	test.GetAppsInternalServerError(t, ctx, service, ctrl, errInternalID)
}

func TestGetAppsBadRequest(t *testing.T) {
	test.GetAppsBadRequest(t, ctx, service, ctrl, badReqID)
}

func TestGetMyAppsAppsOK(t *testing.T) {
	authObj := &auth.Auth{UserID: ID}
	ctx = auth.SetAuth(ctx, authObj)
	test.GetMyAppsAppsOK(t, ctx, service, ctrl)
}

func TestGetMyAppsAppsNotFound(t *testing.T) {
	authObj := &auth.Auth{UserID: notFoundID}
	ctx = auth.SetAuth(ctx, authObj)
	test.GetMyAppsAppsNotFound(t, ctx, service, ctrl)
}

func TestGetMyAppsAppsInternalServerError(t *testing.T) {
	authObj := &auth.Auth{UserID: errInternalID}
	ctx = auth.SetAuth(ctx, authObj)
	test.GetMyAppsAppsInternalServerError(t, ctx, service, ctrl)
}

func TestGetUserAppsAppsOK(t *testing.T) {
	test.GetUserAppsAppsOK(t, ctx, service, ctrl, ID)
}

func TestGetUserAppsAppsNotFound(t *testing.T) {
	test.GetUserAppsAppsNotFound(t, ctx, service, ctrl, notFoundID)
}

func TestGetUserAppsAppsInternalServerError(t *testing.T) {
	test.GetUserAppsAppsInternalServerError(t, ctx, service, ctrl, errInternalID)
}

func TestRegisterAppAppsCreated(t *testing.T) {
	authObj := &auth.Auth{UserID: ID}
	ctx = auth.SetAuth(ctx, authObj)
	test.RegisterAppAppsCreated(t, ctx, service, ctrl, client)
}

func TestRegisterAppAppsBadRequest(t *testing.T) {
	authObj := &auth.Auth{UserID: badReqID}
	ctx = auth.SetAuth(ctx, authObj)
	test.RegisterAppAppsBadRequest(t, ctx, service, ctrl, client)
}

func TestRegisterAppAppsInternalServerError(t *testing.T) {
	authObj := &auth.Auth{UserID: errInternalID}
	ctx = auth.SetAuth(ctx, authObj)
	test.RegisterAppAppsInternalServerError(t, ctx, service, ctrl, client)
}

func TestUpdateAppAppsOK(t *testing.T) {
	test.UpdateAppAppsOK(t, ctx, service, ctrl, ID, client)
}

func TestUpdateAppAppsNotFound(t *testing.T) {
	test.UpdateAppAppsNotFound(t, ctx, service, ctrl, notFoundID, client)
}

func TestUpdateAppAppsInternalServerError(t *testing.T) {
	test.UpdateAppAppsInternalServerError(t, ctx, service, ctrl, errInternalID, client)
}

func TestUpdateAppAppsBadRequest(t *testing.T) {
	test.UpdateAppAppsBadRequest(t, ctx, service, ctrl, badReqID, client)
}

func TestRegenerateClientSecretAppsOK(t *testing.T) {
	test.RegenerateClientSecretAppsOK(t, ctx, service, ctrl, ID)
}

func TestRegenerateClientSecretAppsNotFound(t *testing.T) {
	test.RegenerateClientSecretAppsNotFound(t, ctx, service, ctrl, notFoundID)
}

func TestRegenerateClientSecretAppsInternalServerError(t *testing.T) {
	test.RegenerateClientSecretAppsInternalServerError(t, ctx, service, ctrl, errInternalID)
}

func TestRegenerateClientSecretAppsBadRequest(t *testing.T) {
	test.RegenerateClientSecretAppsBadRequest(t, ctx, service, ctrl, badReqID)
}

func TestDeleteAppAppsOK(t *testing.T) {
	test.DeleteAppAppsOK(t, ctx, service, ctrl, ID)
}

func TestDeleteAppAppsNotFound(t *testing.T) {
	test.DeleteAppAppsNotFound(t, ctx, service, ctrl, notFoundID)
}

func TestDeleteAppAppsInternalServerError(t *testing.T) {
	test.DeleteAppAppsInternalServerError(t, ctx, service, ctrl, errInternalID)
}

func TestDeleteAppAppsBadRequest(t *testing.T) {
	test.DeleteAppAppsBadRequest(t, ctx, service, ctrl, badReqID)
}
