package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/CGA1123/aoc"
)

type Object interface {
	Size() int64
	Name() string
	IsDir() bool
}

type File struct {
	name   string
	size   int64
	parent *Dir
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int64 {
	return f.size
}

func NewFile(name string, size int64, p *Dir) *File {
	return &File{name: name, size: size, parent: p}
}

type Dir struct {
	name     string
	children map[string]Object
	parent   *Dir
	size     int64
}

func (d *Dir) IsDir() bool {
	return true
}

func (d *Dir) AddChild(o Object) {
	d.children[o.Name()] = o
}

func (d *Dir) Name() string {
	return d.name
}

func (d *Dir) Size() int64 {
	if d.size != -1 {
		return d.size
	}

	d.size = 0

	for _, c := range d.children {
		d.size = d.size + c.Size()
	}

	return d.size
}

func (d *Dir) Cd(name string) *Dir {
	if name == ".." {
		return d.parent
	}

	o := d.children[name]
	if o == nil {
		log.Fatalf("%v is not a child", name)
	}

	dir, ok := o.(*Dir)
	if !ok {
		log.Fatalf("%v is not a directory", name)
	}

	return dir
}

func NewDir(name string, p *Dir) *Dir {
	return &Dir{size: -1, name: name, parent: p, children: make(map[string]Object)}
}

func repl(root *Dir) {
	cwd := root

	r := bufio.NewReader(os.Stdin)

	ls := func() {
		fmt.Printf("%-10s name\n\n", "size")
		for _, c := range cwd.children {
			if c.IsDir() {
				fmt.Printf("%v %v/\n", strings.Repeat(" ", 10), c.Name())

				continue
			}

			fmt.Printf("%-10v %v\n", c.Size(), c.Name())
		}
	}

	cd := func(to string) {
		switch to {
		case "/":
			cwd = root
		case "..":
			cwd = cwd.parent
		default:
			if o, ok := cwd.children[to]; ok {
				d, ok := o.(*Dir)
				if !ok {
					fmt.Printf("err: %v is not a directory\n", to)
					return
				}

				cwd = d
				return
			}

			fmt.Printf("err: %v not found\n", to)
		}

	}

	size := func(c string) {
		o, ok := cwd.children[c]
		if !ok {
			fmt.Printf("err: %v not found", c)
		}

		fmt.Printf("%v\n", o.Size())
	}

	for {
		fmt.Printf("%v $ ", cwd.Name())
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Printf("err: %v", err)
			continue
		}

		args := strings.Split(strings.TrimSpace(line), " ")
		switch args[0] {
		case "":
			continue
		case "exit":
			fmt.Println("Bye!")
			return
		case "ls":
			ls()
		case "cd":
			cd(args[1])
		case "size":
			size(args[1])
		default:
			fmt.Printf("err: unknown cmd %v\n", args[0])
		}
	}
}

func main() {
	root := NewDir("/", nil)
	var cwd *Dir

	aoc.EachLine("input.txt", func(s string) {
		if s[0:2] == "$ " {
			args := strings.Split(s[2:], " ")

			switch args[0] {
			case "cd":
				if args[1] == "/" {
					cwd = root
				} else {
					cwd = cwd.Cd(args[1])
				}
			case "ls":
			}

			return
		}

		var child Object

		output := strings.Split(s, " ")
		if output[0] == "dir" {
			child = NewDir(output[1], cwd)
		} else {
			child = NewFile(output[1], aoc.MustParse(output[0]), cwd)
		}

		cwd.AddChild(child)
	})

	root.Size()

	sizes := []int64{}
	traverse(root, func(o Object) {
		if !o.IsDir() {
			return
		}

		if o.Size() <= 100000 {
			sizes = append(sizes, o.Size())
		}
	})

	totalSize := int64(0)
	for _, size := range sizes {
		totalSize = totalSize + size
	}
	log.Printf("%v", totalSize)

	capacity, required := int64(70000000), int64(30000000)
	free := (capacity - root.Size())
	toBeFreed := required - free
	minToBeFreed := int64(math.MaxInt64)

	traverse(root, func(o Object) {
		if !o.IsDir() {
			return
		}
		if o.Size() > minToBeFreed || o.Size() < toBeFreed {
			return
		}

		minToBeFreed = o.Size()
	})

	log.Printf("%v", minToBeFreed)

	if os.Getenv("AOC_REPL") == "1" {
		repl(root)
	}
}

func traverse(d *Dir, f func(Object)) {
	f(d)

	for _, c := range d.children {
		if c.IsDir() {
			traverse(c.(*Dir), f)
		} else {
			f(c)
		}

	}
}
