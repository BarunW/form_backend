package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sonal3323/form-poc/types"
)

/*
============
custom Validator implementation
=============
*/

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()
	cv := &CustomValidator{
		validator: validate,
	}
	validate.RegisterValidation("bpv", cv.billingPeriodValidator)
	return cv
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

/*
==============
validation for billing period
===============
*/
func (cv *CustomValidator) billingPeriodValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == string(types.Month) || value == string(types.Year) {
		return true
	}
	return false
}
