package main

import (
	"log"
	"regexp"

	"github.com/CGA1123/aoc"
)

var TileExp = regexp.MustCompile(`Tile (?P<id>\d+):`)

type Tile struct {
	id    int64
	grid  [][]byte
	edges []string
}

func Edges(t *Tile) []string {
	l := len(t.grid)

	var top, left, right, bottom []byte

	for i := 0; i < l; i++ {
		top = append(top, t.grid[0][i])
		right = append(right, t.grid[i][l-1])

		bottom = append(bottom, t.grid[l-1][i])
		left = append(left, t.grid[i][0])
	}

	return []string{
		string(top),
		string(right),
		string(bottom),
		string(left)}
}

func Flip(t *Tile) *Tile {
	var flip [][]byte

	for i := len(t.grid) - 1; i >= 0; i-- {
		flip = append(flip, t.grid[i])
	}

	r := &Tile{id: t.id, grid: flip}

	r.edges = Edges(r)

	return r
}

func Rotate(t *Tile) *Tile {
	var rotated [][]byte

	for x := 0; x < len(t.grid); x++ {
		var row []byte

		for y := len(t.grid) - 1; y >= 0; y-- {
			row = append(row, t.grid[y][x])
		}

		rotated = append(rotated, row)
	}

	r := &Tile{id: t.id, grid: rotated}
	r.edges = Edges(r)

	return r
}

func Alternatives(t *Tile) []*Tile {
	rt, rf := t, Flip(t)

	tiles := []*Tile{}

	for i := 0; i < 4; i++ {
		rt, rf = Rotate(rt), Rotate(rf)
		tiles = append(tiles, rt, rf)
	}

	return tiles
}

type Image struct {
	current                  *Tile
	left, right, top, bottom *Image
}

func (i *Image) SetDir(direction int, image *Image) {
	switch direction {
	case 0:
		i.top = image
		image.bottom = i
	case 1:
		i.right = image
		image.left = i
	case 2:
		i.bottom = image
		image.top = i
	case 3:
		i.left = image
		image.right = i
	}
}

func (i *Image) Neighbours() []*Image {
	var n []*Image

	for _, i := range []*Image{i.top, i.right, i.bottom, i.left} {
		if i == nil {
			continue
		}

		n = append(n, i)
	}

	return n
}

func Top(i *Image) *Image {
	c := i
	for {
		if c.top == nil {
			return i
		}

		c = c.top
	}
}

func Bottom(i *Image) *Image {
	c := i

	for {
		if c.bottom == nil {
			return c
		}

		c = c.bottom
	}
}

func Right(i *Image) *Image {
	c := i
	for {
		if c.right == nil {
			return c
		}

		c = c.right
	}
}

func Left(i *Image) *Image {
	c := i

	for {
		if c.left == nil {
			return c
		}

		c = c.left
	}
}

func ConstructImage(img *Image, imgs map[int64]*Image, visited *aoc.Set, edges map[int]map[string]*aoc.Set) {
	if visited.Contains(img.current.id) {
		return
	}

	for dir, edge := range img.current.edges {
		opDir := (dir + 2) % 4

		for _, candidate := range edges[opDir][edge].Elements() {
			tile := candidate.(*Tile)

			if tile.id == img.current.id {
				continue
			}

			newImage, ok := imgs[tile.id]
			if !ok {
				newImage = &Image{current: tile}
				imgs[tile.id] = newImage
			}

			img.SetDir(dir, newImage)
		}
	}

	visited.Add(img.current.id)

	for _, neighbour := range img.Neighbours() {
		ConstructImage(neighbour, imgs, visited, edges)
	}
}

func PartOne(i *Image) int64 {
	tl := Top(Left(i)).current.id
	bl := Bottom(Left(i)).current.id
	tr := Top(Right(i)).current.id
	br := Bottom(Right(i)).current.id

	return tl * bl * tr * br
}

func BuildGrid(i *Image) [][]byte {
	left := Top(Left(i))
	current := left

	var grid [][]byte
	var dy, dx int

	for i := 0; i < 96; i++ {
		grid = append(grid, make([]byte, 96))
	}

	for current != nil {
		right := current

		for right != nil {
			for y := 1; y < len(right.current.grid)-1; y++ {
				for x := 1; x < len(right.current.grid)-1; x++ {
					gridX := (x - 1) + (8 * dx)
					gridY := (y - 1) + (8 * dy)

					grid[gridY][gridX] = right.current.grid[y][x]
				}
			}

			right = right.right
			dx++
		}

		left = left.bottom
		current = left
		dy++
		dx = 0
	}

	return grid
}

func ParseMonster() *aoc.Grid {
	var y int64
	grid := aoc.NewGrid()

	aoc.EachLine("monster.txt", func(l string) {
		for x, b := range []byte(l) {
			if b == ' ' {
				continue
			}

			grid.Write(int64(x), y, b)
		}

		y++
	})

	return grid
}

func IsMonster(grid [][]byte, x, y int64, monster *aoc.Grid) bool {
	equal := true

	monster.EachSparse(func(point aoc.Point, i interface{}) {
		xi, yi := x+point.X, y+point.Y

		equal = equal && (grid[yi][xi] == i.(byte))
	})

	return equal
}

func CountMonsters(grid [][]byte, monster *aoc.Grid) int64 {
	var c int64

	for y := int64(0); y < int64(len(grid))-monster.Height(); y++ {
		for x := int64(0); x < (int64(len(grid)) - monster.Width()); x++ {
			if IsMonster(grid, x, y, monster) {
				c++
			}

		}
	}

	return c
}

func PartTwo(i *Image) int64 {
	grids := Alternatives(&Tile{grid: BuildGrid(i)})
	monster := ParseMonster()

	tiles := int64(0)

	for y := 0; y < len(grids[0].grid); y++ {
		for x := 0; x < len(grids[0].grid); x++ {
			if grids[0].grid[y][x] == '#' {
				tiles++
			}
		}
	}

	for _, grid := range grids {
		c := CountMonsters(grid.grid, monster)
		if c == 0 {
			continue
		}

		return tiles - (c * monster.Count())
	}

	return -1
}

func main() {
	var tile *Tile
	edges := map[int]map[string]*aoc.Set{0: {}, 1: {}, 2: {}, 3: {}}
	var id int64
	var grid [][]byte
	var i int

	aoc.EachLine("input.txt", func(line string) {
		if line == "" {
			i = 0
			return
		}

		i++

		if i == 1 {
			capture := aoc.Capture(TileExp, line)
			id = aoc.MustParse(capture["id"])
			return
		}

		grid = append(grid, []byte(line))

		if len(grid) == 10 {
			t := &Tile{id: id, grid: grid}
			t.edges = Edges(t)

			for _, tile := range Alternatives(t) {
				for i, edge := range tile.edges {
					if _, ok := edges[i][edge]; !ok {
						edges[i][edge] = aoc.NewSet()
					}

					edges[i][edge].Add(tile)
				}
			}

			grid = nil
			tile = t
		}
	})

	visited := aoc.NewSet()
	image := &Image{current: tile}
	images := map[int64]*Image{tile.id: image}

	ConstructImage(image, images, visited, edges)

	log.Printf("pt(1): %v", PartOne(image))
	log.Printf("pt(2): %v", PartTwo(image))
}
