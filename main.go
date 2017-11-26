package main

import (
	"flag"
	"fmt"
	"github.com/andybons/gogif"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"os"
)

func loadHeightMap(path string) heightMap {

	img := loadImage(path)
	heightMap := heightMap{}
	for x := 0; x < img.Bounds().Max.X; x++ {
		heightMapLine := []byte{}
		for y := 0; y < img.Bounds().Max.Y; y++ {
			heightMapByte := color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
			heightMapLine = append(heightMapLine, heightMapByte)
		}
		heightMap = append(heightMap, heightMapLine)
	}

	return heightMap
}

func loadImage(path string) image.Image {

	imageFile, err := os.Open(path)

	if err != nil {
		fmt.Printf("Failed to open file: %v Error:%v", path, err.Error())
		os.Exit(1)
	}
	defer imageFile.Close()

	img, err := png.Decode(imageFile)

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

		fmt.Println("Cannot save output file, error: %v\n", err.Error())
		os.Exit(1)
	}

	defer outfile.Close()
	png.Encode(outfile, image)
}

func saveGIF(path string, images []image.Image, frameDelay int) {

	outfile, err := os.Create(path)
	if err != nil {

		fmt.Println("Cannot save output file, error: %v\n", err.Error())
		os.Exit(1)
	}

	defer outfile.Close()

	outGif := &gif.GIF{}
	quantizer := gogif.MedianCutQuantizer{NumColor: 256}
	for _, img := range images {
		bounds := img.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer.Quantize(palettedImage, bounds, img, image.ZP)
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, frameDelay)
	}
	gif.EncodeAll(outfile, outGif)
}

func main() {

	var viewPointX = flag.Float64("x", 0.0, "X coordinate of the render viewpoint")
	var viewPointY = flag.Float64("y", 0.0, "Y coordinate of the render viewpoint")
	var viewPointZ = flag.Float64("z", 120.0, "Z coordinate of the render viewpoint")
	var viewPointPhi = flag.Float64("phi", 0.0, "Phi angle of the render viewpoint")
	var renderGif = flag.Bool("gif", false, "Flag for rendering a 360 degree GIF")
	flag.Parse()

	fmt.Println("Voxels")
	fmt.Println("Loading Heightmap...")
	heightMap := loadHeightMap("height_map.png")
	fmt.Printf("First heightmap byte: %v\n", heightMap[0][0])
	fmt.Println("Loading Colormap...")
	colorMap := loadImage("color_map.png")
	fmt.Printf("First colormap color: %+v\n", colorMap.At(0, 0))
	fmt.Println("Initializing Renderer...")
	options := RenderOptions{horizonHeight: 180.0, heightScale: 200.0, viewDistance: 200, screenWidth: 800, screenHeight: 600}
	renderer := NewVoxelRenderer(heightMap, colorMap, &options)

	if *renderGif {
		images := []image.Image{}
		coneImages := []image.Image{}

		for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
			newImage, newConeImage := renderer.Render(*viewPointX, *viewPointY, *viewPointZ, angle)
			images = append(images, newImage)
			coneImages = append(coneImages, newConeImage)
		}
		saveGIF("out.gif", images, 10)
		saveGIF("cones.gif", coneImages,10)
	} else {
		image, cone := renderer.Render(*viewPointX, *viewPointY, *viewPointZ, *viewPointPhi)
		savePNG("out.png", image)
		savePNG("cone.png", cone)
	}

}
