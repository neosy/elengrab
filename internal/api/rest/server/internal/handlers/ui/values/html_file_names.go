package uivalues

// Pages
var (
	PageIndex = page{
		key:      "index",
		fileName: "index.html",
	}
	PageAuthRegister = page{
		key:      "auth-register",
		fileName: "auth-register.html",
	}
	PageAuthLogin = page{
		key:      "auth-login",
		fileName: "auth-login.html",
	}
	PageWatch = page{
		key:      "watch",
		fileName: "watch.html",
	}
	PageError = page{
		key:      "error",
		fileName: "error.html",
	}
)

const (
	ComponentAccountMenuContentKey = "account-menu-content"
	ComponentRowMenuContentKey     = "row-menu-content"

	ComponentResultRowsKey              = "result-rows"
	ComponentResultNewRowKey            = "result-new-row"
	ComponentResultRowStatusKey         = "result-row-status"
	ComponentResultLoadHistory          = "result-row-load-history"
	ComponentResultShouldLoadHistoryKey = "result-row-should-load-history"
	ComponentResultProgressKey          = "result-progress"
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
