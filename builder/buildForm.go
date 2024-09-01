package builder

import (
	"encoding/json"

	"github.com/a-h/templ"
	"github.com/sonal3323/form-poc/data"
	"github.com/sonal3323/form-poc/integration-react/static"
	"github.com/sonal3323/form-poc/types"
)

type FormBuilder struct {
	ps data.PostgresStore
}

func NewFormBuilder(ps data.PostgresStore) FormBuilder {
	return FormBuilder{
		ps: ps,
	}
}

type anyInterface struct{
   Data any
   Setting any
}

func (f FormBuilder) BuildForm(data types.FormMetaData) (templ.Component, error) {
	formJSON, err := f.ps.GetForm(data)
	if err != nil {
		return nil, err
		//        echo.NewHTTPError(http.StatusInternalServerError, "Unknown Error")
	}
	var contents = []types.Content{}
	var formSettings types.FormSetting_Type

	if err := json.Unmarshal(formJSON.Content, &contents); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(formJSON.Setting, &formSettings); err != nil {
		return nil, err
	}

	var questions []interface{}
	answers := make(types.Answers)
    var fn func(quest types.Question, i *anyInterface ) error

    fn = func(quest types.Question, i *anyInterface ) error {
        if err := json.Unmarshal(quest.Data, i.Data); err != nil {
            return  err
        }
        if err := json.Unmarshal(quest.Setting, i.Setting); err != nil {
            return  err
        }
        return nil
    }
    
	for idx, content := range contents {
        // check the questionTemplate id
		switch content[idx].Id {
		case string(types.MultipleChoice_ID):
			var mcq_s types.MCQ_QuestionModel
            if err := fn(content[idx], &anyInterface{Data: &mcq_s.Data, Setting: &mcq_s.Setting}); err != nil{
                return nil, err
            }
            qId := content[idx].QuestionId
            mcq_s.Title = content[idx].Title
            mcq_s.TemplateId = content[idx].Id
            mcq_s.QuestionUUID = qId

			questions = append(questions, mcq_s)
			answers[int(qId)] = types.ResponseAnswers{} 

		case string(types.ContactInfo_ID):
			var contactinfo types.ContactInfoQuestionModel
            if err := fn(content[idx], &anyInterface{Data: &contactinfo.Data, Setting :&contactinfo.Setting}); err != nil {
				return nil, err
			}
            qId := content[idx].QuestionId
            contactinfo.Title = content[idx].Title
            contactinfo.TemplateId = content[idx].Id
            contactinfo.QuestionUUID = qId
            
			questions = append(questions, contactinfo)
			answers[int(qId)] = types.ResponseAnswers{}
    
		}
	}

	account_id, err := f.ps.GetAccountIdThroughtFormId(formJSON.Id)

	fc := types.FormContents{
		FormId:       formJSON.Id,
		AccountId:    account_id,
		Questions:    questions,
		FormSettings: formSettings,
		Answers:      answers,
	}

	formComponent := static.Form(fc)

	return formComponent, nil
}
