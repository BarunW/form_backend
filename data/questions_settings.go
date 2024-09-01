package data

import (
	"encoding/json"

	"github.com/lib/pq"
	"github.com/sonal3323/form-poc/types"
)

var qsm = types.NewQuestionSettingModels()

var QuestionTemplateSettings = []QuestionTemplateSetting{
	{
		Setting: qsm.GetMCQSettingModel(),
    },
	{
        Setting: qsm.GetContactInfoSettingsModel(),	
    },
    {
        Setting: qsm.GetAddressQuestionSetting(),
    },

}

type QuestionTemplateSetting struct {
	Id      int `json:"-"`
	Setting any `json:"setting"`
}

func (s *PostgresStore) createQTemplSettingsTable() error {
	if exist := s.isTableExist("questempl_settings"); exist {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Failed to begin transaction for questemp_settings", err, nil)
	}

	createTableStmnt := `CREATE TABLE IF NOT EXISTS questempl_settings(
		id SERIAL PRIMARY KEY,
        setting JSON
    );`

	_, err = tx.Exec(createTableStmnt)
	if err != nil {
		return s.handleRollbackAndError("Unable to execute the create table statement for ques_temp_settings", err, tx)
	}

	stmt, err := tx.Prepare(pq.CopyIn("questempl_settings", "setting"))
	if err != nil {
		return s.handleRollbackAndError("Unable to prepare statement for ques_temp_settings", err, tx)
	}
	defer stmt.Close()

	for _, ts := range QuestionTemplateSettings {
		settingsByte, err := json.Marshal(&ts.Setting)
		if err != nil {
			return s.handleRollbackAndError("Unable to Marshal ques_temp_settings ", err, tx)
		}

		_, err = stmt.Exec(
			string(settingsByte),
		)
		if err != nil {
			return s.handleRollbackAndError("Unable to execute the statement for questemp_settings table", err, tx)
		}

	}

	_, err = stmt.Exec()
	if err != nil {
		return s.handleRollbackAndError("Unable to execute the final statement for questemp_settings table", err, tx)
	}

	err = tx.Commit()
	if err != nil {
		return s.handleRollbackAndError("Unable to execute the statement to at questemp_settings table", err, tx)
	}

	return nil
}
