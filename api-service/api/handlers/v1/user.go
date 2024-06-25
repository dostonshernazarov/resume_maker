package v1

import (
	"context"
	"net/http"

	"github.com/dostonshernazarov/resume_maker/api-service/api/models"
	pbu "github.com/dostonshernazarov/resume_maker/api-service/genproto/user_service"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/etc"
	l "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/logger"
	tokens "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/token"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/utils"
	valid "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateUser
// @Summary CREATE
// @Security BearerAuth
// @Description Api for Create
// @Tags USER
// @Accept json
// @Produce json
// @Param User body models.UserReq true "createModel"
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users [post]
func (h *HandlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.UserReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		l.Error(err)
		return
	}

	res := valid.IsValidEmail(body.Email)
	if !res {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})

		h.Logger.Error("Incorrect Email. Try again, error while in Create")
		return
	}

	res = valid.IsValidPassword(body.Password)
	if !res {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongPassword,
		})

		h.Logger.Error("Incorrect Password. Try again, error while in Create")
		return
	}

	isEmail, err := h.Service.UserService().UniqueEmail(context.Background(), &pbu.IsUnique{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})

		h.Logger.Error("Error while check unique email in Create")
		return
	}

	if isEmail.Status {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.EmailAlreadyInUse,
		})

		return
	}

	password, err := etc.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
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
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("error generate new jwt tokens", l.Error(err))
		return
	}

	response, err := h.Service.UserService().CreateUser(context.Background(), &pbu.User{
		Id:          newId,
		Name:        body.FullName,
		Email:       body.Email,
		Password:    password,
		PhoneNumber: body.PhoneNumber,
		Role:        "user",
		Refresh:     refresh,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		l.Error(err)
		return
	}

	c.JSON(http.StatusCreated, &models.UserResCreate{
		Id:           response.Guid,
		FullName:     body.FullName,
		Email:        body.Email,
		PhoneNumber:  body.PhoneNumber,
		Role:         "user",
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// GetUser
// @Summary GET
// @Security BearerAuth
// @Description Api for Get
// @Tags USER
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/{id} [get]
func (h *HandlerV1) GetUser(c *gin.Context) {
	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	id := c.Param("id")

	response, err := h.Service.UserService().GetUser(
		context.Background(), &pbu.Filter{
			Filter: map[string]string{"id": id},
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}

	if response.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Can't get",
		})
		return
	}

	c.JSON(http.StatusOK, &models.UserRes{
		Id:           response.Id,
		FullName:     response.Name,
		Email:        response.Email,
		PhoneNumber:  response.PhoneNumber,
		Role:         response.Role,
		RefreshToken: response.Refresh,
		CreatedAt:    response.CreatedAt,
		UpdatedAt:    response.UpdatedAt,
	})
}

// ListUsers
// @Summary LIST USERS
// @Security BearerAuth
// @Description Api for ListUsers
// @Tags USER
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Success 200 {object} models.Users
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/list [get]
func (h *HandlerV1) ListUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		return
	}

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	response, err := h.Service.UserService().GetAllUsers(
		context.Background(), &pbu.ListUserRequest{
			Limit: int64(params.Limit),
			Page:  int64(params.Page),
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		l.Error(err)
		return
	}

	var users models.Users
	for _, val := range response.Users {
		var respUser models.UserRes
		respUser.Id = val.Id
		respUser.FullName = val.Name
		respUser.Email = val.Email
		respUser.PhoneNumber = val.PhoneNumber
		respUser.CreatedAt = val.CreatedAt
		respUser.Role = val.Role
		respUser.RefreshToken = val.Refresh
		respUser.UpdatedAt = val.UpdatedAt

		users.Users = append(users.Users, &respUser)
	}

	users.Count = response.TotalCount

	c.JSON(http.StatusOK, users)
}

// UpdateUser
// @Summary UPDATE
// @Security BearerAuth
// @Description Api for Update
// @Tags USER
// @Accept json
// @Produce json
// @Param User body models.UserReq true "createModel"
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users [put]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	var (
		body        models.UserUpdateReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	userID, statusCode := GetIdFromToken(c.Request, h.Config)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, models.Error{
			Message: models.InternalMessage,
		})
		return
	}

	getUser, err := h.Service.UserService().GetUser(context.Background(), &pbu.Filter{
		Filter: map[string]string{"id": userID},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("failed to get user in update", l.Error(err))
		return
	}

	if getUser.Role != "user" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		return
	}

	if body.Email != "" {
		emailVal := valid.IsValidEmail(body.Email)
		if !emailVal {
			c.JSON(http.StatusBadRequest, models.Error{
				Message: models.NotAvailable,
			})

			h.Logger.Error("Incorrect Email. Try again, error while in update user")
			return
		}
	}

	response, err := h.Service.UserService().UpdateUser(context.Background(), &pbu.User{
		Id:          userID,
		Name:        body.FullName,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.UserRes{
		Id:           response.Id,
		FullName:     response.Name,
		Email:        response.Email,
		PhoneNumber:  response.PhoneNumber,
		Role:         response.Role,
		RefreshToken: response.Refresh,
		CreatedAt:    response.CreatedAt,
		UpdatedAt:    response.UpdatedAt,
	})
}

// DeleteUser
// @Summary DELETE
// @Security BearerAuth
// @Description Api for Delete
// @Tags USER
// @Accept json
// @Produce json
// @Param id query string true "ID"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/{id} [delete]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	id := c.Query("id")

	user, err := h.Service.UserService().GetUser(context.Background(), &pbu.Filter{
		Filter: map[string]string{"id": id},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to get user in delete", l.Error(err))
		return
	}

	if user.Role != "user" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		return
	}

	_, err = h.Service.UserService().DeleteUser(
		context.Background(), &pbu.UserWithGUID{
			Guid: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
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

// GetByToken
// @Summary GET BY TOKEN
// @Security BearerAuth
// @Description Api for Get user by token
// @Tags USER
// @Accept json
// @Produce json
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/token [get]
func (h *HandlerV1) GetByToken(c *gin.Context) {
	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	userID, statusCode := GetIdFromToken(c.Request, h.Config)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, models.Error{
			Message: models.WrongInfoMessage,
		})
		return
	}

	response, err := h.Service.UserService().GetUser(
		context.Background(), &pbu.Filter{
			Filter: map[string]string{"id": userID},
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		l.Error(err)
		return
	}

	c.JSON(http.StatusOK, &models.UserRes{
		Id:           response.Id,
		FullName:     response.Name,
		Email:        response.Email,
		PhoneNumber:  response.PhoneNumber,
		Role:         response.Role,
		RefreshToken: response.Refresh,
		CreatedAt:    response.CreatedAt,
		UpdatedAt:    response.UpdatedAt,
	})
}
