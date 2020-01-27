package main

import (
  "fmt"
  "net/http"
  "net/url"
  "github.com/google/uuid"
  "github.com/gorilla/mux"
  oidc "github.com/coreos/go-oidc"
  "golang.org/x/oauth2"
  "context"
  "encoding/json"
)

var ctx = context.Background()

// AppToken contains oauth2 Token and ID Token Claims
type AppToken struct {
  OAuth2Token *oauth2.Token
  IDTokenClaims *json.RawMessage
}

func getUUID() string  {
  return uuid.New().String()
}

func homeHandler(state string, conf oauth2.Config) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
      http.Redirect(w, r,conf.AuthCodeURL(state), http.StatusFound)
  }
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
  redirectURL, err := url.Parse("http://172.17.0.1:5000/")
  if err != nil {
    fmt.Println(err.Error())
  }
  logoutURL := fmt.Sprintf("http://172.17.0.2:8080/auth/realms/demo/protocol/openid-connect/logout?redirect_uri=%s", redirectURL)
  fmt.Println(logoutURL)
  http.Redirect(w, r, logoutURL, http.StatusTemporaryRedirect)
}

func callbackHandler(v *oidc.IDTokenVerifier, config oauth2.Config, state string) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    if r.URL.Query().Get("state") != state {
      http.Error(w, "State did not match", http.StatusBadRequest)
      return
    }

    oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
    if err != nil {
      http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
      return
    }

    rawIDToken, ok := oauth2Token.Extra("id_token").(string)
    if !ok {
      http.Error(w, "No id_token was found in oauth2 token", http.StatusInternalServerError)
      return
    }

    idToken, err := v.Verify(ctx, rawIDToken)
    if err != nil {
      http.Error(w, "Failed to verify ID token", http.StatusInternalServerError)
      return
    }

    resp := AppToken {
      OAuth2Token: oauth2Token,
      IDTokenClaims: new(json.RawMessage),
    }

    if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    data, err := json.MarshalIndent(resp, "", "  ")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.Write(data)
  }
}

func main()  {
  provider, err := oidc.NewProvider(ctx, "http://172.17.0.2:8080/auth/realms/demo")
  if err != nil {
    fmt.Printf("Error creating provider:  %v", err)
  }
  oidcConfig := &oidc.Config{
    ClientID: "demo-client",
  }
  verifier := provider.Verifier(oidcConfig)
  config := oauth2.Config{
    ClientID: "demo-client",
    ClientSecret: "07eb1b16-cf7a-44ab-9f16-8826d57feeed",
    Endpoint: provider.Endpoint(),
    RedirectURL: "http://172.17.0.1:5000/auth/demo/callback",
    Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
  }

  state := getUUID()

  r := mux.NewRouter()
  r.HandleFunc("/", homeHandler(state, config))
  r.HandleFunc("/logout", logoutHandler)
  r.HandleFunc("/auth/demo/callback", callbackHandler(verifier, config, state))
  if err := http.ListenAndServe(":5000", r); err != nil {
    fmt.Printf("Error starting server: %s", err.Error())
  }
}
