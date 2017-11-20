package main

import (
	"fmt"

	"github.com/JormungandrK/microservice-apps-management/app"
	"github.com/JormungandrK/microservice-apps-management/db"
	"github.com/JormungandrK/microservice-security/auth"
	"github.com/goadesign/goa"
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

// Return an app by its ID.
func (c *AppsController) Get(ctx *app.GetAppsContext) error {
	res, err := c.Repository.GetApp(ctx.AppID)

	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK(res)
}

// Return a paginated list of all apps for the current user.
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
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK(res)
}

// Return a paginated list of apps for a particular user. Used by system admin users.
func (c *AppsController) GetUserApps(ctx *app.GetUserAppsAppsContext) error {
	res, err := c.Repository.GetUserApps(ctx.UserID)

	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK(res)
}

// Register new app.
func (c *AppsController) RegisterApp(ctx *app.RegisterAppAppsContext) error {
	var authObj *auth.Auth
	hasAuth := auth.HasAuth(ctx)

	if hasAuth {
		authObj = auth.GetAuth(ctx.Context)
	} else {
		return ctx.InternalServerError(goa.ErrInternal("Auth has not been set"))
	}

	userID := authObj.UserID

	res, err := c.Repository.RegisterApp(ctx.Payload, userID)

	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.Created(res)
}

// Delete an app by its id.
func (c *AppsController) DeleteApp(ctx *app.DeleteAppAppsContext) error {
	err := c.Repository.DeleteApp(ctx.AppID)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK([]byte("Application deleted successfully "))
}

// Update an app by its id.
func (c *AppsController) UpdateApp(ctx *app.UpdateAppAppsContext) error {
	res, err := c.Repository.UpdateApp(ctx.Payload, ctx.AppID)

	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK(res)
}

// Regenerate client secret for an app.
func (c *AppsController) RegenerateClientSecret(ctx *app.RegenerateClientSecretAppsContext) error {
	res, err := c.Repository.RegenerateSecret(ctx.AppID)

	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK(res)
}

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
