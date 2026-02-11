package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// ExtendedPMHandler covers PM subtypes not in the core private_messages handler:
// unread, archive, new, warnings, group PMs, PM tags.
type ExtendedPMHandler struct {
	Store *store.Store
}

func (h *ExtendedPMHandler) emptyTopicList(w http.ResponseWriter) {
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         []model.Topic{},
		},
	})
}

// GET /topics/private-messages-unread/{username}
func (h *ExtendedPMHandler) Unread(w http.ResponseWriter, r *http.Request) {
	h.emptyTopicList(w)
}

// GET /topics/private-messages-archive/{username}
func (h *ExtendedPMHandler) Archive(w http.ResponseWriter, r *http.Request) {
	h.emptyTopicList(w)
}

// GET /topics/private-messages-new/{username}
func (h *ExtendedPMHandler) New(w http.ResponseWriter, r *http.Request) {
	h.emptyTopicList(w)
}

// GET /topics/private-messages-warnings/{username}
func (h *ExtendedPMHandler) Warnings(w http.ResponseWriter, r *http.Request) {
	h.emptyTopicList(w)
}

// GET /topics/private-messages-group/{username}/{group_name}
func (h *ExtendedPMHandler) GroupPMs(w http.ResponseWriter, r *http.Request) {
	h.emptyTopicList(w)
}

// GET /topics/private-messages-tags/{username}/{tag}
func (h *ExtendedPMHandler) PMTags(w http.ResponseWriter, r *http.Request) {
	h.emptyTopicList(w)
}

// PUT /topics/private-messages/{username}/archive
func (h *ExtendedPMHandler) MoveToArchive(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /topics/private-messages/{username}/move-to-inbox
func (h *ExtendedPMHandler) MoveToInbox(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}
