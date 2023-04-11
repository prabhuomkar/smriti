package handlers

import (
	"api/internal/models"
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	// UserRequest ...
	UserRequest struct {
		Name     *string `json:"name"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		Features *string `json:"features"`
	}
)

// GetUser ...
func (h *Handler) GetUser(ctx echo.Context) error {
	uid, err := getUserID(ctx)
	if err != nil {
		return err
	}
	user := models.User{}
	result := h.DB.Model(&models.User{}).Where("id=?", uid).First(&user)
	if result.Error != nil {
		log.Printf("error getting user: %+v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, user)
}

// UpdateUser ...
func (h *Handler) UpdateUser(ctx echo.Context) error {
	uid, err := getUserID(ctx)
	if err != nil {
		return err
	}
	user, err := getUser(ctx)
	if err != nil {
		return err
	}
	user.ID = uid
	result := h.DB.Model(&user).Updates(user)
	if result.Error != nil || result.RowsAffected != 1 {
		log.Printf("error updating user: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// DeleteUser ...
func (h *Handler) DeleteUser(ctx echo.Context) error {
	uid, err := getUserID(ctx)
	if err != nil {
		return err
	}
	user := models.User{ID: uid}
	result := h.DB.Delete(&user)
	if result.Error != nil || result.RowsAffected != 1 {
		log.Printf("error deleting user: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetUsers ...
func (h *Handler) GetUsers(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	users := []models.User{}
	result := h.DB.Model(&models.User{}).
		Find(&users).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		log.Printf("error getting users: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, users)
}

// CreateUser ...
func (h *Handler) CreateUser(ctx echo.Context) error {
	user, err := getUser(ctx)
	if err != nil {
		return err
	}
	user.ID = uuid.NewV4()
	if result := h.DB.Create(&user); result.Error != nil {
		log.Printf("error creating user: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusCreated, user)
}

func getUserID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting user id: %+v", err)
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}
	return uid, err
}

func getUser(ctx echo.Context) (*models.User, error) {
	UserRequest := new(UserRequest)
	if err := ctx.Bind(UserRequest); err != nil {
		log.Printf("error getting user: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid user")
	}
	user := models.User{}
	if UserRequest.Name != nil {
		user.Name = *UserRequest.Name
	}
	if UserRequest.Username != nil {
		user.Username = *UserRequest.Username
	}
	if UserRequest.Password != nil {
		user.Password = *UserRequest.Password
	}
	if UserRequest.Features != nil {
		user.Features = *UserRequest.Features
	}
	if reflect.DeepEqual(models.User{}, user) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid user")
	}
	return &user, nil
}
