package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// getPackageAPI retrieves the exported functions and their signatures from a package at a specific version.
func getPackageAPI(repoPath string, version string) (map[string]string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo: %v", err)
	}

	ref, err := r.Tag(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %v", err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %v", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %v", err)
	}

	api := make(map[string]string)
	err = tree.Files().ForEach(func(f *object.File) error {
		content, err := f.Contents()
		if err != nil {
			return err
		}

		if !strings.HasSuffix(f.Name, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, f.Name, content, parser.ParseComments)
		if err != nil {
			return err
		}

		conf := types.Config{Importer: importer.Default()}
		info := &types.Info{
			Defs: make(map[*ast.Ident]types.Object),
		}
		_, err = conf.Check(f.Name, fset, []*ast.File{node}, info)
		if err != nil {
			return err
		}

		fmt.Println(info.Defs)

		for id, obj := range info.Defs {
			if obj != nil && obj.Exported() {
				api[id.Name] = obj.Type().String()
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over files: %v", err)
	}

	return api, nil
}

// compareAPIs compares the APIs of two versions and prints the differences.
func compareAPIs(api1, api2 map[string]string) {
	fmt.Printf("%v", api1)
	fmt.Printf("%v", api2)
	for name, sig1 := range api1 {
		if sig2, ok := api2[name]; ok {
			if sig1 != sig2 {
				fmt.Printf("Signature change: %s: %s -> %s\n", name, sig1, sig2)
			}
		} else {
			fmt.Printf("Removed: %s %s\n", name, sig1)
		}
	}
	for name, sig2 := range api2 {
		if _, ok := api1[name]; !ok {
			fmt.Printf("Added: %s %s\n", name, sig2)
		}
	}
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <repo> <version1> <version2>", os.Args[0])
	}
	repo := os.Args[1]
	version1 := os.Args[2]
	version2 := os.Args[3]

	path := fmt.Sprintf("/tmp/repo/%s", repo)

	// Clone the repository into memory
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      fmt.Sprintf("https://github.com/%s.git", repo),
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed to clone repo: %v", err)
	}
	defer os.RemoveAll(path)

	api1, err := getPackageAPI(path, version1)
	if err != nil {
		log.Fatalf("Failed to get API for version %s: %v", version1, err)
	}
	api2, err := getPackageAPI(path, version2)
	if err != nil {
		log.Fatalf("Failed to get API for version %s: %v", version2, err)
	}

	compareAPIs(api1, api2)
}
