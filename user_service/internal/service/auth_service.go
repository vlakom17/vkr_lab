package service

import (
	"context"
	"errors"
	"time"

	"charts-user-service/internal/domain/auth"
	"charts-user-service/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo user.UserRepository
	authRepo auth.AuthRepository
}

func NewAuthService(
	userRepo user.UserRepository,
	authRepo auth.AuthRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password, about string) (*user.User, error) {

	existingUser, err := s.userRepo.GetByEmail(ctx, email)

	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		About:        about,
	}

	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {

	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if u == nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("invalid credentials")
	}
	oldToken, _ := s.authRepo.GetTokenByUserID(ctx, u.ID)

	if oldToken != "" {
		s.authRepo.Delete(ctx, oldToken)
	}
	token := uuid.New().String()

	err = s.authRepo.Create(ctx, token, u.ID, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.authRepo.Delete(ctx, token)
}

func (s *AuthService) GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error) {
	return s.authRepo.GetUserID(ctx, token)
}
