package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/gabriel-vasile/mimetype"
)

type TemplateData struct {
	Query     string
	NumFiles  int
	Truncated bool
	Result    []FileMeta
}

var fullScan bool
var destinationMeta []FileMeta

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("templates/index.html")
	t, err := template.New("index").Parse(index_html)

	if err != nil {
		fmt.Fprintf(w, "Hata: %s!", err)
	}
	t.ExecuteTemplate(w, "index", nil)
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {

	file := r.URL.Query().Get("file")
	//fmt.Println("Dosya:", file)

	pathSlice := strings.Split(file, string(os.PathSeparator))
	fileName := pathSlice[len(pathSlice)-1]

	infile, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(w, "Dosya bunamadı: %s", file)
		return
	}

	mime, err := mimetype.DetectReader(infile)

	//fileHeader := make([]byte, 512)

	//Copy the headers into the fileHeader buffer

	//infile.Read(fileHeader)
	//Get content type of file
	//fileContentType := http.DetectContentType(fileHeader)

	//Get the file size
	fileStat, _ := infile.Stat()                       //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", mime.String())
	w.Header().Set("Content-Length", fileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	infile.Seek(0, 0)
	io.Copy(w, infile) //'Copy' the file to the client

	if err != nil {
		fmt.Println("cannot open file")
	}

	io.Copy(w, infile)
}
func searchHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{}

	t, err := template.New("index").Parse(index_html)

	q := r.URL.Query().Get("q")
	data.Query = q

	q = strings.TrimSpace(q)
	q = strings.ToLowerSpecial(unicode.TurkishCase, q)
	queryTerms := strings.Split(q, " ")

	var searchResult []FileMeta

	for _, fileMeta := range destinationMeta {
		hasAllKeys := true

		for _, term := range queryTerms {
			if term != "" {
				if !strings.Contains(strings.ToLowerSpecial(unicode.TurkishCase, fileMeta.FileName), term) {
					hasAllKeys = false
				}
			}
		}

		if hasAllKeys {
			searchResult = append(searchResult, fileMeta)
		}
	}

	data.NumFiles = len(searchResult)
	if len(searchResult) > 2000 {
		data.Truncated = true
		data.Result = searchResult[:2000]
	} else {
		data.Truncated = false
		data.Result = searchResult
	}

	err = t.ExecuteTemplate(w, "index", data)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	flag.Parse()
	// 1st argument is the directory location
	arg1 := flag.Arg(0)
	if len(arg1) < 1 {
		fmt.Println("Arama yapmak istediğiniz klasörü girin")
		fmt.Println("Örnek kullanım: filesearch.exe d:\\dosyalar")
		return
	}

	fmt.Println("Klasör:", arg1)
	fmt.Println("Aranacak klasörü indeksliyorum, lütfen bekleyin")

	destinationPath, err := filepath.Abs(arg1)
	if err != nil {
		fmt.Println("Klasör geçersiz")
		return
	}

	//var destinationMeta []FileMeta

	if _, err := os.Stat(destinationPath + string(os.PathSeparator) + "folderMeta.txt"); err == nil {
		// path/to/whatever exists
		destinationMeta, err = loadFolderMeta(destinationPath + string(os.PathSeparator) + "folderMeta.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		fmt.Println("Dosya listesi oluşturuluyor")
		destinationMeta, err = createFolderMeta(destinationPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		saveFolderMeta(destinationMeta, destinationPath+string(os.PathSeparator)+"folderMeta.txt")
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		fmt.Println(err)
		return
	}

	//fmt.Print(destinationMeta)
	//--------------------------
	fmt.Println("Tamamlandı")
	fmt.Println("Uygulamaya erişmek için web tarayıcınız ile http://localhost:8080 adresine gidin")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/open/", serveFileHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	openbrowser("http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
