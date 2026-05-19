package query

import (
	"context"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Service interface {
	Parse(ctx context.Context, raw string) (*ParsedQuery, error)
	Correct(ctx context.Context, query string, dictionary []string) (*Correction, error)
	Tokenize(text string) []Token
	Normalize(text string) string
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

var stopWords = map[string]bool{
	"a": true, "an": true, "the": true, "and": true, "or": true,
	"in": true, "on": true, "at": true, "to": true, "for": true,
	"of": true, "with": true, "is": true, "it": true, "as": true,
	"be": true, "by": true, "from": true, "that": true, "this": true,
	"are": true, "was": true, "were": true, "been": true, "being": true,
	"have": true, "has": true, "had": true, "do": true, "does": true,
	"did": true, "but": true, "not": true, "so": true, "if": true,
	"no": true, "up": true, "out": true, "about": true, "into": true,
	"over": true, "after": true, "then": true, "than": true,
}

func (s *service) Parse(ctx context.Context, raw string) (*ParsedQuery, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, ErrEmptyQuery
	}
	if len(raw) > 500 {
		return nil, ErrQueryTooLong
	}

	tokens := s.Tokenize(raw)
	normalized := s.Normalize(raw)

	return &ParsedQuery{
		Original:   raw,
		Tokens:     tokens,
		Normalized: normalized,
	}, nil
}

func (s *service) Correct(ctx context.Context, query string, dictionary []string) (*Correction, error) {
	words := strings.Fields(strings.ToLower(query))
	for _, word := range words {
		if stopWords[word] {
			continue
		}
		bestDist := 3
		bestWord := ""
		for _, dictWord := range dictionary {
			dist := levenshtein(word, dictWord)
			if dist > 0 && dist <= 2 && dist < bestDist {
				bestDist = dist
				bestWord = dictWord
			}
		}
		if bestWord != "" {
			return &Correction{
				Original:   word,
				Suggestions: []string{bestWord},
				Distance:   bestDist,
			}, nil
		}
	}
	return nil, ErrNoCorrections
}

func (s *service) Tokenize(text string) []Token {
	var tokens []Token
	pos := 0
	runes := []rune(text)

	for pos < len(runes) {
		r := runes[pos]
		start := pos

		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			for pos < len(runes) {
				r = runes[pos]
				if !unicode.IsPunct(r) && !unicode.IsSymbol(r) {
					break
				}
				pos++
			}
			tokens = append(tokens, Token{
				Text:  string(runes[start:pos]),
				Type:  TokenPunct,
				Start: start,
				End:   pos,
			})
		} else if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hiragana, r) ||
			unicode.Is(unicode.Katakana, r) || unicode.Is(unicode.Hangul, r) {
			for pos < len(runes) {
				r = runes[pos]
				if !unicode.Is(unicode.Han, r) && !unicode.Is(unicode.Hiragana, r) &&
					!unicode.Is(unicode.Katakana, r) && !unicode.Is(unicode.Hangul, r) {
					break
				}
				pos++
			}
			tokens = append(tokens, Token{
				Text:  string(runes[start:pos]),
				Type:  TokenCJK,
				Start: start,
				End:   pos,
			})
		} else if unicode.IsDigit(r) {
			for pos < len(runes) {
				r = runes[pos]
				if !unicode.IsDigit(r) {
					break
				}
				pos++
			}
			tokens = append(tokens, Token{
				Text:  string(runes[start:pos]),
				Type:  TokenNumber,
				Start: start,
				End:   pos,
			})
		} else if unicode.IsLetter(r) || r == '\'' {
			for pos < len(runes) {
				r = runes[pos]
				if !unicode.IsLetter(r) && r != '\'' {
					break
				}
				pos++
			}
			tokens = append(tokens, Token{
				Text:  string(runes[start:pos]),
				Type:  TokenWord,
				Start: start,
				End:   pos,
			})
		} else {
			pos++
		}
	}

	var filtered []Token
	for _, t := range tokens {
		if t.Type == TokenPunct {
			continue
		}
		lower := strings.ToLower(t.Text)
		if t.Type == TokenWord && stopWords[lower] {
			continue
		}
		filtered = append(filtered, t)
	}

	return filtered
}

func (s *service) Normalize(text string) string {
	var result strings.Builder
	for _, r := range text {
		if unicode.IsUpper(r) {
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func levenshtein(a, b string) int {
	la := utf8.RuneCountInString(a)
	lb := utf8.RuneCountInString(b)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	matrix := make([][]int, la+1)
	for i := range matrix {
		matrix[i] = make([]int, lb+1)
		matrix[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= la; i++ {
		ai, _ := utf8.DecodeRuneInString(a[i-1:])
		for j := 1; j <= lb; j++ {
			bj, _ := utf8.DecodeRuneInString(b[j-1:])
			cost := 1
			if ai == bj {
				cost = 0
			}
			matrix[i][j] = min3(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}
	return matrix[la][lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
