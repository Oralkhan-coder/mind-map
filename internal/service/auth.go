package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/Oralkhan-coder/mind-map/config"
	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/Oralkhan-coder/mind-map/internal/http/dto"
	"github.com/Oralkhan-coder/mind-map/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService  *UserService
	emailService *EmailService
	secretConfig *config.SecretConfig
}

func NewAuthService(userService *UserService, emailService *EmailService, secretConfig *config.SecretConfig) *AuthService {
	return &AuthService{userService, emailService, secretConfig}
}

func (srv *AuthService) SignUp(ctx context.Context, req dto.SignUpRequest) (string, error) {
	if _, err := srv.userService.GetUserByEmail(ctx, req.Email); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return "", core.InternalServerError(err.Error())
	}

	user, err := model.NewUser(req.Email, req.Nickname, string(hash), false, time.Now(), nil, nil)
	if err != nil {
		return "", err
	}

	id, err := srv.userService.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	token, err := core.GenerateJwtToken(id, srv.secretConfig.JwtSecret)
	if err != nil {
		return "", core.InternalServerError(err.Error())
	}

	tmpl, err := template.ParseFiles("ui/confirmation.page.tmpl")
	if err != nil {
		return "", err
	}
	var body bytes.Buffer
	data := map[string]string{
		"username":     req.Nickname,
		"confirm_link": fmt.Sprintf("http://localhost:8080/confirm?token=%s", token),
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	go func() {
		fmt.Printf("http://localhost:8080/confirm?token=%s \n", token)
		err := srv.emailService.SendHTML("Confirm Your Registration	", body.String(), req.Email)
		if err != nil {
			log.Printf("Background email task failed for %s: %v", req.Email, err)
		}
	}()

	return id, nil
}

func (srv *AuthService) ConfirmEmail(ctx context.Context, token string) error {
	if token == "" {
		return core.BadRequest("empty token")
	}

	tokenStr, err := core.ValidateJwtToken(token, srv.secretConfig.JwtSecret)
	if err != nil {
		return err
	}

	if !tokenStr.Valid {
		return core.Unauthorized("invalid token")
	}

	claims, ok := tokenStr.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to parse claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return core.BadRequest("invalid subject claim")
	}

	err = srv.userService.VerifyUser(ctx, sub)
	if err != nil {
		return err
	}

	return nil
}

func (srv *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := srv.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, core.NotFound("user does not exist")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, core.BadRequest("invalid password")
	}

	token, err := core.GenerateJwtToken(user.ID.Hex(), srv.secretConfig.JwtSecret)
	if err != nil {
		return nil, core.InternalServerError(err.Error())
	}

	return &dto.TokenResponse{AccessToken: token, RefreshToken: ""}, nil
}
