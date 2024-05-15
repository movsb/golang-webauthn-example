package main

import (
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type LoginData struct {
	Name string
}

type Admin struct {
	prefix    string
	wa        *WebAuthn
	store     *Store
	templates map[string]*template.Template
}

func NewAdmin(store *Store, wa *WebAuthn, prefix string) *Admin {
	a := &Admin{
		store:     store,
		prefix:    prefix,
		wa:        wa,
		templates: make(map[string]*template.Template),
	}

	for _, f := range []string{`login.html`, `profile.html`} {
		a.templates[f] = template.Must(template.ParseFiles(f))
	}

	return a
}

func (a *Admin) Handler() http.Handler {
	m := http.NewServeMux()

	m.Handle(`GET /{$}`, a.requireLogin(a.getRoot))
	m.Handle(`/`, http.FileServer(http.Dir(".")))

	m.HandleFunc(`GET /login`, a.getLogin)
	m.HandleFunc(`GET /logout`, a.getLogout)
	m.HandleFunc(`POST /register`, a.postRegister)

	m.Handle(`GET /profile`, a.requireLogin(a.getProfile))

	const webAuthnPrefix = `/login/webauthn/`
	m.Handle(webAuthnPrefix, a.wa.Handler(webAuthnPrefix))

	return http.StripPrefix(strings.TrimSuffix(a.prefix, "/"), m)
}

func (a *Admin) postRegister(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	// TODO validate email
	// TODO validate dup.
	u, err := a.store.AddNewUser(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.store.MakeCookie(u, w, r)
	http.Redirect(w, r, a.prefixed(`profile`), http.StatusFound)
}

func (a *Admin) prefixed(s string) string {
	return path.Join(a.prefix, s)
}

func (a *Admin) redirectToLogin(w http.ResponseWriter, r *http.Request, to string) {
	args := url.Values{}
	args.Set(`u`, to)

	u, err := url.Parse(a.prefixed(`/login`))
	if err != nil {
		panic(err)
	}

	u.RawQuery = args.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (a *Admin) requireLogin(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := a.store.AuthRequest(r); user == nil {
			a.redirectToLogin(w, r, r.RequestURI)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (a *Admin) getRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.prefixed(`/profile`), http.StatusFound)
}

func (a *Admin) executeTemplate(w io.Writer, name string, data any) {
	t := a.templates[name]
	if t == nil {
		panic(`No such template：` + name)
	}
	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func (a *Admin) getLogin(w http.ResponseWriter, r *http.Request) {
	if a.store.AuthRequest(r) != nil {
		to := a.prefixed(`/profile`)
		if u := r.URL.Query().Get(`u`); u != "" {
			to = u
		}
		http.Redirect(w, r, to, http.StatusFound)
		return
	}

	d := LoginData{
		Name: "Golang WebAuthn Example",
	}
	a.executeTemplate(w, `login.html`, &d)
}

func (a *Admin) getLogout(w http.ResponseWriter, r *http.Request) {
	a.store.RemoveCookie(w, r)
	http.Redirect(w, r, a.prefixed(`/login`), http.StatusFound)
}

type ProfileData struct {
	User *User
}

func (d *ProfileData) PublicKeys() []string {
	ss := make([]string, 0, len(d.User.WebAuthnCredentials()))
	for _, c := range d.User.WebAuthnCredentials() {
		ss = append(ss, base64.RawURLEncoding.EncodeToString(c.ID))
	}
	return ss
}

func (a *Admin) getProfile(w http.ResponseWriter, r *http.Request) {
	d := &ProfileData{
		User: a.store.AuthRequest(r),
	}
	a.executeTemplate(w, `profile.html`, &d)
}
