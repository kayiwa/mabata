package httpui

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/kayiwa/mabata/internal/auth"
	"github.com/kayiwa/mabata/internal/config"
	"github.com/kayiwa/mabata/internal/duck"
	"github.com/kayiwa/mabata/internal/models"
)

type Handlers struct {
	cfg  config.Config
	db   *sql.DB
	oidc *auth.OIDC
	tpl  *template.Template
}

type pageData struct {
	User     *models.User
	Queries  []string
	Selected string
	Headers  []string
	Rows     [][]string
	Error    string
}

func New(cfg config.Config, db *sql.DB, oidc *auth.OIDC) *Handlers {
	return &Handlers{
		cfg:  cfg,
		db:   db,
		oidc: oidc,
		tpl:  template.Must(template.New("page").Parse(pageTemplate)),
	}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/login", h.login)
	mux.HandleFunc("/auth/callback", h.callback)
	mux.HandleFunc("/logout", h.logout)
	mux.HandleFunc("/whoami", h.whoami)
}

func (h *Handlers) home(w http.ResponseWriter, r *http.Request) {
	user, _ := h.oidc.Session().Get(r)
	pd := pageData{User: user, Queries: duck.Names()}
	if user != nil {
		selected := r.URL.Query().Get("query")
		if selected == "" && len(pd.Queries) > 0 {
			selected = pd.Queries[0]
		}
		pd.Selected = selected
		if selected != "" {
			headers, rows, err := duck.Run(h.db, selected)
			if err != nil {
				pd.Error = err.Error()
			} else {
				pd.Headers = headers
				pd.Rows = rows
			}
		}
	}
	_ = h.tpl.Execute(w, pd)
}

func (h *Handlers) login(w http.ResponseWriter, r *http.Request) {
	url, err := h.oidc.AuthCodeURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handlers) callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	user, err := h.oidc.Exchange(r.Context(), state, code)
}
