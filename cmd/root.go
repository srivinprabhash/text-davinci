package cmd

import (
	"TextDavinci/lib"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// ROOT COMMAND
var rootCmd = &cobra.Command{
	Use:   "textdavinci",
	Short: "TextDavincii draws your sketch/image using characters in a text file.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Input image path is requried")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// OPEN THE FILE
		inputFile, err := os.Open(args[0])
		if err != nil {
			log.Fatalln(err)
		}
		defer inputFile.Close()

		img, _, err := image.Decode(inputFile)
		if err != nil {
			log.Fatalln(err)
		}

		// Set Flags
		textOptions := lib.TextOptions{}

		// GET FLAGS : white-char
		wc, err := cmd.Flags().GetString("white-char")
		if wc != "" {
			textOptions.Wc = wc
		}

		fmt.Println(wc)

		// GET FLAG : black-char
		bc, err := cmd.Flags().GetString("black-char")
		if bc != "" {
			textOptions.Bc = bc
		}

		outputPath, err := cmd.Flags().GetString("output-path")
		if outputPath != "" {
			textOptions.Output = outputPath
		}

		isFlipped, err := cmd.Flags().GetBool("flip")
		if isFlipped {
			textOptions.Flip = isFlipped
		}

		isThresholdDisabled, err := cmd.Flags().GetBool("disable-threshold")
		if isFlipped {
			textOptions.DisableThresholdOptimization = isThresholdDisabled
		}

		lib.WriteToTxt(img, textOptions)

	},
}

func init() {

	rootCmd.PersistentFlags().StringP("white-char", "w", " ", "Character to draw white pixels")
	rootCmd.PersistentFlags().StringP("black-char", "b", "⠢", "Character to draw black pixels")
	rootCmd.PersistentFlags().StringP("output-path", "o", "output.txt", "Path of the output text file.")
	rootCmd.PersistentFlags().BoolP("disable-threshold", "d", false, "Disables white threshold optimization")
	rootCmd.PersistentFlags().BoolP("flip", "f", false, "Flip white character and black character")

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(0)
	}
}
