package main

import (
	"fmt"
	"image/png"
	"image/color"
	"os"
	"image"
)

func loadColorMap(path string) [][]color.Color {

	img := loadImage(path)
	colorMap := [][]color.Color{}
	for x := 0; x < img.Bounds().Max.X; x++ {
		colorMapLine := []color.Color{}
		for y := 0; y < img.Bounds().Max.Y; y++ {
			colorMapLine = append(colorMapLine, img.At(x,y))
		}
		colorMap = append(colorMap, colorMapLine)
		}
	return colorMap
	}


func loadHeightMap(path string) [][]byte {

	img := loadImage(path)
	heightMap := [][]byte{}
	for x:= 0; x < img.Bounds().Max.X; x++ {
		heightMapLine := []byte{}
		for y := 0; y < img.Bounds().Max.Y; y++ {
			heightMapByte := color.GrayModel.Convert(img.At(x,y)).(color.Gray).Y
			heightMapLine = append(heightMapLine, heightMapByte)
		}
		heightMap = append(heightMap, heightMapLine)
		}
	
	return heightMap
	}


func loadImage(path string) image.Image {
	
	imageFile, err := os.Open(path) 
	
	if (err != nil) {
		fmt.Printf("Failed to open file: %v Error:%v", path, err.Error())
		os.Exit(1)
	}
	defer imageFile.Close()


	img, err := png.Decode(imageFile)
	
	if (err != nil) {
		fmt.Printf("Failed to parse image, Error:%v",err.Error())
		os.Exit(1)
	}

	fmt.Printf("Image bounds: %v\n", img.Bounds())
	return img
}

func main() {
	
	fmt.Println("Voxels")
	fmt.Println("Loading Heightmap...")
	heightMap := loadHeightMap("height_map.png")
	fmt.Printf("First heightmap byte: %v\n", heightMap[0][0])
	fmt.Println("Loading Colormap...")
	colorMap := loadColorMap("color_map.png")
	fmt.Printf("First colormap color: %+v\n", colorMap[0][0])
}
