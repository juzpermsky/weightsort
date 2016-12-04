package main

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"math/rand"

	"fmt"
	"math"
)

type XY struct{ X, Y float64 }
type RAlpha struct{ R, Alpha float64 }

type PolarNode struct {
	From        float64
	To          float64
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
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "points.png"); err != nil {
		panic(err)
	}

	rootNode = &PolarNode{0, 90, RAlpha{}, 0, nil, nil}
	createPolarTree(points)
	printPolarTree(rootNode)
}

func XY2RAlpha(point XY) RAlpha {
	var r float64 = math.Sqrt(point.X*point.X + point.Y*point.Y)
	var alpha float64 = math.Acos(point.X/r) * 180 / math.Pi
	return RAlpha{R: r, Alpha: alpha}
}

func printPolarTree(curNode *PolarNode) {
	fmt.Printf("from %v to %v: R=%v, Alpha=%v, count=%v\n", curNode.From, curNode.To, curNode.MaxPoint.R, curNode.MaxPoint.Alpha, curNode.PointsCount)
	if curNode.Left != nil {
		printPolarTree(curNode.Left)
	}
	if curNode.Right != nil {
		printPolarTree(curNode.Right)
	}
}

func createPolarTree(points plotter.XYs) {
	for _, point := range points {
		rAlpha := XY2RAlpha(point)
		add2PolarTree(rAlpha, rootNode)
	}
}

func add2PolarTree(rAlpha RAlpha, curNode *PolarNode) {
	curNode.PointsCount++
	if rAlpha.R > curNode.MaxPoint.R {
		oldMaxPoint := curNode.MaxPoint
		curNode.MaxPoint = rAlpha
		fmt.Printf("%v added in from %v to %v\n",rAlpha,curNode.From, curNode.To)
		if rAlpha.Alpha < (curNode.To-curNode.From)/2 {
			if curNode.Left == nil {
				curNode.Left = &PolarNode{curNode.From, (curNode.To - curNode.From) / 2, rAlpha, 1, nil, nil}
				fmt.Printf("%v added in from %v to %v\n",rAlpha,curNode.From, (curNode.To - curNode.From) / 2)
			} else {
				add2PolarTree(rAlpha, curNode.Left)
			}
		} else {
			if curNode.Right == nil {
				curNode.Right = &PolarNode{(curNode.To - curNode.From) / 2, curNode.To, rAlpha, 1, nil, nil}
				fmt.Printf("%v added in from %v to %v\n",rAlpha,(curNode.To - curNode.From) / 2, curNode.To)
			} else {
				add2PolarTree(rAlpha, curNode.Right)
			}
		}

		if oldMaxPoint.R != 0 {
			fmt.Printf("%v removed in from %v to %v\n",oldMaxPoint,curNode.From, curNode.To)
			if oldMaxPoint.Alpha < (curNode.To-curNode.From)/2 {
				if curNode.Left == nil {
					curNode.Left = &PolarNode{curNode.From, (curNode.To - curNode.From) / 2, oldMaxPoint, 1, nil, nil}
				} else {
					add2PolarTree(oldMaxPoint, curNode.Left)
				}
			} else {
				if curNode.Right == nil {
					curNode.Right = &PolarNode{(curNode.To - curNode.From) / 2, curNode.To, oldMaxPoint, 1, nil, nil}
				} else {
					add2PolarTree(oldMaxPoint, curNode.Right)
				}
			}
		}
	} else {
		if rAlpha.Alpha < (curNode.To-curNode.From)/2 {
			if curNode.Left == nil {
				curNode.Left = &PolarNode{curNode.From, (curNode.To - curNode.From) / 2, rAlpha, 1, nil, nil}
				fmt.Printf("%v added in from %v to %v\n",rAlpha,curNode.From, (curNode.To - curNode.From) / 2)
			} else {
				add2PolarTree(rAlpha, curNode.Left)
			}
		} else {
			if curNode.Right == nil {
				curNode.Right = &PolarNode{(curNode.To - curNode.From) / 2, curNode.To, rAlpha, 1, nil, nil}
				fmt.Printf("%v added in from %v to %v\n",rAlpha,(curNode.To - curNode.From) / 2, curNode.To)
			} else {
				add2PolarTree(rAlpha, curNode.Right)
			}
		}
	}
	//	fmt.Println(curNode)
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
