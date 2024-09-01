package imageController

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDir(t *testing.T) {
	mfc := NewMediaFilesController()
	path, err := mfc.CreateUserImageStorage(1)
	if err != nil {
		t.Fatal(err)
	}

	newPath := filepath.Join(path, "new.png")
	_, err = mfc.CreateFile(newPath)
	_, err = os.Stat(newPath)
	if os.IsNotExist(err) {
		t.Fatal(err)
	}

	err = mfc.DeleteFile(newPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(newPath)
	if os.IsExist(err) {
		t.Fatal(err)
	}
}
