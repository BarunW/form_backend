package data

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/lib/pq"

	"github.com/sonal3323/form-poc/types"
)

var responseLimitIdsForBasic = map[int]types.ResponseLimitId{
	1: "100-BASIC-RESPLIMIT",
	2: "250-BASIC-RESPLIMIT",
	3: "500-BASIC-RESPLIMIT",
	4: "750-BASIC-RESPLIMIT",
}

var responseLimitIdsForPlus = map[int]types.ResponseLimitId{
	5: "1000-PLUS-RESPLIMIT",
	6: "1500-PLUS-RESPLIMIT",
	7: "2000-PLUS-RESPLIMIT",
	8: "2500-PLUS-RESPLIMIT",
}

var responseLimitIdsForBusiness = map[int]types.ResponseLimitId{
	9:  "10000-BUSINESS-RESPLIMIT",
	10: "15000-BUSINESS-RESPLIMIT",
	11: "20000-BUSINESS-RESPLIMIT",
	12: "25000-BUSINESS-RESPLIMIT",
}

// this is for number of add on response for each plan so that a form can be response by
var cacheResponseLimitId = []map[types.ResponseLimitId]bool{
	{
		responseLimitIdsForBasic[1]: true,
		responseLimitIdsForBasic[2]: true,
		responseLimitIdsForBasic[3]: true,
		responseLimitIdsForBasic[4]: true,
	},
	{
		responseLimitIdsForPlus[5]: true,
		responseLimitIdsForPlus[6]: true,
		responseLimitIdsForPlus[7]: true,
		responseLimitIdsForPlus[8]: true,
	},
	{
		responseLimitIdsForBusiness[9]:  true,
		responseLimitIdsForBusiness[10]: true,
		responseLimitIdsForBusiness[11]: true,
		responseLimitIdsForBusiness[12]: true,
	},
}

func (s *PostgresStore) GetCachedRespLimitId() []map[types.ResponseLimitId]bool {
	return cacheResponseLimitId
}

func getLimit(i int) types.ResponseLimitNum {
	switch i {
	case 1:
		return 100
	case 2:
		return 250
	case 3:
		return 500
	case 4:
		return 750
	case 5:
		return 1000
	case 6:
		return 1500
	case 7:
		return 2000
	case 8:
		return 2500
	case 9:
		return 10000
	case 10:
		return 15000
	case 11:
		return 20000
	case 12:
		return 25000
	}
	return 0
}

// Response Plan for plus
func (s *PostgresStore) createResponseLimitTable() error {
	if exist := s.isTableExist("responseaddon"); exist {
		return nil
	}

	createTableStmnt := `CREATE TABLE IF NOT EXISTS responseaddon(
		id VARCHAR(25) PRIMARY KEY,
        response_limit INTEGER,
        ppr FLOAT,
        plan pricing_plan, 
        currency currency_code
    );`
	// begin the transaction
	tx, err := s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Unable to begin transaction for responseaddon", err, nil)
	}
	_, err = tx.Exec(createTableStmnt)
	if err != nil {
		return s.handleRollbackAndError("Unable to create responseaddon", err, tx)
	}

	err = insertToResponseAddon(s, tx)
	if err != nil {
		//alredy rollback so just return the error
		return err
	}

	err = tx.Commit()
	if err != nil {
		return s.handleRollbackAndError("Failed to commit transaction, responseaddon", err, tx)
	}
	return nil
}

type responseLimitData struct {
	Id            types.ResponseLimitId  `json:"id"`
	ResponseLimit types.ResponseLimitNum `json:"response_limit"`
	PPR           float32                `json:"ppr"`
	Plan          types.PlanType         `json:"plan"`
	Currency      types.CurrencyCode     `json:"Currency"`
}

func insertToResponseAddon(s *PostgresStore, tx *sql.Tx) error {
	//prepare the copy statement
	stmt, err := tx.Prepare(pq.CopyIn("responseaddon", "id", "response_limit", "ppr", "plan", "currency"))
	if err != nil {
		return s.handleRollbackAndError("Unable to begin copy , templates table", err, tx)
	}
	defer stmt.Close()

	for i := 1; i <= 12; i++ {
		respLimitData := &responseLimitData{}
		if i >= 1 && i <= 4 {
			respLimitData.Id = responseLimitIdsForBasic[i]
			respLimitData.ResponseLimit = getLimit(i)
			respLimitData.PPR = 0.1
			respLimitData.Plan = types.Basic
			respLimitData.Currency = types.USD
		}
		if i >= 5 && i <= 8 {
			respLimitData.Id = responseLimitIdsForPlus[i]
			respLimitData.ResponseLimit = getLimit(i)
			respLimitData.PPR = 0.1
			respLimitData.Plan = types.Plus
			respLimitData.Currency = types.USD
		}
		if i >= 9 && i <= 12 {
			respLimitData.Id = responseLimitIdsForBusiness[i]
			respLimitData.ResponseLimit = getLimit(i)
			respLimitData.PPR = 0.1
			respLimitData.Plan = types.Business
			respLimitData.Currency = types.USD
		}
		_, err = stmt.Exec(
			respLimitData.Id,
			respLimitData.ResponseLimit,
			respLimitData.PPR,
			respLimitData.Plan,
			respLimitData.Currency,
		)
		if err != nil {
			return s.handleRollbackAndError("Unable to copy the statement, responseaddon", err, tx)
		}
	}
	// Execute the copy stament
	_, err = stmt.Exec()
	if err != nil {
		return s.handleRollbackAndError("Unable to execute copy statement, responseaddon", err, tx)
	}
	return nil
}

func (s *PostgresStore) GetAllResponseLimit() ([]responseLimitData, error) {
	rows, err := s.db.Query(`SELECT * FROM responseaddon;`)
	if err != nil {
		slog.Error("Unable to perfrom select responseaddon table", "details", err.Error())
		return nil, err
	}
	return scanIntoAllRow(rows)
}

func (s *PostgresStore) getResponseLimit(id types.ResponseLimitId) (*responseLimitData, error) {
	rows, err := s.db.Query(`SELECT * FROM responseaddon WHERE id=$1;`, id)
	if err != nil {
		slog.Error("Unable to perfrom select responseaddon table", "details", err.Error())
		return nil, err
	}
	data, err := scanIntoAllRow(rows)
	if err != nil {
		slog.Error("Unable to scan the rows for getResponseLimit", "details", err.Error())
		return nil, err
	} else if data == nil {
		slog.Error("no data is return from scanIntoAllRow ")
		code := "404"
		return nil, fmt.Errorf("%s", code)
	}
	return &data[0], nil
}

func scanIntoAllRow(rows *sql.Rows) ([]responseLimitData, error) {
	var rld []responseLimitData
	for rows.Next() {
		r := responseLimitData{}
		err := rows.Scan(
			&r.Id,
			&r.ResponseLimit,
			&r.PPR,
			&r.Plan,
			&r.Currency,
		)
		if err != nil {
			slog.Error("Unable to scan the rows for responseaddon", "details", err.Error())
			return nil, err
		}
		rld = append(rld, r)
	}
	return rld, nil
}

func DefaultReponseAddon() (int, int) {
	const default_response_addon int = 10
	const response_collected = 0
	return default_response_addon, response_collected
}
