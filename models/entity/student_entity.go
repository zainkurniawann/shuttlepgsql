package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type Student struct {
	ID        int64          `db:"student_id"`
	UUID      uuid.UUID      `db:"student_uuid"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
	Grade     string         `db:"student_grade"`
	Gender     string         `db:"student_gender"`
	ParentID  sql.NullInt64  `db:"parent_id"`
	ParentUUID sql.NullString `db:"parent_uuid"`
	SchoolID  int64          `db:"school_id"`
	SchoolUUID uuid.UUID     `db:"school_uuid"`
    SchoolName string  
	CreatedAt sql.NullTime   `db:"created_at"`
	CreatedBy sql.NullString `db:"created_by"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

// type Parent struct {
// 	ID        int64          `db:"parent_id"`
// 	UUID      uuid.UUID      `db:"parent_uuid"`
// 	FirstName string         `db:"first_name"`
// 	LastName  string         `db:"last_name"`
// 	Email     string         `db:"email"`
// 	Phone     string         `db:"phone"`
// 	Address   string         `db:"address"`
// 	RoleCode  string         `db:"role_code"` // This can be "parent"
// 	CreatedAt sql.NullTime   `db:"created_at"`
// 	CreatedBy sql.NullString `db:"created_by"`
// 	UpdatedAt sql.NullTime   `db:"updated_at"`
// 	UpdatedBy sql.NullString `db:"updated_by"`
// 	DeletedAt sql.NullTime   `db:"deleted_at"`
// 	DeletedBy sql.NullString `db:"deleted_by"`
// }
