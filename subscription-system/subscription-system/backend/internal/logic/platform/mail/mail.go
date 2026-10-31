// Package mail sends email through the job queue: Enqueue in the business transaction, Work delivers over SMTP
// (MailHog in development) with retries.
package mail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"brightloop/internal/logic/platform"
)

const Queue = "email"

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"` // plain text
}

// Enqueue queues e in tx, so it is sent only if tx commits. dedupeKey "" allows repeats.
func Enqueue(ctx context.Context, tx gdb.TX, dedupeKey string, e Email) error {
	return platform.Enqueue(tx.Ctx(ctx), Queue, dedupeKey, e, time.Now()) // wall clock: the queue runs on NOW(3)
}

// Send delivers one email; tests replace it.
var Send = sendSMTP

// Work claims a batch of due emails and sends them; failures are retried with backoff by platform.Fail.
func Work(ctx context.Context) {
	host, _ := os.Hostname()
	owner := fmt.Sprintf("%s-%d", host, os.Getpid())
	jobs, err := platform.Claim(ctx, Queue, owner, 20)
	if err != nil {
		g.Log().Error(ctx, "mail: claim:", err)
		return
	}
	for _, j := range jobs {
		var e Email
		err := json.Unmarshal([]byte(j.Payload), &e)
		if err == nil {
			err = Send(ctx, e)
		}
		if err != nil {
			g.Log().Warningf(ctx, "mail: job %d attempt %d: %v", j.ID, j.Attempts, err)
			err = platform.Fail(ctx, j.ID, owner, err)
		} else {
			err = platform.Done(ctx, j.ID, owner)
		}
		if err != nil {
			g.Log().Error(ctx, "mail: job", j.ID, err)
		}
	}
}

// Setting reads env first (docker-compose), then config, then def.
func Setting(ctx context.Context, env, key, def string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	if g.Cfg().Available(ctx) { // MustGet panics when there is no config file (tests)
		if v := g.Cfg().MustGet(ctx, key).String(); v != "" {
			return v
		}
	}
	return def
}

func sendSMTP(ctx context.Context, e Email) error {
	if strings.ContainsAny(e.To+e.Subject, "\r\n") {
		return errors.New("header injection: newline in To or Subject") // never sendable, but retried harmlessly until failed
	}
	addr := Setting(ctx, "SMTP_ADDR", "smtp.addr", "127.0.0.1:1025")
	from := Setting(ctx, "SMTP_FROM", "smtp.from", "no-reply@brightloop.test")
	msg := "From: Brightloop <" + from + ">\r\nTo: " + e.To + "\r\nSubject: " + e.Subject +
		"\r\nDate: " + time.Now().UTC().Format(time.RFC1123Z) +
		"\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" +
		strings.ReplaceAll(e.Body, "\n", "\r\n")
	return smtp.SendMail(addr, nil, from, []string{e.To}, []byte(msg)) // shortcut: no SMTP auth/TLS, add both for a real mail server
}
