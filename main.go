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
	MaxPointXY  XY
	PointsCount int64
	Left        *PolarNode
	Right       *PolarNode
}

type NodeQueue []*PolarNode

func (h NodeQueue) Len() int {
	return len(h)
}

func (h *NodeQueue) Push(x *PolarNode) {
	*h = append(*h, x)
}

func (h *NodeQueue) Pop() *PolarNode {
	old := *h
	x := old[0]
	*h = old[1 :]
	return x
}

var rootNode *PolarNode
var curveLen int = 200
var maxR float64
var maxLine plotter.XYs

func main() {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Weight sort test"

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	points := randomPoints(2000)

	rootNode = &PolarNode{0, 90, 45, RAlpha{}, XY{}, 0, nil, nil}
	start := time.Now()
	createPolarTree(points)
	elapsed := time.Since(start)
	fmt.Printf("Tree creation took %s\n", elapsed)
	maxLine = make(plotter.XYs, curveLen)
	//maxLine[0] = rootNode.MaxPointXY
	//maxR = rootNode.MaxPoint.R
	//printPolarTree(rootNode, 3)

	breathFirstTraverse()

	err = plotutil.AddScatters(p, points)
	if err != nil {
		panic(err)
	}

	err = plotutil.AddLines(p, maxLine)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "points.png"); err != nil {
		panic(err)
	}

}

func XY2RAlpha(point XY) RAlpha {
	var r float64 = math.Sqrt(point.X*point.X + point.Y*point.Y)
	var alpha float64 = math.Acos(point.X/r) * 180 / math.Pi
	return RAlpha{R: r, Alpha: alpha}
}

func printPolarTree(curNode *PolarNode, level int) {
	//if index < 499 {
	//	//if curNode.MaxPoint.R < maxR {
	//	//maxR = curNode.MaxPoint.R
	//	index++
	//	maxLine[index] = curNode.MaxPointXY
	//	//}
	//}
	fmt.Printf("from %v to %v: R=%v, Alpha=%v, count=%v\n", curNode.From, curNode.To, curNode.MaxPoint.R, curNode.MaxPoint.Alpha, curNode.PointsCount)
	if (level > 0) || (level < 0) {
		if curNode.Left != nil {
			printPolarTree(curNode.Left, level-1)
		}
		if curNode.Right != nil {
			printPolarTree(curNode.Right, level-1)
		}
	}
}

func createPolarTree(points plotter.XYs) {
	for _, point := range points {
		rAlpha := XY2RAlpha(point)
		add2PolarTree(rAlpha, point, rootNode)
	}
}

func moveDown(rAlpha RAlpha, xy XY, curNode *PolarNode) {
	if rAlpha.Alpha < curNode.Split {
		if curNode.Left == nil {
			curNode.Left = &PolarNode{curNode.From, curNode.Split, rAlpha.Alpha, rAlpha, xy, 1, nil, nil}
		} else {
			add2PolarTree(rAlpha, xy, curNode.Left)
		}
	} else {
		if curNode.Right == nil {
			curNode.Right = &PolarNode{curNode.Split, curNode.To, rAlpha.Alpha, rAlpha, xy, 1, nil, nil}
		} else {
			add2PolarTree(rAlpha, xy, curNode.Right)
		}
	}
}

func add2PolarTree(rAlpha RAlpha, xy XY, curNode *PolarNode) {
	if curNode.PointsCount == 0 {
		curNode.PointsCount = 1
		curNode.Split = rAlpha.Alpha
		curNode.MaxPoint = rAlpha
		curNode.MaxPointXY = xy
	} else {
		curNode.PointsCount++
		oldMaxPoint := curNode.MaxPoint
		oldMaxPointXY := curNode.MaxPointXY
		if rAlpha.R > curNode.MaxPoint.R {
			curNode.MaxPoint = rAlpha
			curNode.MaxPointXY = xy
		}
		if curNode.PointsCount == 2 {
			//curNode.Split = 0.5 * (oldMaxPoint.Alpha + rAlpha.Alpha)
			curNode.Split = 0.5 * (curNode.From + curNode.To)
			moveDown(oldMaxPoint, oldMaxPointXY, curNode)
		}
		if curNode.PointsCount >= 2 {
			moveDown(rAlpha, xy, curNode)
		}
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

// обход дерева "в ширину"
func breathFirstTraverse() {
	queue := &NodeQueue{rootNode}
	index := 0
	for (queue.Len() > 0) && (index < curveLen) {
		curNode := queue.Pop()
		// делаем с узлом что нужно
		//fmt.Println(curNode.MaxPoint)
		//fmt.Printf("from %v to %v: R=%v, Alpha=%v, count=%v\n", curNode.From, curNode.To, curNode.MaxPoint.R, curNode.MaxPoint.Alpha, curNode.PointsCount)
		maxLine[index] = curNode.MaxPointXY
		index++
		// и втыкаем его детей в конец очереди
		if curNode.Left != nil {
			queue.Push(curNode.Left)
			//fmt.Printf("left %v\n",queue.Len())
			//fmt.Println(*queue)
		} else {
			//fmt.Println("left nil")
		}

		if curNode.Right != nil {
			queue.Push(curNode.Right)
			//fmt.Printf("right %v\n",queue.Len())
			//fmt.Println(*queue)
		} else {
			//fmt.Println("right nil")
		}
	}
}
