package pages

import "encoding/json"

type Channel struct {
	EncodedChannelID string `json:"channelId"`

	Title    string `json:"title"`
	URL      string `json:"url"`
	ImageURL string `json:"imageUrl"`

	Ext ChannelExt `json:"ext"`
}

type ChannelExt struct {
	ExtID    string `json:"extId"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
	Username string `json:"username"`
}

type ChannelHeader struct {
	Channel
	Show bool `json:"show"`
}

func (c Channel) JSON() []byte {
	json, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	return json
}

func (c ChannelHeader) JSON() []byte {
	json, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	return json
}
