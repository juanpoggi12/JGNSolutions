package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userAuthRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	UpdatePassword(id primitive.ObjectID, hashed string) (*mongo.UpdateResult, error)
}

type sessionAuthRepository interface {
	Create(ctx context.Context, session *models.Session) error
	FindActive(ctx context.Context) ([]models.Session, error)
	FindActiveByUser(ctx context.Context, userID primitive.ObjectID) ([]models.Session, error)
	MarkRevoked(ctx context.Context, id primitive.ObjectID, revokedAt time.Time, replacedByID *primitive.ObjectID) error
}

// Define the userProfileRepository interface
type userProfileRepository interface {
	Create(ctx context.Context, profile *models.UserProfile) error
}

type AuthService struct {
	users      userAuthRepository
	sessions   sessionAuthRepository
	profiles   userProfileRepository
	logService *LogService
	cfg        utils.Config
}

func NewAuthService(users userAuthRepository, sessions sessionAuthRepository, profiles userProfileRepository, logService *LogService, cfg utils.Config) *AuthService {
	return &AuthService{
		users:      users,
		sessions:   sessions,
		profiles:   profiles,
		logService: logService,
		cfg:        cfg,
	}
}

var (
	ErrEmailExists        = errors.New("el email ya está registrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrInvalidRefresh     = errors.New("refresh token inválido")
	ErrExpiredRefresh     = errors.New("refresh token expirado")
)

func (s *AuthService) Register(ctx context.Context, req dto.RegisterReq) error {

	exists, err := s.users.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailExists
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	now := time.Now()
	user := models.User{
		Email:        req.Email,
		Username:     req.Email,
		PasswordHash: hashed,
		Role:         models.RoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}

	if err := s.users.Create(ctx, &user); err != nil {
		return err
	}

	if s.profiles != nil {
		profile := models.UserProfile{
			UserID:    user.ID,
			FullName:  req.Name,
			Goal:      models.ObjetivoMantenerse,
			Level:     models.NivelPrincipiante,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_ = s.profiles.Create(ctx, &profile)
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {

	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			return errors.New("no se pudo procesar la solicitud para este email")
		}

		return err
	}
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("error al procesar la nueva contraseña")
	}

	_, err = s.users.UpdatePassword(user.ID, hashedPassword)
	if err != nil {
		return errors.New("error al actualizar la contraseña en la base de datos")
	}

	if s.logService != nil {
		s.logService.RecordAction(user.ID, "password_reset")
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginReq, userAgent, ip string) (string, int64, func(http.ResponseWriter), error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", 0, nil, ErrInvalidCredentials
		}
		return "", 0, nil, err
	}

	if !user.IsActive {
		return "", 0, nil, errors.New("usuario inactivo o bloqueado")
	}

	if !utils.VerifyPassword(req.Password, user.PasswordHash) {
		return "", 0, nil, ErrInvalidCredentials
	}

	accessTTL := time.Duration(s.cfg.AccessTTLMinutes) * time.Minute
	accessToken, err := utils.GenerateAccessToken(user.ID.Hex(), string(user.Role), s.cfg.JWTSecret, accessTTL)
	if err != nil {
		return "", 0, nil, err
	}

	plainRefresh, refreshHash, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", 0, nil, err
	}

	now := time.Now()
	session := models.Session{
		UserID:      user.ID,
		RefreshHash: refreshHash,
		ExpiresAt:   now.Add(time.Duration(s.cfg.RefreshTTLDays) * 24 * time.Hour),
		CreatedAt:   now,
		UserAgent:   userAgent,
		IP:          ip,
	}

	if err := s.sessions.Create(ctx, &session); err != nil {
		return "", 0, nil, err
	}

	if s.logService != nil {
		s.logService.RecordAction(user.ID, "login")
	}

	setCookie := func(w http.ResponseWriter) {
		http.SetCookie(w, buildRefreshCookie(plainRefresh, session.ExpiresAt, s.cfg.CookieSecure))
	}

	return accessToken, int64(accessTTL.Seconds()), setCookie, nil
}

func (s *AuthService) Refresh(ctx context.Context, plainRefresh, userAgent, ip string) (string, int64, func(http.ResponseWriter), error) {
	session, err := s.findSessionByRefresh(ctx, plainRefresh)
	if err != nil {
		return "", 0, nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessions.MarkRevoked(ctx, session.ID, time.Now(), nil)
		return "", 0, nil, ErrExpiredRefresh
	}

	user, err := s.users.FindByID(ctx, session.UserID)
	if err != nil {
		return "", 0, nil, err
	}

	accessTTL := time.Duration(s.cfg.AccessTTLMinutes) * time.Minute
	accessToken, err := utils.GenerateAccessToken(user.ID.Hex(), string(user.Role), s.cfg.JWTSecret, accessTTL)
	if err != nil {
		return "", 0, nil, err
	}

	plainRefreshNew, refreshHashNew, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", 0, nil, err
	}

	now := time.Now()
	newSession := models.Session{
		UserID:      session.UserID,
		RefreshHash: refreshHashNew,
		ExpiresAt:   now.Add(time.Duration(s.cfg.RefreshTTLDays) * 24 * time.Hour),
		CreatedAt:   now,
		UserAgent:   userAgent,
		IP:          ip,
	}

	if err := s.sessions.Create(ctx, &newSession); err != nil {
		return "", 0, nil, err
	}

	if err := s.sessions.MarkRevoked(ctx, session.ID, now, &newSession.ID); err != nil {
		return "", 0, nil, err
	}

	if s.logService != nil {
		s.logService.RecordAction(session.UserID, "refresh")
	}

	setCookie := func(w http.ResponseWriter) {
		http.SetCookie(w, buildRefreshCookie(plainRefreshNew, newSession.ExpiresAt, s.cfg.CookieSecure))
	}

	return accessToken, int64(accessTTL.Seconds()), setCookie, nil
}

func (s *AuthService) Logout(ctx context.Context, plainRefresh string) (func(http.ResponseWriter), error) {
	session, err := s.findSessionByRefresh(ctx, plainRefresh)
	if err != nil {
		return nil, err
	}

	if err := s.sessions.MarkRevoked(ctx, session.ID, time.Now(), nil); err != nil {
		return nil, err
	}

	if s.logService != nil {
		s.logService.RecordAction(session.UserID, "logout")
	}

	clear := func(w http.ResponseWriter) {
		http.SetCookie(w, clearRefreshCookie(s.cfg.CookieSecure))
	}

	return clear, nil
}

func (s *AuthService) findSessionByRefresh(ctx context.Context, plain string) (*models.Session, error) {
	sessions, err := s.sessions.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	for _, sess := range sessions {
		if err := utils.CompareRefresh(plain, sess.RefreshHash); err == nil {
			copy := sess
			return &copy, nil
		}
	}

	return nil, ErrInvalidRefresh
}

func buildRefreshCookie(value string, expires time.Time, secure bool) *http.Cookie {
	maxAge := int(time.Until(expires).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Path:     "/api/auth",
		Expires:  expires,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func clearRefreshCookie(secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}
