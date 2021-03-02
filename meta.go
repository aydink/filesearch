package main

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type FileMeta struct {
	FileName string
	Path     string
	Hash     string
	Size     int64
	ModTime  int64
}

func exist(directory string) bool {
	absolutePath, err := filepath.Abs(directory)
	if err != nil {
		fmt.Println("path is not valid")
		return false
	}

	if info, err := os.Stat(absolutePath); !os.IsNotExist(err) {
		if info.IsDir() {
			return true
		} else {
			fmt.Println("path is not a directory!")
			return false
		}
	} else {
		fmt.Println("directory does not exist!")
		return false
	}
}

func createFolderMeta(directory string) ([]FileMeta, error) {

	var folderMeta []FileMeta

	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {

			meta, err := getFileMeta(path)
			//fmt.Println(path)

			if err != nil {
				fmt.Println(err)
			}

			folderMeta = append(folderMeta, meta)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return folderMeta, err
	}

	return folderMeta, nil
}

func loadFolderMeta(filePath string) ([]FileMeta, error) {

	var folderMeta []FileMeta

	// read data from CSV file
	csvFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
		return folderMeta, err
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	// Use tab-delimited instead of comma
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return folderMeta, err
	}

	for _, row := range csvData {
		meta := FileMeta{}

		meta.Hash = row[0]
		meta.FileName = row[1]
		meta.Path = row[2]
		meta.ModTime, err = strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		meta.Size, err = strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		folderMeta = append(folderMeta, meta)
	}

	return folderMeta, nil
}

func getFileMeta(filePath string) (FileMeta, error) {

	meta := FileMeta{}
	absPath, err := filepath.Abs(filePath)

	if err != nil {
		return meta, err
	}

	meta.FileName = filepath.Base(filePath)
	meta.Path = absPath

	//fmt.Println(path)
	/*
		infile, err := os.Open(filePath)

		if err != nil {
			return meta, err
		}
	*/

	stat, err := os.Stat(absPath)
	if err != nil {
		return meta, err
	}

	hash := md5.New()

	/*
		if fullScan {
			io.Copy(hash, infile)
		} else {
			// just read firs 100kb
			limitReader := io.LimitReader(infile, 100*1024)
			io.Copy(hash, limitReader)
		}
	*/
	meta.Hash = hex.EncodeToString(hash.Sum([]byte("")))
	meta.Size = stat.Size()
	meta.ModTime = stat.ModTime().Unix()

	return meta, nil
}

/*
func processFile(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {

		meta, err := getFileMeta(path)

		if err != nil {
			fmt.Println(err)
		}

		folderMeta = append(folderMeta, meta)
	}

	return nil
}
*/

func saveFolderMeta(folderMeta []FileMeta, filePath string) {
	f, err := os.OpenFile(filePath, os.O_CREATE, 0777)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
	}

	for i := range folderMeta {
		fmt.Fprintf(f, "%s\t", folderMeta[i].Hash)
		fmt.Fprintf(f, "%s\t", folderMeta[i].FileName)
		fmt.Fprintf(f, "%s\t", folderMeta[i].Path)
		fmt.Fprintf(f, "%d\t", folderMeta[i].ModTime)
		fmt.Fprintf(f, "%d\n", folderMeta[i].Size)
	}
}

func printFolderMeta(folderMeta []FileMeta) {
	for i := range folderMeta {
		fmt.Printf("%s\t", folderMeta[i].Hash)
		fmt.Printf("%s\t", folderMeta[i].FileName)
		fmt.Printf("%s\t", folderMeta[i].Path)
		fmt.Printf("%d\t", folderMeta[i].ModTime)
		fmt.Printf("%d\n", folderMeta[i].Size)
	}
}

func createMap(fileMeta []FileMeta) map[string]FileMeta {

	metaMap := make(map[string]FileMeta)

	for _, v := range fileMeta {
		metaMap[v.Hash] = v
	}

	return metaMap
}
