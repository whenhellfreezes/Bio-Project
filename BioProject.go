package main

import (
"BioProject/cell"
"BioProject/BioPictures"
"math"
//"fmt"
)

//This is the holder for all the info each block has
/*
type Cell struct {
	Inside bool
	City bool
	JS float64
	JI float64
	JV float64
	AS float64
	AI float64
	AV float64
}
*/

type Direction struct {
	X int
	Y int
}

var UPDIR Direction = Direction{0, -1}  //We assume that the orientation has 0,0 at the top left and we progress to WIDTH,HEIGHT
var DOWNDIR Direction = Direction{0, 1}
var LEFTDIR Direction = Direction{-1, 0}
var RIGHTDIR Direction = Direction{1, 0}
var DMAP = map[int]Direction{
	0: UPDIR,
	1: DOWNDIR,
	2: LEFTDIR,
	3: RIGHTDIR,
}

const (
	HEIGHT = 20 //Size of grid
	WIDTH = 20

	FENCES = true
	LUREEFFECT = false
	STAYSTILL = true

	UNIMPAIRED = 4 //Movement
	IMPAIRED = 1
	LURESCALE = 2

	TIMESTEPS = 20
	OUTPUTFREQ = 1

	INITJS = 9 //The initail values for each square
	INITJI = 2
	INITJV = 0
	INITAS = 9
	INITAI = 2
	INITAV = 0

	JUVENILLETRANSFER = .06
	ADULTTRANSFER = .06

	VACCINATION = .828

	JUVENILLEDEATH = .02
	ADULTDEATH = .01

	JUVENILLEREPO = 3
	ADULTREPO = 4
)

//Returns 
//the neighbor is fenced in or not
//is it in the grid
//Does it and the neighbor are both inside or both outside the fence
var count int
func checkNeighbor(x int, y int, info [][]cell.Cell, d Direction) (bool, bool, bool) {
	w := d.Y
	z := d.X
	var different bool
	var answer bool
	var grid bool
	local := info[x][y].Inside
	if (x+z < WIDTH && y+w < HEIGHT && x+z >= 0 && y+w >= 0) { //check if inside the grid
		answer = info[x+z][y+w].Inside
		grid = true
	} else {
		//Set false if we are outside the grid
		answer = false
		grid = false
	}
	if local == answer {
		different = false
	} else {
		different = true
	}
	return answer, grid, different
}

//Sometimes its conveint to index things with only one variable
//These functions are to help switch between

func split2d(k int) (int, int) {
	x := k % WIDTH
	y := k / WIDTH
	return x, y
}

func join2d(x, y int) int {
	return x + y*WIDTH
}

//Data structure for local neighborhood structure
type Neighborhood struct {
	rates []float64
}


//Give the movement data for the desired cell
func getNeighborhood(x int, y int, info [][]cell.Cell) Neighborhood {
	grid := make([]bool, 4)
	city := make([]bool, 4)
	interior := make([]bool, 4)
	transistion := make([]float64, 4)
	var answer Neighborhood
	var want int
	var ability int
	var ratelist []float64
	var stay float64

	localcity := info[x][y].City //Is the center square a city


	for j, direction := range DMAP {
		//Go through each direction looking at that neighbor (j is index, direction is associate direction)
		//from DMAP
		tempa, tempb, different := checkNeighbor(x, y, info, direction)
		interior[j] = tempa
		grid[j] = tempb
		if !grid[j] { //Worry about the neighbor being out of bounds
			transistion[j] = 0
			city[j] = false
		} else {
			transistion[j] = 1
			city[j] = info[x+direction.X][y+direction.Y].City
		}
		if FENCES { //If there are fences then change ability accordingly
			if different {
				ability = IMPAIRED
			} else {
				ability = UNIMPAIRED
			}
		} else { //No fences then all directions have the same ability
			ability = 1
		}
		if LUREEFFECT { //If there is a lure then change want accordingly
			if city[j] {
				want = LURESCALE
			} else {
				want = 1
			}
		} else { //No lure then all directions have the same want
			want = 1
		}
		transistion[j] = transistion[j]*float64(want*ability) //This is the score for each direction
	}
	//Scale the scores for each direction
	var sum float64
	stay = 0
	if STAYSTILL {
		stay = 1
		if LUREEFFECT {
			if localcity {
				want = LURESCALE
			} else {
				want = 1
			}
		} else {
			want = 1
		}
		if FENCES {
			ability = UNIMPAIRED
		} else {
			ability = 1
		}
	}
	stay = stay*float64(want*ability)
	sum = stay
	for i := 0; i < 4; i++ {
		sum += transistion[i]
	}
	for i := 0; i < 4; i++ {
		transistion[i] = transistion[i]/sum
	}
	stay = stay/sum
	ratelist = []float64{transistion[0], transistion[1], transistion[2], transistion[3], stay}
	answer = Neighborhood{ratelist}
	return answer

}



func main() {
	data := make([][]cell.Cell, WIDTH)
	previousdata := make([][]cell.Cell,WIDTH)
	var i int
	var j int
	//Initialize data 
	for i := 0; i < WIDTH; i++ {
		data[i] = make([]cell.Cell, HEIGHT)
		previousdata[i] = make([]cell.Cell, HEIGHT)
		for j:= 0; j < HEIGHT; j++ {
			data[i][j] = cell.Cell{false, false, INITJS, INITJI, INITJV, INITAS, INITAI, INITAV}
		}
	}

	//Manual Input of the squares to use
	//Note this part could have been done much more neatly

	// City definitions
	// Modeled as (0,0) top left corner, x first, y second

	data[4][12].City = true
	data[4][11].City = true
	data[5][12].City = true
	data[5][11].City = true
	data[6][12].City = true
	data[6][11].City = true

	data[7][12].City = true
	data[7][11].City = true
	data[7][10].City = true
	data[7][9].City = true
	data[7][8].City = true
	data[7][7].City = true

	data[8][12].City = true
	data[8][11].City = true
	data[8][10].City = true
	data[8][9].City = true
	data[8][8].City = true
	data[8][7].City = true

	data[9][12].City = true
	data[9][11].City = true
	data[9][10].City = true
	data[9][9].City = true
	data[9][8].City = true
	data[9][7].City = true

	data[12][12].City = true
	data[12][11].City = true
	data[10][12].City = true
	data[10][11].City = true
	data[11][12].City = true
	data[11][11].City = true


	// Fence 1 Definition - 1 Tile Perimeter of City
	// Modeled as (0,0) top left corner, x first, y second

	data[3][13].Inside = true
	data[3][12].Inside = true
	data[3][11].Inside = true
	data[3][10].Inside = true

	data[4][13].Inside = true
	data[4][12].Inside = true
	data[4][11].Inside = true
	data[4][10].Inside = true

	data[5][13].Inside = true
	data[5][12].Inside = true
	data[5][11].Inside = true
	data[5][10].Inside = true

	data[6][13].Inside = true
	data[6][12].Inside = true
	data[6][11].Inside = true
	data[6][10].Inside = true
	data[6][9].Inside = true
	data[6][8].Inside = true
	data[6][7].Inside = true
	data[6][6].Inside = true

	data[7][13].Inside = true
	data[7][12].Inside = true
	data[7][11].Inside = true
	data[7][10].Inside = true
	data[7][9].Inside = true
	data[7][8].Inside = true
	data[7][7].Inside = true
	data[7][6].Inside = true

	data[8][13].Inside = true
	data[8][12].Inside = true
	data[8][11].Inside = true
	data[8][10].Inside = true
	data[8][9].Inside = true
	data[8][8].Inside = true
	data[8][7].Inside = true
	data[8][6].Inside = true

	data[9][13].Inside = true
	data[9][12].Inside = true
	data[9][11].Inside = true
	data[9][10].Inside = true
	data[9][9].Inside = true
	data[9][8].Inside = true
	data[9][7].Inside = true
	data[9][6].Inside = true

	data[10][13].Inside = true
	data[10][12].Inside = true
	data[10][11].Inside = true
	data[10][10].Inside = true
	data[10][9].Inside = true
	data[10][8].Inside = true
	data[10][7].Inside = true
	data[10][6].Inside = true

	data[11][13].Inside = true
	data[11][12].Inside = true
	data[11][11].Inside = true
	data[11][10].Inside = true

	data[12][13].Inside = true
	data[12][12].Inside = true
	data[12][11].Inside = true
	data[12][10].Inside = true

	data[13][13].Inside = true
	data[13][12].Inside = true
	data[13][11].Inside = true
	data[13][10].Inside = true


	// Fence 2 Definition - 0 Tile Perimeter of City
	// Modeled as (0,0) top left corner, x first, y second
	/*


	data[4][12].Inside = true
	data[4][11].Inside = true
	data[5][12].Inside = true
	data[5][11].Inside = true
	data[6][12].Inside = true
	data[6][11].Inside = true

	data[7][12].Inside = true
	data[7][11].Inside = true
	data[7][10].Inside = true
	data[7][9].Inside = true
	data[7][8].Inside = true
	data[7][7].Inside = true

	data[8][12].Inside = true
	data[8][11].Inside = true
	data[8][10].Inside = true
	data[8][9].Inside = true
	data[8][8].Inside = true
	data[8][7].Inside = true

	data[9][12].Inside = true
	data[9][11].Inside = true
	data[10][12].Inside = true
	data[10][11].Inside = true
	data[11][12].Inside = true
	data[11][11].Inside = true
	*/

	// Fence 3 Definition - 0 Tile Perimeter of City
	// Modeled as (0,0) top left corner, x first, y second
	/*


	data[3][12].Inside = true
	data[3][10].Inside = true

	data[4][13].Inside = true
	data[4][11].Inside = true

	data[5][12].Inside = true
	data[5][10].Inside = true

	data[6][13].Inside = true
	data[6][11].Inside = true
	data[6][9].Inside = true
	data[6][7].Inside = true

	data[7][12].Inside = true
	data[7][10].Inside = true
	data[7][8].Inside = true
	data[7][6].Inside = true

	data[8][13].Inside = true
	data[8][11].Inside = true
	data[8][9].Inside = true
	data[8][7].Inside = true

	data[9][12].Inside = true
	data[9][10].Inside = true
	data[9][8].Inside = true
	data[9][6].Inside = true

	data[10][13].Inside = true
	data[10][11].Inside = true
	data[10][9].Inside = true
	data[10][7].Inside = true

	data[11][12].Inside = true
	data[11][10].Inside = true

	data[12][13].Inside = true
	data[12][11].Inside = true

	data[13][12].Inside = true
	data[13][10].Inside = true
	*/

	//Populate the movement matrix
	//the neighborhood data structure contains the movement values to the neighbors
	movement := make([][]Neighborhood, WIDTH)
	for i := 0; i < WIDTH; i++ {
		movement[i] = make([]Neighborhood, HEIGHT)
		for j := 0; j < HEIGHT; j++ {
			movement[i][j] = getNeighborhood(i, j, data)
		}
	}

	//Run the expiriment
	for year := 0; year < TIMESTEPS; year++ {
		//Output timestep?
		if year % OUTPUTFREQ == 0 || year == (TIMESTEPS-1) {
			BioPictures.Pics(data, year)
		}


		//Step 1 Juvenille Dispersal
		// This block copies old data to previous data
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				previousdata[i][j] = data[i][j]
			}
		}

		//Set all juvenille for next timestep to zero
		for i = 0; i <WIDTH; i++ {
			for j = 0; j < WIDTH; j++ {
				data[i][j] = cell.Cell{previousdata[i][j].Inside, previousdata[i][j].City, 0, 0, 0, previousdata[i][j].AS, previousdata[i][j].AI, previousdata[i][j].AV}
			}
		}
		//Then insert proper value
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				for index, direction := range DMAP {
					_, valid, _ := checkNeighbor(i, j, previousdata, direction)
					if valid {
						data[i+direction.X][j+direction.Y].JS += previousdata[i][j].JS*movement[i][j].rates[index]
						data[i+direction.X][j+direction.Y].JI += previousdata[i][j].JI*movement[i][j].rates[index]
						data[i+direction.X][j+direction.Y].JV += previousdata[i][j].JV*movement[i][j].rates[index]
					}

				}
				//Stay (note this adds 0 if staystill is false)
				data[i][j].JS += previousdata[i][j].JS*movement[i][j].rates[4]
				data[i][j].JI += previousdata[i][j].JS*movement[i][j].rates[4]
				data[i][j].JV += previousdata[i][j].JS*movement[i][j].rates[4]
			}
		}

		//Step 2 Disease Transmission
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				previousdata[i][j] = data[i][j]
			}
		}

		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				data[i][j].JS = previousdata[i][j].JS*math.Exp(-JUVENILLETRANSFER*(previousdata[i][j].JI+previousdata[i][j].AI))
				data[i][j].JI = previousdata[i][j].JS*(1-math.Exp(-JUVENILLETRANSFER*(previousdata[i][j].JI+previousdata[i][j].AI)))
				data[i][j].AS = previousdata[i][j].AS*math.Exp(-ADULTTRANSFER*(previousdata[i][j].JI+previousdata[i][j].AI))
				data[i][j].AI = previousdata[i][j].AS*(1-math.Exp(-ADULTTRANSFER*(previousdata[i][j].JI+previousdata[i][j].AI)))
			}
		}

		//Step 3 Vaccination
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				previousdata[i][j] = data[i][j]
			}
		}

		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				if data[i][j].Inside || data[i][j].City {
					data[i][j].JS = previousdata[i][j].JS*(1-VACCINATION)
					data[i][j].JI = previousdata[i][j].JI*(1-VACCINATION)
					data[i][j].AS = previousdata[i][j].AS*(1-VACCINATION)
					data[i][j].AI = previousdata[i][j].AI*(1-VACCINATION)

					data[i][j].JV = (previousdata[i][j].JI+previousdata[i][j].JS)*VACCINATION+previousdata[i][j].JV
					data[i][j].AV = (previousdata[i][j].AI+previousdata[i][j].AS)*VACCINATION+previousdata[i][j].AV
				}
			}
		}

		//Step 4 Mortality
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				previousdata[i][j] = data[i][j]
			}
		}

		var sum float64 = 0
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				sum = previousdata[i][j].JS+previousdata[i][j].JI+previousdata[i][j].JV+previousdata[i][j].AS+previousdata[i][j].AI+previousdata[i][j].AV
				data[i][j].JS = previousdata[i][j].JS*math.Exp(-JUVENILLEDEATH*sum)
				data[i][j].JI = previousdata[i][j].JI*math.Exp(-JUVENILLEDEATH*sum)
				data[i][j].JV = previousdata[i][j].JV*math.Exp(-JUVENILLEDEATH*sum)
				data[i][j].AS = previousdata[i][j].AS*math.Exp(-ADULTDEATH*sum)
				data[i][j].AI = previousdata[i][j].AI*math.Exp(-ADULTDEATH*sum)
				data[i][j].AV = previousdata[i][j].AV*math.Exp(-ADULTDEATH*sum)
			}
		}

		//Step 5 Reproduction
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				previousdata[i][j] = data[i][j]
			}
		}

		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				data[i][j].JS = JUVENILLEREPO*(previousdata[i][j].JS+previousdata[i][j].JV)+ADULTREPO*(previousdata[i][j].AS+previousdata[i][j].AV)
				data[i][j].JI = 0
				data[i][j].JV = 0
				data[i][j].AS = previousdata[i][j].JS+previousdata[i][j].AS
				data[i][j].AI = previousdata[i][j].JI+previousdata[i][j].AI
				data[i][j].AV = previousdata[i][j].JV+previousdata[i][j].AV
			}
		}


		//Step 6 Disease Again
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				previousdata[i][j] = data[i][j]
			}
		}

		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				data[i][j].JS = previousdata[i][j].JS*math.Exp(-JUVENILLETRANSFER*(previousdata[i][j].JI+previousdata[i][j].AI))
				data[i][j].JI = previousdata[i][j].JS*(1-math.Exp(-JUVENILLETRANSFER*(previousdata[i][j].JI+previousdata[i][j].AI)))
				data[i][j].AS = previousdata[i][j].AS*math.Exp(-ADULTTRANSFER*(previousdata[i][j].JI+previousdata[i][j].AI))
				data[i][j].AI = previousdata[i][j].AS*(1-math.Exp(-ADULTTRANSFER*(previousdata[i][j].JI+previousdata[i][j].AI)))
			}
		}

		//Step 7 Death Again
		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				previousdata[i][j] = data[i][j]
			}
		}

		for i = 0; i < WIDTH; i++ {
			for j = 0; j < HEIGHT; j++ {
				sum = previousdata[i][j].JS+previousdata[i][j].JI+previousdata[i][j].JV+previousdata[i][j].AS+previousdata[i][j].AI+previousdata[i][j].AV
				data[i][j].JS = previousdata[i][j].JS*math.Exp(-JUVENILLEDEATH*sum)
				data[i][j].JI = previousdata[i][j].JI*math.Exp(-JUVENILLEDEATH*sum)
				data[i][j].JV = previousdata[i][j].JV*math.Exp(-JUVENILLEDEATH*sum)
				data[i][j].AS = previousdata[i][j].AS*math.Exp(-ADULTDEATH*sum)
				data[i][j].AI = previousdata[i][j].AI*math.Exp(-ADULTDEATH*sum)
				data[i][j].AV = previousdata[i][j].AV*math.Exp(-ADULTDEATH*sum)
			}
		}
	}
}
