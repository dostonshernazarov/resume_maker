package v1

import (
	"github.com/dostonshernazarov/resume_maker/api/models"
	pbu "github.com/dostonshernazarov/resume_maker/genproto/user-proto"
	"github.com/dostonshernazarov/resume_maker/internal/pkg/etc"
	l "github.com/dostonshernazarov/resume_maker/internal/pkg/logger"
	"github.com/dostonshernazarov/resume_maker/internal/pkg/otlp"
	scode "github.com/dostonshernazarov/resume_maker/internal/pkg/sendcode"
	tokens "github.com/dostonshernazarov/resume_maker/internal/pkg/token"
	val "github.com/dostonshernazarov/resume_maker/internal/pkg/validation"

	// "context"
	"encoding/json"
	// "errors"
	"math/rand"
	"strconv"
	"time"

	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	// "github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
)

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
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) RegisterUser(c *gin.Context) {

	ctx, span := otlp.Start(c, "api", "Register")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var body models.RegisterReq
	var toRedis models.ClientRedis

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Password = strings.TrimSpace(body.Password)
	body.Email = strings.ToLower(body.Email)

	isEmail := val.IsValidEmail(body.Email)
	if !isEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect Email. Try again",
		})

		h.Logger.Error("Incorrect Email. Try again")
		return
	}

	isPassword := val.IsValidPassword(body.Password)
	if !isPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 (numbers and characters) long",
		})

		h.Logger.Error("Password must be at least 8 (numbers and characters) long")
		return
	}

	result, err := h.Service.UserService().CheckUniquess(ctx, &pbu.FV{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("Failed to check email uniquess", l.Error(err))
		return
	}

	if result.Code == 1 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already in use, please use another email address",
		})
		h.Logger.Error("failed to check email unique", l.Error(err))
		return
	}

	// Connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	// Generate code for check email
	code := strconv.Itoa(rand.Int())[:6]

	toRedis.Code = code
	toRedis.Email = body.Email
	toRedis.Fullname = body.Fullname
	toRedis.Password = body.Password

	userByte, err := json.Marshal(toRedis)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("Failed to marshal body", l.Error(err))
		return
	}
	_, err = rdb.Set(ctx, body.Email, userByte, time.Minute*3).Result()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("Failed to set object to redis", l.Error(err))
		return
	}

	scode.SendCode(body.Email, code)

	responsemessage := models.RegisterRes{
		Content: "We send verification password to your email",
	}

	c.JSON(http.StatusOK, responsemessage)
}

// VERIFICATION ...
// @Security BearerAuth
// @Router /v1/users/verify [GET]
// @Summary VERIFICATION
// @Description Api for verify a new user
// @Tags SIGNUP
// @Accept json
// @Produce json
// @Param request query models.Verify true "request"
// @Success 200 {object} models.UserResCreate
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) Verification(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "Verification")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	email := c.Query("email")
	code := c.Query("code")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	val, err := rdb.Get(ctx, email).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email. Try again ..",
		})
		h.Logger.Error("Failed to get user from redis", l.Error(err))
		return
	}

	var userdetail models.ClientRedis
	if err := json.Unmarshal([]byte(val), &userdetail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unmarshiling error",
		})
		h.Logger.Error("Error unmarshalling userdetail", l.Error(err))
		return
	}

	if userdetail.Code != code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect code. Try again",
		})
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while generating uuid",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while generating jwt",
		})
		h.Logger.Error("error generate new jwt tokens", l.Error(err))
		return
	}

	userdetail.Password, err = etc.HashPassword(userdetail.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Oops. Something went wrong with password",
		})
		h.Logger.Error("error in hash password", l.Error(err))
		return
	}

	res, err := h.Service.UserService().Create(ctx, &pbu.User{
		Id:           id.String(),
		FullName:     userdetail.Fullname,
		Email:        userdetail.Email,
		Password:     userdetail.Password,
		DateOfBirth:  "",
		ProfileImg:   "",
		Card:         "",
		Gender:       "",
		PhoneNumber:  "",
		Role:         "user",
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while create user",
		})
		h.Logger.Error("error in create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, &models.UserResCreate{
		Id:           res.Id,
		FullName:     res.FullName,
		Email:        res.Email,
		DateOfBirth:  res.DateOfBirth,
		ProfileImg:   res.ProfileImg,
		Card:         res.Card,
		Gender:       res.Gender,
		PhoneNumber:  res.PhoneNumber,
		Role:         res.Role,
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// LOGIN ...
// @Security BearerAuth
// @Router /v1/users/login [POST]
// @Summary LOGIN
// @Description Api for login user
// @Tags LOGIN
// @Accept json
// @Produce json
// @Param User body models.Login true "Login"
// @Success 200 {object} models.UserResCreate
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) Login(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "Login")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var body models.Login

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	email := body.Email
	password := body.Password

	user, err := h.Service.UserService().Get(ctx, &pbu.Filter{
		Filter: map[string]string{"email": email},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email or password",
		})
		h.Logger.Error("error while get user in login", l.Error(err))
		return
	}

	// println("\n\n", user.User.Email, "\n\n")

	if !etc.CheckPasswordHash(password, user.User.Password) {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Incorrect email or password",
		})
		return
	}

	h.JwtHandler = tokens.JwtHandler{
		Sub:       user.User.Id,
		Role:      user.User.Role,
		SigninKey: h.Config.Token.SignInKey,
		Log:       h.Logger,
		Timeout:   int(h.Config.Token.AccessTTL),
	}

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong",
		})
		h.Logger.Error("error while generate JWT in login", l.Error(err))
		return
	}

	_, err = h.Service.UserService().Update(ctx, &pbu.User{
		Id:           user.User.Id,
		FullName:     user.User.FullName,
		Email:        user.User.Email,
		Password:     user.User.Password,
		DateOfBirth:  user.User.DateOfBirth,
		ProfileImg:   user.User.ProfileImg,
		Card:         user.User.Card,
		Gender:       user.User.Gender,
		PhoneNumber:  user.User.PhoneNumber,
		Role:         user.User.Role,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong",
		})
		h.Logger.Error("error while update user in login", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.UserResCreate{
		Id:           user.User.Id,
		FullName:     user.User.FullName,
		Email:        user.User.Email,
		DateOfBirth:  user.User.DateOfBirth,
		ProfileImg:   user.User.ProfileImg,
		Card:         user.User.Card,
		Gender:       user.User.Gender,
		PhoneNumber:  user.User.PhoneNumber,
		Role:         user.User.Role,
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// LOGIN ADMIN...
// @Security BearerAuth
// @Router /v1/admins/login [POST]
// @Summary LOGIN
// @Description Api for login admin
// @Tags LOGIN
// @Accept json
// @Produce json
// @Param User body models.Login true "Login"
// @Success 200 {object} models.UserResCreate
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) LoginAdmin(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "Login")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var body models.Login

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	email := body.Email
	password := body.Password

	user, err := h.Service.UserService().Get(ctx, &pbu.Filter{
		Filter: map[string]string{"email": email},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email or password",
		})
		h.Logger.Error("error while get user in login admin", l.Error(err))
		return
	}

	if user.User.Role != "admin" {
		if user.User.Role != "sudo" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Permission denied",
			})
			h.Logger.Error("Role not admin")
			return
		}
	}
	// println("\n\n", user.User.Email, "\n\n")

	if !etc.CheckPasswordHash(password, user.User.Password) {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Incorrect email or password",
		})
		return
	}

	h.JwtHandler = tokens.JwtHandler{
		Sub:       user.User.Id,
		Role:      user.User.Role,
		SigninKey: h.Config.Token.SignInKey,
		Log:       h.Logger,
		Timeout:   int(h.Config.Token.AccessTTL),
	}

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong",
		})
		h.Logger.Error("error while generate JWT in login", l.Error(err))
		return
	}

	_, err = h.Service.UserService().Update(ctx, &pbu.User{
		Id:           user.User.Id,
		FullName:     user.User.FullName,
		Email:        user.User.Email,
		Password:     user.User.Password,
		DateOfBirth:  user.User.DateOfBirth,
		ProfileImg:   user.User.ProfileImg,
		Card:         user.User.Card,
		Gender:       user.User.Gender,
		PhoneNumber:  user.User.PhoneNumber,
		Role:         user.User.Role,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong",
		})
		h.Logger.Error("error while update user in login", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.UserResCreate{
		Id:           user.User.Id,
		FullName:     user.User.FullName,
		Email:        user.User.Email,
		DateOfBirth:  user.User.DateOfBirth,
		ProfileImg:   user.User.ProfileImg,
		Card:         user.User.Card,
		Gender:       user.User.Gender,
		PhoneNumber:  user.User.PhoneNumber,
		Role:         user.User.Role,
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// FORGET PASSWORD ...
// @Security BearerAuth
// @Router /v1/users/set/{email} [GET]
// @Summary FORGET PASSWORD
// @Description Api for set new password
// @Tags SET-PASSWORD
// @Accept json
// @Produce json
// @Param email query string true "EMAIL"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) ForgetPassword(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ForgetPassword")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var toRedis models.ForgetPassReq

	email := c.Query("email")

	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	// println("\n\n", email, "\n")
	uniqueCheck, err := h.Service.UserService().CheckUniquess(ctx, &pbu.FV{
		Field: "email",
		Value: email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong",
		})
		h.Logger.Error("error while check unique in forget password", l.Error(err))
		return
	}

	if uniqueCheck.Code == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect email",
		})
		return
	}

	// Connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	// Generate code for check email
	code := strconv.Itoa(rand.Int())[:6]

	toRedis.Code = code
	toRedis.Email = email

	userByte, err := json.Marshal(toRedis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("Failed to marshal body", l.Error(err))
		return
	}
	_, err = rdb.Set(ctx, toRedis.Email, userByte, time.Minute*10).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("Failed to set object to redis", l.Error(err))
		return
	}

	scode.SendCode(toRedis.Email, code)

	responsemessage := models.RegisterRes{
		Content: "We send verification password to your email",
	}

	c.JSON(http.StatusOK, responsemessage)
}

// FORGET PASSWORD CODE ...
// @Security BearerAuth
// @Router /v1/users/code [GET]
// @Summary FORGET PASSWORD CODE
// @Description Api for verify new password code
// @Tags SET-PASSWORD
// @Accept json
// @Produce json
// @Param request query models.Verify true "request"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) ForgetPasswordVerify(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ForgetPassword")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var userdetail models.ForgetPassReq

	email := c.Query("email")
	code := c.Query("code")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	val, err := rdb.Get(ctx, email).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email. Try again ..",
		})
		h.Logger.Error("Failed to get user from redis", l.Error(err))
		return
	}

	if err := json.Unmarshal([]byte(val), &userdetail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unmarshiling error",
		})
		h.Logger.Error("Error unmarshalling userdetail in forget password verify", l.Error(err))
		return
	}

	if userdetail.Code != code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect code. Try again",
		})
		return
	}

	responsemessage := models.RegisterRes{
		Content: "Please enter new password",
	}

	c.JSON(http.StatusOK, responsemessage)
}

// SET NEW PASSWORD ...
// @Security BearerAuth
// @Router /v1/users/password [PUT]
// @Summary SET NEW PASSWORD
// @Description Api for update new password
// @Tags SET-PASSWORD
// @Accept json
// @Produce json
// @Param request query models.Login true "request"
// @Success 200 {object} models.UserResCreate
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) SetNewPassword(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "SetNewPassword")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	email := c.Query("email")
	password := c.Query("password")

	isPassword := val.IsValidPassword(password)
	if !isPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 (numbers and characters) long",
		})

		h.Logger.Error("Password must be at least 8 (numbers and characters) long")
		return
	}

	user, err := h.Service.UserService().Get(ctx, &pbu.Filter{
		Filter: map[string]string{"email": email},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email. Try again ..",
		})
		h.Logger.Error("Failed to get user from set new password", l.Error(err))
		return
	}

	h.JwtHandler = tokens.JwtHandler{
		Sub:       user.User.Id,
		Role:      user.User.Role,
		SigninKey: h.Config.Token.SignInKey,
		Log:       h.Logger,
		Timeout:   int(h.Config.Token.AccessTTL),
	}

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong",
		})
		h.Logger.Error("error while generate JWT in login", l.Error(err))
		return
	}

	password, err = etc.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Oops. Something went wrong with password",
		})
		h.Logger.Error("error while hash password in set new password", l.Error(err))
		return
	}

	updUser, err := h.Service.UserService().Update(ctx, &pbu.User{
		Id:           user.User.Id,
		FullName:     user.User.FullName,
		Email:        user.User.Email,
		Password:     password,
		DateOfBirth:  user.User.DateOfBirth,
		ProfileImg:   user.User.ProfileImg,
		Card:         user.User.Card,
		Gender:       user.User.Gender,
		PhoneNumber:  user.User.PhoneNumber,
		Role:         user.User.Role,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Oops. Something went wrong with password",
		})
		h.Logger.Error("error while hash password in set new password", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.UserResCreate{
		Id:           updUser.Id,
		FullName:     updUser.FullName,
		Email:        updUser.Email,
		DateOfBirth:  updUser.DateOfBirth,
		ProfileImg:   updUser.ProfileImg,
		Card:         updUser.Card,
		Gender:       updUser.Gender,
		PhoneNumber:  updUser.PhoneNumber,
		Role:         updUser.Role,
		AccessToken:  access,
		RefreshToken: updUser.RefreshToken,
	})
}

// UPDATE TOKEN
// @Security BearerAuth
// @Router /v1/token/{refresh} [GET]
// @Summary UPDATE TOKEN
// @Description Api for updated acces token
// @Tags TOKEN
// @Accept json
// @Produce json
// @Param refresh path string true "Refresh Token"
// @Success 200 {object} models.TokenResp
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) UpdateToken(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "SetNewPassword")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	RToken := c.Param("refresh")

	// println("\n\n", RToken, "\n\n")

	user, err := h.Service.UserService().Get(ctx, &pbu.Filter{
		Filter: map[string]string{"refresh_token": RToken},
	})
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Incorrect token.",
		})
		h.Logger.Error("Failed to get user in update token", l.Error(err))
		return
	}

	resClaim, err := tokens.ExtractClaim(RToken, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Went wrong, 1",
		})
		h.Logger.Error("Failed to extract token update token", l.Error(err))
		return
	}

	Now_time := time.Now().Unix()
	exp := (resClaim["exp"])
	if exp.(float64)-float64(Now_time) > 0 {
		h.JwtHandler = tokens.JwtHandler{
			Sub:       user.User.Id,
			Iss:       "client",
			SigninKey: h.Config.Token.SignInKey,
			Role:      user.User.Role,
			Log:       h.Logger,
		}

		accessR, refreshR, err := h.JwtHandler.GenerateJwt()
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Went wrong",
			})
			h.Logger.Error("Failed to generate token update token", l.Error(err))
			return
		}
		_, err = h.Service.UserService().Update(ctx, &pbu.User{
			Id:           user.User.Id,
			FullName:     user.User.FullName,
			Email:        user.User.Email,
			Password:     user.User.Password,
			DateOfBirth:  user.User.DateOfBirth,
			ProfileImg:   user.User.ProfileImg,
			Card:         user.User.Card,
			Gender:       user.User.Gender,
			PhoneNumber:  user.User.PhoneNumber,
			Role:         user.User.Role,
			RefreshToken: refreshR,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Went wrong",
			})
			h.Logger.Error("Failed to update user in update token", l.Error(err))
			return
		}

		respUser := &models.TokenResp{
			ID:      user.User.Id,
			Access:  accessR,
			Refresh: refreshR,
			Role:    user.User.Role,
		}

		c.JSON(http.StatusOK, respUser)

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token expired",
		})
		h.Logger.Error("refresh token expired")
		return
	}

}
