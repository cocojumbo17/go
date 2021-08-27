package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	suphix           = "├───"
	suphix_last      = "└───"
	suphix_level     = "│\t"
	suphix_level_rev = "\t│"
	suphix_empty     = "\t"
)

type node struct {
	is_root  bool
	is_file  bool
	name     string
	size     int64
	children []*node
	is_last  bool
	parent   *node
}

func reverse(in string) string {
	var sb strings.Builder
	runes := []rune(in)
	for i := len(runes) - 1; 0 <= i; i-- {
		sb.WriteRune(runes[i])
	}
	return sb.String()
}
func (n node) String() string {

	var buf strings.Builder
	parent := n.parent
	for parent != nil && !parent.is_root {
		if parent.is_last {
			buf.WriteString(suphix_empty)
		} else {
			buf.WriteString(suphix_level_rev)
		}
		parent = parent.parent
	}
	ttt := buf.String()
	buf.Reset()
	buf.WriteString(reverse(ttt))

	if n.is_last {
		buf.WriteString(suphix_last)
	} else {
		buf.WriteString(suphix)
	}
	buf.WriteString(n.name)
	if n.is_file {
		if n.size == 0 {
			buf.WriteString(" (empty)")
		} else {
			buf.WriteString(fmt.Sprintf(" (%db)", n.size))
		}
	}
	return buf.String()
}

func collectTreeInfo(path string, printFile bool, n *node) error {
	files, err := os.ReadDir(path)
	if err == nil {
		for _, val := range files {

			if !val.IsDir() && !printFile {
				continue
			}

			child := &node{
				is_root:  false,
				is_file:  !val.IsDir(),
				name:     val.Name(),
				children: []*node{},
				parent:   n,
			}

			if val.IsDir() {
				newpath := path + string(os.PathSeparator) + child.name
				err = collectTreeInfo(newpath, printFile, child)
				if err != nil {
					return err
				}
			} else if printFile {
				fi, _ := val.Info()
				child.size = fi.Size()
			}
			n.children = append(n.children, child)
		}
	}
	return err
}

/*
func sortTree(n *node) {
	sort.Slice(n.children, func(i, j int) bool {
		return n.children[i].name < n.children[j].name
	})
	for _, ch := range n.children {
		sortTree(ch)
	}
}
*/
func makrLast(n *node) {
	if len(n.children) > 0 {
		n.children[len(n.children)-1].is_last = true
		for _, ch := range n.children {
			makrLast(ch)
		}
	}

}

func printTree(out io.Writer, n *node) {
	if !n.is_root {
		fmt.Fprintln(out, n)
	}
	for _, ch := range n.children {
		printTree(out, ch)
	}
}

func dirTree(out io.Writer, path string, printFile bool) error {
	root := node{
		is_root: true,
		parent:  nil,
	}
	err := collectTreeInfo(path, printFile, &root)
	if err == nil {
		//sortTree(&root)
		makrLast(&root)
		printTree(out, &root)
	}
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
