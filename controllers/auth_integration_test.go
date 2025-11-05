package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gbadegesintestimony/jwt-authentication/database"
	"github.com/gbadegesintestimony/jwt-authentication/middleware"
	"github.com/gbadegesintestimony/jwt-authentication/models"
	"github.com/gin-gonic/gin"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) *gin.Engine {
	// use in-memory sqlite for tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	database.DB = db
	database.DB.AutoMigrate(&models.User{})

	g := gin.Default()
	// register handlers directly to avoid import cycle with routes
	api := g.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", Register)
			auth.POST("/login", Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/change-password", ChangePassword)
			protected.GET("/me", GetProfile)
			protected.PUT("/me", UpdateProfile)
		}
	}
	return g
}

func TestRegisterLoginChangeProfile(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	g := setupTestServer(t)

	// Register
	regBody := models.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "password123"}
	b, _ := json.Marshal(regBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d: %s", w.Code, w.Body.String())
	}

	// Login
	loginBody := models.LoginRequest{Email: "alice@example.com", Password: "password123"}
	b, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	g.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 ok on login, got %d: %s", w.Code, w.Body.String())
	}
	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["token"]
	if token == "" {
		t.Fatalf("expected token in login response")
	}

	// Change password
	changeBody := models.ChangePasswordRequest{CurrentPassword: "password123", NewPassword: "newpass456"}
	b, _ = json.Marshal(changeBody)
	req = httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	g.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 ok on change-password, got %d: %s", w.Code, w.Body.String())
	}

	// Get profile
	req = httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	g.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 ok on get profile, got %d: %s", w.Code, w.Body.String())
	}
}
