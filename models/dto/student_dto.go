package dto

type AddStudentWithParentRequestDTO struct {
	Student StudentRequestDTO `json:"student"` // Information about the student
	Parent  UserRequestsDTO   `json:"parent"`  // Information about the parent
}

type StudentRequestDTO struct {
	FirstName  string `json:"first_name" validate:"required,max=255"`
	LastName   string `json:"last_name" validate:"required,max=255"`
	Grade      string `json:"grade" validate:"required,max=50"`
	Gender     string `json:"gender" validate:"required,max=50"`
	ParentUUID string `json:"parent_uuid,omitempty" validate:"omitempty,uuid4"` // For linking existing parent
	SchoolUUID string `json:"school_uuid" validate:"required,uuid4"`
}

type StudentResponseDTO struct {
	UUID          string `json:"student_uuid"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Gender        string `json:"gender" validate:"required,max=50"`
	Grade         string `json:"grade" validate:"required,max=50"`
	ParentUUID    string `json:"parent_uuid,omitempty"`
	SchoolUUID    string `json:"school_uuid"`
	SchoolName    string `json:"school_name,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	CreatedBy     string `json:"created_by,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	UpdatedBy     string `json:"updated_by,omitempty"`
}

type ParentRequestDTO struct {
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"` // New field
	Phone     string `json:"phone" validate:"required,phone"`
	Address   string `json:"address" validate:"required,max=255"`
}

type ParentResponseDTO struct {
	UUID      string `json:"parent_uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
