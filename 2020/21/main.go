package main

import (
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/CGA1123/aoc"
)

var RecipeExp = regexp.MustCompile(`(?P<ingredients>.*) \(contains (?P<allergens>.*)\)`)

func PartOne(counts map[string]int64, unsafe map[string]*aoc.Set) int64 {
	var allUnsafe *aoc.Set
	for _, un := range unsafe {
		if allUnsafe == nil {
			allUnsafe = un
		} else {
			allUnsafe = aoc.Union(allUnsafe, un)
		}
	}

	var total int64
	for ingredient, count := range counts {
		if allUnsafe.Contains(ingredient) {
			continue
		}

		total += count
	}

	return total
}

type solve struct {
	Allergen   string
	Ingredient string
}

func PartTwo(unsafe map[string]*aoc.Set) string {
	total := len(unsafe)
	var solution []solve

	for len(solution) != total {
		for allergen, candidates := range unsafe {
			if candidates.Size() != 1 {
				continue
			}

			element := candidates.Elements()[0].(string)
			for _, v := range unsafe {
				v.Remove(element)
			}

			solution = append(solution, solve{Allergen: allergen, Ingredient: element})
		}

	}

	sort.Slice(solution, func(i, j int) bool {
		return solution[i].Allergen < solution[j].Allergen
	})

	var ingredients []string
	for _, v := range solution {
		ingredients = append(ingredients, v.Ingredient)
	}

	return strings.Join(ingredients, ",")
}

func main() {
	counts := map[string]int64{}
	unsafe := map[string]*aoc.Set{}

	aoc.EachLine("input.txt", func(l string) {
		capture := aoc.Capture(RecipeExp, l)
		ingredients := strings.Split(capture["ingredients"], " ")
		allergens := strings.Split(capture["allergens"], ", ")

		for _, ingredient := range ingredients {
			counts[ingredient] += 1
		}

		for _, allergen := range allergens {
			a := aoc.NewSet()

			for _, ingredient := range ingredients {
				a.Add(ingredient)
			}

			if _, ok := unsafe[allergen]; !ok {
				unsafe[allergen] = a
			} else {
				unsafe[allergen] = aoc.Intersection(unsafe[allergen], a)
			}
		}
	})

	log.Printf("pt(1): %v", PartOne(counts, unsafe))
	log.Printf("pt(2): %v", PartTwo(unsafe))
}
