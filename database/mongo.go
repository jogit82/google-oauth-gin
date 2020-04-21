package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jogit82/google-oauth-gin/structs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var cred DBCredentials
var info *mgo.DialInfo

// DBCredentials which stores mongoDB Credentials
type DBCredentials struct {
	Cconnectionstring string `json:"cconnectionstring"`
	Cdbname           string `json:"cdbname"`
	Cdbusername       string `json:"cdbusername"`
	Cdbpassword       string `json:"cdbpassword"`
}

// MongoDBConnection Encapsulates a connection to a database.
type MongoDBConnection struct {
	session *mgo.Session
}

func init() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(file, &cred); err != nil {
		log.Println("unable to marshal DB data")
		return
	}

	// tlsConfig := &tls.Config{}

	// var mongoURI = cred.Cconnectionstring
	// info, err := mgo.ParseURL(mongoURI)
	// if err != nil {
	// 	panic(err)
	// }
	// info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
	// 	conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
	// 	return conn, err
	// }
	// info.Timeout = time.Duration(30)
	// info.Database = cred.Cdbname
	// info.Username = cred.Cdbusername
	// info.Password = cred.Cdbpassword

	info = &mgo.DialInfo{
		// Address if its a local db then the value host=localhost
		Addrs: []string{cred.Cconnectionstring},
		// Timeout when a failure to connect to db
		Timeout: time.Duration(30),
		// Database name
		Database: cred.Cdbname,
		// Database credentials if your db is protected
		Username: cred.Cdbusername,
		Password: cred.Cdbpassword,
	}
}

// SaveUser register a user so we know that we saw that user already.
func (mdb MongoDBConnection) SaveUser(u *structs.User) error {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()
	if _, err := mdb.LoadUser(u.Email); err == nil {
		return fmt.Errorf("user already exists")
	}
	c := mdb.session.DB("webadventure").C("users")
	err := c.Insert(u)
	return err
}

// LoadUser get data from a user.
func (mdb MongoDBConnection) LoadUser(Email string) (result structs.User, err error) {
	mdb.session = mdb.GetSession()
	defer mdb.session.Close()
	c := mdb.session.DB("webadventure").C("users")
	err = c.Find(bson.M{"email": Email}).One(&result)
	return result, err
}

// GetSession return a new session if there is no previous one.
func (mdb *MongoDBConnection) GetSession() *mgo.Session {
	if mdb.session != nil {
		return mdb.session.Copy()
	}
	// log.Println("dialInfo::::", dialInfo)
	// session, err := mgo.DialWithInfo(dialInfo)
	session, err := mgo.Dial("localhost")
	// session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Println("error---", err)
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}
