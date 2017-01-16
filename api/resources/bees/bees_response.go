package bees

import (
	"time"

	"github.com/muesli/beehive/bees"

	"github.com/muesli/smolder"
)

// BeeResponse is the common response to 'bee' requests
type BeeResponse struct {
	smolder.Response

	Bees []beeInfoResponse `json:"bees,omitempty"`
	bees []*bees.BeeInterface
}

type beeInfoResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Namespace   string           `json:"namespace"`
	Description string           `json:"description"`
	LastAction  time.Time        `json:"lastaction"`
	LastEvent   time.Time        `json:"lastevent"`
	Active      bool             `json:"active"`
	Options     []bees.BeeOption `json:"options"`
	Events      []bees.EventDescriptor
	Actions     []bees.ActionDescriptor
}

// Init a new response
func (r *BeeResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Bees = []beeInfoResponse{}
}

// AddBee adds a bee to the response
func (r *BeeResponse) AddBee(bee *bees.BeeInterface) {
	r.bees = append(r.bees, bee)
	r.Bees = append(r.Bees, prepareBeeResponse(r.Context, bee))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *BeeResponse) EmptyResponse() interface{} {
	if len(r.bees) == 0 {
		var out struct {
			Bees interface{} `json:"bees"`
		}
		out.Bees = []beeInfoResponse{}
		return out
	}
	return nil
}

func prepareBeeResponse(context smolder.APIContext, bee *bees.BeeInterface) beeInfoResponse {
	//	ctx := context.(*context.APIContext)
	resp := beeInfoResponse{
		ID:          (*bee).Name(),
		Name:        (*bee).Name(),
		Namespace:   (*bee).Namespace(),
		Description: (*bee).Description(),
		LastAction:  (*bee).LastAction(),
		LastEvent:   (*bee).LastEvent(),
		Active:      (*bee).IsRunning(),
		Options:     (*bee).Options(),
	}

	return resp
}
