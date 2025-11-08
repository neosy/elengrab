package statich

type StaticHandlers struct {
	assetsDir string
}

func NewStaticHandlers(assetsDir string) *StaticHandlers {
	return &StaticHandlers{
		assetsDir: assetsDir,
	}
}
