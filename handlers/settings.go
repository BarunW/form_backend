package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sonal3323/form-poc/types"
)

/*
===========

	Question Settings

============
*/
func (h *Handler) HandleSettingsRequired(c echo.Context) error {
	data := c.Get("data").(*types.UpdateBoolSetting)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	err := h.db.UpdateRequiredSetting(data.FormId, data.QuestionId, data.Setting)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown erro")
	}

	return c.JSON(http.StatusOK, "Updated")
}

/*
   ============
   FormSettings
   ===========
*/

func (h *Handler) HandleUpdateFormSettingProgresBar(c echo.Context) error {
	data := c.Get("data").(*types.UpdateFormBoolSetting)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	err := h.db.UpdateFormSettingProgresBar(data.FormId, data.Setting)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown erro")
	}

	return c.JSON(http.StatusOK, "Updated")
}

func (h *Handler) HandleUpdateFormSettingQNO(c echo.Context) error {
	data := c.Get("data").(*types.UpdateFormBoolSetting)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	err := h.db.UpdateFormSettingQNO(data.FormId, data.Setting)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown erro")
	}

	return c.JSON(http.StatusOK, "Updated")
}

func (h *Handler) HandleFormSettLettersOnAns(c echo.Context) error {
	data := c.Get("data").(*types.UpdateFormBoolSetting)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	err := h.db.UpdateFormSettingLettersOnAns(data.FormId, data.Setting)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown erro")
	}

	return c.JSON(http.StatusOK, "Updated")
}

func (h *Handler) HandleFormSettingFreeNav(c echo.Context) error {
	data := c.Get("data").(*types.UpdateFormBoolSetting)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	err := h.db.UpdateFormSettingFreeNav(data.FormId, data.Setting)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown erro")
	}

	return c.JSON(http.StatusOK, "Updated")
}

func (h *Handler) HandleFormSettingNavArrows(c echo.Context) error {
	data := c.Get("data").(*types.UpdateFormBoolSetting)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	err := h.db.UpdateSettingFormNavArrows(data.FormId, data.Setting)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown erro")
	}

	return c.JSON(http.StatusOK, "Updated")
}
