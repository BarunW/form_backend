package types

type ResponseAnswers struct {
	Answer map[string]string `json:"answer"`
}

// map[questionId]withAnswer
type Answers map[int]any
