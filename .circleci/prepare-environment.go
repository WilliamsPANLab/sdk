package main

import (
	"encoding/json"
	. "log"
	"time"

	"gopkg.in/mgo.v2"

	"flywheel.io/sdk/api"
)

func main() {
	var session *mgo.Session
	var err error

	for i := 1; i <= 3; i++ {
		err = nil
		Println("Connecting to mongo...")
		wait := time.Duration(float64(i) * 5.0 * float64(time.Second))
		session, err = mgo.DialWithTimeout("localhost", wait)
		if err == nil { break }
	}

	if err != nil { Fatalln("Could not connect to mongo:", err) }
	defer session.Close()
	session.SetSafe(&mgo.Safe{})

	root := true
	now := time.Now()
	testUser := &api.User{
		Id: "a@b.c", Email: "a@b.c",
		Firstname: "Test", Lastname: "User",
		Created: &now, Modified: &now,
		RootAccess: &root,
	}

	Println("Inserting test user...")

	// The mongo client does not seem to respect json annotations!
	// Sidestep this by passing it through a json encode, back to a map.
	// Could open an issue on mgo later.
	raw, _ := json.Marshal(testUser)
	var encoded map[string]interface{}
	json.Unmarshal(raw, &encoded)

	// Avatar map is added by the API endpoint, but we don't have that, do we...
	encoded["avatars"] = map[string]interface{}{}

	err = session.DB("scitran").C("users").Insert(encoded)
	if err != nil { Fatalln("Inserting user failed:", err) }
	Println("Test user inserted.")


	Println("Inserting test api key...")

	api_key := map[string]interface{}{
		"_id": "insecure-key",
		"created": &now,
		"last_used": &now,
		"uid": "a@b.c",
		"type": "user",
	}

	err = session.DB("scitran").C("apikeys").Insert(api_key)
	if err != nil { Fatalln("Inserting api key failed:", err) }
	Println("Test api key inserted.")

	var client *api.Client
	var user *api.User

	for i := 1; i <= 15; i++ {
		err = nil
		Println("Connecting to API...")
		client = api.NewApiKeyClient("localhost:8080:insecure-key", api.InsecureNoSSLVerification, api.InsecureUsePlaintext)
		user, _, err = client.GetCurrentUser()
		if err == nil { break }
		time.Sleep(1000 * time.Millisecond)
	}
	if err != nil {	Fatalln("Could not connect to API:", err) }

	Println("Environment is ready with user", user.Firstname, user.Lastname + ".")
}
