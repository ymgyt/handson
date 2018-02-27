package main

import (
	"bytes"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"

	githuboauth "golang.org/x/oauth2/github"
)

const (
	EnvOAuthClientID     = "HANDSON_OAUTH_CLIENT_ID"
	EnvOAuthClientSecret = "HANDSON_OAUTH_CLIENT_SECRET"

	OAuthStateLen = 10
)

var (
	OAuthState string
)

func init() {
	OAuthState = OAuthStateString(OAuthStateLen)
}

func NewOAuthConfig() (*oauth2.Config, error) {
	clientID := os.Getenv(EnvOAuthClientID)
	if clientID == "" {
		return nil, errors.Errorf("environment variable %q is empty", EnvOAuthClientID)
	}

	clientSecret := os.Getenv(EnvOAuthClientSecret)
	if clientSecret == "" {
		return nil, errors.Errorf("environment variable %q is empty", EnvOAuthClientSecret)
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"repo", "gist", "user"},
		Endpoint:     githuboauth.Endpoint,
	}, nil
}

const src = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func OAuthStateString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var b bytes.Buffer
	for i := 0; i < n; i++ {
		b.WriteByte(src[rand.Intn(len(src))])
	}
	return b.String()
}

const indexHTML = `<html><body>
Logged in with <a href="/login">Github</a>
</body></html>
`

func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexHTML))
}

type GitHubLoginHandler struct {
	oauthConfig *oauth2.Config
	logger      *zap.Logger
}

func NewGitHubLoginHandler(cfg *oauth2.Config, logger *zap.Logger) http.Handler {
	return &GitHubLoginHandler{
		oauthConfig: cfg,
		logger:      logger,
	}
}

func (glh *GitHubLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := glh.oauthConfig.AuthCodeURL(OAuthState, oauth2.AccessTypeOnline)
	glh.logger.Info("app",
		zap.String("redirect_to", url),
		zap.String("oauth2", "redirect"),
	)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type GitHubCallbackHandler struct {
	oauthConfig *oauth2.Config
	whenFailed  string
	logger      *zap.Logger
}

func NewGitHubCallbackHandler(cfg *oauth.Config, successfullRedirect, failureRedirect string, logger *zap.Logger) http.Handler {
	return &GitHubCallbackHandler{
		config:              config,
		successfullRedirect: successfullRedirect,
		failureRedirect:     failureRedirect,
		logger:              logger,
	}
}

func (gcbh *GitHubCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != OAuthState {
		gcbh.logger.Warn("app",
			zap.String("oauth2", "invalid callback state"),
			zap.String("expcted_state", state),
			zap.String("actual_state", state),
		)
		http.Redirect(w, r, gcbh.failureRedirect, http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := gcbh.oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		gcbh.logger.Warn("app",
			zap.String("oauth2", "oauth2.Config.Exchange() failed"),
			zap.String("code", code),
			zap.Error(err),
		)
		http.Redirect(w, r, gcbh.failureRedirect, http.StatusTemporaryRedirect)
		return
	}

	oauthClient := gcbh.oauthConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		gcbh.logger.Warn("app",
			zap.String("oauth2", "github.Client.users.Get failed"),
			zap.Error(err),
		)
		http.Redirect(w, r, gcbh.failureRedirect, http.StatusTemporaryRedirect)
		return
	}

	gcbh.logger.Info("app",
		zap.String("oauth2", "successfully logged in"),
		zap.String("state", state),
		zap.String("code", code),
		zap.Any("user", user),
	)

	http.Redirect(w, r, gcbh.successfullRedirect, http.StatusTemporaryRedirect)
}

// debug -1
// info   0
// warn   1
// error  2
func GetLogger(level int) (*zap.Logger, error) {
	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(int8(level))),
		Development:      true,
		Encoding:         "console", // or json
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, // CapitalLevelEncoder
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	option := zap.AddStacktrace(zapcore.ErrorLevel)
	return cfg.Build(option)
}
