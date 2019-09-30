package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

type Person struct {
	UserName string
	Emails []string
	Friends []Friend
}

type Friend struct {
	FName string
}

func EmailDealWith(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	// find the @ symbol
	substrs := strings.Split(s, "@")
	if len(substrs) != 2 {
		return s
	}
	// replace the @ by " at "
	return (substrs[0] + " at " + substrs[1])
}

func main() {
	t := template.New("fieldNameExample")
	t = t.Funcs(template.FuncMap{"emailDeal": EmailDealWith})
	t  = template.Must(t.Parse(`hello {{.UserName}}!
		{{range .Emails}}
			an email {{.|emailDeal}}
		{{end}}
		{{with .Friends}}
			{{range .}}
				my friend name is {{.FName}}
			{{end}}
		{{end}}`))
	p := Person{
		UserName: "Luffy",
		Emails: []string{
			"luffy.monkey@onepiece.com",
			"luffy.monkey@google.com",
			"luffy.monkey@yahoo.com",
		},
		Friends: []Friend{
			{"Sanji"},
			{"Zoro"},
			{"Chopper"},
		},

	}
	err := t.Execute(os.Stdout, p)
	if err != nil {
		log.Fatal(err)
	}
}
