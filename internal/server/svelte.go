package server

import (
	"context"
	"fmt"

	"github.com/yorubanashi/languages/internal/svelte"
)

type SvelteWalkRequest struct {
	Lang string `json:"lang"`
}
type SvelteWalkResponse struct {
	Dir *svelte.Dir `json:"dir"`
}

func (s *Server) svelteWalk(_ context.Context, req *SvelteWalkRequest) (*SvelteWalkResponse, error) {
	key := fmt.Sprintf("svelte+%s", req.Lang)
	if val, ok := s.cache[key]; ok {
		dir, ok := val.(*svelte.Dir)
		if ok {
			return &SvelteWalkResponse{Dir: dir}, nil
		}
	}

	dir, err := svelte.Walk(s.config.Svelte.Pages, req.Lang)
	if err == nil {
		s.cache[key] = dir
	}

	return &SvelteWalkResponse{Dir: dir}, err
}
