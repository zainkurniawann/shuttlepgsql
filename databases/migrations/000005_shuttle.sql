-- +goose Up
-- +goose StatementBegin
CREATE TYPE shuttle_status AS ENUM ('di rumah', 'menunggu dijemput', 'menuju sekolah', 'di sekolah', 'menuju rumah');

CREATE TABLE shuttle (
    shuttle_id BIGINT PRIMARY KEY,
    shuttle_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    student_uuid UUID NOT NULL REFERENCES students(student_uuid) ON DELETE SET NULL,
    driver_uuid UUID NOT NULL REFERENCES users(user_uuid) ON DELETE SET NULL,
    status shuttle_status NOT NULL DEFAULT 'di rumah',
    student_pickup_point JSON,
    student_destination_name VARCHAR (50) NOT NULL,
    student_destination_point JSON,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shuttle;
DROP TYPE IF EXISTS shuttle_status;
-- +goose StatementEnd
