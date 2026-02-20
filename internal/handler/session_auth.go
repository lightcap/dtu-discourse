package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
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

// ---- SSO (DiscourseConnect) ----

// GET /session/sso
func (h *SessionHandler) SSORedirect(w http.ResponseWriter, r *http.Request) {
	if h.Store.SSOSecret == "" || h.Store.SSOCallbackURL == "" {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"sso_url": "/session/sso_login",
		})
		return
	}

	nonce := h.Store.CreateSSONonce()
	payload := fmt.Sprintf("nonce=%s", nonce)
	b64 := base64.StdEncoding.EncodeToString([]byte(payload))

	mac := hmac.New(sha256.New, []byte(h.Store.SSOSecret))
	mac.Write([]byte(b64))
	sig := hex.EncodeToString(mac.Sum(nil))

	dest := fmt.Sprintf("%s?sso=%s&sig=%s",
		h.Store.SSOCallbackURL,
		url.QueryEscape(b64),
		url.QueryEscape(sig),
	)
	http.Redirect(w, r, dest, http.StatusFound)
}

// GET /session/sso_provider
func (h *SessionHandler) SSOProvider(w http.ResponseWriter, r *http.Request) {
	h.SSORedirect(w, r)
}

// GET /session/sso_login — handles the return leg of DiscourseConnect SSO.
func (h *SessionHandler) SSOLogin(w http.ResponseWriter, r *http.Request) {
	if h.Store.SSOSecret == "" {
		writeError(w, http.StatusBadRequest, "SSO not configured")
		return
	}

	ssoPayload := r.URL.Query().Get("sso")
	sig := r.URL.Query().Get("sig")
	if ssoPayload == "" || sig == "" {
		writeError(w, http.StatusBadRequest, "missing sso or sig parameter")
		return
	}

	// Validate HMAC-SHA256 signature
	mac := hmac.New(sha256.New, []byte(h.Store.SSOSecret))
	mac.Write([]byte(ssoPayload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		writeError(w, http.StatusForbidden, "invalid signature")
		return
	}

	// Base64-decode and parse params
	decoded, err := base64.StdEncoding.DecodeString(ssoPayload)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid sso payload encoding")
		return
	}
	params, err := url.ParseQuery(string(decoded))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid sso payload")
		return
	}

	nonce := params.Get("nonce")
	if !h.Store.ValidateSSONonce(nonce) {
		writeError(w, http.StatusForbidden, "invalid or expired nonce")
		return
	}

	externalID := params.Get("external_id")
	email := params.Get("email")
	username := params.Get("username")
	name := params.Get("name")
	if externalID == "" || email == "" {
		writeError(w, http.StatusBadRequest, "external_id and email are required")
		return
	}
	if username == "" {
		username = externalID
	}

	_, err = h.Store.SyncSSO(externalID, email, username, name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
