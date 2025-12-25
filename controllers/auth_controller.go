package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gbadegesintestimony/jwt-authentication/database"
	"github.com/gbadegesintestimony/jwt-authentication/models"
	"github.com/gbadegesintestimony/jwt-authentication/utils"
	"github.com/gin-gonic/gin"
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

	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
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
	token, _ := utils.GenerateToken(32)

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
	response.Success.Message = "User Registered successful"
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
	token, err := utils.GenerateToken(user.ID)
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
	response.Success.Message = "Login successful"
	response.Success.Data = models.AuthData{
		User:  userResponse,
		Token: token,
	}

	c.JSON(http.StatusOK, response)

}

// ChangePassword allows an authenticated user to change their password
func ChangePassword(c *gin.Context) {
	var input models.ChangePasswordRequest
	if err := c.BindJSON(&input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})s
		return
	}

	// get user id from context (set by middleware)
	uid, exists := c.Get("userID")
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

func ForgotPassword(c *gin.Context) {
	type Body struct {
		Email string `json:"email" binding:"required,email"`
	}

	var req Body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// To prevent email enumeration, respond with success even if user not found
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	log.Printf("User found: %v", user.Email)

	otp, err := utils.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate OTP"})
		return
	}

	user.ResetOTP = otp
	user.ResetExpiry = time.Now().Add(15 * time.Minute)
	database.DB.Save(&user)

	if err := utils.SendEmail(
		user.Email,
		"Password Reset OTP",
		"Your OTP for password reset is: "+otp+"\nIt expires in 15 minutes.",
	); err != nil {

		log.Println("SMTP ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send OTP email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP has been sent to your email",
	})
}

func VerifyOTP(c *gin.Context) {
	type Request struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required"`
	}
	var req Request
	if err := c.BindJSON(&req); err != nil {
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or OTP"})
		return
	}

	if user.ResetOTP != req.OTP || time.Now().After(user.ResetExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

func ResetPassword(c *gin.Context) {
	type Input struct {
		Email           string `json:"email" binding:"required,email"`
		OTP             string `json:"otp" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
	}

	var req Input
	if err := c.BindJSON(&req); err != nil {
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new password and confirm password do not match"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or OTP"})
		return
	}

	if user.ResetOTP != req.OTP || time.Now().After(user.ResetExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired OTP"})
		return
	}

	newHashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(newHashed)
	user.ResetOTP = ""
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "password has been reset"})
}
