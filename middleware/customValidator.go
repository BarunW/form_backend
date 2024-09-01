package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ValidatorMiddleware struct{}

func NewValidatorMiddleware() *ValidatorMiddleware {
	return &ValidatorMiddleware{}
}

func bindAndValidate(c echo.Context, target interface{}) error {
	if err := c.Bind(target); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(target); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (vm *ValidatorMiddleware) JsonValidator(i interface{}, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := bindAndValidate(c, i); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		c.Set("data", i)
		return next(c)
	}
}

func (vm *ValidatorMiddleware) CheckSubdomain(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := stripeAndBindSubdomain(c); err != nil {
			return err
		}
		return next(c)
	}
}

// as account_id is subdomain
// we stripe the account
func stripeAndBindSubdomain(c echo.Context) error {
	tokens := strings.Split(c.Request().Host, ".")
	if len(tokens) < 1 {
		return echo.NewHTTPError(http.StatusForbidden, "wrong url")
	}

	accountId := tokens[0]

	if len(accountId) != 11 {
		return echo.NewHTTPError(http.StatusForbidden, "wrong url")
	}

	c.Set("account_id", accountId)

	return nil
}
