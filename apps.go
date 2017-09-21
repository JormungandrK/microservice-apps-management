package main

import (
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

// Get runs the get action.
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

func (c *AppsController) GetMyApps(ctx *app.GetMyAppsAppsContext) error {
	var authObj *auth.Auth
	hasAuth := auth.HasAuth(ctx)

	if hasAuth {
		authObj = auth.GetAuth(ctx.Context)
	} else {
		return ctx.InternalServerError(goa.ErrInternal("Auth has not been set"))
	}

	userID := authObj.UserID

	res, err := c.Repository.GetUserApps(userID)

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

func (c *AppsController) GetUserApps(ctx *app.GetUserAppsAppsContext) error {
	res, err := c.Repository.GetUserApps(ctx.UserID)

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

func (c *AppsController) RegisterApp(ctx *app.RegisterAppAppsContext) error {
	// var authObj *auth.Auth
	// hasAuth := auth.HasAuth(ctx)

	// if hasAuth {
	// 	authObj = auth.GetAuth(ctx.Context)
	// } else {
	// 	return ctx.InternalServerError(goa.ErrInternal("Auth has not been set"))
	// }

	// userID := authObj.UserID

	userID := "5975c461f9f8eb02aae053f3"

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
