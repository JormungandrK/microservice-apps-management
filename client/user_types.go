// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "apps-management": Application User Types
//
// Command:
// $ goagen
// --design=github.com/JormungandrK/microservice-apps-management/design
// --out=$(GOPATH)/src/github.com/JormungandrK/microservice-apps-management
// --version=v1.2.0-dirty

package client

import (
	"github.com/goadesign/goa"
	"unicode/utf8"
)

// Payload for the client apps
type appPayload struct {
	// Description of the app
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	// App domain
	Domain *string `form:"domain,omitempty" json:"domain,omitempty" xml:"domain,omitempty"`
	// Name of the app
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
}

// Validate validates the appPayload type instance.
func (ut *appPayload) Validate() (err error) {
	if ut.Name == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "name"))
	}
	if ut.Description != nil {
		if utf8.RuneCountInString(*ut.Description) > 300 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`request.description`, *ut.Description, utf8.RuneCountInString(*ut.Description), 300, false))
		}
	}
	if ut.Name != nil {
		if utf8.RuneCountInString(*ut.Name) > 50 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`request.name`, *ut.Name, utf8.RuneCountInString(*ut.Name), 50, false))
		}
	}
	return
}

// Publicize creates AppPayload from appPayload
func (ut *appPayload) Publicize() *AppPayload {
	var pub AppPayload
	if ut.Description != nil {
		pub.Description = ut.Description
	}
	if ut.Domain != nil {
		pub.Domain = ut.Domain
	}
	if ut.Name != nil {
		pub.Name = *ut.Name
	}
	return &pub
}

// Payload for the client apps
type AppPayload struct {
	// Description of the app
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	// App domain
	Domain *string `form:"domain,omitempty" json:"domain,omitempty" xml:"domain,omitempty"`
	// Name of the app
	Name string `form:"name" json:"name" xml:"name"`
}

// Validate validates the AppPayload type instance.
func (ut *AppPayload) Validate() (err error) {
	if ut.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "name"))
	}
	if ut.Description != nil {
		if utf8.RuneCountInString(*ut.Description) > 300 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`type.description`, *ut.Description, utf8.RuneCountInString(*ut.Description), 300, false))
		}
	}
	if utf8.RuneCountInString(ut.Name) > 50 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`type.name`, ut.Name, utf8.RuneCountInString(ut.Name), 50, false))
	}
	return
}