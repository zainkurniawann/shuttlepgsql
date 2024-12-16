package dto

type ShuttleRequest struct {
	StudentUUID string `json:"student_uuid" validate:"required,uuid4"`
}

type ShuttleResponse struct {
	StudentName string `db:"student_name"`
	DriverName  string `db:"driver_name"`
	Status      string `db:"status"`
	CreatedAt   string `db:"created_at"` // sudah dalam format string
}

