// +build unit

package http

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/clean-architecture/apperr"
	authormocks "github.com/jayvib/clean-architecture/author/mocks"
	"github.com/jayvib/clean-architecture/config"
	"github.com/jayvib/clean-architecture/model"
	"github.com/jayvib/clean-architecture/user/mocks"
	userusecase "github.com/jayvib/clean-architecture/user/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	mockUser = &model.User{
		ID:        "uniqueid",
		Firstname: "Luffy",
		Lastname:  "Monkey",
		Email:     "luffy.monkey@gmail.com",
		Username:  "luff.monkey",
		Password:  "pirateking",
	}
)

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestMain(m *testing.M) {
	//logrus.SetLevel(logrus.DebugLevel)
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestGetByID(t *testing.T) {
	expected := &model.User{
		ID:        "uniqueid",
		Firstname: "Luffy",
		Lastname:  "Monkey",
		Email:     "luffy.monkey@gmail.com",
		Username:  "luff.monkey",
	}
	e := gin.Default()
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(expected, nil)
	uc := userusecase.New(repo, authorMockRepo)
	RegisterHandlers(nil, e, uc)
	w := performRequest(e, http.MethodGet, "/user/1", nil)

	// Check the status code
	require.Equal(t, http.StatusOK, w.Code)

	// Convert the response into a struct
	var userResp model.User
	err := json.NewDecoder(w.Body).Decode(&userResp)
	require.NoError(t, err, "error found while decoding the json")
	assert.Equal(t, expected, &userResp, "user not equal")
	repo.AssertExpectations(t)
}

func TestGetByEmail(t *testing.T) {
	t.SkipNow()
	expected := &model.User{
		ID:        "uniqueid",
		Firstname: "Luffy",
		Lastname:  "Monkey",
		Email:     "luffy.monkey@gmail.com",
		Username:  "luff.monkey",
	}
	e := gin.Default()
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(expected, nil)
	uc := userusecase.New(repo, authorMockRepo)
	RegisterHandlers(nil, e, uc)
	w := performRequest(e, http.MethodGet, "/user/email/luffy.monkey@gmail.com", nil)

	// Check the status code
	require.Equal(t, http.StatusOK, w.Code)

	// Convert the response into a struct
	var userResp model.User
	err := json.NewDecoder(w.Body).Decode(&userResp)
	require.NoError(t, err, "error found while decoding the json")
	assert.Equal(t, expected, &userResp, "user not equal")
	repo.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	e := gin.Default()
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
	uc := userusecase.New(repo, authorMockRepo)
	RegisterHandlers(nil, e, uc)

	payload, err := json.Marshal(&mockUser)
	require.NoError(t, err)
	r := bytes.NewReader(payload)
	w := performRequest(e, http.MethodPost, "/user/update", r)
	require.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	e := gin.Default()
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	repo.On("Store", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
	repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, apperr.ItemNotFound).Once()
	repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(nil, apperr.ItemNotFound).Once()
	authorMockRepo.On("Store", mock.Anything, mock.AnythingOfType("*model.Author")).Return(nil).Once()
	uc := userusecase.New(repo, authorMockRepo)
	RegisterHandlers(nil, e, uc)

	payload, err := json.Marshal(&mockUser)
	require.NoError(t, err)
	r := bytes.NewReader(payload)
	w := performRequest(e, http.MethodPost, "/user", r)
	require.Equal(t, http.StatusCreated, w.Code)
	repo.AssertExpectations(t)
	authorMockRepo.AssertExpectations(t)

	// Check the body....
	// This should return the user with the id.
	var resultUser *model.User
	err = json.NewDecoder(w.Body).Decode(&resultUser)
	require.NoError(t, err, "error while decoding the response body")
	assert.NotEmpty(t, resultUser.ID)

	// Check also the created and updated date
	assert.NotEqual(t, time.Time{}, resultUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, resultUser.UpdatedAt)
	assert.Empty(t, resultUser.Password)
}
func TestDelete(t *testing.T) {
	e := gin.Default()
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	repo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)
	uc := userusecase.New(repo, authorMockRepo)
	RegisterHandlers(nil, e, uc)
	w := performRequest(e, http.MethodDelete, "/user/1", nil)
	require.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestRegister(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)
	e := gin.Default()
	api := e.Group("/")
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	userCopy := &(*mockUser)
	userCopy.ID = ""
	repo.On("Store", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, apperr.ItemNotFound).Once()
	repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(nil, apperr.ItemNotFound).Once()

	authorMockRepo.On("Store", mock.Anything, mock.AnythingOfType("*model.Author")).Return(nil).Once()
	uc := userusecase.New(repo, authorMockRepo)
	RegisterHandlers(conf, api, uc)
	payload, err := json.Marshal(userCopy)
	require.NoError(t, err)
	r := bytes.NewReader(payload)
	w := performRequest(e, http.MethodPost, "/user/new", r)
	require.Equal(t, http.StatusOK, w.Code)
	var resUser model.User
	err = json.NewDecoder(w.Body).Decode(&resUser)
	require.NoError(t, err)
	assert.Equal(t, userCopy.Username, resUser.Username)
	assert.Equal(t, userCopy.Firstname, resUser.Firstname)
	assert.Equal(t, userCopy.Lastname, resUser.Lastname)
	assert.Equal(t, userCopy.Email, resUser.Email)
	assert.NotEqual(t, time.Time{}, resUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, resUser.UpdatedAt)
	//assert.Empty(t, resUser.Password)
	assert.NotEmpty(t, resUser.ID)

	// Test the authentication token
	assert.NotEmpty(t, resUser.Token)
	require.NotEmpty(t, w.Result().Cookies())
	token := w.Result().Cookies()[0].Value
	claims := new(Claims)
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JWTToken), nil
	})
	require.NoError(t, err)
}

func TestLogin(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)
	require.NotNil(t, conf)
	// Simulate Registration
	e := gin.Default()
	api := e.Group("/")
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	userCopy := &(*mockUser)
	userCopy.ID = ""
	repo.On("Store", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, apperr.ItemNotFound).Once()
	repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(nil, apperr.ItemNotFound).Once()

	authorMockRepo.On("Store", mock.Anything, mock.AnythingOfType("*model.Author")).Return(nil).Once()
	uc := userusecase.New(repo, authorMockRepo)
	RegisterHandlers(conf, api, uc)
	payload, err := json.Marshal(userCopy)
	require.NoError(t, err)
	r := bytes.NewReader(payload)
	w := performRequest(e, http.MethodPost, "/user/new", r)
	require.Equal(t, http.StatusOK, w.Code)
	var resUser model.User
	err = json.NewDecoder(w.Body).Decode(&resUser)
	require.NoError(t, err)
	t.Run("Success", func(t *testing.T) {
		// Start the login request
		repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(&resUser, nil).Once()
		credential := &model.User{
			Username: mockUser.Username,
			Password: mockUser.Password,
		}

		payload, err = json.Marshal(credential)
		require.NoError(t, err)

		// perform request
		w = performRequest(e, http.MethodPost, "/user/login", bytes.NewReader(payload))
		require.Equal(t, http.StatusOK, w.Code)
		var authUser model.User
		err = json.NewDecoder(w.Body).Decode(&authUser)
		require.NoError(t, err)
		assert.Empty(t, authUser.Password)
		assert.NotEmpty(t, authUser.Token)

		// validate the token
		token := w.Result().Cookies()[0].Value
		claims := new(Claims)
		_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.JWTToken), nil
		})
		require.NoError(t, err)
	})

	t.Run("Failed-Incorrect Pass", func(t *testing.T) {
		t.SkipNow()
	})
}
