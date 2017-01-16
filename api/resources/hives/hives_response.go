package hives

import (
	"github.com/muesli/beehive/bees"

	"github.com/muesli/smolder"
)

// HiveResponse is the common response to 'hive' requests
type HiveResponse struct {
	smolder.Response

	Hives []hiveInfoResponse `json:"hives,omitempty"`
	hives []*bees.BeeFactoryInterface
}

type hiveInfoResponse struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Image       string                     `json:"image"`
	Options     []bees.BeeOptionDescriptor `json:"options"`
	Events      []bees.EventDescriptor     `json:"events"`
	Actions     []bees.ActionDescriptor    `json:"actions"`
}

// Init a new response
func (r *HiveResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Hives = []hiveInfoResponse{}
}

// AddHive adds a hive to the response
func (r *HiveResponse) AddHive(hive *bees.BeeFactoryInterface) {
	r.hives = append(r.hives, hive)
	r.Hives = append(r.Hives, prepareHiveResponse(r.Context, hive))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *HiveResponse) EmptyResponse() interface{} {
	if len(r.hives) == 0 {
		var out struct {
			Hives interface{} `json:"hives"`
		}
		out.Hives = []hiveInfoResponse{}
		return out
	}
	return nil
}

func prepareHiveResponse(context smolder.APIContext, hive *bees.BeeFactoryInterface) hiveInfoResponse {
	//	ctx := context.(*context.APIContext)
	resp := hiveInfoResponse{
		ID:          (*hive).Name(),
		Name:        (*hive).Name(),
		Description: (*hive).Description(),
		Image:       "http://localhost:8181/images/" + (*hive).Image(),
		Options:     (*hive).Options(),
		Events:      (*hive).Events(),
		Actions:     (*hive).Actions(),
	}

	return resp
}
