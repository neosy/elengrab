package pages

// Pages
var (
	IndexPage = page{
		key:      "index-page",
		fileName: "index.html",
	}
	AuthRegisterPage = page{
		key:      "auth-register-page",
		fileName: "auth-register.html",
	}
	AuthLoginPage = page{
		key:      "auth-login-page",
		fileName: "auth-login.html",
	}
	WatchPage = page{
		key:      "watch-page",
		fileName: "watch.html",
	}
	AdminPage = page{
		key:      "admin-page",
		fileName: "admin.html",
	}
	ErrorPage = page{
		key:      "error-page",
		fileName: "error.html",
	}
)

type page struct {
	key      string
	fileName string
}

func (p *page) Key() string {
	return p.key
}

func (p *page) FileName() string {
	return p.fileName
}
