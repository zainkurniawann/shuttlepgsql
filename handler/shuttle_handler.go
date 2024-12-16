package handler

import (
	// "log"
	"log"
	"net/http"
	// "shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShuttleHandler struct {
	ShuttleService services.ShuttleServiceInterface
	DB             *sqlx.DB // Add a DB field to the handler
}

func NewShuttleHandler(shuttleService services.ShuttleServiceInterface) *ShuttleHandler {
	return &ShuttleHandler{
		ShuttleService: shuttleService,
	}
}

func (h *ShuttleHandler) GetShuttleStatusByParent(c *fiber.Ctx) error {
	// Ambil userUUID dari token
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok || userUUID == "" {
		return utils.BadRequestResponse(c, "Invalid or missing userUUID", nil)
	}
	log.Printf("awokwawk",userUUID)
	// Parse userUUID ke uuid.UUID
	parentUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid userUUID format", nil)
	}

	log.Println(parentUUID)

	// Panggil service untuk mendapatkan data shuttle
	shuttles, err := h.ShuttleService.GetShuttleStatusByParent(parentUUID)
	if err != nil {
		return utils.NotFoundResponse(c, "Shuttle data not found", nil)
	}
	log.Println("sutels",shuttles)
	// Return response
	return c.Status(http.StatusOK).JSON(shuttles)
}

// shuttleHandler.go
// func (h *ShuttleHandler) AddShuttle(c *fiber.Ctx) error {
// 	// Parse request body into ShuttleRequest DTO
// 	var request dto.ShuttleRequest
// 	if err := c.BodyParser(&request); err != nil {
// 		return utils.BadRequestResponse(c, "Invalid request body", nil)
// 	}

// 	// Parse student_uuid from request body
// 	studentUUID, err := uuid.Parse(request.StudentUUID)
// 	if err != nil {
// 		return utils.BadRequestResponse(c, "Invalid studentUUID format", nil)
// 	}

// 	// Call service to add shuttle for the student
// 	shuttle, err := h.ShuttleService.AddShuttle(studentUUID)
// 	if err != nil {
// 		return utils.InternalServerErrorResponse(c, "Failed to add shuttle", nil)
// 	}

// 	// Return successful response
// 	return c.Status(http.StatusOK).JSON(shuttle)
// }
