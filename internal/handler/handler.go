package handler

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/server"
	"github.com/Mldlr/url-shortener/internal/storage"
	"github.com/Mldlr/url-shortener/internal/utils/validate"
	"io"
	"net/http"
	"strings"
)

type ShortenerHandler struct {
	store *storage.InMemRepo
}

func NewShortenerHandler() *ShortenerHandler {
	return &ShortenerHandler{
		store: storage.NewInMemRepo(),
	}
}

func (h *ShortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleExpand(w, r)
	case http.MethodPost:
		h.handleShorten(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ShortenerHandler) handleExpand(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[1:]
	if len(id) < 1 {
		http.Error(w, "no id", http.StatusBadRequest)
		return
	}
	fmt.Println(id)
	l, err := h.store.Get(id[0])
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", l)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *ShortenerHandler) handleShorten(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request", http.StatusBadRequest)
		return
	}
	long := string(b)
	if !validate.IsURL(long) {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	short := h.store.Add(long)
	w.Header().Set("Content-Type", "text/plain;")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, server.BaseURL+short)
}
