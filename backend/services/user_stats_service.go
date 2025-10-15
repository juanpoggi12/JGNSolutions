package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatsServiceInterface interface {
	GetWorkoutSummary(actor Actor, from, to *time.Time) (dto.WorkoutSummaryResponse, error)
	GetFrequency(actor Actor, from, to *time.Time, granularity string) (dto.UserFrequencyResponse, error)
	GetTopRoutines(actor Actor, limit int) ([]dto.UserTopRoutine, error)
	GetExerciseProgress(actor Actor, exerciseID primitive.ObjectID, from, to *time.Time) (dto.UserExerciseProgressResponse, error)
}

type UserStatsService struct {
	repository repositories.UserStatsRepositoryInterface
}

func NewUserStatsService(repository repositories.UserStatsRepositoryInterface) *UserStatsService {
	return &UserStatsService{repository: repository}
}

func (s *UserStatsService) GetWorkoutSummary(actor Actor, from, to *time.Time) (dto.WorkoutSummaryResponse, error) {
	if actor.UserID.IsZero() {
		return dto.WorkoutSummaryResponse{}, errors.New("usuario no autenticado")
	}

	ctx := context.Background()

	sessions, err := s.repository.FindSessions(ctx, actor.UserID, from, to)
	if err != nil {
		return dto.WorkoutSummaryResponse{}, err
	}

	sessionIDs := make([]primitive.ObjectID, 0, len(sessions))
	for _, session := range sessions {
		sessionIDs = append(sessionIDs, session.ID)
	}

	entries, err := s.repository.FindEntriesBySessions(ctx, sessionIDs)
	if err != nil {
		return dto.WorkoutSummaryResponse{}, err
	}

	type exerciseAccumulator struct {
		totalSets    int
		totalReps    int
		totalWeight  float64
		totalTimeSec int
		maxWeight    float64
		maxReps      int
		maxTime      int
	}

	exerciseAgg := make(map[primitive.ObjectID]*exerciseAccumulator)

	for _, entry := range entries {
		acc := exerciseAgg[entry.ExerciseID]
		if acc == nil {
			acc = &exerciseAccumulator{}
			exerciseAgg[entry.ExerciseID] = acc
		}

		acc.totalSets++
		if entry.RepsDone != nil {
			acc.totalReps += *entry.RepsDone
			if *entry.RepsDone > acc.maxReps {
				acc.maxReps = *entry.RepsDone
			}
		}
		if entry.WeightUsed != nil {
			if entry.RepsDone != nil {
				acc.totalWeight += math.Abs(*entry.WeightUsed) * float64(*entry.RepsDone)
			} else {
				acc.totalWeight += math.Abs(*entry.WeightUsed)
			}
			if *entry.WeightUsed > acc.maxWeight {
				acc.maxWeight = *entry.WeightUsed
			}
		}
		if entry.TimeSec != nil {
			acc.totalTimeSec += *entry.TimeSec
			if *entry.TimeSec > acc.maxTime {
				acc.maxTime = *entry.TimeSec
			}
		}
	}

	exerciseIDs := make([]primitive.ObjectID, 0, len(exerciseAgg))
	for id := range exerciseAgg {
		exerciseIDs = append(exerciseIDs, id)
	}

	exercisesInfo, err := s.repository.FindExercisesByIDs(ctx, exerciseIDs)
	if err != nil {
		return dto.WorkoutSummaryResponse{}, err
	}

	byExercise := make([]dto.WorkoutSummaryExercise, 0, len(exerciseAgg))
	prs := make([]dto.WorkoutPersonalRecord, 0, len(exerciseAgg))
	for id, acc := range exerciseAgg {
		info := exercisesInfo[id]
		byExercise = append(byExercise, dto.WorkoutSummaryExercise{
			ExerciseID:   id.Hex(),
			Name:         info.Name,
			MuscleGroup:  strings.ToLower(string(info.MuscleGroup)),
			TotalSets:    acc.totalSets,
			TotalReps:    acc.totalReps,
			TotalWeight:  acc.totalWeight,
			TotalTimeSec: acc.totalTimeSec,
		})
		prs = append(prs, dto.WorkoutPersonalRecord{
			ExerciseID: id.Hex(),
			Name:       info.Name,
			MaxWeight:  acc.maxWeight,
			MaxReps:    acc.maxReps,
			MaxTimeSec: acc.maxTime,
		})
	}

	sort.Slice(byExercise, func(i, j int) bool {
		if math.Abs(byExercise[i].TotalWeight-byExercise[j].TotalWeight) > 0.01 {
			return byExercise[i].TotalWeight > byExercise[j].TotalWeight
		}
		return byExercise[i].TotalSets > byExercise[j].TotalSets
	})

	sort.Slice(prs, func(i, j int) bool {
		if math.Abs(prs[i].MaxWeight-prs[j].MaxWeight) > 0.01 {
			return prs[i].MaxWeight > prs[j].MaxWeight
		}
		return prs[i].MaxReps > prs[j].MaxReps
	})

	periodMap := map[string]*dto.WorkoutSummaryPeriod{}
	for _, session := range sessions {
		key, label := buildWeekBucket(session.StartTime)
		bucket := periodMap[key]
		if bucket == nil {
			bucket = &dto.WorkoutSummaryPeriod{Period: key, Label: label}
			periodMap[key] = bucket
		}
		bucket.Sessions++
	}

	byPeriod := make([]dto.WorkoutSummaryPeriod, 0, len(periodMap))
	periodKeys := make([]string, 0, len(periodMap))
	for k := range periodMap {
		periodKeys = append(periodKeys, k)
	}
	sort.Strings(periodKeys)
	for _, key := range periodKeys {
		byPeriod = append(byPeriod, *periodMap[key])
	}

	return dto.WorkoutSummaryResponse{
		ByExercise: byExercise,
		ByPeriod:   byPeriod,
		PRs:        prs,
	}, nil
}

func (s *UserStatsService) GetFrequency(actor Actor, from, to *time.Time, granularity string) (dto.UserFrequencyResponse, error) {
	if actor.UserID.IsZero() {
		return dto.UserFrequencyResponse{}, errors.New("usuario no autenticado")
	}

	ctx := context.Background()
	sessions, err := s.repository.FindSessions(ctx, actor.UserID, from, to)
	if err != nil {
		return dto.UserFrequencyResponse{}, err
	}

	mode := strings.ToLower(granularity)
	if mode != "month" {
		mode = "week"
	}

	type bucket struct {
		key   string
		label string
		sort  time.Time
		count int
	}

	buckets := map[string]*bucket{}
	for _, session := range sessions {
		var key, label string
		var sortTime time.Time
		if mode == "month" {
			key, label, sortTime = buildMonthBucket(session.StartTime)
		} else {
			weekKey, weekLabel := buildWeekBucket(session.StartTime)
			key, label = weekKey, weekLabel
			sortTime = startOfISOWeek(session.StartTime)
		}

		b := buckets[key]
		if b == nil {
			b = &bucket{key: key, label: label, sort: sortTime}
			buckets[key] = b
		}
		b.count++
	}

	ordered := make([]*bucket, 0, len(buckets))
	for _, b := range buckets {
		ordered = append(ordered, b)
	}

	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].sort.Before(ordered[j].sort)
	})

	responseBuckets := make([]dto.UserFrequencyBucket, 0, len(ordered))
	for _, b := range ordered {
		responseBuckets = append(responseBuckets, dto.UserFrequencyBucket{
			Period: b.key,
			Label:  b.label,
			Count:  b.count,
		})
	}

	return dto.UserFrequencyResponse{
		Buckets: responseBuckets,
		Total:   len(sessions),
	}, nil
}

func (s *UserStatsService) GetTopRoutines(actor Actor, limit int) ([]dto.UserTopRoutine, error) {
	if actor.UserID.IsZero() {
		return nil, errors.New("usuario no autenticado")
	}
	if limit <= 0 {
		limit = 5
	}

	ctx := context.Background()
	sessions, err := s.repository.FindSessions(ctx, actor.UserID, nil, nil)
	if err != nil {
		return nil, err
	}

	usage := map[primitive.ObjectID]int{}
	routineIDs := make([]primitive.ObjectID, 0)
	seen := map[primitive.ObjectID]struct{}{}
	for _, session := range sessions {
		if session.RoutineID == nil {
			continue
		}
		usage[*session.RoutineID]++
		if _, ok := seen[*session.RoutineID]; !ok {
			seen[*session.RoutineID] = struct{}{}
			routineIDs = append(routineIDs, *session.RoutineID)
		}
	}

	routinesInfo, err := s.repository.FindRoutinesByIDs(ctx, routineIDs)
	if err != nil {
		return nil, err
	}

	out := make([]dto.UserTopRoutine, 0, len(usage))
	for id, count := range usage {
		info := routinesInfo[id]
		out = append(out, dto.UserTopRoutine{
			RoutineID: id.Hex(),
			Name:      info.Name,
			Uses:      count,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Uses == out[j].Uses {
			return out[i].Name < out[j].Name
		}
		return out[i].Uses > out[j].Uses
	})

	if len(out) > limit {
		out = out[:limit]
	}

	return out, nil
}

func (s *UserStatsService) GetExerciseProgress(actor Actor, exerciseID primitive.ObjectID, from, to *time.Time) (dto.UserExerciseProgressResponse, error) {
	if actor.UserID.IsZero() {
		return dto.UserExerciseProgressResponse{}, errors.New("usuario no autenticado")
	}
	if exerciseID.IsZero() {
		return dto.UserExerciseProgressResponse{}, errors.New("exerciseId inválido")
	}

	ctx := context.Background()
	sessions, err := s.repository.FindSessions(ctx, actor.UserID, from, to)
	if err != nil {
		return dto.UserExerciseProgressResponse{}, err
	}

	sessionIDs := make([]primitive.ObjectID, 0, len(sessions))
	sessionInfo := make(map[primitive.ObjectID]models.WorkoutSession, len(sessions))
	for _, s := range sessions {
		sessionIDs = append(sessionIDs, s.ID)
		sessionInfo[s.ID] = s
	}

	entries, err := s.repository.FindEntriesBySessions(ctx, sessionIDs)
	if err != nil {
		return dto.UserExerciseProgressResponse{}, err
	}

	type progressPoint struct {
		sessionID primitive.ObjectID
		date      time.Time
		label     string
		sets      int
		reps      int
		maxWeight float64
		timeSec   int
	}

	points := map[primitive.ObjectID]*progressPoint{}
	for _, entry := range entries {
		if entry.ExerciseID != exerciseID {
			continue
		}

		session := sessionInfo[entry.WorkoutSessionID]
		point := points[entry.WorkoutSessionID]
		if point == nil {
			label := session.StartTime.Format("2006-01-02")
			point = &progressPoint{
				sessionID: entry.WorkoutSessionID,
				date:      session.StartTime,
				label:     label,
			}
			points[entry.WorkoutSessionID] = point
		}

		point.sets++
		if entry.RepsDone != nil {
			point.reps += *entry.RepsDone
		}
		if entry.WeightUsed != nil {
			if *entry.WeightUsed > point.maxWeight {
				point.maxWeight = *entry.WeightUsed
			}
		}
		if entry.TimeSec != nil {
			point.timeSec += *entry.TimeSec
		}
	}

	ordered := make([]*progressPoint, 0, len(points))
	for _, p := range points {
		ordered = append(ordered, p)
	}

	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].date.Before(ordered[j].date)
	})

	resp := dto.UserExerciseProgressResponse{}
	for _, p := range ordered {
		resp.Labels = append(resp.Labels, p.label)
		resp.Sets = append(resp.Sets, p.sets)
		resp.Reps = append(resp.Reps, p.reps)
		resp.Weight = append(resp.Weight, p.maxWeight)
		resp.TimeSec = append(resp.TimeSec, p.timeSec)
	}

	return resp, nil
}

func buildWeekBucket(t time.Time) (string, string) {
	year, week := t.ISOWeek()
	key := fmt.Sprintf("%d-W%02d", year, week)
	return key, key
}

func buildMonthBucket(t time.Time) (string, string, time.Time) {
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	key := start.Format("2006-01")
	return key, key, start
}

func startOfISOWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return start.AddDate(0, 0, 1-weekday)
}
