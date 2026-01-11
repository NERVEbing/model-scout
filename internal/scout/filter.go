package scout

import "strings"

var defaultExcludeSubstrings = []string{
	"image",
	"tts",
	"asr",
	"mt",
	"ocr",
	"rerank",
	"embedding",
	"realtime",
	"livetranslate",
}

func ShouldSkip(id string, excludes []string) bool {
	candidate := strings.ToLower(id)
	for _, substring := range defaultExcludeSubstrings {
		if strings.Contains(candidate, substring) {
			return true
		}
	}
	for _, substring := range excludes {
		if substring == "" {
			continue
		}
		if strings.Contains(candidate, strings.ToLower(substring)) {
			return true
		}
	}
	return false
}

func FilterModels(models []string, excludes []string) []string {
	filtered := make([]string, 0, len(models))
	for _, modelID := range models {
		if ShouldSkip(modelID, excludes) {
			continue
		}
		filtered = append(filtered, modelID)
	}
	return filtered
}
