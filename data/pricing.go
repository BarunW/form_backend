package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lib/pq"

	"github.com/sonal3323/form-poc/types"
)

// planids with monthly billing period pricing
// this constant variables add at latter part of development so there is some
// inconvient to type-cast untl if refactor

const (
	FREE_PLAN        types.PlanId_Type = 1121
	BASIC_MONTHLY    types.PlanId_Type = 1122
	PLUS_MONTHLY     types.PlanId_Type = 1123
	BUSINESS_MONTHLY types.PlanId_Type = 1124
)

// planids with yearly billing period pricing
const (
	BASIC_YEARLY    types.PlanId_Type = 2122
	PLUS_YEARLY     types.PlanId_Type = 2123
	BUSINESS_YEARLY types.PlanId_Type = 2124
)

func getPricingData() []types.Pricing {
	return []types.Pricing{
		// free
		{
			PlanId:        1121,
			Plan:          types.Free,
			Price:         0,
			BillingPeriod: types.Month,
			Currency:      types.USD,
			Features: types.Features{
				Discount:          0,
				Seats:             1,
				MaxResponse:       10,
				UnlimitedForm:     true,
				AcceptPayment:     false,
				RecieveFileUpload: false,
				RemovingBranding:  false,
				CreateOwnBranding: false,
			},
		},
		// Basic plan
		{
			PlanId:        1122,
			Plan:          types.Basic,
			Price:         20,
			BillingPeriod: types.Month,
			Currency:      types.USD,
			Features: types.Features{
				Discount:          58,
				Seats:             1,
				MaxResponse:       100,
				UnlimitedForm:     true,
				AcceptPayment:     true,
				RecieveFileUpload: true,
				RemovingBranding:  false,
				CreateOwnBranding: false,
			},
		},
		// Plus plan
		{
			PlanId:        1123,
			Plan:          types.Plus,
			Price:         70,
			BillingPeriod: types.Month,
			Currency:      types.USD,
			Features: types.Features{
				Discount:          58,
				Seats:             1,
				MaxResponse:       1000,
				UnlimitedForm:     true,
				AcceptPayment:     true,
				RecieveFileUpload: true,
				RemovingBranding:  true,
				CreateOwnBranding: true,
			},
		},
		// Businses plan
		{
			PlanId:        1124,
			Plan:          types.Business,
			Price:         100,
			BillingPeriod: types.Month,
			Currency:      types.USD,
			Features: types.Features{
				Discount:          118,
				Seats:             3,
				MaxResponse:       1000,
				UnlimitedForm:     true,
				AcceptPayment:     true,
				RecieveFileUpload: true,
				RemovingBranding:  true,
				CreateOwnBranding: true,
			},
		},
		// Basic plan yearly
		{
			PlanId:        2122,
			Plan:          types.Basic,
			Price:         20,
			BillingPeriod: types.Year,
			Currency:      types.USD,
			Features: types.Features{
				Discount:          58,
				Seats:             1,
				MaxResponse:       100,
				UnlimitedForm:     true,
				AcceptPayment:     true,
				RecieveFileUpload: true,
				RemovingBranding:  false,
				CreateOwnBranding: false,
			},
		},

		// Plus plan yearly
		{
			PlanId:        2123,
			Plan:          types.Plus,
			Price:         70,
			BillingPeriod: types.Year,
			Currency:      types.USD,
			Features: types.Features{
				Discount:          58,
				Seats:             1,
				MaxResponse:       1000,
				UnlimitedForm:     true,
				AcceptPayment:     true,
				RecieveFileUpload: true,
				RemovingBranding:  true,
				CreateOwnBranding: true,
			},
		},
		// Businses plan yearly
		{
			PlanId:        2124,
			Plan:          types.Business,
			Price:         100,
			BillingPeriod: types.Year,
			Currency:      types.USD,
			Features: types.Features{
				Discount:          118,
				Seats:             3,
				MaxResponse:       1000,
				UnlimitedForm:     true,
				AcceptPayment:     true,
				RecieveFileUpload: true,
				RemovingBranding:  true,
				CreateOwnBranding: true,
			},
		},
	}

}

/*
======================
we use os.exit(1) at error
as guard clause  becuase when server starts its creates pricing table
and insert the data. so it fails first before serving
======================
*/
func (s *PostgresStore) createPricingTable() error {
	if exist := s.isTableExist("pricing"); exist {
		return nil
	}

	createTableStmnt := `CREATE TABLE IF NOT EXISTS pricing(
		id INTEGER PRIMARY KEY,
        plan pricing_plan ,
        price FLOAT , 
        billing_period bp,
        currency currency_code,
        features JSON
    );`

	// Begin the transaction
	tx, err := s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Failed to begin transaction for pricing", err, nil)
	}

	// create table
	_, err = tx.Exec(createTableStmnt)
	if err != nil {
		return s.handleRollbackAndError("Failed to create pricing table", err, tx)
	}

	if err := insertPricingData(s, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return s.handleRollbackAndError("Failed to commit transaction for pricing", err, nil)
	}

	return nil
}

func insertPricingData(s *PostgresStore, tx *sql.Tx) error {
	stmt, err := tx.Prepare(pq.CopyIn("pricing", "id", "plan", "price", "billing_period", "currency", "features"))
	if err != nil {
		return s.handleRollbackAndError("Unabled to prepare the statement for pricing table", err, tx)
	}
	defer stmt.Close()

	plans := getPricingData()
	for _, plan := range plans {
		featureByte, err := json.Marshal(plan.Features)
		if err != nil {
			return s.handleRollbackAndError("Unable to marshal the feature", err, tx)
		}
		_, err = stmt.Exec(
			plan.PlanId,
			plan.Plan,
			plan.Price,
			plan.BillingPeriod,
			plan.Currency,
			string(featureByte),
		)
		if err != nil {
			return s.handleRollbackAndError("Unable to execute the pricing insert statement", err, tx)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return s.handleRollbackAndError("Uable to finalize execution, pricing", err, tx)
	}
	return nil
}

func (s *PostgresStore) getPricing(id int) (*types.Pricing, error) {
	rows, err := s.db.Query("SELECT * FROM pricing WHERE id=$1", id)
	if err != nil {
		slog.Error("Unable to query the data from pricing table", "details", err.Error())
		return nil, err
	}
	defer rows.Close()
	price, err := scanIntoPricing(rows)
	if err != nil {
		return nil, err
	} else if len(price) < 1 {
		return nil, fmt.Errorf("%d not able to get the price data ", id)
	}
	return &price[0], nil

}
func (s *PostgresStore) GetAllPricing() ([]types.Pricing, error) {
	rows, err := s.db.Query("SELECT * FROM pricing")
	if err != nil {
		slog.Error("Unable to query pricing table", "details", err.Error())
		return nil, err
	}
	defer rows.Close()

	allPricing, err := scanIntoPricing(rows)
	if err != nil {
		return nil, err
	}

	return allPricing, nil
}

func scanIntoPricing(rows *sql.Rows) ([]types.Pricing, error) {
	result := []types.Pricing{}
	for rows.Next() {
		var features []byte
		pricing := types.Pricing{
			Features: types.Features{},
		}
		err := rows.Scan(
			&pricing.PlanId,
			&pricing.Plan,
			&pricing.Price,
			&pricing.BillingPeriod,
			&pricing.Currency,
			&features,
		)

		if err != nil {
			slog.Error("Unable to populate the pricing", "details", err.Error())
			return nil, err
		}
		err = json.Unmarshal(features, &pricing.Features)
		if err != nil {
			slog.Error("Unable to umarshal the features in pricing", "details", err.Error())
			return nil, err
		}
		result = append(result, pricing)
	}
	return result, nil
}

func (s *PostgresStore) GetPlanIds() (*[]int, error) {
	ids := []int{}
	rows, err := s.db.Query(`SELECT id from pricing;`)
	if err != nil {
		slog.Error("Unable to query for plan_id", "details", err.Error())
		return nil, err
	}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			slog.Error("Unable to get plan_id", "details", err.Error())
			return nil, err
		}
		ids = append(ids, id)
	}

	return &ids, nil
}
