package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gbadegesintestimony/jwt-authentication/database"
	"github.com/gbadegesintestimony/jwt-authentication/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input models.RegisterRequest
	if err := c.BindJSON(&input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Determine first and last name. Prefer explicit fields; fall back to legacy `name`.
	first := input.FirstName
	last := input.LastName
	if first == "" && last == "" {
		if input.Name != "" {
			parts := strings.Fields(input.Name)
			if len(parts) == 1 {
				first = parts[0]
				last = ""
			} else if len(parts) > 1 {
				first = strings.Join(parts[:len(parts)-1], " ")
				last = parts[len(parts)-1]
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "first_name/last_name or name is required"})
			return
		}
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// create user with first and last name; keep Name for backward compatibility
	// Build user record (keep `Name` for compatibility)
	user := models.User{
		FirstName:    first,
		LastName:     last,
		Name:         strings.TrimSpace(first + " " + last),
		Email:        input.Email,
		PasswordHash: string(hashed),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use or invalid"})
		return
	}
	token, _ := generateToken(user.ID)

	userResponse := models.DetailedUserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		CreatedAt:     user.CreatedAt.Format(time.RFC3339),
		EmailVerified: false,
	}
	response := models.SuccessResponse{}
	response.Success.Status = http.StatusCreated
	response.Success.Data = models.AuthData{
		User:  userResponse,
		Token: token,
	}

	c.JSON(http.StatusCreated, response)
}

func Login(c *gin.Context) {
	var input models.LoginRequest
	if err := c.BindJSON(&input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	userResponse := models.DetailedUserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		CreatedAt:     user.CreatedAt.Format(time.RFC3339),
		EmailVerified: false,
	}
	response := models.SuccessResponse{}
	response.Success.Status = http.StatusOK
	response.Success.Data = models.AuthData{
		User:  userResponse,
		Token: token,
	}

	c.JSON(http.StatusOK, response)

}

// generateToken creates a new JWT token for a user
func generateToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	hoursStr := os.Getenv("JWT_EXPIRATION_HOURS")
	hours := 24
	if hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil {
			hours = h
		}
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ChangePassword allows an authenticated user to change their password
func ChangePassword(c *gin.Context) {
	var input models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get user id from context (set by middleware)
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := uid.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	// hash new password
	newHashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new password"})
		return
	}

	if err := database.DB.Model(&user).Update("password_hash", string(newHashed)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}
