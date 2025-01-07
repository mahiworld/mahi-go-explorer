package response

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// APIResponse generic response
type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse defines a success response
func SuccessResponse(c *gin.Context, code int, data any) {
	resp := new(APIResponse)
	resp.Success = true
	resp.Data = data
	c.JSON(code, resp)
}

// ErrorResponse defines an error response
func ErrorResponse(c *gin.Context, code int, message string) {
	resp := new(APIResponse)
	resp.Success = false
	resp.Message = message
	c.JSON(code, resp)
}

// LogAndErrorResponse logs the error and sends the error response
func LogAndErrorResponse(c *gin.Context, code int, message string, err error) {
	if err != nil {
		logger.Error(message, err.Error())
	}
	ErrorResponse(c, code, message)
}
