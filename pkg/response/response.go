package response

import "github.com/gin-gonic/gin"

// JSON writes a standardized JSON response.
func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

// Error writes an error JSON response.
func Error(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}
