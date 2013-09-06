// Copyright 2013 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

// Recursively batch rename files and directories by regular expression.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var (
	path string // path
	src  string // source regular expression
	repl string // replacement
	R    bool   // recursive
	D    bool   // Rename directories
)

func init() {
	flag.StringVar(&src, "s", "", "Regular expression")
	flag.StringVar(&repl, "r", "", "Replacement")
	flag.BoolVar(&R, "R", true, "Recursively rename")
	flag.BoolVar(&D, "D", true, "Rename directories")

	flag.Usage = func() {
		fmt.Println("\nbrename\n  Recursively batch rename files and directories by regular expression.")
		fmt.Printf("\nUsage: %s -s <regexp> -r <replacement> [-R] [-D] [path...]\n\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("\n  Site: https://github.com/shenwei356/BatchFileRename")
		fmt.Println("Author: Wei Shen (shenwei356@gmail.com)\n")
	}

	flag.Parse()
	if src == "" && repl == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	re, err := regexp.Compile(src)
	if err != nil {
		recover()
		fmt.Println("Bad regular expression!")
		return
	}

	var paths []string
	if len(flag.Args()) == 0 {
		paths = []string{"./"}
	} else {
		paths = flag.Args()
	}
	for _, path := range paths {
		fmt.Printf("%s:\n", path)
		n, err := BatchRename(path, re, repl, R, D)
		if err != nil {
			recover()
			fmt.Println(err)
			continue
		}
		fmt.Printf("%d files be renamed\n\n", n)
	}
}

func BatchRename(path string, re *regexp.Regexp, repl string, recursive bool, D bool) (uint, error) {
	var n uint = 0

	_, err := ioutil.ReadFile(path)
	// it's a file
	if err == nil {
		n, err = Rename(path, re, repl)
		if err != nil {
			recover()
			fmt.Println(err)
			return 0, err
		}
		return n, nil
	}

	// it's a directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		recover()
		if os.IsNotExist(err) {
			return 0, errors.New("Not Exist: " + path)
		}
		return 0, errors.New("ReadDir Error: " + path)
	}

	var filename string
	for _, file := range files {
		filename = file.Name()
		if filename == "." || filename == ".." {
			continue
		}

		fileFullPath := filepath.Join(path, filename)
		// sub directory
		if file.IsDir() {
			if recursive {
				num, err := BatchRename(fileFullPath, re, repl, recursive, D)
				if err != nil {
					recover()
					fmt.Println(err)
					continue
				}
				n += num
			}
			// Rename directories
			if D {
				num, err := Rename(fileFullPath, re, repl)
				if err != nil {
					recover()
					fmt.Println(err)
					continue
				}
				n += num
			}
		} else {
			num, err := Rename(fileFullPath, re, repl)
			if err != nil {
				recover()
				fmt.Println(err)
				continue
			}
			n += num
		}
	}
	return n, nil
}

func Rename(path string, re *regexp.Regexp, repl string) (uint, error) {
	dir, filename := filepath.Split(path)
	// not matched
	if !re.Match([]byte(filename)) {
		return 0, nil
	}

	filename2 := re.ReplaceAllString(filename, repl)
	// not changed
	if filename2 == filename {
		return 0, nil
	}

	err := os.Rename(path, filepath.Join(dir, filename2))
	if err != nil {
		return 0, errors.New("\nRename file error: [" + filename + " -> " + filename2 + "].")
	}

	return 1, nil
}