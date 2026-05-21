package pagination

type Request struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"-"`
}

type Response struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

func NewRequest(page, limit int) Request {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return Request{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func NewResponse(items interface{}, total int64, req Request) Response {
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}
	return Response{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}
}

type CursorRequest struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

type CursorResponse struct {
	Items     []interface{} `json:"items"`
	NextCursor string       `json:"next_cursor,omitempty"`
	HasMore   bool          `json:"has_more"`
}
