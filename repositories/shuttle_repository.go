package repositories

import (
	"shuttle/models/dto"
	"shuttle/models/entity"
	"database/sql"
	// "time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShuttleRepositoryInterface interface {
    GetShuttleStatusByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error)
    SaveShuttle(shuttle entity.Shuttle) error
	UpdateShuttleStatus(shuttleUUID uuid.UUID, status string) error
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
		s.student_uuid,
		s.student_first_name,
		s.student_last_name,
		s.student_gender,
		s.student_grade,
		s.parent_uuid,
		s.school_uuid,
		sc.school_name,
		st.status AS shuttle_status
	FROM students s
	JOIN schools sc ON s.school_uuid = sc.school_uuid
	LEFT JOIN shuttle st ON s.student_uuid = st.student_uuid AND DATE(st.created_at) = CURRENT_DATE
	WHERE s.parent_uuid = $1;

	`

	var shuttles []dto.ShuttleResponse
	err := r.DB.Select(&shuttles, query, parentUUID)
	if err != nil {
		return nil, err
	}

	return shuttles, nil
}

func (r *ShuttleRepository) SaveShuttle(shuttle entity.Shuttle) error {
	query := `
		INSERT INTO shuttle (shuttle_id, shuttle_uuid, student_uuid, driver_uuid, status, created_at)
		VALUES (:shuttle_id, :shuttle_uuid, :student_uuid, :driver_uuid, :status, :created_at)`
	_, err := r.DB.NamedExec(query, shuttle)
	return err
}

func (r *ShuttleRepository) UpdateShuttleStatus(shuttleUUID uuid.UUID, status string) error {
	query := `
		UPDATE shuttle
		SET status = :status, updated_at = NOW()
		WHERE shuttle_uuid = :shuttle_uuid`

	// Data untuk query
	data := map[string]interface{}{
		"status":       status,
		"shuttle_uuid": shuttleUUID,
	}

	// Eksekusi query
	result, err := r.DB.NamedExec(query, data)
	if err != nil {
		return err
	}

	// Periksa apakah row diperbarui
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
