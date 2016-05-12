package player

type Player struct {
	Username string
}

func NewDefaultPlayer() *Player {
	return &Player{Username: "Default"}
}
