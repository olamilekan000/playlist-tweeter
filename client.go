package main

import (
	"context"
	"encoding/base64"
	"errors"
	"sync"

	"github.com/go-resty/resty/v2"
)

type (
	client struct {
		resty *resty.Client

		spotify struct {
			secret string
			client string
		}

		mux sync.Mutex
	}

	Client interface {
		// Configure func configures sengrid client
		Configure(opts ...ClientOption)
	}

	ClientOption func(*client)
)

var c *client

func newClient() *client {
	c := client{
		resty: resty.New(),
	}
	return &c
}

func Get() Client {
	if c != nil {
		return c
	}

	c = newClient()

	return c
}

func SetSpotifyClientKey(key string) ClientOption {
	return func(c *client) {
		c.mux.Lock()
		defer c.mux.Unlock()

		c.spotify.client = key
	}
}

func SetSpotifySecretKey(secret string) ClientOption {
	return func(c *client) {
		c.mux.Lock()
		defer c.mux.Unlock()

		c.spotify.secret = secret
	}
}

func (cl *client) Configure(opts ...ClientOption) {
	for _, o := range opts {
		o(cl)
	}
}

func (cl *client) makeGetRequest(ctx context.Context, url string) ([]byte, error) {
	resp, err := cl.resty.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+cache["access_token"].(string)).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	return resp.Body(), nil
}

func (cl *client) makePostRequest(ctx context.Context, url string, data map[string]string) ([]byte, error) {
	resp, err := cl.resty.R().
		SetContext(ctx).
		SetHeader("Authorization", "Basic "+cl.convertToBase64(cl.spotify.client+":"+c.spotify.secret)).
		SetFormData(data).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	return resp.Body(), nil
}

func (cl *client) convertToBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
