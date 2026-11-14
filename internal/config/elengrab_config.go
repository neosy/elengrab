package iconfig

import "time"

// ElengrabConfig contains the Elengrab application configuration.
type ElengrabConfig struct {
	// Application operating mode.
	Mode string `env:"MODE" envDefault:"public"`
	// Enables demo mode with restricted functionality.
	DemoMode bool `env:"DEMO_MODE" envDefault:"false"`

	// Shows developer information and donation links in the About section.
	AboutDeveloper ElengrabAboutDeveloper `envPrefix:"ABOUT_DEVELOPER_"`

	// Base URL used to build absolute application links.
	BaseURL string `env:"BASE_URL" envDefault:""`

	// Prefix used for generated short links.
	ShortLinkPrefix string `env:"SHORT_LINK_PREFIX" envDefault:"/s"`
	// Length of generated short link identifiers.
	ShortLinkLength uint8 `env:"SHORT_LINK_LENGTH" envDefault:"6"`
	// Lifetime of generated short links in days.
	ShortLinkTTLDays uint16 `env:"SHORT_LINK_TTL_DAYS" envDefault:"180"`

	// AdminLogin is used to create the default administrator at first startup
	AdminLogin string `env:"ADMIN_LOGIN" envDefault:""`
	// AdminPassword is used to create the default administrator at first startup
	AdminPassword string `env:"ADMIN_PASSWORD" envDefault:""`

	// Directory containing downloader binaries such as yt-dlp.
	DownloaderBinDir string `env:"DOWNLOADER_BIN_DIR" envDefault:"/usr/local/bin"`

	// Root directory for application data.
	RootDir string `env:"ROOT_DIR" envDefault:""`
	// Directory containing application assets.
	AssetsDir string `env:"ASSETS_DIR" envDefault:"assets"`
	// Directory containing media files (e.g. thumbnails).
	MediaDir string `env:"MEDIA_DIR" envDefault:"media"`
	// Directory containing downloaded files.
	DownloadsDir string `env:"DOWNLOADS_DIR" envDefault:"downloads"`
	// Directory containing cookies used by the downloader.
	CookiesDir string `env:"COOKIES_DIR" envDefault:"cookies"`

	// Configuration for uploading media files.
	Upload ElengrabUploadConfig `env:"UPLOAD" envDefault:""`

	// Maximum number of concurrent media downloads.
	DownloadWorkers uint32 `env:"DOWNLOAD_WORKERS" envDefault:"3"`
	// Maximum number of concurrent media processing tasks.
	MediaProcessingWorkers uint32 `env:"MEDIA_PROCESSING_WORKERS" envDefault:"3"`
	// Maximum number of concurrent background operations.
	OperationWorkers uint32 `env:"OPERATION_WORKERS" envDefault:"5"`

	// Scope used when checking media uniqueness before deleting duplicates.
	DeleteDuplicatesUniquenessScope string `env:"DELETE_DUPLICATES_UNIQUENESS_SCOPE" envDefault:"per_user"`

	// Enables the use of cookies when downloading media.
	AllowCookies bool `env:"ALLOW_COOKIES" envDefault:"false"`

	// Logging configuration.
	Logging ElengrabLoggingConfig `envPrefix:"LOGGING_"`

	// Behavior configuration.
	Behavior ElengrabBehaviorConfig `envPrefix:"BEHAVIOR_"`

	// Maintenance task configuration.
	Maintenance ElengrabMaintenanceConfig `envPrefix:"MAINTENANCE_"`
}

// ElengrabAboutDeveloper configures the visibility of developer information and donation options.
type ElengrabAboutDeveloper struct {
	// ShowInfo controls whether developer information is displayed.
	ShowInfo bool `env:"SHOW_INFO" envDefault:"true"`

	// ShowDonationLink controls whether the donation link is displayed.
	ShowDonationLink bool `env:"SHOW_DONATION_LINK" envDefault:"false"`
}

// ElengrabUploadConfig holds the settings for uploading media files from the user's device.
type ElengrabUploadConfig struct {
	// How long a temporary uploaded file waits for the user's confirmation before it is removed.
	TempTTLMinutes uint16 `env:"TEMP_TTL_MINUTES" envDefault:"120"`

	// Maximum size of an uploaded file, in MiB.
	MaxFileSizeMiB uint32 `env:"MAX_FILE_SIZE_MIB" envDefault:"10240"`
}

// TTL returns the lifetime of an uploaded file waiting for confirmation.
func (c ElengrabUploadConfig) TempTTL() time.Duration {
	return time.Duration(c.TempTTLMinutes) * time.Minute
}

// MaxFileSize returns the maximum size of an uploaded file, in bytes.
func (c ElengrabUploadConfig) MaxFileSize() int64 {
	return int64(c.MaxFileSizeMiB) << 20 // * 1024 * 1024
}

// ElengrabLoggingConfig contains application logging settings.
type ElengrabLoggingConfig struct {
	// Enables logging during application startup.
	StartupEnabled bool `env:"STARTUP_ENABLED" envDefault:"true"`

	// Enables logging while the application is running.
	RuntimeEnabled bool `env:"RUNTIME_ENABLED" envDefault:"false"`

	// Directory containing application log files.
	Dir string `env:"DIR" envDefault:"logs"`
}

// ElengrabMediaConfig contains media configuration.
type ElengrabBehaviorConfig struct {
	// Visibility assigned to new media.
	// Supported values: private, authenticated, public.
	// If not set, visibility is determined by the application mode.
	NewMediaVisibility string `env:"NEW_MEDIA_VISIBILITY" envDefault:""`
}

// ElengrabMaintenanceConfig contains configuration for periodic maintenance tasks.
type ElengrabMaintenanceConfig struct {
	// Interval for calculating and updating hashes of downloaded files.
	// File hashes are used to identify duplicate files.
	// Defaults to 8h if not specified.
	UpdateHashInterval time.Duration `env:"UPDATE_HASH_INTERVAL" envDefault:"8h"`

	// Interval for deleting duplicate media files.
	// Set to 0 to disable automatic deletion.
	DeleteDuplicatesInterval time.Duration `env:"DELETE_DUPLICATES_INTERVAL" envDefault:"12h"`

	// Interval for deleting database records whose downloaded files are missing.
	// Set to 0 to disable automatic deletion.
	DeleteMissingMediaDownloadsInterval time.Duration `env:"DELETE_MISSING_DOWNLOADS_INTERVAL" envDefault:"12h"`

	// Interval for deleting failed download records.
	// Set to 0 to disable automatic deletion.
	DeleteFailedMediaDownloadsInterval time.Duration `env:"DELETE_FAILED_DOWNLOADS_INTERVAL" envDefault:"24h"`

	// Enables moving unmatched files to the .lost directory during periodic missing media cleanup.
	// This operation runs according to DeleteMissingMediaDownloadsInterval.
	MoveUnmatchedFilesEnabled bool `env:"ENABLE_MOVE_UNMATCHED_FILES" envDefault:"false"`

	// Number of the latest database backups to keep.
	// If set to 0, old backup files are not deleted.
	DatabaseBackupsKeep int `env:"DATABASE_BACKUPS_KEEP" envDefault:"7"`
}
