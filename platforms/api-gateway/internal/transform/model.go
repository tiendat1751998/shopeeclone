package transform

type MatchCondition struct {
	PathPattern string `json:"path_pattern"`
	Method      string `json:"method"`
}

type Action struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value,omitempty"`
}

type Rule struct {
	ID              string          `json:"id"`
	MatchCondition  MatchCondition `json:"match_condition"`
	Actions         []Action       `json:"actions"`
}

type TransformRequest struct {
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Query   map[string]string `json:"query"`
	Body    string            `json:"body"`
}

type TransformResponse struct {
	Headers map[string]string `json:"headers"`
	Query   map[string]string `json:"query"`
	Path    string            `json:"path"`
	Body    string            `json:"body"`
}
