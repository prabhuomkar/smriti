package handlers

import (
	"api/internal/models"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type ( // UserRequest ...
	UserRequest struct {
		Name     *string `json:"name"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		Features *string `json:"features"`
	}
)

const (
	queryGetUser    = `SELECT * FROM users WHERE id=$1`
	queryGetUsers   = `SELECT * FROM users ORDER BY created_at DESC OFFSET $1 LIMIT $2`
	queryCreateUser = `INSERT INTO users (id, name, username, password, features, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5, $6, $7)`
	queryUpdateUser = `UPDATE users SET name = $2, username = $3, password = $4, features = $5,` +
		` updated_at = $6 WHERE id=$1`
	queryDeleteUser = `DELETE FROM users WHERE id=$1`
)

// GetUser ...
func (h *Handler) GetUser(ctx echo.Context) error {
	uid, err := getUserID(ctx)
	if err != nil {
		return err
	}
	user := models.User{}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetUser, uid).Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Features, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		slog.Error("error getting user", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	user.UpdatedAt = time.Now()
	_, err = h.DB.Exec(ctx.Request().Context(), queryUpdateUser, user.ID, user.Name, user.Username, user.Password, user.Features, user.UpdatedAt)
	if err != nil {
		slog.Error("error updating user", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// DeleteUser ...
func (h *Handler) DeleteUser(ctx echo.Context) error {
	uid, err := getUserID(ctx)
	if err != nil {
		return err
	}
	_, err = h.DB.Exec(ctx.Request().Context(), queryDeleteUser, uid)
	if err != nil {
		slog.Error("error deleting user", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// GetUsers ...
func (h *Handler) GetUsers(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	users := []models.User{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetUsers, offset, limit)
	if err != nil {
		slog.Error("error getting users", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Features, &user.CreatedAt, &user.UpdatedAt); err != nil {
			slog.Error("error scanning user", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		users = append(users, user)
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
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	_, err = h.DB.Exec(ctx.Request().Context(), queryCreateUser, user.ID, user.Name,
		user.Username, user.Password, user.Features, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		slog.Error("error creating user", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, user)
}

func getUserID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting user id", "error", err)

		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	return uid, err
}

func getUser(ctx echo.Context) (*models.User, error) {
	UserRequest := new(UserRequest)
	if err := ctx.Bind(UserRequest); err != nil {
		slog.Error("error getting user", "error", err)

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
		user.Password = getPasswordHash(*UserRequest.Password)
	}
	if UserRequest.Features != nil {
		user.Features = *UserRequest.Features
	}
	if reflect.DeepEqual(models.User{}, user) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid user")
	}

	return &user, nil
}

func getPasswordHash(password string) string {
	passwordHash := sha512.New()
	passwordHash.Write([]byte(password))

	return hex.EncodeToString(passwordHash.Sum(nil))
}
