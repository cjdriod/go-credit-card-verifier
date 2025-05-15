package controllers

import (
	"errors"
	"fmt"
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"github.com/cjdriod/go-credit-card-verifier/services"
	"github.com/cjdriod/go-credit-card-verifier/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController(userService services.UserServiceInterface) *UserController {
	return &UserController{userService: userService}
}

func (controller *UserController) RegisterNewUser(c *gin.Context) {
	accountCreateErr := errors.New("failed to create user")
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BindJSON(&body); err != nil || body.Username == "" || body.Password == "" || body.Email == "" {
		utils.BadRequest(c, []string{accountCreateErr.Error()})
		return
	}

	if err := controller.userService.CreateNewUser(body.Username, body.Password, body.Email); err != nil {
		utils.BadRequest(c, []string{accountCreateErr.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  configs.Constant.ApiStatus.Success,
		"message": fmt.Sprintf("User %s successfully created", body.Username),
	})
}

func (controller *UserController) LoginUser(c *gin.Context) {
	loginErr := errors.New("login failed, invalid email or password")
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil || body.Email == "" || body.Password == "" {
		utils.BadRequest(c, []string{loginErr.Error()})
		return
	}

	username, tokenString, err := controller.userService.LoginUser(body.Email, body.Password)

	if err != nil {
		utils.BadRequest(c, []string{loginErr.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 60*60*1, "/", "", configs.Constant.EnableHttpsMode, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  configs.Constant.ApiStatus.Success,
		"message": fmt.Sprintf("User %s successfully login", username),
	})
}

func (controller *UserController) LogoutUser(c *gin.Context) {
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
