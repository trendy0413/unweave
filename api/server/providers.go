package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/unweave/unweave/api/types"
)

type ProviderService struct {
	srv *Service
}

func (p *ProviderService) ListNodeTypes(ctx context.Context, provider types.Provider, filterAvailable bool) ([]types.NodeType, error) {
	if provider != types.LambdaLabsProvider && provider != types.UnweaveProvider {
		return nil, &types.Error{
			Code:       http.StatusBadRequest,
			Message:    "Invalid runtime provider: " + string(provider),
			Suggestion: fmt.Sprintf("Use %q or %q as the runtime provider", types.LambdaLabsProvider, types.UnweaveProvider),
		}
	}

	rt, err := p.srv.InitializeRuntime(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime: %w", err)
	}

	nodeTypes, err := rt.Node.ListNodeTypes(ctx, filterAvailable)
	if err != nil {
		return nil, fmt.Errorf("failed to list node types: %w", err)
	}

	return nodeTypes, nil
}
