package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kubepush",
		Short: "Push an image from Podman to Kind",
		RunE: func(cmd *cobra.Command, args []string) error {
			image, _ := cmd.Flags().GetString("image")
			clusterName, _ := cmd.Flags().GetString("cluster-name")
			removeImage, _ := cmd.Flags().GetBool("remove-image")

			log.Println("Saving image file...")

			tmpFile, err := os.CreateTemp("", "podman-image-*.tar")
			if err != nil {
				return fmt.Errorf("error creating temp file: %d", err)
			}
			defer func(name string) {
				log.Println("Removing image file...")
				err := os.Remove(name)
				if err == nil {
					log.Println("Image file removed successfully")
				} else {
					log.Fatalf("error removing temp file: %d", err)
				}
			}(tmpFile.Name())

			c := exec.Command("sh", "-c", "podman save -o "+tmpFile.Name()+" "+image)
			output, err := c.CombinedOutput()
			if err != nil {
				if output != nil {
					return errors.New(string(output))
				} else {
					return fmt.Errorf("error when saving the image file: %d", err)
				}
			}

			log.Println("Image file saved. Loading image into the cluster...")
			c = exec.Command("sh", "-c", "kind load image-archive --name "+clusterName+" "+tmpFile.Name())
			output, err = c.CombinedOutput()
			if err != nil {
				if output != nil {
					return errors.New(string(output))
				} else {
					return fmt.Errorf("error when pushing the image to the cluster: %d", err)
				}
			}

			log.Println(fmt.Sprintf("%s pushed to the %s cluster.", image, clusterName))

			if removeImage {
				log.Println("Removing the image...")

				c = exec.Command("sh", "-c", "podman rmi "+image)
				output, err = c.CombinedOutput()
				if err != nil {
					if output != nil {
						return errors.New(string(output))
					} else {
						return fmt.Errorf("error when removing the image: %d", err)
					}
				}

				log.Println("Image removed successfully")
			}

			return nil
		},
	}

	rootCmd.Flags().StringP("image", "i", "", "image to push")
	rootCmd.Flags().StringP("cluster-name", "c", "", "the name of the cluster")
	rootCmd.Flags().BoolP("remove-image", "r", false, "remove the image after pushing it to the cluster")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
