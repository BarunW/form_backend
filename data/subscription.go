package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/sonal3323/form-poc/types"
)

type userSubscriptionDeatils struct {
	SubsId        int                   `json:"subs_id"`
	UserId        int                   `json:"user_id"`
	Email         string                `json:"email"`
	Username      string                `json:"username"`
	PlanId        int                   `json:"plan_id"`
	RespLimitId   types.ResponseLimitId `json:"resp_limit_id"`
	StartDate     time.Time             `json:"start_date"`
	EndDate       time.Time             `json:"end_date"`
	Feature       types.Features        `json:"plan_features"`
	BillingPeriod types.BillingMode     `json:"billing_period"`
	Plan          types.PlanType        `json:"plan"`
	Status        types.PLAN_STATUS     `json:"status"`
}

func (s *PostgresStore) createSubscriptionTable() error {
	createTableStmnt := `CREATE TABLE IF NOT EXISTS subscription (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users (id),
            plan_id INTEGER REFERENCES pricing (id),
            resp_limitId VARCHAR REFERENCES responseaddon(id),
            start_date TIMESTAMP WITH TIME ZONE,
            end_date TIMESTAMP WITH TIME ZONE,
            billing_period bp, 
            status plan_status
        );`

	_, err := s.db.Exec(createTableStmnt)

	if err != nil {
		slog.Error("Unable to create the subscription table", "details", err.Error())
		return err
	}
	return nil
}

/*
================
At data level we control the sensitive data like
start Date , EndDate and status
=================
*/

func (s *PostgresStore) Subscribe(subscriber types.Subscriber) error {

	// start day of the subscription
	startDate := time.Now()
	startDate.Format(time.RFC3339)

	// y, m , d
	var endDate time.Time
	if subscriber.PlanId >= int(BASIC_YEARLY) && subscriber.PlanId <= int(BUSINESS_YEARLY) {
		endDate = startDate.AddDate(1, 0, 0)
	} else if subscriber.PlanId >= int(FREE_PLAN) && subscriber.PlanId <= int(BUSINESS_MONTHLY) {
		endDate = startDate.AddDate(0, 1, 0)
	}

	_, err := s.db.Exec(`
                    INSERT INTO subscription(
                    plan_id, user_id, resp_limitId, start_date, end_date, 
                    billing_period, status)
                    VALUES($1, $2, $3, $4, $5, $6, $7)`,
		subscriber.PlanId,
		subscriber.UserId,
		subscriber.ResponseLimitId,
		startDate,
		endDate,
		subscriber.BillingPeriod,
		types.Active,
	)
	if err != nil {
		slog.Error("failed to insert subscriber's data to db", "details", err.Error())
		return err
	}
	return nil
}

func (s *PostgresStore) UserPlanData(uid int) (*userSubscriptionDeatils, error) {
	rows, err := s.db.Query(`
    SELECT subscription.*, 
    users.email, users.username, 
    pricing.features, 
    responseaddon.plan
    FROM subscription 
    JOIN users ON subscription.user_id= users.id 
    JOIN pricing ON subscription.plan_id = pricing.id 
    JOIN responseaddon ON subscription.resp_limitid = responseaddon.id
    WHERE subscription.user_id = $1;`,
		uid,
	)

	if err != nil {
		slog.Error("Unable to get the user plan data", "details", err.Error())
		return nil, err
	}

	for rows.Next() {
		return s.scanIntoSubscription(rows)
	}

	// return the free plan if the user doesn't subscribe any paid plan
	price, err := s.getPricing(1121)
	if err != nil {
		return nil, err
	}
	freeSubs := userSubscriptionDeatils{
		Status:  types.FreePlan,
		Plan:    types.Free,
		Feature: price.Features,
	}
	return &freeSubs, nil
}

func (s *PostgresStore) scanIntoSubscription(rows *sql.Rows) (*userSubscriptionDeatils, error) {
	var feature json.RawMessage
	u := userSubscriptionDeatils{
		Feature: types.Features{},
	}
	err := rows.Scan(
		&u.SubsId,
		&u.UserId,
		&u.PlanId,
		&u.RespLimitId,
		&u.StartDate,
		&u.EndDate,
		&u.BillingPeriod,
		&u.Status,
		&u.Email,
		&u.Username,
		&feature,
		&u.Plan,
	)

	fmt.Printf("%+v", u)

	if err != nil {
		slog.Error("failed to scan into user subscription", "details", err.Error())
		return nil, err
	}

	err = json.Unmarshal(feature, &u.Feature)
	if err != nil {
		slog.Error("Failed to unmarshal the feature in subscription", "details", err.Error())
		return nil, err
	}
	return &u, nil
}
