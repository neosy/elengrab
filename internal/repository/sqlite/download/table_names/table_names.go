package tablenames

const (
	Files           = "media_downloads"
	DownloadTasks   = "download_tasks"
	ProcessingTasks = "processing_tasks"
	DataMigrations  = "data_migrations"
)

var tableNames = []string{
	Files,
	DownloadTasks,
	ProcessingTasks,
	DataMigrations,
}

func TableNames() []string {
	return tableNames
}
