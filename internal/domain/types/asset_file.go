package dtypes

import (
	"path/filepath"
	"slices"
	"strings"
)

type AssetFile struct {
	FilePath string
	FileName string
	Hash     string
	Raw      []byte
}

func NewAssetFile(filePath string, hash string, data []byte) *AssetFile {
	fileName := filepath.Base(filePath)
	return &AssetFile{
		FilePath: filePath,
		FileName: fileName,
		Hash:     hash,
		Raw:      data,
	}
}

func (f *AssetFile) Copy() *AssetFile {
	if f == nil {
		return nil
	}

	file := new(*f)

	file.Raw = slices.Clone(f.Raw)

	return file
}

func (f *AssetFile) FileNameWithHash() string {
	if f.Hash == "" {
		return f.FileName
	}

	ext := filepath.Ext(f.FileName)
	name := strings.TrimSuffix(f.FileName, ext)

	return name + "." + f.Hash + ext
}
