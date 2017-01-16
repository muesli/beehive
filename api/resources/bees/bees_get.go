package bees

import (
	"github.com/muesli/beehive/bees"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *BeeResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *BeeResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *BeeResource) GetDoc() string {
	return "retrieve bees"
}

// GetParams returns the parameters supported by this API endpoint
func (r *BeeResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	// params = append(params, restful.QueryParameter("user_id", "id of a user").DataType("int64"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *BeeResource) GetByIDs(ctx smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := BeeResponse{}
	resp.Init(ctx)

	for _, id := range ids {
		bee := bees.GetBee(id)
		if bee == nil {
			r.NotFound(request, response)
			return
		}

		resp.AddBee(bee)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *BeeResource) Get(ctx smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	//	ctxapi := ctx.(*context.APIContext)
	bees := bees.GetBees()
	if len(bees) == 0 { // err != nil {
		r.NotFound(request, response)
		return
	}

	resp := BeeResponse{}
	resp.Init(ctx)

	for _, bee := range bees {
		resp.AddBee(bee)
	}

	resp.Send(response)
}
