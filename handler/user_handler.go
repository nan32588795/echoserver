package handler

import (
	"net/http"
	"strconv"

	"echoserver/models"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type createUserRequest struct {
	Username string `json:"username" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password"`
}

type createUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to hash password"})
	}

	u := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := u.Insert(c.Request().Context(), boil.GetContextDB(), boil.Infer()); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "insert failed"})
	}

	var res createUserResponse = createUserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
	return c.JSON(http.StatusCreated, res)
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
