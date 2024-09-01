package types

type UpdateQuestionOrDescription struct {
	FormId     string `json:"form_id" validate:"max=26,min=26"`
	Text       string `json:"text"`
	QuestionId int    `json:"question_id"`
}

type UpdateBoolSetting struct {
	FormId     string `json:"form_id" validate:"max=26,min=26"`
	Setting    bool   `json:"setting"`
	QuestionId int    `json:"question_id"`
}

type questionMetaData struct{
    QuestionUUID int64  `json:"quuid"`
    TemplateId   string `json:"template_id"`
    Title        string `json:"title"` 
	Paid         bool   `json:"paid"`
}


type RequiredModel struct {
	Label   string `json:"label"`
	Require bool   `json:"require"`
}


func( q *QuestionSettingModels) getImageOrVideoSettings() ImageOrVideoSettings_Type{ 
    return ImageOrVideoSettings_Type{
        Url: "null",
        Layout: ImageLayout{
            Mobile:        []int{1, 2, 3, 4},
            Desktop:       []int{1, 2, 3, 4, 5, 6},
            MobileLayout:  1,
            DesktopLayout: 1,
        },
    }
}

