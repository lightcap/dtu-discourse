package handler

import (
	"net/http"
	"strings"
)

// AdminUsersSubRouter handles all GET /admin/users/{...} routes to avoid
// Go ServeMux conflicts between /admin/users/{id}/ip_info and
// /admin/users/list/{type}.
type AdminUsersSubRouter struct {
	Users    *UsersHandler
	Extended *ExtendedAdminHandler
}

// ServeGET handles GET /admin/users/{rest...}
func (d *AdminUsersSubRouter) ServeGET(w http.ResponseWriter, r *http.Request) {
	rest := r.PathValue("rest")
	parts := strings.Split(rest, "/")

	if len(parts) == 0 {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	first := parts[0]

	// /admin/users/list/{type}
	if first == "list" && len(parts) >= 2 {
		r.SetPathValue("type", parts[1])
		d.Users.ListUsers(w, r)
		return
	}

	// /admin/users/approve-bulk (no trailing path)
	if first == "approve-bulk" {
		// This is PUT only, handled elsewhere
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// /admin/users/destroy-bulk (no trailing path)
	if first == "destroy-bulk" {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// /admin/users/{id} or /admin/users/{id}/ip_info
	r.SetPathValue("id", strings.TrimSuffix(first, ".json"))

	if len(parts) == 1 {
		// /admin/users/{id}
		d.Users.GetUserByID(w, r)
		return
	}

	second := parts[1]
	switch second {
	case "ip_info":
		d.Extended.IPInfo(w, r)
		return
	}

	writeError(w, http.StatusNotFound, "not found")
}
