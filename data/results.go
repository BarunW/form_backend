package data

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
)


type respondedAnswer map[string]map[string]interface{} 

type AnswerStruct struct{
    Answers []interface{} `json:"answers"`
    Question string `json:"question"`
}

type AnswerType map[string]AnswerStruct

type response struct{ 
	Date    time.Time `json:"date"`
	Answer  []byte `json:"answer"`
}

func (s *PostgresStore) GetResult(formId string) (map[string]AnswerStruct, error) {
	var (
        contents []byte 
        err error
    )

    var rows *sql.Rows
    rows, err = s.db.Query(`SELECT date, answer FROM form_response WHERE form_id=$1`, formId)
    if err != nil{
        return nil, err    
    }

    var responses []response
    responses, err  = scanIntoFormResponses(rows)
    
    var a map[string][]interface{}
    a, err = seperateResponsesByQId(responses) 
    if err != nil{
        return nil, err
    }

	err = s.db.QueryRow(`select content from created_form where id=$1;`, formId).Scan( &contents)  
    if err != nil{
        return nil, err
    }


    var questions []CreatedQuestionTemplate
    questions, err = jsonToContentType(contents)
    
    answers := make(map[string]AnswerStruct)
    
    // the index is correspond to the key of the
    // CreatedQuestionTemplate(i.e quest here)
    // key == index 
    for index, quest := range questions{
        var v map[string]interface{} 
        err := json.Unmarshal(quest[index].Data, &v )
        if err != nil{
            return nil, err
        }
        // as qid is unique 
        qId := strconv.Itoa(int(quest[index].QuestionId))
        if _, ok := a[qId]; ok{
            answers[qId] = AnswerStruct{
                Answers: a[qId],
                Question: v["question"].(string),
            }
        }
    }
    
	return answers, nil
}



func  seperateResponsesByQId(respondedAnswers []response) ( map[string][]interface{}, error){
    var err error 

    answers := make(map[string][]interface{})

    for _, r := range respondedAnswers{    
        ans := make(respondedAnswer) 
        err = json.Unmarshal(r.Answer, &ans)
        if err != nil{
            return nil, nil
        }

        for key, value:= range ans{
            // the key "answer" always exist otherwise failed at handler level
            // when the response is submit by the client
            if value["answer"] == nil {
                continue
            }

            if _, exist := answers[key]; !exist {
                   answers[key] = []interface{}{value}
            } else if exist {
                answers[key] = append(answers[key], value)
            }
            
        }
    }
    return answers, nil
}


func scanIntoFormResponses(rows *sql.Rows) ([]response, error) {
    var responses []response 
    for rows.Next(){
       r := response{}
       err := rows.Scan(&r.Date, &r.Answer) 
       if err != nil{
           return nil, err
       }
       responses = append(responses, r)
    }

    return responses, nil
}
