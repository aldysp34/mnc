// handlers/topup.go
package handlers

import (
	"net/http"
	"time"

	"github.com/aldysp34/mnc_test/rest/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TopUpRequest struct {
	Amount int `json:"amount"`
}

type TopUpResponse struct {
	Status string `json:"status"`
	Result struct {
		TopUpID       string `json:"top_up_id"`
		AmountTopUp   int    `json:"amount_top_up"`
		BalanceBefore int    `json:"balance_before"`
		BalanceAfter  int    `json:"balance_after"`
		CreatedDate   string `json:"created_date"`
	} `json:"result"`
}

// TopUp handles the top-up request
func TopUp(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	// Get user from JWT claims
	user := c.Get("user").(*JWTClaims)

	// Bind request data
	req := new(TopUpRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Validate amount (must be greater than 0)
	if req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid top-up amount"})
	}

	// Retrieve user from the database
	var userModel model.User
	if err := db.Where("id = ?", user.UserID).First(&userModel).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "User not found"})
	}

	// Create a new top-up transaction
	topUpID := uuid.New().String()
	balanceBefore := userModel.Balance
	balanceAfter := balanceBefore + req.Amount

	// Create and save the top-up record
	topUp := model.TopUp{
		ID:            topUpID,
		UserID:        userModel.ID,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     time.Now(),
	}

	if err := db.Create(&topUp).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error processing top-up"})
	}

	// Update user's balance
	userModel.Balance = balanceAfter
	if err := db.Save(&userModel).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating balance"})
	}

	// Prepare response
	response := TopUpResponse{
		Status: "SUCCESS",
	}
	response.Result.TopUpID = topUpID
	response.Result.AmountTopUp = req.Amount
	response.Result.BalanceBefore = balanceBefore
	response.Result.BalanceAfter = balanceAfter
	response.Result.CreatedDate = topUp.CreatedAt.Format("2006-01-02 15:04:05")

	return c.JSON(http.StatusOK, response)
}
