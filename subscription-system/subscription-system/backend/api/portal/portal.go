// Package portal is the request and response shapes of the customer portal. Every handler scopes by auth.CustomerID.
package portal

import (
	"github.com/gogf/gf/v2/frame/g"

	"brightloop/api/admin"
)

type GetProfileReq struct {
	g.Meta `path:"/portal/profile" method:"get" tags:"Portal" summary:"My profile"`
}

type UpdateProfileReq struct {
	g.Meta `path:"/portal/profile" method:"patch" tags:"Portal" summary:"Update my profile"`
	admin.CustomerFields
}
