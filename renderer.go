package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

var blue color.RGBA = color.RGBA{0, 0, 255, 255}

type heightMap [][]byte

type RenderOptions struct {
	viewDistance, screenWidth, screenHeight int
	horizonHeight, heightScale              float64
}

type VoxelRenderer struct {
	hMap    heightMap
	cMap    image.Image
	options *RenderOptions
}

func NewVoxelRenderer(hMap heightMap, cMap image.Image, options *RenderOptions) VoxelRenderer {
	return VoxelRenderer{hMap: hMap, cMap: cMap, options: options}
}

func (renderer VoxelRenderer) Render(x, y, height, phi float64) image.Image {
	renderBounds := image.Rect(0, 0, renderer.options.screenWidth, renderer.options.screenHeight)
	renderImage := image.NewRGBA(renderBounds)

	fillImage(renderImage, blue)

	sinphi := math.Sin(phi)
	cosphi := math.Cos(phi)

	mapRectangle := image.Rect(0, 0, renderer.cMap.Bounds().Max.X, renderer.cMap.Bounds().Max.Y)

	//fmt.Println(sinphi)
	//fmt.Println(cosphi)

	for z := renderer.options.viewDistance; z > 1; z-- {
		fmt.Println(z)
		floatZ := float64(z)

		pointLeft := image.Pt(int((-cosphi*floatZ-sinphi*floatZ)+x), int((sinphi*floatZ-cosphi*floatZ)+y))
		pointRight := image.Pt(int((cosphi*floatZ-sinphi*floatZ)+x), int((-sinphi*floatZ-cosphi*floatZ)+y))
		pointLeft = pointLeft.Mod(mapRectangle)
		pointRight = pointRight.Mod(mapRectangle)
		//fmt.Println(pointLeft, pointRight)
		dx := (pointRight.X - pointLeft.X) / renderer.options.screenWidth
		dy := (pointRight.Y - pointLeft.Y) / renderer.options.screenWidth

		for i := 0; i < renderer.options.screenWidth; i++ {

			heightOnScreen := float64(byte(height)-renderer.hMap[pointLeft.X][pointLeft.Y])/floatZ*renderer.options.heightScale + renderer.options.horizonHeight
			drawVerticalLineFromPoint(renderImage, image.Pt(i, int(heightOnScreen)), renderer.cMap.At(pointLeft.X, pointLeft.Y))
			pointLeft.X += dx
			pointLeft.Y += dy
		}
	}

	return renderImage
}

func drawVerticalLineFromPoint(img *image.RGBA, startPoint image.Point, color color.Color) {

	drawRectangle := image.Rect(startPoint.X, startPoint.Y, startPoint.X+1, img.Bounds().Max.Y)
	draw.Draw(img, drawRectangle, &image.Uniform{color}, image.ZP, draw.Src)
}

func fillImage(img *image.RGBA,color color.Color) {
 draw.Draw(img, img.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
}
