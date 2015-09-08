package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/graphdriver"
	_ "github.com/docker/docker/daemon/graphdriver/overlay"
	"github.com/docker/docker/graph"
)

type GraphTool struct {
	DockerRoot string
}

// NewGraphTool create new graphtool handler
func NewGraphTool(dockerRoot string) *GraphTool {
	return &GraphTool{DockerRoot: dockerRoot}
}

// (g *GraphTool) Mount ...
func (g *GraphTool) Mount(image string, dest string, options []string) {
	driver, err := graphdriver.New(g.DockerRoot, make([]string, 0))
	if err != nil {
		log.Fatal(err.Error())
	}

	graphHandler, err := graph.NewGraph(g.DockerRoot, driver)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info(graphHandler)
}
