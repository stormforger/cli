package api

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/stormforger/cli/buildinfo"
)

type dialer func(network, addr string) (net.Conn, error)

func newDialer(fingerprint []byte, skipCAVerification bool) dialer {
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
		Proxy:              http.ProxyURL(getHTTPProxy()),
	}

	return &http.Client{
		Transport: tr,
		Timeout:   2 * time.Minute,
	}
}

func getHTTPProxy() *url.URL {
	httpProxyURL := os.Getenv("HTTP_PROXY")

	if httpProxyURL == "" {
		return nil
	}

	proxyURL, err := url.Parse(httpProxyURL)
	if err != nil {
		log.Fatal(err)
	}

	return proxyURL
}

// Client represents the API client
type Client struct {
	HTTPClient  *http.Client
	APIEndpoint string
	JWT         string
	UserAgent   string
}

// NewClient returns a new initialized API client
func NewClient(apiEndpoint, jwtToken string) *Client {
	return &Client{
		HTTPClient:  defaultHTTPClient(),
		APIEndpoint: apiEndpoint,
		JWT:         jwtToken,
		UserAgent:   fmt.Sprintf("StormForger-CLI/%v (%v; +https://github.com/stormforger/cli)", buildinfo.BuildInfo.Version, buildinfo.BuildInfo.Commit),
	}
}

// Ping performs an authenticated ping to check if the API
// is working and the user is properly authenticated.
func (c *Client) Ping() (bool, []byte, error) {
	return c.fetch("/authenticated_ping")
}

// PingUnauthenticated performs an unauthenticated ping. This
// can be used to see if a connection is possible and/or the API
// is up in general.
func (c *Client) PingUnauthenticated() (bool, []byte, error) {
	return c.fetch("/ping")
}

func newPatchRequest(uri string, params map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

// createGZIPFormFile creates a "mime/multipart".Writer header similar to CreateFormFile but with content-encoding=gzip.
func createGZIPFormFile(w *multipart.Writer, fieldname, filename, contenttype string) (io.Writer, error) {
	replacer := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			replacer.Replace(fieldname), replacer.Replace(filename)))
	h.Set("Content-Type", contenttype)
	h.Set("Content-Encoding", "gzip")

	return w.CreatePart(h)
}

func fileUploadRequest(uri, method string, params url.Values, fileParamName, fileName, mimeType string, data io.Reader) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := createGZIPFormFile(writer, fileParamName, fileName, mimeType)
	if err != nil {
		return nil, err
	}

	gzipWriter := gzip.NewWriter(part)
	if _, err = io.Copy(gzipWriter, data); err != nil {
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	for key, valueList := range params {
		for _, value := range valueList {
			_ = writer.WriteField(key, value)
		}
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func (c *Client) addDefaultHeaders(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+c.JWT)
	c.setUserAgent(request)
}

func (c *Client) setUserAgent(request *http.Request) {
	request.Header.Set("User-Agent", c.UserAgent)
}

func (c *Client) doRequestRaw(request *http.Request) (*http.Response, error) {
	c.addDefaultHeaders(request)

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// LookupAndFetchResource tries to download a given resource from the API
func (c *Client) LookupAndFetchResource(resourceType, input string) (bool, []byte, error) {
	return c.FetchResource("/lookup?type=" + resourceType + "&q=" + input)
}

// FetchResource tries to download a given resource from the API
func (c *Client) FetchResource(path string) (bool, []byte, error) {
	return c.fetch(path)
}

func (c *Client) fetch(path string) (bool, []byte, error) {
	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return false, nil, err
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		return false, nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, nil, err
	}

	if response.StatusCode != 200 {
		return false, body, nil
	}

	return true, body, nil
}

func close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
