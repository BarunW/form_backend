package data

import (
	"fmt"
	"log/slog"
)

func (s *PostgresStore) UpdateRequiredSetting(formId string, questionId int, setting bool) error {
	updateQuery := fmt.Sprintf("UPDATE created_form SET content=jsonb_set(content, '{%d,%d,setting,required}', '%v'::jsonb, false) WHERE id='%s';",
		questionId, questionId, setting, formId)

	_, err := s.db.Exec(updateQuery)
	if err != nil {
		slog.Error("Unable to update the required setting", "details", err.Error())
		return err
	}

	return nil
}

/*
   ===================
   Update Form Setting
   ==================
*/

func (s *PostgresStore) UpdateFormSettingProgresBar(formId string, setting bool) error {
	updateQuery := fmt.Sprintf("UPDATE created_form SET form_setting=jsonb_set(form_setting, '{progressbar}', '%v'::jsonb, false) WHERE id='%s';",
		setting, formId)

	_, err := s.db.Exec(updateQuery)
	if err != nil {
		slog.Error("Unable to update the form setting progress bar", "details", err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateFormSettingQNO(formId string, setting bool) error {
	updateQuery := fmt.Sprintf("UPDATE created_form SET form_setting=jsonb_set(form_setting, '{question_number}', '%v'::jsonb, false) WHERE id='%s';",
		setting, formId)

	_, err := s.db.Exec(updateQuery)
	if err != nil {
		slog.Error("Unable to update the form setting Qustion Number", "details", err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateFormSettingLettersOnAns(formId string, setting bool) error {
	updateQuery := fmt.Sprintf("UPDATE created_form SET form_setting=jsonb_set(form_setting, '{letters_on_answers}', '%v'::jsonb, false) WHERE id='%s';",
		setting, formId)

	_, err := s.db.Exec(updateQuery)
	if err != nil {
		slog.Error("Unable to update the form setting Qustion Number", "details", err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateFormSettingFreeNav(formId string, setting bool) error {
	updateQuery := fmt.Sprintf("UPDATE created_form SET form_setting=jsonb_set(form_setting, '{free_form_navigation}', '%v'::jsonb, false) WHERE id='%s';",
		setting, formId)

	_, err := s.db.Exec(updateQuery)
	if err != nil {
		slog.Error("Unable to update the form setting Qustion Number", "details", err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateSettingFormNavArrows(formId string, setting bool) error {
	updateQuery := fmt.Sprintf("UPDATE created_form SET form_setting=jsonb_set(form_setting, '{navigation_arrows}', '%v'::jsonb, false) WHERE id='%s';",
		setting, formId)

	_, err := s.db.Exec(updateQuery)
	if err != nil {
		slog.Error("Unable to update the form setting Qustion Number", "details", err.Error())
		return err
	}

	return nil
}
