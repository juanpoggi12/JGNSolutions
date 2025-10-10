package services

import (
	"errors"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"
)

// AuthService maneja la lógica de autenticación
type AuthService struct {
	userRepo    repositories.UserRepositoryInterface
	sessionRepo *repositories.SessionRepository
	cfg         utils.Config
}

// Constructor del servicio de autenticación
func NewAuthService(
	userRepo repositories.UserRepositoryInterface,
	sessionRepo *repositories.SessionRepository,
	cfg utils.Config,
) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo, cfg: cfg}
}

// Login valida credenciales y genera los tokens de sesión (access y refresh)
func (s *AuthService) Login(emailOrUsername, password string) (string, string, error) {
	// 1️⃣ Buscar usuario por email o username
	user, err := s.userRepo.FindByEmailOrUsername(emailOrUsername)
	if err != nil {
		return "", "", errors.New("usuario no encontrado")
	}

	// 2️⃣ Verificar contraseña usando bcrypt
	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", "", errors.New("contraseña incorrecta")
	}

	// 3️⃣ Generar Access Token (JWT)
	// GenerateToken devuelve el token, la fecha de expiración (que acá no usamos), y un posible error
	accessToken, _, err := utils.GenerateToken(user.ID.Hex(), string(user.Role), s.cfg)
	if err != nil {
		return "", "", err
	}

	// 4️⃣ Generar Refresh Token (cadena aleatoria) y guardarla hasheada en la BD
	refreshToken := utils.GenerateRefreshToken()
	hash := utils.HashRefreshToken(refreshToken)

	session := models.Session{
		UserID:           user.ID,
		RefreshTokenHash: hash,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour), // refresco válido por 7 días
		CreatedAt:        time.Now(),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return "", "", err
	}

	// 5️⃣ Devolver ambos tokens
	return accessToken, refreshToken, nil
}
