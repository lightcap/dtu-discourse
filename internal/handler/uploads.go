package handler

import (
	"net/http"
	"path/filepath"

	"github.com/lightcap/dtu-discourse/internal/store"
)

type UploadsHandler struct {
	Store *store.Store
}

// POST /uploads
func (h *UploadsHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB max

	filename := "upload.bin"
	var filesize int

	file, header, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		filename = header.Filename
		filesize = int(header.Size)
	} else if url := r.FormValue("url"); url != "" {
		filename = filepath.Base(url)
		filesize = 0
	}

	ext := filepath.Ext(filename)
	if len(ext) > 0 {
		ext = ext[1:] // remove leading dot
	}

	upload := h.Store.CreateUpload(filename, ext, filesize)
	writeJSON(w, http.StatusOK, upload)
}
