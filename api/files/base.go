package api

import (
	"crypto/sha512"
	"fmt"
	"github.com/h2non/filetype"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileConfigs struct {
	FileType    string
	FileSubType []string
	TypeError   error
}

func ListDirectory(path string) ([]fs.FileInfo, error) {
	folderContent, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	return folderContent, nil
}

func DownloadFile(url string) (string, error) {
	content, err := http.Get(url)

	if err != nil {
		return "", err
	}

	chipper := sha512.New()
	chipper.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))

	fileName := chipper.Sum(nil)

	file, err := os.Create(fmt.Sprintf("./assets/%x", fileName))
	defer file.Close()

	_, err = io.Copy(file, content.Body)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/assets/%x", fileName), nil
}

func HandleFile(files []*multipart.FileHeader, conf FileConfigs) ([]string, error) {

	var output []string

	for _, handler := range files {

		file, err := handler.Open()

		if err != nil {
			return nil, err
		}

		defer file.Close()

		buff := make([]byte, handler.Size)
		_, err = file.Read(buff)

		if err != nil {
			return nil, err
		}

		kind, _ := filetype.Match(buff)

		if conf.FileType != "" && len(strings.Split(kind.MIME.Type, conf.FileType)) == 1 {
			return nil, conf.TypeError
		}

		if conf.FileSubType != nil {
			mismatch := true

			for _, fileType := range conf.FileSubType {

				if kind.MIME.Subtype == fileType {
					mismatch = false
				}

			}

			if mismatch {
				return nil, conf.TypeError
			}

		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}

		chipper := sha512.New()
		chipper.Write([]byte(fmt.Sprintf("%d_%d", time.Now().UnixNano(), filepath.Base(handler.Filename))))

		outputFileName := fmt.Sprintf("%x%s", chipper.Sum(nil), filepath.Ext(handler.Filename))
		serverFile, err := os.Create(fmt.Sprintf("./assets/%s", outputFileName))

		if err != nil {
			return nil, err
		}

		defer serverFile.Close()

		_, err = io.Copy(serverFile, file)

		if err != nil {
			return nil, err
		}

		output = append(output, outputFileName)
	}

	return output, nil
}
