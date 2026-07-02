package analytics

import (
	"sort"
	"strings"
	"unicode"

	"github.com/AlinaKriutchenko/chatsighting/internal/twitch"
)

var stopWords = map[string]bool{
	// English filler
	"the": true, "a": true, "an": true, "is": true, "in": true, "it": true,
	"of": true, "to": true, "and": true, "i": true, "that": true, "im": true,
	"you": true, "he": true, "she": true, "we": true, "they": true, "ur": true,
	"me": true, "my": true, "so": true, "do": true, "for": true, "on": true,
	"at": true, "be": true, "or": true, "no": true, "not": true, "up": true,
	"was": true, "are": true, "but": true, "have": true, "with": true, "as": true,
	"by": true, "from": true, "this": true, "its": true, "just": true,
	// common Twitch emotes that add no info
	"play": true, "lol": true, "lmao": true, "lul": true, "kekw": true,
	"omegalul": true, "pog": true, "pogchamp": true, "monkas": true,
	"pepehands": true, "peepohappy": true, "sadge": true, "copium": true,
	"ez": true, "gg": true,
}

// isUsableWord returns false for invisible chars, pure emoji, or single-rune words.
func isUsableWord(w string) bool {
	if len([]rune(w)) < 2 {
		return false
	}
	for _, r := range w {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func TopWords(messages []twitch.ChatMessage, n int) []string {
	counts := map[string]int{}
	for _, msg := range messages {
		words := strings.Fields(strings.ToLower(msg.Message.Body))
		for _, w := range words {
			w = strings.Trim(w, ".,!?;:'\"@#")
			if !isUsableWord(w) || stopWords[w] {
				continue
			}
			counts[w]++
		}
	}
	type wc struct {
		word  string
		count int
	}
	pairs := []wc{}
	for w, c := range counts {
		pairs = append(pairs, wc{w, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	result := []string{}
	for i := 0; i < n && i < len(pairs); i++ {
		result = append(result, pairs[i].word)
	}
	return result
}
