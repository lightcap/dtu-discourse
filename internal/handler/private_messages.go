package handler

import (
	"net/http"
	"strings"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type PrivateMessagesHandler struct {
	Store *store.Store
}

// GET /topics/private-messages/{username}.json
func (h *PrivateMessagesHandler) Inbox(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	username = strings.TrimSuffix(username, ".json")
	topics := h.Store.GetPrivateMessages(username)
	writeJSON(w, http.StatusOK, model.PrivateMessageListResponse{
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}

// GET /topics/private-messages-sent/{username}.json
func (h *PrivateMessagesHandler) Sent(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	username = strings.TrimSuffix(username, ".json")
	topics := h.Store.GetSentPrivateMessages(username)
	writeJSON(w, http.StatusOK, model.PrivateMessageListResponse{
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}
