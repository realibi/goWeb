package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"html/template"
	"net/http"
	"strconv"
)


type snippet struct{
	id int
	title string
	content string
	created pgtype.Date
	expires pgtype.Date
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound() helper
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	connStr := "user=postgres password=pass dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from snippets where id=" + strconv.Itoa(id))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next(){
		p := snippet{}
		err := rows.Scan(&p.id, &p.title, &p.content, &p.created, &p.expires)
		if err != nil{
			fmt.Println(err)
		}
		result := snippet{id: p.id, title: p.title, content: p.content, created: p.created, expires: p.expires}
		resultJson, err := json.Marshal(result)
		if err != nil{
			app.serverError(w, err)
		}
		fmt.Println(resultJson)

		resultString := "title: " + result.title + "\n" + "content: " + result.content + "\n"

		w.Write([]byte(resultString))
	} else{
		http.NotFound(w, r)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}