package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docopt/docopt-go"
	"strings"
)

func main() {
	usage := `Docker Dist.

Usage:
  docker-dist mount [--options=<mount_options>] [<image>] [<dest>]
  docker-dist umount [--force] <image>
  docker-dist rm <image>

Options:
  -h --help                        This help
  -f --force                       Force unmount
  -o <options> --options=<options> Mount options
`
	arguments, err := docopt.Parse(usage, nil, true, "docker dist 0.1", false)
	if err != nil {
		log.Fatal(err.Error())
	}

	if arguments["mount"].(bool) {
		// Mount docker image
		// - check if cached
		//   - pull if not
		// - create dest
		// - mount
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

		graphTool := NewGraphTool("/var/lib/docker")
		graphTool.Mount(image, dest, options)
	} else if arguments["umount"].(bool) {
		log.Info("Unmount")
	} else if arguments["rm"].(bool) {
		log.Info("Rm")
	}
}
