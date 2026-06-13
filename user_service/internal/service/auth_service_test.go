package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"charts-user-service/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type fakeAuthUserRepository struct {
	user          *user.User
	createdUser   *user.User
	getByEmailErr error
	createErr     error
	createCalled  bool
}

func (r *fakeAuthUserRepository) Create(ctx context.Context, u *user.User) error {
	r.createCalled = true

	if r.createErr != nil {
		return r.createErr
	}

	r.createdUser = u
	r.user = u
	return nil
}

func (r *fakeAuthUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	if r.getByEmailErr != nil {
		return nil, r.getByEmailErr
	}

	return r.user, nil
}

func (r *fakeAuthUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return r.user, nil
}

func (r *fakeAuthUserRepository) Update(ctx context.Context, u *user.User) error {
	r.user = u
	return nil
}

type fakeAuthRepository struct {
	token        string
	userID       uuid.UUID
	createCalled bool
	deleteCalled bool
	createErr    error
	deleteErr    error
	getUserIDErr error
}

func (r *fakeAuthRepository) Create(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error {
	if r.createErr != nil {
		return r.createErr
	}

	r.token = token
	r.userID = userID
	r.createCalled = true
	return nil
}

func (r *fakeAuthRepository) GetUserID(ctx context.Context, token string) (uuid.UUID, error) {
	if r.getUserIDErr != nil {
		return uuid.Nil, r.getUserIDErr
	}

	return r.userID, nil
}

func (r *fakeAuthRepository) Delete(ctx context.Context, token string) error {
	r.deleteCalled = true

	if r.deleteErr != nil {
		return r.deleteErr
	}

	r.token = ""
	return nil
}

func (r *fakeAuthRepository) GetTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	return r.token, nil
}

func TestRegister_CreatesUserWhenEmailIsFree(t *testing.T) {
	userRepo := &fakeAuthUserRepository{
		getByEmailErr: pgx.ErrNoRows,
	}

	authRepo := &fakeAuthRepository{}

	service := NewAuthService(userRepo, authRepo)

	name := "Test User"
	email := "test@test.com"
	password := "password-123"
	about := "about me"

	result, err := service.Register(
		context.Background(),
		name,
		email,
		password,
		about,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected user, got nil")
	}

	if result.ID == uuid.Nil {
		t.Errorf("expected generated user ID")
	}

	if result.Name != name {
		t.Errorf("expected name %s, got %s", name, result.Name)
	}

	if result.Email != email {
		t.Errorf("expected email %s, got %s", email, result.Email)
	}

	if result.About != about {
		t.Errorf("expected about %s, got %s", about, result.About)
	}

	if result.PasswordHash == password {
		t.Errorf("password must not be stored as plain text")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(result.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		t.Errorf("password hash does not match password: %v", err)
	}

	if userRepo.createdUser != result {
		t.Errorf("expected created user to be saved in repository")
	}
}

func TestRegister_ReturnsErrorWhenUserAlreadyExists(t *testing.T) {
	existingUser := &user.User{
		ID:           uuid.New(),
		Name:         "Existing User",
		Email:        "test@test.com",
		PasswordHash: "hash",
		About:        "about",
	}

	userRepo := &fakeAuthUserRepository{
		user: existingUser,
	}

	authRepo := &fakeAuthRepository{}

	service := NewAuthService(userRepo, authRepo)

	result, err := service.Register(
		context.Background(),
		"New User",
		"test@test.com",
		"password-123",
		"new about",
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if userRepo.createCalled {
		t.Errorf("Create should not be called when user already exists")
	}
}

func TestRegister_ReturnsErrorWhenGetByEmailFails(t *testing.T) {
	expectedErr := errors.New("database error")

	userRepo := &fakeAuthUserRepository{
		getByEmailErr: expectedErr,
	}

	authRepo := &fakeAuthRepository{}

	service := NewAuthService(userRepo, authRepo)

	result, err := service.Register(
		context.Background(),
		"Test User",
		"test@test.com",
		"password-123",
		"about",
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

	if userRepo.createCalled {
		t.Errorf("Create should not be called when GetByEmail fails")
	}
}
func TestRegister_ReturnsErrorWhenCreateFails(t *testing.T) {
	expectedErr := errors.New("create failed")

	userRepo := &fakeAuthUserRepository{
		getByEmailErr: pgx.ErrNoRows,
		createErr:     expectedErr,
	}

	authRepo := &fakeAuthRepository{}

	service := NewAuthService(userRepo, authRepo)

	result, err := service.Register(
		context.Background(),
		"Test User",
		"test@test.com",
		"password-123",
		"about",
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

	if !userRepo.createCalled {
		t.Errorf("expected Create to be called")
	}
}

func TestLogin_ReturnsTokenWhenCredentialsAreValid(t *testing.T) {
	id := uuid.New()

	password := "password-123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	userRepo := &fakeAuthUserRepository{
		user: &user.User{
			ID:           id,
			Name:         "Test User",
			Email:        "test@test.com",
			PasswordHash: string(hash),
			About:        "about",
		},
	}

	authRepo := &fakeAuthRepository{}

	service := NewAuthService(userRepo, authRepo)

	token, err := service.Login(
		context.Background(),
		"test@test.com",
		password,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == "" {
		t.Fatalf("expected token, got empty string")
	}

	if !authRepo.createCalled {
		t.Errorf("expected auth token to be created")
	}

	if authRepo.userID != id {
		t.Errorf("expected userID %s, got %s", id, authRepo.userID)
	}

	if authRepo.token != token {
		t.Errorf("expected saved token to match returned token")
	}
}
func TestLogin_ReturnsErrorWhenPasswordIsInvalid(t *testing.T) {
	id := uuid.New()

	correctPassword := "password-123"
	wrongPassword := "wrong-password"

	hash, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	userRepo := &fakeAuthUserRepository{
		user: &user.User{
			ID:           id,
			Name:         "Test User",
			Email:        "test@test.com",
			PasswordHash: string(hash),
			About:        "about",
		},
	}

	authRepo := &fakeAuthRepository{}

	service := NewAuthService(userRepo, authRepo)

	token, err := service.Login(
		context.Background(),
		"test@test.com",
		wrongPassword,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if token != "" {
		t.Errorf("expected empty token, got %s", token)
	}

	if authRepo.createCalled {
		t.Errorf("Create should not be called when password is invalid")
	}
}

func TestLogin_ReturnsErrorWhenUserDoesNotExist(t *testing.T) {
	userRepo := &fakeAuthUserRepository{
		user: nil,
	}

	authRepo := &fakeAuthRepository{}

	service := NewAuthService(userRepo, authRepo)

	token, err := service.Login(
		context.Background(),
		"unknown@test.com",
		"password-123",
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if token != "" {
		t.Errorf("expected empty token, got %s", token)
	}

	if authRepo.createCalled {
		t.Errorf("Create should not be called when user does not exist")
	}
}

func TestLogin_ReturnsErrorWhenGetByEmailFails(t *testing.T) {
	expectedErr := errors.New("database error")

	userRepo := &fakeAuthUserRepository{
		getByEmailErr: expectedErr,
	}

	authRepo := &fakeAuthRepository{}

	service := NewAuthService(userRepo, authRepo)

	token, err := service.Login(
		context.Background(),
		"test@test.com",
		"password-123",
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if token != "" {
		t.Errorf("expected empty token, got %s", token)
	}

	if authRepo.createCalled {
		t.Errorf("Create should not be called when GetByEmail fails")
	}
}

func TestLogout_DeletesToken(t *testing.T) {
	authRepo := &fakeAuthRepository{
		token: "test-token",
	}

	service := NewAuthService(&fakeAuthUserRepository{}, authRepo)

	err := service.Logout(context.Background(), "test-token")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !authRepo.deleteCalled {
		t.Errorf("expected Delete to be called")
	}

	if authRepo.token != "" {
		t.Errorf("expected token to be deleted")
	}
}
func TestLogout_ReturnsErrorWhenDeleteFails(t *testing.T) {
	expectedErr := errors.New("delete failed")

	authRepo := &fakeAuthRepository{
		token:     "test-token",
		deleteErr: expectedErr,
	}

	service := NewAuthService(&fakeAuthUserRepository{}, authRepo)

	err := service.Logout(context.Background(), "test-token")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestGetUserIDByToken_ReturnsUserID(t *testing.T) {
	id := uuid.New()

	authRepo := &fakeAuthRepository{
		token:  "test-token",
		userID: id,
	}

	service := NewAuthService(&fakeAuthUserRepository{}, authRepo)

	result, err := service.GetUserIDByToken(context.Background(), "test-token")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != id {
		t.Errorf("expected userID %s, got %s", id, result)
	}
}

func TestGetUserIDByToken_ReturnsErrorWhenRepositoryFails(t *testing.T) {
	expectedErr := errors.New("token not found")

	authRepo := &fakeAuthRepository{
		getUserIDErr: expectedErr,
	}

	service := NewAuthService(&fakeAuthUserRepository{}, authRepo)

	result, err := service.GetUserIDByToken(context.Background(), "bad-token")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != uuid.Nil {
		t.Errorf("expected nil UUID, got %s", result)
	}
}
