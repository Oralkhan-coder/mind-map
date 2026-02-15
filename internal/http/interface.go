package http

import "github.com/Oralkhan-coder/mind-map/internal/http/transport"

type AuthSrv interface {
	transport.AuthService
}

type MapSrv interface {
	transport.MapService
}
