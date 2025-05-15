package middlewares

import (
	"fmt"
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/database"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"github.com/cjdriod/go-credit-card-verifier/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	unauthorizedErr := []string{"unauthorized action"}

	return func(c *gin.Context) {
		authenticationToken, cookieErr := c.Cookie("Authorization")
		if cookieErr != nil || authenticationToken == "" {
			utils.UnauthorizedRequest(c, unauthorizedErr)
			return
		}

		token, tokenErr := jwt.Parse(authenticationToken, func(token *jwt.Token) (interface{}, error) {
			if len(configs.Constant.JwtSecret) == 0 {
				return nil, fmt.Errorf("JWT secret not found")
			}

			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
			}

			return configs.Constant.JwtSecret, nil
		})

		if tokenErr != nil || !token.Valid {
			utils.UnauthorizedRequest(c, unauthorizedErr)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); !ok {
			utils.UnauthorizedRequest(c, unauthorizedErr)
			return
		} else {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				utils.UnauthorizedRequest(c, unauthorizedErr)
				return
			}

			var user models.User
			database.DB.Find(&user, claims["sub"])
			if user.ID == 0 {
				utils.UnauthorizedRequest(c, unauthorizedErr)
				return
			}

			c.Set("user", user)
		}

		c.Next()
	}
}
