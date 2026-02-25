package pokeapi

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokecache"
)

type mockTransport struct {
	requestCount int
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.requestCount++
	body := `{"count":1,"next":null,"previous":null,"results":[{"name":"canalave-city-area","url":"https://pokeapi.co/api/v2/location-area/1/"}]}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func TestGetLocationAreasCachesResponse(t *testing.T) {
	transport := &mockTransport{}
	httpClient := &http.Client{Transport: transport}
	cache := pokecache.NewCache(5 * time.Minute)
	client := NewClientWithHTTPClient(httpClient, cache)

	// First call (map) — should hit the API
	_, err := client.GetLocationAreas(0)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if transport.requestCount != 1 {
		t.Fatalf("expected 1 HTTP request after first call, got %d", transport.requestCount)
	}

	// Second call at same offset (mapb going back) — should hit cache
	_, err = client.GetLocationAreas(0)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if transport.requestCount != 1 {
		t.Fatalf("expected 1 HTTP request after second call (cache hit), got %d", transport.requestCount)
	}
}
