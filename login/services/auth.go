package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"login_jwt/errors"
	"login_jwt/models"
	"login_jwt/repositorys"
)

type userServiceRedis struct {
	userRepo    repositorys.AuthRepository
	redisClient *redis.Client
}

func NewAuthService(userRepo repositorys.AuthRepository, redisClient *redis.Client) UserServiceRedis {
	return userServiceRedis{userRepo, redisClient}
}

func (s userServiceRedis) Session(c *fiber.Ctx) (users *UserResponse, err error) {
	claims := c.Locals("claims").(jwt.MapClaims)
	userID := claims["iss"].(string)

	// Concatenate the string with the converted userID
	key := "service::Session" + userID

	//Redis GET
	if userJson, err := s.redisClient.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(userJson), &users) == nil {
			fmt.Println("redis")
			return users, nil
		}
	}

	// Repository
	usersDB, err := s.userRepo.Session(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(err.Error())
		}

		// logs.Error(err)
		return nil, errors.NewUnexpectedError(err.Error())
	}

	user := &UserResponse{
		Data: &models.Users{
			ID:        usersDB.Data.ID,
			FirstName: usersDB.Data.FirstName,
			LastName:  usersDB.Data.LastName,
			UserName:  usersDB.Data.UserName,
			UpdatedAt: usersDB.Data.UpdatedAt,
			CreatedAt: usersDB.Data.CreatedAt,
		},
		Messages: usersDB.Messages,
	}

	// Redis SET
	if data, err := json.Marshal(user); err == nil {
		s.redisClient.Set(context.Background(), key, string(data), time.Second*10)
	}

	return user, nil
}

func (s userServiceRedis) Refresh(c *fiber.Ctx) (users *UserResponse, err error) {
	claims := c.Locals("claims").(jwt.MapClaims)
	userID := claims["iss"].(string)

	// Concatenate the string with the converted userID
	key := "service::Refresh" + userID

	//Redis GET
	if productJson, err := s.redisClient.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(productJson), &users) == nil {
			fmt.Println("redis")
			return users, nil
		}
	}
	// Repository
	usersRepo, err := s.userRepo.Refresh(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("user not found")
		}

		// logs.Error(err)
		return nil, errors.NewUnexpectedError("unexpected error")
	}

	user := &UserResponse{
		Messages:     usersRepo.Messages,
		AccessToken:  usersRepo.AccessToken,
		RefreshToken: usersRepo.RefreshToken,
	}

	// Redis SET
	if data, err := json.Marshal(user); err == nil {
		s.redisClient.Set(context.Background(), key, string(data), time.Second*10)
	}

	return user, nil
}

func (s userServiceRedis) Login(c *fiber.Ctx) (users *UserResponse, err error) {

	usersRepo, err := s.userRepo.Login(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("user not found")
		}
		return nil, errors.NewUnexpectedError(err.Error())
	}

	user := &UserResponse{
		Messages:     usersRepo.Messages,
		AccessToken:  usersRepo.AccessToken,
		RefreshToken: usersRepo.RefreshToken,
	}

	return user, nil
}

func (s userServiceRedis) Register(c *fiber.Ctx) (users *UserResponse, err error) {

	usersRepo, err := s.userRepo.Register(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("ไม่พบข้อมูล")
		}
		return nil, errors.NewUnexpectedError(err.Error())
	}

	user := &UserResponse{
		Messages: usersRepo.Messages,
	}

	return user, nil
}
