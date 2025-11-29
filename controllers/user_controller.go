package controllers

import (
	"net/http"

	"github.com/gbadegesintestimony/jwt-authentication/database"
	"github.com/gbadegesintestimony/jwt-authentication/models"
	"github.com/gin-gonic/gin"
)

type UpdateUserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Don't return password hash
	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	})
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	// keep Name in sync
	user.Name = user.FirstName + " " + user.LastName
	database.DB.Save(&user)

	userResponse := models.UpdateResponse{
		Message:       "Profile updated successfully",
		Firstname:     user.FirstName,
		Lastname:      user.LastName,
		UpdateAt:      user.UpdatedAt.Format("2006-01-02 15:04:05"),
		EmailVerified: false, // Placeholder; implement email verification logic as needed
	}
	response := models.SuccessResponse{}
	response.Success.Status = 200
	response.Success.Message = "Profile updated successfully"
	response.Success.Data = userResponse

	c.JSON(http.StatusOK, response)
}
