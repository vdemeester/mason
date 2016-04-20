package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/engine-api/client"
	"github.com/spf13/cobra"
	"github.com/vdemeester/mason/dockerfile"

	log "github.com/Sirupsen/logrus"
)

// dockerfileCmd represents the dockerfile command
var dockerfileCmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "Build client-side using a Dockerfile (like docker build)",
	Long: `Build an image using Dockerfile, just like docker build. It supports
what docker build support, but does the build client-side.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			log.SetLevel(log.DebugLevel)
		}
		// FIXME(vdemeester) more options, and cleaner
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "dockerfile requires 1 argument")
			os.Exit(1)
		}
		contextPath, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while locating the context : %s", err.Error())
			os.Exit(1)
		}
		// FIXME(vdemeester) handle errors
		dockerfileName, _ := cmd.Flags().GetString("file")
		tags, _ := cmd.Flags().GetStringSlice("tag")
		cli, err := client.NewEnvClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting a docker client : %s", err.Error())
			os.Exit(1)
		}
		builder, err := dockerfile.NewBuilder(cli, contextPath, dockerfileName, tags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error constructing the builder : %s", err.Error())
			os.Exit(1)
		}
		if err := builder.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error building : %s", err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(dockerfileCmd)

	dockerfileCmd.Flags().StringP("file", "f", "", "Name of the Dockerfile (Default is `PATH/Dockerfile`)")
	dockerfileCmd.Flags().StringSliceP("tag", "t", []string{}, "Name and optionnaly a tag in the 'name:tag' format")
}
