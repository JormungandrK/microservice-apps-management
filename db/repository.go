package db

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/Microkubes/microservice-apps-management/app"
	"github.com/asaskevich/govalidator"
	"github.com/goadesign/goa"
	"github.com/JormungandrK/backends"
	"time"
	"fmt"
)

// AppRepository defaines the interface for accessing the application data.
type AppRepository interface {
	// GetApp looks up a applications by the app ID.
	GetApp(appID string) (*app.Apps, error)
	GetMyApps(userID, order, sorting string, limit, offset int) ([]byte, error)
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
	apps, err := r.appsRepository.GetOne(backends.NewFilter().Match("id", appID), &ClientApp{})
	if err != nil {
		return nil, err
	}

	regCli := &app.Apps{}
	if err = backends.MapToInterface(apps, regCli); err != nil {
		return nil, err
	}

	return regCli, nil
}

func (r *BackendAppsService) GetMyApps(userID, order, sorting string, limit, offset int) ([]byte, error) {
	var resultsTypeInterface []byte

	// userID, order, sorting string, limit, offset int
	apps, err := r.appsRepository.GetAll(backends.NewFilter().Match("id", userID), resultsTypeInterface, order, sorting, limit, offset)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *BackendAppsService) GetUserApps(userID string) ([]byte, error) { 
	return nil, nil
}

func (r *BackendAppsService) FindApp(id, secret string) (*ClientApp, error) { 
	apps, err := r.appsRepository.GetOne(backends.NewFilter().Match("id", id).Match("secret", secret), &ClientApp{})
	if err != nil {
		return nil, err
	}
	
	regCli := &ClientApp{}
	if err = backends.MapToInterface(apps, regCli); err != nil {
		return nil, err
	}

	return regCli, nil
}

func (r *BackendAppsService) RegisterApp(payload *app.AppPayload, userID string) (*app.RegApps, error) {
	secret, err := GenerateRandomString(42)
	if err != nil {
		return nil, err
	}

	registerClientApp := &ClientApp{
		Name:	payload.Name,
		Description: *payload.Description,
		Domain:	*payload.Domain,
		Secret: secret,
		RegisteredAt: int64(time.Now().Unix()),
		Owner:	userID,
	}

	result, err := r.appsRepository.Save(registerClientApp, nil)
	if err != nil {
		return nil, err
	}

	regCli := &ClientApp{}
	if err = backends.MapToInterface(result, regCli); err != nil {
		return nil, err
	}

	regApp := &app.RegApps {
		ID: regCli.ID,
		Secret: regCli.Secret,
	}

	return regApp, err
}

func (r *BackendAppsService) DeleteApp(appID string) error {
	err := r.appsRepository.DeleteOne(backends.NewFilter().Match("id", appID))
	if err != nil {
		return err
	}

	return err
}

func (r *BackendAppsService) UpdateApp(payload *app.AppPayload, appID string) (*app.Apps, error) {
	result, err := r.appsRepository.Save(payload, backends.NewFilter().Match("id", appID))
	if err != nil {
		return nil, err
	}

	regCli := &app.Apps{}
	if err = backends.MapToInterface(result, regCli); err != nil {
		return nil, err
	}

	return regCli, err
}

func (r *BackendAppsService) RegenerateSecret(appID string) ([]byte, error) {
	app, err := r.appsRepository.GetOne(backends.NewFilter().Match("id", appID), &ClientApp{})
	if err != nil {
		return nil, err
	}

	regCli := &ClientApp{}
	if err = backends.MapToInterface(app, regCli); err != nil {
		return nil, err
	}

	regCli.Secret, err = GenerateRandomString(42)
	if err != nil {
		return nil, err
	}

	_, err = r.appsRepository.Save(regCli, backends.NewFilter().Match("id", appID))
	if err != nil {
		return nil, err
	}

	return []byte(regCli.Secret), err
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