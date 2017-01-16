package context

import (
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// APIContext is polly's central context
type APIContext struct {
}

// NewAPIContext returns a new polly context
func (context *APIContext) NewAPIContext() smolder.APIContext {
	ctx := &APIContext{}
	return ctx
}

// Authentication parses the request for an access-/authtoken and returns the matching user
func (context *APIContext) Authentication(request *restful.Request) (interface{}, error) {
	t := request.QueryParameter("accesstoken")
	if len(t) == 0 {
		t = request.HeaderParameter("authorization")
		if strings.Index(t, " ") > 0 {
			t = strings.TrimSpace(strings.Split(t, " ")[1])
		}
	}

	return nil, nil // context.GetUserByAccessToken(t)
}

// LogSummary logs out the current context stats
func (context *APIContext) LogSummary() {
}
