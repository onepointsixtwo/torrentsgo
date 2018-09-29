package parser

import (
	"github.com/onepointsixtwo/torrentsgo/bencoding"
	"github.com/onepointsixtwo/torrentsgo/model"
)

func ParseMetaInfo(reader io.Reader) (*model.MetaInfo, error) {
	decoded, err := bencoding.DecodeBencoding(reader)
	if err != nil {
		return nil, err
	}

	return parseMetaInfoFromDecodedData(decoded)
}

func parseMetaInfoFromDecodedData(data map[string]interface{}) (*model.MetaInfo, error) {
	//TODO: parse the meta info struct from the data!
	return nil, nil
}
