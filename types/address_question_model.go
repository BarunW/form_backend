package types

// data
type AddressDataModel struct {
    Question                string           `json:"question"`
    Description             string           `json:"description"`
    Address                 string           `json:"address"`
    AddressLine2            string           `json:"address_line_2"`
    City                    string           `json:"city"`
    State                   string           `json:"state"`
    Zipcode                 string           `json:"zipcode"`
    Country                 string           `json:"country"`
}

// Settings

type AddressSettingModel struct {
	QuestionSetting      map[int]RequiredModel `json:"question_setting"`
	ImageOrVideoSettings ImageOrVideoSettings_Type      `json:"image_or_videoSettings"`
}


type AddresQuestionModel struct {
    questionMetaData     
	Data    ContactInfoModel        `json:"data"`
	Setting ContactInfoSettingModel `json:"setting"`
    
}

func(q QuestionSettingModels ) GetAddressQuestionSetting() AddressSettingModel{
    return AddressSettingModel{
			QuestionSetting: map[int]RequiredModel{
				1: {
					Label:   "Address",
					Require: false,
				},
				2: {
					Label:   "Address Line 2",
					Require: false,
				},
				3: {
					Label:   "City/Town",
					Require: false,
				},
				4: {
					Label:   "State/Region/Province",
					Require: false,
				},
				5: {
					Label:   "Zipcode/Postcode",
					Require: false,
				},
				6: {
					Label:   "Country",
					Require: false,
				},
			},
			ImageOrVideoSettings: q.getImageOrVideoSettings(),
    }

}

func(q QuestionDataModels) GetAddressQuestionModelData() AddressDataModel{
    return AddressDataModel{
        Question: "...",
        Description: "",
        Address: "Wangoi",
        AddressLine2: "Wahengbam Leikai",
        City: "Imphal",
        Zipcode: "795009",
        Country: "India",
    }
}


// Update
