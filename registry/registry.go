package registry

type Template struct {
	ID          string
	Label       string
	Description string
	Repo        string
	DefaultRef  string
	Variables   []string
}

var templates = []Template{
	{
		ID:          "mars",
		Label:       "API + Alembic (recommended)",
		Description: "FastAPI backend with Alembic migrations, PostgreSQL, Redis, JWT auth, and optional Google OIDC",
		Repo:        "https://github.com/ayoub3bidi/mars",
		DefaultRef:  "v0.1.3",
		Variables:   []string{"project_name", "app_title", "api_version"},
	},
}

func List() []Template {
	return templates
}

func ByID(id string) (Template, bool) {
	for _, t := range templates {
		if t.ID == id {
			return t, true
		}
	}
	return Template{}, false
}
