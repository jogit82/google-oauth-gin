package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jogit82/google-oauth-gin/database"
	"github.com/jogit82/google-oauth-gin/structs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cred Credentials
var conf *oauth2.Config

func init() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(file, &cred); err != nil {
		log.Println("unable to marshal data")
		return
	}

	conf = &oauth2.Config{
		ClientID:     cred.Cid,     // application's ID.
		ClientSecret: cred.Csecret, // application's secret.
		RedirectURL:  "http://127.0.0.1:9090/auth",
		Scopes: []string{ // specifies optional requested permissions.
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		// Endpoint contains the resource server's token endpoint URLs.
		// These are constants specific to each server and are
		// often available via site-specific packages, such as
		// google.Endpoint or github.Endpoint.
		Endpoint: google.Endpoint,
	}
}

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func getLoginURL(state string) string {
	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return conf.AuthCodeURL(state)
}

// AuthHandler handles authentication of a user and initiates a session.
func AuthHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.

	// shortcut to get session
	session := sessions.Default(c)

	// Get returns the session value associated to the given key.
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.HTML(http.StatusUnauthorized, "error", gin.H{"message": "Invalid session state. retrievedState != queryState"})
		return
	}
	code := c.Request.URL.Query().Get("code")
	// Exchange converts an authorization code into a token.
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error", gin.H{"message": "Login failed. Please try again."})
		return
	}

	// Client returns an HTTP client using the provided token. The token will auto-refresh as necessary. The underlying HTTP transport will be obtained using the provided context. The returned client and its Transport should not be modified.
	client := conf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	log.Println("Resp body: ", string(data))
	u := structs.User{}
	if err = json.Unmarshal(data, &u); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error", gin.H{"message": "Error marshalling response. Please try agian."})
		return
	}

	// Set sets the session value associated to the given key.
	session.Set("user-id", u.Email)
	session.Set("user-img", u.Picture)

	// Save saves all sessions used during the current request.
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error", gin.H{"message": "Error while saving session. Please try again."})
		return
	}
	seen := false
	db := database.MongoDBConnection{}
	if _, mongoErr := db.LoadUser(u.Email); mongoErr == nil {
		seen = true
	} else {
		err = db.SaveUser(&u)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error", gin.H{"message": "Error while saving user. Please try again."})
			return
		}
	}
	// c.HTML(http.StatusOK, "battle.tmpl", gin.H{"email": u.Email, "seen": seen})
	c.HTML(http.StatusOK, "user", gin.H{"user": u.Email, "userImg": u.Picture, "seen": seen})
}

// LoginHandler handles the login procedure.
func LoginHandler(c *gin.Context) {
	log.Println("client IP>>>>>", c.ClientIP())
	state, err := RandToken(32)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"message": "Error while generating random data."})
		return
	}
	session := sessions.Default(c)
	session.Set("state", state)
	err = session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"message": "Error while saving session."})
		return
	}
	link := getLoginURL(state)
	c.HTML(http.StatusOK, "auth", gin.H{"link": link})
}
