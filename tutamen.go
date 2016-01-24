package tutamen

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	APPROVED = "approved"
	PENDING  = "pending"
	DENIED   = "denied"
	DEFAULT_RETRY_INTERVAL = 5 * time.Second
)

type client struct {
	http      *http.Client
	certpath   string
	keypath    string
	ac_base    string
	ss_base    string
}

func NewClientV1(certpath, keypath, ac_server, ss_server string) (*client, error) {

	c := new(client)

	// TODO: validation
	c.certpath  = certpath
	c.keypath   = keypath
	c.ac_base   = "https://" + ac_server + "/api/v1/"
	c.ss_base   = "https://" + ss_server + "/api/v1/"

	x509, err := tls.LoadX509KeyPair(c.certpath, c.keypath)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport {
		TLSClientConfig: &tls.Config {
			Certificates: []tls.Certificate{ x509 },
		},
	}

	c.http = &http.Client{Transport: tr}

	return c, nil
}

func (c *client) GetAuthorization(collection string) (string, error) {

	type auth_request struct {
		Objperm string `json:"objperm"`
		Objtype string `json:"objtype"`
		Objuid  string `json:"objuid"`
	}

	type auth_reply struct {
		Authorizations []string `json:"authorizations"`
	}

	// Get

	url := c.ac_base + "authorizations/"

	tx_body, err := json.Marshal(auth_request{
		Objperm: "read",
		Objtype: "collection",
		Objuid:  collection,
	})
	if err != nil {
		return "", err
	}

	resp, err := c.http.Post(url, "application/json", bytes.NewBuffer(tx_body))
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	rx_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse

	var m auth_reply
	err = json.Unmarshal(rx_body, &m)
	if err != nil {
		return "", err
	}

	if len(m.Authorizations) == 0 {
		return "", errors.New("0 authorizations returned")
	} else if len(m.Authorizations) > 1 {
		return "", errors.New("multiple authorizations not yet implemented")
	}

	// Return

	return m.Authorizations[0], nil
}

func (c *client) WaitForToken(authorization string) (string, error) {

	url := c.ac_base + "authorizations/" + authorization

	type token_reply struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	}
	var m token_reply

	for {
		resp, err := c.http.Get(url)
		if err != nil {
			return "", err
		} else {
			// TODO: is this okay, when does it happen?
			defer resp.Body.Close()
		}

		rx_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(rx_body, &m)
		if err != nil {
			return "", err
		}

		switch m.Status {
		case APPROVED:
			return m.Token, nil
		case PENDING:
			time.Sleep(DEFAULT_RETRY_INTERVAL)
		case DENIED:
			return "", errors.New("Token denied")
		default:
			return "", errors.New("Unknown token status: " + m.Status)
		}
	}
}

func (c *client) GetSecret(collection, secret, token string) (string, error) {

	url := c.ss_base + "collections/" + collection + "/secrets/" + secret + "/versions/latest/"

	type secret_reply struct {
		Data string `json:"data"`
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("tutamen-tokens", token)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	rx_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var m secret_reply
	err = json.Unmarshal(rx_body, &m)
	if err != nil {
		return "", err
	}

	return m.Data, nil
}

