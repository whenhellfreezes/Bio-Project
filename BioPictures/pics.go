package BioPictures

import(
"BioProject/cell"
"fmt"
"math"
"image"
"image/color"
"image/png"
"image/draw"
"strconv"
"os"
"bufio"
)


const(
OUTPUTNAME = "Output"
SUFFIX = ".png"
ORIGINAL = "/home/nicholas/go/src/BioProject/BioPictures/template.png"
PICWIDTH = 346
PICHEIGHT = 346
HEIGHT = 20 //Size of grid
WIDTH = 20
COLORRANGE = 37
PICTURESPECIFIC = .56 //Tinker this to make things look better
DESIREDLEVEL = .2
)


func Pics(input [][]cell.Cell, number int) {
	file, _ := os.Open(ORIGINAL) //Open file
	reader := bufio.NewReader(file) //Load template image to io
	src, _ := png.Decode(reader) //Decode png
	stringnum := strconv.Itoa(number) //Convert number into string
	f, _ := os.Create(OUTPUTNAME+stringnum+SUFFIX) //Open file to be created
	//Transform template png into a RGBA
	b := src.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, src, b.Min, draw.Src)
	//Make the RGBA to overlay
	var max float64 = 0
	var min float64 = 999999999999
	var sum float64
	var passed bool = true
	for i := 0; i < WIDTH; i++ {
		for j := 0; j < HEIGHT; j++ {
			if input[i][j].City {
				if input[i][j].JI+input[i][j].AI >= DESIREDLEVEL {
					passed = false
				}
			}
			sum = input[i][j].JI+input[i][j].AI
			if sum < min {
				min = sum
			}
			if sum > max {
				max = sum
			}
		}
	}
	//difference := max-min
	fmt.Println(passed)
	overlay := image.NewRGBA(b)
	for i := 0; i < PICWIDTH; i++ { //Almost literally magic is what follows
		for j := 0; j < PICHEIGHT; j++ {
			cellx := int(math.Floor(float64(i)/PICWIDTH*WIDTH))
			celly := int(math.Floor(float64(j)/PICHEIGHT*HEIGHT))
			ratio := (input[cellx][celly].JI+input[cellx][celly].AI-PICTURESPECIFIC)
			overlay.Set(i, j, color.RGBA{uint8(math.Floor(ratio*COLORRANGE+0.5)),uint8(math.Floor(ratio*COLORRANGE+0.5)),138-uint8(math.Floor(ratio*COLORRANGE+0.5)),140})
		}
	}
	draw.Draw(m, b, overlay, b.Min, draw.Over)
	if err := png.Encode(f,m); err != nil {
		fmt.Println("OH NOESS!!!!")
	}
}
