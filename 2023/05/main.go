package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

type MappingPart struct {
	Range            int64
	SourceStart      int64
	DestinationStart int64
}

type Mapping struct {
	Name  string
	Parts []MappingPart
}

func (m Mapping) Call(source int64) int64 {
	for _, mp := range m.Parts {
		if source < mp.SourceStart || source >= (mp.SourceStart+mp.Range) {
			continue
		}

		return mp.DestinationStart + (source - mp.SourceStart)
	}

	return source
}

func parseNumbers(in string) []int64 {
	out := []int64{}

	for _, n := range strings.Split(in, " ") {
		nn := strings.TrimSpace(n)
		if nn == "" {
			continue
		}

		out = append(out, aoc.MustParse(nn))
	}

	return out
}

func main() {
	groups := [][]string{}
	lines := []string{}

	aoc.EachLine("input.txt", func(line string) {
		if line == "" {
			groups = append(groups, lines)
			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	})
	groups = append(groups, lines)

	seeds := parseNumbers(strings.TrimPrefix(groups[0][0], "seeds: "))

	mappings := map[string]Mapping{}

	for i := 1; i < len(groups); i++ {
		mapping := Mapping{
			Parts: []MappingPart{},
		}

		for i, line := range groups[i] {
			if i == 0 {
				mapping.Name = strings.TrimSuffix(line, " map:")
				continue
			}

			parts := parseNumbers(line)
			mapping.Parts = append(mapping.Parts, MappingPart{
				DestinationStart: parts[0],
				SourceStart:      parts[1],
				Range:            parts[2],
			})
		}

		mappings[mapping.Name] = mapping
	}

	c := map[int64]int64{}
	lowestLocation := seedToLocation(c, mappings, seeds[0])
	for _, seed := range seeds {
		loc := seedToLocation(c, mappings, seed)
		if loc < lowestLocation {
			lowestLocation = loc
		}
	}

	log.Printf("1: %v", lowestLocation)
}

func seedToLocation(c map[int64]int64, mappings map[string]Mapping, seed int64) int64 {
	if r, ok := c[seed]; ok {
		return r
	}

	result := seed
	for _, fn := range []string{
		"seed-to-soil",
		"soil-to-fertilizer",
		"fertilizer-to-water",
		"water-to-light",
		"light-to-temperature",
		"temperature-to-humidity",
		"humidity-to-location",
	} {
		result = mappings[fn].Call(result)
	}

	c[seed] = result

	return result
}
