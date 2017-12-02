package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

var white color.RGBA = color.RGBA{255, 255, 255, 255}
var red color.RGBA = color.RGBA{255, 0, 0, 255}
var green color.RGBA = color.RGBA{0, 255, 0, 255}
var blue color.RGBA = color.RGBA{0, 0, 255, 255}

type RenderOptions struct {
	renderDistance             int
	renderBounds               image.Rectangle
	horizonHeight, heightScale float64
}

type VoxelRenderer struct {
	hMap    image.Image
	cMap    image.Image
	skyBox  image.Image
	options *RenderOptions
}

func NewVoxelRenderer(hMap, cMap, skyBox image.Image, options *RenderOptions) VoxelRenderer {

	return VoxelRenderer{
		hMap:    multiplyImage(hMap, 3, true),
		cMap:    multiplyImage(cMap, 3, false),
		skyBox:  skyBox,
		options: options,
	}
}

func (renderer VoxelRenderer) Render(x, y, height, phi float64) image.Image {
	sinphi := math.Sin(phi)
	cosphi := math.Cos(phi)

	renderedFrame := image.NewRGBA(renderer.options.renderBounds)
	draw.Draw(renderedFrame, renderer.options.renderBounds, renderer.skyBox, renderer.options.renderBounds.Min, draw.Src)

	floatScreenWidth := float64(renderer.options.renderBounds.Max.X)

	for z := renderer.options.renderDistance; z > 1; z-- {
		floatZ := float64(z)

		leftX := (-cosphi*floatZ - sinphi*floatZ) + x
		leftY := (sinphi*floatZ - cosphi*floatZ) + y
		rightX := (cosphi*floatZ - sinphi*floatZ) + x
		rightY := (-sinphi*floatZ - cosphi*floatZ) + y

		dx := (rightX - leftX) / floatScreenWidth
		dy := (rightY - leftY) / floatScreenWidth
		for i := 0; i < renderer.options.renderBounds.Max.X; i++ {

			mapHeight := renderer.hMap.At(int(leftX), int(leftY)).(color.Gray).Y
			heightOnScreen := float64(uint8(height)-mapHeight)/floatZ*renderer.options.heightScale + renderer.options.horizonHeight
			drawVerticalLineFromPoint(renderedFrame, image.Pt(i, int(heightOnScreen)), renderer.cMap.At(int(leftX), int(leftY)))
			leftX += dx
			leftY += dy
		}
	}

	return renderedFrame
}

func drawVerticalLineFromPoint(img draw.Image, startPoint image.Point, color color.Color) {

	drawRectangle := image.Rect(startPoint.X, startPoint.Y, startPoint.X+1, img.Bounds().Max.Y)
	draw.Draw(img, drawRectangle, &image.Uniform{color}, image.ZP, draw.Src)
}

func multiplyImage(img image.Image, factor int, grayscale bool) image.Image {

	outputBounds := image.Rect(0, 0, img.Bounds().Max.X*factor, img.Bounds().Max.Y*factor)
	var outputImage draw.Image
	if grayscale {
		outputImage = image.NewGray(outputBounds)
	} else {
		outputImage = image.NewRGBA(outputBounds)
	}
	for i := 0; i < factor; i++ {
		for j := 0; j < factor; j++ {
			startPoint := image.Pt(j*img.Bounds().Max.X, i*img.Bounds().Max.Y)
			copyBounds := image.Rectangle{startPoint, startPoint.Add(img.Bounds().Size())}
			draw.Draw(outputImage, copyBounds, img, image.ZP, draw.Src)
		}
	}
	return outputImage
}
