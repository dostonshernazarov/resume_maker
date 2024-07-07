package v1

import (
	"context"

	"github.com/dostonshernazarov/resume_maker/api-service/api/models"
	pbu "github.com/dostonshernazarov/resume_maker/api-service/genproto/user_service"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/etc"
	l "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/logger"
	scode "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/sendcode"
	tokens "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/token"
	val "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/validation"

	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterUser
// REGISTER USER ...
// @Security BearerAuth
// @Router /v1/users/register [POST]
// @Summary REGISTER USER
// @Description Api for register a new user
// @Tags SIGNUP
// @Accept json
// @Produce json
// @Param User body models.RegisterReq true "RegisterUser"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) RegisterUser(c *gin.Context) {
	var body models.RegisterReq
	var toRedis models.ClientRedis

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		},
		)
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Password = strings.TrimSpace(body.Password)
	body.Email = strings.ToLower(body.Email)

	isEmail := val.IsValidEmail(body.Email)
	if !isEmail {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("Incorrect Email. Try again")
		return
	}

	isPassword := val.IsValidPassword(body.Password)
	if !isPassword {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongPassword,
		})

		h.Logger.Error("Password must be at least 8 (numbers and characters) long")
		return
	}

	result, err := h.Service.UserService().UniqueEmail(context.Background(), &pbu.IsUnique{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		},
		)
		h.Logger.Error("Failed to check email uniques", l.Error(err))
		return
	}

	if result.Status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.EmailAlreadyInUse,
		})
		h.Logger.Error("failed to check email unique", l.Error(err))
		return
	}

	// Generate code for check email
	code := strconv.Itoa(rand.Int())[:6]

	toRedis.Code = code
	toRedis.Email = body.Email
	toRedis.Fullname = body.Fullname
	toRedis.Password = body.Password

	userByte, err := json.Marshal(toRedis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to marshal body", l.Error(err))
		return
	}
	err = h.redisStorage.Set(context.Background(), body.Email, userByte, time.Minute*5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to set object to redis", l.Error(err))
		return
	}

	scode.SendCode(body.Email, code)

	responseMessage := models.RegisterRes{
		Content: "We send verification password to your email",
	}

	c.JSON(http.StatusOK, responseMessage)
}

// Verification
// @Security BearerAuth
// @Router /v1/users/verify [GET]
// @Summary VERIFICATION
// @Description Api for verify a new user
// @Tags SIGNUP
// @Accept json
// @Produce json
// @Param request query models.Verify true "request"
// @Success 200 {object} models.UserResCreate
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) Verification(c *gin.Context) {
	email := c.Query("email")
	code := c.Query("code")

	value, err := h.redisStorage.Get(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("Failed to get user from redis", l.Error(err))
		return
	}

	var userDetail models.ClientRedis
	if err := json.Unmarshal([]byte(value), &userDetail); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Error unmarshalling user detail", l.Error(err))
		return
	}

	if userDetail.Code != code {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Error generate new uuid", l.Error(err))
		return
	}

	h.JwtHandler = tokens.JwtHandler{
		Sub:       id.String(),
		Iss:       "client",
		SigninKey: h.Config.Token.SignInKey,
		Role:      "user",
		Log:       h.Logger,
	}

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error generate new jwt tokens", l.Error(err))
		return
	}

	userDetail.Password, err = etc.HashPassword(userDetail.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error in hash password", l.Error(err))
		return
	}

	res, err := h.Service.UserService().CreateUser(context.Background(), &pbu.User{
		Id:       id.String(),
		Name:     userDetail.Fullname,
		Email:    userDetail.Email,
		Refresh:  refresh,
		Password: userDetail.Password,
		Role:     "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error in create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, &models.UserResCreate{
		Id:           res.Guid,
		FullName:     userDetail.Fullname,
		Email:        userDetail.Email,
		ProfileImg:   "",
		PhoneNumber:  "",
		Role:         "user",
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// Login
// @Security BearerAuth
// @Router /v1/users/login [POST]
// @Summary LOGIN
// @Description Api for login user
// @Tags LOGIN
// @Accept json
// @Produce json
// @Param User body models.Login true "Login"
// @Success 200 {object} models.UserResCreate
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) Login(c *gin.Context) {
	var body models.Login

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	email := body.Email
	password := body.Password

	user, err := h.Service.UserService().GetUser(context.Background(), &pbu.Filter{
		Filter: map[string]string{"email": email},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("error while get user in login", l.Error(err))
		return
	}

	if !etc.CheckPasswordHash(password, user.Password) {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		return
	}

	h.JwtHandler = tokens.JwtHandler{
		Sub:       user.Id,
		Role:      user.Role,
		SigninKey: h.Config.Token.SignInKey,
		Log:       h.Logger,
		Timeout:   int(h.Config.Token.AccessTTL),
	}

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error while generate JWT in login", l.Error(err))
		return
	}

	_, err = h.Service.UserService().UpdateRefresh(context.Background(), &pbu.RefreshRequest{
		UserId:       user.Id,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error while update user in login", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.UserResCreate{
		Id:           user.Id,
		FullName:     user.Name,
		Email:        user.Email,
		ProfileImg:   user.Image,
		PhoneNumber:  user.PhoneNumber,
		Role:         user.Role,
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// ForgetPassword
// @Security BearerAuth
// @Router /v1/users/set/{email} [GET]
// @Summary FORGET PASSWORD
// @Description Api for set new password
// @Tags SET-PASSWORD
// @Accept json
// @Produce json
// @Param email query string true "EMAIL"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) ForgetPassword(c *gin.Context) {
	var toRedis models.ForgetPassReq

	email := c.Query("email")

	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	uniqueCheck, err := h.Service.UserService().UniqueEmail(context.Background(), &pbu.IsUnique{
		Email: email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error while check unique in forget password", l.Error(err))
		return
	}

	if !uniqueCheck.Status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		return
	}

	// Generate code for check email
	code := strconv.Itoa(rand.Int())[:6]

	toRedis.Code = code
	toRedis.Email = email

	userByte, err := json.Marshal(toRedis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to marshal body", l.Error(err))
		return
	}
	err = h.redisStorage.Set(context.Background(), toRedis.Email, userByte, time.Minute*10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to set object to redis", l.Error(err))
		return
	}

	scode.SendCode(toRedis.Email, code)

	responseMessage := models.RegisterRes{
		Content: "We send verification password to your email",
	}

	c.JSON(http.StatusOK, responseMessage)
}

// ForgetPasswordVerify
// @Security BearerAuth
// @Router /v1/users/code [GET]
// @Summary FORGET PASSWORD CODE
// @Description Api for verify new password code
// @Tags SET-PASSWORD
// @Accept json
// @Produce json
// @Param request query models.Verify true "request"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) ForgetPasswordVerify(c *gin.Context) {
	var userDetail models.ForgetPassReq

	email := c.Query("email")
	code := c.Query("code")

	value, err := h.redisStorage.Get(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email. Try again ..",
		})
		h.Logger.Error("Failed to get user from redis", l.Error(err))
		return
	}

	if err := json.Unmarshal([]byte(value), &userDetail); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Error unmarshalling user detail in forget password verify", l.Error(err))
		return
	}

	if userDetail.Code != code {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		return
	}

	responseMessage := models.RegisterRes{
		Content: "Please enter new password",
	}

	c.JSON(http.StatusOK, responseMessage)
}

// SetNewPassword
// @Security BearerAuth
// @Router /v1/users/password [PUT]
// @Summary SET NEW PASSWORD
// @Description Api for update new password
// @Tags SET-PASSWORD
// @Accept json
// @Produce json
// @Param request query models.Login true "request"
// @Success 200 {object} models.UserResCreate
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) SetNewPassword(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")

	isPassword := val.IsValidPassword(password)
	if !isPassword {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongPassword,
		})

		h.Logger.Error("Password must be at least 8 (numbers and characters) long")
		return
	}

	user, err := h.Service.UserService().GetUser(context.Background(), &pbu.Filter{
		Filter: map[string]string{"email": email},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to get user from set new password", l.Error(err))
		return
	}

	h.JwtHandler = tokens.JwtHandler{
		Sub:       user.Id,
		Role:      user.Role,
		SigninKey: h.Config.Token.SignInKey,
		Log:       h.Logger,
		Timeout:   int(h.Config.Token.AccessTTL),
	}

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error while generate JWT in login", l.Error(err))
		return
	}

	password, err = etc.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error while hash password in set new password", l.Error(err))
		return
	}

	updUser, err := h.Service.UserService().UpdatePassword(context.Background(), &pbu.UpdatePasswordRequest{
		UserId:      user.Id,
		NewPassword: password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error while hash password in set new password", l.Error(err))
		return
	}

	if updUser.Status {
		c.JSON(http.StatusOK, &models.UserResCreate{
			Id:           user.Id,
			FullName:     user.Name,
			Email:        user.Email,
			ProfileImg:   user.Image,
			PhoneNumber:  user.PhoneNumber,
			Role:         user.Role,
			AccessToken:  access,
			RefreshToken: refresh,
		})
	} else {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
	}

}

// UpdateToken
// @Security BearerAuth
// @Router /v1/token/{refresh} [GET]
// @Summary UPDATE TOKEN
// @Description Api for updated access token
// @Tags TOKEN
// @Accept json
// @Produce json
// @Param refresh path string true "Refresh Token"
// @Success 200 {object} models.TokenResp
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) UpdateToken(c *gin.Context) {
	RToken := c.Param("refresh")

	user, err := h.Service.UserService().GetUser(context.Background(), &pbu.Filter{
		Filter: map[string]string{"refresh": RToken},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to get user in update token", l.Error(err))
		return
	}

	resClaim, err := tokens.ExtractClaim(RToken, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: models.TokenExpired,
		})
		h.Logger.Error("Failed to extract token update token", l.Error(err))
		return
	}

	NowTime := time.Now().Unix()
	exp := resClaim["exp"]
	if exp.(float64)-float64(NowTime) > 0 {
		h.JwtHandler = tokens.JwtHandler{
			Sub:       user.Id,
			Iss:       "client",
			SigninKey: h.Config.Token.SignInKey,
			Role:      user.Role,
			Log:       h.Logger,
		}

		accessR, refreshR, err := h.JwtHandler.GenerateJwt()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error("Failed to generate token update token", l.Error(err))
			return
		}
		_, err = h.Service.UserService().UpdateRefresh(context.Background(), &pbu.RefreshRequest{
			UserId:       user.Id,
			RefreshToken: refreshR,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error("Failed to update user in update token", l.Error(err))
			return
		}

		respUser := &models.TokenResp{
			ID:      user.Id,
			Access:  accessR,
			Refresh: refreshR,
			Role:    user.Role,
		}

		c.JSON(http.StatusOK, respUser)

	} else {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: "Token expired",
		})
		h.Logger.Error("refresh token expired")
		return
	}
}
