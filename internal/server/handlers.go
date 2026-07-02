package server

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/AlinaKriutchenko/chatsighting/internal/analytics"
)

var durationRe = regexp.MustCompile(`(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?`)

func parseDuration(s string) int {
	m := durationRe.FindStringSubmatch(s)
	if m == nil {
		return 0
	}
	h, _ := strconv.Atoi(m[1])
	min, _ := strconv.Atoi(m[2])
	sec, _ := strconv.Atoi(m[3])
	return h*3600 + min*60 + sec
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleVODs(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}

	user, err := s.twitch.GetUser(username)
	if err != nil {
		http.Error(w, "Streamer \""+username+"\" not found on Twitch. Check the spelling and try again.", http.StatusNotFound)
		return
	}

	vods, err := s.twitch.GetVODs(user.ID, 5)
	if err != nil {
		http.Error(w, user.DisplayName+" doesn't have any saved VODs. They may not have streamed recently, or their past broadcasts aren't saved.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
		"vods": vods,
	})
}

func spikeCount(r *http.Request) int {
	n, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil || n < 1 || n > 20 {
		return 5
	}
	return n
}

func (s *Server) handleAnalysis(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	vodID := r.URL.Query().Get("vod_id")
	count := spikeCount(r)

	if username == "" && vodID == "" {
		http.Error(w, "username or vod_id required", http.StatusBadRequest)
		return
	}

	if vodID != "" {
		v, err := s.twitch.GetVODByID(vodID)
		if err != nil {
			http.Error(w, "VOD not found.", http.StatusNotFound)
			return
		}
		messages, err := s.twitch.GetAllChatMessages(v.ID, parseDuration(v.Duration))
		if err != nil {
			http.Error(w, "Failed to fetch chat: "+err.Error(), http.StatusInternalServerError)
			return
		}
		windows := analytics.BuildWindows(messages)
		spikes := analytics.FindTopSpikes(windows, count)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"vod":            v,
			"spikes":         spikes,
			"stream_average": analytics.StreamAverage(windows),
			"windows":        analytics.SlimWindows(windows),
		})
		return
	}

	user, err := s.twitch.GetUser(username)
	if err != nil {
		http.Error(w, "Streamer \""+username+"\" not found on Twitch. Check the spelling and try again.", http.StatusNotFound)
		return
	}

	v, err := s.twitch.GetLatestVOD(user.ID)
	if err != nil {
		http.Error(w, user.DisplayName+" doesn't have any saved VODs. They may not have streamed recently, or their past broadcasts aren't saved.", http.StatusNotFound)
		return
	}

	messages, err := s.twitch.GetAllChatMessages(v.ID, parseDuration(v.Duration))
	if err != nil {
		http.Error(w, "Failed to fetch chat: "+err.Error(), http.StatusInternalServerError)
		return
	}

	windows := analytics.BuildWindows(messages)
	spikes := analytics.FindTopSpikes(windows, count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":           user,
		"vod":            v,
		"spikes":         spikes,
		"stream_average": analytics.StreamAverage(windows),
		"windows":        analytics.SlimWindows(windows),
	})
}
