package types

/*
   ============
   Multiple Choice question Update Data Types
   =============
*/

type MCQAddChoice struct {
	FormId     string `json:"form_id" validate:"max=26,min=26"`
	Text       string `json:"text"`
	QuestionId int    `json:"question_id"`
}

// Types for upadting choice
type MCQAnswerUpdate struct {
	FormId     string `json:"form_id" validate:"max=26,min=26"`
	Text       string `json:"text"`
	QuestionId int    `json:"question_id"`
	ChoiceId   string `json:"choice_id"`
}

/*
===============
MCQ Data and settings
===============
*/
type Choice struct {
	Label string `json:"label"`
}

type MutipleChoiceData struct {
	QuestionText string            `json:"question"`
	Description  string            `json:"description"`
	Choices      map[string]Choice `json:"choices"`
	TotalChoice  int               `json:"total_choice"`
}

// MCQ settings
type MultipleSelection_Type struct {
	Type         []string `json:"type"`
	SelectedType string   `json:"selected_type"`
	Number       int      `json:"number"`
}

type Settings_MultipleChoice struct {
	Required             bool                      `json:"required"`
	MultipleSelection    MultipleSelection_Type    `json:"multiple_selection"`
	OtherOption          bool                      `json:"other_option"`
	VerticalAlignment    bool                      `json:"vertial_aligment"`
	ImageOrVideoSettings ImageOrVideoSettings_Type `json:"image_or_videoSettings"`
}

// MCQ structure
type MCQ_QuestionModel struct {
    questionMetaData 
	Data    MutipleChoiceData       `json:"data"`
	Setting Settings_MultipleChoice `json:"setting"`
}

func(q QuestionDataModels) GetMCQDataModel() MutipleChoiceData {
   return MutipleChoiceData{
			QuestionText: "...",
			Description:  "",
			Choices:      map[string]Choice{"01HQ5PJ2XBK9SBE2D58BAV570S": {Label: "Choice One"}},
			TotalChoice:  1,
		}
}


func(q QuestionSettingModels) GetMCQSettingModel() Settings_MultipleChoice{
    return Settings_MultipleChoice{
			Required: false,
			MultipleSelection: MultipleSelection_Type{
				Type:         []string{"Ulimited", "Custom Number", "Exact Number"},
				Number:       1,
				SelectedType: "Ulimited",
			},
			OtherOption:       false,
			VerticalAlignment: true,
			ImageOrVideoSettings: ImageOrVideoSettings_Type{
				Url: "null",
				Layout: ImageLayout{
					Mobile:        []int{1, 2, 3, 4},
					Desktop:       []int{1, 2, 3, 4, 5, 6},
					MobileLayout:  1,
					DesktopLayout: 1,
				},
			},
		}

}
