package handlers

import (
	"net/http"
	"time"

	"github.com/aldysp34/mnc_test/rest/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Pin         string `json:"pin"`
}

func Register(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var existingUser model.User
	if err := db.Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Phone Number already registered"})
	}

	newUser := model.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Pin:         req.Pin,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&newUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating user"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "SUCCESS",
		"result": newUser,
	})
}
