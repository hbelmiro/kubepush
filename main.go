package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kubepush",
		Short: "Push an image from Podman to Kind",
		RunE: func(cmd *cobra.Command, args []string) error {
			image, _ := cmd.Flags().GetString("image")
			clusterName, _ := cmd.Flags().GetString("cluster-name")

			log.Println("Saving image file...")
			c := exec.Command("sh", "-c", "podman save -o ~/Downloads/image.tar "+image)
			output, err := c.CombinedOutput()
			if err != nil {
				if output != nil {
					return errors.New(string(output))
				} else {
					return fmt.Errorf("error when saving the image file: %d", err)
				}
			}

			log.Println("Image file saved. Loading image into the cluster...")
			c = exec.Command("sh", "-c", "kind load image-archive --name "+clusterName+" ~/Downloads/image.tar")
			output, err = c.CombinedOutput()
			if err != nil {
				if output != nil {
					return errors.New(string(output))
				} else {
					return fmt.Errorf("error when pushing the image to the cluster: %d", err)
				}
			}

			println(fmt.Sprintf("%s pushed to the %s cluster.", image, clusterName))
			return nil
		},
	}

	rootCmd.Flags().StringP("image", "i", "", "image to push")
	rootCmd.Flags().StringP("cluster-name", "c", "", "the name of the cluster")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
