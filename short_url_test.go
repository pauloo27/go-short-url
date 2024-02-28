package goshorturl_test

import (
	goshorturl "go-short-url"
	"os"
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	urlShortener  *goshorturl.URLShortener
	aliasAlphabet = []rune("_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func TestMain(m *testing.M) {
	urlShortener = goshorturl.NewURLShortener(os.Getenv("SHORT_URL_BASE_URL"))
	os.Exit(m.Run())
}

func TestURLShortener_ShortenWithAlias(t *testing.T) {
	alias, err := gonanoid.Generate(string(aliasAlphabet), 8)
	if err != nil {
		t.Fatal(err)
	}

	url := "https://example.com"

	shortURL, err := urlShortener.ShortenWithAlias(url, alias)
	if err != nil {
		t.Fatal(err)
	}

	if shortURL.Alias != alias {
		t.Errorf("expected alias %s, got %s", alias, shortURL.Alias)
	}

	if shortURL.OriginalURL != url {
		t.Errorf("expected original URL %s, got %s", url, shortURL.OriginalURL)
	}
}

func TestURLShortener_Shorten(t *testing.T) {
	url := "https://example.com"

	shortURL, err := urlShortener.Shorten(url)
	if err != nil {
		t.Fatal(err)
	}

	if shortURL.Alias == "" {
		t.Errorf("expected non-empty alias, got %s", shortURL.Alias)
	}

	if shortURL.OriginalURL != url {
		t.Errorf("expected original URL %s, got %s", url, shortURL.OriginalURL)
	}
}

func TestURLShortener_GetURL(t *testing.T) {
	url := "https://example.com"

	shortURL, err := urlShortener.Shorten(url)
	if err != nil {
		t.Fatal(err)
	}

	redirectURL, err := urlShortener.GetURL(shortURL.Alias)
	if err != nil {
		t.Fatal(err)
	}

	if redirectURL != url {
		t.Errorf("expected redirect URL %s, got %s", url, redirectURL)
	}
}

func TestURLShortener_GetMostAccessedURLs(t *testing.T) {
	urls, err := urlShortener.GetMostAccessedURLs(5)
	if err != nil {
		t.Fatal(err)
	}

	if urls.Limit != 5 {
		t.Errorf("expected limit 5, got %d", urls.Limit)
	}

	if len(urls.Item) > 5 {
		t.Errorf("expected at most 5 items, got %d", len(urls.Item))
	}
}
