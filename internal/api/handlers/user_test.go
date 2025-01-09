package handlers

import (
	"bytes"
	"encoding/json"
	userpkg "mahi-go-explorer/pkg/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockUserService struct {
	userpkg.Service
	CreateUserMock func(user *userpkg.User) (any, error)
	GetUsersMock   func(conds bson.M, opts *options.FindOptions) ([]userpkg.User, error)
	GetUserMock    func(conds bson.M, opts *options.FindOneOptions) (*userpkg.User, error)
	UpdateUserMock func(conds bson.M, update bson.M, opts *options.UpdateOptions) (any, error)
	DeleteUserMock func(conds bson.M) (any, error)
}

func (m *mockUserService) CreateUser(req *userpkg.User) (any, error) {
	return m.CreateUserMock(req)
}

func TestCreateUserHandler(t *testing.T) {
	// Create a mock user service
	mockUserService := &mockUserService{
		CreateUserMock: func(user *userpkg.User) (any, error) {
			return "1234567890abcdef12345678", nil
		},
	}

	type testCase struct {
		Name               string
		RequestBody        userpkg.CreateRequest
		ExpectedStatusCode int
		ExpectedError      bool
		ExpectedMessage    string
	}

	tests := []testCase{
		{
			Name: "Create user",
			RequestBody: userpkg.CreateRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "XXXXXXXXXXXXXXXXX",
				Phone:     "1234567890",
				Role:      "admin",
				Password:  "XXXXXXXX",
			},
			ExpectedStatusCode: http.StatusCreated,
			ExpectedError:      false,
		},
		{
			Name: "Create user with missing email",
			RequestBody: userpkg.CreateRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "",
				Phone:     "1234567890",
				Role:      "admin",
				Password:  "strongpassword",
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      true,
			ExpectedMessage:    "Email and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Set up mock behavior based on test case
			mockUserService.CreateUserMock = func(user *userpkg.User) (any, error) {
				if tt.ExpectedStatusCode == http.StatusInternalServerError {
					return nil, assert.AnError
				}
				return "1234567890abcdef12345678", nil
			}

			// Create Gin engine and set it to Test mode
			gin.SetMode(gin.TestMode)

			// Create a new HTTP request
			reqBody, _ := json.Marshal(tt.RequestBody)
			req, err := http.NewRequest("POST", "/api/user", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			// Record the response
			rr := httptest.NewRecorder()

			// Create a Gin router and set the route
			router := gin.Default()
			router.POST("/api/user", createUserHandler(mockUserService))

			// Serve the HTTP request and get the response
			router.ServeHTTP(rr, req)

			// Assert the response status code
			assert.Equal(t, tt.ExpectedStatusCode, rr.Code)

			// Decode the response body
			var response map[string]interface{}
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Fatal(err)
			}

			// If an error is expected, verify the error message
			if tt.ExpectedError {
				assert.Contains(t, response["message"], tt.ExpectedMessage)
			} else {
				// If no error, verify the success response
				assert.Equal(t, "1234567890abcdef12345678", response["data"])
			}
		})
	}
}
