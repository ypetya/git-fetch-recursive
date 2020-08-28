package walkdir

import (
	"os"
)

func WalkDir(folder string,
	visit chan<- string,
	isExcluded func(string, string) bool) error {

	f, err := os.Open(folder)
	if err != nil {
		return err
	}

	visit <- folder

	list, err := f.Readdir(-1)
	f.Close()
	for _, fileInfo := range list {
		if fileInfo.IsDir() {
			folderName := fileInfo.Name()
			fN := folder + "/" + folderName
			if !isExcluded(folderName, fN) {
				err := WalkDir(fN, visit, isExcluded)
				if err != nil {
					return err
				}
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}
