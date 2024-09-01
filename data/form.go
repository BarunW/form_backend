package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"

	"github.com/sonal3323/form-poc/types"
	"github.com/sonal3323/form-poc/utils"
)

type FormData struct {
	FormId       string    `json:"form_id"`
	Title        string    `json:"title"`
	CreatedOn    time.Time `json:"created_on"`
	Questions    int       `json:"questions"`
	Responses    int       `json:"responses"`
	Completion   float32   `json:"completion"`
	UpdatedOn    time.Time `json:"updated_on"`
	Integrations []string  `json:"integrations"`
    
}

type Form struct {
	Id         string          `json:"id"`
	UserId     int             `json:"user_id"`
	WokspaceId string          `json:"workspace_id"`
	Title      string          `json:"title"`
	Content    []byte          `json:"content"`
	Setting    []byte          `json:"form_setting"`
}

func (s *PostgresStore) createdFormTable() error {

	if exist := s.isTableExist("created_form"); exist {
		return nil
	}
	stmnt := `CREATE TABLE IF NOT EXISTS created_form(
            id VARCHAR(26) PRIMARY KEY,
            user_id INTEGER REFERENCES users(id) ,
            workspace_id VARCHAR REFERENCES workspace(id) ON DELETE CASCADE,
            title VARCHAR(100),
            content JSONB,
            form_setting JSONB,
            created_on TIMESTAMP WITH TIME ZONE, 
            questions INTEGER,
            responses INTEGER,
            completion REAL,
            updated_on TIMESTAMP WITH TIME ZONE,
            integrations text[] 
    );`
	_, err := s.db.Exec(stmnt)
	if err != nil {
		return s.handleRollbackAndError("Unable to create table created_form ", err, nil)
	}

	return nil

}

func newFormData() FormData {

	createdOn := time.Now()

	return FormData{
		CreatedOn:    createdOn,
		Questions:    1,
		Responses:    0,
		Completion:   0.0,
		UpdatedOn:    createdOn,
		Integrations: []string{""},
	}
}

type CreatedQuestionTemplate map[int]QuestionTemplate

type ContentsType struct {
	Contents []CreatedQuestionTemplate `json:"contents"`
}

// this function
// assigned the default question number 
// title and other
func createContent(qt QuestionTemplate) ([]byte, error) {
	// default questions key is 0
	var cqt CreatedQuestionTemplate = map[int]QuestionTemplate{
		0: qt,
	}

	createdQuestionTemplates := ContentsType{
		Contents: []CreatedQuestionTemplate{cqt},
	}

	contentByte, err := json.Marshal(&createdQuestionTemplates.Contents)
	if err != nil {
		slog.Error("Unable to Marshal the content", "details", err.Error())
		return nil, err
	}

	return contentByte, nil
}

func isFormBelongsToUser(tx *sql.Tx, formId string, userId int) error {
	var isTrue bool
	err := tx.QueryRow(`SELECT EXISTS ( SELECT * FROM created_form WHERE id=$1 AND user_id=$2);`,
		formId, userId).Scan(&isTrue)
	if err != nil {
		return err
	}

	if !isTrue {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostgresStore) CreateForm(cf types.CreateForm) error {
	// this id is for multiple choice question template use as default when user create
	var err error
	var qt *QuestionTemplate
	qt, err = s.GetQuestionTemplate(cf.TemplateId)

	if err != nil {
		slog.Error("Unable to get question template", "details", err.Error())
		return err

	}
    
    // array of question templates including
    // settings for the question
	var contentByte []byte
	contentByte, err = createContent(*qt)
	if err != nil {
		return err
	}

    // settings for the form
	formSetting := &types.FormSetting_Type{}
	formSettingByte, err := json.Marshal(formSetting)
	if err != nil {
		slog.Error("Unable to Marshal the form settings", "details", err.Error())
		return err
	}

	// generate ulid for form id
	ulid, err := utils.GenerateULID()
	if err != nil {
		slog.Error("Unable to generate ULID", "details", err.Error())
		return err
	}

	// get new form data
	fd := newFormData()
	formatTime := fd.CreatedOn.Format(time.RFC3339)

	var tx *sql.Tx
	if tx, err = s.db.Begin(); err != nil {
		return s.handleRollbackAndError("Failed to begin transaction for created_form", err, nil)
	}

	err = isWorkspaceBelongToUser(tx, cf.WorkspaceId, cf.UserId)
	if err != nil {
		return s.handleRollbackAndError("Workspace and userID doesn't match", err, tx)
	}
	// Insert the from in created_form table
	_, err = tx.Exec(`INSERT INTO created_form(
            id, user_id, workspace_id, title, content, form_setting, 
            created_on,
            questions,
            responses,
            completion,
            updated_on,
            integrations
        ) 
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`,
		ulid,
		cf.UserId,
		cf.WorkspaceId,
		cf.FormName,
		contentByte,
		formSettingByte,
		formatTime,
		fd.Questions,
		fd.Responses,
		fd.Completion,
		formatTime,
		pq.Array(fd.Integrations),
	)

	if err != nil {
		return s.handleRollbackAndError("Unable to insert the data on created_form table ", err, tx)
	}

	err = insertOnFormRespData(tx, ulid, cf.WorkspaceId)
	if err != nil {
		return s.handleRollbackAndError("Failed to insert on form_response_data table", err, tx)
	}

	if err = tx.Commit(); err != nil {
		return s.handleRollbackAndError("Unable to commit transaction for created_form table ", err, tx)
	}

	return nil
}

func (s *PostgresStore) AddQuestion(aqt types.AddQuestionTempl) error {
	var length int
	err := s.db.QueryRow(`SELECT questions FROM created_form WHERE id = $1`, aqt.Id).Scan(&length)
	if err != nil {
		slog.Error("Unable to query created_form for content length", "details", err.Error())
		return err
	}

	qt, err := s.GetQuestionTemplate(aqt.TemplateId)
	if err != nil {
		slog.Error("Unable to get question template", "details", err.Error())
		return err
	}

	var cqt CreatedQuestionTemplate = map[int]QuestionTemplate{
		// its map ( 0 : {data} )
		length: *qt,
	}
	newContentsByte, err := json.Marshal(cqt)
	if err != nil {
		slog.Error("Unable to Marshal newContents to json", "details", err.Error())
	}

	updateTime := time.Now()
	updateTime.Format(time.RFC3339)

	_, err = s.db.Exec("UPDATE created_form SET content = content || $1::jsonb, questions=questions+$2, updated_on=$3 WHERE id=$4",
		newContentsByte, 1, updateTime, aqt.Id)
	if err != nil {
		slog.Error("Unable to update the created_form table", "details", err.Error())
		return err
	}

	return nil
}

/*
fm = FormMetaData [ Id, UserId]
*/
func (s *PostgresStore) GetForm(fm types.FormMetaData) (*Form, error) {
	var form Form
	err := s.db.QueryRow("SELECT id, user_id, workspace_id, title, content, form_setting  FROM created_form WHERE id=$1",
		fm.Id).Scan(&form.Id, &form.UserId, &form.WokspaceId, &form.Title, &form.Content, &form.Setting)

	if err != nil {
		slog.Error("Unable to query the created_form table for the user", "details", err.Error())
		return nil, err
	}

	if form.Id == "" || form.UserId == 0 {
		return nil, fmt.Errorf("%s", "404")
	}

	err = updateTotalStart(fm.Id, s)
	if err != nil {
		return nil, err
	}

	return &form, nil
}


// FormData
func (s *PostgresStore) GetFormData(fm types.FormMetaData) (*FormData, error) {
	fd := FormData{}
    var (
        totalSubmissions int 
        totalStart float32
        err error
    )
	err = s.db.QueryRow(`SELECT 
            form_response_data.submissions,
            form_response_data.total_start,
            id, 
            created_form.title,
            created_form.created_on, 
            created_form.questions, 
            created_form.responses,
            created_form.completion,
            created_form.updated_on,
            FROM form_response_data 
            LEFT JOIN created_form ON form_response_data.form_id=created_form.id
            WHERE created_form.id=$1 and created_form.user_id=$2`,
		fm.Id, fm.UserId).Scan(totalSubmissions, totalStart, &fd.FormId, &fd.Title, &fd.CreatedOn, &fd.Questions, &fd.Responses, &fd.Completion, &fd.UpdatedOn)

	if err != nil {
		slog.Error("Unable to query the response data created_form table for the user", "details", err.Error())
		return nil, err
	}
        
    fd.Completion = (float32(totalSubmissions)/totalStart)  * 100
    fd.Responses = totalSubmissions

	if fd.FormId == "" {
		return nil, fmt.Errorf("%s", "404")
	}

	return &fd, nil

}

func (s *PostgresStore) GetAllFormData(userId int, workspaceId string) (*[]FormData, error) {

	var formDatas []FormData
	rows, err := s.db.Query(`SELECT 
            form_response_data.total_submissions,
            form_response_data.total_start,
            created_form.id, 
            created_form.title,
            created_form.created_on, 
            created_form.questions, 
            created_form.responses,
            created_form.completion,
            created_form.updated_on
            FROM form_response_data 
            JOIN created_form ON form_response_data.workspace_id=created_form.workspace_id
            WHERE created_form.user_id=$1 and created_form.workspace_id=$2`,
		userId, workspaceId)
    
	if err != nil {
		slog.Error("Unable to query the created_form table for the user", "details", err.Error())
		return nil, err
	}

	for rows.Next() {
		fd := FormData{}
        var totalSubmissions int 
        var totalStart float32
		err := rows.Scan(&totalSubmissions, &totalStart, &fd.FormId, &fd.Title, &fd.CreatedOn, &fd.Questions, &fd.Responses, &fd.Completion, &fd.UpdatedOn)
		if err != nil {
			slog.Error("Unable to scan the rows for geeting form data from form table", "details", err.Error())
			return nil, err
		}
        fd.Completion = (float32(totalSubmissions)/totalStart)  * 100
        fd.Responses = totalSubmissions
		formDatas = append(formDatas, fd)
	}


	return &formDatas, nil

}

func (s *PostgresStore) DeleteQuestion(formId string, questionId int) error {

	var content []byte
	err := s.db.QueryRow(`SELECT content FROM created_form WHERE id = $1`, formId).Scan(&content)
	if err != nil {
		slog.Error("Unable to query created_form for content length", "details", err.Error())
		return err
	}

	// deleting question
	// requires updating question id
	// as it is manually incremented

	updatedContent, err := updateIdAndDeleteQuestion(content, questionId)
	if err != nil {
		slog.Error("Unable to delete the question", "details", err.Error())
		return err
	}

	updateTime := time.Now()
	updateTime.Format(time.RFC3339)

	_, err = s.db.Exec("UPDATE created_form SET content=$1, questions=questions-$2, updated_on=$3 WHERE id=$4", updatedContent, 1, updateTime, formId)
	if err != nil {
		slog.Error("Unable to delete the question", "details", err.Error())
		return err
	}
	return err
}

/* Delete Question Helper Function */
func updateIdAndDeleteQuestion(data []byte, qId int) ([]byte, error) {

	createdQuestionTemplates := ContentsType{
		Contents: []CreatedQuestionTemplate{},
	}

	err := json.Unmarshal(data, &createdQuestionTemplates.Contents)
	if err != nil {
		slog.Error("Unable to UnMarshal the content", "details", err.Error())
		return nil, err
	}

	if qId >= len(createdQuestionTemplates.Contents) || len(createdQuestionTemplates.Contents) == 1 {
		return nil, fmt.Errorf("Invalid Id")
	}

grepId:
	for idx, question := range createdQuestionTemplates.Contents {
		if _, ok := question[qId]; ok {
			// to iterate and delete the questions and change the ids to the next question
			for idx < len(createdQuestionTemplates.Contents)-1 {
				createdQuestionTemplates.Contents[idx][idx] = createdQuestionTemplates.Contents[idx+1][idx+1]
				idx++
			}
			// delete the last element
			fmt.Println(idx)
			createdQuestionTemplates.Contents = createdQuestionTemplates.Contents[:len(createdQuestionTemplates.Contents)-1]
			break grepId
		}
	}

	var newData []byte
	newData, err = json.Marshal(&createdQuestionTemplates.Contents)
	if err != nil {
		slog.Error("Unable to Marshal the content", "details", err.Error())
		return newData, err
	}

	return newData, nil
}

func (s *PostgresStore) ReorderQuestion(formId string, fromPositionId, toPositionId int) error {
	var content []byte
	var noOfQuestions int
	err := s.db.QueryRow(`SELECT content, questions FROM created_form WHERE id = $1`, formId).Scan(&content, &noOfQuestions)
	if err != nil {
		slog.Error("Unable to query created_form for content and questions ", "details", err.Error())
		return err
	}

	// we store the questions in an array and num of questions = len(array) | but the position starts from 0 to questions -1
	if fromPositionId > noOfQuestions-1 || fromPositionId < 0 || toPositionId > noOfQuestions-1 || toPositionId < 0 {
		return fmt.Errorf("Invalid position")
	}

	var reOrderedContent []byte
	if fromPositionId > toPositionId {
		reOrderedContent, err = pushQuestionToTop(content, fromPositionId, toPositionId)
		if err != nil {
			return err
		}

	} else if fromPositionId < toPositionId {
		reOrderedContent, err = pushQuestionToBottom(content, fromPositionId, toPositionId)
		if err != nil {
			return err
		}

	}
	_, err = s.db.Exec("UPDATE created_form SET content=$1 WHERE id = $2", reOrderedContent, formId)
	if err != nil {
		slog.Error("Unable to update created_form for reordering of content ", "details", err.Error())
		return err
	}
	return nil
}

func jsonToContentType(data []byte) ([]CreatedQuestionTemplate, error) {

	cqt := ContentsType{
		Contents: []CreatedQuestionTemplate{},
	}

	err := json.Unmarshal(data, &cqt.Contents)
	if err != nil {
		slog.Error("Unable to Marshal the content", "details", err.Error())
		return nil, err
	}

	return cqt.Contents, nil
}

func pushQuestionToTop(data []byte, fromPos, toPos int) ([]byte, error) {
	content, err := jsonToContentType(data)
	if err != nil {
		return nil, err
	}

	length := fromPos - toPos
	i := 0
	temp := content[fromPos][fromPos]
	for i < length {
		content[fromPos-i][fromPos-i] = content[fromPos-i-1][fromPos-i-1]
		i++
	}
	content[toPos][toPos] = temp

	return json.Marshal(content)
}

func pushQuestionToBottom(data []byte, fromPos, toPos int) ([]byte, error) {

	content, err := jsonToContentType(data)
	if err != nil {
		return nil, err
	}

	length := toPos - fromPos
	i := 0
	temp := content[fromPos][fromPos]
	for i < length {
		content[fromPos+i][fromPos+i] = content[fromPos+i+1][fromPos+i+1]
		i++
	}
	content[toPos][toPos] = temp

	return json.Marshal(content)
}

func (s *PostgresStore) DeleteForm(formId string, userId int) error {
	// first upadte the foreign keys
	// user_id and workspace id
	// and the form
	var (
		tx  *sql.Tx
		err error
	)

	tx, err = s.db.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction before deleting form", "DETAILS", err.Error())
		return err
	}

	err = isFormBelongsToUser(tx, formId, userId)
	if err != nil {
		slog.Error("This form doesn't belong to the user", "DETAILS", err.Error())
		return err
	}

	// delete the constraint
	var deleteResult sql.Result
	deleteResult, err = tx.Exec("DELETE FROM created_form WHERE id=$1", formId)
	if err != nil {
		return s.handleRollbackAndError("Failed to execute the delete statement", err, tx)
	}
	// check the result
	err = s.checkRowsAffected(deleteResult)
	if err != nil {
		if err == sql.ErrNoRows {
			return s.handleRollbackAndError("0 rows affected after updating create_form", err, tx)
		}
		return s.handleRollbackAndError("Unable to get the result after updating created_form", err, tx)
	}

	err = tx.Commit()
	if err != nil {
		return s.handleRollbackAndError("Failed to commit transaction after deleting the form", err, tx)
	}

	return nil
}

/*
==================
Update Questions And Description
================
*/
func (s *PostgresStore) UpdateQuestion(questionId int, text, formId string) error {

	updateStatment := fmt.Sprintf(`UPDATE created_form
    SET content = jsonb_set(content, '{%d,%d,data,question}', to_jsonb('%s'::text), false) 
    WHERE id='%s';`, questionId, questionId, text, formId)
	fmt.Println(updateStatment)
	_, err := s.db.Exec(updateStatment)

	if err != nil {
		slog.Error("Unable to update the question", "details", err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateDescription(questionId int, text, formId string) error {

	updateStatment := fmt.Sprintf(`UPDATE created_form
    SET content = jsonb_set(content, '{%d,%d,data,description}', to_jsonb('%s'::text), false) 
    WHERE id='%s';`, questionId, questionId, text, formId)
	fmt.Println(updateStatment)
	_, err := s.db.Exec(updateStatment)
	if err != nil {
		slog.Error("Unable to update the question", "details", err.Error())
		return err
	}

	return nil
}

func (s *PostgresStore) GetFormLink(formId string) (string, error) {
	account_id, err := s.GetAccountIdThroughtFormId(formId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}

	return createdFormLink(account_id, formId), nil

}

func createdFormLink(subDomain, param string) string {
	const protocol string = "http"
	const domain string = "localhost"

	return fmt.Sprintf("%s://%s.%s:8080/to/%s", protocol, subDomain, domain, param)
}
