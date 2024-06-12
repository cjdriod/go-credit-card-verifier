package controllers

import (
	"errors"
	"fmt"
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/database"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"github.com/cjdriod/go-credit-card-verifier/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func RegisterUserController(c *gin.Context) {
	accountCreateErr := errors.New("failed to create user")
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := c.BindJSON(&body)
	if err != nil || body.Username == "" || body.Password == "" || body.Email == "" {
		utils.BadRequest(c, []string{accountCreateErr.Error()})
		return
	}

	hashPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if hashErr != nil {
		fmt.Println("Bcrypt failed:", err)
		utils.BadRequest(c, []string{accountCreateErr.Error()})
		return
	}

	newUser := models.User{Email: body.Email, Username: body.Username, Password: string(hashPassword)}
	result := database.DB.Create(&newUser)

	if result.Error != nil {
		fmt.Println("Database error:", result.Error)
		utils.BadRequest(c, []string{accountCreateErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  configs.Constant.ApiStatus.Success,
		"message": fmt.Sprintf("User %s successfully created", body.Username),
	})
}

func LoginUserController(c *gin.Context) {
	loginErr := errors.New("login failed, invalid email or password")
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.BindJSON(&body)
	if err != nil {
		utils.BadRequest(c, []string{loginErr.Error()})
		return
	}

	var user models.User
	database.DB.First(&user, "email = ?", body.Email)

	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if user.ID == 0 || bcryptErr != nil {
		utils.UnauthorizedRequest(c, []string{loginErr.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, signErr := token.SignedString(configs.Constant.JwtSecret)

	if signErr != nil {
		fmt.Println("JWT signature error:", signErr)
		utils.UnauthorizedRequest(c, []string{loginErr.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 60*60*1, "/", "", configs.Constant.EnableHttpsMode, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  configs.Constant.ApiStatus.Success,
		"message": fmt.Sprintf("User %s successfully login", user.Username),
	})
}

func LogoutUserController(c *gin.Context) {
	unauthorizedErr := errors.New("unauthorized action")
	user, err := c.Get("user")

	if !err {
		utils.UnauthorizedRequest(c, []string{unauthorizedErr.Error()})
		return
	}
	c.SetCookie("Authorization", "", -1, "/", "", configs.Constant.EnableHttpsMode, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  configs.Constant.ApiStatus.Success,
		"message": fmt.Sprintf("User %s successfully login", user.(models.User).Username),
	})
}
