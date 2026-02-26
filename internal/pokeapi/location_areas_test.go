package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokecache"
)

type mockTransport struct {
	body         string
	statusCode   int
	requestCount int
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.requestCount++
	return &http.Response{
		StatusCode: t.statusCode,
		Body:       io.NopCloser(strings.NewReader(t.body)),
	}, nil
}

func newTestClient(transport *mockTransport) *Client {
	httpClient := &http.Client{Transport: transport}
	cache := pokecache.NewCache(5 * time.Minute)
	return NewClientWithHTTPClient(httpClient, cache)
}

func TestGetLocationAreasCachesResponse(t *testing.T) {
	transport := &mockTransport{
		statusCode: http.StatusOK,
		body:       `{"next":null,"previous":null,"results":[{"name":"canalave-city-area"}]}`,
	}
	client := newTestClient(transport)

	areas, _, _, err := client.GetLocationAreas("")
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if len(areas) != 1 || areas[0].Name != "canalave-city-area" {
		t.Fatalf("unexpected areas: %v", areas)
	}
	if transport.requestCount != 1 {
		t.Fatalf("expected 1 HTTP request after first call, got %d", transport.requestCount)
	}

	_, _, _, err = client.GetLocationAreas("")
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if transport.requestCount != 1 {
		t.Fatalf("expected 1 HTTP request after second call (cache hit), got %d", transport.requestCount)
	}
}

func TestGetLocationAreasPagination(t *testing.T) {
	nextURL := "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"
	prevURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"

	transport := &mockTransport{
		statusCode: http.StatusOK,
		body: fmt.Sprintf(
			`{"next":%q,"previous":%q,"results":[{"name":"area-1"}]}`,
			nextURL, prevURL,
		),
	}
	client := newTestClient(transport)

	areas, next, prev, err := client.GetLocationAreas("")
	if err != nil {
		t.Fatalf("call failed: %v", err)
	}
	if len(areas) != 1 || areas[0].Name != "area-1" {
		t.Fatalf("unexpected areas: %v", areas)
	}
	if next != nextURL {
		t.Errorf("next = %q, want %q", next, nextURL)
	}
	if prev != prevURL {
		t.Errorf("prev = %q, want %q", prev, prevURL)
	}
}

func TestGetLocationAreasNullPagination(t *testing.T) {
	transport := &mockTransport{
		statusCode: http.StatusOK,
		body:       `{"next":null,"previous":null,"results":[]}`,
	}
	client := newTestClient(transport)

	areas, next, prev, err := client.GetLocationAreas("")
	if err != nil {
		t.Fatalf("call failed: %v", err)
	}
	if len(areas) != 0 {
		t.Errorf("expected empty areas, got %v", areas)
	}
	if next != "" {
		t.Errorf("next = %q, want empty string", next)
	}
	if prev != "" {
		t.Errorf("prev = %q, want empty string", prev)
	}
}

func TestGetHTTPError(t *testing.T) {
	transport := &mockTransport{
		statusCode: http.StatusNotFound,
		body:       `Not Found`,
	}
	client := newTestClient(transport)

	_, _, _, err := client.GetLocationAreas("")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "404")
	}
}
