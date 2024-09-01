package data

import (
	"database/sql"
	"log/slog"

	"github.com/lib/pq"

	"github.com/sonal3323/form-poc/types"
)

const (
	Paid types.IsPaid = true
)

type questionMetaData struct {
	Id     string       `json:"id"`
	Title  string       `json:"title"`
	IsPaid types.IsPaid `json:"paid"`
}


var questionsMetaDatas = []questionMetaData{
	{
		Id:     "01HPK5P93048XX5KQTMDFKP9SM",
		Title:  "Multiple Choice",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P930BF7AXZZNXFTTM8RB",
		Title:  "Contact Info",
		IsPaid: !Paid,
	},
	{

		Id:     "01HPK5P930VPBGTHC13X4Z0M7V",
		Title:  "Phone Number",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P930SQWX8562KPGREVWK",
		Title:  "Short Text",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P9302D8FXP9EMER9229B",
		Title:  "Long Text",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P930CZET4A5V03WY2K14",
		Title:  "Statement",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P930XRGSXTK0DFF1N3KC",
		Title:  "Picture Choice",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P9306GN80JXS1BP0BFCZ",
		Title:  "Ranking",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P930MG7ZW43NWK3YZ85B",
		Title:  "Yes/No",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P9304B0Q7TW2YW8P5Q93",
		Title:  "Email",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P930CXS1N7QSY90HASP7",
		Title:  "Opinion Scale",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P931KQZK33HB1H7ZB65A",
		Title:  "Net Promoter Score®",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P9316NG25HV4ZH4WTPHA",
		Title:  "Rating",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P931C2B8CTS1EAGPXF65",
		Title:  "Matrix",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P931SVN96PM8FDC83MMX",
		Title:  "Date",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P931YH7BEE704W79CG4C",
		Title:  "Number",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P931PH4SFWP94ZV3478Q",
		Title:  "Dropdown",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P931ACQSYAC550T51N4A",
		Title:  "Legal",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P931DPGA7KMPVGTYM931",
		Title:  "File Upload",
		IsPaid: Paid,
	},
	{
		Id:     "01HPK5P931EMYPJ2GAPAH10Q7K",
		Title:  "Payment",
		IsPaid: Paid,
	},
	{
		Id:     "01HPK5P931N5SXPKVMVDWKZ79W",
		Title:  "Website",
		IsPaid: !Paid,
	},
	{
		Id:     "01HPK5P931P6YTDEQZ5969WYQB",
		Title:  "Calendly",
		IsPaid: !Paid,
	},
}

func (s *PostgresStore) createQuestionTemplatesTable() error {
	if exist := s.isTableExist("question_templates"); exist {
		return nil
	}

    var(
        tx *sql.Tx
        stmt *sql.Stmt
        err error
    )
    
	// begin the transaction
	tx, err = s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Unable to begin transaction, templates table", err, nil)
	}
	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS question_templates(
        id VARCHAR(26) PRIMARY KEY,
        title VARCHAR(30),
        is_paid BOOLEAN 
    );`)

	if err != nil {
		return s.handleRollbackAndError("Unable to create table, templates table", err, tx)
	}

	// Prepare the COPY statement
	stmt, err = tx.Prepare(pq.CopyIn("question_templates", "id", "title", "is_paid"))
	if err != nil {
		return s.handleRollbackAndError("Unable to prepare for CopyIn , templates table", err, tx)
	}
	defer stmt.Close()

	// Iterate over the templates and add rows to the COPY statement
	for _, temp := range questionsMetaDatas {
		_, err = stmt.Exec(temp.Id, temp.Title, temp.IsPaid)
		if err != nil {
			return s.handleRollbackAndError("Unable to copy the statement", err, tx)
		}
	}

	// Execute the COPY statement
	_, err = stmt.Exec()
	if err != nil {
		return s.handleRollbackAndError("Unable to execute copy statement", err, tx)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return s.handleRollbackAndError("Failed to commit", err, tx)
	}
	return nil
}

func (s *PostgresStore) GetQuestionTemplates() (*[]questionMetaData, error) {
	questionTempls := []questionMetaData{}
	rows, err := s.db.Query(`SELECT * from question_templates;`)
	if err != nil {
		slog.Error("Unable to query for Question-Templates", "details", err.Error())
		return nil, err
	}
	for rows.Next() {
		questionTempl := questionMetaData{}
		err := rows.Scan(&questionTempl.Id, &questionTempl.Title, &questionTempl.IsPaid)
		if err != nil {
			slog.Error("Unable to get Question-template ", "details", err.Error())
			return nil, err
		}
		questionTempls = append(questionTempls, questionTempl)
	}

	return &questionTempls, nil
}
