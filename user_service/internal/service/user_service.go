package service

import (
	"context"

	"charts-user-service/internal/domain/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo user.UserRepository
}

func NewUserService(
	userRepo user.UserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	name *string,
	email *string,
	password *string,
	about *string,
) (*user.User, error) {

	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if name != nil {
		u.Name = *name
	}

	if email != nil {
		u.Email = *email
	}

	if about != nil {
		u.About = *about
	}

	if password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		u.PasswordHash = string(hash)
	}

	err = s.userRepo.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
