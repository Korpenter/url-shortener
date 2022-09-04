package handlers

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
	"io"
	"log"
	"net/http"
	"strings"
)

type ShortenerHandler struct {
	urls storage.Repositories
	cfg  *config.Config
}

func NewShortenerHandler(r storage.Repositories, c *config.Config) *ShortenerHandler {
	return &ShortenerHandler{
		urls: r,
		cfg:  c,
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
	l, err := h.urls.Get(id[0])
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
	if !validators.IsURL(long) {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	short := h.urls.Add(long)
	w.Header().Set("Content-Type", "text/plain;")
	w.WriteHeader(http.StatusCreated)
	if _, err = io.WriteString(w, h.cfg.Prefix+short); err != nil {
		log.Println(err)
	}
}
