package types

type Question struct {
	Id          string          `json:"template_id"`
	Title       string          `json:"title"`
	IsPaid      bool            `json:"paid"`
    QuestionId  int64             `json:"qId"` 
	Data        []byte `json:"data"`
	Setting     []byte `json:"setting"`
}

type Content map[int]Question

type FormContents struct {
	FormId       string           `json:"form_id"`
	AccountId    string           `json:"account_id"`
	Questions    interface{}      `json:"questions"`
	FormSettings FormSetting_Type `json:"form_settings"`
	Answers      Answers          `json:"answers"`
}
