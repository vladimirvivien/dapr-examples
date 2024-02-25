package emojis

import (
	"embed"
	"encoding/json"
)

// Data (derived from https://github.com/milesj/emojibase)
//
//go:embed data.json
var filesys embed.FS

// Emoji represents an Emoji data entry
type Emoji struct {
	Unicode string
	Label   string
	Tags    []string
	Order   int
	Hexcode string
}

var emojis []Emoji

func init() {
	data, err := filesys.ReadFile("data.json")
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(data, &emojis); err != nil {
		panic(err.Error())
	}
}
