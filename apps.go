package main

import (
	"github.com/JormungandrK/microservice-apps-management/app"
	"github.com/goadesign/goa"
)

// AppsController implements the apps resource.
type AppsController struct {
	*goa.Controller
}

// NewAppsController creates a apps controller.
func NewAppsController(service *goa.Service) *AppsController {
	return &AppsController{Controller: service.NewController("AppsController")}
}

// Get runs the get action.
func (c *AppsController) Get(ctx *app.GetAppsContext) error {
	// AppsController_Get: start_implement

	// Put your logic here

	// AppsController_Get: end_implement
	res := &app.Apps{}
	return ctx.OK(res)
}
