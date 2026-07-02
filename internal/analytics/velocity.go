package analytics

import (
	"math"

	"github.com/AlinaKriutchenko/chatsighting/internal/twitch"
)

const WindowSeconds = 30

type Window struct {
	StartSeconds float64
	EndSeconds   float64
	MessageCount int
	Messages     []twitch.ChatMessage
}

func BuildWindows(messages []twitch.ChatMessage) []Window {
	if len(messages) == 0 {
		return nil
	}
	windows := []Window{}
	current := Window{
		StartSeconds: math.Floor(messages[0].ContentOffset/WindowSeconds) * WindowSeconds,
	}
	current.EndSeconds = current.StartSeconds + WindowSeconds

	for _, msg := range messages {
		if msg.ContentOffset >= current.EndSeconds {
			windows = append(windows, current)
			current = Window{
				StartSeconds: math.Floor(msg.ContentOffset/WindowSeconds) * WindowSeconds,
			}
			current.EndSeconds = current.StartSeconds + WindowSeconds
		}
		current.Messages = append(current.Messages, msg)
		current.MessageCount++
	}
	windows = append(windows, current)
	return windows
}
