package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type ByName []string

func (f ByName) Len() int { return len(f) }

func (f ByName) Less(i, j int) bool { return f[i] < f[j] }

func (f ByName) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func CountDir(files []string, PrintFiles bool, path string) (int, error) {
	count := 0
	for _, f := range files {
		filename := path + "/" + f
		file, err := os.Open(filename)
		if err != nil {
			return 0, err
		}
		info, err := file.Stat()
		if err != nil {
			return 0, err
		}
		if info.IsDir() == true {
			count = count + 1
		} else if PrintFiles == true {
			count = count + 1
		}
	}
	return count, nil
}

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
		files, err := file.Readdirnames(0)
		if err != nil {
			return err
		}
		Count, err := CountDir(files, printFiles, path)
		sort.Sort(ByName(files))
		for _, f := range files {
			tempfile, err := os.Open(path + "/" + f)
			infotemp, err := tempfile.Stat()
			if printFiles == false && infotemp.IsDir() == false {
				continue
			}
			if Count > 1 {
				fmt.Fprint(Out, tab+"├───")
			} else {
				fmt.Fprint(Out, tab+"└───")
			}
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			if infotemp.IsDir() == true {
				fmt.Fprintln(Out, infotemp.Name())
				if Count > 1 {
					err = RealDir(Out, path+"/"+infotemp.Name(), printFiles, tab+"│\t")
				} else {
					err = RealDir(Out, path+"/"+infotemp.Name(), printFiles, tab+"\t")
				}
				if err != nil {
					return err
				}
			} else if printFiles == true {
				var str string
				if infotemp.Size() == 0 {
					str = "(empty)"
				} else {
					str = "(" + strconv.Itoa(int(infotemp.Size())) + "b)"
				}
				fmt.Fprintln(Out, infotemp.Name(), str)
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
