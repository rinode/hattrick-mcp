// Command hattrick-mcp is an MCP server that exposes Hattrick CHPP data to AI assistants.
package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	_ "time/tzdata" // the hattrick library needs Europe/Stockholm, which Windows and minimal Linux lack

	"github.com/lucianoq/hattrick"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/mrjones/oauth"
)

// Set at build time: go build -ldflags "-X main.consumerKey=... -X main.consumerSecret=...".
var (
	version        = "dev"
	consumerKey    string
	consumerSecret string
)

type credentials struct {
	AccessToken  string `json:"access_token"`
	AccessSecret string `json:"access_secret"`
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr) // stdout carries the MCP protocol

	if v := os.Getenv("HATTRICK_CONSUMER_KEY"); v != "" {
		consumerKey = v
	}
	if v := os.Getenv("HATTRICK_CONSUMER_SECRET"); v != "" {
		consumerSecret = v
	}
	if consumerKey == "" || consumerSecret == "" {
		log.Fatal("missing CHPP consumer key: set HATTRICK_CONSUMER_KEY and HATTRICK_CONSUMER_SECRET, or build with -ldflags")
	}

	path, err := credentialsPath()
	if err != nil {
		log.Fatal(err)
	}

	cmd := "serve"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "serve":
		err = serve(path)
	case "http":
		err = serveHTTP(path)
	case "version":
		fmt.Println(version)
	default:
		err = fmt.Errorf("unknown command %q (use serve, http or version)", cmd)
	}
	if err != nil {
		log.Fatal(err)
	}
}

// serveHTTP runs the MCP server over Streamable HTTP for hosting, behind a
// bearer token. Put TLS in front (a reverse proxy) for anything public.
func serveHTTP(path string) error {
	token := os.Getenv("HATTRICK_MCP_TOKEN")
	if token == "" {
		return errors.New("set HATTRICK_MCP_TOKEN to a secret value to run the HTTP server")
	}
	addr := os.Getenv("HATTRICK_MCP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	a, err := newAccount(path)
	if err != nil {
		return err
	}

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return newServer(a) }, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.Handle("/", bearerAuth(token, handler))

	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errc := make(chan error, 1)
	go func() {
		log.Printf("hattrick-mcp HTTP server listening on %s", addr)
		errc <- srv.ListenAndServe()
	}()
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutCtx)
	}
}

// bearerAuth rejects any request without the exact "Bearer <token>" header.
func bearerAuth(token string, next http.Handler) http.Handler {
	want := []byte("Bearer " + token)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := []byte(r.Header.Get("Authorization"))
		if subtle.ConstantTimeCompare(got, want) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func credentialsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "hattrick-mcp", "credentials.json"), nil
}

// account holds the CHPP client, which exists once the manager has logged in.
type account struct {
	path string
	mu   sync.Mutex
	ht   hattrick.Client
	req  *oauth.RequestToken
}

func newAccount(path string) (*account, error) {
	a := &account{path: path}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return a, nil
	}
	if err != nil {
		return nil, err
	}
	var cred credentials
	if err := json.Unmarshal(b, &cred); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	a.ht, err = newClient(cred)
	return a, err
}

func newClient(cred credentials) (hattrick.Client, error) {
	return hattrick.NewClient(hattrick.AuthConfig{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken:    cred.AccessToken,
		AccessSecret:   cred.AccessSecret,
	})
}

func newConsumer() *oauth.Consumer {
	return oauth.NewConsumer(consumerKey, consumerSecret, oauth.ServiceProvider{
		RequestTokenUrl:   "https://chpp.hattrick.org/oauth/request_token.ashx",
		AuthorizeTokenUrl: "https://chpp.hattrick.org/oauth/authorize.aspx",
		AccessTokenUrl:    "https://chpp.hattrick.org/oauth/access_token.ashx",
	})
}

func (a *account) client() (hattrick.Client, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.ht == nil {
		return nil, errors.New("not logged in: call the login tool")
	}
	return a.ht, nil
}

// scopes are the write permissions requested at login. Read access needs none.
var scopes = []string{"set_matchorder", "set_training", "place_bid", "manage_youthplayers", "manage_challenges"}

// login starts the CHPP OAuth 1.0a out-of-band flow and returns the URL where the manager grants access.
func (a *account) login() (string, error) {
	req, url, err := newConsumer().GetRequestTokenAndUrl("oob")
	if err != nil {
		return "", err
	}
	url += "&scope=" + strings.Join(scopes, ",")
	a.mu.Lock()
	a.req = req
	a.mu.Unlock()
	return url, nil
}

// completeLogin exchanges the verification code for an access token and stores it.
func (a *account) completeLogin(code string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.req == nil {
		return errors.New("no login in progress: call the login tool first")
	}
	tok, err := newConsumer().AuthorizeToken(a.req, strings.TrimSpace(code))
	if err != nil {
		return err
	}
	cred := credentials{tok.Token, tok.Secret}
	b, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(a.path), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(a.path, b, 0o600); err != nil {
		return err
	}
	ht, err := newClient(cred)
	if err != nil {
		return err
	}
	a.ht, a.req = ht, nil
	return nil
}

func serve(path string) error {
	a, err := newAccount(path)
	if err != nil {
		return err
	}
	return newServer(a).Run(context.Background(), &mcp.StdioTransport{})
}
