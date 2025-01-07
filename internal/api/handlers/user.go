package handlers

import (
	"mahi-go-explorer/internal/api/middleware"
	"mahi-go-explorer/internal/api/response"
	userpkg "mahi-go-explorer/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRoutes defnies user service routes
func UserRoutes(r *gin.Engine, s userpkg.Service) {
	user := r.Group("/api/user")
	user.Use(middleware.Authenticate())
	{
		user.POST("", createUserHandler(s))
		user.GET("", getUsersHandler(s))
		user.GET("/:id", getUserHandler(s))
		user.PUT("/:id", updateUserHandler(s))
		user.DELETE("/:id", deleteUserHandler(s))
	}
}

func createUserHandler(s userpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req userpkg.CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Bad Request", err)
			return
		}

		if req.Email == "" || req.Password == "" {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Email and password are required", nil)
			return
		}

		u, _ := req.CreateUser()

		res, err := s.CreateUser(u)
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err)
			return
		}

		response.SuccessResponse(c, http.StatusCreated, res)
	}
}

func getUsersHandler(s userpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		conds := bson.M{}
		users, err := s.GetUsers(conds, nil)
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusInternalServerError, "Failed to get users", err)
			return
		}

		response.SuccessResponse(c, http.StatusOK, users)
	}
}

func getUserHandler(s userpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		cu, err := getUserContext(c)
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusInternalServerError, "Failed to get user context", err)
			return
		}

		var uID primitive.ObjectID
		idstr := c.Param("id")
		if idstr == "me" {
			uID = cu.ID
		} else {
			objID, err := primitive.ObjectIDFromHex(idstr)
			if err != nil {
				response.LogAndErrorResponse(c, http.StatusBadRequest, "Invalid ID", err)
				return
			}
			uID = objID
		}

		user, err := s.GetUser(bson.M{"_id": uID}, nil)
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusInternalServerError, "Failed to get user", err)
			return
		}

		response.SuccessResponse(c, http.StatusOK, user)
	}
}

func updateUserHandler(s userpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(idstr)
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Invalid ID", err)
			return
		}

		var req userpkg.UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Bad Request", err)
			return
		}

		update := bson.M{}
		if req.FirstName != "" {
			update["firstName"] = req.FirstName
		}
		if req.LastName != "" {
			update["lastName"] = req.LastName
		}
		if req.Email != "" {
			update["email"] = req.Email
		}
		if req.Role != "" {
			update["role"] = req.Role
		}

		res, err := s.UpdateUser(bson.M{"_id": objID}, update, nil)
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
			return
		}

		response.SuccessResponse(c, http.StatusOK, res)
	}
}

func deleteUserHandler(s userpkg.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(idstr)
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusBadRequest, "Invalid ID", err)
			return
		}

		res, err := s.DeleteUser(bson.M{"_id": objID})
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusInternalServerError, "Failed to delete user", err)
			return
		}

		response.SuccessResponse(c, http.StatusOK, res)
	}
}
