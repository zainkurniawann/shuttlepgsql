package services

import (
	// "time"
	// "log"
	"shuttle/models/dto"
	// "shuttle/models/entity"
	"shuttle/repositories"

	"github.com/google/uuid"
	// "github.com/jmoiron/sqlx"
)

type ShuttleServiceInterface interface {
	GetShuttleStatusByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error)
	// AddShuttle(studentUUID uuid.UUID) (*dto.ShuttleResponse, error)
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

// func (service *SchoolService) AddShuttle(req dto.ShuttleRequest, username string) error {
// 	school := entity.School{
// 		ID:          time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
// 		UUID:        uuid.New(),
// 		Name:        req.Name,
// 		Address:     req.Address,
// 		Contact:     req.Contact,
// 		Email:       req.Email,
// 		Description: req.Description,
// 		CreatedBy:   toNullString(username),
// 	}

// 	if err := service.schoolRepository.SaveSchool(school); err != nil {
// 		return err
// 	}

// 	return nil
// }