package oidcc

import (
	"fmt"
	"net/url"
	"strings"
)

func NewPlansAll(issuer, secret string, publish Publish) (plans []*PlanMetadata, err error) {
	var plan *PlanMetadata

	if plan, err = NewCertificationProfileBasicDiscoveryPlan("certification-profile-basic", "Certification Profile: Basic", secret, issuer, publish); err != nil {
		return nil, err
	}

	plans = append(plans, plan)

	if plan, err = NewCertificationProfileHybridDiscoveryPlan("certification-profile-hybrid", "Certification Profile: Hybrid", secret, issuer, publish); err != nil {
		return nil, err
	}

	plans = append(plans, plan)

	if plan, err = NewCertificationProfileImplicitDiscoveryPlan("certification-profile-implicit", "Certification Profile: Implicit", secret, issuer, publish); err != nil {
		return nil, err
	}

	plans = append(plans, plan)

	if plan, err = NewCertificationProfileFormPostBasicDiscoveryPlan("certification-profile-formpost-basic", "Certification Profile: Form Post Basic", secret, issuer, publish); err != nil {
		return nil, err
	}

	plans = append(plans, plan)

	if plan, err = NewCertificationProfileFormPostHybridDiscoveryPlan("certification-profile-formpost-hybrid", "Certification Profile: Form Post Hybrid", secret, issuer, publish); err != nil {
		return nil, err
	}

	plans = append(plans, plan)

	if plan, err = NewCertificationProfileFormPostImplicitDiscoveryPlan("certification-profile-formpost-implicit", "Certification Profile: Form Post Implicit", secret, issuer, publish); err != nil {
		return nil, err
	}

	plans = append(plans, plan)

	if plan, err = NewCertificationProfileConfigDiscoveryPlan("certification-profile-config", "Certification Profile: Config", issuer, publish); err != nil {
		return nil, err
	}

	plans = append(plans, plan)

	return plans, nil
}

func NewPlanDiscovery(name string, variant *PlanVariant, publish Publish, alias, description, issuer string, client, client2, clientSecretPost *PlanClient) (plan *PlanMetadata, err error) {
	var (
		discoveryURI *url.URL
	)

	if discoveryURI, err = url.ParseRequestURI(issuer); err != nil {
		return nil, err
	}

	discoveryURI = discoveryURI.JoinPath(".well-known", "openid-configuration")

	plan = &PlanMetadata{
		Name: name,
		Config: &PlanConfig{
			Alias:       alias,
			Description: description,
			Server: &PlanServer{
				DiscoveryURL: discoveryURI.String(),
			},
			Client:           client,
			Client2:          client2,
			ClientSecretPost: clientSecretPost,
		},
		Publish: publish.String(),
		Variant: variant,
	}

	return plan, nil
}

func NewCertificationProfileStandardDiscoveryPlan(name, alias, description, secret, issuer string, publish Publish) (plan *PlanMetadata, err error) {
	client := &PlanClient{
		ClientID:     "conformance-" + alias + "-1",
		ClientSecret: secret,
	}

	client2 := &PlanClient{
		ClientID:     "conformance-" + alias + "-2",
		ClientSecret: secret,
	}

	clientPost := &PlanClient{
		ClientID:     "conformance-" + alias + "-post",
		ClientSecret: secret,
	}

	variant := &PlanVariant{
		ServerMetadata:     "discovery",
		ClientRegistration: "static_client",
	}

	return NewPlanDiscovery(name, variant, publish, alias, description, issuer, client, client2, clientPost)
}

func NewCertificationProfileBasicDiscoveryPlan(alias, description, secret, issuer string, publish Publish) (plan *PlanMetadata, err error) {
	return NewCertificationProfileStandardDiscoveryPlan("oidcc-basic-certification-test-plan", alias, description, secret, issuer, publish)
}

func NewCertificationProfileFormPostBasicDiscoveryPlan(alias, description, secret, issuer string, publish Publish) (plan *PlanMetadata, err error) {
	return NewCertificationProfileStandardDiscoveryPlan("oidcc-formpost-basic-certification-test-plan", alias, description, secret, issuer, publish)
}

func NewCertificationProfileFormPostHybridDiscoveryPlan(alias, description, secret, issuer string, publish Publish) (plan *PlanMetadata, err error) {
	return NewCertificationProfileStandardDiscoveryPlan("oidcc-formpost-hybrid-certification-test-plan", alias, description, secret, issuer, publish)
}

func NewCertificationProfileFormPostImplicitDiscoveryPlan(alias, description, secret, issuer string, publish Publish) (plan *PlanMetadata, err error) {
	client := &PlanClient{
		ClientID:     "conformance-" + alias + "-1",
		ClientSecret: secret,
	}

	variant := &PlanVariant{
		ServerMetadata:     "discovery",
		ClientRegistration: "static_client",
	}

	return NewPlanDiscovery("oidcc-formpost-implicit-certification-test-plan", variant, publish, alias, description, issuer, client, nil, nil)
}

func NewCertificationProfileHybridDiscoveryPlan(alias, description, secret, issuer string, publish Publish) (plan *PlanMetadata, err error) {
	return NewCertificationProfileStandardDiscoveryPlan("oidcc-hybrid-certification-test-plan", alias, description, secret, issuer, publish)
}

func NewCertificationProfileImplicitDiscoveryPlan(alias, description, secret, issuer string, publish Publish) (plan *PlanMetadata, err error) {
	client := &PlanClient{
		ClientID:     "conformance-" + alias + "-1",
		ClientSecret: secret,
	}

	variant := &PlanVariant{
		ServerMetadata:     "discovery",
		ClientRegistration: "static_client",
	}

	return NewPlanDiscovery("oidcc-implicit-certification-test-plan", variant, publish, alias, description, issuer, client, nil, nil)
}

func NewCertificationProfileConfigDiscoveryPlan(alias, description, issuer string, publish Publish) (plan *PlanMetadata, err error) {
	return NewPlanDiscovery("oidcc-config-certification-test-plan", nil, publish, alias, description, issuer, nil, nil, nil)
}

var clientAuthTypes = []string{"none", "client_secret_basic", "client_secret_post", "client_secret_jwt"}
var responseTypes = []string{"code", "id_token", "id_token token", "code id_token", "code token", "code id_token token"}
var responseModes = []string{"default", "form_post"}

func responseTypeToFlowDescription(in string) string {
	switch in {
	case "code":
		return "Authorization Code"
	case "id_token":
		return "Implicit"
	case "id_token token":
		return "Implicit (Token)"
	case "code id_token":
		return "Hybrid (ID Token)"
	case "code token":
		return "Hybrid (Token)"
	case "code id_token token":
		return "Hybrid (Both)"
	default:
		return ""
	}
}

func clientAuthTypeToDescription(in string) string {
	switch in {
	case "none":
		return "Public"
	case "client_secret_basic":
		return "Basic"
	case "client_secret_post":
		return "Post"
	case "client_secret_jwt":
		return "JWT"
	default:
		return ""
	}
}

func NewComprehensiveDiscoveryPlanAll(secret, issuer string, publish Publish) (plans []*PlanMetadata, err error) {
	for _, clientAuthType := range clientAuthTypes {
		for _, responseType := range responseTypes {
			for _, responseMode := range responseModes {
				alias := fmt.Sprintf("conformance-%s-%s", strings.ReplaceAll(clientAuthType, "client_secret_", ""), strings.ReplaceAll(responseType, " ", "-"))
				fp := ""
				alg := ""

				switch responseMode {
				case "form_post":
					alias += "formpost"
					fp = " Form Post"
				}

				s := secret
				switch clientAuthType {
				case "none":
					s = ""
				case "client_secret_jwt":
					alg = "HS256"
				}

				plan, err := NewComprehensiveDiscoveryPlan(alias, fmt.Sprintf("Comprehensive: %s %s%s", responseTypeToFlowDescription(responseType), clientAuthTypeToDescription(clientAuthType), fp), s, alg, issuer, clientAuthType, responseType, responseMode, publish)
				if err != nil {
					return nil, err
				}

				plan.Variant.ServerMetadata = ""

				plans = append(plans, plan)
			}
		}
	}

	return plans, nil
}

func NewComprehensiveDiscoveryPlan(alias, description, secret, secretAlg, issuer, clientAuthType, responseType, responseMode string, publish Publish) (plan *PlanMetadata, err error) {
	client := &PlanClient{
		ClientID:           "conformance-" + alias + "-1",
		ClientSecret:       secret,
		ClientSecretJWTAlg: secretAlg,
	}

	client2 := &PlanClient{
		ClientID:           "conformance-" + alias + "-2",
		ClientSecret:       secret,
		ClientSecretJWTAlg: secretAlg,
	}

	variant := &PlanVariant{
		ServerMetadata:     "discovery",
		ClientRegistration: "static_client",
		ClientAuthType:     clientAuthType,
		ResponseType:       responseType,
		ResponseMode:       responseMode,
	}

	return NewPlanDiscovery("oidcc-test-plan", variant, publish, alias, description, issuer, client, client2, nil)
}
