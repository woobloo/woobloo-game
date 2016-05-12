package tile

const (
	TileWidth  = uint(32)
	TileHeight = uint(32)
)

type Tile struct {
	X, Y    uint
	Terrain TerrainType
	Foliage FoliageType
	Mineral MineralType
}

func NewDefaultMap(width uint, height uint) [][]*Tile {
	tiles := [][]*Tile{}
	for h := uint(0); h < height; h++ {
		row := make([]*Tile, width)
		for w := uint(0); w < width; w++ {
			row[w] = &Tile{
				X:       w * TileWidth,
				Y:       h * TileHeight,
				Terrain: Plains,
				Foliage: Grassy,
				Mineral: None}
		}
		tiles = append(tiles, row)
	}
	return tiles
}

type TerrainType byte

const (
	Plains TerrainType = iota
	Hilly
	Mountainous
	Coastal
)

type FoliageType byte

const (
	Desert FoliageType = iota
	Grassy
	Forest
	Jungle
)

type MineralType byte

const (
	None MineralType = iota
	Generic
	Iron
	Silver
	Gold
	Gemstones
)
