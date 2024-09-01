package questionsUpdates

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/sonal3323/form-poc/types"
	"github.com/sonal3323/form-poc/utils"
)

type QuestionCRUD struct {
	db *sql.DB
}

func NewQuestionCRUD(db *sql.DB) *QuestionCRUD {
	return &QuestionCRUD{
		db: db,
	}
}

func (s *QuestionCRUD) handleRollbackAndError(msg string, err error, tx *sql.Tx) error {
	if err != nil {
		slog.Error(msg, "details", err.Error())
		if tx != nil {
			tx.Rollback()
		}
	}
	return err
}

type Choice struct {
	Label string `json:"label"`
}

func (q *QuestionCRUD) IsPathMatchWithQuestionId(questionId int, formId, path string) error {
	var qtId string
	query := fmt.Sprintf("SELECT content->%d->'%d'->'template_id' FROM created_form WHERE id=$1", questionId, questionId)
	err := q.db.QueryRow(query, formId).Scan(&qtId)
	if err != nil {
		slog.Error("Unable to query the row", "details", err.Error())
		return err
	}
	fmt.Println(qtId, types.MultipleChoice_ID)
	if path != types.GetQTID_Path(qtId) {
		return sql.ErrNoRows
	}

	return nil
}

func (q *QuestionCRUD) AddChoice(questionId int, text, formId string) error {
	tx, err := q.db.Begin()
	if err != nil {
		slog.Error("Unable to begin transaction for add choice", "detaisl", err.Error())
		return err
	}

	// query the total choice
	queryString := fmt.Sprintf("select content->%d->'%d'->'data'->'total_choice' from created_form WHERE id='%s';", questionId, questionId, formId)
	var length int
	err = tx.QueryRow(queryString).Scan(&length)
	if err != nil {
		return q.handleRollbackAndError("Unable to get the length of total choice", err, tx)
	}

	// create a new choice
	newChoice, err := json.Marshal(Choice{Label: text})
	if err != nil {
		slog.Error("Unable to marshal newChoice", "details", err.Error())
		return err
	}

	// generating new Id
	newChoiceId, err := utils.GenerateULID()
	if err != nil {
		return q.handleRollbackAndError("Unable to generate new ULID", err, tx)
	}

	// add choice
	updatedStmnt := fmt.Sprintf(`UPDATE created_form SET content=jsonb_set(content, '{%d,%d,data,choices,%s}', '%s') 
                    WHERE id = '%s';`, questionId, questionId, newChoiceId, string(newChoice), formId)

	_, err = tx.Exec(updatedStmnt)
	if err != nil {
		return q.handleRollbackAndError("Unable to update the the choices", err, tx)
	}

	// update the total choice
	updateTotalChoice := fmt.Sprintf(`UPDATE created_form
    SET content = jsonb_set(content, '{%d,%d,data,total_choice}','%d' ,false) 
    WHERE id='%s';`, questionId, questionId, length+1, formId)

	_, err = tx.Exec(updateTotalChoice)
	if err != nil {
		return q.handleRollbackAndError("Unable to update the the total_choice", err, tx)
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return q.handleRollbackAndError("Unable to commit for add choice", err, tx)
	}

	return nil
}

func (q *QuestionCRUD) UpdateChoiceAnswer(questionId int, choiceId, text, formId string) error {

	updateStatment := fmt.Sprintf(`UPDATE created_form
    SET content = jsonb_set(content, '{%d,%d,data,choices,%s,label}', to_jsonb('%s'::text), false) 
    WHERE id='%s';`, questionId, questionId, choiceId, text, formId)
	fmt.Println(updateStatment)
	_, err := q.db.Exec(updateStatment)
	if err != nil {
		slog.Error("Unable to update the question", "details", err.Error())
		return err
	}

	return nil
}

func (q *QuestionCRUD) DeleteChoice(questionId int, choiceId, formId string) error {
	tx, err := q.db.Begin()
	if err != nil {
		slog.Error("Unable to begin transaction for add choice", "detail", err.Error())
		return err
	}
	// delete the choice
	queryString := fmt.Sprintf(`UPDATE created_form SET content=jsonb_set(content,'{%d,%d,data,choices}', 
                    to_jsonb(content->%d->'%d'->'data'->'choices'::text)-'%s', false) WHERE id='%s';`,
		questionId, questionId, questionId, questionId, choiceId, formId)

	fmt.Println(queryString)
	_, err = tx.Exec(queryString)
	if err != nil {
		return q.handleRollbackAndError("Unable to get the remove choice", err, tx)
	}

	var length int
	// get the total  choice
	getTotalChoice := fmt.Sprintf("select content->%d->'%d'->'data'->'total_choice' from created_form WHERE id='%s';", questionId, questionId, formId)
	err = tx.QueryRow(getTotalChoice).Scan(&length)
	if err != nil {
		return q.handleRollbackAndError("Unable to get the the total_choice", err, tx)
	}

	// update the total choice
	updateTotalChoice := fmt.Sprintf(`UPDATE created_form
    SET content = jsonb_set(content, '{%d,%d,data,total_choice}','%d' ,false) 
    WHERE id='%s';`, questionId, questionId, length-1, formId)

	_, err = tx.Exec(updateTotalChoice)
	if err != nil {
		return q.handleRollbackAndError("Unable to update the the total_choice", err, tx)
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return q.handleRollbackAndError("Unable to commit for add choice", err, tx)
	}
	return nil
}

/*
===============
 Settings
===============
*/

func (s *QuestionCRUD) UpdateOtherOption(formId string, questionId int, setting bool) error {
	updateQuery := fmt.Sprintf("UPDATE created_form SET content=jsonb_set(content, '{%d,%d,setting,other_option}', '%v'::jsonb, false) WHERE id='%s';",
		questionId, questionId, setting, formId)

	_, err := s.db.Exec(updateQuery)
	if err != nil {
		slog.Error("Unable to update the required setting", "details", err.Error())
		return err
	}

	return nil
}
