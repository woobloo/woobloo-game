package tile

const (
	TileWidth  = uint(32)
	TileHeight = uint(32)
)

type Tile struct {
	X, Y uint
}

func NewDefaultMap(width uint, height uint) [][]*Tile {
	tiles := [][]*Tile{}
	for h := uint(0); h < height; h++ {
		row := make([]*Tile, width)
		for w := uint(0); w < width; w++ {
			row[w] = &Tile{X: w * TileWidth, Y: h * TileHeight}
		}
		tiles = append(tiles, row)
	}
	return tiles
}
