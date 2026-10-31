// Command hattrick-mcp is an MCP server that exposes Hattrick CHPP data to AI assistants.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	case "login":
		err = login(path)
	case "serve":
		err = serve(path)
	case "version":
		fmt.Println(version)
	default:
		err = fmt.Errorf("unknown command %q (use login, serve or version)", cmd)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func credentialsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "hattrick-mcp", "credentials.json"), nil
}

// login runs the CHPP OAuth 1.0a out-of-band flow and stores the access token.
func login(path string) error {
	c := oauth.NewConsumer(consumerKey, consumerSecret, oauth.ServiceProvider{
		RequestTokenUrl:   "https://chpp.hattrick.org/oauth/request_token.ashx",
		AuthorizeTokenUrl: "https://chpp.hattrick.org/oauth/authorize.aspx",
		AccessTokenUrl:    "https://chpp.hattrick.org/oauth/access_token.ashx",
	})
	reqToken, url, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		return err
	}
	fmt.Println("1. Open this URL and grant access:")
	fmt.Println("   " + url)
	fmt.Print("2. Paste the verification code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return err
	}
	tok, err := c.AuthorizeToken(reqToken, strings.TrimSpace(code))
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(credentials{tok.Token, tok.Secret}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		return err
	}
	fmt.Println("Saved credentials to " + path)
	return nil
}

func serve(path string) error {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("not logged in: run `hattrick-mcp login` first")
	}
	if err != nil {
		return err
	}
	var cred credentials
	if err := json.Unmarshal(b, &cred); err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	ht, err := hattrick.NewClient(hattrick.AuthConfig{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken:    cred.AccessToken,
		AccessSecret:   cred.AccessSecret,
	})
	if err != nil {
		return err
	}
	return newServer(ht).Run(context.Background(), &mcp.StdioTransport{})
}
