package design

// Use . imports to enable the DSL
import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
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
	DefaultMedia(AppMedia)

	Action("get", func() {
		Description("Get app by id")
		Routing(GET("/:appId"))
		Params(func() {
			Param("appId", String, "App ID")
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
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

// AppsPayload defines the payload for the client apps.
var AppPayload = Type("AppPayload", func() {
	Description("Payload for the client apps")

	Attribute("name", String, "Name of the app", func() {
		MaxLength(50)
	})
	Attribute("description", String, "Description of the app", func() {
		MaxLength(300)
	})
	Attribute("domain", String, "App domain", func() {
		Format("uri")
	})

	Required("name", "description", "domain")
})

// Swagger UI
var _ = Resource("swagger", func() {
	Description("The API swagger specification")

	Files("swagger.json", "swagger/swagger.json")
	Files("swagger-ui/*filepath", "swagger-ui/dist")
})
