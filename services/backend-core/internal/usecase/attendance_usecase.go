package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/pkg/aiservice"
	"github.com/google/uuid"
)

type AttendanceUsecase interface {
	CheckIn(ctx context.Context, req domain.CheckInRequest) (*domain.CheckInResult, error)
	GetRecentHistory(ctx context.Context, employeeID uuid.UUID) ([]domain.Attendance, error)
}

type attendanceUsecase struct {
	attRepo      domain.AttendanceRepository
	branchRepo   domain.BranchRepository
	empRepo      domain.EmployeeRepository
	aiClient     *aiservice.Client
	minThreshold float64
}

func NewAttendanceUsecase(
	attRepo domain.AttendanceRepository,
	branchRepo domain.BranchRepository,
	empRepo domain.EmployeeRepository,
	aiClient *aiservice.Client,
) AttendanceUsecase {
	return &attendanceUsecase{
		attRepo:      attRepo,
		branchRepo:   branchRepo,
		empRepo:      empRepo,
		aiClient:     aiClient,
		minThreshold: 0.80, // Default 80% cosine similarity
	}
}

func (u *attendanceUsecase) CheckIn(ctx context.Context, req domain.CheckInRequest) (*domain.CheckInResult, error) {
	// 1. Check if employee exists
	emp, err := u.empRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	today := time.Now().Format("2006-01-02")
	existing, _ := u.attRepo.GetTodayAttendance(ctx, req.EmployeeID, today)
	if existing != nil && existing.CheckInTime != nil {
		return nil, errors.New("anda sudah melakukan presensi masuk untuk hari ini")
	}

	// 2. Call AI Biometrics Worker
	aiResult, err := u.aiClient.ExtractFaceAndVerifyLiveness(ctx, req.PhotoBytes, req.Filename)
	if err != nil {
		return nil, fmt.Errorf("biometric processing failed: %w", err)
	}

	if !aiResult.FaceDetected {
		return nil, errors.New("wajah tidak terdeteksi pada foto. Hadapkan wajah tegak lurus ke kamera.")
	}

	if !aiResult.IsLive {
		return nil, fmt.Errorf("anti-spoofing gagal (skor liveness: %.2f). Terindikasi foto layar atau tidak alami", aiResult.LivenessScore)
	}

	// 3. Face Similarity Verification with pgvector
	similarity, err := u.empRepo.CheckFaceSimilarity(ctx, req.EmployeeID, aiResult.Embedding)
	if err != nil {
		return nil, fmt.Errorf("face comparison query failed: %w", err)
	}

	// If no master photo is registered yet, treat current photo as the enrolled master!
	if similarity == -1.0 {
		bio := &domain.EmployeeBiometrics{
			CompanyID:           emp.CompanyID,
			EmployeeID:          emp.ID,
			MasterImageURL:      "local://initial_enrollment.jpg",
			FaceEmbedding:       aiResult.Embedding,
			LivenessScoreSample: aiResult.LivenessScore,
		}
		_ = u.empRepo.SaveBiometrics(ctx, bio)
		similarity = 1.0 // First registration counts as matched
	} else if similarity < u.minThreshold {
		return nil, fmt.Errorf("kemiripan wajah tidak sesuai (skor: %.2f, minimum: %.2f)", similarity, u.minThreshold)
	}

	// 4. Geofencing Verification with PostGIS
	var isInside bool
	var distance float64
	geoCheck, err := u.branchRepo.CheckGeofence(ctx, req.BranchID, req.Latitude, req.Longitude)
	if err == nil && geoCheck != nil {
		isInside = geoCheck.IsInside
		distance = geoCheck.DistanceMeters
	} else {
		// Fallback if branch has no strict geofence or employee has WFH rights
		if emp.AllowWFH {
			isInside = true
			distance = 0.0
		}
	}

	status := domain.StatusOnTime
	if !isInside {
		if !emp.AllowWFH {
			status = domain.StatusOutOfGeofence
		} else {
			status = domain.StatusWFH
		}
	}

	// 5. Save Attendance Record
	att := &domain.Attendance{
		CompanyID:              emp.CompanyID,
		EmployeeID:             emp.ID,
		BranchID:               &req.BranchID,
		AttendanceDate:         today,
		CheckInLat:             req.Latitude,
		CheckInLng:             req.Longitude,
		CheckInDistanceMeters:  distance,
		CheckInSimilarityScore: similarity,
		CheckInStatus:          status,
	}

	if err := u.attRepo.RecordCheckIn(ctx, att); err != nil {
		return nil, fmt.Errorf("gagal menyimpan data presensi: %w", err)
	}

	return &domain.CheckInResult{
		AttendanceID:    att.ID,
		Status:          status,
		CheckInTime:     time.Now().Format("15:04:05 WIB"),
		SimilarityScore: similarity,
		LivenessScore:   aiResult.LivenessScore,
		DistanceMeters:  distance,
		InsideGeofence:  isInside,
		Message:         "Presensi masuk berhasil dicatat!",
	}, nil
}

func (u *attendanceUsecase) GetRecentHistory(ctx context.Context, employeeID uuid.UUID) ([]domain.Attendance, error) {
	return u.attRepo.ListRecent(ctx, employeeID, 20)
}
