package data

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/sonal3323/form-poc/types"
	"github.com/sonal3323/form-poc/utils"
)

type Workspace struct {
	Name string `json:"name" validate:"required, max=30, min=1"`
	//Foreign key
	UserId int `json:"user_id" validate:"required, gte=0, lte=1000000000"`
}

// form meta data
type FormMetaData struct {
	WorspaceId    string `json:"workspace_id" validate:"required, max=26, min=26"`
	WorkspaceName string `json:"form_name" validate:"required"`
}

func (s *PostgresStore) createWorkspaceTable() error {
	createTableStmnt := `CREATE TABLE IF NOT EXISTS workspace(
            id VARCHAR (26) PRIMARY KEY,
            user_id INTEGER REFERENCES users (id),
            name VARCHAR(69)
        );`

	_, err := s.db.Exec(createTableStmnt)

	if err != nil {
		slog.Error("Unable to create the subscription table", "details", err.Error())
		return err
	}
	return nil
}

func (s *PostgresStore) createWorkspaceWithTx(tx *sql.Tx, userId int, workspaceName string) error {

	workspaceId, err := utils.GenerateULID()
	if err != nil {
		fmt.Println("Failed to generate ulid to create a new workspace")
		return err
	}

	_, err = tx.Exec(`INSERT INTO workspace (id, user_id, name) VALUES($1, $2, $3);`, workspaceId, userId, workspaceName)
	if err != nil {
		fmt.Println("Failed to insert into workspace")
		return err
	}

	return nil
}

func (s *PostgresStore) CreateWorkspace(ws types.Workspace) error {
	id, err := utils.GenerateULID()
	if err != nil {
		slog.Error("Error", "details", err.Error())
		return err
	}
	_, err = s.db.Exec(`INSERT INTO workspace (id, user_id, name) VALUES($1, $2, $3);`, id, ws.UserId, ws.Name)
	if err != nil {
		slog.Error("Failed to insert into workspace ", "details", err.Error())
		return err
	}
	return nil
}

func (s *PostgresStore) GetWorkspace(userId int) (*[]types.Workspace, error) {
	var workspaces []types.Workspace
	rows, err := s.db.Query("SELECT * FROM workspace WHERE user_id=$1", userId)
	if err != nil {
		slog.Error("Unable to query the workspace", "details", err.Error())
	}
	for rows.Next() {
		workspace := new(types.Workspace)
		err := rows.Scan(
			&workspace.WorkspaceId,
			&workspace.UserId,
			&workspace.Name,
		)

		if err != nil {
			slog.Error("Unable to scan the workspace rows", "details", err.Error())
			return nil, err
		}
		workspaces = append(workspaces, *workspace)
	}
	return &workspaces, nil
}

func (s *PostgresStore) GetWorkSpaceFormsData(userId int, workspaceId string) (*[]FormData, error) {
	return s.GetAllFormData(userId, workspaceId)
}

func isWorkspaceBelongToUser(tx *sql.Tx, workspaceId string, userId int) error {
	var isTrue bool
	err := tx.QueryRow(`SELECT EXISTS ( SELECT * FROM workspace WHERE id=$1 AND user_id=$2);`,
		workspaceId, userId).Scan(&isTrue)
	if err != nil {
		return err
	}

	if !isTrue {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostgresStore) DeleteWorkspace(userId int, workspaceId string) error {
	var (
		tx  *sql.Tx
		err error
	)

	tx, err = s.db.Begin()
	if err != nil {
		slog.Error("Failed to begin db transaction for deleting workspace", "DETAILS", err.Error())
		return err
	}
	err = isWorkspaceBelongToUser(tx, workspaceId, userId)
	if err != nil {
		return s.handleRollbackAndError("workspaceId and userId doesn't match", err, tx)
	}

	// After successfully updated
	// delete the workspace
	var deleteResult sql.Result
	deleteResult, err = tx.Exec("DELETE FROM workspace WHERE id=$1", workspaceId)
	if err != nil {
		return s.handleRollbackAndError("Failed to delete the workspace", err, tx)
	}

	err = s.checkRowsAffected(deleteResult)
	if err != nil {
		if err == sql.ErrNoRows {
			return s.handleRollbackAndError("Failed to delete the workspace 0 rows affected", sql.ErrNoRows, tx)
		}
		return s.handleRollbackAndError("Failed to get result after deleting the workspace", err, tx)
	}

	err = tx.Commit()
	if err != nil {
		return s.handleRollbackAndError("Failed to commit the transaction after deleting the workpsace", err, tx)
	}

	return nil
}
