package services

import (
	"login_jwt/models"
	"github.com/gofiber/fiber/v2"
)



type UserResponse struct {
	Data         *models.Users `json:"data"`
	Messages     string        `json:"messages"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserServiceRedis interface {
	Login(c *fiber.Ctx) (*UserResponse, error)
	Register(c *fiber.Ctx) (*UserResponse, error)
	Refresh(c *fiber.Ctx) (*UserResponse, error)
	Session(c *fiber.Ctx) (*UserResponse, error)
}
