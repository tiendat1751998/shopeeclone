package query

type TokenType string

const (
	TokenWord   TokenType = "word"
	TokenCJK    TokenType = "cjk"
	TokenNumber TokenType = "number"
	TokenPunct  TokenType = "punctuation"
)

type Token struct {
	Text  string    `json:"text"`
	Type  TokenType `json:"type"`
	Start int       `json:"start"`
	End   int       `json:"end"`
}

type ParsedQuery struct {
	Original   string       `json:"original"`
	Tokens     []Token      `json:"tokens"`
	Normalized string       `json:"normalized"`
	Correction *Correction  `json:"correction,omitempty"`
}

type Correction struct {
	Original string   `json:"original"`
	Suggestions []string `json:"suggestions"`
	Distance   int      `json:"distance"`
}
