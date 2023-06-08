package handler

import (
	"WebProject_part7/internal/core"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type UserService interface {
	GetAll(ctx context.Context) ([]*core.User, error)
	GetById(ctx context.Context, id string) (*core.User, error)
	CreateUser(ctx context.Context, user *core.User) (*core.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (handler *UserHandler) InitRoutes(app *fiber.App) {
	app.Get("/users", handler.GetAll)
	app.Get("/users/:userId", handler.GetById)
	app.Post("/users", handler.CreateUser)
}

func (handler *UserHandler) GetAll(ctx *fiber.Ctx) error {
	ctxTimeout, cancel := context.WithTimeout(ctx.UserContext(), time.Second*2)
	defer cancel()

	usersChannel := make(chan []*core.User, 0)
	var err error
	var users []*core.User

	go func(channel chan<- []*core.User) {
		users, err = handler.service.GetAll(ctx.UserContext())
		channel <- users

	}(usersChannel)

	if err != nil {
		return err
	}

	select {
	case <-ctxTimeout.Done():
		fmt.Println("Processing timeot in Handler")
		break
	case users = <-usersChannel:
		fmt.Println("Finished processing in Handler")
	}

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"users": users,
		})
}

func (handler *UserHandler) GetById(ctx *fiber.Ctx) error {
	ctxTimeout, cancel := context.WithTimeout(ctx.UserContext(), time.Second*2)
	defer cancel()

	userChannel := make(chan *core.User, 0)
	var err error
	var user *core.User

	go func(channel chan<- *core.User) {
		user, err = handler.service.GetById(ctxTimeout, ctx.Params("userId"))
		channel <- user
	}(userChannel)

	if err != nil {
		return err
	}

	select {
	case <-ctxTimeout.Done():
		fmt.Println("Processing timeot in Handler")
		break
	case user = <-userChannel:
		fmt.Println("Finished processing in Handler")
	}

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"user": user,
		})
}

func (handler *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	user := &core.User{}

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	savedUser, err := handler.service.CreateUser(ctx.UserContext(), user)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	return ctx.Status(http.StatusCreated).JSON(
		fiber.Map{
			"user": savedUser,
		})

}
