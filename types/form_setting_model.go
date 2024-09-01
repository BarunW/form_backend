package types

type AccessSchedulingSetting struct {
	CloseForNewUser     bool
	ScheduleCloseDate   bool
	ResponseLimit       bool
	ShowCustomClosedMsg bool
}

type FormMessageSettings struct {
	AnswerButton                string `json:"answer_button" validate:"max=50,min=1"`
	KeyForNextQuestion          string `json:"key_fornext_question" validate:"max=200,min=1"`
	HintForMultipleSelection    string `json:"hint_for_multiple_selection" validate:"max=165,min=0"`
	DropDownInstructionQuestion string `json:"drop_down_instruction_ques" validate:"max=100, min=0"`
	OtherOptionLabel            string `json:"other_option_label" validate:"max=100, min=0"`
}

type FormSetting_Type struct {
	AccessAndScheduling AccessSchedulingSetting `json:"access_scheduling"`
	Messages            FormMessageSettings     `json:"messages"`
	CookieConstent      bool                    `json:"cookie_constent"`
	Progressbar         bool                    `json:"progressbar"`
	QuestionNumber      bool                    `json:"question_number"`
	LettersOnAnswers    bool                    `json:"letters_on_answers"`
	Branding            bool                    `json:"branding"`
	NavigationArrows    bool                    `json:"navigation_arrows"`
	FreeFormNavigation  bool                    `json:"free_form_navigation"`
	AutosaveProgress    bool                    `json:"auto_save_progress"`
	UTMTracking         bool                    `json:"utm_tracking"`
}

type UpdateFormBoolSetting struct {
	FormId  string `json:"form_id" validate:"max=26,min=26"`
	Setting bool   `json:"setting"`
}
