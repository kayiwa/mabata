package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/kayiwa/mabata/internal/config"
	"github.com/kayiwa/mabata/internal/models"
)

type OIDC struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   oauth2.Config
	session  *SessionManager
	allowed  map[string]struct{}
	states   sync.Map
}

type Claims struct {
	Email             string   `json:"email"`
	Name              string   `json:"name"`
	Subject           string   `json:"sub"`
	PreferredUsername string   `json:"preferred_username"`
	Groups            []string `json:"groups"`
}

func New(cfg config.Config) (*OIDC, error) {
	ctx := context.Background()
	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", cfg.AzureTenantID)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.AzureClientID})
	o := &OIDC{
		provider: provider,
		verifier: verifier,
		oauth2: oauth2.Config{
			ClientID:     cfg.AzureClientID,
			ClientSecret: cfg.AzureClientSecret,
			RedirectURL:  cfg.AzureRedirectURL,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		session: NewSessionManager(cfg.SessionSecret),
		allowed: map[string]struct{}{},
	}
	for _, g := range cfg.AzureAllowedGroupIDs {
		o.allowed[g] = struct{}{}
	}
	return o, nil
}

func (o *OIDC) AuthCodeURL() (string, error) {
	state, err := randomString(24)
	if err != nil {
		return "", err
	}
	nonce, err := randomString(24)
	if err != nil {
		return "", err
	}
	o.states.Store(state, nonce)
	return o.oauth2.AuthCodeURL(state, oidc.Nonce(nonce)), nil
}

func (o *OIDC) Exchange(ctx context.Context, state, code string) (*models.User, error) {
	v, ok := o.states.Load(state)
	if !ok {
		return nil, fmt.Errorf("invalid state")
	}
	o.states.Delete(state)
	expectedNonce, _ := v.(string)

	tok, err := o.oauth2.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	rawIDToken, ok := tok.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("missing id_token")
	}
	idToken, err := o.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}
	if idToken.Nonce != expectedNonce {
		return nil, fmt.Errorf("invalid nonce")
	}
	var c Claims
	if err := idToken.Claims(&c); err != nil {
		return nil, err
	}
	user := &models.User{
		Email:  firstNonEmpty(c.Email, c.PreferredUsername),
		Name:   c.Name,
		Sub:    c.Subject,
		Groups: c.Groups,
	}
	if !o.authorized(user.Groups) {
		return nil, fmt.Errorf("user is not in an allowed group")
	}
	return user, nil
}

func (o *OIDC) Session() *SessionManager {
	return o.session
}

func (o *OIDC) authorized(groups []string) bool {
	if len(o.allowed) == 0 {
		return true
	}
	for _, g := range groups {
		if _, ok := o.allowed[strings.TrimSpace(g)]; ok {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
