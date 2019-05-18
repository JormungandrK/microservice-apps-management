package design

// Use . imports to enable the DSL
import (
	. "github.com/keitaroinc/goa/design"
	. "github.com/keitaroinc/goa/design/apidsl"
)

// Define default description and default global property values
var _ = API("apps-management", func() {
	Title("The apps management microservice")
	Description("A service that provides basic access to the applications management")
	Version("1.0")
	Scheme("http")
	Host("localhost:8080")
})

// Resources group related API endpoints together.
var _ = Resource("apps", func() {
	BasePath("/apps")

	// Alows preflight request - HTTP OPTIONS
	Origin("*", func() {
		Methods("OPTIONS")
	})

	Action("get", func() {
		Description("Get app by id")
		Routing(GET("/:appId"))
		Params(func() {
			Param("appId", String, "App ID")
		})
		Response(OK, AppMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getMyApps", func() {
		Description("Get all user's apps")
		Routing(GET("/my"))
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getUserApps", func() {
		Description("Get app by id")
		Routing(GET("/users/:userId/all"))
		Params(func() {
			Param("userId", String, "User ID")
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("registerApp", func() {
		Description("Register new app")
		Routing(POST(""))
		Payload(AppPayload)
		Response(Created, RegAppMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("deleteApp", func() {
		Description("Delete an app")
		Routing(DELETE("/:appId"))
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("updateApp", func() {
		Description("Register new app")
		Routing(PUT("/:appId"))
		Payload(AppPayload)
		Response(OK, AppMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("regenerateClientSecret", func() {
		Description("Regenerate client secret")
		Routing(PUT("/:appId/regenerate-secret"))
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("verifyApp", func() {
		Description("Verify an application by its ID and secret")
		Routing(POST("/verify"))
		Payload(AppCredentialsPayload)
		Response(OK, AppMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})

// AppMedia defines the media type used to render client apps.
var AppMedia = MediaType("application/vnd.goa.apps+json", func() {
	TypeName("apps")
	Reference(AppPayload)

	Attributes(func() {
		Attribute("id", String, "Unique app ID")
		Attribute("name")
		Attribute("description")
		Attribute("domain")
		Attribute("owner", String, "User ID")
		Attribute("secret", String, "Client secret")
		Attribute("registeredAt", Integer, "Time when app is registered")
		Required("id", "name", "description", "domain", "owner", "registeredAt")
	})

	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("description")
		Attribute("domain")
		Attribute("owner")
		Attribute("registeredAt")
	})
})

// RegAppMedia defines the media type used to render client apps.
var RegAppMedia = MediaType("application/vnd.goa.reg.apps+json", func() {
	TypeName("reg-apps")

	Attributes(func() {
		Attribute("id", String, "App ID")
		Attribute("secret", String, "Client secret")
		Required("id", "secret")
	})

	View("default", func() {
		Attribute("id")
		Attribute("secret")
	})
})

// AppsPayload defines the payload for the client apps.
var AppPayload = Type("AppPayload", func() {
	Description("Payload for the client apps")

	Attribute("name", String, "Name of the app", func() {
		MaxLength(50)
	})
	Attribute("description", String, "Description of the app", func() {
		MaxLength(300)
	})
	Attribute("domain", String, "App domain")

	Required("name")
})

// AppCredentialsPayload holds the app credentials: app ID and app secret.
var AppCredentialsPayload = Type("AppCredentialsPayload", func() {
	Description("App ID+secret credentials")
	Attribute("id", String, "The app ID")
	Attribute("secret", String, "The app secret")
	Required("id", "secret")
})

// Swagger UI
var _ = Resource("swagger", func() {
	Description("The API swagger specification")

	Files("swagger.json", "swagger/swagger.json")
	Files("swagger-ui/*filepath", "swagger-ui/dist")
})
