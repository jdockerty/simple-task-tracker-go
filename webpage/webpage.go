package webpage

import (
	"html/template"
	"net/http"
)


func PopulateTemplate(){
	template := template.Must(template.ParseFiles(`webpage\tasks.html`))
	MyString := "hi"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, MyString)
	})

	http.ListenAndServe(":8080", nil)
	}


