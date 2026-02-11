package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// ============================================================================
// Extended Uploads
// ============================================================================

type ExtendedUploadsHandler struct {
	Store *store.Store
}

// GET /uploads/lookup-metadata
func (h *ExtendedUploadsHandler) LookupMetadata(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// GET /uploads/lookup-urls
func (h *ExtendedUploadsHandler) LookupURLs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// POST /uploads/generate-presigned-put
func (h *ExtendedUploadsHandler) GeneratePresignedPut(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"url":         "/uploads/default/original/placeholder.png",
		"key":         "uploads/default/original/placeholder.png",
		"unique_identifier": "upload_1",
	})
}

// POST /uploads/complete-external-upload
func (h *ExtendedUploadsHandler) CompleteExternalUpload(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":                1,
		"url":               "/uploads/default/original/placeholder.png",
		"original_filename": "placeholder.png",
		"filesize":          1024,
		"extension":         "png",
		"short_url":         "upload://placeholder",
		"short_path":        "/uploads/short-url/placeholder.png",
		"human_filesize":    "1 KB",
	})
}

// POST /uploads/create-multipart
func (h *ExtendedUploadsHandler) CreateMultipart(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"external_upload_identifier": "multipart_1",
		"key":                        "uploads/default/original/multipart.bin",
		"unique_identifier":          "multipart_1",
	})
}

// POST /uploads/batch-presign-multipart-parts
func (h *ExtendedUploadsHandler) BatchPresignParts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"presigned_urls": map[string]interface{}{},
	})
}

// POST /uploads/complete-multipart
func (h *ExtendedUploadsHandler) CompleteMultipart(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":                1,
		"url":               "/uploads/default/original/multipart.bin",
		"original_filename": "multipart.bin",
		"filesize":          10240,
		"extension":         "bin",
		"short_url":         "upload://multipart",
		"human_filesize":    "10 KB",
	})
}

// POST /uploads/abort-multipart
func (h *ExtendedUploadsHandler) AbortMultipart(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ============================================================================
// Extended Backups
// ============================================================================

type ExtendedBackupsHandler struct {
	Store *store.Store
}

// GET /admin/backups/{filename}/restore
func (h *ExtendedBackupsHandler) Restore(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /admin/backups/{filename}
func (h *ExtendedBackupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /admin/backups/status
func (h *ExtendedBackupsHandler) Status(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"is_operation_running": false,
		"allow_restore":       true,
	})
}

// GET /admin/backups/logs
func (h *ExtendedBackupsHandler) Logs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// DELETE /admin/backups/cancel
func (h *ExtendedBackupsHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /admin/backups/rollback
func (h *ExtendedBackupsHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/backups/readonly
func (h *ExtendedBackupsHandler) SetReadonly(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /admin/backups/is-backup-restore-running
func (h *ExtendedBackupsHandler) IsRunning(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"is_operation_running": false,
	})
}
