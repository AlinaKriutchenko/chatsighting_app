package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// gqlClientID is Twitch's own frontend client ID used by their web player.
const gqlClientID = "kimne78kx3ncx6brgo4mv6wki5h1ko"

type ChatMessage struct {
	ContentOffset float64
	Commenter     struct {
		DisplayName string
	}
	Message struct {
		Body string
	}
}

type gqlEdge struct {
	Node struct {
		ContentOffsetSeconds float64 `json:"contentOffsetSeconds"`
		Commenter            *struct {
			DisplayName string `json:"displayName"`
		} `json:"commenter"`
		Message struct {
			Fragments []struct {
				Text string `json:"text"`
			} `json:"fragments"`
		} `json:"message"`
	} `json:"node"`
}

func (c *Client) fetchCommentPage(videoID string, offsetSeconds int) ([]gqlEdge, bool, error) {
	payload := []map[string]interface{}{{
		"operationName": "VideoCommentsByOffsetOrCursor",
		"variables": map[string]interface{}{
			"videoID":              videoID,
			"contentOffsetSeconds": offsetSeconds,
		},
		"extensions": map[string]interface{}{
			"persistedQuery": map[string]interface{}{
				"version":    1,
				"sha256Hash": "b70a3591ff0f4e0313d126c6a1502d79a1c02baebb288227c582044aa76adf6a",
			},
		},
	}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, false, err
	}

	req, err := http.NewRequest("POST", "https://gql.twitch.tv/gql", bytes.NewReader(body))
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Client-Id", gqlClientID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}

	var result []struct {
		Data struct {
			Video *struct {
				Comments struct {
					Edges    []gqlEdge `json:"edges"`
					PageInfo struct {
						HasNextPage bool `json:"hasNextPage"`
					} `json:"pageInfo"`
				} `json:"comments"`
			} `json:"video"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, false, fmt.Errorf("GQL parse error: %w (body: %.500s)", err, string(data))
	}
	if len(result) == 0 {
		return nil, false, fmt.Errorf("empty GQL response")
	}
	if len(result[0].Errors) > 0 {
		return nil, false, fmt.Errorf("GQL error: %s", result[0].Errors[0].Message)
	}
	if result[0].Data.Video == nil {
		return nil, false, fmt.Errorf("video not found or chat unavailable")
	}

	comments := result[0].Data.Video.Comments
	return comments.Edges, comments.PageInfo.HasNextPage, nil
}

// GetAllChatMessages fetches all chat messages for a VOD by stepping through
// offsets rather than using cursor pagination (which Twitch blocks server-side).
// Each page returns ~59 messages; we advance to just past the last message's
// offset to get the next non-overlapping chunk.
// vodDurationSeconds is used to skip past sections where Twitch returns empty
// pages prematurely (a known Twitch API inconsistency).
func (c *Client) GetAllChatMessages(videoID string, vodDurationSeconds int) ([]ChatMessage, error) {
	const maxMessages = 50000
	const skipStep = 30 // seconds to jump when Twitch returns no progress
	var all []ChatMessage
	nextOffset := 0
	lastOffset := -1.0

	for {
		edges, _, err := c.fetchCommentPage(videoID, nextOffset)
		if err != nil {
			return nil, err
		}

		var pageLastOffset float64
		for _, e := range edges {
			offset := e.Node.ContentOffsetSeconds
			if offset <= lastOffset {
				continue
			}
			var msg ChatMessage
			msg.ContentOffset = offset
			if e.Node.Commenter != nil {
				msg.Commenter.DisplayName = e.Node.Commenter.DisplayName
			}
			var parts []string
			for _, f := range e.Node.Message.Fragments {
				parts = append(parts, f.Text)
			}
			msg.Message.Body = strings.Join(parts, "")
			all = append(all, msg)
			pageLastOffset = offset
		}

		if pageLastOffset <= lastOffset {
			// No progress from Twitch — skip ahead if there's still VOD left.
			if nextOffset+skipStep < vodDurationSeconds {
				nextOffset += skipStep
				continue
			}
			break
		}
		lastOffset = pageLastOffset
		nextOffset = int(pageLastOffset) + 1

		if len(all) >= maxMessages {
			break
		}
	}
	return all, nil
}
