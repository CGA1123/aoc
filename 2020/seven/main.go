package main

import (
	"bufio"
	"container/list"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Edge struct {
	to     string
	weight int
}

var edgeRegex = regexp.MustCompile(`^(\d+) (.*) bags?\.?$`)

func partOne(graph map[string][]string, node string) int {
	visited := map[string]struct{}{}

	queue := list.New()
	queue.PushBack(node)

	for queue.Len() > 0 {
		el := queue.Front()
		queue.Remove(el)

		node := el.Value.(string)
		visited[node] = struct{}{}

		for _, edge := range graph[node] {
			if _, alreadyVisited := visited[edge]; !alreadyVisited {
				queue.PushBack(edge)
			}
		}

	}

	return len(visited) - 1
}

func partTwo(graph map[string][]Edge, node string) int {
	var count int
	for _, e := range graph[node] {
		count = count + (e.weight * (1 + partTwo(graph, e.to)))
	}

	return count
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))

	graph := map[string][]Edge{}
	inverted := map[string][]string{}

	for s.Scan() {
		line := strings.Split(s.Text(), " bags contain ")
		node, edgesStr := line[0], line[1]

		for _, edgeStr := range strings.Split(edgesStr, ", ") {
			for _, match := range edgeRegex.FindAllStringSubmatch(edgeStr, 2) {
				weight, _ := strconv.Atoi(match[1])
				to := match[2]

				inverted[to] = append(inverted[to], node)
				graph[node] = append(graph[node], Edge{to: to, weight: weight})
			}
		}

	}
	if s.Err() != nil {
		log.Printf("error scanning: %v", s.Err())
		return
	}

	log.Printf("pt1: %v", partOne(inverted, "shiny gold"))
	log.Printf("pt2: %v", partTwo(graph, "shiny gold"))
}
