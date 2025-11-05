package controllers

import (
	"net/http"

	"github.com/gbadegesintestimony/jwt-authentication/database"
	"github.com/gbadegesintestimony/jwt-authentication/models"
	"github.com/gin-gonic/gin"
)

type UpdateUserInput struct {
	Name string `json:"name"`
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
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
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

	user.Name = input.Name
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}
