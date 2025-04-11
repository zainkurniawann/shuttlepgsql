package routes

import (
	"shuttle/handler"
	"shuttle/middleware"
	"shuttle/repositories"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Route(r *fiber.App, db *sqlx.DB) {
	authRepository := repositories.NewAuthRepository(db)
	userRepository := repositories.NewUserRepository(db)
	schoolRepository := repositories.NewSchoolRepository(db)
	vehicleRepository := repositories.NewVehicleRepository(db)
	studentRepository := repositories.NewStudentRepository(db)
	childernRepository := repositories.NewChildernRepository(db)
	shuttleRepository := repositories.NewShuttleRepository(db)

	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(authRepository, userRepository)
	schoolService := services.NewSchoolService(schoolRepository, userRepository)
	vehicleService := services.NewVehicleService(vehicleRepository)
	studentService := services.NewStudentService(studentRepository, userRepository)
	childernService := services.NewChildernService(childernRepository)
	shuttleService := services.NewShuttleService(shuttleRepository)

	authHandler := handler.NewAuthHttpHandler(authService)
	userHandler := handler.NewUserHttpHandler(userService, schoolService)
	schoolHandler := handler.NewSchoolHttpHandler(schoolService)
	vehicleHandler := handler.NewVehicleHttpHandler(vehicleService)
	studentHandler := handler.NewStudentHttpHandler(studentService)
	childernHandler := handler.NewChildernHandler(childernService)
	shuttleHandler := handler.NewShuttleHandler(shuttleService)

	wsService := utils.NewWebSocketService(userRepository, authRepository)

	// FOR PUBLIC
	r.Post("login", authHandler.Login)
	r.Post("/refresh-token", authHandler.IssueNewAccessToken)
	r.Static("/assets", "./assets")

	r.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	r.Get("/ws/:id", websocket.New(wsService.HandleWebSocketConnection))

	// FOR AUTHENTICATED
	protected := r.Group("/api")
	protected.Use(middleware.AuthenticationMiddleware())
	protected.Use(middleware.AuthorizationMiddleware([]string{"SA", "AS", "D", "P"}))

	protected.Get("/my/profile", authHandler.GetMyProfile)
	protected.Post("/logout", authHandler.Logout)

	protectedSuperAdmin := protected.Group("/superadmin")
	protectedSuperAdmin.Use(middleware.AuthorizationMiddleware([]string{"SA"}))

	protectedSchoolAdmin := protected.Group("/school")
	protectedSchoolAdmin.Use(middleware.AuthorizationMiddleware([]string{"AS"}))
	protectedSchoolAdmin.Use(middleware.SchoolAdminMiddleware(userService))

	protectedParent := protected.Group("/parent")
	protectedParent.Use(middleware.AuthorizationMiddleware([]string{"P"}))

	protectedDriver := protected.Group("/driver")
	protectedDriver.Use(middleware.AuthorizationMiddleware([]string{"D"}))

	// USER FOR SUPERADMIN
	protectedSuperAdmin.Get("/user/sa/all", userHandler.GetAllSuperAdmin)
	protectedSuperAdmin.Get("/user/as/all", userHandler.GetAllSchoolAdmin)
	protectedSuperAdmin.Get("/user/driver/all", userHandler.GetAllPermittedDriver)
	protectedSuperAdmin.Get("/user/sa/:id", userHandler.GetSpecSuperAdmin)
	protectedSuperAdmin.Get("/user/as/:id", userHandler.GetSpecSchoolAdmin)
	protectedSuperAdmin.Get("/user/driver/:id", userHandler.GetSpecPermittedDriver)
	protectedSuperAdmin.Post("/user/add", userHandler.AddUser)
	protectedSuperAdmin.Put("/user/update/:id", userHandler.UpdateUser)
	protectedSuperAdmin.Delete("/user/sa/delete/:id", userHandler.DeleteSuperAdmin)
	protectedSuperAdmin.Delete("/user/as/delete/:id", userHandler.DeleteSchoolAdmin)
	protectedSuperAdmin.Delete("/user/driver/delete/:id", userHandler.DeleteDriver)

	// SCHOOL FOR SUPERADMIN
	protectedSuperAdmin.Get("/school/all", schoolHandler.GetAllSchools)
	protectedSuperAdmin.Get("/school/:id", schoolHandler.GetSpecSchool)
	protectedSuperAdmin.Post("/school/add", schoolHandler.AddSchool)
	protectedSuperAdmin.Put("/school/update/:id", schoolHandler.UpdateSchool)
	protectedSuperAdmin.Delete("/school/delete/:id", schoolHandler.DeleteSchool)

	// VEHICLE FOR SUPERADMIN
	protectedSuperAdmin.Get("/vehicle/all", vehicleHandler.GetAllVehicles)
	protectedSuperAdmin.Get("/vehicle/:id", vehicleHandler.GetSpecVehicle)
	protectedSuperAdmin.Post("/vehicle/add", vehicleHandler.AddVehicle)
	protectedSuperAdmin.Put("/vehicle/update/:id", vehicleHandler.UpdateVehicle)
	protectedSuperAdmin.Delete("/vehicle/delete/:id", vehicleHandler.DeleteVehicle)

	////////////////////////////////////// SCHOOL ADMIN //////////////////////////////////////

	protectedSchoolAdmin.Get("/user/driver/all", userHandler.GetAllPermittedDriver)
	// protectedSchoolAdmin.Post("/user/driver/add", handler.AddSchoolDriver)
	// protectedSchoolAdmin.Put("/user/driver/update/:id", handler.UpdateSchoolDriver)
	//protectedSchoolAdmin.Delete("/user/driver/delete/:id", handler.DeleteSchoolDriver)

	// protectedSchoolAdmin.Get("/student/all", handler.GetAllStudentWithParents)
	protectedSchoolAdmin.Post("/student/add", studentHandler.AddStudentWithParent)
	// protectedSchoolAdmin.Put("/student/update/:id", handler.UpdateSchoolStudentWithParents)
	// protectedSchoolAdmin.Delete("/student/delete/:id", handler.DeleteSchoolStudentWithParents)

	protectedSchoolAdmin.Get("/route/all", handler.GetAllRoutes)
	protectedSchoolAdmin.Get("/route/:id", handler.GetSpecRoute)
	protectedSchoolAdmin.Post("/route/add", handler.AddRoute)

	///////////////////////////////// PARENT ///////////////////////////////////

	protectedParent.Get("/my/childern/all", childernHandler.GetAllChilderns)
	protectedParent.Get("/my/childern/:id", childernHandler.GetSpecChildern)
	protectedParent.Put("/my/childern/update/:id", childernHandler.UpdateChildern)
	protectedParent.Get("/my/childern/shuttle/inf", shuttleHandler.GetShuttleStatusByParent)

	////////////////////////////// DRIVER😂 /////////////////////////////////////

	protectedDriver.Post("/shuttle/add", shuttleHandler.AddShuttle)
	protectedDriver.Put("/shuttle/update/:id", shuttleHandler.EditShuttle)

}
