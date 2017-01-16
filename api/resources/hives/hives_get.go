package hives

import (
	"github.com/muesli/beehive/bees"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *HiveResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *HiveResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *HiveResource) GetDoc() string {
	return "retrieve hives"
}

// GetParams returns the parameters supported by this API endpoint
func (r *HiveResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	// params = append(params, restful.QueryParameter("user_id", "id of a user").DataType("int64"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *HiveResource) GetByIDs(ctx smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := HiveResponse{}
	resp.Init(ctx)

	for _, id := range ids {
		hive := bees.GetBeeFactory(id)
		if hive == nil {
			r.NotFound(request, response)
			return
		}

		resp.AddHive(hive)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *HiveResource) Get(ctx smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	//	ctxapi := ctx.(*context.APIContext)
	hives := bees.GetBeeFactories()
	if len(hives) == 0 { // err != nil {
		r.NotFound(request, response)
		return
	}

	resp := HiveResponse{}
	resp.Init(ctx)

	for _, hive := range hives {
		resp.AddHive(hive)
	}

	resp.Send(response)
}
