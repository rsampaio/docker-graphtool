package main

import (
	"github.com/docopt/docopt-go"
	"strings"
)

func main() {
	usage := `Docker graphtool.

Usage:
  dg mount [--options=<mount_options>] [<image>] [<dest>]
  dg umount [--force] <temp_image>
  dg bundle <image> <bundle_file>

Options:
  -h --help                        This help
  -f --force                       Force unmount
  -o <options> --options=<options> Mount options
`
	arguments, err := docopt.Parse(usage, nil, true, "docker dist 0.1", false)
	if err != nil {
		panic(err.Error())
	}

	graphtool := NewGraphTool("/var/lib/docker")

	if arguments["mount"].(bool) {
		image := arguments["<image>"].(string)
		dest := arguments["<dest>"].(string)
		options := []string{""}
		if arguments["--options"] != nil {
			options = strings.Split(arguments["--options"].(string), ",")
		}

		if err := graphtool.Mount(image, dest, options); err != nil {
			graphtool.logger.Fatal(err.Error())
		}
	} else if arguments["umount"].(bool) {
		graphtool.Unmount(arguments["<mount_point>"].(string))
	} else if arguments["bundle"].(bool) {
		graphtool.Bundle(arguments["<image>"].(string), arguments["<bundle_file>"].(string))
	}
}
