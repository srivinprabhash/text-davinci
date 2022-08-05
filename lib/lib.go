package lib

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/nfnt/resize"
)

// Converts an image to grayscale
func toGrayScale(ogImage image.Image) image.Image {

	size := ogImage.Bounds().Size()
	rectangle := image.Rect(0, 0, size.X, size.Y)

	grayScaledImg := image.NewRGBA(rectangle)

	for x := 0; x < size.X; x++ {

		for y := 0; y < size.Y; y++ {

			ogPixel := ogImage.At(x, y)
			pixel := color.GrayModel.Convert(ogPixel)
			grayScaledImg.Set(x, y, pixel)

		}

	}
	return grayScaledImg

}

// Convert the gray scale version of the image
// to a black and white image.
//
// whiteThreshold - determines the hexa color position to
// divide the black and white
func toBlackAndWhite(ogImage image.Image, whiteThreashold uint8) image.Image {

	size := ogImage.Bounds().Size()
	rectangle := image.Rect(0, 0, size.X, size.Y)

	blackAndWhiteImg := image.NewRGBA(rectangle)

	// Loop through each pixel and determines
	// if it should be black or white
	// depends on the color value and given
	// white threshold.
	for x := 0; x < size.X; x++ {

		for y := 0; y < size.Y; y++ {

			ogPixel := ogImage.At(x, y)
			originalColor := color.RGBAModel.Convert(ogPixel).(color.RGBA)
			modifiedColorValue := originalColor.R

			if modifiedColorValue >= uint8(whiteThreashold) {
				modifiedColorValue = 255
			} else {
				modifiedColorValue = 0
			}

			pixel := color.RGBA{
				R: modifiedColorValue,
				G: modifiedColorValue,
				B: modifiedColorValue,
				A: originalColor.A,
			}

			blackAndWhiteImg.Set(x, y, pixel)

		}

	}

	return blackAndWhiteImg

}

// Finds the optimal white threshold value so
// that the number of black pixels and white
// pixels are close.
func toOptimalBlackAndWhite(ogImage image.Image) image.Image {

	size := ogImage.Bounds().Size()

	totalPixels := size.X * size.Y
	halfPixels := int(math.Round(float64(totalPixels) / 2))

	var closestPixelDiff int
	var closestImg image.Image

	// Repeats all the possible white threshold values
	// until the optimal threshold value is found.
	for i := 0; i < 256; i++ {

		blackAndWhiteImg := toBlackAndWhite(ogImage, uint8(i))
		var blackPixels int

		for x := 0; x < size.X; x++ {
			for y := 0; y < size.Y; y++ {
				pixel := blackAndWhiteImg.At(x, y)
				pixelColor := color.RGBAModel.Convert(pixel).(color.RGBA)

				if pixelColor.R > 100 {
					blackPixels++
				}
			}
		}

		pixelDiff := int(math.Abs(float64(halfPixels - blackPixels)))

		if closestPixelDiff == 0 || pixelDiff+1000 < closestPixelDiff {
			closestPixelDiff = pixelDiff
			closestImg = blackAndWhiteImg
		}
	}

	return closestImg

}

type TextOptions struct {
	Wc                           string `default:" "`
	Bc                           string `default:"⠢"`
	Output                       string `default:"output.txt"`
	DisableThresholdOptimization bool
	Flip                         bool
}

func WriteToTxt(ogImage image.Image, options ...TextOptions) {

	// Resize the image to 60x60
	resizedImage := resize.Resize(60, 60, ogImage, resize.Lanczos3)

	// Convert the image to grayscale
	grayScaledImg := toGrayScale(resizedImage)

	// Convert the grayscaled image to a
	// fully black and white image. If white threshold
	// optimization is disabled, uses 128 as the threshold
	var blackAndWhiteImage image.Image

	if options[0].DisableThresholdOptimization {
		blackAndWhiteImage = toBlackAndWhite(grayScaledImg, 128)
	} else {
		blackAndWhiteImage = toOptimalBlackAndWhite(grayScaledImg)
	}

	// Create the oupt text file
	textFile, err := os.Create(options[0].Output)
	if err != nil {
		log.Fatalln(err)
	}
	defer textFile.Close()

	// Loop through each pixel of the image and
	// write characters to the text file depending
	// on the color of the pixel.
	size := blackAndWhiteImage.Bounds().Size()

	for x := 0; x < size.X; x++ {

		// New line
		if x != 0 {
			textFile.WriteString("\n")
		}

		for y := 0; y < size.Y; y++ {

			pixel := blackAndWhiteImage.At(y, x)
			pixelColor := color.RGBAModel.Convert(pixel).(color.RGBA)

			if pixelColor.R == 0 {
				if options[0].Flip {
					textFile.WriteString(options[0].Bc)
				} else {
					textFile.WriteString(options[0].Wc)
				}
			} else {
				if options[0].Flip {
					textFile.WriteString(options[0].Wc)
				} else {
					textFile.WriteString(options[0].Bc)
				}
			}

		}
	}

}
