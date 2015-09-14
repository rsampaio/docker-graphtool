package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docopt/docopt-go"
	"strings"
)

func main() {
	usage := `Docker graphtool.

Usage:
  dg mount [--options=<mount_options>] [<image>] [<dest>]
  dg umount [--force] <temp_image>

Options:
  -h --help                        This help
  -f --force                       Force unmount
  -o <options> --options=<options> Mount options
`
	arguments, err := docopt.Parse(usage, nil, true, "docker dist 0.1", false)
	if err != nil {
		log.Fatal(err.Error())
	}

	graphTool := NewGraphTool("/var/lib/docker")

	if arguments["mount"].(bool) {
		image := arguments["<image>"].(string)
		dest := arguments["<dest>"].(string)
		options := []string{""}
		if arguments["--options"] != nil {
			options = strings.Split(arguments["--options"].(string), ",")
		}

		log.WithFields(log.Fields{
			"options": options,
			"image":   image,
			"dest":    dest,
		}).Info("Mount")

		tempImage, err := graphTool.Mount(image, dest, options)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Printf("%s\n", tempImage)
	} else if arguments["umount"].(bool) {
		log.Info("Unmount")
		graphTool.Unmount(arguments["<temp_image>"].(string))
	}
}
