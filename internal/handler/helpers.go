package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]interface{}{
		"errors":     []string{msg},
		"error_type": errorType(status),
	})
}

func errorType(status int) string {
	switch status {
	case http.StatusNotFound:
		return "not_found"
	case http.StatusForbidden:
		return "invalid_access"
	case http.StatusUnprocessableEntity:
		return "invalid_parameters"
	case http.StatusTooManyRequests:
		return "rate_limit"
	default:
		return "invalid"
	}
}

func decodeBody(r *http.Request) (map[string]interface{}, error) {
	ct := r.Header.Get("Content-Type")
	data := make(map[string]interface{})

	if strings.Contains(ct, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return nil, err
		}
		return data, nil
	}

	// Form-encoded (default for many SDK clients)
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	for k, v := range r.Form {
		val := interface{}(v[0])
		if len(v) > 1 {
			val = v
		}
		// Parse bracket-notation keys like "post[raw]" into nested maps.
		if idx := strings.Index(k, "["); idx > 0 && strings.HasSuffix(k, "]") {
			outer := k[:idx]
			inner := k[idx+1 : len(k)-1]
			nested, ok := data[outer].(map[string]interface{})
			if !ok {
				nested = make(map[string]interface{})
				data[outer] = nested
			}
			nested[inner] = val
		} else {
			data[k] = val
		}
	}
	return data, nil
}

func pathParam(r *http.Request, name string) string {
	return r.PathValue(name)
}

func pathParamInt(r *http.Request, name string) (int, bool) {
	s := r.PathValue(name)
	s = strings.TrimSuffix(s, ".json")
	id, err := strconv.Atoi(s)
	return id, err == nil
}

func queryInt(r *http.Request, name string, def int) int {
	s := r.URL.Query().Get(name)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
