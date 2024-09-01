package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func (s *PostgresStore) createFormResponesTable() error {
	tx, err := s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Unable to begin transaction for FromResponse Table", err, nil)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS form_response(
                id VARCHAR(13) PRIMARY KEY,
                workspace_id VARCHAR(26) REFERENCES workspace(id) ON DELETE CASCADE,
                form_id VARCHAR(26) REFERENCES created_form(id) ON DELETE CASCADE,
                date TIMESTAMP WITH TIME ZONE, 
                answer JSONB
            );`)

	if err != nil {
		return s.handleRollbackAndError("Unable to form_response Table", err, tx)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS form_response_data(
        id SERIAL PRIMARY KEY, 
        workspace_id VARCHAR(26) REFERENCES workspace(id) ON DELETE CASCADE,
        form_id VARCHAR(26) REFERENCES created_form(id) ON DELETE CASCADE,
        total_submissions INTEGER,
        total_start INTEGER
        );`)

	if err != nil {
		return s.handleRollbackAndError("Unable to form_response_data Table", err, tx)
	}

	if err := tx.Commit(); err != nil {
		s.handleRollbackAndError("Unable to commit for createing form repsonse tables", err, tx)
	}
	return nil
}

type ansType struct {
	resId   string
	formId  string
	resTime string
	answers []byte
}

func saveAnswer(tx *sql.Tx, ans ansType, s *PostgresStore) error {
	var workspaceId string
	if err := tx.QueryRow("SELECT workspace_id from created_form WHERE id=$1",
		ans.formId).Scan(&workspaceId); err != nil {
		return s.handleRollbackAndError("Unable to query created_form to get workspace_id", err, tx)
	}

	if result, err := tx.Exec(`INSERT INTO form_response(
        id, workspace_id, form_id, date, answer) VALUES($1, $2, $3, $4, $5);`,
		ans.resId,
		workspaceId,
		ans.formId,
		ans.resTime,
		ans.answers,
	); err != nil {
		return s.handleRollbackAndError("Unable to insert data to form_response table", err, tx)
	} else if n, resultErr := result.RowsAffected(); resultErr != nil || n == 0 {
		if n == 0 {
			return s.handleRollbackAndError("Unable to insert the data to form response_table", sql.ErrNoRows, tx)
		}
		return s.handleRollbackAndError("Unable to get the rows affected after inserting", resultErr, tx)
	}

	return nil
}

func updateTotalSubmissions(tx *sql.Tx, formId string, s *PostgresStore) error {
	//  update the total submissions
	if result, err := tx.Exec("UPDATE form_response_data SET total_submissions=total_submissions+1  WHERE form_id=$1",
		formId,
	); err != nil {
		return s.handleRollbackAndError("Unable to update the submission in form_response_data table", err, tx)
	} else if n, resultErr := result.RowsAffected(); resultErr != nil || n == 0 {
		if n == 0 {
			return s.handleRollbackAndError("Unable to update the submissions in form_response_data table zero rows affected", sql.ErrNoRows, tx)
		}
		return s.handleRollbackAndError("Errors while getting results.RowsAffected() for form_response_data", resultErr, tx)
	}

	return nil
}

func updateTotalStart(formId string, s *PostgresStore) error {
	//  update the total submissions
	if result, err := s.db.Exec("UPDATE form_response_data SET total_start=total_start+1  WHERE form_id=$1",
		formId,
	); err != nil {
		slog.Error("Failed to update the total_stat in form_response_data", "details", err.Error())
		return err
	} else if n, resultErr := result.RowsAffected(); resultErr != nil || n == 0 {
		if n == 0 {
			slog.Error("Zero row affected", "details", "rows not found to be updated")
			return sql.ErrNoRows
		}
		slog.Error("Unable to get the no. of rows affected", "details", resultErr.Error())
		return err
	}

	return nil
}

func insertOnFormRespData(tx *sql.Tx, formId, workSpaceId string) error {
	_, err := tx.Exec(`
                        INSERT INTO form_response_data(
                        workspace_id, form_id, total_submissions, total_start) 
                        VALUES($1, $2, $3, $4);`, workSpaceId, formId, 0, 0)
	if err != nil {
		slog.Error("Unable to insert the form_response_data", "details", err.Error())
		return err
	}

	return nil
}

func updateResponseCollected(tx *sql.Tx, accId string, s *PostgresStore) error {
	if result, err := tx.Exec("UPDATE responses SET response_collected=response_collected+1 WHERE account_id=$1",
		accId,
	); err != nil {
		return s.handleRollbackAndError("Unable to update the response_collected in response ", err, tx)
	} else if n, resultErr := result.RowsAffected(); resultErr != nil || n == 0 {
		if n == 0 {
			return s.handleRollbackAndError("Unable to update the in respones_collected table zero rows affected", sql.ErrNoRows, tx)
		}
		return s.handleRollbackAndError("Errors while getting results.RowsAffected() for responses", resultErr, tx)

	}

	return nil
}

func (s *PostgresStore) Submit(formId, accountId string, answers []byte) error {
	responseId := fmt.Sprintf("%d", time.Now().UnixMilli())
	responseTime := time.Now().Format(time.RFC3339)

	tx, err := s.db.Begin()
	if err != nil {
		slog.Error("Unable to begin db transaction", "details", err.Error())
		return err
	}
	err = saveAnswer(tx, ansType{resId: responseId, formId: formId, resTime: responseTime, answers: answers}, s)
	if err != nil {
		// already rollback the transaction
		return err
	}

	err = updateTotalSubmissions(tx, formId, s)
	if err != nil {
		return err
	}

	// as users submits the forms
	// update the response collected
	// increment by one for each submission
	err = updateResponseCollected(tx, accountId, s)

	err = tx.Commit()
	if err != nil {
		return s.handleRollbackAndError("Failed to commit for adding answer to the form_response table", err, tx)
	}

	return nil
}
