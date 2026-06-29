package auth

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Rowkash/go-gin-auth/internal/common"
	"github.com/Rowkash/go-gin-auth/internal/sessions"
	"github.com/Rowkash/go-gin-auth/internal/users"
)

type Service struct {
	usersService    *users.Service
	tokenService    *TokenService
	sessionsService *sessions.Service
}

func NewService(usersService *users.Service, tokenService *TokenService, sessionsService *sessions.Service) *Service {
	return &Service{
		usersService:    usersService,
		tokenService:    tokenService,
		sessionsService: sessionsService,
	}
}

func (s *Service) register(ctx context.Context, data *RegisterDto) (*Tokens, error) {
	candidate, _ := s.usersService.FindOne(users.UserFilter{
		Email: &data.Email,
	})
	if candidate != nil {
		return nil, common.NewBadRequestError("Email is already taken")
	}
	hashPass, err := generateFromPassword(data.Password)
	if err != nil {
		return nil, common.NewInternalServerError()
	}
	user := users.UserInsertInput{Email: data.Email, Password: hashPass, Name: data.Name}
	newUser, err := s.usersService.Create(user)
	if err != nil {
		return nil, common.NewInternalServerError()
	}
	tokens, err := s.tokenService.generate(*newUser)
	if err != nil {
		return nil, err
	}

	sessionData := s.buildSessionData(newUser.ID, tokens.RefreshToken)
	err = s.sessionsService.Create(ctx, tokens.RefreshToken, sessionData)
	if err != nil {
		return nil, common.NewInternalServerError() // maybe change later
	}

	return tokens, nil
}

func (s *Service) login(ctx context.Context, data *LoginDto) (*Tokens, error) {
	user, err := s.validateUser(data.Email, data.Password)
	if err != nil {
		return nil, err
	}

	tokens, err := s.tokenService.generate(*user)
	if err != nil {
		return nil, err
	}

	sessionData := s.buildSessionData(user.ID, tokens.RefreshToken)
	err = s.sessionsService.Create(ctx, tokens.RefreshToken, sessionData)
	if err != nil {
		return nil, common.NewInternalServerError() // maybe change later
	}

	return tokens, nil
}

func (s *Service) logout(ctx context.Context, refreshToken string) error {
	return s.sessionsService.Delete(ctx, refreshToken)
}

func (s *Service) validateUser(email string, password string) (*users.User, error) {
	filter := users.UserFilter{
		Email: &email,
	}
	user, err := s.usersService.FindOne(filter)
	if err != nil {
		return nil, common.NewBadRequestError("Wrong email or password")
	}

	match, err := comparePasswordAndHash(password, user.Password)
	if err != nil {
		log.Println("Password hashing error:", err)
		return nil, err
	}
	if !match {
		return nil, common.NewBadRequestError("Wrong email or password")
	}
	return user, nil
}

func (s *Service) buildSessionData(userID uint, refreshToken string) sessions.UserSession {
	return sessions.UserSession{
		UserId:       strconv.FormatUint(uint64(userID), 10),
		RefreshToken: refreshToken,
		CreatedAt:    time.Now().UnixMilli(),
		ExpiresIn:    30 * 24 * time.Hour,
	}
}
