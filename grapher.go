package gocalc

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/blldr/gosolve"
)

type (
	Graph struct {
		width      int
		points     [][][]float64
		resolution int
	}
)

func NewGraph(width int, resolution int) *Graph {
	return &Graph{
		width,
		make([][][]float64, 0, 10),
		resolution,
	}
}

func (g *Graph) AddEquation(eqString string, iStart, iEnd float64, prefixData [][]float64) (int, error) {
	pointsSlice, err := g.createPointsSlice(eqString, iStart, iEnd)
	if err != nil {
		return 0, err
	}
	pointsSlice = append(prefixData, pointsSlice...)
	g.points = append(g.points, pointsSlice)
	return len(g.points) - 1, nil
}

func (g *Graph) RemoveEquation(eqId int) {
	if eqId >= len(g.points) {
		return
	}
	g.points = append(g.points[:eqId], g.points[eqId+1:]...)
}

func (g *Graph) ChangeEquation(eqId int, eqString string, iStart, iEnd float64, prefixData [][]float64) error {
	if eqId >= len(g.points) {
		return nil
	}
	pointsSlice, err := g.createPointsSlice(eqString, iStart, iEnd)
	if err != nil {
		return err
	}
	pointsSlice = append(prefixData, pointsSlice...)
	g.points[eqId] = pointsSlice
	return nil
}

func (g *Graph) GetPointsSlice() [][][]float64 {
	return g.points
}

func (g *Graph) createPointsSlice(eqString string, iStart, iEnd float64) ([][]float64, error) {
	xCount := g.width * g.resolution
	pointsSlice := make([][]float64, 0, xCount)
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
		x = math.Round(x*100) / 100
		if x < iStart || x > iEnd {
			pointsSlice = append(pointsSlice, make([]float64, 0))
			continue
		}
		var eq string
		xFloat := strconv.FormatFloat(x, 'f', 2, 64)
		eq = strings.ReplaceAll(eqString, "x", fmt.Sprintf("(%s)", xFloat))

		yValues, err := gosolve.FindRoots(eq, -float64(g.width)/2, float64(g.width)/2, 'y')
		if err != nil {
			return nil, err
		}
		pointsSlice = append(pointsSlice, yValues)
	}
	for i := 0.0; i < float64(xCount); i++ {
		y := float64(g.width/2) - i/float64(g.resolution)
		y = math.Round(y*100) / 100
		var eq string
		yFloat := strconv.FormatFloat(y, 'f', 2, 64)
		eq = strings.ReplaceAll(eqString, "y", fmt.Sprintf("(%s)", yFloat))
		xValues, err := gosolve.FindRoots(eq, iStart, iEnd, 'x')
		if err != nil {
			return nil, err
		}
		pointsSlice = append(pointsSlice, xValues)

	}
	return pointsSlice, nil
}
