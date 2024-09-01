package imageController

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type MediaFilesController struct{}

func NewMediaFilesController() *MediaFilesController {
	return &MediaFilesController{}
}

func (m *MediaFilesController) GetBasePath() (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("Unable to get the current working dir", "details", err.Error())
		return "", err
	}
	newPath, err := filepath.Abs(filepath.Join(cwd, "image"))
	if err != nil {
		slog.Error("Unable to get the absolute path", "details", err.Error())
		return "", err
	}
	return newPath, err
}

func (m *MediaFilesController) CreateFile(fileName string) (*os.File, error) {

	file, err := os.Create(fileName)
	if err != nil {
		slog.Error("Failed to create file", "details", err.Error())
		return nil, err
	}

	return file, err

}

func (m *MediaFilesController) createDir() (string, error) {

	// get the path of the imageVideo
	path, err := m.GetBasePath()
	if err != nil {
		return "", err
	}
	fmt.Println(path)
	//check if is already exist or not
	_, err = os.Stat(path)
	if os.IsExist(err) {
		return path, nil
	}

	// if dir is not exist create
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		slog.Error("Unable to create the imageVideo Dir", "details", err.Error())
	}

	return path, nil
}

func (m *MediaFilesController) DeleteFile(fileName string) error {

	if err := os.Remove(fileName); err != nil {
		slog.Error("Failed to remove the file", "details", err.Error())
		return err
	}

	return nil
}

func (m *MediaFilesController) CreateUserImageStorage(userId int) (string, error) {
	path, err := m.createDir()
	if err != nil {
		return "", err
	}

	newPath := filepath.Join(path, strconv.Itoa(userId))
	_, err = os.Stat(newPath)
	if !os.IsNotExist(err) {
		return newPath, nil
	}

	err = os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		slog.Error("Failed to created user image storage dir", "details", err.Error())
		return "", err
	}

	return newPath, nil
}
