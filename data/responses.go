package data

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type ResponsesModel struct {
	UserId            int    `json:"user_id"`
	AccountId         string `json:"account_id"`
	ResponseLimit     int    `json:"response_limit"`
	ResponseCollected int    `json:"response_collected"`
}

func (s *PostgresStore) createResponsesTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS responses(
                            id SERIAL PRIMARY KEY,
                            user_id INTEGER REFERENCES users(id),
                            account_id VARCHAR(11),
                            response_limit INTEGER,
                            response_collected INTEGER
                        );`)
	if err != nil {
		slog.Error("Unable to create table for responses", "details", err.Error())
		return err
	}

	return nil
}

func InsertResponsesWithTx(tx *sql.Tx, rm ResponsesModel) error {
	_, err := tx.Exec(`INSERT INTO responses( user_id, account_id, response_limit, response_collected)
                        VALUES($1, $2, $3, $4);`, rm.UserId, rm.AccountId, rm.ResponseLimit, rm.ResponseCollected)
	if err != nil {
		fmt.Println("Unable to execute the insert statement")
		return err
	}

	return nil
}

func (s *PostgresStore) GetResponses(userId int) (*ResponsesModel, error) {
	var responses ResponsesModel
	err := s.db.QueryRow("SELECT user_id, response_limit, response_collected FROM responses WHERE user_id=$1",
		userId).Scan(&responses.UserId, &responses.ResponseLimit, &responses.ResponseCollected)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &responses, nil
}
