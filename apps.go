package main

import (
	"github.com/JormungandrK/microservice-apps-management/app"
	"github.com/JormungandrK/microservice-apps-management/db"
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
