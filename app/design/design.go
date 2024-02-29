package design

import (
	"encoding/json"
	"fmt"
	"os"
)

type Font string
type Color string
type Texture string
type Pattern string

const (
	SystemFont  Font = ""
	CrimsonText Font = "font-crimson"
	Pacifico    Font = "font-pacifico"
)

const (
	Citron Color = "citron"
	Peach  Color = "peach"
	Gray   Color = "gray"
)

const (
	Noise     Texture = "noise.png"
	NoTexture Texture = ""
)

const (
	Bubbles   Pattern = "bubbles.svg"
	NoPattern Pattern = ""
)

type CardDesign struct {
	Color   Color   `json:"color"`
	Font    Font    `json:"font"`
	Texture Texture `json:"texture"`
	/* 0-100 */
	TextureOpacity int     `json:"textureOpacity"`
	Pattern        Pattern `json:"pattern"`
}

func (cd CardDesign) ToJSON() string {
	b, err := json.Marshal(cd)
	if err != nil {
		fmt.Println("Error marshalling card design:", err)
	}
	return string(b)
}

func DefaultCardDesign() CardDesign {
	return CardDesign{
		Color:          Gray,
		Font:           SystemFont,
		Texture:        NoTexture,
		TextureOpacity: 0,
		Pattern:        NoPattern,
	}
}

func ParseCardDesign(design string) (CardDesign, error) {
	if design == "" {
		return DefaultCardDesign(), nil
	}

	var cd CardDesign
	err := json.Unmarshal([]byte(design), &cd)
	if cd.Color == "" {
		cd.Color = Gray
	}
	return cd, err
}

type DesignOption map[string]string
type DesignSetting map[string]DesignOption
type DesignConfig map[string]DesignSetting

func ReadDesignConfig() (DesignConfig, error) {
	f, err := os.ReadFile("./design.settings.json")
	if err != nil {
		return nil, err
	}
	var conf DesignConfig
	err = json.Unmarshal([]byte(f), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
