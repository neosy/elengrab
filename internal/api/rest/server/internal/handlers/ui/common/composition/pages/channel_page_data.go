package pages

import (
	"encoding/json"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
)

// Channel page
type ChannelPageData struct {
	BasePaths  paths.HttpPaths
	BaseValues baseValues
	Paths      PagePaths
	Values     ChannelPageValues
	Extra      map[string]any
}

type ChannelPageValues struct {
	PagesListValues

	ChannelHeader ChannelHeader
}

type ChannelHeader struct {
	Channel
	Show bool `json:"show"`
}

func (c ChannelHeader) JSON() []byte {
	json, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	return json
}
