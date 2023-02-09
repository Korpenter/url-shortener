package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
)

func ExampleShorten() {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := storage.NewInMemRepo()
	r := NewRouter(repo, cfg)
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("github.com"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	body := w.Result().Body
	defer body.Close()
	bodyResult, _ := io.ReadAll(body)
	fmt.Println(string(bodyResult))
	// Output:
	// http://localhost:8080/aAE3t8nGJ9A
}

func ExampleAPIShorten() {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := storage.NewInMemRepo()
	r := NewRouter(repo, cfg)
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"github.com/"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	body := w.Result().Body
	defer body.Close()
	bodyResult, _ := io.ReadAll(body)
	fmt.Println(string(bodyResult))
	// Output:
	// {"result":"http://localhost:8080/GaSgGCXYQ18"}

}

func ExampleExpand() {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := storage.NewMockRepo()
	r := NewRouter(repo, cfg)
	request := httptest.NewRequest(http.MethodGet, "/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	header := w.Header()
	fmt.Println(header["Location"])
	// Output:
	// [https://github.com/Mldlr/url-shortener/internal/app/utils/encoders]
}

func ExampleAPIShortenBatch() {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := storage.NewInMemRepo()
	r := NewRouter(repo, cfg)
	payload := `[{"correlation_id":"TestCorrelationID1","original_url":"https://github.com/"},{"correlation_id":"TestCorrelationID2","original_url":"https://yandex.com/"}]`
	request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	body := w.Result().Body
	defer body.Close()
	bodyResult, _ := io.ReadAll(body)
	fmt.Println(string(bodyResult))
	// Output:
	// [{"correlation_id":"TestCorrelationID1","short_url":"http://localhost:8080/vRveliyDLz8"},{"correlation_id":"TestCorrelationID2","short_url":"http://localhost:8080/BlbEuA4l5GJ"}]
}

func ExampleAPIDeleteBatch() {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := storage.NewInMemRepo()
	r := NewRouter(repo, cfg)
	request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["c", "b"]`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	status := w.Code
	fmt.Println(status)
	// Output:
	// 202
}

func ExamplePing() {
	cfg := &config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := storage.NewInMemRepo()
	r := NewRouter(repo, cfg)
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	status := w.Code
	fmt.Println(status)
	// Output:
	// 200
}

func ExampleAPIInternalStats() {
	exampleSubnet := "192.168.1.0/24"
	examplePrefix, _ := netip.ParsePrefix(exampleSubnet)
	cfg := &config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: exampleSubnet,
		SubnetPrefix:  examplePrefix,
	}
	repo := storage.NewInMemRepo()
	r := NewRouter(repo, cfg)
	request := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
	request.Header.Set("X-Real-IP", "192.168.1.3")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	fmt.Println(w.Body)
	// Output:
	// {"urls":0,"users":0}
}
