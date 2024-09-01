package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sonal3323/form-poc/data/questionsUpdates"
	"github.com/sonal3323/form-poc/types"
)

type HandleQuestionUpdates struct {
	qu *questionsUpdates.QuestionCRUD
}

func NewHandleQuestionUpdates(qu *questionsUpdates.QuestionCRUD) *HandleQuestionUpdates {
	return &HandleQuestionUpdates{
		qu: qu,
	}
}

func (h *Handler) HandleGetAllQuestionTemplates(c echo.Context) error {
	qt, err := h.db.GetQuestionTemplates()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error failed to get questions")
	}

	return c.JSON(http.StatusOK, qt)
}

func (h *Handler) ChangQuestionPosition(c echo.Context) error {
	data := c.Get("data").(*types.ChangeQuestionPostition)
	if err := h.db.ReorderQuestion(data.Id, data.FromPositionId, data.ToPositionId); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "snap back to reality, here comes error")
	}
	return c.JSON(http.StatusOK, "Success")
}

func stripeURL(c echo.Context) (string, error){
    url := c.Request().URL.Path
    paths := strings.Split(url,"/")
    if len(paths) < 1 {
        return "", fmt.Errorf("%d",http.StatusNotFound)
    }
    return paths[1], nil
}

/*
=============
Multiple Choice Question Updates Handlers
===========
*/
func (h *HandleQuestionUpdates) HandleMCQAddChoice(c echo.Context) error { 
    
    path, err := stripeURL(c)
    if err != nil{
		return echo.NewHTTPError(http.StatusNotFound, "Invalid id")
    }

	data := c.Get("data").(*types.MCQAddChoice)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid id")
	}

    err = h.qu.IsPathMatchWithQuestionId(data.QuestionId, data.FormId, path)
    if err != nil{
		return echo.NewHTTPError(http.StatusNotFound, "question id doesn't match with template id")
    }

	if err := h.qu.AddChoice(data.QuestionId, data.Text, data.FormId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Snap back to reality here comes error")
	}
	return c.JSON(http.StatusOK, "Successfully Updated")
}

func (h *HandleQuestionUpdates) HandleMCQDeleteChoice(c echo.Context) error {

	formId := c.Param("form_id")
	if !IsValidULID(formId) {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid id")
	}

	choiceId := c.Param("choice_id")

	questionId, err := strconv.Atoi(c.Param("question_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	if err := h.qu.DeleteChoice(questionId, choiceId, formId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops Error cause by soya dev")
	}

	return c.JSON(http.StatusOK, "Successfully Remove")
}

func (h *HandleQuestionUpdates) HandleUpdateChoiceAnswer(c echo.Context) error {
	data := c.Get("data").(*types.MCQAnswerUpdate)

	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid id")
	}

	if err := h.qu.UpdateChoiceAnswer(data.QuestionId, data.ChoiceId, data.Text, data.FormId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Snap back to reality here comes error")
	}
	return c.JSON(http.StatusOK, "Successfully Updated")
}

/*
=========
Settings Handler
=========
*/
func (h *HandleQuestionUpdates) HandleMCQOtherOptionSetting(c echo.Context) error {
	data := c.Get("data").(*types.UpdateBoolSetting)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	err := h.qu.UpdateOtherOption(data.FormId, data.QuestionId, data.Setting)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown erro")
	}

	return c.JSON(http.StatusOK, "Updated")
}
