package post

import (
	"database/sql"
	"net/http"
	"strconv"
)

type Page struct {
	ID    int
	Title string
	Desc  string
	Body  string
}

func GetPage(db *sql.DB, id int) (*Page, error) {
	var err error

	row := db.QueryRow("SELECT title, desc, body FROM posts WHERE id IS ?", id)
	var p *Page = new(Page)
	p.ID = id
	err = row.Scan(&p.Title, &p.Desc, &p.Body)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func GetAllPages(db *sql.DB) (*[]*Page, error) {
	var err error

	rows, err := db.Query("SELECT id, title, desc FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages *[]*Page = new([]*Page)

	for rows.Next() != false {
		var page *Page = new(Page)
		err = rows.Scan(&page.ID, &page.Title, &page.Desc)
		if err != nil {
			return nil, err
		}
		*pages = append(*pages, page)
	}

	return pages, nil
}

func (p *Page) LoadForm(r *http.Request) error {
	var err error

	p.ID, err = strconv.Atoi(r.URL.Path[len("/save/"):])
	if err != nil {
		return err
	}

	err = r.ParseForm()
	if err != nil {
		return err
	}

	p.Title = r.PostFormValue("post-title")
	p.Desc = r.PostFormValue("post-desc")
	p.Body = r.PostFormValue("post-body")

	return nil
}

func (p *Page) Save(db *sql.DB) error {
	var err error

	_, err = db.Exec("UPDATE posts SET title = ?, desc = ?, body = ? WHERE id IS ?", p.Title, p.Desc, p.Body, p.ID)
	if err != nil {
		return err
	}

	return nil
}
