package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"
)

type QuestionTemplate struct {
	Id        string          `json:"template_id"`
	DataId    int             `json:"-"`
	SettingId int             `json:"-"`
	Title     string          `json:"title"`
	IsPaid    bool            `json:"paid"`
	Data      []byte          `json:"data"`
	Setting   []byte          `json:"setting"`
    QuestionId int64          `json:"qId"` 
}

func (s *PostgresStore) createTemplDataSettingsTable() error {
	if exist := s.isTableExist("template_data_settings"); exist {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Failed to begin db transaction for template_plan table", err, nil)
	}
    
    // there only there fields 
    // the questionTemplate id 
    // and the its structure and setting 
    // for a particular question
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS template_data_settings(     
            question_template_id VARCHAR(26),
            data_id INTEGER,
            setting_id INTEGER,  
            PRIMARY KEY (question_template_id, data_id, setting_id),
            FOREIGN KEY (question_template_id) REFERENCES question_templates(id),
            FOREIGN KEY (data_id) REFERENCES questempl_data(id),
            FOREIGN KEY (setting_id) REFERENCES questempl_settings(id)
        );`)

	if err != nil {
		return s.handleRollbackAndError("Failed to create template_data_settings table", err, tx)
	}

	err = insertQuestionTemplateDataSettings(s, tx)
	if err != nil {
	    fmt.Println(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return s.handleRollbackAndError("Failed to commit transaction template_data_settings table", err, tx)
	}
	return nil

}

func insertQuestionTemplateDataSettings(s *PostgresStore, tx *sql.Tx) error {
    
    var err error
    
    var stmt *sql.Stmt
	stmt, err = tx.Prepare(pq.CopyIn("template_data_settings", "question_template_id", "data_id", "setting_id",))
	if err != nil {
		return s.handleRollbackAndError("Unable to prepare statement for templatePlan table", err, tx)
	}

	defer stmt.Close()

    templs, err_ := s.GetQuestionTemplates()
	if err != nil {
		return s.handleRollbackAndError("Unable to prepare statement for temp_data_settings table", err_, tx)
	}

	for i, temp := range *templs {
		id := i + 1
		_, err := stmt.Exec(temp.Id, id, id)
		if err != nil {
			return s.handleRollbackAndError("Unable to execute statment for temp_data_settings table ", err, tx)
		}
		// use until all forms are aviable in templs
		if i == 2 {
			break
		}
	}
	if _, err := stmt.Exec(); err != nil {
		return s.handleRollbackAndError("Unable to finalize statement for temp_data_settings table", err, tx)
	}
	return nil
}

func (s *PostgresStore) GetQuestionTemplate(quesTemplId string) (*QuestionTemplate, error) {
	rows, err := s.db.Query(`
    SELECT question_template_id, data_id, setting_id, question_templates.title, question_templates.is_paid, questempl_data.data, questempl_settings.setting FROM template_data_settings 
    LEFT JOIN question_templates ON template_data_settings.question_template_id = question_templates.id 
    JOIN questempl_data ON template_data_settings.data_id = questempl_data.id 
    JOIN questempl_settings ON template_data_settings.setting_id = questempl_settings.id
    WHERE question_template_id = $1;`, quesTemplId)

	if err != nil {
		slog.Error(fmt.Sprintf("Unable to query template_plan for plan_id %v", quesTemplId), "details", err.Error())
		return nil, err
	}

	var qt QuestionTemplate
	for rows.Next() {
		err := rows.Scan(
			&qt.Id,
			&qt.DataId,
			&qt.SettingId,
			&qt.Title,
			&qt.IsPaid,
			&qt.Data,
			&qt.Setting,
		)
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to scan rows for template_data_settings for quesTemplId %v", quesTemplId), "details", err.Error())
			return nil, err
		}
	}

    qt.QuestionId = time.Now().UnixMilli()
	return &qt, nil
}
