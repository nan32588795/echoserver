package validator

import (
	"net/http"
	"regexp"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type customValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *customValidator {
	v := validator.New()
	v.RegisterValidation("password", passwordValidator)
	return &customValidator{
		validator: v,
	}
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[!@#\$%\^&\*(),.?":{}|<>]`).MatchString(password)

	return hasLetter && hasNumber && hasSymbol
}
