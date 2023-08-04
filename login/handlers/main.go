package handlers

import (

	"github.com/gofiber/fiber/v2"

)


type CatalogHandler interface {
	GetCustomers(c *fiber.Ctx) error
}

type UserHandler interface {
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Session(c *fiber.Ctx) error
}
