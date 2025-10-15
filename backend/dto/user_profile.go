package dto

type UserProfileCreateRequest struct {
	FullName  string  `json:"full_name" binding:"required,min=2,max=80"`
	BirthDate string  `json:"birth_date" binding:"required,datetime=2006-01-02"` // yyyy-mm-dd
	WeightKg  float64 `json:"weight_kg" binding:"required,gt=0,lte=500"`
	HeightCm  int     `json:"height_cm" binding:"required,gt=0,lte=300"`
	Level     string  `json:"level" binding:"required,oneof=PRINCIPIANTE INTERMEDIO AVANZADO"`
	Goal      string  `json:"goal" binding:"required,oneof=PERDER_PESO GANAR_MUSCULO MANTENERSE"`
}

type UserProfileUpdateRequest struct {
	FullName  *string  `json:"full_name" binding:"omitempty,min=2,max=80"`
	BirthDate *string  `json:"birth_date" binding:"omitempty,datetime=2006-01-02"`
	WeightKg  *float64 `json:"weight_kg" binding:"omitempty,gt=0,lte=500"`
	HeightCm  *int     `json:"height_cm" binding:"omitempty,gt=0,lte=300"`
	Level     *string  `json:"level" binding:"omitempty,oneof=PRINCIPIANTE INTERMEDIO AVANZADO"`
	Goal      *string  `json:"goal" binding:"omitempty,oneof=PERDER_PESO GANAR_MUSCULO MANTENERSE"`
}

type UserProfileSearchRequest struct {
	Name  string `form:"name"`
	Level string `form:"level" binding:"omitempty,oneof=PRINCIPIANTE INTERMEDIO AVANZADO"`
	Goal  string `form:"goal" binding:"omitempty,oneof=PERDER_PESO GANAR_MUSCULO MANTENERSE"`
}

type UserProfileResponse struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	BirthDate string  `json:"birth_date"` // yyyy-mm-dd in mapper
	WeightKg  float64 `json:"weight_kg"`
	HeightCm  int     `json:"height_cm"`
	Level     string  `json:"level"`
	Goal      string  `json:"goal"`
	CreatedAt string  `json:"created_at"` // 👈 nuevo campo (RFC3339)
	UpdatedAt string  `json:"updated_at"` // RFC3339
}

type UserProfileListResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Level    string  `json:"level"`
	Goal     string  `json:"goal"`
	WeightKg float64 `json:"weight_kg"`
	HeightCm int     `json:"height_cm"`
}

// --- DTOs para alias de perfil (/api/profile) ---

type ProfileResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	BirthDate string  `json:"birthDate"`
	Weight    float64 `json:"weight"`
	Height    float64 `json:"height"`
	Level     string  `json:"level"`
	Goal      string  `json:"goal"`
	UpdatedAt string  `json:"updatedAt"`
}

type ProfileUpdateRequest struct {
	Name      *string  `json:"name" binding:"omitempty,min=2,max=80"`
	BirthDate *string  `json:"birthDate" binding:"omitempty,datetime=2006-01-02"`
	Weight    *float64 `json:"weight" binding:"omitempty,gte=0,lte=500"`
	Height    *float64 `json:"height" binding:"omitempty,gte=0,lte=300"`
	Level     *string  `json:"level" binding:"omitempty,oneof=beginner intermediate advanced PRINCIPIANTE INTERMEDIO AVANZADO"`
	Goal      *string  `json:"goal" binding:"omitempty"`
}

type ProfileChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}
