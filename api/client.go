package api

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Dialer func(network, addr string) (net.Conn, error)

func newDialer(fingerprint []byte, skipCAVerification bool) Dialer {
	return func(network, addr string) (net.Conn, error) {
		c, err := tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: skipCAVerification})
		if err != nil {
			return c, err
		}
		connstate := c.ConnectionState()
		keyPinValid := false
		for _, peercert := range connstate.PeerCertificates {
			der, err := x509.MarshalPKIXPublicKey(peercert.PublicKey)
			hash := sha256.Sum256(der)
			// log.Println(peercert.Issuer)
			// log.Printf("%#v\n\n", hash)
			if err != nil {
				log.Fatal(err)
			}
			if bytes.Compare(hash[0:], fingerprint) == 0 {
				// log.Println("Pinned Key found")
				keyPinValid = true
			}
		}
		if keyPinValid == false {
			log.Fatal("TLS Public Key could not be verified!")
		}
		return c, nil
	}
}

func defaultHTTPClient() *http.Client {
	fingerprint := []byte{0x5b, 0x15, 0x6c, 0xda, 0x7b, 0xc3, 0xd, 0x8b, 0xe8, 0x88, 0x57, 0x75, 0xbc, 0x30, 0xc1, 0x84, 0x18, 0x75, 0x6f, 0x2d, 0x3b, 0x81, 0x91, 0xff, 0x34, 0x10, 0xda, 0x13, 0x4a, 0x83, 0x23, 0x9d}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		DialTLS:            newDialer(fingerprint, false),
	}

	return &http.Client{
		Transport: tr,
		Timeout:   2 * time.Minute,
	}
}

type Client struct {
	HTTPClient  *http.Client
	APIEndpoint string
	JWT         string
	UserAgent   string
}

func NewClient(apiEndpoint string, jwtToken string) *Client {
	return &Client{
		HTTPClient:  defaultHTTPClient(),
		APIEndpoint: apiEndpoint,
		JWT:         jwtToken,
		UserAgent:   "StormForger CLI (https://stormforger.com)",
	}
}

// FIXME would be nice to return a struct
//       where we see the status and in case of
//       success also the email address of the
//       authenticated user (useful) to check
//       if we are authenticated as the correct user
func (c *Client) Ping() (bool, error) {
	req, err := http.NewRequest("GET", c.APIEndpoint+"/authenticated_ping", nil)

	// TODO how to set these on all requests?
	req.Header.Add("Authorization", "Bearer "+c.JWT)
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	return (resp.StatusCode == 200), nil
}

func (c *Client) Har(file string) (string, error) {
	// TODO how to pass options, like --skip-assets here?
	//      definiting a struct maybe, but where?
	//      finally: add options here
	extraParams := map[string]string{}

	req, err := newfileUploadRequest(c.APIEndpoint+"/har", extraParams, "har_file", file)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	return string(body), nil
}

func (c *Client) Login(email string, password string) (string, error) {
	data := map[string]string{"email": email, "password": password}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(data)

	req, err := http.NewRequest("POST", c.APIEndpoint+"/beta/user/token", body)

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("Login unsuccessful!")
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(responseBody, &dat); err != nil {
		return "", err
	}

	return dat["jwt"].(string), nil
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
