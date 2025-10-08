package services

import (
	"errors"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"
)

type AuthService struct {
	userRepo    repositories.UserRepositoryInterface
	sessionRepo *repositories.SessionRepository
	cfg         utils.Config
}

// Constructor
func NewAuthService(userRepo repositories.UserRepositoryInterface, sessionRepo *repositories.SessionRepository, cfg utils.Config) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo, cfg: cfg}
}

// Login genera access y refresh tokens
func (s *AuthService) Login(emailOrUsername, password string) (string, string, error) {
	// Buscar usuario
	user, err := s.userRepo.FindByEmailOrUsername(emailOrUsername)
	if err != nil {
		return "", "", errors.New("usuario no encontrado")
	}

	// Verificar contraseña (usa bcrypt)
	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", "", errors.New("contraseña incorrecta")
	}

	// Crear Access Token (JWT)
	accessToken, _, err := utils.GenerateToken(user.ID.Hex(), string(user.Role), s.cfg)
	if err != nil {
		return "", "", err
	}

	// Crear Refresh Token (guardado en BD)
	refreshToken := utils.GenerateRefreshToken()
	hash := utils.HashRefreshToken(refreshToken)

	session := models.Session{
		UserID:           user.ID,
		RefreshTokenHash: hash,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:        time.Now(),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
