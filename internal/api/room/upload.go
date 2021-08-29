package room

import (
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"demodesk/neko/internal/utils"
)

const (
	// Maximum upload of 32 MB files.
	MAX_UPLOAD_SIZE = 32 << 20
)

func (h *RoomHandler) uploadDrop(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		utils.HttpBadRequest(w, "failed to parse multipart form")
		return
	}

	//nolint
	defer r.MultipartForm.RemoveAll()

	X, err := strconv.Atoi(r.FormValue("x"))
	if err != nil {
		utils.HttpBadRequest(w, "no X coordinate received")
		return
	}

	Y, err := strconv.Atoi(r.FormValue("y"))
	if err != nil {
		utils.HttpBadRequest(w, "no Y coordinate received")
		return
	}

	req_files := r.MultipartForm.File["files"]
	if len(req_files) == 0 {
		utils.HttpBadRequest(w, "no files received")
		return
	}

	dir, err := os.MkdirTemp("", "neko-drop-*")
	if err != nil {
		utils.HttpInternalServerError(w, err)
		return
	}

	files := []string{}
	for _, req_file := range req_files {
		path := path.Join(dir, req_file.Filename)

		srcFile, err := req_file.Open()
		if err != nil {
			utils.HttpInternalServerError(w, err)
			return
		}

		defer srcFile.Close()

		dstFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			utils.HttpInternalServerError(w, err)
			return
		}

		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			utils.HttpInternalServerError(w, err)
			return
		}

		files = append(files, path)
	}

	if !h.desktop.DropFiles(X, Y, files) {
		utils.HttpInternalServerError(w, "unable to drop files")
		return
	}

	utils.HttpSuccess(w)
}

func (h *RoomHandler) uploadDialogPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		utils.HttpBadRequest(w, "failed to parse multipart form")
		return
	}

	//nolint
	defer r.MultipartForm.RemoveAll()

	if !h.desktop.IsFileChooserDialogOpened() {
		utils.HttpBadRequest(w, "open file chooser dialog first")
		return
	}

	req_files := r.MultipartForm.File["files"]
	if len(req_files) == 0 {
		utils.HttpBadRequest(w, "no files received")
		return
	}

	dir, err := os.MkdirTemp("", "neko-dialog-*")
	if err != nil {
		utils.HttpInternalServerError(w, err)
		return
	}

	for _, req_file := range req_files {
		path := path.Join(dir, req_file.Filename)

		srcFile, err := req_file.Open()
		if err != nil {
			utils.HttpInternalServerError(w, err)
			return
		}

		defer srcFile.Close()

		dstFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			utils.HttpInternalServerError(w, err)
			return
		}

		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			utils.HttpInternalServerError(w, err)
			return
		}
	}

	if err := h.desktop.HandleFileChooserDialog(dir); err != nil {
		utils.HttpInternalServerError(w, "unable to handle file chooser dialog")
		return
	}

	utils.HttpSuccess(w)
}

func (h *RoomHandler) uploadDialogClose(w http.ResponseWriter, r *http.Request) {
	if !h.desktop.IsFileChooserDialogOpened() {
		utils.HttpBadRequest(w, "file chooser dialog is not open")
		return
	}

	h.desktop.CloseFileChooserDialog()
	utils.HttpSuccess(w)
}
