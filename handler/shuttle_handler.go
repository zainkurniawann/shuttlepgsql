package handler

import (
	// "log"
	"log"
	"net/http"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strings"
	"errors"
	"database/sql"

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

func (h *ShuttleHandler) AddShuttle(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok || userUUID == "" {
		return utils.BadRequestResponse(c, "Invalid or missing userUUID", nil)
	}
	username := c.Locals("user_name").(string)
	log.Println("Username:", username)
	driverUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid userUUID format", nil)
	}
	log.Println("Driver UUID:", driverUUID)
	shuttleReq := new(dto.ShuttleRequest)
	if err := c.BodyParser(shuttleReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	// Set default status if empty
	if shuttleReq.Status == "" {
		shuttleReq.Status = "menunggu dijemput"
	}

	log.Println("Adding shuttle with data:", shuttleReq)
	if err := utils.ValidateStruct(c, shuttleReq); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}
	log.Println("woi", shuttleReq)
	if err := h.ShuttleService.AddShuttle(*shuttleReq, driverUUID.String(), username); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to add shuttle", nil)
	}

	return utils.SuccessResponse(c, "Shuttle added successfully", nil)
}

func (h *ShuttleHandler) EditShuttle(c *fiber.Ctx) error {
	// Ambil shuttleUUID dari parameter
	id := c.Params("id")
	if id == "" {
		return utils.BadRequestResponse(c, "Missing shuttleUUID in URL", nil)
	}

	// Parse request body untuk status
	var statusReq struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.BodyParser(&statusReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", nil)
	}

	// Validasi status
	if err := utils.ValidateStruct(c, statusReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid status: "+err.Error(), nil)
	}

	log.Println("Editing shuttle:", id, "with status:", statusReq.Status)

	// Panggil service untuk update
	if err := h.ShuttleService.EditShuttleStatus(id, statusReq.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundResponse(c, "Shuttle not found", nil)
		}
		return utils.InternalServerErrorResponse(c, "Failed to edit shuttle", nil)
	}

	return utils.SuccessResponse(c, "Shuttle status updated successfully", nil)
}