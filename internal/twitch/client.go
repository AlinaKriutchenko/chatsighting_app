package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	HTTPClient   *http.Client
}

func NewClient(clientID, clientSecret string) (*Client, error) {
	c := &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{},
	}
	if err := c.refreshToken(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) refreshToken() error {
	url := fmt.Sprintf(
		"https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials",
		c.ClientID, c.ClientSecret,
	)
	resp, err := c.HTTPClient.Post(url, "application/json", strings.NewReader(""))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	c.AccessToken = data["access_token"].(string)
	return nil
}

func (c *Client) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-Id", c.ClientID)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
