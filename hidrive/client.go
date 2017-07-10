package hidrive

import (
	"net/http"

	"github.com/dghubble/sling"
)

type Client struct {
	Users *UserService
}

func NewClient(httpClient *http.Client, baseUrl string) *Client {
	baseClient := sling.New().Client(httpClient).Base(baseUrl)

	return &Client{
		Users: newUserService(baseClient.New()),
	}
}
