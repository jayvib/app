package main

import (
	"html/template"
	"os"
)

var tmpl = `
<html>
<head>
<title>Template Argument</title>
</head>
<body>
<h2>One Piece Characters</h2>
{{ range $key, $value := . }}
	The key is {{ $key }} and the value is {{ $value }}
{{ end }}
</body>
</html>`

func main() {
	t := template.Must(template.New("tmpl").Parse(tmpl))
	onePiece := map[string]string{
		"Luffy": "Monkey",
		"Sanji": "Vinsmoke",
		"Nami": "The Navigator",
		"Chopper": "The Doctor",
	}

	t.Execute(os.Stdout, onePiece)
}
