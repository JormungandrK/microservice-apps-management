package main

import (
	"fmt"

	"github.com/Microkubes/microservice-apps-management/app"
	"github.com/Microkubes/microservice-apps-management/db"
	"github.com/Microkubes/microservice-security/auth"
	"github.com/goadesign/goa"
	errors "github.com/JormungandrK/backends"
)

// AppsController implements the apps resource.
type AppsController struct {
	*goa.Controller
	Repository db.AppRepository
}

// NewAppsController creates a apps controller.
func NewAppsController(service *goa.Service, repository db.AppRepository) *AppsController {
	return &AppsController{
		Controller: service.NewController("AppsController"),
		Repository: repository,
	}
}

// Get returns an app by its ID.
func (c *AppsController) Get(ctx *app.GetAppsContext) error {
	res, err := c.Repository.GetApp(ctx.AppID)

	if err != nil {
		if errors.IsErrNotFound(err) {
			return ctx.NotFound(err)
		}
		
		if errors.IsErrAlreadyExists(err) {
			return ctx.BadRequest(err)
		}
		
		return ctx.InternalServerError(err)
	}

	return ctx.OK(res)
}

// GetMyApps returns a paginated list of all apps for the current user.
func (c *AppsController) GetMyApps(ctx *app.GetMyAppsAppsContext) error {
	var authObj *auth.Auth
	hasAuth := auth.HasAuth(ctx)

	if hasAuth {
		authObj = auth.GetAuth(ctx.Context)
	} else {
		return ctx.InternalServerError(goa.ErrInternal("Auth has not been set"))
	}

	userID := authObj.UserID

	res, err := c.Repository.GetMyApps(userID)

	if err != nil {
		if errors.IsErrNotFound(err) {
			return ctx.NotFound(err)
		}
		
		return ctx.InternalServerError(err)
	}

	return ctx.OK(res)
}

// GetUserApps returns a paginated list of apps for a particular user. Used by system admin users.
func (c *AppsController) GetUserApps(ctx *app.GetUserAppsAppsContext) error {
	res, err := c.Repository.GetUserApps(ctx.UserID)

	if err != nil {
		if errors.IsErrNotFound(err) {
			return ctx.NotFound(err)
		}
		
		return ctx.InternalServerError(err)
	}

	return ctx.OK(res)
}

// RegisterApp registers new app.
func (c *AppsController) RegisterApp(ctx *app.RegisterAppAppsContext) error {
	// var authObj *auth.Auth
	// hasAuth := auth.HasAuth(ctx)

	// if hasAuth {
	// 	authObj = auth.GetAuth(ctx.Context)
	// } else {
	// 	return ctx.InternalServerError(goa.ErrInternal("Auth has not been set"))
	// }

	// userID := authObj.UserID
	userID := "test user"

	res, err := c.Repository.RegisterApp(ctx.Payload, userID)
	if err != nil {
		if errors.IsErrAlreadyExists(err) {
			return ctx.BadRequest(err)
		}
		
		return ctx.InternalServerError(err)
	}

	return ctx.Created(res)
}

// DeleteApp deletes an app by its id.
func (c *AppsController) DeleteApp(ctx *app.DeleteAppAppsContext) error {
	err := c.Repository.DeleteApp(ctx.AppID)
	if err != nil {
		if errors.IsErrNotFound(err) {
			return ctx.NotFound(err)
		}
		
		if errors.IsErrAlreadyExists(err) {
			return ctx.BadRequest(err)
		}
		
		return ctx.InternalServerError(err)
	}

	return ctx.OK([]byte("Application deleted successfully "))
}

// UpdateApp updates an app by its id.
func (c *AppsController) UpdateApp(ctx *app.UpdateAppAppsContext) error {
	res, err := c.Repository.UpdateApp(ctx.Payload, ctx.AppID)

	if err != nil {
		if errors.IsErrNotFound(err) {
			return ctx.NotFound(err)
		}
		
		if errors.IsErrAlreadyExists(err) {
			return ctx.BadRequest(err)
		}
		
		return ctx.InternalServerError(err)
	}

	return ctx.OK(res)
}

// RegenerateClientSecret regenerates the client secret for an app.
func (c *AppsController) RegenerateClientSecret(ctx *app.RegenerateClientSecretAppsContext) error {
	res, err := c.Repository.RegenerateSecret(ctx.AppID)

	if err != nil {
		if errors.IsErrNotFound(err) {
			return ctx.NotFound(err)
		}
		
		if errors.IsErrAlreadyExists(err) {
			return ctx.BadRequest(err)
		}
		
		return ctx.InternalServerError(err)
	}

	return ctx.OK(res)
}

// VerifyApp check if an app with the supplied credentials exists.
func (c *AppsController) VerifyApp(ctx *app.VerifyAppAppsContext) error {
	clientApp, err := c.Repository.FindApp(ctx.Payload.ID, ctx.Payload.Secret)
	if err != nil {
		return ctx.InternalServerError(err)
	}
	if clientApp == nil {
		return ctx.NotFound(fmt.Errorf("not-found"))
	}

	return ctx.OK(&app.Apps{
		ID:           clientApp.ID,
		Description:  clientApp.Description,
		Domain:       clientApp.Domain,
		Name:         clientApp.Name,
		Owner:        clientApp.Owner,
		RegisteredAt: int(clientApp.RegisteredAt),
	})
}
