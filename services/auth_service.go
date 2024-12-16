package services

import (
	"context"
	"path/filepath"
	"time"

	"shuttle/errors"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/repositories"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Login(email, password string) (userDataa dto.UserDataOnLoginDTO, err error)
	GetMyProfile(userUUID, roleCode string) (interface{}, error)
	CheckStoredRefreshToken(userID string, refreshToken string) error
	DeleteRefreshTokenOnLogout(ctx context.Context, userID string) error
	UpdateUserStatus(userUUID, status string, lastActive time.Time) error
	UpdateRefreshToken(userUUID, refreshToken string) error
}

type AuthService struct {
	authRepository repositories.AuthRepositoryInterface
	userRepository repositories.UserRepositoryInterface
}

func NewAuthService(authRepository repositories.AuthRepositoryInterface, userRepository repositories.UserRepositoryInterface) AuthService {
	return AuthService{
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (service AuthService) Login(email, password string) (userData dto.UserDataOnLoginDTO, err error) {
	user, err := service.authRepository.Login(email)
	if err != nil {
		logger.LogError(err, "Failed to login", map[string]interface{}{
			"email": email,
		})
		return dto.UserDataOnLoginDTO{}, errors.New("invalid email or password", 0)
	}

	userDataOnLogin := dto.UserDataOnLoginDTO{
		UserID:    user.ID,
		UserUUID:  user.UUID,
		Username:  user.Username,
		RoleCode:  user.RoleCode,
		Password:  user.Password,
	}

	if !validatePassword(password, userDataOnLogin.Password) {
		return dto.UserDataOnLoginDTO{}, errors.New("invalid email or password", 0)
	}

	return userDataOnLogin, nil
}

func (service *AuthService) GetMyProfile(userUUID, roleCode string) (interface{}, error) {

	user, err := service.userRepository.FetchSpecificUser(userUUID)
	if err != nil {
		return nil, err
	}

	parsedUserUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, errors.New("invalid user UUID format", 0)
	}

	switch user.RoleCode {
	case "SA":
		superAdminDetails, err := service.userRepository.FetchSuperAdminDetails(parsedUserUUID)
		if err != nil {
			return nil, err
		}

		picture := superAdminDetails.Picture
		if picture != "" {
			imageURL, err := generateImageURL(picture)
			if err != nil {
				return nil, err
			}
			superAdminDetails.Picture = imageURL
		}

		details := dto.SuperAdminDetailsResponseDTO{
			Picture:   superAdminDetails.Picture,
			FirstName: superAdminDetails.FirstName,
			LastName:  superAdminDetails.LastName,
			Gender:    dto.Gender(superAdminDetails.Gender),
			Phone:     superAdminDetails.Phone,
			Address:   superAdminDetails.Address,
		}

		result := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Role:       dto.Role(user.Role),
			RoleCode:   user.RoleCode,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			Details:    details,
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}

		return result, nil
	case "AS":
		schoolAdminDetails, err := service.userRepository.FetchSchoolAdminDetails(parsedUserUUID)
		if err != nil {
			return nil, err
		}

		picture := schoolAdminDetails.Picture
		if picture != "" {
			imageURL, err := generateImageURL(picture)
			if err != nil {
				return nil, err
			}
			schoolAdminDetails.Picture = imageURL
		}

		details := dto.SchoolAdminDetailsResponseDTO{
			FirstName: schoolAdminDetails.FirstName,
			LastName:  schoolAdminDetails.LastName,
			Gender:    dto.Gender(schoolAdminDetails.Gender),
			Phone:     schoolAdminDetails.Phone,
			Address:   schoolAdminDetails.Address,
		}

		result := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Role:       dto.Role(user.Role),
			RoleCode:   user.RoleCode,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			Details:    details,
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}

		return result, nil
	case "P":
		parentDetails, err := service.userRepository.FetchParentDetails(parsedUserUUID)
		if err != nil {
			return nil, err
		}

		picture := parentDetails.Picture
		if picture != "" {
			imageURL, err := generateImageURL(picture)
			if err != nil {
				return nil, err
			}
			parentDetails.Picture = imageURL
		}

		details := dto.ParentDetailsResponseDTO{
			FirstName: parentDetails.FirstName,
			LastName:  parentDetails.LastName,
			Gender:    dto.Gender(parentDetails.Gender),
			Phone:     parentDetails.Phone,
			Address:   parentDetails.Address,
		}

		result := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Role:       dto.Role(user.Role),
			RoleCode:   user.RoleCode,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			Details:    details,
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}

		return result, nil
	case "D":
		driverDetails, err := service.userRepository.FetchDriverDetails(parsedUserUUID)
		if err != nil {
			return nil, err
		}

		picture := driverDetails.Picture
		if picture != "" {
			imageURL, err := generateImageURL(picture)
			if err != nil {
				return nil, err
			}
			driverDetails.Picture = imageURL
		}

		details := dto.DriverDetailsResponseDTO{
			FirstName: driverDetails.FirstName,
			LastName:  driverDetails.LastName,
			Gender:    dto.Gender(driverDetails.Gender),
			Phone:     driverDetails.Phone,
			Address:   driverDetails.Address,
		}

		result := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Role:       dto.Role(user.Role),
			RoleCode:   user.RoleCode,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			Details:    details,
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}

		return result, nil
	default:
		return nil, errors.New("invalid role code", 0)
	}
}

func (service *AuthService) CheckStoredRefreshToken(userUUID string, refreshToken string) error {
	refreshTokenEntity, err := service.authRepository.CheckRefreshTokenData(userUUID, refreshToken)
	if err != nil {
		return err
	}

	if refreshTokenEntity.RefreshToken != refreshToken {
		return errors.New("invalid refresh token", 0)
	}

	if refreshTokenEntity.ExpiredAt.Before(time.Now()) {
		return errors.New("refresh token has expired", 0)
	}

	return nil
}

func (service *AuthService) DeleteRefreshTokenOnLogout(ctx context.Context, userUUID string) error {
	err := service.authRepository.DeleteRefreshToken(ctx, userUUID)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) UpdateUserStatus(userUUID, status string, lastActive time.Time) error {
	err := service.authRepository.UpdateUserStatus(userUUID, status, lastActive)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) UpdateRefreshToken(userUUID, refreshToken string) error {
	tokendata, err := service.authRepository.CheckRefreshTokenData(userUUID, refreshToken)
	if err != nil {
		return err
	}

	if tokendata.LastUsedAt != nil && time.Since(*tokendata.LastUsedAt) > time.Hour {
		_, err := service.authRepository.UpdateRefreshToken(userUUID, refreshToken)
		if err != nil {
			return err
		}
	} else {
		return errors.New("cannot reissue a new access token yet", 0)
	}

	return nil
}

func generateImageURL(imagePath string) (string, error) {
	fileName := filepath.Base(imagePath)
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}

	ext := filepath.Ext(fileName)
	if !contains(allowedExtensions, ext) {
		return "", nil
	}

	baseURL := "http://" + viper.GetString("BASE_URL") + "/assets/images/"
	return baseURL + fileName, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func validatePassword(providedPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	return err == nil
}