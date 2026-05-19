package moderation

import (
	"strings"
	"sync"
)

type Filter struct {
	mu           sync.RWMutex
	spamPatterns []string
	bannedWords  []string
	maxLength    int
}

func NewFilter() *Filter {
	return &Filter{
		spamPatterns: []string{
			"http://", "https://", "www.",
			"@everyone", "@here",
			"buy now", "click here", "free money",
		},
		bannedWords: []string{},
		maxLength:   500,
	}
}

func (f *Filter) AddSpamPattern(pattern string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.spamPatterns = append(f.spamPatterns, strings.ToLower(pattern))
}

func (f *Filter) AddBannedWord(word string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.bannedWords = append(f.bannedWords, strings.ToLower(word))
}

func (f *Filter) IsSpam(content string) bool {
	lower := strings.ToLower(content)
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, pattern := range f.spamPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

func (f *Filter) ContainsBannedWords(content string) (bool, string) {
	lower := strings.ToLower(content)
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, word := range f.bannedWords {
		if strings.Contains(lower, word) {
			return true, word
		}
	}
	return false, ""
}

func (f *Filter) ValidateContent(content string) (bool, string) {
	if len(content) > f.maxLength {
		return false, "message_too_long"
	}
	if strings.TrimSpace(content) == "" {
		return false, "empty_message"
	}
	if f.IsSpam(content) {
		return false, "spam_detected"
	}
	if banned, word := f.ContainsBannedWords(content); banned {
		return false, "banned_word:" + word
	}
	return true, ""
}
