package twitch

import (
	"encoding/json"
	"fmt"
)

type VOD struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Duration  string `json:"duration"`
	CreatedAt string `json:"created_at"`
	ViewCount int    `json:"view_count"`
}

func (c *Client) GetVODs(broadcasterID string, count int) ([]VOD, error) {
	url := fmt.Sprintf(
		"https://api.twitch.tv/helix/videos?user_id=%s&type=archive&first=%d",
		broadcasterID, count,
	)
	data, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []VOD `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no VODs found")
	}
	return resp.Data, nil
}

func (c *Client) GetLatestVOD(broadcasterID string) (*VOD, error) {
	vods, err := c.GetVODs(broadcasterID, 1)
	if err != nil {
		return nil, err
	}
	return &vods[0], nil
}

func (c *Client) GetVODByID(vodID string) (*VOD, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/videos?id=%s", vodID)
	data, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []VOD `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("VOD not found")
	}
	return &resp.Data[0], nil
}
