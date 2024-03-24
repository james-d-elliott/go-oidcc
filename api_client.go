package oidcc

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func NewClient(tlsConfig *tls.Config) *http.Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client
}

func NewAPIClient(root *url.URL, headers http.Header, tlsConfig *tls.Config) *APIClient {
	if root == nil {
		root = &url.URL{Scheme: "https", Host: "localhost:8443", Path: "/api"}
	}

	if tlsConfig == nil {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &APIClient{
		root:    root,
		client:  NewClient(tlsConfig),
		headers: headers,
	}
}

type APIClient struct {
	root    *url.URL
	client  *http.Client
	headers http.Header
}

func (c *APIClient) NewRequestURI(query url.Values, path ...string) (uri *url.URL) {
	uri = c.root.JoinPath(path...)

	if len(query) != 0 {
		uri.RawQuery = query.Encode()
	}

	return uri
}

func (c *APIClient) NewRequestWithContext(ctx context.Context, method string, uri *url.URL, body io.Reader) (r *http.Request, err error) {
	return http.NewRequestWithContext(ctx, method, uri.String(), body)
}

func (c *APIClient) DoContext(ctx context.Context, method string, body io.Reader, query url.Values, path ...string) (resp *http.Response, err error) {
	uri := c.NewRequestURI(query, path...)

	req, err := c.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}

	for key, values := range c.headers {
		req.Header[key] = values
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.client.Do(req)
}

// https://auth.jameselliott.dev/api/oidc/authorization?client_id=conformance-certification-profile-basic-1&redirect_uri=https://localhost:8443/test/a/certification-profile-basic/callback&scope=openid&state=pI4BXVlX9f&nonce=2Qfy3oeEc5&response_type=code
func (c *APIClient) PostPlans(plans ...*PlanMetadata) (responses []*PlanCreateResponse, err error) {
	var lastErr error
	for _, plan := range plans {
		response, err := c.PostPlan(plan)
		if err != nil {
			lastErr = err

			continue
		}

		responses = append(responses, response)
	}

	return responses, lastErr
}

func (c *APIClient) PostPlan(plan *PlanMetadata) (response *PlanCreateResponse, err error) {
	query := url.Values{}

	query.Set("planName", plan.Name)

	if plan.Variant != nil {
		var data []byte

		if data, err = json.Marshal(plan.Variant); err != nil {
			return nil, err
		}

		query.Add("variant", string(data))
	}

	var form []byte

	if form, err = json.Marshal(plan.Config); err != nil {
		return nil, err
	}

	resp, err := c.DoContext(context.Background(), http.MethodPost, bytes.NewReader(form), query, "plan")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("request failed and failed to read body with error: %w", err)
		}

		return nil, fmt.Errorf("request failed with data: %s", data)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	response = &PlanCreateResponse{}

	if err = decoder.Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *APIClient) DeletePlan(plan PlanMetadata) (ok bool, err error) {
	if plan.ID == "" {
		return false, fmt.Errorf("plan has no id")
	}

	resp, err := c.DoContext(context.Background(), http.MethodDelete, nil, nil, "plan", plan.ID)
	if err != nil {
		return false, err
	}

	return resp.StatusCode >= 200 && resp.StatusCode < 300, nil
}

func (c *APIClient) DeletePlans() (ok bool, err error) {
	plans, err := c.GetPlans(0, 0, 40, false, "", "")
	if err != nil {
		panic(err)
	}

	for _, plan := range plans.Data {
		if _, err = c.DeletePlan(plan); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (c *APIClient) GetPlans(draw, start, length int, public bool, search, order string) (plans *PlanMetadataResponse, err error) {
	query := url.Values{}

	query.Set("public", strconv.FormatBool(public))

	if draw != 0 {
		query.Set("draw", strconv.Itoa(draw))
	}

	if start != 0 {
		query.Set("start", strconv.Itoa(start))
	}

	if length != 0 {
		query.Set("length", strconv.Itoa(length))
	}

	if len(search) != 0 {
		query.Set("search", search)
	}

	if len(order) != 0 {
		query.Set("order", order)
	}

	r, err := c.DoContext(context.Background(), http.MethodGet, nil, query, "plan")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	plans = &PlanMetadataResponse{}

	if err = decoder.Decode(plans); err != nil {
		return nil, err
	}

	return plans, nil
}
