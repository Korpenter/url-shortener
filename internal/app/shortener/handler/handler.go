package handler

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/shortener/db"
	"github.com/Mldlr/url-shortener/internal/app/shortener/utils"
	"io"
	"net/http"
	"strings"
)

type ShortenerHandler struct {
	store *db.InMemRepo
}

func NewShortenerHandler() *ShortenerHandler {
	return &ShortenerHandler{
		store: db.NewInMemRepo(),
	}
}

func (h *ShortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ShortenerHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[1:]
	if len(id) < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	fmt.Println(id)
	l, err := h.store.GetByID(id[0])
	if err != nil {
		http.Error(w, "unknown id", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Location", l)
	fmt.Println(l)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *ShortenerHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request", http.StatusBadRequest)
		return
	}
	long := string(b)
	if !utils.IsValid(long) {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	short := utils.MD5(long)[:8]
	h.store.AddLink(long, short)
	w.Header().Set("Content-Type", "text/plain;")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, short)
}
