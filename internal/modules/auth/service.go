package auth

import (
	"caronago/internal/platform/apierror"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpiry int //in hours
}

func NewService(db *gorm.DB, jwtSecret string, jwtExpiry int) *Service {

	return &Service{
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	var count int64

	if err := s.db.WithContext(ctx).Model(&User{}).Where("email = ? ", input.Email).Count(&count).Error; err != nil {
		return nil, apierror.Internal("Database error checking email", err)
	}

	if count > 0 {
		return nil, apierror.Conflict("Email is already registered", nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apierror.Internal("Failed to hash password", err)
	}

	user := User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Gender:    input.Gender,
		Role:      RoleUser,
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, apierror.Internal("Failed to create user", err)
	}

	token, err := s.generateToken(&user)

	if err != nil {
		return nil, apierror.Internal("Failed to generate authentication token", err)
	}

	return &AuthResponse{
		Token: token,
		User:  &user,
	}, nil

}

func (s *Service) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"role": string(user.Role),
		"exp":  time.Now().Add(time.Duration(s.jwtExpiry) * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
