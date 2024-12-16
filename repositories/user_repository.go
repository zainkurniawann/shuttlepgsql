package repositories

import (
	"database/sql"
	"fmt"
	"shuttle/models/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryInterface interface {
	// Might need to move this to a different repository
	FetchAllDriversForPermittedSchool(offset int, limit int, sortField string, sortDirection string, userUUID string) ([]entity.User, entity.School, entity.Vehicle, error)
	FetchPermittedSchoolAccess(userUUID string) (string, error)
	
	BeginTransaction() (*sqlx.Tx, error)
	FetchSpecificUser(userUUID string) (entity.User, error)
	CheckEmailExist(uuid string, email string) (bool, error)
	CheckUsernameExist(uuid string, username string) (bool, error)
	CountSuperAdmin() (int, error)
	CountSchoolAdmin() (int, error)
	CountAllDriver() (int, error)

	FetchAllSuperAdmins(offset, limit int, sortField, sortDirection string) ([]entity.User, error)
	FetchAllSchoolAdmins(offset, limit int, sortField, sortDirection string) ([]entity.User, entity.School, error)
	FetchAllDrivers(offset int, limit int, sortField string, sortDirection string) ([]entity.User, entity.School, entity.Vehicle, error)
	FetchSpecDriverFromAllSchools(userUUID string) (entity.User, entity.School, entity.Vehicle, error)

	FetchSpecSuperAdmin(userUUID string) (entity.User, error)
	FetchSpecSchoolAdmin(userUUID string) (entity.User, entity.School, error)

	FetchSuperAdminDetails(userUUID uuid.UUID) (entity.SuperAdminDetails, error)
	FetchSchoolAdminDetails(userUUID uuid.UUID) (entity.SchoolAdminDetails, error)
	FetchParentDetails(userUUID uuid.UUID) (entity.ParentDetails, error)
	FetchDriverDetails(userUUID uuid.UUID) (entity.DriverDetails, error)

	SaveUser(tx *sqlx.Tx, user entity.User) (uuid.UUID, error)
	SaveSuperAdminDetails(tx *sqlx.Tx, details entity.SuperAdminDetails, userUUID uuid.UUID, params interface{}) error
	SaveSchoolAdminDetails(tx *sqlx.Tx, details entity.SchoolAdminDetails, userUUID uuid.UUID, params interface{}) error
	SaveParentDetails(tx *sqlx.Tx, details entity.ParentDetails, userUUID uuid.UUID, params interface{}) error
	SaveDriverDetails(tx *sqlx.Tx, details entity.DriverDetails, userUUID uuid.UUID, params interface{}) error

	UpdateUser(tx *sqlx.Tx, user entity.User, userUUID string) error
	UpdateSuperAdminDetails(tx *sqlx.Tx, details entity.SuperAdminDetails, userUUID string) error
	UpdateSchoolAdminDetails(tx *sqlx.Tx, details entity.SchoolAdminDetails, userUUID string) error
	UpdateParentDetails(tx *sqlx.Tx, details entity.ParentDetails, userUUID string) error
	UpdateDriverDetails(tx *sqlx.Tx, details entity.DriverDetails, userUUID uuid.UUID) error

	DeleteSuperAdmin(tx *sqlx.Tx, userUUID uuid.UUID, user_name string) error
	DeleteSchoolAdmin(tx *sqlx.Tx, userUUID uuid.UUID, user_name string) error
	DeleteDriver(tx *sqlx.Tx, userUUID uuid.UUID, user_name string) error
}

type userRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(DB *sqlx.DB) UserRepositoryInterface {
	return &userRepository{
		DB: DB,
	}
}

func (r *userRepository) FetchAllDriversForPermittedSchool(offset int, limit int, sortField string, sortDirection string, userUUID string) ([]entity.User, entity.School, entity.Vehicle, error) {
    var users []entity.User
    var user entity.User
    var details entity.DriverDetails
    var school entity.School
    var vehicle entity.Vehicle

    query := fmt.Sprintf(`
        SELECT
            u.user_uuid, u.user_username, u.user_email, u.user_status, u.user_last_active, u.created_at, u.created_by,
            d.school_uuid, d.vehicle_uuid, d.user_picture, d.user_first_name, d.user_last_name, d.user_gender, d.user_phone, d.user_address, d.user_license_number,
            s.school_name, v.vehicle_number
        FROM users u
        LEFT JOIN driver_details d ON u.user_uuid = d.user_uuid
        LEFT JOIN schools s ON d.school_uuid = s.school_uuid
        LEFT JOIN vehicles v ON d.vehicle_uuid = v.vehicle_uuid
        WHERE u.user_role = 'driver' AND u.deleted_at IS NULL AND d.school_uuid = (SELECT school_uuid FROM school_admin_details WHERE user_uuid = $1)
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortField, sortDirection)

    rows, err := r.DB.Queryx(query, userUUID, limit, offset)
    if err != nil {
        return nil, entity.School{}, entity.Vehicle{}, err
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(
            &user.UUID, &user.Username, &user.Email, &user.Status, &user.LastActive, &user.CreatedAt, &user.CreatedBy,
            &details.SchoolUUID, &details.VehicleUUID, &details.Picture, &details.FirstName, &details.LastName,
            &details.Gender, &details.Phone, &details.Address, &details.LicenseNumber,
            &school.Name, &vehicle.VehicleNumber,
        )
        if err != nil {
            return nil, entity.School{}, entity.Vehicle{}, err
        }

        details.UserUUID = user.UUID
        user.Details = details
        users = append(users, user)
    }

    return users, school, vehicle, nil
}

func (r *userRepository) FetchPermittedSchoolAccess(userUUID string) (string, error) {
    query := `SELECT school_uuid FROM school_admin_details WHERE user_uuid = $1`
    var schoolUUID string
    err := r.DB.Get(&schoolUUID, query, userUUID)
    if err != nil {
        return "", err
    }

    return schoolUUID, nil
}



func (r *userRepository) BeginTransaction() (*sqlx.Tx, error) {
	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *userRepository) FetchSpecificUser(userUUID string) (entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE user_uuid = $1 AND deleted_at IS NULL`
	if err := r.DB.Get(&user, query, userUUID); err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) CheckEmailExist(uuid string, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(user_id) FROM users WHERE user_email = $1 AND deleted_at IS NULL`

	if uuid != "" {
		query += ` AND user_uuid != $2`
		if err := r.DB.Get(&count, query, email, uuid); err != nil {
			return false, err
		}
	} else {
		if err := r.DB.Get(&count, query, email); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (r *userRepository) CheckUsernameExist(uuid string, username string) (bool, error) {
	var count int
	query := `SELECT COUNT(user_id) FROM users WHERE user_username = $1 AND deleted_at IS NULL`
	
	if uuid != "" {
		query += ` AND user_uuid != $2`
		if err := r.DB.Get(&count, query, username, uuid); err != nil {
			return false, err
		}
	} else {
		if err := r.DB.Get(&count, query, username); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (r *userRepository) CountSuperAdmin() (int, error) {
	query := `
        SELECT COUNT(*) 
        FROM users
        WHERE user_role = 'superadmin' AND deleted_at IS NULL
    `
	var total int
	err := r.DB.Get(&total, query)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepository) CountSchoolAdmin() (int, error) {
	query := `
        SELECT COUNT(*)
        FROM users
        WHERE user_role = 'schooladmin' AND deleted_at IS NULL
    `
	var total int
	err := r.DB.Get(&total, query)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepository) CountAllDriver() (int, error) {
	query := `
		SELECT COUNT(user_id)
		FROM users
		WHERE user_role = 'driver' AND deleted_at IS NULL
	`
	var total int
	err := r.DB.Get(&total, query)
	if err != nil {
		return 0, err
	}

	return total, nil
}



func (r *userRepository) FetchAllSuperAdmins(offset int, limit int, sortField string, sortDirection string) ([]entity.User, error) {
	var users []entity.User
	var user entity.User
	var details entity.SuperAdminDetails

	query := fmt.Sprintf(`
        SELECT 
            u.user_uuid, u.user_username, u.user_email, u.user_status, 
            u.user_last_active, u.created_at, u.created_by,
            d.user_picture, d.user_first_name, d.user_last_name, d.user_gender, d.user_phone, d.user_address
        FROM users u
        LEFT JOIN super_admin_details d ON u.user_uuid = d.user_uuid
        WHERE u.user_role = 'superadmin' AND u.deleted_at IS NULL
        ORDER BY %s %s
        LIMIT $1 OFFSET $2
    `, sortField, sortDirection)

	rows, err := r.DB.Queryx(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.UUID, &user.Username, &user.Email,
			&user.Status, &user.LastActive, &user.CreatedAt, &user.CreatedBy,
			&details.Picture, &details.FirstName, &details.LastName,
			&details.Gender, &details.Phone, &details.Address,
		)
		if err != nil {
			return nil, err
		}

		details.UserUUID = user.UUID
		user.Details = details
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) FetchAllSchoolAdmins(offset int, limit int, sortField string, sortDirection string) ([]entity.User, entity.School, error) {
	var users []entity.User
	var user entity.User
	var details entity.SchoolAdminDetails
    var school entity.School

	query := fmt.Sprintf(`
        SELECT
            u.user_uuid, u.user_username, u.user_email, u.user_status,
            u.user_last_active, u.created_at, u.created_by,
            d.school_uuid, d.user_picture, d.user_first_name, d.user_last_name, d.user_gender, d.user_phone,
            s.school_name
        FROM users u
        LEFT JOIN school_admin_details d ON u.user_uuid = d.user_uuid
        LEFT JOIN schools s ON d.school_uuid = s.school_uuid
        WHERE u.user_role = 'schooladmin' AND u.deleted_at IS NULL
        ORDER BY %s %s
        LIMIT $1 OFFSET $2
    `, sortField, sortDirection)

	rows, err := r.DB.Queryx(query, limit, offset)
	if err != nil {
		return nil, entity.School{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.UUID, &user.Username, &user.Email,
			&user.Status, &user.LastActive, &user.CreatedAt, &user.CreatedBy,
			&details.SchoolUUID, &details.Picture, &details.FirstName, &details.LastName,
			&details.Gender, &details.Phone, &school.Name,
		)
		if err != nil {
            return nil, entity.School{}, err
        }

        details.UserUUID = user.UUID
        user.Details = details
        users = append(users, user)
    }

    return users, school, nil
}

func (r *userRepository) FetchAllDrivers(offset int, limit int, sortField string, sortDirection string) ([]entity.User, entity.School, entity.Vehicle, error) {
    var users []entity.User
    var user entity.User
    var details entity.DriverDetails
    var school entity.School
    var vehicle entity.Vehicle

    query := fmt.Sprintf(`
        SELECT
            u.user_uuid, u.user_username, u.user_email, u.user_status, u.user_last_active, u.created_at, u.created_by,
            d.school_uuid, d.vehicle_uuid, d.user_picture, d.user_first_name, d.user_last_name, d.user_gender, d.user_phone, d.user_address, d.user_license_number,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_name
				END,
				'N/A'
			) AS school_name,
			COALESCE(
				CASE
					WHEN v.deleted_at IS NULL THEN v.driver_uuid
				END,
				NULL
			) AS driver_uuid,
			COALESCE(
				CASE
					WHEN v.deleted_at IS NULL THEN v.vehicle_number
				END,
				'N/A'
			) AS vehicle_number,
			COALESCE(
				CASE
					WHEN v.deleted_at IS NULL THEN v.vehicle_name
				END,
				'N/A'
			) AS vehicle_name
        FROM users u
        LEFT JOIN driver_details d ON u.user_uuid = d.user_uuid
        LEFT JOIN schools s ON d.school_uuid = s.school_uuid
        LEFT JOIN vehicles v ON d.vehicle_uuid = v.vehicle_uuid
        WHERE u.user_role = 'driver' AND u.deleted_at IS NULL
        ORDER BY %s %s
        LIMIT $1 OFFSET $2
    `, sortField, sortDirection)

    rows, err := r.DB.Queryx(query, limit, offset)
    if err != nil {
        return nil, entity.School{}, entity.Vehicle{}, err
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(
            &user.UUID, &user.Username, &user.Email, &user.Status, &user.LastActive, &user.CreatedAt, &user.CreatedBy,
            &details.SchoolUUID, &details.VehicleUUID, &details.Picture, &details.FirstName, &details.LastName,
            &details.Gender, &details.Phone, &details.Address, &details.LicenseNumber,
            &school.Name, &vehicle.UUID, &vehicle.VehicleNumber, &vehicle.VehicleName,
        )
        if err != nil {
            return nil, entity.School{}, entity.Vehicle{}, err
        }

        details.UserUUID = user.UUID
        user.Details = details
        users = append(users, user)
    }

    return users, school, vehicle, nil
}

func (r *userRepository) FetchSpecDriverFromAllSchools(userUUID string) (entity.User, entity.School, entity.Vehicle, error) {
	var user entity.User
	var details entity.DriverDetails
	var school entity.School
	var vehicle entity.Vehicle

	query := `
		SELECT
			u.user_uuid, u.user_username, u.user_email, u.user_status, u.user_last_active, u.created_at, u.created_by,
			d.school_uuid, d.vehicle_uuid, d.user_picture, d.user_first_name, d.user_last_name, d.user_gender, d.user_phone, d.user_address, d.user_license_number,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_uuid
				END,
				NULL
			) AS school_uuid,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_name
				END,
				'N/A'
			) AS school_name,
			COALESCE(
				CASE
					WHEN v.deleted_at IS NULL THEN v.vehicle_uuid
				END,
				NULL
			) AS vehicle_uuid,
			COALESCE(
				CASE
					WHEN v.deleted_at IS NULL THEN v.vehicle_number
				END,
				'N/A'
			) AS vehicle_number,
			COALESCE(
				CASE
					WHEN v.deleted_at IS NULL THEN v.vehicle_name
				END,
				'N/A'
			) AS vehicle_name
		FROM users u
		LEFT JOIN driver_details d ON u.user_uuid = d.user_uuid
		LEFT JOIN schools s ON d.school_uuid = s.school_uuid
		LEFT JOIN vehicles v ON d.vehicle_uuid = v.vehicle_uuid
		WHERE u.user_role = 'driver' AND u.deleted_at IS NULL AND u.user_uuid = $1
	`

	err := r.DB.QueryRowx(query, userUUID).Scan(
		&user.UUID, &user.Username, &user.Email, &user.Status, &user.LastActive, &user.CreatedAt, &user.CreatedBy,
		&details.SchoolUUID, &details.VehicleUUID, &details.Picture, &details.FirstName, &details.LastName,
		&details.Gender, &details.Phone, &details.Address, &details.LicenseNumber,
		&school.UUID, &school.Name, &vehicle.UUID, &vehicle.VehicleNumber, &vehicle.VehicleName,
	)
	
	if err != nil {
		return user, school, vehicle, err
	}

	details.UserUUID = user.UUID
	user.Details = details
	return user, school, vehicle, nil
}



func (r *userRepository) FetchSpecSuperAdmin(userUUID string) (entity.User, error) {
	var user entity.User
	var details entity.SuperAdminDetails

	query := `
        SELECT
            u.user_uuid, u.user_username, u.user_email, u.user_status,
            u.user_last_active, u.created_at, u.created_by, u.updated_at, u.updated_by, u.deleted_at, u.deleted_by,
            d.user_picture, d.user_first_name, d.user_last_name, d.user_gender, d.user_phone, d.user_address
        FROM users u
        LEFT JOIN super_admin_details d ON u.user_uuid = d.user_uuid
        WHERE u.user_uuid = $1 AND u.user_role = 'superadmin' AND u.deleted_at IS NULL
    `

	err := r.DB.QueryRowx(query, userUUID).Scan(
		&user.UUID, &user.Username, &user.Email,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.CreatedBy,
		&user.UpdatedAt, &user.UpdatedBy, &user.DeletedAt, &user.DeletedBy,
		&details.Picture, &details.FirstName, &details.LastName,
		&details.Gender, &details.Phone, &details.Address,
	)

	if err != nil {
		return user, err
	}

	details.UserUUID = user.UUID
	user.Details = details
	return user, nil
}

func (r *userRepository) FetchSpecSchoolAdmin(userUUID string) (entity.User, entity.School, error) {
	var user entity.User
	var details entity.SchoolAdminDetails
    var school entity.School

	query := `
        SELECT
            u.user_uuid, u.user_username, u.user_email, u.user_status,
            u.user_last_active, u.created_at, u.created_by, u.updated_at, u.updated_by, u.deleted_at, u.deleted_by,
            d.school_uuid, d.user_picture, d.user_first_name, d.user_last_name, d.user_gender, d.user_phone, d.user_address,
            s.school_name
        FROM users u
        LEFT JOIN school_admin_details d ON u.user_uuid = d.user_uuid
        LEFT JOIN schools s ON d.school_uuid = s.school_uuid
        WHERE u.user_uuid = $1 AND u.user_role = 'schooladmin' AND u.deleted_at IS NULL
    `
	err := r.DB.QueryRowx(query, userUUID).Scan(
		&user.UUID, &user.Username, &user.Email,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.CreatedBy,
		&user.UpdatedAt, &user.UpdatedBy, &user.DeletedAt, &user.DeletedBy,
		&details.SchoolUUID, &details.Picture, &details.FirstName, &details.LastName,
		&details.Gender, &details.Phone, &details.Address, &school.Name,
	)

    if err != nil {
        return user, school, err
    }

    details.UserUUID = user.UUID
    user.Details = details
    return user, school, nil
}



func (r *userRepository) FetchSuperAdminDetails(userUUID uuid.UUID) (entity.SuperAdminDetails, error) {
	var superAdminDetails entity.SuperAdminDetails
	query := `SELECT user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address
			  FROM super_admin_details WHERE user_uuid = $1`
	if err := r.DB.Get(&superAdminDetails, query, userUUID); err != nil {
		return superAdminDetails, err
	}

	return superAdminDetails, nil
}

func (r *userRepository) FetchSchoolAdminDetails(userUUID uuid.UUID) (entity.SchoolAdminDetails, error) {
	var schoolAdminDetails entity.SchoolAdminDetails
	query := `SELECT school_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address
			  FROM school_admin_details WHERE user_uuid = $1`
	if err := r.DB.Get(&schoolAdminDetails, query, userUUID); err != nil {
		return schoolAdminDetails, err
	}

	return schoolAdminDetails, nil
}

func (r *userRepository) FetchParentDetails(userUUID uuid.UUID) (entity.ParentDetails, error) {
	var parentDetails entity.ParentDetails
	query := `SELECT user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address
			  FROM parent_details WHERE user_uuid = $1`
	if err := r.DB.Get(&parentDetails, query, userUUID); err != nil {
		return parentDetails, err
	}

	return parentDetails, nil
}

func (r *userRepository) FetchDriverDetails(userUUID uuid.UUID) (entity.DriverDetails, error) {
	var driverDetails entity.DriverDetails
	query := `SELECT school_uuid, school_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address, user_license_number
			  FROM driver_details WHERE user_uuid = $1`
	if err := r.DB.Get(&driverDetails, query, userUUID); err != nil {
		return driverDetails, err
	}

	return driverDetails, nil
}



func (r *userRepository) SaveUser(tx *sqlx.Tx, userEntity entity.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (user_id, user_uuid, user_username, user_email, user_password, user_role, user_role_code, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING user_uuid`
	var userUUID uuid.UUID
	err := tx.QueryRow(query, userEntity.ID, userEntity.UUID, userEntity.Username, userEntity.Email, userEntity.Password, userEntity.Role, userEntity.RoleCode, userEntity.CreatedBy).Scan(&userUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

func (r *userRepository) SaveSuperAdminDetails(tx *sqlx.Tx, details entity.SuperAdminDetails, userUUID uuid.UUID, params interface{}) error {
	details.UserUUID = userUUID
	query := `
        INSERT INTO super_admin_details 
        (user_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address) 
        VALUES (:user_uuid, :user_picture, :user_first_name, :user_last_name, :user_gender, :user_phone, :user_address)
    `
	params = details
	_, err := tx.NamedExec(query, params)

	return err
}

func (r *userRepository) SaveSchoolAdminDetails(tx *sqlx.Tx, details entity.SchoolAdminDetails, userUUID uuid.UUID, params interface{}) error {
	details.UserUUID = userUUID
	query := `
        INSERT INTO school_admin_details 
        (user_uuid, school_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address) 
        VALUES (:user_uuid, :school_uuid, :user_picture, :user_first_name, :user_last_name, :user_gender, :user_phone, :user_address)
    `
	params = details
	_, err := tx.NamedExec(query, params)
	return err
}

func (r *userRepository) SaveParentDetails(tx *sqlx.Tx, details entity.ParentDetails, userUUID uuid.UUID, params interface{}) error {
	details.UserUUID = userUUID
	query := `
        INSERT INTO parent_details 
        (user_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address) 
        VALUES (:user_uuid, :user_picture, :user_first_name, :user_last_name, :user_gender, :user_phone, :user_address)
    `
	params = details
	_, err := tx.NamedExec(query, params)
	return err
}

func (r *userRepository) SaveDriverDetails(tx *sqlx.Tx, details entity.DriverDetails, userUUID uuid.UUID, params interface{}) error {
	details.UserUUID = userUUID

	if details.SchoolUUID == nil || *details.SchoolUUID == uuid.Nil {
		details.SchoolUUID = nil
	}
	if details.VehicleUUID == nil || *details.VehicleUUID == uuid.Nil {
		details.VehicleUUID = nil
	}

	query := `
		INSERT INTO driver_details 
		(user_uuid, school_uuid, vehicle_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address, user_license_number) 
		VALUES (:user_uuid, :school_uuid, :vehicle_uuid, :user_picture, :user_first_name, :user_last_name, :user_gender, :user_phone, :user_address, :user_license_number)
	`
	params = details
	_, err := tx.NamedExec(query, params)
	if err != nil {
		return err
	}

	if details.VehicleUUID != nil {
		return r.UpdateDriverUUIDInVehicles(tx, userUUID, *details.VehicleUUID)
	}
	return nil
}

func (r *userRepository) UpdateDriverUUIDInVehicles(tx *sqlx.Tx, userUUID uuid.UUID, vehicleUUID uuid.UUID) error {
	var userUUIDParam interface{}
    if userUUID == uuid.Nil {
        userUUIDParam = nil
    } else {
        userUUIDParam = userUUID
    }

    query := `
        UPDATE vehicles
        SET driver_uuid = $1
        WHERE vehicle_uuid = $2
		`
    _, err := tx.Exec(query, userUUIDParam, vehicleUUID)
    return err
}

func (r *userRepository) UpdateUser(tx *sqlx.Tx, user entity.User, userUUID string) error {
	query := `
        UPDATE users
        SET user_username = $1, user_email = $2, user_role = $3, updated_at = NOW(), updated_by = $4
        WHERE user_uuid = $5`
	_, err := tx.Exec(query, user.Username, user.Email, user.Role, user.UpdatedBy, userUUID)
	return err
}

func (r *userRepository) UpdateSuperAdminDetails(tx *sqlx.Tx, details entity.SuperAdminDetails, userUUID string) error {
	query := `
        UPDATE super_admin_details
        SET user_picture = $1, user_first_name = $2, user_last_name = $3, user_gender = $4, user_phone = $5, user_address = $6
        WHERE user_uuid = $7`
	res, err := tx.Exec(query, details.Picture, details.FirstName, details.LastName, details.Gender, details.Phone, details.Address, userUUID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *userRepository) UpdateSchoolAdminDetails(tx *sqlx.Tx, details entity.SchoolAdminDetails, userUUID string) error {
	query := `
        UPDATE school_admin_details
        SET school_uuid = $1, user_picture = $2, user_first_name = $3, user_last_name = $4, user_gender = $5, user_phone = $6, user_address = $7
        WHERE user_uuid = $8`
	res, err := tx.Exec(query, details.SchoolUUID, details.Picture, details.FirstName, details.LastName, details.Gender, details.Phone, details.Address, userUUID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *userRepository) UpdateParentDetails(tx *sqlx.Tx, details entity.ParentDetails, userUUID string) error {
	query := `
        UPDATE parent_details
        SET user_picture = :user_picture, user_first_name = :user_first_name, user_last_name = :user_last_name, user_gender = :user_gender, user_phone = :user_phone, user_address = :user_address
        WHERE user_uuid = :user_uuid`
	res, err := tx.NamedExec(query, details)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *userRepository) UpdateDriverDetails(tx *sqlx.Tx, details entity.DriverDetails, userUUID uuid.UUID) error {
	details.UserUUID = userUUID

    if details.SchoolUUID == nil || *details.SchoolUUID == uuid.Nil {
        details.SchoolUUID = nil
    }
    if details.VehicleUUID == nil || *details.VehicleUUID == uuid.Nil {
        details.VehicleUUID = nil
    }

    var currentVehicleUUID *uuid.UUID
    if details.VehicleUUID == nil {
        err := tx.Get(&currentVehicleUUID, `SELECT vehicle_uuid FROM driver_details WHERE user_uuid = $1`, userUUID)
        if err != nil && err != sql.ErrNoRows {
            return err
        }
    }
		
	query := `
        UPDATE driver_details
        SET school_uuid = $1, vehicle_uuid = $2, user_first_name = $3, user_last_name = $4,
		user_gender = $5, user_phone = $6, user_address = $7, user_license_number = $8
		WHERE user_uuid = $9`
	res, err := tx.Exec(query, details.SchoolUUID, details.VehicleUUID, details.FirstName, details.LastName, details.Gender, details.Phone, details.Address, details.LicenseNumber, details.UserUUID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	
	if details.VehicleUUID != nil {
        return r.UpdateDriverUUIDInVehicles(tx, userUUID, *details.VehicleUUID)
    }

    if currentVehicleUUID != nil {
        return r.UpdateDriverUUIDInVehicles(tx, uuid.Nil, *currentVehicleUUID)
    }
	
	return nil
}



func (r *userRepository) DeleteSuperAdmin(tx *sqlx.Tx, userUUID uuid.UUID, user_name string) error {
	query := `UPDATE users SET deleted_at = NOW(), deleted_by = $1 WHERE user_uuid = $2`
	res, err := tx.Exec(query, user_name, userUUID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) DeleteSchoolAdmin(tx *sqlx.Tx, userUUID uuid.UUID, user_name string) error {
	query := `UPDATE users SET deleted_at = NOW(), deleted_by = $1 WHERE user_uuid = $2`
	res, err := tx.Exec(query, user_name, userUUID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) DeleteDriver(tx *sqlx.Tx, userUUID uuid.UUID, user_name string) error {
	query := `UPDATE users SET deleted_at = NOW(), deleted_by = $1 WHERE user_uuid = $2`
	res, err := tx.Exec(query, user_name, userUUID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}