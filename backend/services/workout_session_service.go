package services

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutSessionService struct {
	repo *repositories.WorkoutSessionRepository
}

func NewWorkoutSessionService(repo *repositories.WorkoutSessionRepository) *WorkoutSessionService {
	return &WorkoutSessionService{repo: repo}
}

func (s *WorkoutSessionService) Create(ctx context.Context, session *models.WorkoutSession) error {
	return s.repo.Create(ctx, session)
}

func (s *WorkoutSessionService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.WorkoutSession, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *WorkoutSessionService) Update(ctx context.Context, session *models.WorkoutSession) error {
	return s.repo.Update(ctx, session)
}

func (s *WorkoutSessionService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *WorkoutSessionService) Search(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error) {
	return s.repo.Search(ctx, filter, opts...)
}
