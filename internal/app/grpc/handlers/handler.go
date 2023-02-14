package handler

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Mldlr/url-shortener/internal/app/grpc/proto"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenerHandler is an gRPC server shortener handler.
type ShortenerHandler struct {
	pb.UnimplementedShortenerServer
	shortener service.ShortenerService
}

// NewShortenerHandler creates new ShortenerHandler instance.
func NewShortenerHandler(service service.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{
		shortener: service,
	}
}

// Shorten shortens the url
func (h *ShortenerHandler) Shorten(ctx context.Context, in *pb.ShortenURLRequest) (*pb.ShortenURLResponse, error) {
	userID, ok := helpers.CheckMDValue(ctx, "user_id")
	if !ok {
		return nil, status.Error(codes.Internal, "error getting user cookie")
	}
	var exitStatus codes.Code
	url, err := h.shortener.Shorten(ctx, &models.URL{LongURL: in.OriginalURL, UserID: userID})
	if err != nil {
		// If there is an error, and its not a duplicate url
		if errors.Is(err, models.ErrInvalidURL) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else if !errors.Is(err, models.ErrDuplicate) {
			fmt.Println(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
		// If it is a duplicate url
		exitStatus = codes.AlreadyExists
	} else {
		// If URL is new, return a created status.
		exitStatus = codes.OK
	}
	var resp pb.ShortenURLResponse
	resp.ShortURL = url.ShortURL
	return &resp, status.Error(exitStatus, "")
}

// Expand return original url for short.
func (h *ShortenerHandler) Expand(ctx context.Context, in *pb.ExpandURLRequest) (*pb.ExpandURLResponse, error) {
	var resp pb.ExpandURLResponse
	url, err := h.shortener.Expand(ctx, in.ShortURL)
	if err != nil {
		// If the URL has been deleted, return Gone status.
		if !errors.Is(err, models.ErrURLDeleted) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	resp.OriginalURL = url.LongURL
	return &resp, nil
}

// Ping checks the status of shortener service.
func (h *ShortenerHandler) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	// Ping the repository, cancel the operation if request is cancelled.
	err := h.shortener.Ping(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, status.Error(codes.OK, "")
}

// ShortenBatch shortens a batch of URLs.
func (h *ShortenerHandler) ShortenBatch(ctx context.Context, in *pb.BatchLinksRequest) (*pb.BatchLinksResponse, error) {
	userID, ok := helpers.CheckMDValue(ctx, "user_id")
	if !ok {
		return nil, status.Error(codes.Internal, "error getting user cookie")
	}
	urls := make([]*models.URL, len(in.BatchLinkRequestItem))
	for i, v := range in.BatchLinkRequestItem {
		// Create a new URL model and add it to the URLs map.
		urls[i] = &models.URL{
			LongURL: v.OriginalURL,
			UserID:  userID,
		}
	}
	var statusCode codes.Code
	shortenedURLs, err := h.shortener.ShortenBatch(ctx, userID, urls)
	if err != nil {
		// If there is an error, and its not a duplicate url
		if !errors.Is(err, models.ErrDuplicate) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		// If it is a duplicate url
		statusCode = codes.AlreadyExists
	} else {
		// If URLs are new, return a created status.
		statusCode = codes.OK
	}
	resp := &pb.BatchLinksResponse{BatchLinkResponseItem: make([]*pb.BatchResponseItem, len(in.BatchLinkRequestItem))}
	for i, v := range in.BatchLinkRequestItem {
		// Create response item
		resp.BatchLinkResponseItem[i] = &pb.BatchResponseItem{
			CorrelationId: v.CorrelationId,
			ShortURL:      shortenedURLs[i].ShortURL,
		}
	}
	return resp, status.Error(statusCode, "")
}

// InternalStats returns the amount of registered users and stored urls
func (h *ShortenerHandler) InternalStats(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	stats, err := h.shortener.Stats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var resp pb.StatsResponse
	resp.UrlCount = int32(stats.URLCount)
	resp.UserCount = int32(stats.UserCount)
	return &resp, nil
}

// APIUserExpand retrieves the list of shortened URLs
// created by a user and returns them as a JSON array.
func (h *ShortenerHandler) UserExpand(ctx context.Context, in *pb.UserURLRequest) (*pb.UserURLResponse, error) {
	// Get the user ID from the request context.
	userID, ok := helpers.CheckMDValue(ctx, "user_id")
	if !ok {
		return nil, status.Error(codes.Internal, "error getting user cookie")
	}
	// Get the list of URLs created by user.
	urls, err := h.shortener.ExpandUser(ctx, userID)
	if err != nil {
		if !errors.Is(err, models.ErrNoContent) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.NotFound, err.Error())
	}
	resp := &pb.UserURLResponse{Urls: make([]*pb.UserLink, len(urls))}
	// Build the response
	for i, v := range urls {
		resp.Urls[i] = &pb.UserLink{
			ShortURL:    v.ShortURL,
			OriginalURL: v.LongURL,
		}
	}
	return resp, nil
}

// APIDeleteBatch processes a batch request to delete multiple shortened URLs.
func (h *ShortenerHandler) DeleteBatch(ctx context.Context, in *pb.DeleteURLRequest) (*pb.DeleteURLResponse, error) {
	// Get the user ID from the request context.
	userID, ok := helpers.CheckMDValue(ctx, "user_id")
	if !ok {
		return nil, status.Error(codes.Internal, "error getting user cookie")
	}
	if len(in.Urls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	go h.shortener.DeleteBatch(in.Urls, userID)
	var resp pb.DeleteURLResponse
	// Return an accepted status to indicate that the request has been
	// received and is being processed.
	return &resp, nil
}
