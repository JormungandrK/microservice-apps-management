package db

import (
	"crypto/rand"
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/JormungandrK/microservice-apps-management/app"
	"github.com/goadesign/goa"

	"golang.org/x/crypto/bcrypt"
)

// AppRepository defaines the interface for accessing the application data.
type AppRepository interface {
	// GetApp looks up a applications by the app ID.
	GetApp(appID string) (*app.Apps, error)
	GetUserApps(userID string) ([]byte, error)
	RegisterApp(payload *app.AppPayload, userID string) (*app.RegApps, error)
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
			return nil, goa.ErrNotFound(err)
		} else {
			return nil, goa.ErrInternal(err)
		}
	}

	return res, nil
}

func (c *MongoCollection) GetUserApps(userID string) ([]byte, error) {
	objectID, err := hexToObjectID(userID)
	if err != nil {
		return nil, err
	}

	var apps []map[string]interface{}
	if err := c.Find(bson.M{"owner": objectID}).Sort("registeredat").All(&apps); err != nil {
		return nil, goa.ErrInternal(err)
	}

	if len(apps) == 0 {
		return nil, goa.ErrNotFound("No apps found!")
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
	b, _ := generateSecret(42)
	hashedSecret, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return nil, goa.ErrInternal(err)
	}
	secret := string(hashedSecret)

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

func generateSecret(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
