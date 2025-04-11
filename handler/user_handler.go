package handler

import (
	"net/http"
	"strconv"

	"echoserver/models"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	err := u.Insert(c.Request().Context(), boil.GetContextDB(), boil.Infer())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "insert failed"})
	}
	return c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := models.Users(qm.OrderBy("id")).All(c.Request().Context(), boil.GetContextDB())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "fetch failed"})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := models.FindUser(c.Request().Context(), boil.GetContextDB(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := models.FindUser(c.Request().Context(), boil.GetContextDB(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	_, err = user.Update(c.Request().Context(), boil.GetContextDB(), boil.Infer())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := models.FindUser(c.Request().Context(), boil.GetContextDB(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}
	_, err = user.Delete(c.Request().Context(), boil.GetContextDB())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "delete failed"})
	}
	return c.NoContent(http.StatusNoContent)
}
