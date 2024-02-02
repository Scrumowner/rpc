package logic

import (
	"fmt"
	"hugoproxy-main/proxy/worker-content"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type AVLNode struct {
	Value  int
	Left   *AVLNode
	Right  *AVLNode
	Height int
}

type Node struct {
	ID    int
	Name  string
	Form  string
	Links []*Node
}

func WorkerTest() {
	t := time.NewTicker(1 * time.Hour)
	var b byte = 0
	for {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		select {
		case <-t.C:
			err := os.WriteFile("/app/worker-content/tasks/_index.md", []byte(fmt.Sprintf(worker_content.Content, currentTime, b)), 0644)
			if err != nil {
				log.Println(err)
			}
			b++
		}
	}
}

func height(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getBalance(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

// rotateRight выполняет правый поворот узла.
func rotateRight(y *AVLNode) *AVLNode {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	y.Height = max(height(y.Left), height(y.Right)) + 1
	x.Height = max(height(x.Left), height(x.Right)) + 1

	return x
}

func rotateLeft(x *AVLNode) *AVLNode {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	x.Height = max(height(x.Left), height(x.Right)) + 1
	y.Height = max(height(y.Left), height(y.Right)) + 1

	return y
}

func insert(root *AVLNode, value int) *AVLNode {
	if root == nil {
		return &AVLNode{Value: value, Height: 1}
	}

	if value < root.Value {
		root.Left = insert(root.Left, value)
	} else if value > root.Value {
		root.Right = insert(root.Right, value)
	} else {
		return root
	}

	root.Height = 1 + max(height(root.Left), height(root.Right))

	balance := getBalance(root)

	if balance > 1 {
		if value < root.Left.Value {
			return rotateRight(root)
		} else if value > root.Left.Value {
			root.Left = rotateLeft(root.Left)
			return rotateRight(root)
		}
	}

	if balance < -1 {
		if value > root.Right.Value {
			return rotateLeft(root)
		} else if value < root.Right.Value {
			root.Right = rotateRight(root.Right)
			return rotateLeft(root)
		}
	}

	return root
}

func generateMermaidCode(root *AVLNode) string {
	if root == nil {
		return ""
	}

	var mermaidCode strings.Builder
	mermaidCode.WriteString("graph TD;\n")
	traverseTreeForMermaid(root, &mermaidCode)

	return mermaidCode.String()
}

func traverseTreeForMermaid(node *AVLNode, mermaidCode *strings.Builder) {
	if node != nil {
		if node.Left != nil {
			mermaidCode.WriteString(fmt.Sprintf("  %d --> %d;\n", node.Value, node.Left.Value))
		}
		if node.Right != nil {
			mermaidCode.WriteString(fmt.Sprintf("  %d --> %d;\n", node.Value, node.Right.Value))
		}
		traverseTreeForMermaid(node.Left, mermaidCode)
		traverseTreeForMermaid(node.Right, mermaidCode)
	}
}

func BinaryTreeWorker() {
	t := time.NewTicker(1 * time.Hour)
	var treeRoot *AVLNode
	var b int

	for {
		select {
		case <-t.C:

			newValue := int(rand.Intn(100))
			treeRoot = insert(treeRoot, newValue)

			if b%100 == 0 {
				treeRoot = nil
				for i := 0; i < 5; i++ {
					newValue := rand.Intn(100)
					treeRoot = insert(treeRoot, newValue)
				}
			}

			mermaidCode := generateMermaidCode(treeRoot)

			updatedTree := fmt.Sprintf(worker_content.Tree, mermaidCode)

			err := os.WriteFile("/app/worker-content/tasks/binary.md", []byte(updatedTree), 0664)
			if err != nil {
				log.Println(err)
			}
			b++
		}
	}
}

func generateGraph() string {
	var nodesBuilder strings.Builder

	nodesCount := rand.Intn(26) + 5

	// Генерация узлов
	for i := 1; i <= nodesCount; i++ {
		nodeID := i
		nodeName := fmt.Sprintf("Node%d", nodeID)
		nodeForm := getRandomForm()
		nodesBuilder.WriteString(fmt.Sprintf("    %s[%s %s]\n", nodeName, nodeForm, getRandomText()))
	}

	for i := 1; i <= nodesCount; i++ {
		fromNode := fmt.Sprintf("Node%d", i)
		toNode := fmt.Sprintf("Node%d", rand.Intn(nodesCount)+1)
		nodesBuilder.WriteString(fmt.Sprintf("    %s --> %s\n", fromNode, toNode))
	}

	return nodesBuilder.String()
}

func getRandomForm() string {
	forms := []string{"circle", "rect", "square", "ellipse", "round-rect", "rhombus"}
	return forms[rand.Intn(len(forms))]
}

func getRandomText() string {
	texts := []string{"Square", "Circle", "Round Rect", "Rhombus", "Ellipse"}
	return texts[rand.Intn(len(texts))]
}

func GraphWorker() {
	t := time.NewTicker(1200 * time.Second)
	for {
		select {
		case <-t.C:
			graphContent := fmt.Sprintf(worker_content.Graph, generateGraph())
			err := os.WriteFile("/app/worker-content/tasks/graph.md", []byte(graphContent), 0664)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
