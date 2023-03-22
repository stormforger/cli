package api

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/stormforger/cli/buildinfo"
)

func defaultHTTPClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
		Proxy:           http.ProxyURL(getHTTPProxy()),
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
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
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

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
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

func (c *Client) doRequest(request *http.Request) (bool, []byte, error) {
	response, err := c.doRequestRaw(request)
	if err != nil {
		return false, nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, nil, err
	}

	if response.StatusCode != 200 {
		return false, body, nil
	}

	return true, body, nil
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

	return c.doRequest(req)
}

func (c *Client) put(path string, body []byte) (bool, []byte, error) {
	req, err := http.NewRequest("PUT", c.APIEndpoint+path, bytes.NewReader(body))
	if err != nil {
		return false, nil, err
	}

	return c.doRequest(req)
}

func close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
