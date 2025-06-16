package auth

import (
	"github.com/cebas/go-util/util"
	"github.com/skratchdot/open-golang/open"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	config *oauth2.Config = nil
	token  *oauth2.Token  = nil
	ctx                   = context.TODO()
)

const (
	httpPort     = "4242"
	callbackPath = "/oauth/callback"
)

type Gauth struct {
	credentialsFile string
	tokenFile       string
	scope           string
	server          *http.Server
}

func NewGauth(credentialsFile string, scope string) Gauth {
	var tokenFile string

	// if credentialsFile ends with .json, remove it
	if filepath.Ext(credentialsFile) == ".json" {
		tokenFile = credentialsFile[:len(credentialsFile)-len(".json")] + "-token.json"
	} else {
		tokenFile = credentialsFile + "-token.json"
	}

	return Gauth{
		credentialsFile: credentialsFile,
		tokenFile:       tokenFile,
		scope:           scope,
		server:          nil,
	}
}

func (ga *Gauth) HttpClient() (client *http.Client, err error) {
	var credentials []byte
	credentials, err = os.ReadFile(ga.credentialsFile)
	if err != nil {
		return
	}

	config, err = google.ConfigFromJSON(credentials, ga.scope)
	if err != nil {
		return
	}

	config.RedirectURL = "http://localhost:" + httpPort + callbackPath

	// The file tokenFile stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	err = ga.tokenFromFile()
	if err != nil {
		log.Println("google auth token file not found")
		err = ga.tokenFromWeb()
		if err == nil && token != nil {
			err = ga.saveToken()
			if err != nil {
				return
			}
		} else {
			err = errors.New("error getting google auth token")
			return
		}
	}

	client = config.Client(ctx, token)
	return
}

// Saves a token to a file
func (ga *Gauth) saveToken() (err error) {
	log.Printf("saving google auth token to: [%s]\n", ga.tokenFile)

	file, err := os.OpenFile(ga.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if !util.WarningErrorCheck(err) {
		return
	}

	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	err = json.NewEncoder(file).Encode(token)
	return
}

// Retrieves a token from a local file.
func (ga *Gauth) tokenFromFile() (err error) {
	log.Printf("reading google auth token from: [%s]\n", ga.tokenFile)

	file, err := os.Open(ga.tokenFile)
	if err != nil {
		return
	}

	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	token = &oauth2.Token{}
	return json.NewDecoder(file).Decode(token)
}

func (ga *Gauth) shutdownServer() {
	log.Println("http server shutting down")
	err := ga.server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func (ga *Gauth) waitForResponse() {
	for {
		time.Sleep(1 * time.Second)
		if token != nil {
			time.Sleep(1 * time.Second)
			ga.shutdownServer()
			break
		}
	}
}

func (ga *Gauth) tokenFromWeb() (err error) {
	log.Println("getting new google auth token from web")

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	err = open.Run(authURL)
	util.FatalErrorCheck(err)

	http.HandleFunc(callbackPath, callbackHandler)

	ga.server = &http.Server{Addr: ":" + httpPort, Handler: nil}

	token = nil
	go ga.waitForResponse()

	err = ga.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return
	}

	return nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	log.Println("request received")

	queryParts, err := url.ParseQuery(r.URL.RawQuery)
	util.FatalErrorCheck(err)

	// Use the authorization code that is pushed to the redirect URL.
	code := queryParts["code"][0]

	log.Println("got token")

	// Exchange will do the handshake to retrieve the initial access token.
	token, err = config.Exchange(ctx, code)
	util.FatalErrorCheck(err)

	// show success page
	msg := "<p><strong>Autenticado!</strong></p>"
	_, _ = w.Write([]byte(msg))
}
