package data

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/lib/pq"
	"github.com/sonal3323/form-poc/types"
)

func (s *PostgresStore) createTemplPlanTable() error {
	if exist := s.isTableExist("template_plan"); exist {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return s.handleRollbackAndError("Failed to begin db transaction for template_plan table", err, nil)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS template_plan(     
            plan_id INTEGER,
            question_template_id VARCHAR(26),
            PRIMARY KEY (plan_id, question_template_id),
            FOREIGN KEY (question_template_id) REFERENCES question_templates(id),
            FOREIGN KEY (plan_id) REFERENCES pricing(id)
        );`)
	if err != nil {
		return s.handleRollbackAndError("Failed to create template_plan table", err, tx)
	}

	err = insertTemplateAndPlanId(s, tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return s.handleRollbackAndError("Failed to commit templateAndPlanId transaction", err, tx)
	}
	return nil

}

func insertTemplateAndPlanId(s *PostgresStore, tx *sql.Tx) error {

	stmt, err := tx.Prepare(pq.CopyIn("template_plan", "plan_id", "question_template_id"))
	if err != nil {
		return s.handleRollbackAndError("Unable to prepare statement for templatePlan table", err, tx)
	}

	defer stmt.Close()

	planIds, err := s.GetPlanIds()
	if err != nil {
		return s.handleRollbackAndError("Unable to prepare statement for templPlan and failed to get plan id", err, tx)
	}

	var f func(string, int) error
	f = func(templId string, planId int) error {
		_, err := stmt.Exec(planId, templId)
		return s.handleRollbackAndError("Unable to execute statment for templates_planId", err, tx)
	}

	templs, err := s.GetQuestionTemplates()
	if err != nil {
		return s.handleRollbackAndError("Unable to prepare statement for templatePlan table", err, tx)
	}
	for _, id := range *planIds {
		for _, temp := range *templs {
			if id == int(FREE_PLAN) && temp.IsPaid == !Paid {
				if err := f(temp.Id, id); err != nil {
					return err
				}
				continue
			}

			if err := f(temp.Id, id); err != nil {
				return err
			}
		}
	}
	if _, err := stmt.Exec(); err != nil {
		s.handleRollbackAndError("Unable to finalize statement for templateAndPlanId", err, tx)
	}
	return nil
}

type templPlanId struct {
	PlanId             types.PlanId_Type
	QuestionTemplateId string
}

func (s *PostgresStore) GetTemplatePlan(planId types.PlanId_Type) (*[]templPlanId, error) {
	rows, err := s.db.Query("SELECT * FROM template_plan WHERE plan_id=$1;", planId)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to query template_plan for plan_id %v", planId), "details", err.Error())
		return nil, err
	}

	var templPlanIds []templPlanId

	for rows.Next() {
		templPlanId_ := templPlanId{}

		err := rows.Scan(
			&templPlanId_.PlanId,
			&templPlanId_.QuestionTemplateId,
		)
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to scan rows for template_plan for plan_id %v", planId), "details", err.Error())
			return nil, err
		}

		templPlanIds = append(templPlanIds, templPlanId_)
	}

	return &templPlanIds, nil
}
