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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	collection   *mongo.Collection
	emailService *EmailService
	secretConfig *config.SecretConfig
}

func NewAuthService(collection *mongo.Collection, emailService *EmailService, secretConfig *config.SecretConfig) *AuthService {
	return &AuthService{collection, emailService, secretConfig}
}

func (srv *AuthService) SignUp(ctx context.Context, req dto.SignUpRequest) (string, error) {
	var existingUser model.User
	err := srv.collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return "", core.BadRequest("user with this email already exists")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return "", core.InternalServerError(err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return "", core.InternalServerError(err.Error())
	}

	user, err := model.NewUser(req.Email, req.Nickname, string(hash), false, time.Now(), nil, nil)
	if err != nil {
		return "", err
	}

	res, err := srv.collection.InsertOne(ctx, user)
	if err != nil {
		return "", core.InternalServerError(err.Error())
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", core.InternalServerError("failed to get inserted ID")
	}

	token, err := core.GenerateJwtToken(oid.Hex(), srv.secretConfig.JwtSecret)
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
		err := srv.emailService.SendHTML("Confirm Your Registration	", body.String(), req.Email)
		if err != nil {
			log.Printf("Background email task failed for %s: %v", req.Email, err)
		}
	}()

	return oid.Hex(), nil
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

	objID, err := primitive.ObjectIDFromHex(sub)
	if err != nil {
		return core.BadRequest("invalid user id")
	}

	res, err := srv.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID, "is_verified": false},
		bson.M{"$set": bson.M{"is_verified": true}},
	)
	if err != nil {
		return core.InternalServerError(err.Error())
	}

	if res.MatchedCount == 0 {
		return core.NotFound("user not found or already verified")
	}
	return nil
}
