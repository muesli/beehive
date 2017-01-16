package bees

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/muesli/beehive/bees"
	"github.com/muesli/smolder"
)

// BeePostStruct holds all values of an incoming POST request
type BeePostStruct struct {
	Bee struct {
		Description string          `json:"description"`
		Active      bool            `json:"active"`
		Options     bees.BeeOptions `json:"options"`
	} `json:"bee"`
}

// BeePutStruct holds all values of an incoming PUT request
type BeePutStruct struct {
	BeePostStruct
}

// PutAuthRequired returns true because all requests need authentication
func (r *BeeResource) PutAuthRequired() bool {
	return false
}

// PutDoc returns the description of this API endpoint
func (r *BeeResource) PutDoc() string {
	return "update an existing bee"
}

// PutParams returns the parameters supported by this API endpoint
func (r *BeeResource) PutParams() []*restful.Parameter {
	return nil
}

// Put processes an incoming PUT (update) request
func (r *BeeResource) Put(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	resp := BeeResponse{}
	resp.Init(context)

	pps := BeePutStruct{}
	err := request.ReadEntity(&pps)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse PUT data",
			"BeeResource PUT"))
		return
	}

	id := request.PathParameter("bee-id")
	bee := bees.GetBee(id)
	if bee == nil {
		r.NotFound(request, response)
		return
	}

	(*bee).SetDescription(pps.Bee.Description)
	(*bee).SetOptions(pps.Bee.Options)

	if pps.Bee.Active {
		bees.RestartBee(bee)
	} else {
		(*bee).Stop()
	}

	resp.AddBee(bee)
	resp.Send(response)
}
