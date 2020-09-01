package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ecrcp",
		Usage: "copy images from docker hub to AWS ECR",
		Action: func(c *cli.Context) error {
			// call function that pulls image from docker hub and pushed to ecr.
			srcimage := c.Args().Get(0)
			destimage := c.Args().Get(1)
			pullpush(srcimage, destimage)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

} // main function ends
