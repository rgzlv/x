package main

import (
	"crypto/rand"
	"database/sql"
	"dtla/internal/post"
	"dtla/internal/util"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type handlerData struct {
	sstate    *ServerState
	cleanPath string
	tmpl      util.TmplData
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *handlerData), sstate *ServerState, wantAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		var hd handlerData

		hd.tmpl.Auth.ASDefault = util.ASDefault
		hd.tmpl.Auth.ASError = util.ASError
		hd.tmpl.Auth.ASOk = util.ASOk

		hd.cleanPath, err = util.CleanPath(r.URL.Path)
		if err != nil {
			util.LogHTTPError(w, err)
			return
		}
		hd.sstate = sstate
		hd.tmpl.URLPath = r.URL.Path

		if !wantAuth {
			fn(w, r, &hd)
			return
		}

		hd.tmpl.Auth.Status = util.ASDefault

		sidCookie, err := r.Cookie("sid")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				fn(w, r, &hd)
				return
			} else {
				handlerAuthError(fn, w, r, &hd, err)
				return
			}
		}
		cookieSidBytes, err := hex.DecodeString(sidCookie.Value)
		if err != nil {
			handlerAuthError(fn, w, r, &hd, err)
			return
		}

		err = getUserData(r, hd.sstate.DB, &hd.tmpl.Auth)
		if err != nil {
			if errors.Is(err, util.ErrSessionExpired) {
				util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, err.Error())
				return
			}
			handlerAuthError(fn, w, r, &hd, err)
			return
		}

		err = bcrypt.CompareHashAndPassword(hd.tmpl.Auth.SID, cookieSidBytes)
		if err != nil {
			handlerAuthError(fn, w, r, &hd, err)
			return
		}

		hd.tmpl.Auth.Status = util.ASOk
		fn(w, r, &hd)
	}
}

// Call handler with Auth.Status and Auth.Error set
// so they can be checked in handler and templates
func handlerAuthError(fn func(http.ResponseWriter, *http.Request, *handlerData), w http.ResponseWriter, r *http.Request, hd *handlerData, err error) {
	hd.tmpl.Auth.Status = util.ASError
	hd.tmpl.Auth.Error = err.Error()
	util.LogError(err.Error())
	fn(w, r, hd)
}

func rootHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	err = util.ExecuteTemplate(w, r, hd.cleanPath, *hd.sstate.TmplDir, &hd.tmpl)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}
}

func viewAllHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	hd.tmpl.Data, err = post.GetAllPages(hd.sstate.DB)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}

	err = util.ExecuteTemplate(w, r, "view-all.html", *hd.sstate.TmplDir, &hd.tmpl)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	pageID, err := strconv.Atoi(r.URL.Path[len("/view/"):])
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}

	hd.tmpl.Data, err = post.GetPage(hd.sstate.DB, pageID)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}

	hd.tmpl.URLPath = "/view/"
	err = util.ExecuteTemplateHTML(w, r, "view.html", *hd.sstate.TmplDir, &hd.tmpl)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	if hd.tmpl.Auth.Status != util.ASOk {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, "Nav atļauts rediģēt ieteikumu")
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/edit/"):])
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}

	hd.tmpl.Data, err = post.GetPage(hd.sstate.DB, id)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}

	hd.tmpl.URLPath = "/edit/"
	err = util.ExecuteTemplate(w, r, "edit.html", *hd.sstate.TmplDir, &hd.tmpl)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	if hd.tmpl.Auth.Status != util.ASOk {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, "Nav atļauts rediģēt ieteikumu")
		return
	}

	var page post.Page

	err = page.LoadForm(r)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}

	err = page.Save(hd.sstate.DB)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}

	http.Redirect(w, r, "/view/"+strconv.Itoa(page.ID), http.StatusSeeOther)
}

func toolsHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	if r.URL.Path == "/tools/sockets" && runtime.GOOS == "windows" {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, "Nav atbalstīts uz Windows servera")
		return
	}

	filename := hd.cleanPath + ".html"
	hd.tmpl.URLPath = "/tools/"
	err = util.ExecuteTemplate(w, r, filename, *hd.sstate.TmplDir, &hd.tmpl)
	if err != nil {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, err.Error())
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	if hd.tmpl.Auth.Status != util.ASOk {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, "Nav atļauts dzēst ieteikumu")
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/delete/"):])
	if err != nil {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, err.Error())
		return
	}

	_, err = hd.sstate.DB.Exec("DELETE FROM posts WHERE id IS ?", id)
	if err != nil {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, err.Error())
		return
	}

	http.Redirect(w, r, "/view/", http.StatusTemporaryRedirect)
}

func newHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	if hd.tmpl.Auth.Status != util.ASOk {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, "Nav atļauts izveidot jaunu ieteikumu")
		return
	}

	res, err := hd.sstate.DB.Exec("INSERT INTO posts (title, desc, body) VALUES ('Virsraksts', 'Apraksts', '<p>Saturs</p>')")
	if err != nil {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, err.Error())
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		util.ExecuteTemplateError(w, r, *hd.sstate.TmplDir, &hd.tmpl, err.Error())
		return
	}

	http.Redirect(w, r, "/edit/"+strconv.FormatInt(id, 10), http.StatusTemporaryRedirect)
}

func loginHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	if hd.tmpl.Auth.Status == util.ASOk {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err = util.ExecuteTemplate(w, r, "login.html", *hd.sstate.TmplDir, &hd.tmpl)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}
}

func loginPostHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	err = r.ParseForm()
	if err != nil {
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, err)
		return
	}

	user := r.PostFormValue("login-name")
	pswd := r.PostFormValue("login-pswd")

	if user == "" || pswd == "" {
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, errors.New("Lietotājvārds un parole nedrīkst būt neaizpildīti"))
		return
	}

	row := hd.sstate.DB.QueryRow("SELECT id FROM users WHERE user IS ?", user)
	var id uint
	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, errors.New("Lietotājs neeksistē"))
			return
		}
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, err)
		return
	}

	row = hd.sstate.DB.QueryRow("SELECT pswd FROM users WHERE id IS ?", id)
	var pswdHashDB string
	err = row.Scan(&pswdHashDB)
	if err != nil {
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, err)
		return
	}

	pswdHashDBBytes := []byte(pswdHashDB)
	pswdBytes := []byte(pswd)

	err = bcrypt.CompareHashAndPassword(pswdHashDBBytes, pswdBytes)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, errors.New("Nepareiza parole"))
			return
		}
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, err)
		return
	}

	var sidBytes []byte = make([]byte, 72)
	_, err = rand.Read(sidBytes)
	if err != nil {
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, err)
		return
	}

	sidHashBytes, err := bcrypt.GenerateFromPassword(sidBytes, 6)
	if err != nil {
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, err)
		return
	}

	_, err = hd.sstate.DB.Exec("UPDATE users SET sid = ? WHERE id IS ?", string(sidHashBytes), id)
	if err != nil {
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, err)
		return
	}

	hd.tmpl.Auth.SStart = time.Now()
	hd.tmpl.Auth.SAge = time.Second * 300

	_, err = hd.sstate.DB.Exec("UPDATE users SET sstart = ?, sage = ? WHERE id IS ?", hd.tmpl.Auth.SStart.Unix(), int(hd.tmpl.Auth.SAge.Seconds()), id)
	if err != nil {
		util.ExecuteTemplateLoginWithError(w, r, &hd.tmpl, *hd.sstate.TmplDir, err)
		return
	}

	// When checking sid in requests, hex will be transformed to []byte with hex.DecodeString()
	// and compared using bcrypt.CompareHashAndPassword()
	sidCookie := http.Cookie{
		Name:   "sid",
		MaxAge: int(hd.tmpl.Auth.SAge.Seconds()),
		Value:  hex.EncodeToString(sidBytes),
	}
	http.SetCookie(w, &sidCookie)

	// bcrypt.GenerateFromPassword() uses a random salt so to compare the sid in later requests
	// the id is needed to find the user for which to compare it against
	idCookie := http.Cookie{
		Name:   "id",
		MaxAge: int(hd.tmpl.Auth.SAge.Seconds()),
		Value:  fmt.Sprintf("%d", id),
	}
	http.SetCookie(w, &idCookie)

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("refresh", "1;url=/")
	hd.tmpl.Auth.Status = util.ASOk
	err = util.ExecuteTemplate(w, r, "login.html", *hd.sstate.TmplDir, &hd.tmpl)
	if err != nil {
		util.LogHTTPError(w, err)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	sidCookie := http.Cookie{
		Name:    "sid",
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, &sidCookie)

	idCookie := http.Cookie{
		Name:    "id",
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, &idCookie)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	var err error

	err = util.ServeFile(w, r, hd.cleanPath)
	if err != nil {
		util.LogHTTPError(w, err)
		return
	}
}

func licenseHandler(w http.ResponseWriter, r *http.Request, hd *handlerData) {
	http.ServeFile(w, r, "LICENSE")
}
