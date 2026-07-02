package twitch

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}

func (c *Client) GetUser(username string) (*User, error) {
	url := fmt.Sprintf(
		"https://api.twitch.tv/helix/users?login=%s",
		username,
	)
	data, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []User `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("user %s not found", username)
	}
	return &resp.Data[0], nil
}
