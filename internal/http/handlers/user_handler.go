package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/example/golang-rest-boilerplate/internal/service"
	"github.com/example/golang-rest-boilerplate/pkg/response"
)

// UserHandler manages user CRUD endpoints.
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler constructs a new UserHandler.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List returns all users.
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userService.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"users": users})
}

// Get retrieves a single user by ID.
func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.Get(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"user": user})
}

type updateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// Update modifies the user's name.
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, req.Name)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"user": user})
}

// Delete removes a user account.
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
