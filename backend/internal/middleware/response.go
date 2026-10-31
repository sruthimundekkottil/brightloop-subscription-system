// Package middleware holds the HTTP middleware every route uses.
package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"

	"brightloop/api/common"
	"brightloop/internal/apperr"
	"brightloop/internal/logic/platform"
)

// Response writes { "data", "meta" } on success and { "error": {...} } on failure.
// Handlers that write the body themselves (PDF downloads) are left alone.
func Response(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		e := toAppError(err)
		if e.Status >= 500 {
			g.Log().Error(r.Context(), err) // full detail only in logs
		}
		r.Response.ClearBuffer() // drops GoFrame's panic text
		r.Response.WriteHeader(e.Status)
		r.Response.WriteJson(g.Map{"error": e})
		return
	}
	if r.Response.BufferLength() > 0 {
		return
	}
	res := r.GetHandlerResponse()
	if p, ok := res.(interface{ Paged() (any, common.PageMeta) }); ok {
		items, meta := p.Paged()
		r.Response.WriteJson(g.Map{"data": items, "meta": meta})
		return
	}
	r.Response.WriteJson(g.Map{"data": res})
}

func toAppError(err error) *apperr.Error {
	var ae *apperr.Error
	if errors.As(err, &ae) {
		return ae
	}
	var ve gvalid.Error
	if errors.As(err, &ve) {
		fields := map[string]string{}
		for field, rules := range ve.Maps() {
			for _, msg := range rules {
				fields[jsonName(field)] = msg.Error()
				break
			}
		}
		return apperr.Validation(fields)
	}
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		switch {
		case me.Number == 1062:
			return apperr.Conflict("This record already exists.")
		case me.Number == 3819:
			return &apperr.Error{Status: http.StatusUnprocessableEntity, Code: "VALIDATION_ERROR",
				Message: "A value is outside the allowed range."}
		case platform.IsRetryable(err): // platform.Tx already retried once
			return apperr.Conflict("Someone else is changing this record. Please try again.")
		}
	}
	return apperr.Internal()
}

// jsonName turns GoFrame's Go field name into our JSON name: request fields are tagged lowerCamelCase.
func jsonName(field string) string {
	if field == "" {
		return field
	}
	return strings.ToLower(field[:1]) + field[1:]
}
