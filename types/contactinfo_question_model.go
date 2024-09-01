package types

// data
type ContactInfoModel struct {
	Question    string           `json:"question"`
	Description string           `json:"description"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	PhoneNumber PhoneNumberModel `json:"phone_number"`
	Email       string           `json:"email"`
	Company     string           `json:"company"`
}

// Settings

type ExcludeIncludeReqModel struct {
	Label   string `json:"label"`
	Require bool   `json:"require"`
	Include bool   `json:"include"`
}

type ContactInfoSettingModel struct {
	QuestionSetting      map[int]ExcludeIncludeReqModel `json:"question_setting"`
	ImageOrVideoSettings ImageOrVideoSettings_Type      `json:"image_or_videoSettings"`
}

type ContactInfoQuestionModel struct {
    questionMetaData     
	Data    ContactInfoModel        `json:"data"`
	Setting ContactInfoSettingModel `json:"setting"`
    
}

func(qd *QuestionDataModels) GetContactInfoDataModel() ContactInfoModel{
    return ContactInfoModel{
			Question:    "...",
			Description: "",
			FirstName:   "Sush",
			LastName:    "M",
			PhoneNumber: PhoneNumberModel{
				Country:     "IN",
				CountryCode: "91",
				Number:      "0123456789",
			},
			Email:   "sush@mail.com",
			Company: "TheStartUp",
		}

}

func(qs *QuestionSettingModels) GetContactInfoSettingsModel() ContactInfoSettingModel{
    return  ContactInfoSettingModel{
			QuestionSetting: map[int]ExcludeIncludeReqModel{
				1: {
					Label:   "First Name",
					Include: true,
					Require: false,
				},
				2: {
					Label:   "Last Name",
					Include: true,
					Require: false,
				},
				3: {
					Label:   "Contact",
					Include: true,
					Require: false,
				},
				4: {
					Label:   "Email",
					Include: true,
					Require: false,
				},
				5: {
					Label:   "Company",
					Include: true,
					Require: false,
				},
			},
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

// Update
