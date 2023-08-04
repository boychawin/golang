package repositorys

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"login_jwt/errors"
	"login_jwt/models"
	"login_jwt/utils"
)

func NewAuthRepositoryDB(db *gorm.DB) AuthRepository {
	return &authRepositoryDB{db}
}

type authRepositoryDB struct {
	db *gorm.DB
}

func (u authRepositoryDB) Session(c *fiber.Ctx) (*AuthResponse, error) {
	claims := c.Locals("claims").(jwt.MapClaims)
	iss := claims["iss"].(string)

	user := models.Users{}
	err := u.db.Where("id = ?", iss).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No user found with the given username
			return nil, errors.NewUnexpectedError("ไม่พบข้อมูล")
		}
		// Other error occurred
		return nil, errors.NewUnexpectedError("unexpected error")
	}


	users := &AuthResponse{
		Data:     &user,
		Messages: "",
	}

	return users, err

}

func (u authRepositoryDB) Refresh(c *fiber.Ctx) (*AuthResponse, error) {
	claims := c.Locals("claims").(jwt.MapClaims)
	iss := claims["iss"].(string)

	user := models.Users{}
	err := u.db.Where("id = ?", iss).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No user found with the given ID
			return nil, errors.NewUnexpectedError("ไม่พบข้อมูล")
		}
		// Other error occurred
		return nil, errors.NewUnexpectedError("unexpected error: " + err.Error())
	}

	accessString, refreshString, err := utils.CreateTokens(user)
	if err != nil {
		return nil, errors.NewUnexpectedError("unexpected error: " + err.Error())
	}

	users := &AuthResponse{
		Messages:     "",
		AccessToken:  accessString,
		RefreshToken: refreshString,
	}

	return users, nil
}



func (u authRepositoryDB) Login(c *fiber.Ctx) (*AuthResponse, error) {
	request := SignupRequest{}

	err := c.BodyParser(&request)
	if err != nil {
		return nil, err
	}

	if request.UserName == "" || request.PassWord == "" {
		return nil, errors.NewUnexpectedError("ต้องระบุชื่อผู้ใช้หรือรหัสผ่าน")
	}

	user := models.Users{}

	err = u.db.Where("user_name = ? ", request.UserName).First(&user).Error
	if err != nil {
		return nil, errors.NewUnexpectedError("ไม่พบข้อมูล")
	}

	// Compare the provided password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(request.PassWord))
	if err != nil {
		return nil, errors.NewUnexpectedError("ชื่อผู้ใช้ หรือ รหัสผ่าน ไม่ถูกต้อง")
	}

	accessString, refreshString, err := utils.CreateTokens(user)
	if err != nil {
		return nil, errors.NewUnexpectedError("ชื่อผู้ใช้ หรือ รหัสผ่าน ไม่ถูกต้อง")
	}

	users := &AuthResponse{
		Messages:     "",
		AccessToken:  accessString,
		RefreshToken: refreshString,
	}
	return users, err
}


func (u authRepositoryDB) Register(c *fiber.Ctx) (*AuthResponse, error) {

	request := models.Users{}
	err := c.BodyParser(&request)
	if err != nil {
		return nil, err
	}

	if request.UserName == "" || request.PassWord == "" || request.FirstName == "" || request.LastName == "" {
		return nil, errors.NewUnexpectedError("ต้องระบุชื่อผู้ใช้หรือรหัสผ่าน")
	}

	var count int64
	u.db.Model(&models.Users{}).Where("user_name = ?", request.UserName).Count(&count)

	if count >= 1 {
		// Multiple users with the specified username exist, return an error
		return nil, errors.NewUnexpectedError("มีข้อมูลแล้ว")
	}

	// Encode and hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.PassWord), bcrypt.DefaultCost)
	if err != nil {
		// Handle the error if password hashing fails
		return nil, errors.NewUnexpectedError(err.Error())
	}

	// Create the user object
	user := models.Users{
		UserName:  request.UserName,
		PassWord:  string(hashedPassword),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Role:      "client", // Assuming a default role value
		Status:    1,        // Assuming a default status value
	}

	// Insert the user into the database
	tx := u.db.Create(&user)
	if tx.Error != nil {
		// Handle the error that occurred during the database insert
		return nil, errors.NewUnexpectedError(tx.Error.Error())
	}

	users2 := &AuthResponse{
		Data:     &user,
		Messages: "",
	}
	return users2, nil
}
