package opsman

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"

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

func NewClient(baseUrl string, username, secret, client, client_secret string) Client {
	return Client{baseUrl, username, secret, client, client_secret}
}

func (c Client) execute(method string, endpoint string, data string, timeout time.Duration) ([]byte, error) {
	t := timeout
	if t == 0 {
		t = 30 * time.Second
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
	)
	if err != nil {
		return []byte(""), err
	}

	stdout := new(bytes.Buffer)
	stdoutLogger := log.New(stdout, "", 0)
	devNull := ioutil.Discard
	stderrLogger := log.New(devNull, "", 0)
	requestService := api.NewRequestService(oAuthClient)

	curlCommand := commands.NewCurl(requestService, stdoutLogger, stderrLogger)
	switch method {
	case "GET":
		err = curlCommand.Execute([]string{"-path", endpoint})
	case "POST":
		err = curlCommand.Execute([]string{"-path", endpoint, "-x", "POST", "-d", data})
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
