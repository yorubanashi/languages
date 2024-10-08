package server

import (
	"context"

	"github.com/yorubanashi/languages/internal/svelte"
)

type SvelteWalkRequest struct {
	Lang string `json:"lang"`
}
type SvelteWalkResponse struct {
}

func (s *Server) svelteWalk(ctx context.Context, req *SvelteWalkRequest) (*SvelteWalkResponse, error) {
	svelte.Walk(s.config.Svelte.Pages, req.Lang)
	return nil, nil
}
func (s *Server) svelteWalkHandler(ctx context.Context, decode func(interface{}) error) (interface{}, error) {
	in := &SvelteWalkRequest{}
	err := decode(in)
	if err != nil {
		return nil, err
	}
	return s.svelteWalk(ctx, in)
}
