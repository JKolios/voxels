package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/andybons/gogif"
	"github.com/icza/mjpeg"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

func loadImage(path string) image.Image {

	imageFile, err := os.Open(path)

	if err != nil {
		fmt.Printf("Failed to open file: %v Error:%v", path, err.Error())
		os.Exit(1)
	}
	defer imageFile.Close()

	img, imgType, err := image.Decode(imageFile)
	fmt.Printf("Image Type: %s", imgType)
	if err != nil {
		fmt.Printf("Failed to parse image, Error:%v", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Image bounds: %v\n", img.Bounds())
	return img
}

func savePNG(path string, image image.Image) {
	outfile, err := os.Create(path)
	if err != nil {

		fmt.Printf("Cannot save output file, error: %v\n", err.Error())
		os.Exit(1)
	}

	defer outfile.Close()
	png.Encode(outfile, image)
}

func saveGIF(path string, images []image.Image, frameDelay int) {

	outfile, err := os.Create(path)
	if err != nil {

		fmt.Printf("Cannot save output file, error: %v\n", err.Error())
		os.Exit(1)
	}

	defer outfile.Close()

	outGif := &gif.GIF{}
	quantizer := gogif.MedianCutQuantizer{NumColor: 1024}
	for count, img := range images {
		fmt.Printf("Quantizing image %v\n", count)
		bounds := img.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer.Quantize(palettedImage, bounds, img, image.ZP)
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, frameDelay)
	}
	gif.EncodeAll(outfile, outGif)
}

func saveMJPEG(path string, images []image.Image, fps int) {

	frameSize := images[0].Bounds().Max
	outFile, err := mjpeg.New(path, int32(frameSize.X), int32(frameSize.Y), int32(fps))
	defer outFile.Close()
	if err != nil {

		fmt.Printf("Cannot create output file, error: %v\n", err.Error())
		os.Exit(1)
	}

	for _, image := range images {

		buf := &bytes.Buffer{}
		if err := jpeg.Encode(buf, image, nil); err != nil {

			fmt.Printf("Cannot encode image into JPEG, error: %v\n", err.Error())
			os.Exit(1)
		}

		if err := outFile.AddFrame(buf.Bytes()); err != nil {

			fmt.Printf("Cannot add frame to a JPEG file, error: %v\n", err.Error())
			os.Exit(1)
		}
	}
}

func main() {

	var viewPointX = flag.Float64("x", 0.0, "X coordinate of the render viewpoint")
	var viewPointY = flag.Float64("y", 0.0, "Y coordinate of the render viewpoint")
	var viewPointZ = flag.Float64("z", 120.0, "Z coordinate of the render viewpoint")
	var viewPointPhi = flag.Float64("phi", 0.0, "Phi angle of the render viewpoint")
	var renderGif = flag.Bool("gif", false, "Flag for rendering a 360 degree GIF")
	var renderPng = flag.Bool("png", true, "Flag for rendering a single PNG frame")
	var renderMjpeg = flag.Bool("mjpeg", true, "Flag for rendering a 360 degree MJPEG")
	flag.Parse()

	fmt.Println("Voxels")
	fmt.Println("Loading Heightmap...")
	heightMap := loadImage("height_map.png")
	fmt.Printf("First heightmap byte: %v\n", heightMap.At(0, 0))
	fmt.Println("Loading Colormap...")
	colorMap := loadImage("color_map.png")
	fmt.Printf("First colormap color: %+v\n", colorMap.At(0, 0))
	fmt.Println("Loading Skybox...")
	skyBox := loadImage("skybox.png")
	fmt.Println("Initializing Renderer...")
	options := RenderOptions{
		horizonHeight:  180.0,
		heightScale:    200.0,
		renderDistance: 200,
		renderBounds:   image.Rect(0, 0, 800, 600),
	}
	renderer := NewVoxelRenderer(heightMap, colorMap, skyBox, &options)

	if *renderPng {
		image := renderer.Render(*viewPointX, *viewPointY, *viewPointZ, *viewPointPhi)
		savePNG("out.png", image)

	}
	if *renderGif {
		images := []image.Image{}

		for angle := 0.0; angle < 2*math.Pi; angle += 0.2 {
			newImage := renderer.Render(*viewPointX, *viewPointY, *viewPointZ, angle)
			fmt.Printf("Render complete for angle: %v\n", angle)
			images = append(images, newImage)
		}
		saveGIF("out.gif", images, 10)
	}
	if *renderMjpeg {

		images := []image.Image{}

		for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
			newImage := renderer.Render(*viewPointX, *viewPointY, *viewPointZ, angle)
			fmt.Printf("Render complete for angle: %v\n", angle)
			images = append(images, newImage)
		}
		saveMJPEG("out.mjpeg", images, 24)
	}

}
