package file

import (
	"fmt"
	"github.com/logansua/nfl_app/utils"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const (
	maxUploadSize = 2 * 1024 * 1024 // 2 mb
)

func UploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			utils.RenderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		// parse and validate file and post parameters
		file, _, err := r.FormFile("uploadFile")

		if err != nil {
			utils.RenderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)

		if err != nil {
			utils.RenderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		fileType := http.DetectContentType(fileBytes)

		switch fileType {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
		case "application/pdf":
			break
		default:
			utils.RenderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}

		fileName := utils.RandToken(12)
		fileEndings, err := mime.ExtensionsByType(fileType)

		if err != nil {
			utils.RenderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}

		var uploadPath = os.Getenv("APP_UPLOADS_PATH")
		if uploadPath == "" {
			utils.RenderError(w, "UPLOAD_PATH_IS_NOT_DEFINED", http.StatusInternalServerError)
			return
		}

		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

		newFile, err := os.Create(newPath)

		if err != nil {
			utils.RenderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}

		defer newFile.Close() // idempotent, okay to call twice

		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			utils.RenderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
