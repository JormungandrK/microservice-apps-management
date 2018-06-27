package db

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/Microkubes/microservice-apps-management/app"
	"github.com/asaskevich/govalidator"
	"github.com/goadesign/goa"
	"github.com/JormungandrK/backends"
)

// AppRepository defaines the interface for accessing the application data.
type AppRepository interface {
	// GetApp looks up a applications by the app ID.
	GetApp(appID string) (*app.Apps, error)
	GetMyApps(userID string) ([]byte, error)
	GetUserApps(userID string) ([]byte, error)
	RegisterApp(payload *app.AppPayload, userID string) (*app.RegApps, error)
	DeleteApp(appID string) error
	UpdateApp(payload *app.AppPayload, appID string) (*app.Apps, error)
	RegenerateSecret(appID string) ([]byte, error)
	FindApp(id, secret string) (*ClientApp, error)
}

// ClientApp holds the data for a registered application (client).
type ClientApp struct {
	ID           string `json:"id" bson:"_id"`
	Name         string `json:"name" bson:"name"`
	Description  string `json:"description,omitempty" bson:"description"`
	Domain       string `json:"domain,omitempty" bson:"domain"`
	Owner        string `json:"owner" bson:"owner"`
	RegisteredAt int64  `json:"registeredAt" bson:"registeredAt"`
	Secret       string `json:"secret" bson:"secret"`
}

func (r *BackendAppsService) GetApp(appID string) (*app.Apps, error) {
	return nil, nil
}

func (r *BackendAppsService) GetMyApps(userID string) ([]byte, error) {
	return nil, nil
}

func (r *BackendAppsService) GetUserApps(userID string) ([]byte, error) {
	return nil, nil
}

func (r *BackendAppsService) RegisterApp(payload *app.AppPayload, userID string) (*app.RegApps, error) {
	return nil, nil
}

func (r *BackendAppsService) DeleteApp(appID string) error {
	return nil
}

func (r *BackendAppsService) UpdateApp(payload *app.AppPayload, appID string) (*app.Apps, error) {
	return nil, nil
}

func (r *BackendAppsService) RegenerateSecret(appID string) ([]byte, error) {
	return nil, nil
}

func (r *BackendAppsService) FindApp(id, secret string) (*ClientApp, error) {
	return nil, nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded securely generated random string.
// It will return an error if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Validate domain
func validateDomain(domain string) error {
	if ok := govalidator.IsURL(domain); !ok {
		return goa.ErrBadRequest("invalid domain name")
	}

	return nil
}


// BackendAppsService holds data for implementation of the AppRepository interface.
type BackendAppsService struct {
	appsRepository backends.Repository
}

// NewAppssService creates new AppsService.
func NewAppsService(appsRepository backends.Repository) AppRepository {
	return &BackendAppsService{
		appsRepository: appsRepository,
	}
}