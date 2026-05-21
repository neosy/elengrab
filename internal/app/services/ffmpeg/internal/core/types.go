package core

type videoInfoSetter struct {
	codec      bool
	resolution bool
	bitrate    bool
}

type audioInfoSetter struct {
	codec      bool
	sampleRate bool
	bitrate    bool
}

type ffprobeStream struct {
	CodecType string `json:"codec_type"`
	CodecName string `json:"codec_name"`

	Width  int `json:"width"`
	Height int `json:"height"`

	SampleRate string `json:"sample_rate"`
	BitRate    string `json:"bit_rate"`

	Duration string `json:"duration"`
}
