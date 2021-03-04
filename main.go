package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type TemplateData struct {
	Query     string
	NumFiles  int
	Truncated bool
	Result    []FileMeta
}

var fullScan bool
var destinationMeta []FileMeta

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
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	openbrowser("http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
