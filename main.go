package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
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

func saveImage(path string, image image.Image) {
	outfile, err := os.Create(path)
	if err != nil {

		fmt.Println("Cannot save output file, error: %v\n", err.Error())
		os.Exit(1)
	}

	defer outfile.Close()
	png.Encode(outfile, image)
}

func main() {

	fmt.Println("Voxels")
	fmt.Println("Loading Heightmap...")
	heightMap := loadHeightMap("height_map.png")
	fmt.Printf("First heightmap byte: %v\n", heightMap[0][0])
	fmt.Println("Loading Colormap...")
	colorMap := loadImage("color_map.png")
	fmt.Printf("First colormap color: %+v\n", colorMap.At(0, 0))
	fmt.Println("Initializing Renderer...")
	options := RenderOptions{horizonHeight: 120.0, heightScale: 120.0, viewDistance: 100, screenWidth: 800, screenHeight: 600}
	renderer := NewVoxelRenderer(heightMap, colorMap, &options)
	image := renderer.Render(0.0, 0.0, 280.0, 200.0)
	saveImage("out.png", image)
}
