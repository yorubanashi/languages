package server

import (
	"context"

	"github.com/yorubanashi/languages/internal/svelte"
)

type SvelteWalkRequest struct {
	Lang string `json:"lang"`
}
type SvelteWalkResponse struct {
	Dir *svelte.Dir `json:"dir"`
}

func (s *Server) svelteWalk(ctx context.Context, req *SvelteWalkRequest) (*SvelteWalkResponse, error) {
	dir, err := svelte.Walk(s.config.Svelte.Pages, req.Lang)
	return &SvelteWalkResponse{Dir: dir}, err
}
func (s *Server) svelteWalkHandler(ctx context.Context, decode func(interface{}) error) (interface{}, error) {
	in := &SvelteWalkRequest{}
	err := decode(in)
	if err != nil {
		return nil, err
	}
	return s.svelteWalk(ctx, in)
}
