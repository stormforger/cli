package api

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"fmt"
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
	HTTPClient *http.Client
	APIDomain string
	JWT string
	UserAgent string
}

func NewClient(apiDomain string, jwtToken string) *Client {
	return &Client{
		HTTPClient: defaultHTTPClient(),
		APIDomain: apiDomain,
		JWT: jwtToken,
		UserAgent: "StormForger CLI (https://stormforger.com)",
	}
}

func (c *Client) Ping() {
	req, err := http.NewRequest("GET", c.APIDomain + "/authenticated_ping", nil)

	// TODO how to set these on all requests?
	req.Header.Add("Authorization", "Bearer "+ c.JWT)
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "  ")

	fmt.Printf("%s\n\n %s", resp.Status, prettyJSON.Bytes())
}

