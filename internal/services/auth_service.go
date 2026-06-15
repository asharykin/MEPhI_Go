package services

import (
	"banksystem/internal/model"
	"banksystem/internal/repositories"
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repositories.UserRepository
	jwtService *JWTService
	db         *sql.DB
}

func NewAuthService(db *sql.DB, userRepo *repositories.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		db:         db,
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	exists, err := s.userRepo.CheckEmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return model.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = s.userRepo.Create(ctx, tx, user)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", model.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", model.ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(int64(user.ID))
	if err != nil {
		return "", err
	}

	return token, nil
}
