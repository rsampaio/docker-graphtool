package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/graphdriver"
	_ "github.com/docker/docker/daemon/graphdriver/aufs"
	_ "github.com/docker/docker/daemon/graphdriver/overlay"
	"github.com/docker/docker/graph"
	"syscall"
)

type GraphTool struct {
	DockerRoot string
}

// NewGraphTool create new graphtool handler
func NewGraphTool(dockerRoot string) *GraphTool {
	return &GraphTool{DockerRoot: dockerRoot}
}

// (g *GraphTool) Mount ...
func (g *GraphTool) Mount(imageName string, dest string, options []string) {
	driver, err := graphdriver.New(g.DockerRoot, make([]string, 0))
	if err != nil {
		log.Fatal(err.Error())
	}

	graphHandler, err := graph.NewGraph(g.DockerRoot+"/graph", driver)
	if err != nil {
		log.Fatal(err.Error())
	}

	tagCfg := &graph.TagStoreConfig{
		Graph: graphHandler,
	}

	tagStore, err := graph.NewTagStore(g.DockerRoot+"/repositories-"+driver.String(), tagCfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	image, err := tagStore.LookupImage(imageName)
	if image == nil {
		log.Fatal(err.Error())
	}

	path, _ := driver.Get(image.ID, "graphtool")
	log.Infof("%+v", path)

	if err = syscall.Mount(path, dest, "none",
		syscall.MS_BIND|syscall.MS_NOATIME|syscall.MS_NOSUID, ""); err != nil {
		log.Fatal(err.Error())
	}

}
