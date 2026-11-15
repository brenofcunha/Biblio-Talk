package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"lume_api/internal/models"
	"lume_api/internal/repositories"
)

type BookService struct {
	books   *repositories.BookRepository
	http    *http.Client
	openURL string
	gbURL   string
	gbKey   string
}

func NewBookService(books *repositories.BookRepository, openURL, gbURL, gbKey string) *BookService {
	return &BookService{books: books, http: &http.Client{Timeout: 6 * time.Second}, openURL: openURL, gbURL: gbURL, gbKey: gbKey}
}

func (s *BookService) Search(ctx context.Context, query string, limit int) ([]models.Book, error) {
	local, err := s.books.SearchLocal(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	if len(local) > 0 {
		return local, nil
	}
	external, err := s.searchOpenLibrary(ctx, query, limit)
	if err != nil {
		external, err = s.searchGoogleBooks(ctx, query, limit)
		if err != nil {
			return nil, err
		}
	}
	persisted := make([]models.Book, 0, len(external))
	for _, b := range external {
		saved, err := s.books.UpsertByExternalID(ctx, b)
		if err != nil {
			continue
		}
		persisted = append(persisted, *saved)
	}
	return persisted, nil
}

func (s *BookService) searchOpenLibrary(ctx context.Context, query string, limit int) ([]models.Book, error) {
	u := fmt.Sprintf("%s/search.json?q=%s&limit=%d", strings.TrimRight(s.openURL, "/"), url.QueryEscape(query), limit)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("open library indisponível")
	}
	var payload struct {
		Docs []struct {
			Key          string   `json:"key"`
			Title        string   `json:"title"`
			AuthorName   []string `json:"author_name"`
			ISBN         []string `json:"isbn"`
			Publisher    []string `json:"publisher"`
			FirstPublish int      `json:"first_publish_year"`
			CoverID      int      `json:"cover_i"`
		} `json:"docs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	out := []models.Book{}
	for _, d := range payload.Docs {
		var isbn *string
		if len(d.ISBN) > 0 {
			isbn = &d.ISBN[0]
		}
		var publisher *string
		if len(d.Publisher) > 0 {
			publisher = &d.Publisher[0]
		}
		var year *int
		if d.FirstPublish > 0 {
			y := d.FirstPublish
			year = &y
		}
		var cover *string
		if d.CoverID > 0 {
			cu := fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-L.jpg", d.CoverID)
			cover = &cu
		}
		ext := d.Key
		out = append(out, models.Book{ExternalAPIID: &ext, Title: d.Title, Authors: d.AuthorName, ISBN: isbn, Publisher: publisher, PublishYear: year, CoverURL: cover, Genres: []string{}})
	}
	return out, nil
}

func (s *BookService) searchGoogleBooks(ctx context.Context, query string, limit int) ([]models.Book, error) {
	if s.gbKey == "" {
		return nil, fmt.Errorf("google books sem chave")
	}
	u := fmt.Sprintf("%s/volumes?q=%s&maxResults=%d&key=%s", strings.TrimRight(s.gbURL, "/"), url.QueryEscape(query), limit, s.gbKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("google books indisponível")
	}
	var payload struct {
		Items []struct {
			ID         string `json:"id"`
			VolumeInfo struct {
				Title               string   `json:"title"`
				Authors             []string `json:"authors"`
				Publisher           string   `json:"publisher"`
				PublishedDate       string   `json:"publishedDate"`
				PageCount           int      `json:"pageCount"`
				Categories          []string `json:"categories"`
				IndustryIdentifiers []struct {
					Identifier string `json:"identifier"`
				} `json:"industryIdentifiers"`
				ImageLinks struct {
					Thumbnail string `json:"thumbnail"`
				} `json:"imageLinks"`
			} `json:"volumeInfo"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	out := []models.Book{}
	for _, item := range payload.Items {
		vi := item.VolumeInfo
		var year *int
		if len(vi.PublishedDate) >= 4 {
			y, _ := strconv.Atoi(vi.PublishedDate[:4])
			if y > 0 {
				year = &y
			}
		}
		var pageCount *int
		if vi.PageCount > 0 {
			p := vi.PageCount
			pageCount = &p
		}
		var isbn *string
		if len(vi.IndustryIdentifiers) > 0 {
			isbn = &vi.IndustryIdentifiers[0].Identifier
		}
		var publisher *string
		if vi.Publisher != "" {
			publisher = &vi.Publisher
		}
		var cover *string
		if vi.ImageLinks.Thumbnail != "" {
			cover = &vi.ImageLinks.Thumbnail
		}
		ext := item.ID
		out = append(out, models.Book{ExternalAPIID: &ext, Title: vi.Title, Authors: vi.Authors, ISBN: isbn, Publisher: publisher, PublishYear: year, PageCount: pageCount, CoverURL: cover, Genres: vi.Categories})
	}
	return out, nil
}

func (s *BookService) GetDetail(ctx context.Context, id string) (map[string]interface{}, error) {
	return s.books.GetBookDetail(ctx, id)
}

func (s *BookService) Popular(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	return s.books.PopularWeek(ctx, limit)
}
