package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"fmt"
)

var white color.RGBA = color.RGBA{255,255,255, 255}
var red color.RGBA = color.RGBA{255, 0, 0, 255}
var green color.RGBA = color.RGBA{0,255,0,255}
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

func (renderer VoxelRenderer) Render(x, y, height, phi float64) (image.Image, image.Image) {
	renderBounds := image.Rect(0, 0, renderer.options.screenWidth, renderer.options.screenHeight)
	renderImage := image.NewRGBA(renderBounds)

	fillImage(renderImage, blue)

	sinphi := math.Sin(phi)
	cosphi := math.Cos(phi)

	mapRectangle := image.Rect(0, 0, renderer.cMap.Bounds().Max.X, renderer.cMap.Bounds().Max.Y)

	
	coneBounds := renderer.cMap.Bounds()
  	cone := image.NewRGBA(coneBounds) 
  	draw.Draw(cone, coneBounds, renderer.cMap, coneBounds.Min, draw.Src) 

	for z := renderer.options.viewDistance; z > 1; z-- {
		floatZ := float64(z)

		pointLeft := image.Pt(int((-cosphi*floatZ-sinphi*floatZ)+x), int((sinphi*floatZ-cosphi*floatZ)+y))
		pointRight := image.Pt(int((cosphi*floatZ-sinphi*floatZ)+x), int((-sinphi*floatZ-cosphi*floatZ)+y))
		fmt.Println("Before")
		fmt.Println(pointLeft, pointRight)
		pointLeft = pointLeft.Mod(mapRectangle)
		pointRight = pointRight.Mod(mapRectangle)
		
		cone.Set(pointLeft.X, pointLeft.Y, red)
		cone.Set(pointRight.X, pointRight.Y, green)

		fmt.Println("After")
		fmt.Println(pointLeft, pointRight)
		dx := float64(pointRight.X - pointLeft.X) / float64(renderer.options.screenWidth)
		fmt.Println(dx)
		dy := float64(pointRight.Y - pointLeft.Y) / float64(renderer.options.screenWidth)
		fmt.Println(dy)
		for i := 0; i < renderer.options.screenWidth; i++ {

			cone.Set(pointLeft.X, pointLeft.Y, blue)
			heightOnScreen := float64(byte(height)-renderer.hMap[pointLeft.X][pointLeft.Y])/floatZ*renderer.options.heightScale + renderer.options.horizonHeight
			drawVerticalLineFromPoint(renderImage, image.Pt(i, int(heightOnScreen)), renderer.cMap.At(pointLeft.X, pointLeft.Y))
			pointLeft.X += int(dx)
			pointLeft.Y += int(dy)
		}
	}

	return renderImage, cone
}

func drawVerticalLineFromPoint(img *image.RGBA, startPoint image.Point, color color.Color) {

	drawRectangle := image.Rect(startPoint.X, startPoint.Y, startPoint.X+1, img.Bounds().Max.Y)
	draw.Draw(img, drawRectangle, &image.Uniform{color}, image.ZP, draw.Src)
}

func fillImage(img *image.RGBA,color color.Color) {
 draw.Draw(img, img.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
}
