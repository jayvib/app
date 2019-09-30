package main

import (
	"html/template"
	"os"
)

func main() {
	tEmpty := template.New("templateTest")
	tEmpty = template.Must(tEmpty.Parse(`Empty pipeline if demo: {{if ""}} will not be outputed{{end}}\n`))
	tEmpty.Execute(os.Stdout, nil)

	tWithValue := template.New("templateTest")
	tWithValue = template.Must(tWithValue.Parse(`WithValue pipeline if demo: {{if "anything"}} will be outputed{{end}}`))
	tWithValue.Execute(os.Stdout, nil)

	tifelse := template.New("templateTest")
	tifelse = template.Must(tifelse.Parse(`ifelse pipeline if demo: {{if ""}} will be outputed{{else}} else outputed{{end}}`))
	tifelse.Execute(os.Stdout, nil)
}
