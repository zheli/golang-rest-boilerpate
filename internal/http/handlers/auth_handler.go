package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/example/golang-rest-boilerplate/internal/service"
	"github.com/example/golang-rest-boilerplate/pkg/response"
)

// AuthHandler handles authentication related HTTP requests.
type AuthHandler struct {
	authService   *service.AuthService
	googleService *service.GoogleOAuthService
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(authService *service.AuthService, googleService *service.GoogleOAuthService) *AuthHandler {
	return &AuthHandler{authService: authService, googleService: googleService}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register creates a new user account.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, gin.H{"user": user})
}

// Login authenticates a user with email/password.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"token": token, "user": user})
}

// GoogleLogin initiates the Google OAuth flow.
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	if h.googleService == nil {
		response.Error(c, http.StatusServiceUnavailable, "google oauth is not configured")
		return
	}

	state := uuid.NewString()
	url := h.googleService.AuthCodeURL(state)
	response.JSON(c, http.StatusOK, gin.H{"auth_url": url, "state": state})
}

// GoogleCallback handles the OAuth callback from Google.
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	if h.googleService == nil {
		response.Error(c, http.StatusServiceUnavailable, "google oauth is not configured")
		return
	}

	code := c.Query("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "code query param missing")
		return
	}

	token, err := h.googleService.Exchange(c.Request.Context(), code)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := h.googleService.FetchUserInfo(c.Request.Context(), token)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.FindOrCreateOAuthUser(c.Request.Context(), userInfo.Name, userInfo.Email, "google", userInfo.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	jwtToken, err := h.authService.GenerateToken(user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"token": jwtToken, "user": user})
}
