package imagemap

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Area struct {
	Id  string `json:"id"`
	Alt string `json:"alt,omitempty"`
	Min Point  `json:"min"`
	Max Point  `json:"max"`
}

type XMLArea struct {
	Shape  string `xml:"shape"`
	Coords string `xml:"coords"`
	XRef   struct {
		Href string `xml:"href,attr"`
		Alt  string `xml:",chardata"`
	} `xml:"xref"`
}

type XML struct {
	Image struct {
		Href string `xml:"href,attr"`
	} `xml:"image"`
	Area []XMLArea `xml:"area"`
}

func New(id string, image string, size Point, areas []Area) kb.Item {
	return kb.Item{
		"type":  "image-map",
		"id":    id,
		"image": image,
		"areas": areas,
		"size":  size,
	}
}

func extractID(href string) string {
	href = strings.TrimPrefix(href, "#")
	i := strings.LastIndexByte(href, '/')
	if i < 0 {
		return href
	}
	return href[i+1:]
}

func FromXML(m *XML) (kb.Item, error) {
	var err error
	areas := []Area{}
	for _, area := range m.Area {
		switch area.Shape {
		case "rect":
			tokens := strings.Split(area.Coords, ",")
			if len(tokens) != 4 {
				return nil, errors.New("invalid image-map coords \"" + area.Coords + "\"")
			}

			x0, err0 := strconv.Atoi(tokens[0])
			y0, err1 := strconv.Atoi(tokens[1])
			x1, err2 := strconv.Atoi(tokens[2])
			y1, err3 := strconv.Atoi(tokens[3])

			if err0 != nil || err1 != nil || err2 != nil || err3 != nil {
				return nil, errors.New("invalid image-map coords \"" + area.Coords + "\"")
			}

			areas = append(areas, Area{
				Id:  extractID(area.XRef.Href),
				Alt: area.XRef.Alt,
				Min: Point{x0, y0},
				Max: Point{x1, y1},
			})
		default:
			return nil, errors.New("unhandled image-map shape \"" + area.Shape + "\"")
		}
	}

	const marker = ";base64,"

	i := strings.Index(m.Image.Href, marker)
	if i < 0 {
		return New("", m.Image.Href, Point{0, 0}, areas), err
	}

	rd := base64.NewDecoder(base64.StdEncoding,
		strings.NewReader(m.Image.Href[i+len(marker):]))
	img, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}
	size := Point{
		X: img.Bounds().Dx(),
		Y: img.Bounds().Dy(),
	}
	return New("", m.Image.Href, size, areas), err
}
