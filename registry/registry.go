package registry

type Template struct {
	ID          string
	Label       string
	Description string
	Repo        string
	DefaultRef  string
	Variables   []string
	Features    []string
	Kind        string // "api" or "docs"
}

var templates = []Template{
	{
		ID:          "mars",
		Label:       "API + Alembic (recommended)",
		Description: "FastAPI backend with Alembic migrations, PostgreSQL, Redis, JWT auth, and optional Google OIDC",
		Repo:        "https://github.com/ayoub3bidi/mars",
		DefaultRef:  "v0.1.4",
		Variables:   []string{"project_name", "app_title", "api_version"},
		Features:    []string{"oidc"},
		Kind:        "api",
	},
	{
		ID:          "mercury",
		Label:       "API + Flyway",
		Description: "FastAPI backend with Flyway SQL migrations, PostgreSQL, Redis, JWT auth, and optional Google OIDC",
		Repo:        "https://github.com/ayoub3bidi/mercury",
		DefaultRef:  "v0.4.1",
		Variables:   []string{"project_name", "app_title", "api_version"},
		Features:    []string{"oidc"},
		Kind:        "api",
	},
	{
		ID:          "venus",
		Label:       "Docs site",
		Description: "Docusaurus documentation site with Tailwind CSS, optional i18n and local search",
		Repo:        "https://github.com/ayoub3bidi/venus",
		DefaultRef:  "v1.1.2",
		Variables:   []string{"project_name", "site_title"},
		Features:    []string{"i18n", "search"},
		Kind:        "docs",
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

func APIByMigration(useAlembic bool) (Template, bool) {
	if useAlembic {
		return ByID("mars")
	}
	return ByID("mercury")
}
