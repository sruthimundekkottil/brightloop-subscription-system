package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/url"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"brightloop/internal/apperr"
	"brightloop/internal/logic/platform"
	"brightloop/internal/logic/platform/audit"
	"brightloop/internal/logic/platform/mail"
)

// ResetTTLMinutes is how long a reset link works. Expiry uses database wall time, not the business test clock.
const ResetTTLMinutes = 30

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// ForgotPassword emails a reset link if an active account has this email. The caller always answers 200,
// so the reply never reveals whether the account exists.
func ForgotPassword(ctx context.Context, email string) error {
	email = NormalizeEmail(email)
	uid, err := g.DB().GetValue(ctx, "SELECT id FROM users WHERE email = ? AND is_active", email)
	if err != nil || uid.IsEmpty() {
		return err
	}
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	link := mail.Setting(ctx, "APP_URL", "app.url", "http://localhost:5173") + "/reset-password?token=" + url.QueryEscape(token)
	return platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		r, err := tx.Exec("INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES (?, ?, NOW(3) + INTERVAL ? MINUTE)",
			uid.Uint64(), hashToken(token), ResetTTLMinutes)
		if err != nil {
			return err
		}
		id, _ := r.LastInsertId()
		if err = audit.Log(ctx, tx, "password_reset_tokens", uint64(id), "request", "", g.Map{"user_id": uid.Uint64()}); err != nil {
			return err
		}
		return mail.Enqueue(ctx, tx, "", mail.Email{To: email, Subject: "Reset your Brightloop password",
			Body: "Someone asked to reset the password for this account.\n\nSet a new password (the link works for 30 minutes, once):\n" +
				link + "\n\nIf this wasn't you, ignore this email; your password stays the same."})
	})
}

func invalidToken() *apperr.Error {
	return apperr.Validation(map[string]string{"token": "This reset link is invalid or has expired. Request a new one."})
}

// ResetPassword sets a new password for a valid, unused, unexpired token, then voids every open token of that user.
func ResetPassword(ctx context.Context, token, password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	return platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE: two requests with the same token serialise here, and only the first sees used_at IS NULL.
		row, err := tx.GetOne(`SELECT id, user_id FROM password_reset_tokens
			WHERE token_hash = ? AND used_at IS NULL AND expires_at > NOW(3) FOR UPDATE`, hashToken(token))
		if err != nil {
			return err
		}
		if row.IsEmpty() {
			return invalidToken()
		}
		uid := row["user_id"].Uint64()
		if _, err = tx.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hash, uid); err != nil {
			return err
		}
		if _, err = tx.Exec("UPDATE password_reset_tokens SET used_at = NOW(3) WHERE user_id = ? AND used_at IS NULL", uid); err != nil {
			return err
		}
		return audit.Log(platform.WithActor(ctx, platform.Actor{UserID: uid}), tx, "users", uid, "reset-password", "", nil)
	})
}
