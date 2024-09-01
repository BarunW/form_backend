package types

type QuestionDataModels struct{}

func NewQuestionDataModels() QuestionDataModels{
    return QuestionDataModels{}
}

type QuestionSettingModels struct{}

func NewQuestionSettingModels() QuestionSettingModels{
    return QuestionSettingModels{}
}

type Account struct {
	ID       int    `json:"id" db:"user_id" validate:"-"`
	Username string `json:"username" db:"username" validate:"required,min=5,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

/*
================
Pricing Model
================
*/

type PlanType string
type BillingMode string
type CurrencyCode string
type ResponseLimitNum int
type PlanId_Type int

// billing period either monthly or yearly
const (
	Month BillingMode = "Monthly"
	Year  BillingMode = "Yearly"
)

// currency code
const (
	USD CurrencyCode = "USD"
	INR CurrencyCode = "INR"
)

// avaible plan
const (
	Free       PlanType = "free"
	Basic      PlanType = "basic"
	Plus       PlanType = "plus"
	Business   PlanType = "business"
	Enterprise PlanType = "enterprise"
)

// plan features
type Features struct {
	Discount float32 `json:"discount" db:"discount"`
	// basic feature
	Seats             int              `json:"seats" db:"seats"`
	MaxResponse       ResponseLimitNum `json:"max_response" db:"max_response"`
	UnlimitedForm     bool             `json:"unlimited" db:"unlimited"`
	AcceptPayment     bool             `json:"acceptpayment" db:"acceptpayment"`
	RecieveFileUpload bool             `json:"recieve_file" db:"file_upload"`
	RemovingBranding  bool             `json:"remove_brandig" db:"remove_branding"`
	CreateOwnBranding bool             `json:"own_branding" db:"own_branding"`

	// integration feature
	IntegrationFeatures struct {
	} `json:"intergration_feature"`
}

// details
type Pricing struct {
	PlanId        int          `json:"plan_id"`
	Plan          PlanType     `json:"plan"`
	Price         float32      `json:"price"`
	Currency      CurrencyCode `json:"currency"`
	BillingPeriod BillingMode  `json:"billing_period"`
	Features
}

// price calculation variables
type PriceCalculationVariables struct {
	PlanId          int             `json:"plan_id" validate:"required"`
	ResponseLimitID ResponseLimitId `json:"response_limit_id" validate:"required"`
}

/*
=================
Plan subscriptions
*rpv - response plan validator
*bpv - billing period validator
=================
*/
// Response Plan for business
type ResponseLimitId string

// Response Plan for basic

type Subscriber struct {
	UserId          int             `json:"user_id" validate:"required,gte=0,lte=100000000"`
	PlanId          int             `json:"plan_id" validate:"required,gte=1120,lte=1124"`
	ResponseLimitId ResponseLimitId `json:"response_limit_id" validate:"required"`
	BillingPeriod   BillingMode     `json:"billing_period" validate:"required,bpv"`
}

type PLAN_STATUS string

const (
	Active   PLAN_STATUS = "ACTIVE"
	Pause    PLAN_STATUS = "PAUSE"
	FreePlan PLAN_STATUS = "ON FREE"
)

/*
==================
Templates
=================
*/
type IsPaid bool

/*
===================

	Form

==================
*/
type CreateForm struct {
	FormName    string `json:"form_name" validate:"required,min=1,max=100"`
	UserId      int    `json:"user_id" validate:"required,gte=0,lte=99999999"`
	TemplateId  string `json:"template_id" validate:"required,min=26,max=26"`
	WorkspaceId string `json:"workspace_id" validate:"required,min=26,max=26"`
}

type AddQuestionTempl struct {
	Id         string `json:"id" validate:"required,min=26,max=26"`
	TemplateId string `json:"template_id" validate:"required,max=26,min=26"`
}

type FormMetaData struct {
	Id     string `json:"id" validate:"required,min=26,max=26"`
	UserId int    `json:"user_id" validate:"required,lte=99999999,gte=0"`
}

type ChangeQuestionPostition struct {
	Id             string `json:"id" validate:"required,min=26,max=26"`
	FromPositionId int    `json:"from" validate:"lte=999999"`
	ToPositionId   int    `json:"to" validate:"lte=999999"`
}

/*
===================

	Workspace

==================
*/
type Workspace struct {
	WorkspaceId string `json:"id"`
	Name        string `json:"name" validate:"required,min=0,max=69"`
	UserId      int    `json:"user_id" validate:"required"`
}

/*
=================
    Image
================
*/

type MediaType string

const (
	IMG  MediaType = "IMG"
	VID  MediaType = "VID"
	NONE MediaType = "null"
)

type UpdateImageorVideoURL struct {
	FormId     string    `json:"form_id" validate:"required,min=26,max=26"`
	QuestionId int       `json:"question_id"`
	URL        string    `json:"url"`
	MediaType  MediaType `json:"type"`
}
