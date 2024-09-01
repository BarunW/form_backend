package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/sonal3323/form-poc/builder"
	"github.com/sonal3323/form-poc/data"
	"github.com/sonal3323/form-poc/imageController"
	"github.com/sonal3323/form-poc/types"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *data.PostgresStore
}

func NewHandler(psql *data.PostgresStore) *Handler {
	return &Handler{
		db: psql,
	}
}

/*
==============
Validator
=============
*/

func IsValidULID(id string) bool {
	var validULID = regexp.MustCompile(`[0-7][0-9A-HJKMNP-TV-Z]{25}`)
	return validULID.MatchString(id)
}

func isValidFileName(fileName string) bool {
	return true
}

/*
=====================
auth handler
*Creating user
*Login
=====================
*/

type userData struct {
    user_name string
    user_id int
    email string
    access_token string
    refresh_token string
}
func (h *Handler) CreateUser(c echo.Context) error {
	newUsr := c.Get("data").(*types.Account)
	// check the email and usrename is already exist

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUsr.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("bcrypt hasing password failed", "details", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	newUsr.Password = string(hashedPassword)

    ud, cerr := h.db.CreateAccount(newUsr)
	if cerr != nil {
		if cerr.Error() == "409" {
			return echo.NewHTTPError(http.StatusConflict, "failed to create account email already exists for a user")
		}
		slog.Error("Account creation Failed", "details", cerr.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to create account right now")
	}
	// sign jwt token at this point
	return h.handleUserData(c, userData{
        user_id: ud.ID,
        email: ud.Email,
        user_name: ud.Username,
        access_token: ud.Tokens.AccessToken,
        refresh_token: ud.Tokens.RefreshToken,
    })

}

func (h *Handler) Login(c echo.Context) error {
	credentials := c.Get("data").(*types.Login)
	acc, err := h.db.GetAccount(credentials)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return echo.NewHTTPError(http.StatusForbidden, "wrong email or password")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
    
	return h.handleUserData(c, userData{
        user_id: acc.ID,
        email: acc.Email,
        user_name: acc.Username,
        access_token: acc.Tokens.AccessToken,
        refresh_token: acc.Tokens.RefreshToken,
    })
}

func (h *Handler) handleUserData( c echo.Context, u userData) error { 

	upi, uerr := h.db.UserPlanData(u.user_id)
	if uerr != nil {
		slog.Error("unable to get User Plan Data", "details", uerr.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to login")
	}

	if upi.UserId == 0 || upi.SubsId == 0 || upi.Username == "" {
		upi.UserId = u.user_id
		upi.SubsId = -69
		upi.Username = u.user_name 
		upi.Email = u.email
	}

    c.SetCookie(&http.Cookie{
        Name: "form_poc_token_",
        Value: u.refresh_token,
        HttpOnly: true,
        Secure: false,
        MaxAge: int(time.Now().AddDate(0,1,0).Unix()),
    })

    dataWithAccessToken := struct{
        Data interface{}   `json:"data"`
        AccessToken string `json:"access_token"`
    }{
        Data: upi,
        AccessToken: u.access_token,
    }

	// generate jwt or validate jwt
	return c.JSON(http.StatusCreated, dataWithAccessToken)
}

/*
=====================
Pricing handler
*Getting all pricing/plans
=====================
*/

func (h *Handler) PricingHandler(c echo.Context) error {
	resp, err := h.db.GetAllPricing()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get pricing data")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) PriceCalculation(c echo.Context) error {
	variables := c.Get("data").(*types.PriceCalculationVariables)

	if err := mapRespLimitIdWithPlanId(h, variables.PlanId, variables.ResponseLimitID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := h.db.PriceCalculation(variables)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch the pricing details")
	}
	return c.JSON(http.StatusOK, result)
}

/*
======================
ResponseLimiIds Mapping with price plan
=====================
*/
func mapRespLimitIdWithPlanId(h *Handler, planId int, rlId types.ResponseLimitId) error {
	ids := h.db.GetCachedRespLimitId()
	/*
	   --------------
	   1122 for basic monthly plan
	   1123 for plus monthly plan
	   1124 for business monthly plan

	   2122 for basic yearly plan
	   2123 for plus yearly plan
	   2124 for business yearly plan
	   ---------------
	*/
	switch planId {
	case 1122, 2122:
		_, ok := ids[0][rlId]
		if !ok {
			return fmt.Errorf("planId %d, responseId %v doesn't match ", planId, rlId)
		}
	case 1123, 2123:
		_, ok := ids[1][rlId]
		if !ok {
			return fmt.Errorf("planId %d, responseId %v doesn't match ", planId, rlId)
		}

	case 1124, 2124:
		_, ok := ids[2][rlId]
		if !ok {
			return fmt.Errorf("planId %d, responseId %v doesn't match ", planId, rlId)
		}
	default:
		return fmt.Errorf("please check the planid %d", planId)
	}

	return nil

}

/*
=====================
Subscription handler
*Subscribe to plan
* Get the subscribe details about a user
=====================
*/
func (h *Handler) SubscribePlan(c echo.Context) error {
	subData := c.Get("data").(*types.Subscriber)

	err := mapRespLimitIdWithPlanId(h, subData.PlanId, subData.ResponseLimitId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.db.Subscribe(*subData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to subscribed")
	}
	return c.JSON(http.StatusOK, "sucessfully subscribe")
}

/*
=================
User Response limits to form
=================
*/

func (h *Handler) ResponseLimit(c echo.Context) error {
	s, err := h.db.GetAllResponseLimit()

	if err != nil {
		slog.Error("Failed to response response limit data", "d", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "snap back to reality, here comes error")
	}
	return c.JSON(http.StatusOK, s)
}

/*
==============
Form And Template handler
=============
*/

func (h *Handler) HandleCreateForm(c echo.Context) error {
	cf := c.Get("data").(*types.CreateForm)
	err := h.db.CreateForm(*cf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "snap back to reality, here comes error")
	}
	return c.JSON(http.StatusOK, "Created")
}

// To add more question to the form
func (h *Handler) AddQuestionToForm(c echo.Context) error {
	addForm := c.Get("data").(*types.AddQuestionTempl)
	err := h.db.AddQuestion(*addForm)
	if err != nil && err.Error() != "404" {
		return echo.NewHTTPError(http.StatusInternalServerError, "snap back to reality, here comes error")
	}
	return c.JSON(http.StatusOK, "Updated")
}

// Upadte the question for a question in a form
func (h *Handler) HandleUpdateQuestion(c echo.Context) error {
	data := c.Get("data").(*types.UpdateQuestionOrDescription)
	if err := h.db.UpdateQuestion(data.QuestionId, data.Text, data.FormId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Snap back to reality here comes error")
	}
	return c.JSON(http.StatusOK, "Successfully Updated")
}

func (h *Handler) HandleUpdateDescription(c echo.Context) error {
	data := c.Get("data").(*types.UpdateQuestionOrDescription)
	if err := h.db.UpdateDescription(data.QuestionId, data.Text, data.FormId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Snap back to reality here comes error")
	}

	return c.JSON(http.StatusOK, "Successfully Updated")

}

func (h *Handler) HandleGetForm(c echo.Context) error {
	formMetaData := c.Get("data").(*types.FormMetaData)
	form, err := h.db.GetForm(*formMetaData)
	if err != nil {
        if err == sql.ErrNoRows{
            return echo.NewHTTPError(http.StatusNotFound, "Form doesn't exist")
        }
		return echo.NewHTTPError(http.StatusBadRequest, "snap back to reality, here comes error")
	}
	return c.JSON(http.StatusOK, form)
}

func (h *Handler) HandleDeleteForm(c echo.Context) error {
    userId, err := strconv.Atoi(c.Param("userId"))
    if  err != nil{
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
    }

	formId := c.Param("id")
	if !IsValidULID(formId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}
	if err := h.db.DeleteForm(formId, userId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "snap back to reality here comes error")
	}
	return c.JSON(http.StatusOK, "Successfully Deleted")
}

func (h *Handler) HandleDeleteQuestion(c echo.Context) error {
	var formId string = c.Param("formId")
	questionId, err := strconv.Atoi(c.Param("questionId"))

	if err != nil || !IsValidULID(formId) {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "snap back to reality, you provided invalid id")
	}

	err = h.db.DeleteQuestion(formId, questionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "snap back to reality, you provided invalid id")
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (h *Handler) HandleGetFormLink(c echo.Context) error {

	formId := c.Param("formId")
	if !IsValidULID(formId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	var link string
	link, err := h.db.GetFormLink(formId)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "404")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "snap back to reality here comes error")
	}

	return c.JSON(http.StatusOK, link)
}

/*
=================
Workspace
=================
*/

func (h *Handler) HandleCreateWorkspace(c echo.Context) error {
	workspace := c.Get("data").(*types.Workspace)
	if err := h.db.CreateWorkspace(*workspace); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid User Id")
	}
	return c.JSON(http.StatusCreated, "Workspace created")
}

func (h *Handler) HandleGetWorkspace(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userId == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid User Id")
	}

	workspaces, werr := h.db.GetWorkspace(userId)
	if werr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something Went Wrong")
	}
	return c.JSON(http.StatusOK, workspaces)
}

func (h *Handler) HandleGetWorkspaceForms(c echo.Context) error {

	workspaceId := c.Param("workspace_id")

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || !IsValidULID(workspaceId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	formDatas, err := h.db.GetWorkSpaceFormsData(userId, workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(http.StatusOK, formDatas)
}

func (h *Handler) HandleDeleteWorkspace(c echo.Context) error {
    userId, err := strconv.Atoi(c.Param("userId"))
    if  err != nil{
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
    }
	workspaceId := c.Param("id")
	if !IsValidULID(workspaceId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	err = h.db.DeleteWorkspace(userId, workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong")
	}
	return c.JSON(http.StatusOK, "Successfully Deleted")
}

/*
=================
ImageAndVideo settings Handler
================
*/
type ImageVideoRouteParams struct {
	UserId     int
	FormId     string
	QuestionId int
	FileName   string
}

func validateParams(c echo.Context) (*ImageVideoRouteParams, error) {

	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return nil, fmt.Errorf("Error")
	}

	form_id := c.Param("form_id")
	if !IsValidULID(form_id) {
		fmt.Println("Invalid form_id", form_id)
		return nil, fmt.Errorf("Error")
	}

	questionId, err := strconv.Atoi(c.Param("question_id"))
	if err != nil {
		fmt.Println("Invalid questionId", questionId)
		return nil, fmt.Errorf("Error")
	}

	fileName := c.Param("filename")
	if !isValidFileName(fileName) {
		fmt.Println("Invalid fileName", fileName)
		return nil, fmt.Errorf("Error")
	}

	return &ImageVideoRouteParams{
		UserId:     user_id,
		FormId:     form_id,
		QuestionId: questionId,
		FileName:   fileName,
	}, nil

}

func (h *Handler) HandleImageUpload(c echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to the get file")
	}

	file, ferr := fileHeader.Open()
	if ferr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to the open file")
	}

	params, perr := validateParams(c)
	if perr != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Error Inavlid Id")
	}

	err = h.db.UploadImage(params.UserId, params.FormId, params.QuestionId, params.FileName, file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error Inavlid Id")
	}

	return c.HTML(http.StatusOK, "<h1>Test successfull </h1>")

}

func (h *Handler) HandleGetImage(c echo.Context) error {
	mfc := imageController.NewMediaFilesController()
	userId := c.Param("user_id")
	fileName := c.Param("filename")

	basePath, err := mfc.GetBasePath()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Whose gonna the error")
	}
	fmt.Println(basePath)

	fullPath := filepath.Join(basePath, userId, fileName)
	return c.File(fullPath)
}

func (h *Handler) HandleUpdateImageOrVideoURL(c echo.Context) error {
	data := c.Get("data").(*types.UpdateImageorVideoURL)
	fmt.Println(data)
	if !IsValidULID(data.FormId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}
	err := h.db.UpdateImageOrVideoURL(data.FormId, data.QuestionId, data.URL, data.MediaType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to add the url :(")
	}
	return c.JSON(http.StatusOK, "successfully added the url")
}

//func (h *Handler) HandleDeleteImageFile(c echo.Context) error {
//	formId := c.Param("imageId")
//	if !isValidULID(formId) {
//		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
//	}
//
//	questionId, err := strconv.Atoi(c.Param("userId"))
//	if err != nil {
//		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
//	}
//
//	err = h.db.DeleteImageFile(userId, imageId)
//	if err != nil {
//		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to delete the file")
//	}
//
//	return c.JSON(http.StatusOK, "Successfully Removed")
//}

/*
   ==================
    GET STATIC/ PUBLISHED FORM
   =================
*/

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func (h *Handler) HandleGetPublishedForm(c echo.Context) error {
	accountId := c.Get("account_id").(string)

	formId := c.Param("formId")
	if !IsValidULID(formId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}
	fd := types.FormMetaData{
		Id: formId,
	}

	// check the account exist for that particular formID`
	if accountId_InDb, err := h.db.GetAccountIdThroughtFormId(formId); err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "Snap back to reality here comes error")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Snap back to reality internal server error")
	} else if accountId_InDb != accountId {
		return echo.NewHTTPError(http.StatusBadRequest, "Snap back to reality here comes error")
	}

	fb := builder.NewFormBuilder(*h.db)
	comp, err := fb.BuildForm(fd)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Snap back to reality here it comes 404")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Sorry Unable to get the form")
	}

	return Render(c, http.StatusOK, comp)
}

// form submition
func (h *Handler) HandleSubmitForm(c echo.Context) error {

	accountId := c.Get("account_id").(string)
	formId := c.Param("formId")
	if !IsValidULID(formId) {
		return echo.NewHTTPError(http.StatusBadGateway, "Invalid url")
	}

	// extract the answers from the request body
	bdy := c.Request().Body
	bdyByt, err := io.ReadAll(bdy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to get request body")
	}

	err = h.db.Submit(formId, accountId, bdyByt)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "404")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, "successfully submitted")
}
// Result
func (s *Handler) HandleGetResult(c echo.Context)  error{
    formId := c.Param("formId")
    if !IsValidULID(formId){
       return echo.NewHTTPError(http.StatusForbidden, "Invalid" )
    }
    
    result, err := s.db.GetResult(formId)
    if err != nil{
        if err == sql.ErrNoRows{
            return echo.NewHTTPError(http.StatusNotFound, "Form Not found" )
        }
        return echo.NewHTTPError(http.StatusInternalServerError, "Unknown error")
    }

    return  c.JSON(http.StatusOK, result) 
}



/*
=================
Responses
@return ResponsesModel and error
=================
*/

func (s *Handler) HandleGetResponses(c echo.Context) error {

	userId, uerr := strconv.Atoi(c.Param("userId"))
	if uerr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid userId")
	}

	r, err := s.db.GetResponses(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "oops internal server error")
	}

	return c.JSON(http.StatusOK, r)
}


