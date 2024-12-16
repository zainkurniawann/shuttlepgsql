package repositories

import (
	"shuttle/models/dto"
	"shuttle/models/entity"
	// "time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShuttleRepositoryInterface interface {
    GetShuttleStatusByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error)
    SaveShuttle(shuttle *entity.Shuttle) (*entity.Shuttle, error)
}

type ShuttleRepository struct {
    DB *sqlx.DB
}

func NewShuttleRepository(DB *sqlx.DB) ShuttleRepositoryInterface {
	return &ShuttleRepository{
		DB: DB,
	}
}
func (r *ShuttleRepository) GetShuttleStatusByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error) {
	query := `
		SELECT 
			sd.student_first_name || ' ' || sd.student_last_name AS student_name,
			dd.user_first_name || ' ' || dd.user_last_name AS driver_name,
			s.status,
			TO_CHAR(s.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
		FROM 
			shuttle s
		JOIN 
			students sd ON s.student_uuid = sd.student_uuid
		LEFT JOIN 
			driver_details dd ON s.driver_uuid = dd.user_uuid
		WHERE 
			sd.parent_uuid = $1;
	`

	var shuttles []dto.ShuttleResponse
	err := r.DB.Select(&shuttles, query, parentUUID)
	if err != nil {
		return nil, err
	}

	return shuttles, nil
}





// shuttleRepository.go
func (r *ShuttleRepository) SaveShuttle(shuttle *entity.Shuttle) (*entity.Shuttle, error) {
	// Query to insert a new shuttle record
	query := `
		INSERT INTO shuttle (student_uuid, status, created_at)
		VALUES ($1, $2, NOW()) RETURNING shuttle_id, created_at
	`

	// Save shuttle and retrieve the generated ID and created_at
	err := r.DB.QueryRow(query, shuttle.StudentUUID, shuttle.Status).
		Scan(&shuttle.ShuttleID, &shuttle.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Retrieve additional details (e.g., student name, driver name) if needed
	// You can join with the students and driver_details table to get the necessary fields

	return shuttle, nil
}
