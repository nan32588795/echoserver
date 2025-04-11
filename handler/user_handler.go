package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"echoserver/repository"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	u := new(repository.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	if err := repository.CreateUser(h.DB, u); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := repository.GetAllUsers(h.DB)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := repository.GetUserByID(h.DB, id)
	if err != nil {
		return err
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	u := new(repository.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	if err := repository.UpdateUser(h.DB, id, u); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User updated"})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repository.DeleteUser(h.DB, id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
}
