package v1

import (
	"net/http"

	models "github.com/dostonshernazarov/resume_maker/api/models"
	pbu "github.com/dostonshernazarov/resume_maker/genproto/user-proto"
	"github.com/dostonshernazarov/resume_maker/internal/pkg/etc"
	l "github.com/dostonshernazarov/resume_maker/internal/pkg/logger"
	"github.com/dostonshernazarov/resume_maker/internal/pkg/otlp"
	tokens "github.com/dostonshernazarov/resume_maker/internal/pkg/token"
	"github.com/dostonshernazarov/resume_maker/internal/pkg/utils"
	valid "github.com/dostonshernazarov/resume_maker/internal/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/encoding/protojson"
)

// CREATE
// @Summary CREATE
// @Security BearerAuth
// @Description Api for Create
// @Tags USER
// @Accept json
// @Produce json
// @Param User body models.UserReq true "createModel"
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/users [post]
func (h *HandlerV1) Create(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "CreateUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body        models.UserReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error binding JSON",
		})
		l.Error(err)
		return
	}

	res := valid.IsValidEmail(body.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect Email. Try again",
		})

		h.Logger.Error("Incorrect Email. Try again, error while in Create")
		return
	}

	res = valid.IsValidPassword(body.Password)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect Password. Try again",
		})

		h.Logger.Error("Incorrect Password. Try again, error while in Create")
		return
	}

	isEmail, err := h.Service.UserService().CheckUniquess(ctx, &pbu.FV{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong",
		})

		h.Logger.Error("Error while check unique email in Create")
		return
	}

	if isEmail.Code != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email already in use",
		})

		return
	}

	password, err := etc.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong",
		})

		h.Logger.Error("Error while hash password in Create")
		return
	}

	newId := uuid.NewString()

	h.JwtHandler = tokens.JwtHandler{
		Sub:       newId,
		Iss:       "client",
		Role:      "user",
		SigninKey: h.Config.Token.SignInKey,
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

	response, err := h.Service.UserService().Create(ctx, &pbu.User{
		Id:           newId,
		FullName:     body.FullName,
		Email:        body.Email,
		Password:     password,
		DateOfBirth:  body.DateOfBirth,
		ProfileImg:   "",
		Card:         body.Card,
		Gender:       body.Gender,
		PhoneNumber:  body.PhoneNumber,
		Role:         "user",
		RefreshToken: refresh,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}

	c.JSON(http.StatusCreated, &models.UserResCreate{
		Id:           response.Id,
		FullName:     response.FullName,
		Email:        response.Email,
		DateOfBirth:  response.DateOfBirth,
		ProfileImg:   response.ProfileImg,
		Card:         response.Card,
		Gender:       response.Gender,
		PhoneNumber:  response.PhoneNumber,
		Role:         response.Role,
		AccessToken:  access,
		RefreshToken: response.RefreshToken,
	})
}

// GET
// @Summary GET
// @Security BearerAuth
// @Description Api for Get
// @Tags USER
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/users/{id} [get]
func (h *HandlerV1) Get(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "GetUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	id := c.Param("id")

	response, err := h.Service.UserService().Get(
		ctx, &pbu.Filter{
			Filter: map[string]string{"id": id},
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}

	if response.User.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Can't get",
		})
		return
	}

	c.JSON(http.StatusOK, &models.UserRes{
		Id:           response.User.Id,
		FullName:     response.User.FullName,
		Email:        response.User.Email,
		DateOfBirth:  response.User.DateOfBirth,
		ProfileImg:   response.User.ProfileImg,
		Card:         response.User.Card,
		Gender:       response.User.Gender,
		PhoneNumber:  response.User.PhoneNumber,
		Role:         response.User.Role,
		RefreshToken: response.User.RefreshToken,
		CreatedAt:    response.User.CreatedAt,
		UpdatedAt:    response.User.UpdatedAt,
		DeletedAt:    response.User.DeletedAt,
	})
}

// LIST USERS
// @Summary LIST USERS
// @Security BearerAuth
// @Description Api for ListUsers
// @Tags USER
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.FieldValues true "request"
// @Success 200 {object} models.ListUsersRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/users/list [get]
func (h *HandlerV1) ListUsers(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		return
	}

	columnQ := c.Query("column")
	valueQ := c.Query("value")

	if columnQ == "" {
		columnQ = "email"
	}

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	response, err := h.Service.UserService().ListUsers(
		ctx, &pbu.ListUsersReq{
			Limit:  params.Limit,
			Offset: (params.Page - 1) * params.Limit,
			Fv: &pbu.FV{
				Field: columnQ,
				Value: valueQ,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// LIST DELETED USERS
// @Summary LIST DELETED USERS
// @Security BearerAuth
// @Description Api for ListDeletedUsers
// @Tags USER
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.FieldValues true "request"
// @Success 200 {object} models.ListUsersRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/users/list/deleted [get]
func (h *HandlerV1) ListDeletedUsers(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListDeletedUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		return
	}

	columnQ := c.Query("column")
	valueQ := c.Query("value")

	if columnQ == "" {
		columnQ = "email"
	}

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	response, err := h.Service.UserService().ListDeletedUsers(
		ctx, &pbu.ListUsersReq{
			Limit:  params.Limit,
			Offset: (params.Page - 1) * params.Limit,
			Fv: &pbu.FV{
				Field: columnQ,
				Value: valueQ,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UPDATE
// @Summary UPDATE
// @Security BearerAuth
// @Description Api for Update
// @Tags USER
// @Accept json
// @Produce json
// @Param User body models.UserReq true "createModel"
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/users [put]
func (h *HandlerV1) Update(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body        models.UserReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	userID, statusCode := GetIdFromToken(c.Request, h.Config)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{
			"error": "Can't get",
		})
		return
	}

	getUser, err := h.Service.UserService().Get(ctx, &pbu.Filter{
		Filter: map[string]string{"id": userID},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Went wrong",
		})
		h.Logger.Error("failed to get user in update", l.Error(err))
		return
	}

	if getUser.User.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Can't update",
		})
		return
	}

	if body.Email != "" {
		emailVal := valid.IsValidEmail(body.Email)
		if !emailVal {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Incorrect Email. Try again",
			})

			h.Logger.Error("Incorrect Email. Try again, error while in update user")
			return
		}
	}

	if body.Password != "" {
		validpas := valid.IsValidPassword(body.Password)
		if !validpas {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Incorrect Password. Try again",
			})

			h.Logger.Error("Incorrect Password. Try again, error while in update user")
			return
		}
		body.Password, err = etc.HashPassword(body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Went wrong",
			})
			h.Logger.Error("failed to hash password in update", l.Error(err))
			return
		}
	}

	if body.FullName == "" {
		body.FullName = getUser.User.FullName
	}

	if body.Email == "" {
		body.Email = getUser.User.Email
	}

	if body.Password == "" {
		body.Password = getUser.User.Password
	}

	if body.DateOfBirth == "" {
		body.DateOfBirth = getUser.User.DateOfBirth
	}

	if body.Card == "" {
		body.Card = getUser.User.Card
	}

	if body.Gender == "" {
		body.Gender = getUser.User.Gender
	}

	if body.PhoneNumber == "" {
		body.PhoneNumber = getUser.User.PhoneNumber
	}

	response, err := h.Service.UserService().Update(ctx, &pbu.User{
		Id:          userID,
		FullName:    body.FullName,
		Email:       body.Email,
		Password:    body.Password,
		DateOfBirth: body.DateOfBirth,
		ProfileImg:  "",
		Card:        body.Card,
		Gender:      body.Gender,
		PhoneNumber: body.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.Logger.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.UserRes{
		Id:           response.Id,
		FullName:     response.FullName,
		Email:        response.Email,
		DateOfBirth:  response.DateOfBirth,
		ProfileImg:   response.ProfileImg,
		Card:         response.Card,
		Gender:       response.Gender,
		PhoneNumber:  response.PhoneNumber,
		Role:         response.Role,
		RefreshToken: response.RefreshToken,
		CreatedAt:    response.CreatedAt,
		UpdatedAt:    response.UpdatedAt,
		DeletedAt:    response.DeletedAt,
	})
}

// DELETE
// @Summary DELETE
// @Security BearerAuth
// @Description Api for Delete
// @Tags USER
// @Accept json
// @Produce json
// @Param id query string true "ID"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/users/{id} [delete]
func (h *HandlerV1) Delete(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "DeleteUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	id := c.Query("id")

	user, err := h.Service.UserService().Get(ctx, &pbu.Filter{
		Filter: map[string]string{"id": id},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Went wrong, error",
		})
		h.Logger.Error("failed to get user in delete", l.Error(err))
		return
	}

	if user.User.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Can't delete",
		})
		return
	}

	_, err = h.Service.UserService().SoftDelete(
		ctx, &pbu.Id{
			Id: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Went wrong, error",
		})
		h.Logger.Error("failed to delete user", l.Error(err))
		return
	}

	// if response != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Went wrong",
	// 	})
	// 	h.Logger.Error("failed to delete user", l.Error(err))
	// 	return
	// }

	c.JSON(http.StatusOK, &models.RegisterRes{
		Content: "User has been deleted",
	})
}

// GET BY TOKEN
// @Summary GET BY TOKEN
// @Security BearerAuth
// @Description Api for Get user by token
// @Tags USER
// @Accept json
// @Produce json
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/users/token [get]
func (h *HandlerV1) GetByToken(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "GetUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	// println("\n", c.Request.Header.Get("Authorization"), "\n")

	userID, statusCode := GetIdFromToken(c.Request, h.Config)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{
			"error": "Can't get",
		})
		return
	}

	response, err := h.Service.UserService().Get(
		ctx, &pbu.Filter{
			Filter: map[string]string{"id": userID},
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}

	c.JSON(http.StatusOK, &models.UserRes{
		Id:           response.User.Id,
		FullName:     response.User.FullName,
		Email:        response.User.Email,
		DateOfBirth:  response.User.DateOfBirth,
		ProfileImg:   response.User.ProfileImg,
		Card:         response.User.Card,
		Gender:       response.User.Gender,
		PhoneNumber:  response.User.PhoneNumber,
		Role:         response.User.Role,
		RefreshToken: response.User.RefreshToken,
		CreatedAt:    response.User.CreatedAt,
		UpdatedAt:    response.User.UpdatedAt,
		DeletedAt:    response.User.DeletedAt,
	})
}
