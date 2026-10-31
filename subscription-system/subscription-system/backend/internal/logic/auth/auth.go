// Package auth handles registration, login and sessions. Sessions are a signed JWT in an httpOnly cookie
// holding only the user id; role, active flag and customer id are re-read on every request.
package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"brightloop/internal/apperr"
	"brightloop/internal/logic/platform"
	"brightloop/internal/logic/platform/audit"
)

const (
	CookieName = "bl_session"
	// Session lifetime is wall-clock time, not the business test clock.
	SessionTTL = 12 * time.Hour
)

// Secret signs session tokens. main sets it from config jwt.secret; signing refuses an empty secret.
var Secret []byte

var errBadLogin = &apperr.Error{Status: http.StatusUnauthorized, Code: "UNAUTHENTICATED", Message: "Email or password is incorrect."}

// dummyHash makes a login for an unknown email cost the same bcrypt time as a wrong password.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("not-a-real-password"), bcrypt.DefaultCost)

// CustomerID is the logged-in customer's id, 0 for staff and admins. Every portal query must filter by it.
func CustomerID(ctx context.Context) uint64 { return platform.ActorFrom(ctx).CustomerID }

func NormalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }

// HashPassword checks the password rule and returns the bcrypt hash, or a VALIDATION_ERROR on "password".
func HashPassword(p string) (string, error) {
	if msg := passwordError(p); msg != "" {
		return "", apperr.Validation(map[string]string{"password": msg})
	}
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h), err
}

type RegisterInput struct {
	Email, Password, Name, CompanyName, Phone string
}

// Register creates the customer login and its customer profile in one transaction.
func Register(ctx context.Context, in RegisterInput) (platform.Actor, error) {
	hash, err := HashPassword(in.Password)
	if err != nil {
		return platform.Actor{}, err
	}
	email := NormalizeEmail(in.Email)
	var a platform.Actor
	err = platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		r, err := tx.Exec("INSERT INTO users (email, password_hash, role) VALUES (?, ?, 'customer')", email, hash)
		if err != nil {
			return err
		}
		uid, _ := r.LastInsertId()
		r, err = tx.Exec("INSERT INTO customers (user_id, name, company_name, phone) VALUES (?, ?, NULLIF(?, ''), NULLIF(?, ''))",
			uid, strings.TrimSpace(in.Name), strings.TrimSpace(in.CompanyName), strings.TrimSpace(in.Phone))
		if err != nil {
			return err
		}
		cid, _ := r.LastInsertId()
		a = platform.Actor{UserID: uint64(uid), Role: "customer", CustomerID: uint64(cid)}
		actx := platform.WithActor(ctx, a) // a self-registration is changed by the new user
		if err = audit.Log(actx, tx, "users", a.UserID, "register", "", g.Map{"email": email}); err != nil {
			return err
		}
		return audit.Log(actx, tx, "customers", a.CustomerID, "create", "registration", g.Map{"name": in.Name})
	})
	if IsDuplicate(err) {
		return a, EmailTaken()
	}
	return a, err
}

func IsDuplicate(err error) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}

// EmailTaken is the duplicate-email error, attached to the email field.
func EmailTaken() *apperr.Error {
	e := apperr.Conflict("This email is already registered.")
	e.Fields = map[string]string{"email": e.Message}
	return e
}

// Login returns the user for a correct email and password. Wrong email and wrong password get the same error.
func Login(ctx context.Context, email, password string) (platform.Actor, error) {
	row, err := g.DB().GetOne(ctx, "SELECT id, password_hash, is_active FROM users WHERE email = ?", NormalizeEmail(email))
	if err != nil {
		return platform.Actor{}, err
	}
	hash := dummyHash
	if !row.IsEmpty() {
		hash = row["password_hash"].Bytes()
	}
	if bcrypt.CompareHashAndPassword(hash, []byte(password)) != nil || row.IsEmpty() {
		return platform.Actor{}, errBadLogin
	}
	if !row["is_active"].Bool() { // only revealed to someone who knows the password
		return platform.Actor{}, &apperr.Error{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: "This account is deactivated."}
	}
	return LoadActor(ctx, row["id"].Uint64())
}

// LoadActor reads the current role and customer id of an active user; an inactive or missing user is UNAUTHENTICATED.
func LoadActor(ctx context.Context, userID uint64) (platform.Actor, error) {
	row, err := g.DB().GetOne(ctx, `SELECT u.role, c.id AS customer_id FROM users u
		LEFT JOIN customers c ON c.user_id = u.id WHERE u.id = ? AND u.is_active`, userID)
	if err != nil {
		return platform.Actor{}, err
	}
	if row.IsEmpty() {
		return platform.Actor{}, apperr.Unauthenticated()
	}
	a := platform.Actor{UserID: userID, Role: row["role"].String()}
	if a.Role == "customer" {
		a.CustomerID = row["customer_id"].Uint64()
	}
	return a, nil
}

func IssueToken(userID uint64) (string, error) {
	if len(Secret) == 0 {
		return "", errors.New("auth.Secret is not set")
	}
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(SessionTTL)),
	}).SignedString(Secret)
}

// ParseToken returns the user id of a valid, unexpired HS256 token, or UNAUTHENTICATED.
func ParseToken(token string) (uint64, error) {
	var c jwt.RegisteredClaims
	if len(Secret) == 0 || token == "" {
		return 0, apperr.Unauthenticated()
	}
	_, err := jwt.ParseWithClaims(token, &c, func(*jwt.Token) (any, error) { return Secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		return 0, apperr.Unauthenticated()
	}
	id, err := strconv.ParseUint(c.Subject, 10, 64)
	if err != nil {
		return 0, apperr.Unauthenticated()
	}
	return id, nil
}
