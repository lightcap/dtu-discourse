package handler

import (
	"net/http"
	"strconv"
	"strings"
)

// PostSubRouter handles all GET /posts/{...} routes to avoid Go ServeMux
// ambiguity between patterns like /posts/{id}/revisions/latest and
// /posts/by_number/{topic_id}/{post_number}.
type PostSubRouter struct {
	Posts    *PostsHandler
	Extended *ExtendedPostsHandler
}

// ServeGET handles GET /posts/{rest...}
func (d *PostSubRouter) ServeGET(w http.ResponseWriter, r *http.Request) {
	rest := r.PathValue("rest")
	parts := strings.Split(rest, "/")

	if len(parts) == 0 {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	first := parts[0]

	// /posts/by_number/{topic_id}/{post_number}
	if first == "by_number" && len(parts) >= 3 {
		r.SetPathValue("topic_id", parts[1])
		r.SetPathValue("post_number", parts[2])
		d.Extended.ByNumber(w, r)
		return
	}

	// /posts/{id}/revisions/latest
	idStr := strings.TrimSuffix(first, ".json")
	_, err := strconv.Atoi(idStr)

	// If first segment is not numeric, it's a username-based route
	if err != nil {
		r.SetPathValue("username", first)
		if len(parts) >= 2 {
			switch parts[1] {
			case "deleted":
				d.Extended.Deleted(w, r)
				return
			case "pending":
				d.Extended.Pending(w, r)
				return
			}
		}
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	r.SetPathValue("id", idStr)

	if len(parts) == 1 {
		// /posts/{id}
		d.Posts.Get(w, r)
		return
	}

	second := parts[1]

	switch second {
	case "revisions":
		if len(parts) == 3 {
			third := parts[2]
			if third == "latest" {
				d.Extended.LatestRevision(w, r)
				return
			}
			// /posts/{id}/revisions/{revision}
			r.SetPathValue("revision", third)
			d.Extended.Revision(w, r)
			return
		}
		writeError(w, http.StatusNotFound, "not found")
		return
	case "reply-history":
		d.Extended.ReplyHistory(w, r)
		return
	case "reply-ids":
		d.Extended.ReplyIDs(w, r)
		return
	case "cooked":
		d.Extended.Cooked(w, r)
		return
	case "raw":
		d.Extended.Raw(w, r)
		return
	case "replies":
		d.Extended.Replies(w, r)
		return
	}

	writeError(w, http.StatusNotFound, "not found")
}
