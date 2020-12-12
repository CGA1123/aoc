package main

import (
	"bufio"
	"container/list"
	"log"
	"os"
	"strings"
)

type visit struct {
	depth int
	node  string
}

func total(graph map[string]map[string]struct{}, parent string) int {
	visited := map[string]struct{}{}

	var count int
	queue := list.New()
	queue.PushBack(visit{node: parent, depth: 0})

	for queue.Len() > 0 {
		el := queue.Front()
		queue.Remove(el)

		node := el.Value.(visit)
		visited[node.node] = struct{}{}
		count = count + node.depth

		for edge := range graph[node.node] {
			if _, alreadyVisited := visited[edge]; alreadyVisited {
				continue
			}

			queue.PushFront(visit{node: edge, depth: node.depth + 1})
		}
	}

	return count
}

func find(graph map[string]map[string]struct{}, start, finish string) int {
	visited := map[string]struct{}{}

	queue := list.New()
	queue.PushBack(visit{node: start, depth: 0})

	var depth int

	for queue.Len() > 0 {
		el := queue.Front()
		queue.Remove(el)

		node := el.Value.(visit)
		visited[node.node] = struct{}{}

		if node.node == finish {
			depth = node.depth
			break
		}

		var x []string
		for edge := range graph[node.node] {
			x = append(x, edge)
			if _, alreadyVisited := visited[edge]; !alreadyVisited {
				queue.PushBack(visit{node: edge, depth: node.depth + 1})
			}
		}
	}

	return depth - 2
}

type Node struct {
	value    string
	siblings []*Node
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))

	orbits := map[string]map[string]struct{}{}

	for s.Scan() {
		p := strings.Split(s.Text(), ")")
		if len(p) != 2 {
			log.Printf("error: malformed orbit [%v]", s.Text())
		}

		if _, ok := orbits[p[0]]; !ok {
			orbits[p[0]] = map[string]struct{}{}
		}

		if _, ok := orbits[p[1]]; !ok {
			orbits[p[1]] = map[string]struct{}{}
		}

		orbits[p[0]][p[1]] = struct{}{}
		orbits[p[1]][p[0]] = struct{}{}
	}

	log.Printf("%v", total(orbits, "COM"))
	log.Printf("%v", find(orbits, "YOU", "SAN"))
}
