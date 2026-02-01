package ffmpeginfo

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
