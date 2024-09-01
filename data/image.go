package data

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sonal3323/form-poc/imageController"
	"github.com/sonal3323/form-poc/types"
)

func (s *PostgresStore) UploadImage(userId int, formId string, questionId int, fileName string, r io.Reader) error {
	mfc := imageController.NewMediaFilesController()

	userImgStoragePath, err := mfc.CreateUserImageStorage(userId)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	fileName = fmt.Sprintf("%s_%d", fileName, now)
	fullPath := filepath.Join(userImgStoragePath, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer file.Close()

	nn, err := io.Copy(file, r)
	if err != nil || nn == 0 {
		slog.Error("Unable to copy the content from the given file", "details", err.Error())
		mfc.DeleteFile(fullPath)
		return err
	}

	// create the url
	url := fmt.Sprintf("http://localhost:8080/image/%d/%s", userId, fileName)
	fmt.Println(url)
	updateType := fmt.Sprintf(`'{%d, %d, setting, image_or_videoSettings, type}', '"IMG"'::jsonb, false) WHERE id='%s';`, questionId, questionId, formId)
	updateURL := fmt.Sprintf(`UPDATE created_form SET content=jsonb_set(jsonb_set(content, '{%d, %d, setting, image_or_videoSettings, url}', to_jsonb('%s'::text) , false),%s`,
		questionId, questionId, url, updateType)
	fmt.Println(updateURL)
	_, err = s.db.Exec(updateURL)
	if err != nil {
		slog.Error("Unable to update the url", "details", err.Error())
		mfc.DeleteFile(fullPath)
		return err
	}

	return nil

}

func DeleteTheFile(url string) error {
	fmt.Println(url)
	// spilt the url into chunks and get the last two element as is userId and fileName
	substr := strings.Split(url, "/")
	length := len(substr)

	mfc := imageController.NewMediaFilesController()

	// get the base path where the image is store
	basePath, err := mfc.GetBasePath()
	if err != nil {
		return err
	}
	fmt.Printf("%+v", substr)
	// create new path

	fullPath := filepath.Join(basePath, substr[length-2], substr[length-1])
	fmt.Println(fullPath)

	// delete the file
	return mfc.DeleteFile(fullPath)
}

//
// change question id to imageId
//func (s *PostgresStore) DeleteImageFile(formId string, imageId string ) error {
//
//	// begin the db transaction
//	tx, err := s.db.Begin()
//	if err != nil {
//		slog.Error("Unable to begi db transaction for remove of image", "details", err.Error())
//		return err
//	}
//
//	var url string
//	getURLQUery := fmt.Sprintf("SELECT content->%d->'%d'->'setting'->'image_or_videoSettings'->'url' from created_form WHERE id='%s'", questionId, questionId, formId)
//	err = tx.QueryRow(getURLQUery).Scan(&url)
//
//	if err != nil {
//		return s.handleRollbackAndError("Unable to query the url", err, tx)
//	}
//    fmt.Println(url)
//
//    // update the url to "" and type to none
//    updateType := fmt.Sprintf(`'{%d, %d, setting, image_or_videoSettings, type}', 'null'::jsonb) WHERE id='%s';`, questionId, questionId, formId)
//    updateURL  := fmt.Sprintf(`UPDATE created_form SET content=jsonb_set(jsonb_set(content, '{%d, %d, setting, image_or_videoSettings, url}', 'null'::jsonb , false),%s`,
//                    questionId, questionId, updateType)
//
//	fmt.Println(updateURL)
//
//	_, err = tx.Exec(updateURL)
//	if err != nil {
//		return s.handleRollbackAndError("Unable to remove the url", err, tx)
//	}
//
//
//	err = DeleteTheFile(url)
//	if err != nil {
//		return s.handleRollbackAndError("Unable to delete the file", err, tx)
//	}
//
//	return tx.Commit()
//}
//

func (s *PostgresStore) UpdateImageOrVideoURL(formId string, questionId int, url string, t types.MediaType) error {

	if url == "" {
		url = "null"
		t = "null"
	}

	updateType := fmt.Sprintf(`'{%d, %d, setting, image_or_videoSettings, type}', '"%s"'::jsonb) WHERE id='%s';`, questionId, questionId, string(t), formId)

	updateURL := fmt.Sprintf(`UPDATE created_form SET content=jsonb_set(jsonb_set(content, '{%d, %d, setting, image_or_videoSettings, url}', '"%s"'::jsonb, false),%s`,
		questionId, questionId, url, updateType)

	_, err := s.db.Exec(updateURL)
	fmt.Println(updateURL)

	if err != nil {
		slog.Error("Unable to add the url", "details", err.Error())
		return err
	}

	return nil
}
