package dto

type Format struct {
	FormatID   string  `json:"format_id"`
	Ext        string  `json:"ext"`
	Height     int     `json:"height"`
	Width      int     `json:"width"`
	FPS        float32 `json:"fps"`
	Format     string  `json:"format"`
	FormatNote string  `json:"format_note"`
	Resolution string  `json:"resolution"`
	VCodec     string  `json:"vcodec"`
	ACodec     string  `json:"acodec"`
	Vbr        float32 `json:"vbr"`
	Abr        float32 `json:"abr"`
	Asr        *int    `json:"asr"`
	Filesize   *int    `json:"filesize"`
}

type YouTubeInfo struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Formats []Format `json:"formats"`
}
