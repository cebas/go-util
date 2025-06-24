package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/cebas/go-util/util"
	"github.com/skratchdot/open-golang/open"
)

var (
	config *oauth2.Config = nil
	token  *oauth2.Token  = nil
	log    *util.Log
)

const (
	httpPort = "4242"
)

type Gauth struct {
	ctx             context.Context
	credentialsFile string
	tokenFile       string
	callbackPath    string
	server          *http.Server
}

func NewGauth(ctx context.Context, credentialsFile string, aLog *util.Log) Gauth {
	log = aLog

	return Gauth{
		ctx:             ctx,
		credentialsFile: credentialsFile,
		tokenFile:       "token." + credentialsFile,
		callbackPath:    "/oauth/callback-" + util.RandomString(8),
		server:          nil,
	}
}

func (ga *Gauth) HttpClient(scope string) (client *http.Client, err error) {
	var credentials []byte
	credentials, err = os.ReadFile(ga.credentialsFile)
	if err != nil {
		return
	}

	config, err = google.ConfigFromJSON(credentials, scope)
	if err != nil {
		return
	}

	config.RedirectURL = "http://localhost:" + httpPort + ga.callbackPath

	// The file tokenFile stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	err = ga.tokenFromFile()
	if err != nil {
		log.Println(2, "google auth token file not found")
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

	client = config.Client(ga.ctx, token)
	return
}

// Saves a token to a file
func (ga *Gauth) saveToken() (err error) {
	log.Printf(2, "saving google auth token to: [%s]\n", ga.tokenFile)

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
	log.Printf(2, "reading google auth token from: [%s]\n", ga.tokenFile)

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
	log.Println(2, "http server shutting down")
	err := ga.server.Shutdown(ga.ctx)
	util.FatalErrorCheck(err)
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
	log.Println(2, "getting new google auth token from web")

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	err = open.Run(authURL)
	util.FatalErrorCheck(err)

	http.HandleFunc(ga.callbackPath, ga.callbackHandler)

	ga.server = &http.Server{Addr: ":" + httpPort, Handler: nil}

	token = nil
	go ga.waitForResponse()

	err = ga.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return
	}

	return nil
}

func (ga *Gauth) callbackHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	log.Println(3, "request received")

	queryParts, err := url.ParseQuery(r.URL.RawQuery)
	util.FatalErrorCheck(err)

	// Use the authorization code that is pushed to the redirect URL.
	code := queryParts["code"][0]

	log.Println(3, "got token")

	// Exchange will do the handshake to retrieve the initial access token.
	token, err = config.Exchange(ga.ctx, code)
	util.FatalErrorCheck(err)

	// show success page
	msg := "<p><strong>Autenticado!</strong></p>"
	_, _ = w.Write([]byte(msg))
}
