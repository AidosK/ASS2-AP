package main

import (
	"aidoskanatbay.net/snippetbox/pkg/forms"
	"aidoskanatbay.net/snippetbox/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle form submission

		// Access form values
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		// Do something with the form data, e.g., send an email, save to database, etc.

		// For now, let's just print the form data
		fmt.Printf("Name: %s\nEmail: %s\nMessage: %s\n", name, email, message)

		// Redirect to a thank you page or display a success message
		http.Redirect(w, r, "/thank-you", http.StatusSeeOther)
		return
	}

	// Render the contact form page for GET requests
	app.render(w, r, "contact.page.tmpl", nil)
}
func (app *application) students(w http.ResponseWriter, r *http.Request) {
	// Retrieve snippets with the category "students"
	snippets, err := app.snippets.Student("students")
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "student.page.tmpl", &templateData{
		Snippets: snippets,
	})
}

func (app *application) staff(w http.ResponseWriter, r *http.Request) {
	// Retrieve snippets with the category "students"
	snippets, err := app.snippets.Staff("staff")
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "staff.page.tmpl", &templateData{
		Snippets: snippets,
	})
}
func (app *application) applicant(w http.ResponseWriter, r *http.Request) {
	// Retrieve snippets with the category "students"
	snippets, err := app.snippets.Applicant("staff")
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "applicant.page.tmpl", &templateData{
		Snippets: snippets,
	})
}

func (app *application) researcher(w http.ResponseWriter, r *http.Request) {
	// Retrieve snippets with the category "students"
	snippets, err := app.snippets.Researcher("researcher")
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "researcher.page.tmpl", &templateData{
		Snippets: snippets,
	})
}
