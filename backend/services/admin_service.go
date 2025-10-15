package services

import (
	"errors"
	"strings"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
)

type AdminServiceInterface interface {
	CountUsers(actor Actor) (int64, error)
	CountExercises(actor Actor) (int64, error)
	CountRoutines(actor Actor) (int64, error)
	CountWorkoutSessions(actor Actor) (int64, error)
	ListUsers(actor Actor) ([]dto.UserListResponse, error)
	TopExercises(actor Actor, limit int) ([]dto.ExerciseStatResponse, error)
	TopRoutines(actor Actor, limit int) ([]dto.RoutineStatResponse, error)
	ListProfiles(actor Actor) ([]dto.UserProfileListResponse, error)
	CountProfilesByLevel(actor Actor) ([]dto.ProfileLevelStat, error)
	CountProfilesByGoal(actor Actor) ([]dto.ProfileGoalStat, error)
	ListLogs(actor Actor) ([]models.Log, error) // 👈 nuevo método
	UpdateUserRole(actor Actor, id string, role string) (dto.AdminUserResponse, error)
	GetMetricsSummary(actor Actor) (dto.AdminMetricsSummaryResponse, error)
}

type AdminService struct {
	adminRepository       repositories.AdminRepositoryInterface
	userRepository        repositories.UserRepositoryInterface
	userProfileRepository repositories.UserProfileRepositoryInterface
	logRepository         *repositories.LogRepository
}

func NewAdminService(
	adminRepo repositories.AdminRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	profileRepo repositories.UserProfileRepositoryInterface,
	logRepo *repositories.LogRepository,
) *AdminService {
	return &AdminService{
		adminRepository:       adminRepo,
		userRepository:        userRepo,
		userProfileRepository: profileRepo,
		logRepository:         logRepo,
	}
}

func (service *AdminService) CountUsers(actor Actor) (int64, error) {
	if actor.Role != "admin" {
		return 0, errors.New("no tienes permiso para acceder a estas estadísticas")
	}
	return service.adminRepository.ContarDocumentos("users")
}

func (service *AdminService) CountExercises(actor Actor) (int64, error) {
	if actor.Role != "admin" {
		return 0, errors.New("no tienes permiso para acceder a estas estadísticas")
	}
	return service.adminRepository.ContarDocumentos("exercises")
}

func (service *AdminService) CountRoutines(actor Actor) (int64, error) {
	if actor.Role != "admin" {
		return 0, errors.New("no tienes permiso para acceder a estas estadísticas")
	}
	return service.adminRepository.ContarDocumentos("routines")
}

func (service *AdminService) CountWorkoutSessions(actor Actor) (int64, error) {
	if actor.Role != "admin" {
		return 0, errors.New("no tienes permiso para acceder a estas estadísticas")
	}
	return service.adminRepository.ContarDocumentos("workoutSessions")
}

func (s *AdminService) ListUsers(actor Actor) ([]dto.UserListResponse, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver usuarios")
	}
	users, err := s.userRepository.ListarUsuariosBasico()
	if err != nil {
		return nil, err
	}

	out := make([]dto.UserListResponse, 0, len(users))
	for _, u := range users {
		out = append(out, dto.UserListResponse{
			ID:    u.ID.Hex(),
			Email: u.Email,
			Role:  string(u.Role), // 👈 conversión explícita
		})
	}
	return out, nil
}

func (s *AdminService) TopExercises(actor Actor, limit int) ([]dto.ExerciseStatResponse, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver estadísticas")
	}
	rows, err := s.adminRepository.TopExercisesByEntries(limit)
	if err != nil {
		return nil, err
	}

	out := make([]dto.ExerciseStatResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.ExerciseStatResponse{
			ExerciseID: r.ID.Hex(),
			Name:       r.Name,
			UsageCount: r.UsageCount,
		})
	}
	return out, nil
}

func (s *AdminService) TopRoutines(actor Actor, limit int) ([]dto.RoutineStatResponse, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver estadísticas")
	}
	rows, err := s.adminRepository.TopRoutinesBySessions(limit)
	if err != nil {
		return nil, err
	}

	out := make([]dto.RoutineStatResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.RoutineStatResponse{
			RoutineID:  r.ID.Hex(),
			Name:       r.Name,
			UsageCount: r.UsageCount,
		})
	}
	return out, nil
}

// --- 👤 LISTADO DE PERFILES DE USUARIO (solo admin) ---
func (s *AdminService) ListProfiles(actor Actor) ([]dto.UserProfileListResponse, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver perfiles de usuario")
	}

	perfiles, err := s.userProfileRepository.ListarPerfiles()
	if err != nil {
		return nil, err
	}

	out := make([]dto.UserProfileListResponse, 0, len(perfiles))
	for _, p := range perfiles {
		out = append(out, dto.UserProfileListResponse{
			ID:       p.ID.Hex(),
			FullName: p.FullName,
			Level:    string(p.Level),
			Goal:     string(p.Goal),
			WeightKg: p.WeightKg,
			HeightCm: p.HeightCm,
		})
	}

	return out, nil
}

// --- 📊 ESTADÍSTICAS DE PERFILES (solo admin) ---

func (s *AdminService) CountProfilesByLevel(actor Actor) ([]dto.ProfileLevelStat, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver estadísticas de perfiles")
	}

	perfiles, err := s.userProfileRepository.ListarPerfiles()
	if err != nil {
		return nil, err
	}

	stats := map[string]int{}
	for _, p := range perfiles {
		stats[string(p.Level)]++
	}

	out := []dto.ProfileLevelStat{}
	for level, count := range stats {
		out = append(out, dto.ProfileLevelStat{
			Level: level,
			Count: count,
		})
	}

	return out, nil
}

func (s *AdminService) CountProfilesByGoal(actor Actor) ([]dto.ProfileGoalStat, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver estadísticas de perfiles")
	}

	perfiles, err := s.userProfileRepository.ListarPerfiles()
	if err != nil {
		return nil, err
	}

	stats := map[string]int{}
	for _, p := range perfiles {
		stats[string(p.Goal)]++
	}

	out := []dto.ProfileGoalStat{}
	for goal, count := range stats {
		out = append(out, dto.ProfileGoalStat{
			Goal:  goal,
			Count: count,
		})
	}

	return out, nil
}

// --- 📜 LISTADO DE LOGS (solo admin) ---
func (s *AdminService) ListLogs(actor Actor) ([]models.Log, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver los logs del sistema")
	}

	logs, err := s.logRepository.ListarLogs()
	if err != nil {
		return nil, err
	}

	return logs, nil
}
func (s *AdminService) UpdateUserRole(actor Actor, id string, role string) (dto.AdminUserResponse, error) {
	if actor.Role != "admin" {
		return dto.AdminUserResponse{}, errors.New("no tienes permiso para modificar roles")
	}

	usuario, err := s.userRepository.ObtenerUsuarioPorID(id)
	if err != nil {
		return dto.AdminUserResponse{}, errors.New("usuario no encontrado")
	}

	normalized := strings.ToLower(strings.TrimSpace(role))
	if normalized != "admin" && normalized != "user" {
		return dto.AdminUserResponse{}, errors.New("rol inválido")
	}

	nuevoRol := models.Role(normalized)
	usuario.Role = nuevoRol
	usuario.UpdatedAt = time.Now()

	if _, err := s.userRepository.ModificarUsuario(usuario); err != nil {
		return dto.AdminUserResponse{}, err
	}

	return dto.AdminUserResponse{
		ID:        usuario.ID.Hex(),
		Email:     usuario.Email,
		Username:  usuario.Username,
		Role:      string(usuario.Role),
		UpdatedAt: usuario.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *AdminService) GetMetricsSummary(actor Actor) (dto.AdminMetricsSummaryResponse, error) {
	if actor.Role != "admin" {
		return dto.AdminMetricsSummaryResponse{}, errors.New("no tienes permiso para acceder a estas estadísticas")
	}

	usersCount, err := s.adminRepository.ContarDocumentos("users")
	if err != nil {
		return dto.AdminMetricsSummaryResponse{}, err
	}

	exercisesCount, err := s.adminRepository.ContarDocumentos("exercises")
	if err != nil {
		return dto.AdminMetricsSummaryResponse{}, err
	}

	routinesCount, err := s.adminRepository.ContarDocumentos("routines")
	if err != nil {
		return dto.AdminMetricsSummaryResponse{}, err
	}

	workoutSessionsCount, err := s.adminRepository.ContarDocumentos("workoutSessions")
	if err != nil {
		return dto.AdminMetricsSummaryResponse{}, err
	}

	return dto.AdminMetricsSummaryResponse{
		UsersCount:           usersCount,
		ExercisesCount:       exercisesCount,
		RoutinesCount:        routinesCount,
		WorkoutSessionsCount: workoutSessionsCount,
	}, nil
}
