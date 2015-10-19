package main

import (
	"github.com/Sirupsen/logrus"

	"github.com/docker/docker/daemon/graphdriver"
	_ "github.com/docker/docker/daemon/graphdriver/aufs"
	_ "github.com/docker/docker/daemon/graphdriver/overlay"
	"github.com/docker/docker/graph"
	"github.com/docker/docker/image"
)

type GraphTool struct {
	DockerRoot   string
	graphDriver  graphdriver.Driver
	graphHandler *graph.Graph
	logger       *logrus.Logger
}

// NewGraphTool create new graphtool handler
func NewGraphTool(dockerRoot string) *GraphTool {
	return &GraphTool{
		DockerRoot: dockerRoot,
		logger: logrus.WithFields(
			logrus.Fields{
				"docker_root": dockerRoot,
			},
		).Logger,
	}
}

// initDriver ...
func (g *GraphTool) InitDriver() error {
	var err error
	g.graphDriver, err = graphdriver.New(g.DockerRoot, make([]string, 0))
	if err != nil {
		return err
	}
	g.graphHandler, err = graph.NewGraph(g.DockerRoot+"/graph", g.graphDriver)
	if err != nil {
		return err
	}
	return nil
}

// lookupImage ...
func (g *GraphTool) LookupImage(imageName string) (*image.Image, error) {
	tagCfg := &graph.TagStoreConfig{
		Graph: g.graphHandler,
	}
	tagStore, err := graph.NewTagStore(g.DockerRoot+"/repositories-"+g.graphDriver.String(), tagCfg)
	if err != nil {
		return nil, err
	}
	image, err := tagStore.LookupImage(imageName)
	if err != nil {
		return nil, err
	}
	return image, err
}
