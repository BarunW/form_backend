package types

type ImageLayout struct {
	Mobile        []int `json:"mobile"`
	Desktop       []int `json:"desktop"`
	MobileLayout  int   `json:"mobile_layout"`
	DesktopLayout int   `json:"desktop_layout"`
}

type ImageOrVideoSettings_Type struct {
	Url        string      `json:"url"`
	Type       string      `json:"type"`
	Layout     ImageLayout `json:"layout"`
	Brightness float32     `json:"brightness"`
	FocalPoint float32     `json:"focal_point"`
}
