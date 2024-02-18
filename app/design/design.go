package design

type Color string
type Texture string
type Pattern string

const (
	Citron Color = "citron"
	Peach  Color = "peach"
	Gray   Color = "gray"
)

const (
	NoiseLight Texture = "noise-light.png"
	NoiseDark  Texture = "noise-dark.png"
)

const (
	Bubbles Pattern = "bubbles.svg"
)

type CardDesign struct {
	Color          Color   `json:"color"`
	Texture        Texture `json:"texture"`
	TextureOpacity float64 `json:"textureOpacity"`
	Pattern        Pattern `json:"pattern"`
}

func DefaultCardDesign() CardDesign {
	return CardDesign{
		Color:          Gray,
		Texture:        NoiseLight,
		TextureOpacity: 0.1,
		Pattern:        "",
	}
}
