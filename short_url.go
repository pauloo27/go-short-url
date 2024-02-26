package goshorturl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	httpClient *http.Client
)

type URLShortener struct {
	baseURL string
}

type ShortURL struct {
	Alias       string `json:"alias"`
	OriginalURL string `json:"original_url"`
	AccessCount int    `json:"access_count"`
}

type ShortURLList struct {
	Item  []*ShortURL `json:"items"`
	Limit int         `json:"limit"`
	Count int         `json:"count"`
}

func NewURLShortener(baseURL string) *URLShortener {
	return &URLShortener{baseURL: baseURL}
}

func (s *URLShortener) Shorten(url string) (*ShortURL, error) {
	return s.doShort(url, nil)
}

func (s *URLShortener) ShortenWithAlias(url, alias string) (*ShortURL, error) {
	return s.doShort(url, &alias)
}

func (s *URLShortener) doShort(url string, alias *string) (*ShortURL, error) {
	uri := fmt.Sprintf("%s/urls", s.baseURL)
	body := map[string]any{
		"original_url": url,
	}

	if alias != nil {
		body["alias"] = alias
	}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(uri, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var shortURL ShortURL
	err = json.NewDecoder(res.Body).Decode(&shortURL)
	if err != nil {
		return nil, err
	}

	return &shortURL, nil
}

func (s *URLShortener) GetURL(alias string) (url string, err error) {
	// the default client will follow redirects
	if httpClient == nil {
		httpClient = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	uri := fmt.Sprintf("%s/urls/%s", s.baseURL, alias)
	res, err := httpClient.Get(uri)
	if err != nil {
		return "", err
	}

	if res.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("alias no found")
	}

	if res.StatusCode != http.StatusFound {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res.Header.Get("Location"), nil
}

func (s *URLShortener) GetMostAccessedURLs(limit int) (*ShortURLList, error) {
	uri := fmt.Sprintf("%s/urls?limit=%d", s.baseURL, limit)
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	var shortURLList ShortURLList

	err = json.NewDecoder(res.Body).Decode(&shortURLList)
	if err != nil {
		return nil, err
	}

	return &shortURLList, nil
}
