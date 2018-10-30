package opsman

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"

	"net/http"

	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/network"
	"github.com/pkg/errors"
)

type Client struct {
	baseUrl      string
	username     string
	secret       string
	clientID     string
	clientSecret string
}

var defaultRequestTimeout = 30 * time.Second
var defaultConnectTimeout = 5 * time.Second

func NewClient(baseUrl string, username, secret, client, clientSecret string) Client {
	return Client{baseUrl, username, secret, client, clientSecret}
}

func (c Client) execute(method string, endpoint string, data string, timeout time.Duration) ([]byte, error) {
	t := timeout
	if t == 0 {
		t = defaultRequestTimeout
	}
	oAuthClient, err := network.NewOAuthClient(
		c.baseUrl,
		c.username,
		c.secret,
		c.clientID,
		c.clientSecret,
		true,
		false,
		t,
		defaultConnectTimeout,
	)
	if err != nil {
		return []byte(""), err
	}

	stdout := new(bytes.Buffer)
	stdoutLogger := log.New(stdout, "", 0)
	devNull := ioutil.Discard
	stderrLogger := log.New(devNull, "", 0)
	requestService := api.New(api.ApiInput{
		Client: oAuthClient,
	})

	curlCommand := commands.NewCurl(requestService, stdoutLogger, stderrLogger)
	switch method {
	case "GET":
		err = curlCommand.Execute([]string{"-path", endpoint})
	case "POST":
		err = curlCommand.Execute([]string{"-path", endpoint, "-x", "POST", "-d", data})
	case "DELETE":
		err = curlCommand.Execute([]string{"-path", endpoint, "-x", "DELETE"})
	}
	return stdout.Bytes(), err
}

func (c Client) Get(endpoint string, timeout time.Duration) ([]byte, error) {
	body, err := c.execute("GET", endpoint, "", timeout)
	return body, errors.Wrap(err, string(body))
}

func (c Client) Post(endpoint, data string, timeout time.Duration) ([]byte, error) {
	body, err := c.execute("POST", endpoint, data, timeout)
	return body, errors.Wrap(err, string(body))
}

func (c Client) Delete(endpoint string, timeout time.Duration) error {
	_, err := c.execute("DELETE", endpoint, "", timeout)
	return err
}

func (c Client) Do(request *http.Request) (*http.Response, error) {
	client, err := network.NewOAuthClient(
		c.baseUrl,
		c.username,
		c.secret,
		c.clientID,
		c.clientSecret,
		true,
		false,
		defaultRequestTimeout,
		defaultConnectTimeout,
	)
	if err != nil {
		return nil, err
	}
	return client.Do(request)
}