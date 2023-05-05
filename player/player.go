package player

type Player struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Points int    `json:"points"`
}

type PublicPlayer struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	IsMod  bool   `json:"isMod"`
	IsBoss bool   `json:"isBoss"`
	Points int    `json:"points"`
}

func New(name string) *Player {
	return &Player{
		Name:   name,
		Active: true,
		Points: 0,
	}
}

func (p *Player) ToPublic(mod string, boss string) *PublicPlayer {
	return &PublicPlayer{
		Id:     p.Name,
		Name:   p.Name,
		Active: p.Active,
		IsMod:  p.Name == mod,
		IsBoss: p.Name == boss,
		Points: p.Points,
	}
}
