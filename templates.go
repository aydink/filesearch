package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

var index_html = `
<!doctype html>
<html lang="en">

<head>
    <title>Basit dosya arama motoru</title>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <style>
		.search-box {
            font-size: 1.5em;
            width: 40%;
        }

        .search-button {
            font-size: 1.5em;            
        }
        .table { display:table; border-spacing: 5px; border-collapse:separate;  }
        .row { display: table-row }
        .cell { display: table-cell; margin-left: 15px; }

    </style>
   <script>
       document.addEventListener("DOMContentLoaded", function(event) {
            document.getElementById("q").focus();
       });
   </script>
</head>

<body>
    <div class="container">
        <nav class="navbar fixed-top navbar-light bg-light py-2">
            <a><b>Dosyaları bulmak için anahtar kelimeler girin</b></a>            
        </nav>
    </div>

    <div class="container">
        <form class="" method="GET" action="/search/">
            <div class="input-group">
                <input type="text" class="search-box" name="q" id="q" value="{{.Query}}">&nbsp;
                <button class="search-button" type="submit">&nbsp;&nbsp; Ara &nbsp;&nbsp;</button>&nbsp;                        
            </div>
        </form>
    </div>
	{{if .NumFiles}}
		<b>{{ $.NumFiles}}</b> dosya bulundu
		{{if $.Truncated}}, yalnızca ilk 2000 sonuç gösterilmiştir.{{end}}
		<br/>
	{{end}}

    <div class="table">

    {{if .Result}}        
        <div class="row">
            <div class="cell"><b><a href="/search?q={{- $.Query -}}&order=mtime">Tarih</a></b></div>            
            <div class="cell"><b><a href="/search?q={{- $.Query -}}&order=name">Dosya Adı</a></b></div>
            <div class="cell"><b><a href="/search?q={{- $.Query -}}&order=size">Boyut</a></b></div>
        </div>
    {{end}}

    {{range .Result}}
        <div class="row">
            <div class="cell"> {{printModTime .ModTime}} </div>            
            <div class="cell"> <a href="/open?file={{.Path}}">{{.FileName}}</a></div>
            <div class="cell"> {{fileSize .Size}}</div>
        </div>    
    {{end}}
    </div>

</body>
</html>
`
var funcMap template.FuncMap
var tpl *template.Template
var err error

func printModTime(modTime int64) string {
	return time.Unix(modTime, 0).Format("2006-01-02 15:04:05")
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func init() {

	funcMap = template.FuncMap{
		"printModTime": printModTime,
		"fileSize":     ByteCountSI,
	}

	tpl, err = template.New("index").Funcs(funcMap).Parse(index_html)

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}
