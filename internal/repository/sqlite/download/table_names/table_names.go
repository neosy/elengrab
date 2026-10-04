package tablenames

const (
	Files          = "media_downloads"
	DownloadTasks  = "download_tasks"
	DataMigrations = "data_migrations"
)

var tableNames = []string{
	Files,
	DownloadTasks,
	DataMigrations,
}

func TableNames() []string {
	return tableNames
}
