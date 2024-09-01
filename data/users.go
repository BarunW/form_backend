package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"github.com/sonal3323/form-poc/types"
	"github.com/sonal3323/form-poc/utils"
	"golang.org/x/crypto/bcrypt"
)

type userAccount struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AccountId string `json:"account_id"` 
    Tokens    tokens
}

func (s *PostgresStore) createAccountTable() error {
	createTableQuery := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
        account_id VARCHAR(11),
        username VARCHAR(50),
        email VARCHAR(50),
        password BYTEA 
	)`
	tx, err := s.db.Begin()
    if err != nil{
        slog.Error("Failed to begin db transaction before creating users table", "details", err.Error())
        return err
    }
    _, err = tx.Exec(createTableQuery)
    if err != nil{
       return s.handleRollbackAndError("Error While creating users table", err, tx)
    }

    createTokenTable := `CREATE TABLE IF NOT EXISTS jwt_tokens_(
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        token TEXT[] 
    )`
    _, err = tx.Exec(createTokenTable)
    if err != nil{
       return s.handleRollbackAndError("Error While creating jwt_tokens table", err, tx)
    }

    if err := tx.Commit(); err != nil{ 
       return s.handleRollbackAndError("Failed to commit the tx users and jwt_tokens table", err, tx)
    }
	return nil 
}


func (s *PostgresStore) CreateAccount(acc *types.Account) (*userAccount, error){
	if exist, err := s.isUserExist(acc.Email); err != nil {
		fmt.Println("Error in userExist", err)
		return nil, err
	} else if exist == true {
		return nil, fmt.Errorf("409")
	}

	acc_id, err := utils.GenerateAccountId()
	if err != nil {
		return nil, err
	}

	var tx *sql.Tx
	tx, err = s.db.Begin()
	if err != nil {
        rErr := s.handleRollbackAndError("Failed to begin db transaction for create account", err, nil)
        return nil, rErr
	}
	_, err = tx.Exec("INSERT INTO users (account_id, username, email, password)VALUES ($1, $2, $3, $4)",
		acc_id,
		acc.Username,
		acc.Email,
		acc.Password,
	)
	if err != nil {
        rErr := s.handleRollbackAndError("Failed to insert users", err, tx)
        return nil, rErr
	}

	var userId int
	// query the user id
	err = tx.QueryRow("SELECT id from users WHERE email=$1", acc.Email).Scan(&userId)
	if err != nil {
		
        rErr := s.handleRollbackAndError("Unable to get userid from users", err, tx)
        return nil, rErr
	}

	// create a default workspace for the create user
	if err := s.createWorkspaceWithTx(tx, userId, "My Workspace"); err != nil {
        rErr := s.handleRollbackAndError("Failed to created workspace", err, tx)
        return nil, rErr
	}

	// create the responses data too
	// response limit and responses collected
	// 10 response limit
	rl, rc := DefaultReponseAddon()
	rm := ResponsesModel{
		UserId:            userId,
		AccountId:         acc_id,
		ResponseLimit:     rl,
		ResponseCollected: rc,
	}

	err = InsertResponsesWithTx(tx, rm)
	if err != nil {
        rErr := s.handleRollbackAndError("Unable to insert response", err, tx) 
		return nil, rErr
	}

	if err := tx.Commit(); err != nil {
        rErr:=s.handleRollbackAndError("Failed to commit the transaction while creating user", err, tx)
        return nil, rErr 
    }

    tk, err := s.authorize(userId, acc.Email) 
    if err != nil{
        slog.Error("Unable to authorize the user", "details", err)
        return nil, err
    }
    return &userAccount{
        ID: userId,
        Email: acc.Email,
        Username: acc.Username,
        Tokens: tk,
    }, nil
}

func(s *PostgresStore) authorize(userId int, email string) (tokens, error){

    access_token , err := newAccessToken(email)
    if err != nil{
        slog.Error("Unable to get new access token", "details", err)
        return tokens{}, err
    }

    var refresh_token string 
    refresh_token, err = s.setUserRefreshToken(userId, email)
    if err != nil{
        slog.Error("Unable to get new refresh token", "details", err)
        return tokens{}, err
    }

    return tokens{AccessToken: access_token, RefreshToken: refresh_token}, nil
    
}


func (s *PostgresStore) GetAccount(credentials *types.Login) (*userAccount, error) {
	var hasehdPassword []byte
	err := s.db.QueryRow("SELECT password FROM users WHERE email = $1", credentials.Email).Scan(&hasehdPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, bcrypt.ErrMismatchedHashAndPassword
		}
		slog.Error("unable to query", "details", err.Error())
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(hasehdPassword, []byte(credentials.Password)); err != nil {
		return nil, bcrypt.ErrMismatchedHashAndPassword
	}

	var ua userAccount
	err = s.db.QueryRow("SELECT id, account_id, username, email FROM users WHERE email = $1", credentials.Email).Scan(&ua.ID, &ua.AccountId, &ua.Username, &ua.Email)
	if err != nil {
		slog.Error("unable to query", "details", err.Error())
		return nil, err
	}
    
    var tk tokens
    tk, err = s.authorize(ua.ID, ua.Email)
    if err != nil{
        slog.Error("unable to authorize the token", "details", err.Error())
        return nil, err
    }
    ua.Tokens = tk 
	return &ua, nil
}

func (s *PostgresStore) isUserExist(email string) (bool, error) {
	var user_email string
	err := s.db.QueryRow("SELECT email from users WHERE email=$1", email).Scan(&user_email)
	if err != nil {
		fmt.Println(sql.ErrNoRows.Error(), user_email)
		if err.Error() == sql.ErrNoRows.Error() && user_email == "" {
			return false, nil
		}
	}

	return true, err
}

func (s *PostgresStore) GetAccountIdThroughtFormId(formId string) (string, error) {
	var account_id string
	err := s.db.QueryRow(`select users.account_id  from 
                    created_form LEFT JOIN users ON created_form.user_id=users.id WHERE created_form.id=$1;`, formId).Scan(&account_id)
	return account_id, err
}
