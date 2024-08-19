package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"

	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"snippetbox.formme.net/internal/models"
)

// Create a new application type

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, 200, "home.tmpl", data)
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, 200, "view.tmpl", data)
}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		w.Header().Set("Allow", http.MethodPost)
	// 		app.clientError(w, http.StatusMethodNotAllowed)
	// 		return
	// 	}
	// 	title := "O snail"
	// 	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// 	expires := 7
	// 	id, err := app.snippets.Insert(title, content, expires)
	// 	if err != nil {
	// 		app.serverError(w, err)
	// 		return
	// 	}
	// 	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.tmpl", data)
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// title := "O snail"
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// expires := 7
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//Get data from client
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//Data validation
	fieldErrors := make(map[string]string)
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal 1,7 or 365"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}
	//Insert into database
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
