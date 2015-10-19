package main

import (
	"github.com/docker/docker/runconfig"
	"syscall"
)

// (g *GraphTool) Mount ...
func (g *GraphTool) Mount(imageName string, dest string, options []string) error {
	var err error
	if err = g.InitDriver(); err != nil {
		return err
	}

	image, err := g.LookupImage(imageName)
	if image == nil {
		return err
	}

	fake_image, err := g.graphHandler.Create(nil, "daedbeef", image.ID, "", "", &runconfig.Config{}, &runconfig.Config{})
	if err != nil {
		return err
	}

	path, _ := g.graphDriver.Get(fake_image.ID, "graphtool")
	if err = syscall.Mount(path, dest, "none", syscall.MS_BIND, ""); err != nil {
		return err
	}

	// We don't need the original reference anymore
	// the bind mount is still the last reference to the filesystem
	g.graphDriver.Put(fake_image.ID)

	return nil
}

// (g *GraphTool) Unmount ...
func (g *GraphTool) Unmount(target string) error {
	return syscall.Unmount(target, 0)
}
