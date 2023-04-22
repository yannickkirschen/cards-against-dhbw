package player

type Player struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
}

type PublicPlayer struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	IsMod  bool   `json:"isMod"`
	IsBoss bool   `json:"isBoss"`
	Points int    `json:"points"`
}

func New(name string) *Player {
	return &Player{
		Name:   name,
		Points: 0,
	}
}

func (p *Player) ToPublic(mod *Player, boss *Player) *PublicPlayer {
	return &PublicPlayer{
		Id:     p.Name,
		Name:   p.Name,
		IsMod:  p == mod,
		IsBoss: p == boss,
		Points: p.Points,
	}
}
