package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/gum/pager"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

var (
	enumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#353ec5")).MarginRight(1)
	rootStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#0cff04")).Bold(true)
	itemStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffbbe8"))
)

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

func getTree(root string, maxDepth int, showHidden bool) *tree.Tree {
	t := tree.New()
	last_node := make(map[int]*tree.Tree)
	last_node[0] = t

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		depth := strings.Count(path, string(os.PathSeparator)) - strings.Count(root, string(os.PathSeparator))
		name := filepath.Base(path)
		if maxDepth > 0 && depth > maxDepth {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			if !showHidden && isHidden(name) {
				return fs.SkipDir
			}
			t_new := tree.Root(rootStyle.Render(name))
			last_node[depth].Child(t_new)
			last_node[depth+1] = t_new
		} else {
			last_node[depth].Child(itemStyle.Render(name))
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return t
}

type cliArgs struct {
	root       string
	maxDepth   int
	showHidden bool
}

func argParse() cliArgs {
	args := cliArgs{}
	flag.Usage = func() {
		fmt.Println(`Usage: tree [options] dir

dir is the root directory. Default is "."`)
		flag.PrintDefaults()
	}
	flag.IntVar(&args.maxDepth, "depth", 0, "Maximum depth to display")
	flag.IntVar(&args.maxDepth, "d", 0, "Maximum depth to display")
	flag.BoolVar(&args.showHidden, "show-hidden", false, "Show hidden directories. Default false")
	flag.Parse()

	args.root = "."
	if flag.NArg() > 0 {
		args.root = flag.Arg(0)
	}

	return args
}

func main() {
	args := argParse()
	t := getTree(args.root, args.maxDepth, args.showHidden).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle)

	opts := pager.Options{
		Content: t.String(),
	}
	if err := opts.Run(); err != nil {
		log.Fatal(err)
	}
}
