package emojis

import (
	"slices"
	"strings"
)

// Params used to specify search parameters
type Params struct {
	Include []string // Include slice of strings to include in search
	Exclude []string // Exclude slice of strings to exclude in search
	Order   int
	Hexcode string
}

// Search searches for all emojis using the provided params
func Search(params Params) (result []Emoji) {

	for _, emo := range emojis {
		if shouldExclude(emo, params.Exclude) {
			continue
		}

		if params.Order > 0 && params.Order != emo.Order {
			continue
		}

		if params.Hexcode != "" && emo.Hexcode != strings.ToUpper(params.Hexcode) {
			continue
		}

		for _, include := range params.Include {
			include = strings.ToLower(include)
			if strings.Contains(emo.Label, include) || slices.Contains(emo.Tags, include) {
				result = append(result, emo)
			}
		}
	}

	return
}

// shouldExclude checks emoji tags and labels for exclusions
func shouldExclude(emo Emoji, excludes []string) bool {
	for _, exclude := range excludes {
		exclude = strings.ToLower(exclude)
		if strings.Contains(emo.Label, exclude) || slices.Contains(emo.Tags, exclude) {
			return true
		}
	}

	return false
}
