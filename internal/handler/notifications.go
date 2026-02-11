package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/middleware"
	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type NotificationsHandler struct {
	Store *store.Store
}

// GET /notifications.json
func (h *NotificationsHandler) List(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsername(r)
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeJSON(w, http.StatusOK, model.NotificationListResponse{
			Notifications: []model.Notification{},
		})
		return
	}
	notifs := h.Store.GetNotifications(u.ID)
	writeJSON(w, http.StatusOK, model.NotificationListResponse{
		Notifications:            notifs,
		TotalRowsNotifications:   len(notifs),
		SeenNotificationID:       0,
	})
}
