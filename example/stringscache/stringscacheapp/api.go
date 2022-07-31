package stringscacheapp

import (
	"context"

	"github.com/IlyaFloppy/mockaka/example/stringscache/api/scpb"
	"github.com/IlyaFloppy/mockaka/example/stringscache/api/spb"
)

// Reverse a string. Makes a request to StringsService unless request is cached.
func (s *Service) Reverse(ctx context.Context, req *scpb.Message) (*scpb.Message, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if val, ok := s.cache[req.Message]; ok {
		return &scpb.Message{
			Message: val,
		}, nil
	}

	res, err := s.stringsClient.Reverse(ctx, &spb.Message{
		Message: req.Message,
	})

	if err != nil {
		// assume stringcacheapp should return errors as is and not cache them.
		return nil, err //nolint:wrapcheck
	}

	s.cache[req.Message] = res.Message

	return &scpb.Message{
		Message: res.Message,
	}, nil
}

// Invalidate clears cache.
func (s *Service) Invalidate(context.Context, *scpb.InvalidateRequest) (*scpb.InvalidateResponse, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.cache = make(map[string]string)

	return &scpb.InvalidateResponse{}, nil
}
