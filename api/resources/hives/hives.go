package hives

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// HiveResource is the resource responsible for /hives
type HiveResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &HiveResource{}
	_ smolder.GetSupported   = &HiveResource{}
)

// Register this resource with the container to setup all the routes
func (r *HiveResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "HiveResource"
	r.TypeName = "hive"
	r.Endpoint = "hives"
	r.Doc = "Manage hives"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}
