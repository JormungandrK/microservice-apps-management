package db

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/JormungandrK/microservice-apps-management/app"
	"github.com/asaskevich/govalidator"
	"github.com/goadesign/goa"
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
}

// MongoCollection wraps a mgo.Collection to embed methods in models.
type MongoCollection struct {
	*mgo.Collection
}

// NewSession returns a new Mongo Session.
func NewSession(Host string, Username string, Password string, Database string) *mgo.Session {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{Host},
		Username: Username,
		Password: Password,
		Database: Database,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	// SetMode - consistency mode for the session.
	session.SetMode(mgo.Monotonic, true)

	return session
}

// PrepareDB ensure presence of persistent and immutable data in the DB.
func PrepareDB(session *mgo.Session, db string, dbCollection string, indexes [][]string) *mgo.Collection {
	// Create collection
	collection := session.DB(db).C(dbCollection)

	// Define indexes
	for _, elem := range indexes {
		i := elem
		index := mgo.Index{
			Key:        i,
			Unique:     true,
			DropDups:   false,
			Background: true,
			Sparse:     true,
		}

		// Create indexes
		if err := collection.EnsureIndex(index); err != nil {
			panic(err)
		}
	}

	return collection
}

func (c *MongoCollection) GetApp(appID string) (*app.Apps, error) {
	objectID, err := hexToObjectID(appID)
	if err != nil {
		return nil, err
	}

	res := &app.Apps{}
	if err := c.FindId(objectID).One(&res); err != nil {
		if err.Error() == "not found" {
			return nil, goa.ErrNotFound("app not found.")
		} else {
			return nil, goa.ErrInternal(err)
		}
	}

	res.ID = appID

	return res, nil
}

func (c *MongoCollection) GetMyApps(userID string) ([]byte, error) {
	var apps []map[string]interface{}
	if err := c.Find(bson.M{"owner": userID}).Sort("registeredat").All(&apps); err != nil {
		return nil, goa.ErrInternal(err)
	}

	if len(apps) == 0 {
		return nil, goa.ErrNotFound("no apps found!")
	}

	for _, client := range apps {
		client["id"] = client["_id"]
		delete(client, "_id")
		delete(client, "secret")
	}

	res, err := json.Marshal(apps)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return res, nil
}

func (c *MongoCollection) GetUserApps(userID string) ([]byte, error) {
	var apps []map[string]interface{}
	if err := c.Find(bson.M{"owner": userID}).Sort("registeredat").All(&apps); err != nil {
		return nil, goa.ErrInternal(err)
	}

	if len(apps) == 0 {
		return nil, goa.ErrNotFound("no apps found!")
	}

	for _, client := range apps {
		client["id"] = client["_id"]
		delete(client, "_id")
	}

	res, err := json.Marshal(apps)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return res, nil
}

func (c *MongoCollection) RegisterApp(payload *app.AppPayload, userID string) (*app.RegApps, error) {
	appID := bson.NewObjectIdWithTime(time.Now())
	registeredAt := int(time.Now().Unix())
	secret, err := GenerateRandomString(42)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}
	if err := validateDomain(*payload.Domain); err != nil {
		return nil, err
	}

	err = c.Insert(bson.M{
		"_id":          appID,
		"name":         payload.Name,
		"description":  payload.Description,
		"domain":       payload.Domain,
		"owner":        userID,
		"secret":       secret,
		"registeredat": registeredAt,
	})

	if err != nil {
		if mgo.IsDup(err) {
			return nil, goa.ErrBadRequest("application already exists in the database")
		}
		return nil, goa.ErrInternal(err)
	}

	res := &app.RegApps{
		ID:     appID.Hex(),
		Secret: secret,
	}

	return res, nil
}

func (c *MongoCollection) DeleteApp(appID string) error {
	objectID, err := hexToObjectID(appID)
	if err != nil {
		return err
	}

	err = c.RemoveId(objectID)
	if err != nil {
		if err.Error() == "not found" {
			return goa.ErrNotFound("no apps found!")
		} else {
			return goa.ErrInternal(err)
		}
	}

	return nil
}

func (c *MongoCollection) UpdateApp(payload *app.AppPayload, appID string) (*app.Apps, error) {
	objectID, err := hexToObjectID(appID)
	if err != nil {
		return nil, err
	}

	err = c.Update(
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"name":        payload.Name,
			"description": payload.Description,
			"domain":      payload.Domain,
		},
		})

	if err != nil {
		if err.Error() == "not found" {
			return nil, goa.ErrNotFound("application not found.")
		} else {
			return nil, goa.ErrInternal(err)
		}
	}

	client, err := c.GetApp(appID)
	if err != nil {
		return nil, err
	}
	client.ID = appID

	return client, nil
}

func (c *MongoCollection) RegenerateSecret(appID string) ([]byte, error) {
	objectID, err := hexToObjectID(appID)
	if err != nil {
		return nil, err
	}

	secret, err := GenerateRandomString(42)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	err = c.Update(
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"secret": secret,
		},
		})

	if err != nil {
		if err.Error() == "not found" {
			return nil, goa.ErrNotFound("application not found")
		} else {
			return nil, goa.ErrInternal(err)
		}
	}

	var client map[string]interface{}
	if err := c.FindId(objectID).One(&client); err != nil {
		return nil, goa.ErrInternal(err)
	}

	client["id"] = client["_id"]
	client["secret"] = secret
	delete(client, "_id")

	res, err := json.Marshal(client)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}

	return res, nil
}

// Convert hex representation of object id to bson object id
func hexToObjectID(hexID string) (bson.ObjectId, error) {
	// Return whether userID is a valid hex representation of object id.
	if bson.IsObjectIdHex(hexID) != true {
		return "", goa.ErrBadRequest("invalid ID")
	}

	// Return an ObjectId from the provided hex representation.
	objectID := bson.ObjectIdHex(hexID)

	// Return true if objectID is valid. A valid objectID must contain exactly 12 bytes.
	if objectID.Valid() != true {
		return "", goa.ErrInternal("invalid object ID")
	}

	return objectID, nil
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
