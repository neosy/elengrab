package indexh

type IndexHandlers struct {
	assetsDir string
}

func NewIndexHandlers(assetsDir string) *IndexHandlers {
	return &IndexHandlers{
		assetsDir: assetsDir,
	}
}
