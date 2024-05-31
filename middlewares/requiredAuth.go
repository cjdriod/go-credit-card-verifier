package middlewares

import (
	"errors"
	"fmt"
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/database"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func abortUnauthorizedRequest(c *gin.Context) {
	unauthorizedErr := errors.New("unauthorized action")

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status": configs.Constant.ApiStatus.Fail,
		"errors": []string{unauthorizedErr.Error()},
	})
}

func RequireAuthentication(c *gin.Context) {
	authenticationToken, cookieErr := c.Cookie("Authorization")
	if cookieErr != nil || authenticationToken == "" {

		abortUnauthorizedRequest(c)
		return
	}

	token, tokenErr := jwt.Parse(authenticationToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return configs.Constant.JwtSecret, nil
	})

	if tokenErr != nil || !token.Valid {
		abortUnauthorizedRequest(c)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		abortUnauthorizedRequest(c)
		return
	} else {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			abortUnauthorizedRequest(c)
			return
		}

		var user models.User
		database.DB.Find(&user, claims["sub"])
		if user.ID == 0 {
			abortUnauthorizedRequest(c)
			return
		}

		c.Set("user", user)
	}

	c.Next()
}
