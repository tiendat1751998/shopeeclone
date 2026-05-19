package search

import (
	"context"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Repository interface {
	Search(ctx context.Context, query SearchQuery) (*SearchResult, error)
	FacetedSearch(ctx context.Context, query SearchQuery) (*SearchResult, error)
	GetByID(ctx context.Context, id string) (*ProductDocument, error)
	Index(ctx context.Context, doc *ProductDocument) error
	BulkIndex(ctx context.Context, docs []*ProductDocument) error
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	docs map[string]*ProductDocument
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		docs: make(map[string]*ProductDocument),
	}
}

func (r *InMemoryRepository) Search(ctx context.Context, query SearchQuery) (*SearchResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	start := time.Now()

	var results []ProductDocument
	for _, doc := range r.docs {
		if matchesQuery(doc, query) {
			results = append(results, *doc)
		}
	}

	sortResults(results, query.SortBy)

	total := int64(len(results))
	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = 20
	}

	startIdx := (page - 1) * limit
	if startIdx >= len(results) {
		return &SearchResult{
			Products: []ProductDocument{},
			Total:    total,
			Page:     page,
			Limit:    limit,
			TookMs:   time.Since(start).Milliseconds(),
		}, nil
	}

	endIdx := startIdx + limit
	if endIdx > len(results) {
		endIdx = len(results)
	}

	return &SearchResult{
		Products: results[startIdx:endIdx],
		Total:    total,
		Page:     page,
		Limit:    limit,
		TookMs:   time.Since(start).Milliseconds(),
	}, nil
}

func (r *InMemoryRepository) FacetedSearch(ctx context.Context, query SearchQuery) (*SearchResult, error) {
	result, err := r.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	categoryFacet := Facet{
		Field:  "category",
		Values: make(map[string]int64),
	}
	brandFacet := Facet{
		Field:  "brand",
		Values: make(map[string]int64),
	}

	for _, doc := range r.docs {
		if matchesQuery(doc, query) {
			if doc.Category != "" {
				categoryFacet.Values[doc.Category]++
			}
			if len(doc.Tags) > 0 {
				for _, tag := range doc.Tags {
					brandFacet.Values[tag]++
				}
			}
		}
	}

	result.Facets = []Facet{categoryFacet, brandFacet}
	return result, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*ProductDocument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	doc, ok := r.docs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return doc, nil
}

func (r *InMemoryRepository) Index(ctx context.Context, doc *ProductDocument) error {
	if doc.ID == "" {
		return ErrInvalidQuery
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.docs[doc.ID] = doc
	return nil
}

func (r *InMemoryRepository) BulkIndex(ctx context.Context, docs []*ProductDocument) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, doc := range docs {
		if doc.ID != "" {
			r.docs[doc.ID] = doc
		}
	}
	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.docs, id)
	return nil
}

func matchesQuery(doc *ProductDocument, q SearchQuery) bool {
	if q.Query != "" {
		qLower := strings.ToLower(q.Query)
		terms := strings.Fields(qLower)
		docText := strings.ToLower(doc.Title + " " + doc.Description + " " + doc.Category)
		matched := false
		for _, term := range terms {
			if strings.Contains(docText, term) || levenshteinMatch(docText, term, 2) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if q.Category != "" && !strings.EqualFold(doc.Category, q.Category) {
		return false
	}

	if q.MinPrice > 0 && doc.Price < q.MinPrice {
		return false
	}
	if q.MaxPrice > 0 && doc.Price > q.MaxPrice {
		return false
	}

	if q.MinRating > 0 && doc.Rating < q.MinRating {
		return false
	}

	return true
}

func levenshteinMatch(text, term string, maxDistance int) bool {
	words := strings.Fields(text)
	for _, word := range words {
		dist := levenshteinDistance(word, term)
		if dist <= maxDistance {
			return true
		}
	}
	return false
}

func levenshteinDistance(a, b string) int {
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
		for j := 1; j <= lb; j++ {
			cost := 1
			ai, _ := utf8.DecodeRuneInString(a[i-1:])
			bj, _ := utf8.DecodeRuneInString(b[j-1:])
			if ai == bj {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}
	return matrix[la][lb]
}

func sortResults(results []ProductDocument, sortBy string) {
	switch sortBy {
	case "price_asc":
		sort.Slice(results, func(i, j int) bool {
			return results[i].Price < results[j].Price
		})
	case "price_desc":
		sort.Slice(results, func(i, j int) bool {
			return results[i].Price > results[j].Price
		})
	case "rating":
		sort.Slice(results, func(i, j int) bool {
			return results[i].Rating > results[j].Rating
		})
	case "newest":
		sort.Slice(results, func(i, j int) bool {
			return results[i].CreatedAt.After(results[j].CreatedAt)
		})
	default:
		sort.Slice(results, func(i, j int) bool {
			return results[i].Rating*0.6 + relevanceScore(results[i])*0.4 >
				results[j].Rating*0.6 + relevanceScore(results[j])*0.4
		})
	}
}

func relevanceScore(doc ProductDocument) float64 {
	hours := time.Since(doc.CreatedAt).Hours()
	recency := math.Exp(-hours / 720)
	rating := doc.Rating / 5.0
	stock := 1.0
	if doc.Stock > 0 {
		stock = math.Min(1.0, float64(doc.Stock)/1000.0)
	}
	return recency*0.3 + rating*0.5 + stock*0.2
}
