package gocalc

import (
	"fmt"
	// "math"
	"strings"

	"github.com/blldr/gosolve"
)

type (
	Graph struct {
		width      int
		points     []EqPoints
		resolution int
	}
	Point struct {
		X float64
		Y float64
	}
	EqPoints struct {
		EqType EqType
		Color  string
		Points []Point
	}
	EqType int
)

const (
	equality EqType = iota
	greater
	smaller
)

func NewGraph(width int, resolution int) *Graph {
	return &Graph{
		width,
		make([]EqPoints, 0, 10),
		resolution,
	}
}

func (g *Graph) AddEquation(eqString string, iStart, iEnd float64, EqType EqType, color string) (int, error) {
	pointsSlice, err := g.createPointsSlice(eqString, iStart, iEnd)
	if err != nil {
		return 0, err
	}
	g.points = append(g.points, EqPoints{EqType, color, pointsSlice})
	return len(g.points) - 1, nil
}

func (g *Graph) RemoveEquation(eqId int) {
	if eqId >= len(g.points) {
		return
	}
	g.points = append(g.points[:eqId], g.points[eqId+1:]...)
}

func (g *Graph) ChangeEquation(eqId int, eqString string, iStart, iEnd float64, eqType EqType, color string) error {
	if eqId >= len(g.points) {
		return nil
	}
	pointsSlice, err := g.createPointsSlice(eqString, iStart, iEnd)
	if err != nil {
		return err
	}
	g.points[eqId] = EqPoints{eqType, color, pointsSlice}
	return nil
}

func (g *Graph) GetPoints() []EqPoints {
	return g.points
}

func (g *Graph) createPointsSlice(eqString string, iStart, iEnd float64) ([]Point, error) {
	xCount := g.width * g.resolution
	pointsSlice := make([]Point, 0, xCount)
	if iStart > iEnd {
		iStart, iEnd = iEnd, iStart
	}
	if iStart < float64(-g.width/2) {
		iStart = float64(-g.width / 2)
	}
	if iEnd > float64(g.width/2) {
		iEnd = float64(g.width / 2)
	}
	for i := 0.0; i < float64(xCount); i++ {
		x := i/float64(g.resolution) - float64(g.width/2)
		// x = math.Round(x*1000) / 1000
		if x < iStart || x > iEnd {
			continue
		}
		var eq string
		eq = strings.ReplaceAll(eqString, "x", fmt.Sprintf("(%f)", x))

		var yValues []float64
		var err error
		yValues, err = gosolve.FindRoots(eq, iStart, iEnd, 'y')

		if err != nil {
			return nil, err
		}
		for j := 0; j < len(yValues); j++ {
			pointsSlice = append(pointsSlice, Point{x, yValues[j]})
		}
	}
	for i := 0.0; i < float64(xCount); i++ {
		y := float64(g.width/2) - i/float64(g.resolution)
		// y = math.Round(y*100) / 100
		var eq string
		eq = strings.ReplaceAll(eqString, "y", fmt.Sprintf("(%f)", y))
		var xValues []float64
		var err error
		xValues, err = gosolve.FindRoots(eq, iStart, iEnd, 'x')

		if err != nil {
			return nil, err
		}
		for j := 0; j < len(xValues); j++ {
			pointsSlice = append(pointsSlice, Point{xValues[j], y})

		}

	}
	return pointsSlice, nil
}
