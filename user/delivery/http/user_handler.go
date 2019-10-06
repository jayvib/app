package http

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/apperr"
	"github.com/jayvib/app/config"
	"github.com/jayvib/app/log"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/user"
	"github.com/jayvib/app/utils/crypto"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func RegisterHandlers(conf *config.Config, r gin.IRouter, us user.Usecase) {
	handler := NewUserHandler(conf, us)
	registerHandlers(r, handler)
}

func registerHandlers(r gin.IRouter, handler *UserHandler) {
	r.GET("/user/:id", handler.GetByID)
	r.POST("/user/update", handler.Update)
	r.POST("/user", handler.Store)
	r.DELETE("/user/:id", handler.Delete)
	r.POST("/user/login", handler.Login)
	r.POST("/user/new", handler.Register)
}

func NewUserHandler(conf *config.Config, u user.Usecase) *UserHandler {
	return &UserHandler{
		config:   conf,
		UUsecase: u,
	}
}

type UserHandler struct {
	config   *config.Config
	UUsecase user.Usecase
}

type RequestError struct {
	Msg        string `json:"msg,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (h *UserHandler) GetByID(c *gin.Context) {
	_id := c.Param("id")
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := h.UUsecase.GetByID(ctx, _id)
	if err != nil {
		log.Infof("%#v", getStatusCode(err))
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) GetByEmail(c *gin.Context) {
	_email := c.Param("email")
	logrus.Println("Email:", _email)
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := h.UUsecase.GetByEmail(ctx, _email)
	if err != nil {
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) GetByUsername(c *gin.Context) {}

func (h *UserHandler) Update(c *gin.Context) {
	var _user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&_user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	logrus.Debugf("%#v\n", _user)
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.UUsecase.Update(ctx, &_user)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Store(c *gin.Context) {
	var _user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&_user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if err := h.UUsecase.Store(ctx, &_user); err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	logrus.Trace(&_user)
	_user.Password = ""
	c.JSON(http.StatusCreated, &_user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	_id := c.Param("id")

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	err := h.UUsecase.Delete(ctx, _id)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Login(c *gin.Context) {
	var creds model.User
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, RequestError{Msg: err.Error()})
		return
	}
	ctx := c.Request.Context()
	u, err := h.UUsecase.GetByUsername(ctx, creds.Username)
	if err != nil {
		c.JSON(getStatusCode(err), RequestError{Msg: err.Error()})
		return
	}

	ok, err := crypto.IsPasswordMatch(u.Password, creds.Password)
	if err != nil {
		c.JSON(getStatusCode(err), RequestError{Msg: err.Error()})
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, RequestError{Msg: "Invalid Password"})
		return
	}

	u.Password = ""

	// Create JWT Claims
	// reference: https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/
	expirationTime := time.Now().Add(3 * time.Hour)
	tk := &Claims{
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(h.config.JWTToken))
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	u.Token = tokenString
	c.SetCookie("token", tokenString,
		int(expirationTime.Unix()), "/", "", false, true)

	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) Register(c *gin.Context) {
	var _user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&_user)
	if err != nil {
		c.JSON(http.StatusBadRequest, RequestError{Msg: err.Error()})
		return
	}
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.UUsecase.Store(ctx, &_user)
	if err != nil {
		c.JSON(getStatusCode(err), RequestError{Msg: err.Error()})
		return
	}

	//_user.Password = ""

	// Create JWT Claims
	expirationTime := time.Now().Add(3 * time.Hour)
	tk := Claims{
		Username: _user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(h.config.JWTToken))
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	_user.Token = tokenString

	// Best Practice to store the JWT token
	// in the HttpOnly Cookie.
	// see https://logrocket.com/blog/jwt-authentication-best-practices/
	c.SetCookie("token", tokenString,
		int(expirationTime.Unix()), "/", "", false, true)

	c.JSON(http.StatusOK, &_user)
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case apperr.ItemNotFound:
		return http.StatusNotFound
	case apperr.UsernameAlreadyExist:
		return http.StatusConflict
	case apperr.EmailAlreadyExist:
		return http.StatusConflict
	default:
		if aerr, ok := err.(apperr.Error); ok {
			switch aerr.Code() {
			case apperr.NoItemFound:
				return http.StatusNotFound
			case apperr.EmptyID:
				return http.StatusBadRequest
			case apperr.ValidationErr:
				return http.StatusBadRequest
			case apperr.BadParameter:
				return http.StatusBadRequest
			}
		}
		return http.StatusInternalServerError
	}
}
