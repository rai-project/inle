package utils

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

func ClangParse(fileName string) ([]string, error) {
	// 1. Preprocess
	pp, err := exec.Command("clang", "-E", fileName).Output()
	if err != nil {
		return nil, err
	}

	ppFilePath := "/tmp/pp.c"
	err = ioutil.WriteFile(ppFilePath, pp, 0644)
	if err != nil {
		return nil, err
	}

	// 2. Generate JSON from AST
	astPP, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", ppFilePath).Output()
	if err != nil {
		return nil, err
	}

	lines := readAST(astPP)
	return lines, nil
}

func readAST(data []byte) []string {
	uncolored := regexp.MustCompile(`\x1b\[[\d;]+m`).ReplaceAll(data, []byte{})
	return strings.Split(string(uncolored), "\n")
}
