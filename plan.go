package oidcc

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type PlanMetadataResponse struct {
	Draw            int            `json:"draw"`
	RecordsTotal    int            `json:"recordsTotal"`
	RecordsFiltered int            `json:"recordsFiltered"`
	Data            []PlanMetadata `json:"data"`
}

type PlanMetadata struct {
	ID                       string       `json:"_id"`
	Name                     string       `json:"planName"`
	Variant                  *PlanVariant `json:"variant,omitempty"`
	Config                   *PlanConfig  `json:"config,omitempty"`
	Started                  time.Time    `json:"started,omitempty"`
	Owner                    *PlanOwner   `json:"owner,omitempty"`
	Description              string       `json:"description,omitempty"`
	CertificationProfileName any          `json:"certificationProfileName,omitempty"`
	Modules                  []PlanModule `json:"modules,omitempty"`
	Version                  string       `json:"version,omitempty"`
	Summary                  string       `json:"summary,omitempty"`
	Publish                  string       `json:"publish,omitempty"`
	Immutable                any          `json:"immutable,omitempty"`
}

type PlanConfig struct {
	Alias            string      `json:"alias,omitempty"`
	Description      string      `json:"description,omitempty"`
	Publish          string      `json:"publish,omitempty"`
	Server           *PlanServer `json:"server,omitempty"`
	Client           *PlanClient `json:"client,omitempty"`
	Client2          *PlanClient `json:"client2,omitempty"`
	ClientSecretPost *PlanClient `json:"client_secret_post,omitempty"`
}

type PlanOwner struct {
	Sub string `json:"sub"`
	Iss string `json:"iss"`
}

type PlanModule struct {
	TestModule string       `json:"testModule"`
	Variant    *PlanVariant `json:"variant,omitempty"`
	Instances  []any        `json:"instances"`
}

type PlanPage struct {
	Draw   int    `json:"draw"`
	Start  int    `json:"start"`
	Length int    `json:"length"`
	Search string `json:"search,omitempty"`
	Order  string `json:"order,omitempty"`
}

type PlanCreateResponse struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Modules []PlanModule `json:"modules,omitempty"`
}

type PlanCreateResponseModule struct {
	TestModule string `json:"testModule"`
}

type Client struct {
	ClientID                string   `yaml:"client_id"`
	ClientSecret            string   `yaml:"client_secret,omitempty"`
	Public                  bool     `yaml:"public"`
	RedirectURIs            []string `yaml:"redirect_uris"`
	Scopes                  []string `yaml:"scopes"`
	GrantTypes              []string `yaml:"grant_types"`
	ResponseTypes           []string `yaml:"response_types"`
	ResponseModes           []string `yaml:"response_modes"`
	AuthorizationPolicy     string   `yaml:"authorization_policy"`
	ConsentMode             string   `yaml:"consent_mode"`
	RequestObjectSigningAlg string   `yaml:"request_object_signing_alg"`
	TokenEndpointAuthMethod string   `yaml:"token_endpoint_auth_method"`
}

type PlanServer struct {
	ACRValues             string `json:"acr_values,omitempty"`
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
	Issuer                string `json:"issuer,omitempty"`
	JSONWebKeysURI        string `json:"jwks_uri,omitempty"`
	DiscoveryURL          string `json:"discoveryUrl,omitempty"`
	LoginHint             string `json:"login_hint,omitempty"`
}

type PlanClient struct {
	ClientID           string `json:"client_id,omitempty"`
	ClientSecret       string `json:"client_secret,omitempty"`
	ClientSecretJWTAlg string `json:"client_secret_jwt_alg,omitempty"`
}

type PlanVariant struct {
	ServerMetadata     string `json:"server_metadata,omitempty"`
	ClientRegistration string `json:"client_registration,omitempty"`
	ClientAuthType     string `json:"client_auth_type,omitempty"`
	ResponseType       string `json:"response_type,omitempty"`
	ResponseMode       string `json:"response_mode,omitempty"`
}

func (p PlanMetadata) GetClients(root *url.URL) []Client {
	var clients []Client

	redirectURI := root.JoinPath("test", "a", p.Config.Alias, "callback")

	if p.Config.Client != nil {
		client := Client{
			ClientID:     p.Config.Client.ClientID,
			ClientSecret: fmt.Sprintf("$plaintext$%s", p.Config.Client.ClientSecret),
			Public:       false,
			RedirectURIs: []string{
				redirectURI.String(),
			},
			Scopes: []string{
				"openid",
				"offline_access",
				"profile",
				"email",
				"phone",
				"address",
				"all",
			},
			GrantTypes: []string{
				"authorization_code",
				"refresh_token",
			},
			ResponseTypes: []string{
				"code",
			},
			ResponseModes: []string{
				"form_post",
				"query",
				"fragment",
				"jwt",
				"form_post.jwt",
				"query.jwt",
				"fragment.jwt",
			},
			AuthorizationPolicy:     "one_factor",
			ConsentMode:             "implicit",
			TokenEndpointAuthMethod: "client_secret_basic",
			RequestObjectSigningAlg: "none",
		}

		switch {
		case strings.Contains(p.Name, "-hybrid-"):
			client.ResponseTypes = append(client.ResponseTypes, "code id_token", "code token", "code id_token token")
			client.GrantTypes = append(client.GrantTypes, "implicit")
		case strings.Contains(p.Name, "-implicit-"):
			client.ResponseTypes = append(client.ResponseTypes, "id_token", "token", "id_token token")
			client.GrantTypes = append(client.GrantTypes, "implicit")
		}

		clients = append(clients, client)
	}

	if p.Config.Client2 != nil {
		client := Client{
			ClientID:     p.Config.Client2.ClientID,
			ClientSecret: fmt.Sprintf("$plaintext$%s", p.Config.Client2.ClientSecret),
			Public:       false,
			RedirectURIs: []string{
				redirectURI.String(),
			},
			Scopes: []string{
				"openid",
				"offline_access",
				"profile",
				"email",
				"phone",
				"address",
				"all",
			},
			GrantTypes: []string{
				"authorization_code",
				"refresh_token",
			},
			ResponseTypes: []string{
				"code",
			},
			ResponseModes: []string{
				"form_post",
				"query",
				"fragment",
				"jwt",
				"form_post.jwt",
				"query.jwt",
				"fragment.jwt",
			},
			AuthorizationPolicy:     "one_factor",
			ConsentMode:             "implicit",
			TokenEndpointAuthMethod: "client_secret_basic",
			RequestObjectSigningAlg: "none",
		}

		switch {
		case strings.Contains(p.Name, "-hybrid-"):
			client.ResponseTypes = append(client.ResponseTypes, "code id_token", "code token", "code id_token token")
			client.GrantTypes = append(client.GrantTypes, "implicit")
		case strings.Contains(p.Name, "-implicit-"):
			client.ResponseTypes = append(client.ResponseTypes, "id_token", "token", "id_token token")
			client.GrantTypes = append(client.GrantTypes, "implicit")
		}

		clients = append(clients, client)
	}

	if p.Config.ClientSecretPost != nil {
		client := Client{
			ClientID:     p.Config.ClientSecretPost.ClientID,
			ClientSecret: fmt.Sprintf("$plaintext$%s", p.Config.ClientSecretPost.ClientSecret),
			Public:       false,
			RedirectURIs: []string{
				redirectURI.String(),
			},
			Scopes: []string{
				"openid",
				"offline_access",
				"profile",
				"email",
				"phone",
				"address",
				"all",
			},
			GrantTypes: []string{
				"authorization_code",
				"refresh_token",
			},
			ResponseTypes: []string{
				"code",
			},
			ResponseModes: []string{
				"form_post",
				"query",
				"fragment",
				"jwt",
				"form_post.jwt",
				"query.jwt",
				"fragment.jwt",
			},
			AuthorizationPolicy:     "one_factor",
			ConsentMode:             "implicit",
			TokenEndpointAuthMethod: "client_secret_post",
			RequestObjectSigningAlg: "none",
		}

		switch {
		case strings.Contains(p.Name, "-hybrid-"):
			client.ResponseTypes = append(client.ResponseTypes, "code id_token", "code token", "code id_token token")
			client.GrantTypes = append(client.GrantTypes, "implicit")
		case strings.Contains(p.Name, "-implicit-"):
			client.ResponseTypes = append(client.ResponseTypes, "id_token", "token", "id_token token")
			client.GrantTypes = append(client.GrantTypes, "implicit")
		}

		clients = append(clients, client)
	}

	return clients
}
