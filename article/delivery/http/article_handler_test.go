// +build unit

package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	apperrors "github.com/jayvib/clean-architecture/apperr"
	"github.com/jayvib/clean-architecture/article/mocks"
	"github.com/jayvib/clean-architecture/article/usecase"
	authormocks "github.com/jayvib/clean-architecture/author/mocks"
	"github.com/jayvib/clean-architecture/config"
	jwtmiddleware "github.com/jayvib/clean-architecture/middleware/jwt"
	"github.com/jayvib/clean-architecture/model"
	userhttp "github.com/jayvib/clean-architecture/user/delivery/http"
	usermocks "github.com/jayvib/clean-architecture/user/mocks"
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

var mockAuthor = &model.Author{
	ID:   "unqueid",
	Name: "Luffy Monkey",
}

func setCookieOpt(token string) func(*http.Request) {
	return func(req *http.Request) {
		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: true,
		}
		req.AddCookie(cookie)
	}
}

var mockArticle = &model.Article{
	ID:        "uniqueid",
	Title:     "Pirate King",
	Content:   "Luffy will be the next Pirate King",
	Author:    mockAuthor,
	UpdatedAt: time.Now(),
	CreatedAt: time.Now(),
}

type reqOptionFn func(req *http.Request)

func performRequest(r http.Handler, method, path string, body io.Reader, opt ...reqOptionFn) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if opt != nil {
		for _, o := range opt {
			o(req)
		}
	}
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
	e := gin.Default()
	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)

	t.Run("StatusOK", func(t *testing.T) {
		articleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockArticle, nil).Once()
		authorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockAuthor, nil).Once()
		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, e, u)
		w := performRequest(e, http.MethodGet, "/articles/1", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var articleRes model.Article
		err := json.NewDecoder(w.Body).Decode(&articleRes)
		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, articleRes.Title)
		articleRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("Item Not Found", func(t *testing.T) {
		// TODO: Implement the test
		t.SkipNow()
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		t.SkipNow()
	})
}

func TestGetByIDAuthenticated(t *testing.T) {
	e := gin.Default()
	secret := "dontforgettobeawesome"

	// Create a Group of API

	// Not authentication
	//api := e.Group("/")
	authapi := e.Group("/api/v1")
	authapi.Use(jwtmiddleware.Authenticate(secret))

	token, err := getValidAuthToken()
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		// The User is already authenticated.
		// Simulate request
		articleRepo := new(mocks.Repository)
		authorRepo := new(authormocks.Repository)
		articleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockArticle, nil).Once()
		authorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockAuthor, nil).Once()
		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		articleHandler := ArticleHandler{AUsecase: u}
		authapi.GET("/articles/:id", articleHandler.GetByID)
		setCookieOpt := func(req *http.Request) {
			cookie := &http.Cookie{
				Name:     "token",
				Value:    token,
				Path:     "/",
				Domain:   "",
				Secure:   false,
				HttpOnly: true,
			}
			req.AddCookie(cookie)
		}
		w := performRequest(e, http.MethodGet, "/api/v1/articles/1", nil, setCookieOpt)
		require.Equal(t, http.StatusOK, w.Code)
		var articleRes model.Article
		err = json.NewDecoder(w.Body).Decode(&articleRes)
		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, articleRes.Title)
		articleRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("Authentication Failed", func(t *testing.T) {
		t.SkipNow()
	})
}

func getValidAuthToken() (string, error) {
	conf, err := config.New()
	if err != nil {
		return "", err
	}
	// Do Registration simulation
	e := gin.Default()
	userRepo := new(usermocks.Repository)
	authorRepo := new(authormocks.Repository)
	user := &model.User{
		Firstname: "Luffy",
		Lastname:  "Monkey",
		Email:     "luffy.monkey@gmail.com",
		Username:  "luff.monkey",
		Password:  "pirateking",
	}
	userRepo.On("Store", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	userRepo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(nil, apperrors.ItemNotFound).Once()
	userRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, apperrors.ItemNotFound).Once()
	authorRepo.On("Store", mock.Anything, mock.AnythingOfType("*model.Author")).Return(nil).Once()
	userUsecase := userusecase.New(userRepo, authorRepo)
	userHandler := userhttp.NewUserHandler(conf, userUsecase)
	e.POST("/user/new", userHandler.Register)
	payload, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	payloadReader := bytes.NewReader(payload)
	w := performRequest(e, http.MethodPost, "/user/new", payloadReader)
	if w.Code != http.StatusOK {
		return "", errors.New("request failed")
	}
	var resUser model.User
	err = json.NewDecoder(w.Body).Decode(&resUser)
	if err != nil {
		return "", err
	}
	return w.Result().Cookies()[0].Value, nil
}

func TestStore(t *testing.T) {
	e := gin.Default()
	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)
	t.Run("StatusOK", func(t *testing.T) {
		articleRepo.On("Store",
			mock.Anything, mock.AnythingOfType("*model.Article")).Return(nil).Once()
		articleRepo.On("GetByTitle",
			mock.Anything, mock.AnythingOfType("string")).Return(nil, apperrors.ItemNotFound)

		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, e, u)
		copyArticle := &(*mockArticle)
		copyArticle.ID = "unqueid"
		payload, err := json.Marshal(mockArticle)
		require.NoError(t, err)
		w := performRequest(e, http.MethodPost, "/articles", bytes.NewReader(payload))
		assert.Equal(t, http.StatusCreated, w.Code)

		var resArticle model.Article
		err = json.NewDecoder(w.Body).Decode(&resArticle)
		assert.NoError(t, err)
		assert.NotEmpty(t, resArticle.ID)
		articleRepo.AssertExpectations(t)
	})
	t.Run("Already Exist", func(t *testing.T) {
		// TODO: Need to implement the test
		t.SkipNow()
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		// TODO: Need to implement the test
		t.SkipNow()
	})
}

func TestStoreAuthenticated(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)
	e := gin.Default()
	authapi := e.Group("/api/v1")
	authapi.Use(jwtmiddleware.Authenticate(conf.JWTToken))
	token, err := getValidAuthToken()
	require.NoError(t, err)
	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)

	t.Run("Success", func(t *testing.T) {
		articleRepo.On("Store",
			mock.Anything, mock.AnythingOfType("*model.Article")).Return(nil).Once()
		articleRepo.On("GetByTitle",
			mock.Anything, mock.AnythingOfType("string")).Return(nil, apperrors.ItemNotFound)

		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, authapi, u)
		copyArticle := &(*mockArticle)
		copyArticle.ID = "unqueid"
		payload, err := json.Marshal(mockArticle)
		require.NoError(t, err)
		w := performRequest(e, http.MethodPost, "/api/v1/articles", bytes.NewReader(payload), setCookieOpt(token))
		assert.Equal(t, http.StatusCreated, w.Code)
		var resArticle model.Article
		err = json.NewDecoder(w.Body).Decode(&resArticle)
		assert.NoError(t, err)
		assert.NotEmpty(t, resArticle.ID)
		articleRepo.AssertExpectations(t)
	})
	t.Run("Invalid Token", func(t *testing.T) {
		token := "invalidtoken"
		copyArticle := &(*mockArticle)
		copyArticle.ID = "uniqueid"
		payload, err := json.Marshal(mockArticle)
		require.NoError(t, err)
		w := performRequest(e, http.MethodPost, "/api/v1/articles", bytes.NewReader(payload), setCookieOpt(token))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUpdate(t *testing.T) {
	e := gin.Default()
	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)
	t.Run("StatusOK", func(t *testing.T) {
		articleRepo.On("Update",
			mock.Anything, mock.AnythingOfType("*model.Article")).Return(nil).Once()

		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, e, u)
		payload, err := json.Marshal(mockArticle)
		require.NoError(t, err)
		w := performRequest(e, http.MethodPost, "/articles/update", bytes.NewReader(payload))
		assert.Equal(t, http.StatusOK, w.Code)
		articleRepo.AssertExpectations(t)
	})
	t.Run("unexpected error", func(t *testing.T) {
		// TODO: Immplement the test
		t.SkipNow()
	})
}

func TestUpdateAuthenticated(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)
	e := gin.Default()
	authapi := e.Group("/api/v1")
	authapi.Use(jwtmiddleware.Authenticate(conf.JWTToken))
	token, err := getValidAuthToken()
	require.NoError(t, err)

	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)
	t.Run("Success", func(t *testing.T) {
		articleRepo.On("Update",
			mock.Anything, mock.AnythingOfType("*model.Article")).Return(nil).Once()

		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, authapi, u)
		payload, err := json.Marshal(mockArticle)
		require.NoError(t, err)
		w := performRequest(e, http.MethodPost, "/api/v1/articles/update", bytes.NewReader(payload), setCookieOpt(token))
		assert.Equal(t, http.StatusOK, w.Code)
		articleRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		token := "whosyourdaddy?"
		payload, err := json.Marshal(mockArticle)
		require.NoError(t, err)
		w := performRequest(e, http.MethodPost, "/api/v1/articles/update", bytes.NewReader(payload), setCookieOpt(token))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestDelete(t *testing.T) {
	e := gin.Default()
	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)
	t.Run("StatusOK", func(t *testing.T) {
		articleRepo.On("Delete",
			mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
		articleRepo.On("GetByID",
			mock.Anything, mock.AnythingOfType("string")).Return(mockArticle, nil).Once()
		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, e, u)
		w := performRequest(e, http.MethodDelete, "/articles/1", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		articleRepo.AssertExpectations(t)
	})
}

func TestDeleteAuthenticated(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)
	// TODO(JAYSON): Implement the test.
	e := gin.Default()
	authapi := e.Group("/api/v1")
	authapi.Use(jwtmiddleware.Authenticate(conf.JWTToken))
	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)
	token, err := getValidAuthToken()
	require.NoError(t, err)
	t.Run("StatusOK", func(t *testing.T) {
		articleRepo.On("Delete",
			mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
		articleRepo.On("GetByID",
			mock.Anything, mock.AnythingOfType("string")).Return(mockArticle, nil).Once()
		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, authapi, u)
		w := performRequest(e, http.MethodDelete, "/api/v1/articles/1", nil, setCookieOpt(token))
		assert.Equal(t, http.StatusOK, w.Code)
		articleRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		token := "whosyourdaddy?"
		w := performRequest(e, http.MethodDelete, "/api/v1/articles/1", nil, setCookieOpt(token))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestFetch(t *testing.T) {
	mockArticles := []*model.Article{
		mockArticle,
	}
	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)
	t.Run("StatusOK", func(t *testing.T) {
		//logrus.SetLevel(logrus.DebugLevel)
		num := 1
		cursor := "2"
		articleRepo.On("Fetch", mock.Anything, cursor, num).Return(mockArticles, "10", nil).Once()
		authorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockAuthor, nil).Once()
		e := gin.Default()
		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, e, u)

		// Perform Request
		w := performRequest(e, http.MethodGet, "/articles?num=1&cursor="+cursor, nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "10", w.Header().Get("X-Cursor"))
		articleRepo.AssertExpectations(t)
	})

	t.Run("unexpected error", func(t *testing.T) {
		// TODO: Implement the test
		t.SkipNow()
	})
}

func TestFetchAuthenticated(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)
	e := gin.Default()
	authapi := e.Group("/api/v1")
	authapi.Use(jwtmiddleware.Authenticate(conf.JWTToken))
	token, err := getValidAuthToken()
	require.NoError(t, err)

	mockArticles := []*model.Article{
		mockArticle,
	}
	articleRepo := new(mocks.Repository)
	authorRepo := new(authormocks.Repository)
	//logrus.SetLevel(logrus.DebugLevel)
	num := 1
	cursor := "2"
	t.Run("StatusOK", func(t *testing.T) {
		articleRepo.On("Fetch", mock.Anything, cursor, num).Return(mockArticles, "10", nil).Once()
		authorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockAuthor, nil).Once()
		u := usecase.New(articleRepo, authorRepo, time.Second*2)
		RegisterHandlers(nil, authapi, u)

		// Perform Request
		w := performRequest(e, http.MethodGet, "/api/v1/articles?num=1&cursor="+cursor, nil, setCookieOpt(token))
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "10", w.Header().Get("X-Cursor"))
		articleRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		token := "whosyourdaddy?"
		w := performRequest(e, http.MethodGet, "/api/v1/articles?num=1&cursor="+cursor, nil, setCookieOpt(token))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
