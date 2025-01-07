package handlers

import (
	"errors"
	userpkg "mahi-go-explorer/pkg/user"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all routes
func RegisterRoutes(
	r *gin.Engine,
	userService userpkg.Service,
) {
	AuthRoutes(r, userService)
	UserRoutes(r, userService)
}

func getUserContext(c *gin.Context) (*userpkg.UserContext, error) {
	userContext, ok := c.Get("user")
	if !ok {
		return nil, errors.New("Context key doesn't exist")
	}

	currentUser, ok := userContext.(*userpkg.UserContext)
	if !ok {
		return nil, errors.New("Context key is not of type UserContext")
	}

	return currentUser, nil
}
