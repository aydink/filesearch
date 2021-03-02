package main

var index_html = `
<!doctype html>
<html lang="en">

<head>
    <title>Resmi Yayınlar</title>
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
            <a><b>Resmi yayınlar</b></a>            
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

    {{range .Result}}    
        <a href="/open?file={{.Path}}">{{.FileName}}</a><br/>    
    {{end}}

</body>
</html>
`
