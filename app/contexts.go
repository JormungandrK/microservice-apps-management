// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "apps-management": Application Contexts
//
// Command:
// $ goagen
// --design=github.com/Microkubes/microservice-apps-management/design
// --out=$(GOPATH)/src/github.com/Microkubes/microservice-apps-management
// --version=v1.2.0-dirty

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// DeleteAppAppsContext provides the apps deleteApp action context.
type DeleteAppAppsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AppID string
}

// NewDeleteAppAppsContext parses the incoming request URL and body, performs validations and creates the
// context used by the apps controller deleteApp action.
func NewDeleteAppAppsContext(ctx context.Context, r *http.Request, service *goa.Service) (*DeleteAppAppsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := DeleteAppAppsContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAppID := req.Params["appId"]
	if len(paramAppID) > 0 {
		rawAppID := paramAppID[0]
		rctx.AppID = rawAppID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *DeleteAppAppsContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *DeleteAppAppsContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *DeleteAppAppsContext) NotFound(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 404, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *DeleteAppAppsContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}

// GetAppsContext provides the apps get action context.
type GetAppsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AppID string
}

// NewGetAppsContext parses the incoming request URL and body, performs validations and creates the
// context used by the apps controller get action.
func NewGetAppsContext(ctx context.Context, r *http.Request, service *goa.Service) (*GetAppsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := GetAppsContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAppID := req.Params["appId"]
	if len(paramAppID) > 0 {
		rawAppID := paramAppID[0]
		rctx.AppID = rawAppID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *GetAppsContext) OK(r *Apps) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.apps+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *GetAppsContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *GetAppsContext) NotFound(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 404, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *GetAppsContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}

// GetMyAppsAppsContext provides the apps getMyApps action context.
type GetMyAppsAppsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Payload *AppPayload
}

// NewGetMyAppsAppsContext parses the incoming request URL and body, performs validations and creates the
// context used by the apps controller getMyApps action.
func NewGetMyAppsAppsContext(ctx context.Context, r *http.Request, service *goa.Service) (*GetMyAppsAppsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := GetMyAppsAppsContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *GetMyAppsAppsContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotFound sends a HTTP response with status code 404.
func (ctx *GetMyAppsAppsContext) NotFound(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 404, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *GetMyAppsAppsContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}

// GetUserAppsAppsContext provides the apps getUserApps action context.
type GetUserAppsAppsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	UserID string
}

// NewGetUserAppsAppsContext parses the incoming request URL and body, performs validations and creates the
// context used by the apps controller getUserApps action.
func NewGetUserAppsAppsContext(ctx context.Context, r *http.Request, service *goa.Service) (*GetUserAppsAppsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := GetUserAppsAppsContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramUserID := req.Params["userId"]
	if len(paramUserID) > 0 {
		rawUserID := paramUserID[0]
		rctx.UserID = rawUserID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *GetUserAppsAppsContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotFound sends a HTTP response with status code 404.
func (ctx *GetUserAppsAppsContext) NotFound(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 404, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *GetUserAppsAppsContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}

// RegenerateClientSecretAppsContext provides the apps regenerateClientSecret action context.
type RegenerateClientSecretAppsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AppID string
}

// NewRegenerateClientSecretAppsContext parses the incoming request URL and body, performs validations and creates the
// context used by the apps controller regenerateClientSecret action.
func NewRegenerateClientSecretAppsContext(ctx context.Context, r *http.Request, service *goa.Service) (*RegenerateClientSecretAppsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := RegenerateClientSecretAppsContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAppID := req.Params["appId"]
	if len(paramAppID) > 0 {
		rawAppID := paramAppID[0]
		rctx.AppID = rawAppID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *RegenerateClientSecretAppsContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *RegenerateClientSecretAppsContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *RegenerateClientSecretAppsContext) NotFound(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 404, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *RegenerateClientSecretAppsContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}

// RegisterAppAppsContext provides the apps registerApp action context.
type RegisterAppAppsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
}

// NewRegisterAppAppsContext parses the incoming request URL and body, performs validations and creates the
// context used by the apps controller registerApp action.
func NewRegisterAppAppsContext(ctx context.Context, r *http.Request, service *goa.Service) (*RegisterAppAppsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := RegisterAppAppsContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// Created sends a HTTP response with status code 201.
func (ctx *RegisterAppAppsContext) Created(r *RegApps) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.reg.apps+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 201, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *RegisterAppAppsContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *RegisterAppAppsContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}

// UpdateAppAppsContext provides the apps updateApp action context.
type UpdateAppAppsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	AppID   string
	Payload *AppPayload
}

// NewUpdateAppAppsContext parses the incoming request URL and body, performs validations and creates the
// context used by the apps controller updateApp action.
func NewUpdateAppAppsContext(ctx context.Context, r *http.Request, service *goa.Service) (*UpdateAppAppsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := UpdateAppAppsContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramAppID := req.Params["appId"]
	if len(paramAppID) > 0 {
		rawAppID := paramAppID[0]
		rctx.AppID = rawAppID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *UpdateAppAppsContext) OK(r *Apps) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.apps+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *UpdateAppAppsContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *UpdateAppAppsContext) NotFound(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 404, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *UpdateAppAppsContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}

// VerifyAppAppsContext provides the apps verifyApp action context.
type VerifyAppAppsContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Payload *AppCredentialsPayload
}

// NewVerifyAppAppsContext parses the incoming request URL and body, performs validations and creates the
// context used by the apps controller verifyApp action.
func NewVerifyAppAppsContext(ctx context.Context, r *http.Request, service *goa.Service) (*VerifyAppAppsContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := VerifyAppAppsContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *VerifyAppAppsContext) OK(r *Apps) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.apps+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *VerifyAppAppsContext) NotFound(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 404, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *VerifyAppAppsContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}
