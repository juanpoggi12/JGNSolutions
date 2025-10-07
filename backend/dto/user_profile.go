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
	UpdatedAt string  `json:"updated_at"` // RFC3339
}
