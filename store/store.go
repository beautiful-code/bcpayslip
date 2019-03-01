package store

import (
	"os"

	"bcpayslip/helpers"
	"bcpayslip/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetSession Dial to database and Return mgo session ...
func GetSession(collection string, pk string) *mgo.Session {
	var session *mgo.Session
	var err error
	if os.Getenv("bc_env") == "development" {
		session, err = mgo.Dial("127.0.0.1")
	} else {
		session, err = mgo.Dial(os.Getenv("MONGO_URI"))
	}
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	EnsureIndex(collection, pk, session)
	return session
}

// EnsureIndex Ensure an index on the collection, why? ...
func EnsureIndex(collection string, pk string, s *mgo.Session) {
	session := s.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C(collection)
	index := mgo.Index{
		Key:        []string{pk},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// GetUser get user data ...
func GetUser(userID string) (models.User, error) {
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("User")
	var user models.User
	err := c.Find(bson.M{"userid": userID}).One(&user)
	return user, err
}

// SaveUser Create user data ...
func SaveUser(userID string, firstName string, lastName string, email string, accessToken string, avatar string) error {
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("User")
	_, err := GetUser(userID)
	if err == nil {
		err = c.Update(
			bson.M{"userid": userID},
			bson.M{"$set": bson.M{
				"userid": userID, "firstname": firstName,
				"lastname": lastName, "email": email,
				"accesstoken": accessToken, "avatar": helpers.ImageToBase64(avatar),
			}},
		)
	} else {
		var user models.User
		user.UserID = userID
		user.FirstName = firstName
		user.LastName = lastName
		user.Email = email
		user.AccessToken = accessToken
		user.Avatar = helpers.ImageToBase64(avatar)
		err = c.Insert(user)
	}
	return err
}
