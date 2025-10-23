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

	// --- CAMBIO 1: Ajustar fechas para cálculo de mejora ---
	now := time.Now()
	// Período actual (últimos 30 días por defecto si no se especifican from/to)
	currentTo := now
	if to != nil {
		currentTo = *to
	}
	currentFrom := currentTo.AddDate(0, 0, -30)
	if from != nil {
		currentFrom = *from
		// Asegura que el período no sea mayor a 30 días para el cálculo simple
		if currentTo.Sub(currentFrom) > 30*24*time.Hour {
			currentFrom = currentTo.AddDate(0, 0, -30) // Limita a 30 días si es mayor
		}
	}

	// Período anterior (los 30 días antes del período actual)
	previousTo := currentFrom
	previousFrom := previousTo.AddDate(0, 0, -30)

	// Fechas para buscar en la BD (abarca ambos períodos)
	dbSearchFrom := previousFrom
	dbSearchTo := currentTo
	// --- FIN CAMBIO 1 ---

	// Busca sesiones en el rango extendido (60 días)
	sessions, err := s.repository.FindSessions(ctx, actor.UserID, &dbSearchFrom, &dbSearchTo)
	if err != nil {
		return dto.WorkoutSummaryResponse{}, err
	}

	sessionIDs := make([]primitive.ObjectID, 0, len(sessions))
	sessionDateMap := make(map[primitive.ObjectID]time.Time) // Guardamos fecha de cada sesión
	for _, session := range sessions {
		sessionIDs = append(sessionIDs, session.ID)
		sessionDateMap[session.ID] = session.StartTime // Usamos StartTime para asignar a período
	}

	entries, err := s.repository.FindEntriesBySessions(ctx, sessionIDs)
	if err != nil {
		return dto.WorkoutSummaryResponse{}, err
	}

	// --- CAMBIO 2: Agregación separada por período para mejora ---
	type periodMaxWeights struct {
		currentMax  float64
		previousMax float64
	}
	exerciseMaxWeights := make(map[primitive.ObjectID]*periodMaxWeights)
	// --- FIN CAMBIO 2 ---

	type exerciseAccumulator struct {
		totalSets    int
		totalReps    int
		totalWeight  float64
		totalTimeSec int
		maxWeight    float64 // Máximo del período 'current'
		maxReps      int     // Máximo del período 'current'
		maxTime      int     // Máximo del período 'current'
	}
	exerciseAggCurrent := make(map[primitive.ObjectID]*exerciseAccumulator) // Para los KPIs normales

	for _, entry := range entries {
		if entry.WeightUsed == nil || *entry.WeightUsed <= 0 { // Solo consideramos entradas con peso positivo
			continue
		}

		sessionTime := sessionDateMap[entry.WorkoutSessionID]
		isCurrentPeriod := !sessionTime.Before(currentFrom) && sessionTime.Before(currentTo)
		isPreviousPeriod := !sessionTime.Before(previousFrom) && sessionTime.Before(previousTo)

		// --- CAMBIO 3: Calcular máximos por período ---
		maxData := exerciseMaxWeights[entry.ExerciseID]
		if maxData == nil {
			maxData = &periodMaxWeights{}
			exerciseMaxWeights[entry.ExerciseID] = maxData
		}
		if isCurrentPeriod && *entry.WeightUsed > maxData.currentMax {
			maxData.currentMax = *entry.WeightUsed
		}
		if isPreviousPeriod && *entry.WeightUsed > maxData.previousMax {
			maxData.previousMax = *entry.WeightUsed
		}
		// --- FIN CAMBIO 3 ---

		// --- Agregación normal (solo para el período actual 'from' a 'to') ---
		if isCurrentPeriod {
			acc := exerciseAggCurrent[entry.ExerciseID]
			if acc == nil {
				acc = &exerciseAccumulator{}
				exerciseAggCurrent[entry.ExerciseID] = acc
			}

			acc.totalSets++
			if entry.RepsDone != nil {
				acc.totalReps += *entry.RepsDone
				if *entry.RepsDone > acc.maxReps {
					acc.maxReps = *entry.RepsDone
				}
			}
			// Calculo de volumen total (sin cambios)
			weightValue := math.Abs(*entry.WeightUsed)
			if entry.RepsDone != nil {
				acc.totalWeight += weightValue * float64(*entry.RepsDone)
			} else {
				acc.totalWeight += weightValue
			}
			// Max weight *dentro* del período actual
			if *entry.WeightUsed > acc.maxWeight {
				acc.maxWeight = *entry.WeightUsed
			}
			if entry.TimeSec != nil {
				acc.totalTimeSec += *entry.TimeSec
				if *entry.TimeSec > acc.maxTime {
					acc.maxTime = *entry.TimeSec
				}
			}
		}
	}

	// Obtener info de ejercicios (sin cambios)
	exerciseIDsCurrent := make([]primitive.ObjectID, 0, len(exerciseAggCurrent))
	for id := range exerciseAggCurrent {
		exerciseIDsCurrent = append(exerciseIDsCurrent, id)
	}
	exercisesInfo, err := s.repository.FindExercisesByIDs(ctx, exerciseIDsCurrent)
	if err != nil {
		return dto.WorkoutSummaryResponse{}, err
	}

	// Construir respuesta de ByExercise y PRs (solo con datos del período actual)
	byExercise := make([]dto.WorkoutSummaryExercise, 0, len(exerciseAggCurrent))
	prs := make([]dto.WorkoutPersonalRecord, 0, len(exerciseAggCurrent))
	for id, acc := range exerciseAggCurrent {
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
			MaxWeight:  acc.maxWeight, // Max del período actual
			MaxReps:    acc.maxReps,
			MaxTimeSec: acc.maxTime,
		})
	}
	// Ordenar (sin cambios)
	sort.Slice(byExercise, func(i, j int) bool {
		// Order by total weight desc, then name asc
		if byExercise[i].TotalWeight == byExercise[j].TotalWeight {
			return byExercise[i].Name < byExercise[j].Name
		}
		return byExercise[i].TotalWeight > byExercise[j].TotalWeight
	})
	sort.Slice(prs, func(i, j int) bool {
		// Order by max weight desc, then name asc
		if prs[i].MaxWeight == prs[j].MaxWeight {
			return prs[i].Name < prs[j].Name
		}
		return prs[i].MaxWeight > prs[j].MaxWeight
	})

	// Agrupar por período (ByPeriod - solo con sesiones del período actual 'from' a 'to')
	periodMap := map[string]*dto.WorkoutSummaryPeriod{}
	for _, session := range sessions {
		// Filtra sesiones fuera del rango 'currentFrom' a 'currentTo'
		if session.StartTime.Before(currentFrom) || !session.StartTime.Before(currentTo) {
			continue
		}
		key, label := buildWeekBucket(session.StartTime) // O usa buildMonthBucket si prefieres
		bucket := periodMap[key]
		if bucket == nil {
			bucket = &dto.WorkoutSummaryPeriod{Period: key, Label: label}
			periodMap[key] = bucket
		}
		bucket.Sessions++
	}
	// Construir byPeriod (sin cambios en la lógica de ordenación)
	byPeriod := make([]dto.WorkoutSummaryPeriod, 0, len(periodMap))
	periodKeys := make([]string, 0, len(periodMap))
	for k := range periodMap {
		periodKeys = append(periodKeys, k)
	}
	sort.Strings(periodKeys)
	for _, key := range periodKeys {
		byPeriod = append(byPeriod, *periodMap[key])
	}

	// --- CAMBIO 4: Calcular Mejora Promedio ---
	var totalPercentageChange float64
	var improvementCount int
	var improvementPercent *float64 // Usa puntero para poder ser nil

	for _, weights := range exerciseMaxWeights {
		if weights.previousMax > 0 && weights.currentMax > 0 { // Solo si hay datos en ambos períodos
			percentageChange := ((weights.currentMax / weights.previousMax) - 1) * 100
			totalPercentageChange += percentageChange
			improvementCount++
		}
	}

	if improvementCount > 0 {
		avgImprovement := totalPercentageChange / float64(improvementCount)
		improvementPercent = &avgImprovement // Asigna el valor calculado
	}
	// --- FIN CAMBIO 4 ---

	return dto.WorkoutSummaryResponse{
		ByExercise:         byExercise,
		ByPeriod:           byPeriod,
		PRs:                prs,
		ImprovementPercent: improvementPercent, // <-- Añade el resultado
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
	startOfWeek := startOfISOWeek(t)
	// Podrías formatear la etiqueta de forma más amigable si quieres
	label := fmt.Sprintf("Sem %d (%s)", week, startOfWeek.Format("02 Jan"))
	return key, label
}

func buildMonthBucket(t time.Time) (string, string, time.Time) {
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	key := start.Format("2006-01")
	return key, key, start
}

func startOfISOWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	// Monday is 1, Sunday is 7
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return start.AddDate(0, 0, 1-weekday)
}
