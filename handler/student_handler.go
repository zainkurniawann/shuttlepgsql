package handler

import (
	// "fmt"

	// "strconv"
	"strings"

	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
)

type StudentHandlerInterface interface {
	AddStudentWithParent(c *fiber.Ctx) error
}

type studentHandler struct {
	studentService services.StudentService
}

func NewStudentHttpHandler(studentService services.StudentService) StudentHandlerInterface {
	return &studentHandler{
		studentService: studentService,
	}
}

func (handler *studentHandler) AddStudentWithParent(c *fiber.Ctx) error {
	// Mengambil username dari context (asumsi dari middleware)
	createdBy := c.Locals("user_name").(string)

	// Parsing body request
	req := new(dto.AddStudentWithParentRequestDTO)
	if err := c.BodyParser(req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	// Validasi data menggunakan utils
	if err := utils.ValidateStruct(c, req); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	// Memanggil service untuk menambahkan student dan parent
	studentUUID, parentUUID, err := handler.studentService.AddPermittedSchoolStudentWithParents(*req, createdBy)
	if err != nil {
		logger.LogError(err, "Failed to add student and parent", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	// Menyusun response sukses
	response := fiber.Map{
		"message":      "Student and parent added successfully",
		"student_uuid": studentUUID.String(),
		"parent_uuid":  parentUUID.String(),
	}
	return utils.SuccessResponse(c, "Student and parent added successfully", response)
}











// import (
// 	"shuttle/errors"
// 	"shuttle/logger"
// 	"shuttle/models"
// 	"shuttle/services"
// 	"shuttle/utils"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// func GetAllStudentWithParents(c *fiber.Ctx) error {
// 	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
// 	if err != nil {
// 		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
// 			"school_id": c.Locals("schoolId"),
// 		})
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	students, err := services.GetAllPermitedSchoolStudentsWithParents(SchoolObjID)
// 	if err != nil {
// 		logger.LogError(err, "Failed to fetch all students", nil)
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	return c.Status(fiber.StatusOK).JSON(students)
// }

// func AddSchoolStudentWithParents(c *fiber.Ctx) error {
// 	username := c.Locals("username").(string)
// 	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
// 	if err != nil {
// 		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
// 			"school_id": c.Locals("schoolId"),
// 		})
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}
	
// 	student := new(models.SchoolStudentRequest)
// 	if err := c.BodyParser(student); err != nil {
// 		return utils.BadRequestResponse(c, "Invalid request data", nil)
// 	}

// 	if err := utils.ValidateStruct(c, student); err != nil {
// 		return utils.BadRequestResponse(c, err.Error(), nil)
// 	}

// 	if (models.User{}) == student.Parent {
// 		return utils.BadRequestResponse(c, "Parent details are required", nil)
// 	}

// 	if student.Parent.Phone == "" || student.Parent.Address == "" || student.Parent.Email == "" {
// 		return utils.BadRequestResponse(c, "Parent details are required", nil)
// 	}

// 	if err := services.AddPermittedSchoolStudentWithParents(*student, SchoolObjID, username); err != nil {
// 		if customErr, ok := err.(*errors.CustomError); ok {
// 			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
// 		}
// 		logger.LogError(err, "Failed to add student", nil)
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	return utils.SuccessResponse(c, "Student created successfully", nil)
// }

// func UpdateSchoolStudentWithParents(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
// 	if err != nil {
// 		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
// 			"school_id": c.Locals("schoolId"),
// 		})
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	student := new(models.SchoolStudentRequest)
// 	if err := c.BodyParser(student); err != nil {
// 		return utils.BadRequestResponse(c, "Invalid request data", nil)
// 	}

// 	if err := utils.ValidateStruct(c, student); err != nil {
// 		return utils.BadRequestResponse(c, err.Error(), nil)
// 	}

// 	if err := services.UpdatePermittedSchoolStudentWithParents(id, *student, SchoolObjID); err != nil {
// 		if customErr, ok := err.(*errors.CustomError); ok {
// 			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
// 		}
// 		logger.LogError(err, "Failed to update student", nil)
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	return utils.SuccessResponse(c, "Student updated successfully", nil)
// }

// func DeleteSchoolStudentWithParents(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
// 	if err != nil {
// 		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
// 			"school_id": c.Locals("schoolId"),
// 		})
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	if err := services.DeletePermittedSchoolStudentWithParents(id, SchoolObjID); err != nil {
// 		if customErr, ok := err.(*errors.CustomError); ok {
// 			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
// 		}
// 		logger.LogError(err, "Failed to delete student", nil)
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	return utils.SuccessResponse(c, "Student deleted successfully", nil)
// }