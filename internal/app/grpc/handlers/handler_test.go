package handler

import (
	"context"
	"testing"
	"time"

	"github.com/Mldlr/url-shortener/internal/app/config"
	pb "github.com/Mldlr/url-shortener/internal/app/grpc/proto"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestShorten(t *testing.T) {
	repo := storage.NewMockRepo()
	cfg := &config.Config{}
	shortener := service.NewShortenerImpl(repo, cfg)
	shortenerHandler := NewShortenerHandler(shortener)
	tests := []struct {
		name     string
		url      string
		errCode  codes.Code
		shortURL string
	}{
		{
			name:     "Post correct url",
			url:      "https://github.com",
			errCode:  codes.OK,
			shortURL: "wCg6zOhmhh1",
		},
		{
			name:     "Post Duplicate url",
			url:      "https://github.com",
			errCode:  codes.AlreadyExists,
			shortURL: "wCg6zOhmhh1",
		},
		{
			name:     "Post invalid url",
			url:      "",
			errCode:  codes.InvalidArgument,
			shortURL: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &pb.ShortenURLRequest{
				OriginalURL: tt.url,
			}
			incCtx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"user_id": "1324"}))
			rsp, err := shortenerHandler.Shorten(incCtx, reg)
			if tt.shortURL != "" {
				assert.Equal(t, tt.shortURL, rsp.ShortURL)
			}
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestExpand(t *testing.T) {
	repo := storage.NewMockRepo()
	cfg := &config.Config{}
	shortener := service.NewShortenerImpl(repo, cfg)
	shortenerHandler := NewShortenerHandler(shortener)
	tests := []struct {
		name        string
		url         string
		errCode     codes.Code
		originalURL string
	}{
		{
			name:        "Get existing url 1",
			url:         "3S93m80EGmF",
			errCode:     codes.OK,
			originalURL: "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders",
		},
		{
			name:        "Get existing url 2",
			url:         "aQqomlSbUsE",
			errCode:     codes.OK,
			originalURL: "https://yandex.ru/",
		},
		{
			name:        "Get non existent url",
			url:         "dfsfsdfsdfsdf",
			errCode:     codes.NotFound,
			originalURL: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &pb.ExpandURLRequest{
				ShortURL: tt.url,
			}
			incCtx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"user_id": "1324"}))
			rsp, err := shortenerHandler.Expand(incCtx, reg)
			if tt.originalURL != "" {
				assert.Equal(t, tt.originalURL, rsp.OriginalURL)
			}
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestPing(t *testing.T) {
	repo := storage.NewMockRepo()
	cfg := &config.Config{}
	shortener := service.NewShortenerImpl(repo, cfg)
	shortenerHandler := NewShortenerHandler(shortener)
	tests := []struct {
		name    string
		errCode codes.Code
	}{
		{
			name:    "Get ping",
			errCode: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &pb.PingRequest{}
			incCtx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"user_id": "1324"}))
			_, err := shortenerHandler.Ping(incCtx, reg)
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestShortenBatch(t *testing.T) {
	repo := storage.NewMockRepo()
	cfg := &config.Config{}
	shortener := service.NewShortenerImpl(repo, cfg)
	shortenerHandler := NewShortenerHandler(shortener)
	tests := []struct {
		name     string
		urls     []*pb.BatchRequstItem
		errCode  codes.Code
		response []*pb.BatchResponseItem
	}{
		{
			name: "Post correct url",
			urls: []*pb.BatchRequstItem{
				{
					CorrelationId: "1",
					OriginalURL:   "https://github.com/",
				},
				{
					CorrelationId: "2",
					OriginalURL:   "https://yandex.com/",
				},
			},
			errCode: codes.OK,
			response: []*pb.BatchResponseItem{
				{
					CorrelationId: "1",
					ShortURL:      "vRveliyDLz8",
				},
				{
					CorrelationId: "2",
					ShortURL:      "BlbEuA4l5GJ",
				},
			},
		},
		{
			name: "Post with already existing urls",
			urls: []*pb.BatchRequstItem{
				{
					CorrelationId: "1",
					OriginalURL:   "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders",
				},
				{
					CorrelationId: "2",
					OriginalURL:   "https://yandex.ru/",
				},
			},
			errCode: codes.AlreadyExists,
			response: []*pb.BatchResponseItem{
				{
					CorrelationId: "1",
					ShortURL:      "3S93m80EGmF",
				},
				{
					CorrelationId: "2",
					ShortURL:      "aQqomlSbUsE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.BatchLinksRequest{BatchLinkRequestItem: tt.urls}
			incCtx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"user_id": "1324"}))
			rsp, err := shortenerHandler.ShortenBatch(incCtx, req)
			assert.EqualValues(t, tt.response, rsp.BatchLinkResponseItem)
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestInternalStats(t *testing.T) {
	repo := storage.NewMockRepo()
	cfg := &config.Config{}
	shortener := service.NewShortenerImpl(repo, cfg)
	shortenerHandler := NewShortenerHandler(shortener)
	tests := []struct {
		name     string
		errCode  codes.Code
		urlCount int32
		usrCount int32
	}{
		{
			name:     "Get stats",
			errCode:  codes.OK,
			urlCount: 2,
			usrCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.StatsRequest{}
			incCtx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"user_id": "1324"}))
			rsp, err := shortenerHandler.InternalStats(incCtx, req)
			assert.Equal(t, tt.urlCount, rsp.UrlCount)
			assert.Equal(t, tt.usrCount, rsp.UserCount)
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestUserExpand(t *testing.T) {
	repo := storage.NewMockRepo()
	cfg := &config.Config{}
	shortener := service.NewShortenerImpl(repo, cfg)
	shortenerHandler := NewShortenerHandler(shortener)
	tests := []struct {
		name     string
		userID   string
		errCode  codes.Code
		response []*pb.UserLink
	}{
		{
			name:    "user with urls",
			userID:  "KS097f1lS&F",
			errCode: codes.OK,
			response: []*pb.UserLink{
				{
					ShortURL:    "3S93m80EGmF",
					OriginalURL: "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders",
				},
				{
					ShortURL:    "aQqomlSbUsE",
					OriginalURL: "https://yandex.ru/",
				},
			},
		},
		{
			name:     "user without urls",
			userID:   "asdasasdasd",
			errCode:  codes.NotFound,
			response: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &pb.UserURLRequest{}
			incCtx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"user_id": tt.userID}))
			rsp, err := shortenerHandler.UserExpand(incCtx, reg)
			if tt.response != nil {
				assert.EqualValues(t, tt.response, rsp.Urls)
			}
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestDeleteBatch(t *testing.T) {
	repo := storage.NewMockRepo()
	cfg := &config.Config{}
	shortener := service.NewShortenerImpl(repo, cfg)
	shortenerHandler := NewShortenerHandler(shortener)
	tests := []struct {
		name         string
		userID       string
		request      []string
		errCode      codes.Code
		checkErrCode codes.Code
	}{
		{
			name:         "try to delete url not from user",
			userID:       "1234",
			request:      []string{"aQqomlSbUsE"},
			errCode:      codes.OK,
			checkErrCode: codes.OK,
		},
		{
			name:         "try to delete 0 urls",
			userID:       "1324",
			request:      []string{},
			errCode:      codes.InvalidArgument,
			checkErrCode: codes.Unavailable,
		},
		{
			name:         "delete user url",
			userID:       "KS097f1lS&F",
			request:      []string{"aQqomlSbUsE"},
			errCode:      codes.OK,
			checkErrCode: codes.Unavailable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &pb.DeleteURLRequest{Urls: tt.request}
			incCtx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"user_id": tt.userID}))
			_, err := shortenerHandler.DeleteBatch(incCtx, reg)
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
			if len(tt.request) == 0 {
				return
			}
			checkReq := &pb.ExpandURLRequest{
				ShortURL: tt.request[0],
			}
			var statusCheckErr *status.Status
			time.Sleep(time.Second * 15)
			_, err = shortenerHandler.Expand(incCtx, checkReq)
			statusCheckErr, _ = status.FromError(err)
			assert.Equal(t, tt.checkErrCode.String(), statusCheckErr.Code().String())

		})
	}
}
