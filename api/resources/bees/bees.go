package bees

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// BeeResource is the resource responsible for /bees
type BeeResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &BeeResource{}
	_ smolder.GetSupported   = &BeeResource{}
	_ smolder.PutSupported   = &BeeResource{}
)

// Register this resource with the container to setup all the routes
func (r *BeeResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "BeeResource"
	r.TypeName = "bee"
	r.Endpoint = "bees"
	r.Doc = "Manage bees"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}
