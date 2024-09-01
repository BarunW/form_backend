package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sonal3323/form-poc/data/questionsUpdates"
	"github.com/sonal3323/form-poc/types"
)

type PostgresStore struct {
	db *sql.DB
	QU *questionsUpdates.QuestionCRUD
}

/*
==================
While creating table put table which has foreign key
at bottom
=================
*/

func (s *PostgresStore) Init() error {
	err := s.createENUMS()
	if err != nil {
		return err
	}

	err = s.createAccountTable()
	if err != nil {
		return err
	}

	err = s.createWorkspaceTable()
	if err != nil {
		return err
	}

	err = s.createPricingTable()
	if err != nil {
		return err
	}

	err = s.createResponseLimitTable()
	if err != nil {
		return err
	}

	err = s.createResponsesTable()
	if err != nil {
		return err
	}

	// subscription depends on three foreign key above tables
	err = s.createSubscriptionTable()
	if err != nil {
		return err
	}

	if err := s.handleQuestionTemplates(); err != nil {
		return err
	}

	// all the created form will be on this table
	err = s.createdFormTable()
	if err != nil {
		return err
	}

	if err := s.createFormResponesTable(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) handleQuestionTemplates() error {

	err := s.createQuestionTemplatesTable()
	if err != nil {
		return err
	}

	err = s.createQuestTemplData()
	if err != nil {
		return err
	}

	err = s.createQTemplSettingsTable()
	if err != nil {
		return err
	}

	// plan and question_template junction table
	err = s.createTemplPlanTable()
	if err != nil {
		return err
	}

	// question_template, data, settings  junction table
	err = s.createTemplDataSettingsTable()
	if err != nil {
		return err
	}

	return nil
}

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "mysecretpassword"
	dbName     = "postgres"
)

func NewPostgresStore() (*PostgresStore, error) {
	constStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	var db *sql.DB
	var err error
    
	//	<-time.After(40 * time.Second)
	// Retry loop
	for retries := 0; retries < 5; retries++ {
		db, err = sql.Open("postgres", constStr)
		if err != nil {
			fmt.Println("Error connecting to database: Retrying", retries)
			<-time.After(12 * time.Second)
		} else {
			break // Connected successfully
		}
	}

	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		slog.Error("ERROR AT pinging Db", "details", err.Error())
		fmt.Println(err)
		return nil, err
	}

	slog.Info("postgres", "database connection", "successfull")
	qu := questionsUpdates.NewQuestionCRUD(db)

	return &PostgresStore{
		db: db,
		// QU updates the question fields
		QU: qu,
	}, nil

}

func (s *PostgresStore) isTableExist(tableName string) bool {
	var exists bool
	err := s.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_name = $1
		);
	`, tableName).Scan(&exists)

	if err != nil {
		slog.Error("unable to check the table", "details", err.Error())
		os.Exit(1)
	}

	return exists
}

func (s *PostgresStore) handleRollbackAndError(msg string, err error, tx *sql.Tx) error {
	if err != nil {
		slog.Error(msg, "details", err.Error())
		if tx != nil {
			tx.Rollback()
		}
	}
	return err
}

func (s *PostgresStore) createENUMS() error {
	tx, err := s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Unable to create begin the transaction to create enum", err, nil)
	}

	// Plan status
	_, err = tx.Exec(`DO $$ BEGIN CREATE TYPE plan_status AS ENUM ('PAUSE', 'ACTIVE'); EXCEPTION WHEN duplicate_object THEN null; END $$`)
	if err != nil {
		return s.handleRollbackAndError("Failed to create plan_status enum type", err, tx)
	}

	// Billing period enum
	_, err = tx.Exec(`DO $$ BEGIN CREATE TYPE bp AS ENUM ('Monthly', 'Yearly'); EXCEPTION WHEN duplicate_object THEN null; END $$`)
	if err != nil {
		return s.handleRollbackAndError("Failed to create the Billing period enum", err, tx)
	}

	// creating currency enum
	_, err = tx.Exec("DO $$ BEGIN CREATE TYPE currency_code AS ENUM('USD', 'INR');EXCEPTION WHEN duplicate_object THEN null; END $$")
	if err != nil {
		return s.handleRollbackAndError("Unable to create currency Enum", err, tx)
	}

	// creating pricing plan enum( free, basic, plus, businses, enterprise)
	createEnum := fmt.Sprintf("DO $$ BEGIN CREATE TYPE pricing_plan AS ENUM('%s', '%s', '%s', '%s', '%s');EXCEPTION WHEN duplicate_object THEN null; END $$",
		string(types.Free), string(types.Basic), string(types.Plus), string(types.Business), string(types.Enterprise))
	_, err = tx.Exec(createEnum)
	if err != nil {
		return s.handleRollbackAndError("Unable to create enum", err, tx)
	}

	err = tx.Commit()
	if err != nil {
		return s.handleRollbackAndError("Unable to commit transaction for enum ", err, tx)
	}

	return nil
}

func (s *PostgresStore) checkRowsAffected(result sql.Result) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}
