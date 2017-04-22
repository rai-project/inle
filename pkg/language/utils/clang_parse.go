package utils

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/elliotchance/c2go/ast"
)

type treeNode struct {
	indent int
	node   ast.Node
}

// see https://github.com/elliotchance/c2go
func ClangParse(fileName string) (string, error) {
	// 1. Preprocess
	pp, err := exec.Command("clang", "-E", fileName).Output()
	if err != nil {
		return "", err
	}

	ppFilePath := "/tmp/pp.c"
	err = ioutil.WriteFile(ppFilePath, pp, 0644)
	if err != nil {
		return "", err
	}

	// 2. Generate JSON from AST
	astPP, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", ppFilePath).Output()
	if err != nil {
		return "", err
	}

	lines := readAST(astPP)

	nodes := convertLinesToNodes(lines)
	tree := buildTree(nodes, 0)

	// Render(go_out, tree[0], "", 0, "")
	astTree := ast.NewAst()
	goOut := ast.Render(astTree, tree[0].(ast.Node))

	// Format the code
	goOutFmt, err := format.Source([]byte(goOut))
	if err != nil {
		return "", err
	}

	// Put together the whole file
	all := "package main\n\nimport (\n"

	for _, importName := range astTree.Imports() {
		all += fmt.Sprintf("\t\"%s\"\n", importName)
	}

	all += ")\n\n" + string(goOutFmt)

	return all, nil
}

func readAST(data []byte) []string {
	uncolored := regexp.MustCompile(`\x1b\[[\d;]+m`).ReplaceAll(data, []byte{})
	return strings.Split(string(uncolored), "\n")
}

func convertLinesToNodes(lines []string) []treeNode {
	nodes := []treeNode{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// It is tempting to discard null AST nodes, but these may
		// have semantic importance: for example, they represent omitted
		// for-loop conditions, as in for(;;).
		line = strings.Replace(line, "<<<NULL>>>", "NullStmt", 1)

		indentAndType := regexp.MustCompile("^([|\\- `]*)(\\w+)").FindStringSubmatch(line)
		if len(indentAndType) == 0 {
			panic(fmt.Sprintf("Cannot understand line '%s'", line))
		}

		offset := len(indentAndType[1])
		node := ast.Parse(line[offset:])

		indentLevel := len(indentAndType[1]) / 2
		nodes = append(nodes, treeNode{indentLevel, node})
	}

	return nodes
}

// buildTree convert an array of nodes, each prefixed with a depth into a tree.
func buildTree(nodes []treeNode, depth int) []ast.Node {
	if len(nodes) == 0 {
		return []ast.Node{}
	}

	// Split the list into sections, treat each section as a a tree with its own root.
	sections := [][]treeNode{}
	for _, node := range nodes {
		if node.indent == depth {
			sections = append(sections, []treeNode{node})
		} else {
			sections[len(sections)-1] = append(sections[len(sections)-1], node)
		}
	}

	results := []ast.Node{}
	for _, section := range sections {
		slice := []treeNode{}
		for _, n := range section {
			if n.indent > depth {
				slice = append(slice, n)
			}
		}

		children := buildTree(slice, depth+1)
		for _, child := range children {
			section[0].node.AddChild(child)
		}
		results = append(results, section[0].node)
	}

	return results
}
