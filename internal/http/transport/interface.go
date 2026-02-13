package transport

import (
	"context"

	"github.com/Oralkhan-coder/mind-map/internal/http/dto"
)

type AuthService interface {
	SignUp(ctx context.Context, req dto.SignUpRequest) (string, error)
	ConfirmEmail(ctx context.Context, token string) error
}
