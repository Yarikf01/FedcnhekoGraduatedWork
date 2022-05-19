package log

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

var IsSuppressErrors bool

type ValidationError struct {
	Field string
}

type ForbiddenError struct {
}

type ReconError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error, %s", e.Field)
}

func (e ForbiddenError) Error() string {
	return "access forbidden error"
}


func LogicErr(err error) *echo.HTTPError {
	switch {
	case errors.As(err, &ValidationError{}):
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	case errors.As(err, &ForbiddenError{}):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	case errors.As(err, &ReconError{}):
		return echo.NewHTTPError(http.StatusConflict, err)
	}

	if IsSuppressErrors {
		err = errors.New("something went wrong")
	}
	return echo.NewHTTPError(http.StatusConflict, err.Error())
}
