package handlers

import (
	"mahi-go-explorer/internal/api/response"
	userpkg "mahi-go-explorer/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthRoutes defines auth routes
func AuthRoutes(r *gin.Engine, s userpkg.Service) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/signup", signupHandler(s))
		auth.POST("/login", loginHandler(s))
	}
}

func signupHandler(s userpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		//bind the request
		var req userpkg.CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Bad Request", err)
			return
		}

		//validate required fields
		if req.Email == "" || req.Password == "" {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Bad Request", nil)
			return
		}

		u, err := req.CreateUser()
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Bad Request", err)
			return
		}

		resp, err := s.CreateUser(u)
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}

		response.SuccessResponse(c, http.StatusCreated, resp)
	}
}

func loginHandler(s userpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		//bind the request
		var req userpkg.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Bad Request", err)
			return
		}

		//validate required fields
		if req.Email == "" || req.Password == "" {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Bad Request", nil)
			return
		}

		token, err := s.LoginUser(&req)
		if err != nil {
			switch err.Error() {
			case "user not found":
				response.LogAndErrorResponse(c, http.StatusNotFound, "User Not Found", err)
				return
			case "invalid password":
				response.LogAndErrorResponse(c, http.StatusUnauthorized, "Invalid Password", err)
				return
			default:
				response.LogAndErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", err)
				return
			}
		}

		response.SuccessResponse(c, http.StatusOK, gin.H{
			"accessToken": token,
		})
	}
}
