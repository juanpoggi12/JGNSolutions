package dto

type AdminRoleUpdateRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user ADMIN USER"`
}

type AdminUserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	UpdatedAt string `json:"updatedAt"`
}

type AdminMetricsSummaryResponse struct {
	UsersCount           int64 `json:"usersCount"`
	ExercisesCount       int64 `json:"exercisesCount"`
	RoutinesCount        int64 `json:"routinesCount"`
	WorkoutSessionsCount int64 `json:"workoutSessionsCount"`
}
