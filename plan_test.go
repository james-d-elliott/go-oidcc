package oidcc

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/url"
	"slices"
	"testing"
)

func TestDeleteComprehensive(t *testing.T) {
	client := NewAPIClient(nil, nil, nil)

	plans, err := client.GetPlans(0, 0, 100, false, "Comprehensive:", "")
	if err != nil {
		panic(err)
	}

	for _, plan := range plans.Data {
		if _, err = client.DeletePlan(plan); err != nil {

		}
	}
}

func TestCreateComprehensive(t *testing.T) {
	var (
		plans []*PlanMetadata
		err   error
	)

	issuer := "https://auth.jameselliott.dev"
	secret := "dzT3itYvQMLAPPdvYmfRLyyVgvQe9tswEsDNw4dPpMzooSCL72ucYBpxLhPo"

	client := NewAPIClient(nil, nil, nil)

	if plans, err = NewComprehensiveDiscoveryPlanAll(secret, issuer, SummaryPublish); err != nil {
		panic(err)
	}

	responses, err := client.PostPlans(plans...)
	if err != nil {
		panic(err)
	}

	fmt.Println(responses)
}

func TestCreate(t *testing.T) {
	var (
		plans []*PlanMetadata
		err   error
	)

	issuer := "https://auth.jameselliott.dev"
	secret := "dzT3itYvQMLAPPdvYmfRLyyVgvQe9tswEsDNw4dPpMzooSCL72ucYBpxLhPo"

	client := NewAPIClient(nil, nil, nil)

	if _, err = client.DeletePlans(); err != nil {
		panic(err)
	}

	if plans, err = NewPlansAll(issuer, secret, SummaryPublish); err != nil {
		panic(err)
	}

	slices.Reverse(plans)

	responses, err := client.PostPlans(plans...)
	if err != nil {
		panic(err)
	}

	fmt.Println(responses)
}

func TestGetAndMarshalClients(t *testing.T) {
	client := NewAPIClient(nil, nil, nil)

	plans, err := client.GetPlans(0, 0, 40, false, "Comprehensive: ", "")
	if err != nil {
		panic(err)
	}

	var clients []Client

	for _, plan := range plans.Data {
		clients = append(clients, plan.GetClients(&url.URL{Scheme: "https", Host: "localhost:8443"})...)
	}

	out, err := yaml.Marshal(&ClientData{IdentityProviders: ClientDataIdentityProviders{OpenIDConnect: ClientDataIdentityProvidersOpenIDConnect{Clients: clients}}})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}

type ClientDataIdentityProviders struct {
	OpenIDConnect ClientDataIdentityProvidersOpenIDConnect `yaml:"oidc"`
}

type ClientDataIdentityProvidersOpenIDConnect struct {
	Clients []Client `yaml:"clients"`
}

type ClientData struct {
	IdentityProviders ClientDataIdentityProviders `yaml:"identity_providers"`
}
