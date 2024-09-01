package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sonal3323/form-poc/data"
	"github.com/sonal3323/form-poc/handlers"
	"github.com/sonal3323/form-poc/middleware"
	"github.com/sonal3323/form-poc/types"
	"github.com/sonal3323/form-poc/utils"
)

/*
==============================
Setting up postgres database connection
=============================
*/

func dbConn() (*data.PostgresStore, error) {
	psql, err := data.NewPostgresStore()
	if err != nil {
		slog.Error("Unable to set up database connection", "details", err.Error())
		return nil, err
	}

	// this will create all necessary table in db
	err = psql.Init()
	if err != nil {
		return nil, err
	}

	return psql, err

}

/*
===========================
Setting up Server
==========================
*/

func initServer(e *echo.Echo) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      e,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Server Unable to listen", "details", err.Error())
			os.Exit(1)
		}
	}()

	slog.Info("Server running", "PORT", "8080")

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	<-sigChan
	slog.Info("Server", "gracefully shutting down", "")
	tc, f := context.WithTimeout(context.Background(), 30*time.Second)
	f()
	server.Shutdown(tc)
}

// Handle all routes related to question
func handleQuestionUpdatesRoutes(e *echo.Echo, h *handlers.HandleQuestionUpdates, cvm *middleware.ValidatorMiddleware) {
	/*
	   ==============
	   Multiple Choice
	   ==============
	*/
	e.PUT("/mcq/addChoice", cvm.JsonValidator(
		&types.MCQAddChoice{}, echo.HandlerFunc(h.HandleMCQAddChoice)))

	e.PUT("/mcq/updateChoice", cvm.JsonValidator(
		&types.MCQAnswerUpdate{}, echo.HandlerFunc(h.HandleUpdateChoiceAnswer)))

	e.DELETE("/mcq/removeChoice/:form_id/:question_id/:choice_id", echo.HandlerFunc(h.HandleMCQDeleteChoice))

	//Settings
	e.PUT("/mcq/setting/otherOption", cvm.JsonValidator(
		&types.UpdateBoolSetting{}, echo.HandlerFunc(h.HandleMCQOtherOptionSetting)))
	/*
	   ==============
	   Contact Info
	   =============
	*/

}

/*
======================
Main func Program Entry Point
======================
*/

func main() {
	// db connection
	psql, err := dbConn()
	if err != nil {
		os.Exit(1)
	}

	//handlers
	h := handlers.NewHandler(psql)

	// echo router
	e := echo.New()

	//set echo validator to this customValidator
	e.Validator = utils.NewCustomValidator()

	// all validation middleware will be use from this struct
	cvm := middleware.NewValidatorMiddleware()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	// middleware for serving static file
	e.Static("/static", "./integration-react/static")
	e.Use(echoMiddleware.Gzip())
    e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:3000"}, // Adjust the origin to match your web application's origin
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Access-Control-Allow-Origin"},
		ExposeHeaders: []string{"Access-Control-Allow-Origin"},	}))
	/*
	   -----------------
	   Handle GET Routes
	   -----------------
	*/
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/pricing", echo.HandlerFunc(h.PricingHandler))
	e.GET("/response-limit", echo.HandlerFunc(h.ResponseLimit))
	e.GET("/workspace/:user_id", echo.HandlerFunc(h.HandleGetWorkspace))
	e.GET("/workspace/forms/:user_id/:workspace_id", echo.HandlerFunc(h.HandleGetWorkspaceForms))
	e.GET("/questions", echo.HandlerFunc(h.HandleGetAllQuestionTemplates))
	e.GET("/responses/:userId", echo.HandlerFunc(h.HandleGetResponses))

	/*
	   -----------------
	   Handle POST Routes
	   -----------------
	*/
	e.POST("/pricing-details", cvm.JsonValidator(&types.PriceCalculationVariables{}, echo.HandlerFunc(h.PriceCalculation)))
	e.POST("/signup", cvm.JsonValidator(&types.Account{}, echo.HandlerFunc(h.CreateUser)))
	e.POST("/login", cvm.JsonValidator(&types.Login{}, echo.HandlerFunc(h.Login)))
	e.POST("/subscribe", cvm.JsonValidator(&types.Subscriber{}, echo.HandlerFunc(h.SubscribePlan)))
	e.POST("/create-form", cvm.JsonValidator(&types.CreateForm{}, echo.HandlerFunc(h.HandleCreateForm)))
	e.POST("/form/add-question", cvm.JsonValidator(&types.AddQuestionTempl{}, echo.HandlerFunc(h.AddQuestionToForm)))
	e.POST("/get-form", cvm.JsonValidator(&types.FormMetaData{}, echo.HandlerFunc(h.HandleGetForm)))
	e.POST("/workspace", cvm.JsonValidator(&types.Workspace{}, echo.HandlerFunc(h.HandleCreateWorkspace)))

	/*
	   -----------------
	   Handle UPDATE Routes
	   -----------------
	*/
	e.PATCH("/question/updatePosition", cvm.JsonValidator(&types.ChangeQuestionPostition{}, echo.HandlerFunc(h.ChangQuestionPosition)))

	e.PUT("/question/updateQuestion", cvm.JsonValidator(
		&types.UpdateQuestionOrDescription{}, echo.HandlerFunc(h.HandleUpdateQuestion)))

	e.PUT("/question/updateDescription", cvm.JsonValidator(
		&types.UpdateQuestionOrDescription{}, echo.HandlerFunc(h.HandleUpdateDescription)))

	e.PUT("/question/setting/required", cvm.JsonValidator(
		&types.UpdateBoolSetting{}, echo.HandlerFunc(h.HandleSettingsRequired)))

	/*
	   -----------------
	   Handle DELETE Routes
	   -----------------
	*/

	e.DELETE("/question/:formId/:questionId", echo.HandlerFunc(h.HandleDeleteQuestion))
    e.DELETE("/form/:userId/:id", echo.HandlerFunc(h.HandleDeleteForm))
    e.DELETE("/workspace/:userId/:id", echo.HandlerFunc(h.HandleDeleteWorkspace))

	/*
	   Form Setting Update
	   @Progressbar
	   @QuestionNumber
	   @LettersOnAnswers
	   @FreeNaviagtion
	   @Navigation Arrow
	*/

	e.PUT("/form/setting/pb", cvm.JsonValidator(
		&types.UpdateFormBoolSetting{}, echo.HandlerFunc(h.HandleUpdateFormSettingProgresBar)))

	e.PUT("form/setting/qno", cvm.JsonValidator(
		&types.UpdateFormBoolSetting{}, echo.HandlerFunc(h.HandleUpdateFormSettingQNO)))

	e.PUT("form/setting/letters", cvm.JsonValidator(
		&types.UpdateFormBoolSetting{}, echo.HandlerFunc(h.HandleFormSettLettersOnAns)))

	e.PUT("form/setting/freeNav", cvm.JsonValidator(
		&types.UpdateFormBoolSetting{}, echo.HandlerFunc(h.HandleFormSettingFreeNav)))

	e.PUT("form/setting/navArrow", cvm.JsonValidator(
		&types.UpdateFormBoolSetting{}, echo.HandlerFunc(h.HandleFormSettingNavArrows)))

	/*
	   --------------
	   Image
	   --------------
	*/

	e.GET("/image/:user_id/:filename", echo.HandlerFunc(h.HandleGetImage))
	e.POST("/image/upload/:user_id/:form_id/:question_id/:filename", echo.HandlerFunc(h.HandleImageUpload))
	e.PUT("/image/updateURL", cvm.JsonValidator(&types.UpdateImageorVideoURL{}, echo.HandlerFunc(h.HandleUpdateImageOrVideoURL)))
	//	e.DELETE("/image/deleteFile/:form_id/:file_id", echo.HandlerFunc(h.HandleDeleteImageFile))
	/*
	   ---------------
	   Serve the created form
	   ---------------
	*/
	e.GET("/to/:formId", echo.HandlerFunc(h.HandleGetPublishedForm), cvm.CheckSubdomain)
	e.GET("form-link/:formId", h.HandleGetFormLink)
	e.POST("/form/:formId/submissions", echo.HandlerFunc(h.HandleSubmitForm), cvm.CheckSubdomain)
	e.GET("/form-result/:formId", echo.HandlerFunc(h.HandleGetResult))

	/*
	   -------------------
	   start the server using *echo.Echo(type of struct)  as handler
	   -------------------
	*/

	// handler for updating questions templates fields
	quesUpdatesHandler := handlers.NewHandleQuestionUpdates(psql.QU)
	handleQuestionUpdatesRoutes(e, quesUpdatesHandler, cvm)

	// start the server and wait
	initServer(e)

}
