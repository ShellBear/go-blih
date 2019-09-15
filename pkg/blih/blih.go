package blih

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	DefaultApiBaseURL = "https://blih.epitech.eu/"
	DefaultUserAgent  = "go-blih"
)

type Service struct {
	User    string
	token   string
	ctx     context.Context
	BaseURL string
	client  *http.Client

	Utils      *UtilsService
	Repository *RepositoryService
	SSHKey     *SSHKeyService
}

func New(user, token string, ctx context.Context, baseURL string) *Service {
	h := sha512.New()
	h.Write([]byte(token))

	s := &Service{
		User:    user,
		token:   hex.EncodeToString(h.Sum(nil)),
		BaseURL: baseURL,
		ctx:     ctx,
		client:  http.DefaultClient,
	}

	s.Utils = NewUtilsService(s)
	s.Repository = NewRepositoryService(s)
	s.SSHKey = NewSSHKeyService(s)

	return s
}

type SignedData struct {
	User      string      `json:"user"`
	Signature string      `json:"signature"`
	Data      interface{} `json:"data,omitempty"`
}

func (c *Service) SignData(data interface{}) (io.Reader, error) {
	signedData := SignedData{User: c.User, Data: data}
	hash := hmac.New(sha512.New, []byte(c.token))

	hash.Write([]byte(signedData.User))
	if data != nil {
		b, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return nil, err
		}

		hash.Write(b)
	}

	signedData.Signature = hex.EncodeToString(hash.Sum(nil))

	b, err := json.Marshal(signedData)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func (c *Service) NewRequest(endpoint, method string, body interface{}) (*http.Request, error) {
	reqURL, err := url.ParseRequestURI(c.BaseURL + endpoint)
	if err != nil {
		return nil, err
	}

	signedData, err := c.SignData(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(c.ctx, method, reqURL.String(), signedData)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", DefaultUserAgent)

	return req, nil
}

func (c *Service) SendRequest(endpoint, method string, body interface{}) (*http.Response, error) {
	req, err := c.NewRequest(endpoint, method, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
