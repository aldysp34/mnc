// handlers/login.go
package handlers

import (
	"net/http"
	"time"

	"github.com/aldysp34/mnc_test/rest/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Pin         string `json:"pin"`
}

type JWTClaims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key") // Replace with a more secure secret

// Generate JWT tokens (access and refresh)
func generateTokens(user *model.User) (string, string, error) {
	// Access token
	accessTokenClaims := &JWTClaims{
		UserID:      user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15 minutes expiry
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshTokenClaims := &JWTClaims{
		UserID:      user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days expiry
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func Login(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Check if user exists
	var user model.User
	if err := db.Where("phone_number = ? AND pin = ?", req.PhoneNumber, req.Pin).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Phone Number and PIN doesn’t match."})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error during login"})
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := generateTokens(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error generating token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "SUCCESS",
		"result": map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}
