package handler

import (
	"net/http"
	"time"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// SessionHandler covers session/auth related endpoints: login, logout,
// forgot-password, passkeys, 2FA, user-api-keys, and session info.
type SessionHandler struct {
	Store *store.Store
}

// ---- Session ----

// POST /session
func (h *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, _ := decodeBody(r)
	login, _ := body["login"].(string)
	if login == "" {
		login, _ = body["username"].(string)
	}
	u := h.Store.GetUserByUsername(login)
	if u == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"error":  "Invalid credentials",
			"reason": "invalid_credentials",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"name":     u.Name,
		},
	})
}

// DELETE /session/{username}
func (h *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /session/current.json
func (h *SessionHandler) Current(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"current_user": map[string]interface{}{
			"id":              1,
			"username":        "admin",
			"name":            "Admin User",
			"avatar_template": "/letter_avatar_proxy/v4/letter/a/bbce88/{size}.png",
			"admin":           true,
			"moderator":       true,
			"trust_level":     4,
			"can_create_topic": true,
			"can_review":       true,
			"unread_notifications": 0,
			"unread_high_priority_notifications": 0,
		},
	})
}

// ---- Forgot Password ----

// POST /session/forgot_password
func (h *SessionHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_found": true,
	})
}

// ---- User API Keys ----

// POST /user-api-key/new
func (h *SessionHandler) NewUserAPIKey(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"key": "user_api_key_placeholder_" + time.Now().Format("20060102150405"),
	})
}

// POST /user-api-key
func (h *SessionHandler) CreateUserAPIKey(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"key": "user_api_key_" + time.Now().Format("20060102150405"),
	})
}

// POST /user-api-key/revoke
func (h *SessionHandler) RevokeUserAPIKey(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /user-api-key/undo-revoke
func (h *SessionHandler) UndoRevokeUserAPIKey(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Passkeys ----

// GET /session/passkey/challenge
func (h *SessionHandler) PasskeyChallenge(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"challenge": "passkey_challenge_placeholder",
	})
}

// POST /session/passkey/auth
func (h *SessionHandler) PasskeyAuth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":       1,
			"username": "admin",
		},
	})
}

// ---- 2FA ----

// POST /session/2fa
func (h *SessionHandler) TwoFactorAuth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"ok":   true,
		"user": map[string]interface{}{"id": 1, "username": "admin"},
	})
}

// GET /session/2fa.json
func (h *SessionHandler) TwoFactorStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"totp_enabled":          false,
		"security_key_enabled":  false,
		"backup_enabled":        false,
		"allowed_methods":       []string{"totp", "security_key", "backup_code"},
	})
}

// ---- HP (Honeypot) ----

// GET /session/hp
func (h *SessionHandler) Honeypot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"value":     "honeypot_value",
		"challenge": "honeypot_challenge",
	})
}

// ---- Email Login ----

// POST /session/email-login/{token}
func (h *SessionHandler) EmailLogin(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{"id": 1, "username": "admin"},
	})
}

// GET /session/email-login/{token}
func (h *SessionHandler) EmailLoginInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"can_login":  true,
		"token_email": "admin@example.com",
	})
}

// ---- SSO (client side) ----

// GET /session/sso
func (h *SessionHandler) SSORedirect(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"sso_url": "/session/sso_login",
	})
}

// GET /session/sso_provider
func (h *SessionHandler) SSOProvider(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"sso_url": "/session/sso_login",
	})
}
