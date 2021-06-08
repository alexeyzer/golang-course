package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strconv"
)

type ByName []fs.DirEntry

func CountDir(files []fs.DirEntry, PrintFiles bool) int {
	count := 0
	for _, f := range files {
		if f.IsDir() == true {
			count = count + 1
		} else if PrintFiles == true {
			count = count + 1
		}
	}
	return count
}

func (f ByName) Len() int { return CountDir(f, true) }

func (f ByName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }

func (f ByName) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func RealDir(Out io.Writer, path string, printFiles bool, tab string) (err error) {

	file, err := os.Open(path)

	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() == true {
		files, err := file.ReadDir(-1)
		if err != nil {
			return err
		}
		Count := CountDir(files, printFiles)
		sort.Sort(ByName(files))
		for _, f := range files {
			if printFiles == false && f.IsDir() == false {
				continue
			}
			if Count > 1 {
				fmt.Fprint(Out, tab+"├───")
			} else {
				fmt.Fprint(Out, tab+"└───")
			}
			if f.IsDir() == true {
				fmt.Fprintln(Out, f.Name())
				if Count > 1 {
					err = RealDir(Out, path+"/"+f.Name(), printFiles, tab+"│\t")
				} else {
					err = RealDir(Out, path+"/"+f.Name(), printFiles, tab+"\t")
				}
				if err != nil {
					return err
				}
			} else if printFiles == true {
				info, err := f.Info()
				if err != nil {
					return nil
				}
				var str string
				if info.Size() == 0 {
					str = "(empty)"
				} else {
					str = "(" + strconv.Itoa(int(info.Size())) + "b)"
				}
				fmt.Fprintln(Out, f.Name(), str)
			}
			Count = Count - 1
		}
	}
	return nil
}

func dirTree(Out io.Writer, path string, printFiles bool) (err error) {
	return RealDir(Out, path, printFiles, "")
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
