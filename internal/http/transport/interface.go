package transport

import (
	"context"

	"github.com/Oralkhan-coder/mind-map/internal/http/dto"
	"github.com/Oralkhan-coder/mind-map/internal/model"
)

type AuthService interface {
	SignUp(ctx context.Context, req dto.SignUpRequest) (string, error)
	ConfirmEmail(ctx context.Context, token string) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error)
}

type MapService interface {
	GetMaps(ctx context.Context, userId string) ([]*model.Map, error)
	GetByID(ctx context.Context, mapId, userId string) (*model.Map, error)
	CreateMap(ctx context.Context, req *dto.MapCURequest, userId string) (string, error)
	UpdateMap(ctx context.Context, mapId, userId string, update *dto.MapCURequest) error
	DeleteMap(ctx context.Context, mapId, userId string) error
}

type NodeService interface {
	GetByMapId(ctx context.Context, mapId, userId string) ([]*model.Node, error)
	CreateNode(ctx context.Context, req *dto.NodeCreateRequest, userId, mapId string) (string, error)
}
