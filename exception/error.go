package exception

import "github.com/gin-gonic/gin"

// fungsi untuk mengirimkan error
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
