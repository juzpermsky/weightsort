package main

import (
	"fmt"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"math"
	"math/rand"
	"time"
)

type XY struct{ X, Y float64 }
type RAlpha struct{ R, Alpha float64 }

type PolarNode struct {
	From        float64
	To          float64
	Split       float64
	MaxPoint    RAlpha
	PointsCount int64
	Left        *PolarNode
	Right       *PolarNode
}

var rootNode *PolarNode

func main() {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Weight sort test"

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	points := randomPoints(10)

	err = plotutil.AddScatters(p, points)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(6 * vg.Inch, 6 * vg.Inch, "points.png"); err != nil {
		panic(err)
	}

	rootNode = &PolarNode{0, 90, 45, RAlpha{}, 0, nil, nil}
	start := time.Now()
	createPolarTree(points)
	elapsed := time.Since(start)
	fmt.Printf("Tree creation took %s\n", elapsed)
	printPolarTree(rootNode, 3)

}

func XY2RAlpha(point XY) RAlpha {
	var r float64 = math.Sqrt(point.X * point.X + point.Y * point.Y)
	var alpha float64 = math.Acos(point.X / r) * 180 / math.Pi
	return RAlpha{R: r, Alpha: alpha}
}

func printPolarTree(curNode *PolarNode, level int) {
	fmt.Printf("from %v to %v: R=%v, Alpha=%v, count=%v\n", curNode.From, curNode.To, curNode.MaxPoint.R, curNode.MaxPoint.Alpha, curNode.PointsCount)
	if (level > 0) || (level < 0) {
		if curNode.Left != nil {
			printPolarTree(curNode.Left, level - 1)
		}
		if curNode.Right != nil {
			printPolarTree(curNode.Right, level - 1)
		}
	}
}

func createPolarTree(points plotter.XYs) {
	for _, point := range points {
		rAlpha := XY2RAlpha(point)
		//fmt.Println(rAlpha)
		add2PolarTree(rAlpha, rootNode)
	}
}

func moveDown(rAlpha RAlpha, curNode *PolarNode) {
	if rAlpha.Alpha < curNode.Split {
		if curNode.Left == nil {
			curNode.Left = &PolarNode{curNode.From, curNode.Split, rAlpha.Alpha, rAlpha, 1, nil, nil}
		} else {
			add2PolarTree(rAlpha, curNode.Left)
		}
	} else {
		if curNode.Right == nil {
			curNode.Right = &PolarNode{curNode.Split, curNode.To, rAlpha.Alpha, rAlpha, 1, nil, nil}
		} else {
			add2PolarTree(rAlpha, curNode.Right)
		}
	}
}

func add2PolarTree(rAlpha RAlpha, curNode *PolarNode) {

	curNode.PointsCount++
	oldMaxPoint := curNode.MaxPoint
	if rAlpha.R > curNode.MaxPoint.R {
		curNode.MaxPoint = rAlpha
		if oldMaxPoint.R != 0 {
			moveDown(oldMaxPoint, curNode)
		}
	}
	if curNode.PointsCount > 1 {
		moveDown(rAlpha, curNode)
	}
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		pts[i].X = rand.Float64() * 5
		pts[i].Y = rand.Float64() * 5
	}
	return pts
}
