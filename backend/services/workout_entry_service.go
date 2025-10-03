package services

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutEntryService struct {
	repo *repositories.WorkoutEntryRepository
}

func NewWorkoutEntryService(repo *repositories.WorkoutEntryRepository) *WorkoutEntryService {
	return &WorkoutEntryService{repo: repo}
}

func (s *WorkoutEntryService) Create(ctx context.Context, entry *models.WorkoutEntry) error {
	return s.repo.Create(ctx, entry)
}

func (s *WorkoutEntryService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.WorkoutEntry, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *WorkoutEntryService) Update(ctx context.Context, entry *models.WorkoutEntry) error {
	return s.repo.Update(ctx, entry)
}

func (s *WorkoutEntryService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *WorkoutEntryService) Search(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error) {
	return s.repo.Search(ctx, filter, opts...)
}
