package tablenames

const (
	Files          = "files"
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
