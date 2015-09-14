package main

import (
	"github.com/docker/docker/runconfig"
	"syscall"
)

// (g *GraphTool) Mount ...
func (g *GraphTool) Mount(imageName string, dest string, options []string) (string, error) {
	var err error
	if err = g.InitDriver(); err != nil {
		return "", err
	}

	image, err := g.LookupImage(imageName)
	if image == nil {
		return "", err
	}

	fake_image, err := g.graphHandler.Create(nil, "daedbeef", image.ID, "", "", &runconfig.Config{}, &runconfig.Config{})
	if err != nil {
		return "", err
	}

	path, _ := g.graphDriver.Get(fake_image.ID, "graphtool")
	if err = syscall.Mount(path, dest, "none", syscall.MS_BIND, ""); err != nil {
		return "", err
	}
	defer g.graphDriver.Put(fake_image.ID)
	return fake_image.ID, nil
}

// (g *GraphTool) Unmount ...
func (g *GraphTool) Unmount(imageId string) error {
	return nil
}
