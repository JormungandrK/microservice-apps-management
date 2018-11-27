// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "apps-management": Application Media Types
//
// Command:
// $ goagen
// --design=github.com/Microkubes/microservice-apps-management/design
// --out=$(GOPATH)/src/github.com/Microkubes/microservice-apps-management
// --version=v1.3.1

package app

import (
	"github.com/goadesign/goa"
	"unicode/utf8"
)

// apps media type (default view)
//
// Identifier: application/vnd.goa.apps+json; view=default
type Apps struct {
	// Description of the app
	Description string `form:"description" json:"description" yaml:"description" xml:"description"`
	// App domain
	Domain string `form:"domain" json:"domain" yaml:"domain" xml:"domain"`
	// Unique app ID
	ID string `form:"id" json:"id" yaml:"id" xml:"id"`
	// Name of the app
	Name string `form:"name" json:"name" yaml:"name" xml:"name"`
	// User ID
	Owner string `form:"owner" json:"owner" yaml:"owner" xml:"owner"`
	// Time when app is registered
	RegisteredAt int `form:"registeredAt" json:"registeredAt" yaml:"registeredAt" xml:"registeredAt"`
}

// Validate validates the Apps media type instance.
func (mt *Apps) Validate() (err error) {
	if mt.ID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "id"))
	}
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if mt.Domain == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "domain"))
	}
	if mt.Owner == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "owner"))
	}

	if utf8.RuneCountInString(mt.Description) > 300 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.description`, mt.Description, utf8.RuneCountInString(mt.Description), 300, false))
	}
	if utf8.RuneCountInString(mt.Name) > 50 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, utf8.RuneCountInString(mt.Name), 50, false))
	}
	return
}

// reg-apps media type (default view)
//
// Identifier: application/vnd.goa.reg.apps+json; view=default
type RegApps struct {
	// App ID
	ID string `form:"id" json:"id" yaml:"id" xml:"id"`
	// Client secret
	Secret string `form:"secret" json:"secret" yaml:"secret" xml:"secret"`
}

// Validate validates the RegApps media type instance.
func (mt *RegApps) Validate() (err error) {
	if mt.ID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "id"))
	}
	if mt.Secret == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "secret"))
	}
	return
}
