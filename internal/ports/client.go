package ports

import (
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpErrors HTTPErrors
	client     *http.Client
}

type HTTPError struct {
	Err error
}

type HTTPErrors []HTTPError

func (he HTTPError) Error() string {
	return he.Err.Error()
}

func (he HTTPErrors) Error() string {
	if len(he) == 0 {
		return ""
	}
	var es string
	es += he[0].Error()
	for i := 1; i < len(he); i++ {
		es += "\n" + he[i].Error()
	}
	return es
}

func (c *Client) CheckErrors() string {
	if len(c.httpErrors) > 0 {
		return c.httpErrors.Error()
	} else {
		return "No client errors"
	}
}

func (c *Client) AddError(err error) {
	c.httpErrors = append(c.httpErrors, err.(HTTPError))
}

func (c *Client) GetWithParam(URL string, param string, value string) (*http.Response, error) {

	u, err := url.Parse(URL)
	if err != nil {
		return nil, HTTPError{err}
	}
	v := url.Values{}
	v.Add(param, value)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, HTTPError{err}
	}
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, HTTPError{err}
	}
	return resp, nil
}

func (c *Client) GetData(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}

func NewClient() Client {
	return Client{client: http.DefaultClient}
}
