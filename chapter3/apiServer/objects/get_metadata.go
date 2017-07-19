package objects

import (
	"lib/es"
)

func getMetaData(name string, version int) (es.Metadata, error) {
	if version == -1 {
		return es.SearchLatestVersion(name)
	}
	return es.GetMetadata(name, version)
}
