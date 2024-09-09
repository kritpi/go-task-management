package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kritpi/go-task-management/domain/exceptions"
	"github.com/kritpi/go-task-management/domain/requests"
	"github.com/kritpi/go-task-management/domain/usecases"
	"github.com/kritpi/go-task-management/utils"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type userHandler struct {
	service usecases.UserUseCase
}

func NewUserHandler(service usecases.UserUseCase) UserHandler {
	return &userHandler{
		service: service,
	}
}

func (u *userHandler) Register(c *fiber.Ctx) error {
	// Parse request
	var req *requests.UserRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	// Register user
	if err := u.service.Register(c.Context(), req); err != nil {
		switch err {
		case exceptions.ErrDuplicatedEmail:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email already registered",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func (u *userHandler) Login(c *fiber.Ctx) error {
	// Parse request
	var req *requests.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	// Login user
	user, err := u.service.Login(c.Context(), req)
	if err != nil {
		switch err {
		case exceptions.ErrLoginFailed:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Login failed",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(user)
}