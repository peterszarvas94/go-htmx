package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Name struct {
	Name string
}

type Item struct {
	Id   int
	Text string
}

type ItemsData struct {
	Items []Item
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/items/", itemshandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func getData() []Item {
	return []Item{
		{Id: 1, Text: "Item 1"},
		{Id: 2, Text: "Item 2"},
		{Id: 3, Text: "Item 3"},
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed index", http.StatusMethodNotAllowed)
		return
	}

	index := "templates/index.html"
	items := "templates/items.html"
	item := "templates/item.html"
	tmpl, err := template.ParseFiles(index, items, item)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	data := ItemsData{
		Items: getData(),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func itemshandler(w http.ResponseWriter, r *http.Request) {
	rawPath := strings.Split(r.URL.Path, "/")
	var path []string
	for _, element := range rawPath {
		if element != "" {
			path = append(path, element)
		}
	}
	item := "templates/item.html"
	items := "templates/items.html"

	// get all items
	if len(path) == 1 && r.Method == "GET" {
		tmpl, err := template.ParseFiles(items, item)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		data := ItemsData{
			Items: getData(),
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// add new item
	if len(path) == 1 && r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		text := r.FormValue("text")
		fmt.Println("Text:", text)

		tmpl, err := template.ParseFiles(item)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		data := Item{
			Text: text,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// get item by id
	if len(path) == 2 && r.Method == "GET" {
		id, err := strconv.Atoi(path[1])
		if err != nil {
			http.Error(w, "Invalid Id", http.StatusInternalServerError)
			return
		}

		if id > len(getData()) || id < 1 {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles(item)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		data := Item{
			Text: getData()[id-1].Text,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// delete item by id
	if len(path) == 2 && r.Method == "DELETE" {
		id, err := strconv.Atoi(path[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: delete item from database
		fmt.Println("Delete item:", id)

		return
	}

}
