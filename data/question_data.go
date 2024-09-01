package data

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/sonal3323/form-poc/types"
)

type QuestionTemplateContent struct {
	Id   int `json:"-"`
	Data any `json:"data"`
}

var qdm = types.NewQuestionDataModels()

var QuestionTemplatesContent = []QuestionTemplateContent{
	{
		Data: qdm.GetMCQDataModel(),
	},
	{
		Data: qdm.GetContactInfoDataModel(),	
    },
    {
        Data: qdm.GetAddressQuestionModelData(),     
    },
}

func (s *PostgresStore) createQuestTemplData() error {
    var err error
	if exist := s.isTableExist("questempl_data"); exist {
		return nil
	}

    var tx *sql.Tx
	tx, err = s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Failed to begin transaction for questemp_data", err, nil)
	}

	createTableStmnt := `CREATE TABLE IF NOT EXISTS questempl_data(
		id SERIAL PRIMARY KEY,
        data JSON
    );`

	_, err = tx.Exec(createTableStmnt)
	if err != nil {
		return s.handleRollbackAndError("Unable to execute the create table statement for ques", err, tx)
	}
    
    var stmt *sql.Stmt
	stmt, err = tx.Prepare(pq.CopyIn("questempl_data", "data"))
	if err != nil {
		return s.handleRollbackAndError("Unable to prepare statement for questemp_data", err, tx)
	}

	for _, quesData := range QuestionTemplatesContent {
		dataByte, err := json.Marshal(&quesData.Data)
		fmt.Println(string(dataByte))
		if err != nil {
			return s.handleRollbackAndError("Unable to Marshal questemp_data ", err, tx)
		}
		_, err = stmt.Exec(
			string(dataByte),
		)
		if err != nil {
			return s.handleRollbackAndError("Unable to execute the statement to at questemp_data table", err, tx)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return s.handleRollbackAndError("Unable to do final execution of statement   for questemp_data table ", err, tx)
	}

	err = tx.Commit()
	if err != nil {
		return s.handleRollbackAndError("Unable to commit the transaction for questemp_data table", err, tx)
	}

	return nil
}
