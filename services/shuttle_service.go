package services

import (
	"time"
	"log"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"database/sql"

	"github.com/google/uuid"
	// "github.com/jmoiron/sqlx"
)

type ShuttleServiceInterface interface {
	GetShuttleStatusByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error)
	AddShuttle(req dto.ShuttleRequest, driverUUID, createdBy string) error
	EditShuttleStatus(shuttleUUID, status string) error
}

type ShuttleService struct {
	shuttleRepository repositories.ShuttleRepositoryInterface
}

// NewShuttleService creates a new ShuttleService
func NewShuttleService(shuttleRepository repositories.ShuttleRepositoryInterface) ShuttleServiceInterface {
	return &ShuttleService{
		shuttleRepository: shuttleRepository,
	}
}
func (s *ShuttleService) GetShuttleStatusByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error) {
	// Ambil data dari repository
	shuttles, err := s.shuttleRepository.GetShuttleStatusByParent(parentUUID)
	if err != nil {
		return nil, err
	}

	// Proses data ke DTO
	responses := make([]dto.ShuttleResponse, 0, len(shuttles)) // Pre-allocate slice untuk efisiensi
	for _, shuttle := range shuttles {
		// Membuat DTO response
		response := &dto.ShuttleResponse{
			StudentName: shuttle.StudentName,
			DriverName:  shuttle.DriverName,
			Status:      shuttle.Status,
			CreatedAt:   shuttle.CreatedAt,
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

func (s *ShuttleService) AddShuttle(req dto.ShuttleRequest, driverUUID, createdBy string) error {
	// Validasi StudentUUID
	studentUUID, err := uuid.Parse(req.StudentUUID)
	if err != nil {
		log.Println("Invalid StudentUUID:", req.StudentUUID)
		return err
	}

	// Validasi DriverUUID
	driverUUIDParsed, err := uuid.Parse(driverUUID)
	if err != nil {
		log.Println("Invalid DriverUUID:", driverUUID)
		return err
	}

	// Menetapkan status default jika tidak diberikan
	if req.Status == "" {
		req.Status = "menunggu dijemput"
	}

	// Membuat shuttle entity
	shuttle := entity.Shuttle{
		ShuttleID:   time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		ShuttleUUID: uuid.New(),
		StudentUUID: studentUUID,
		DriverUUID:  driverUUIDParsed,
		Status:      req.Status,
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},  // Perubahan di sini
	}

	// Logging untuk memeriksa data shuttle
	log.Println("Shuttle to be added:", shuttle)

	// Menyimpan shuttle menggunakan repository
	err = s.shuttleRepository.SaveShuttle(shuttle)
	if err != nil {
		log.Println("Error saving shuttle:", err)
		return err
	}

	// Berhasil
	log.Println("Shuttle successfully added:", shuttle)
	return nil
}

func (s *ShuttleService) EditShuttleStatus(shuttleUUID, status string) error {
	// Parse UUID
	shuttleUUIDParsed, err := uuid.Parse(shuttleUUID)
	if err != nil {
		log.Println("Invalid ShuttleUUID:", shuttleUUID)
		return err
	}

	// Update status melalui repository
	if err := s.shuttleRepository.UpdateShuttleStatus(shuttleUUIDParsed, status); err != nil {
		log.Println("Failed to update shuttle status:", err)
		return err
	}

	log.Println("Shuttle status updated:", shuttleUUIDParsed, "New status:", status)
	return nil
}

