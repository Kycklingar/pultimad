package main

import (
	"html/template"
	"log"
	"os"
)

type index struct {
	Creators []creator
}

func main() {
	tmpl := template.New("")
	tmpl.Funcs(
		template.FuncMap{
				"unescape": func(x string) template.HTML { return template.HTML(x)},
		},
	)

	_, err := tmpl.ParseFiles("index.html", "posts.html")
	if err != nil {
		log.Fatal(err)
	}


	err = connect("user=postgres dbname=yp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	//var ind index
	//ind.Creators, err = archivedCreators()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = tmpl.ExecuteTemplate(os.Stdout, "index.html", ind)
	//if err != nil {
	//	log.Fatal(err)
	//}

	posts, err := posts()
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "posts.html", posts)
	if err != nil {
		log.Fatal(err)
	}

}
