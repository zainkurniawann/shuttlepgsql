package repositories

import (
	// "database/sql"
	// "fmt"
	"shuttle/models/entity"
	// "strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StudentRepositoryInterface interface {
	BeginTransaction() (*sqlx.Tx, error)
	SaveStudent(tx *sqlx.Tx, student entity.Student) (uuid.UUID, error)
}

type studentRepository struct {
	DB *sqlx.DB
}

func NewStudentRepository(DB *sqlx.DB) StudentRepositoryInterface {
	return &studentRepository{
		DB: DB,
	}
}

func (r *studentRepository) BeginTransaction() (*sqlx.Tx, error) {
	tx, err := r.DB.Beginx() // Beginx is used for sqlx transactions
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *studentRepository) SaveStudent(tx *sqlx.Tx, student entity.Student) (uuid.UUID, error) {
	query := `
		INSERT INTO students (
			student_id, student_uuid, student_first_name, student_last_name, student_gender, student_grade, parent_uuid, 
			school_uuid, created_at, created_by
		) VALUES (
			:student_id, :student_uuid, :first_name, :last_name, :student_gender, :student_grade, :parent_uuid, 
			:school_uuid, NOW(), :created_by
		)
		RETURNING student_uuid
	`

	// Using sqlx.NamedExec for named parameters
	rows, err := tx.NamedQuery(query, student)
	if err != nil {
		return uuid.Nil, err
	}
	defer rows.Close()

	var studentUUID uuid.UUID
	if rows.Next() {
		err = rows.Scan(&studentUUID)
		if err != nil {
			return uuid.Nil, err
		}
	}

	return studentUUID, nil
}