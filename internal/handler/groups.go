package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type GroupsHandler struct {
	Store *store.Store
}

// GET /groups.json
func (h *GroupsHandler) List(w http.ResponseWriter, r *http.Request) {
	groups := h.Store.ListGroups()
	writeJSON(w, http.StatusOK, model.GroupListResponse{
		Groups:          groups,
		TotalRowsGroups: len(groups),
	})
}

// GET /groups/{group_name}.json
func (h *GroupsHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "group_name")
	name = strings.TrimSuffix(name, ".json")
	g := h.Store.GetGroup(name)
	if g == nil {
		writeError(w, http.StatusNotFound, "group not found")
		return
	}
	writeJSON(w, http.StatusOK, model.GroupResponse{Group: *g})
}

// POST /admin/groups
func (h *GroupsHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// Handle nested "group" key
	groupData := body
	if nested, ok := body["group"].(map[string]interface{}); ok {
		groupData = nested
	}
	name, _ := groupData["name"].(string)
	if name == "" {
		writeError(w, http.StatusUnprocessableEntity, "name is required")
		return
	}
	g, err := h.Store.CreateGroup(name, groupData)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"basic_group": g})
}

// PUT /groups/{group_id}
func (h *GroupsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "group_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if nested, ok := body["group"].(map[string]interface{}); ok {
		body = nested
	}
	g, err := h.Store.UpdateGroup(id, body)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.GroupResponse{Group: *g})
}

// DELETE /admin/groups/{group_id}.json
func (h *GroupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "group_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	if err := h.Store.DeleteGroup(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/groups/{group_id}/members.json
func (h *GroupsHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "group_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	body, _ := decodeBody(r)
	userIDs := extractIntSlice(body, "usernames", h.Store)
	if err := h.Store.AddGroupMembers(id, userIDs); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /admin/groups/{group_id}/members.json
func (h *GroupsHandler) RemoveMembers(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "group_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	body, _ := decodeBody(r)
	userIDs := extractIntSlice(body, "usernames", h.Store)
	if err := h.Store.RemoveGroupMembers(id, userIDs); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/groups/{group_id}/owners.json
func (h *GroupsHandler) AddOwners(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "group_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	body, _ := decodeBody(r)
	userIDs := extractIntSlice(body, "usernames", h.Store)
	if err := h.Store.AddGroupOwners(id, userIDs); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /admin/groups/{group_id}/owners.json
func (h *GroupsHandler) RemoveOwners(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "group_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	body, _ := decodeBody(r)
	userIDs := extractIntSlice(body, "usernames", h.Store)
	if err := h.Store.RemoveGroupOwners(id, userIDs); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /groups/{group_name}/members.json
func (h *GroupsHandler) Members(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "group_name")
	name = strings.TrimSuffix(name, ".json")
	g := h.Store.GetGroup(name)
	if g == nil {
		writeError(w, http.StatusNotFound, "group not found")
		return
	}
	members, owners := h.Store.GetGroupMembers(g.ID)
	writeJSON(w, http.StatusOK, model.GroupMembersResponse{
		Members: members,
		Owners:  owners,
		Meta: model.GroupMembersMeta{
			Total:  len(members),
			Limit:  50,
			Offset: 0,
		},
	})
}

// POST /groups/{group}/notifications
func (h *GroupsHandler) SetNotificationLevel(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func extractIntSlice(body map[string]interface{}, key string, s *store.Store) []int {
	var ids []int

	// Try "user_ids" first
	if v, ok := body["user_ids"].([]interface{}); ok {
		for _, item := range v {
			switch val := item.(type) {
			case float64:
				ids = append(ids, int(val))
			case string:
				if id, err := strconv.Atoi(val); err == nil {
					ids = append(ids, id)
				}
			}
		}
		return ids
	}

	// Try usernames (comma-separated string)
	if v, ok := body[key].(string); ok {
		for _, name := range strings.Split(v, ",") {
			name = strings.TrimSpace(name)
			if u := s.GetUserByUsername(name); u != nil {
				ids = append(ids, u.ID)
			}
		}
	}

	return ids
}
