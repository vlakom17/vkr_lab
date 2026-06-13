package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"charts-user-service/internal/domain/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	user         *user.User
	getByIDErr   error
	updateErr    error
	updateCalled bool
}

func (r *fakeUserRepository) Create(ctx context.Context, u *user.User) error {
	return nil
}

func (r *fakeUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}

func (r *fakeUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}

	return r.user, nil
}

func (r *fakeUserRepository) Update(ctx context.Context, u *user.User) error {
	r.updateCalled = true

	if r.updateErr != nil {
		return r.updateErr
	}

	r.user = u
	return nil
}

func TestGetUserProfile_ReturnsUser(t *testing.T) {
	id := uuid.New()

	expectedUser := &user.User{
		ID:           id,
		Name:         "Test User",
		Email:        "test@test.com",
		PasswordHash: "hash",
		About:        "about",
		CreatedAt:    time.Now(),
	}

	repo := &fakeUserRepository{
		user: expectedUser,
	}

	service := NewUserService(repo)

	result, err := service.GetUserProfile(context.Background(), id)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expectedUser {
		t.Errorf("expected user %v, got %v", expectedUser, result)
	}
}
func TestGetUserProfile_ReturnsErrorWhenGetByIDFails(t *testing.T) {
	id := uuid.New()

	expectedErr := errors.New("user not found")

	repo := &fakeUserRepository{
		getByIDErr: expectedErr,
	}

	service := NewUserService(repo)

	result, err := service.GetUserProfile(context.Background(), id)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
func TestUpdateUser_UpdatesNameEmailAndAbout(t *testing.T) {
	id := uuid.New()

	repo := &fakeUserRepository{
		user: &user.User{
			ID:           id,
			Name:         "Old Name",
			Email:        "old@test.com",
			PasswordHash: "hash",
			About:        "old about",
			CreatedAt:    time.Now(),
		},
	}

	service := NewUserService(repo)

	newName := "New Name"
	newEmail := "new@test.com"
	newAbout := "new about"

	result, err := service.UpdateUser(
		context.Background(),
		id,
		&newName,
		&newEmail,
		nil,
		&newAbout,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != newName {
		t.Errorf("expected name %s, got %s", newName, result.Name)
	}

	if result.Email != newEmail {
		t.Errorf("expected email %s, got %s", newEmail, result.Email)
	}

	if result.About != newAbout {
		t.Errorf("expected about %s, got %s", newAbout, result.About)
	}
}

func TestUpdateUser_UpdatesPasswordHash(t *testing.T) {
	id := uuid.New()

	oldHash := "old_hash"

	repo := &fakeUserRepository{
		user: &user.User{
			ID:           id,
			Name:         "Test User",
			Email:        "test@test.com",
			PasswordHash: oldHash,
			About:        "about",
			CreatedAt:    time.Now(),
		},
	}

	service := NewUserService(repo)

	newPassword := "new-password-123"

	result, err := service.UpdateUser(
		context.Background(),
		id,
		nil,
		nil,
		&newPassword,
		nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.PasswordHash == oldHash {
		t.Errorf("expected password hash to change")
	}

	if result.PasswordHash == newPassword {
		t.Errorf("password must not be stored as plain text")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(result.PasswordHash),
		[]byte(newPassword),
	)

	if err != nil {
		t.Errorf("password hash does not match new password: %v", err)
	}
}

func TestUpdateUser_ReturnsErrorWhenGetByIDFails(t *testing.T) {
	id := uuid.New()

	expectedErr := errors.New("user not found")

	repo := &fakeUserRepository{
		getByIDErr: expectedErr,
	}

	service := NewUserService(repo)

	newName := "New Name"

	result, err := service.UpdateUser(
		context.Background(),
		id,
		&newName,
		nil,
		nil,
		nil,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.updateCalled {
		t.Errorf("Update should not be called when GetByID fails")
	}
}

func TestUpdateUser_ReturnsErrorWhenUpdateFails(t *testing.T) {
	id := uuid.New()

	expectedErr := errors.New("update failed")

	repo := &fakeUserRepository{
		user: &user.User{
			ID:           id,
			Name:         "Old Name",
			Email:        "old@test.com",
			PasswordHash: "hash",
			About:        "old about",
			CreatedAt:    time.Now(),
		},
		updateErr: expectedErr,
	}

	service := NewUserService(repo)

	newName := "New Name"

	result, err := service.UpdateUser(
		context.Background(),
		id,
		&newName,
		nil,
		nil,
		nil,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if !repo.updateCalled {
		t.Errorf("expected Update to be called")
	}
}
