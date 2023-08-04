package repositorys

import (
	"github.com/gofiber/fiber/v2"
	"login_jwt/models"
)



type AuthResponse struct {
	Data         *models.Users `json:"data"`
	Messages     string        `json:"messages"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
}


type AuthRepository interface {
	Register(c *fiber.Ctx) (*AuthResponse, error)
	Login(c *fiber.Ctx) (*AuthResponse, error)
	Refresh(c *fiber.Ctx) (*AuthResponse, error)
	Session(c *fiber.Ctx) (*AuthResponse, error)
}

type Users struct {
	ID       uint
	Username string `gorm:"unique;size(50)"`
}
type SignupRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}
