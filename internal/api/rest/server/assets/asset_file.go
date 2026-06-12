package assets

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/assetx"
)

func FileFromData(name string, data []byte) *dtypes.AssetFile {
	hash := HashFromData(data)
	file := dtypes.NewAssetFile(name, hash, data)
	return file
}

func HashFromData(data []byte) string {
	return assetx.AssetFingerprintHex32(data)
}
