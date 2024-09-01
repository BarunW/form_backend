package types


// DO NOT EDIT THIS CODE
type qt_id string
type path string


const (
	MultipleChoice_ID   qt_id = "01HPK5P93048XX5KQTMDFKP9SM"
	ContactInfo_ID      qt_id = "01HPK5P930BF7AXZZNXFTTM8RB"
	PhoneNumber_ID      qt_id = "01HPK5P930VPBGTHC13X4Z0M7V"
	Short_ID            qt_id = "01HPK5P930SQWX8562KPGREVWK"
	LongText_ID         qt_id = "01HPK5P9302D8FXP9EMER9229B"
	Statement_ID        qt_id = "01HPK5P930CZET4A5V03WY2K14"
	PictureChoice_ID    qt_id = "01HPK5P930XRGSXTK0DFF1N3KC"
	Ranking_ID          qt_id = "01HPK5P9306GN80JXS1BP0BFCZ"
	YesNO_ID            qt_id = "01HPK5P930MG7ZW43NWK3YZ85B"
	Email_ID            qt_id = "01HPK5P9304B0Q7TW2YW8P5Q93"
	OpinionScale_ID     qt_id = "01HPK5P930CXS1N7QSY90HASP7"
	NetPromoterScore_ID qt_id = "01HPK5P931KQZK33HB1H7ZB65A"
	Rating_ID           qt_id = "01HPK5P9316NG25HV4ZH4WTPHA"
	Matrix_ID           qt_id = "01HPK5P931C2B8CTS1EAGPXF65"
	Date_ID             qt_id = "01HPK5P931SVN96PM8FDC83MMX"
	Number_ID           qt_id = "01HPK5P931YH7BEE704W79CG4C"
	Dropdown_ID         qt_id = "01HPK5P931PH4SFWP94ZV3478Q"
	Legal_ID            qt_id = "01HPK5P931ACQSYAC550T51N4A"
	FileUpload_ID       qt_id = "01HPK5P931DPGA7KMPVGTYM931"
	Payment_ID          qt_id = "01HPK5P931EMYPJ2GAPAH10Q7K"
	Website_ID          qt_id = "01HPK5P931N5SXPKVMVDWKZ79W"
	Calendly_ID         qt_id = "01HPK5P931P6YTDEQZ5969WYQB"
)

var qt_idPathTable = map[qt_id]path{
    MultipleChoice_ID  : "mcq",  
    ContactInfo_ID     : "contact-info",  
    PhoneNumber_ID     : "phone-number", 
    Short_ID           : "short", 
    LongText_ID        : "long-text", 
    Statement_ID       : "statement",
    PictureChoice_ID   : "picture-choice", 
    Ranking_ID         : "ranking", 
    YesNO_ID           : "yes-no", 
    Email_ID           : "email", 
    OpinionScale_ID    : "opinion-scale", 
    NetPromoterScore_ID: "promoter-score", 
    Rating_ID          : "rating", 
    Matrix_ID          : "maxtrix", 
    Date_ID            : "date", 
    Number_ID          : "number", 
    Dropdown_ID        : "drop-down", 
    Legal_ID           : "legal", 
    FileUpload_ID      : "file-upload", 
    Payment_ID         : "payment", 
    Website_ID         : "website", 
    Calendly_ID        : "calendly", 
}

func GetQTID_Path( id string) string{  
    newId := id[1:27] 
    urlPath, exist := qt_idPathTable[qt_id(newId)]
    if !exist{
        return ""
    }
    return string(urlPath)
}
