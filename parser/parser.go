package parser

import (
	"fmt"
	"github.com/onepointsixtwo/torrentsgo/bencoding"
	"github.com/onepointsixtwo/torrentsgo/model"
	"io"
	"net/url"
	"reflect"
	"time"
)

/*
	NOTE: Any fields which do not have a parser func returning an error _as well as_ a value are optional fields
	in Bittorrent Metainfo.
*/

// Public parser func

func ParseMetaInfo(reader io.Reader) (*model.MetaInfo, error) {
	decoded, err := bencoding.DecodeBencoding(reader)
	if err != nil {
		return nil, fmt.Errorf("Unable to decode bencoded data - %v", err)
	}

	return parseMetaInfoFromDecodedData(decoded)
}

// MetaInfo parsing

func parseMetaInfoFromDecodedData(data *model.OrderedMap) (*model.MetaInfo, error) {
	announceUrls, announceUrlsError := parseAnnounceUrlsFromDecodedData(data)
	if announceUrlsError != nil {
		return nil, announceUrlsError
	}
	creationDate := parseCreationDateFromDecodedData(data)
	comment := parseCommentFromDecodedData(data)
	createdBy := parseCreatedByFromDecodedData(data)
	encoding := parseEncodingFromDecodedData(data)
	info, err := parseInfoFromDecodedData(data)
	if err != nil {
		return nil, err
	}

	return model.NewMetaInfo(announceUrls, creationDate, comment, createdBy, encoding, info), nil
}

func parseAnnounceUrlsFromDecodedData(data *model.OrderedMap) ([]*url.URL, error) {
	// NOTE: this is currently not supporting the newer file extension of 'announce-list' in addition to announce,
	// but intentionally returns array of URLs so this can easily be supported later
	announce, err := readStringValueFromMap(data, "announce")
	if err != nil {
		return nil, err
	}

	u, parseUrlErr := url.Parse(announce)
	if parseUrlErr != nil {
		return nil, err
	}

	urls := make([]*url.URL, 1)
	urls[0] = u
	return urls, nil
}

func parseCreationDateFromDecodedData(data *model.OrderedMap) time.Time {
	timestamp, _ := readIntValueFromMap(data, "creation date")
	return time.Unix(int64(timestamp), 0)
}

func parseCommentFromDecodedData(data *model.OrderedMap) string {
	str, _ := readStringValueFromMap(data, "comment")
	return str
}

func parseCreatedByFromDecodedData(data *model.OrderedMap) string {
	str, _ := readStringValueFromMap(data, "created by")
	return str
}

func parseEncodingFromDecodedData(data *model.OrderedMap) string {
	str, _ := readStringValueFromMap(data, "encoding")
	return str
}

// Info parsing

func parseInfoFromDecodedData(data *model.OrderedMap) (*model.Info, error) {
	infoData, err := readDictionaryValueFromMap(data, "info")
	if err != nil {
		return nil, err
	}

	pieceLength, errPieceLength := parsePieceLengthFromDecodedInfoData(infoData)
	if errPieceLength != nil {
		return nil, errPieceLength
	}
	pieces, errPieces := parsePiecesDataFromDecodedInfoData(infoData)
	if errPieces != nil {
		return nil, errPieces
	}
	private := parsePrivateFromDecodedInfoData(infoData)
	files, filesErr := parseFilesFromDecodedInfoData(infoData)
	if filesErr != nil {
		return nil, filesErr
	}
	directoryName, directoryNameError := parseDirectoryNameFromDecodedInfoData(infoData)
	if directoryNameError != nil {
		return nil, directoryNameError
	}

	return model.NewInfo(pieceLength, pieces, private, files, directoryName), nil
}

func parsePieceLengthFromDecodedInfoData(infoData *model.OrderedMap) (int, error) {
	return readIntValueFromMap(infoData, "piece length")
}

func parsePiecesDataFromDecodedInfoData(infoData *model.OrderedMap) ([]byte, error) {
	return readBytesValueFromMap(infoData, "pieces")
}

func parsePrivateFromDecodedInfoData(infoData *model.OrderedMap) int {
	value, _ := readIntValueFromMap(infoData, "private")
	return value
}

func parseFilesFromDecodedInfoData(infoData *model.OrderedMap) ([]*model.File, error) {
	if isMultiFileMode(infoData) {
		return parseMultiFileModeFilesFromDecodedInfoData(infoData)
	}
	return parseSingleFileModeFilesFromDecodedInfoData(infoData)
}

func parseSingleFileModeFilesFromDecodedInfoData(infoData *model.OrderedMap) ([]*model.File, error) {
	file, err := parseFileFromOuterMap(infoData)
	if err != nil {
		return nil, err
	}

	files := make([]*model.File, 1)
	files[0] = file
	return files, nil
}

func parseMultiFileModeFilesFromDecodedInfoData(infoData *model.OrderedMap) ([]*model.File, error) {
	filesList, err := readListFromMap(infoData, "files")
	if err != nil {
		return nil, err
	}

	filesLength := len(filesList)
	files := make([]*model.File, 0)
	for i := 0; i < filesLength; i++ {
		maybeDict := filesList[i]
		dict, ok := maybeDict.(*model.OrderedMap)
		if !ok {
			continue
		}

		file, err := parseFileFromFilesMap(dict)
		if err == nil {
			files = append(files, file)
		}
	}

	if len(files) > 0 {
		return files, nil
	} else {
		return nil, fmt.Errorf("Expected at least one file when parsing multi-files")
	}
}

func parseFileFromOuterMap(mp *model.OrderedMap) (*model.File, error) {
	fileName, err := readStringValueFromMap(mp, "name")
	if err != nil {
		return nil, err
	}
	length, err2 := readIntValueFromMap(mp, "length")
	if err2 != nil {
		return nil, err2
	}
	md5Sum, _ := readStringValueFromMap(mp, "md5sum")

	return model.NewFile(fileName, length, md5Sum), nil
}

func parseFileFromFilesMap(mp *model.OrderedMap) (*model.File, error) {
	pathList, pathErr := readListFromMap(mp, "path")
	if pathErr != nil {
		return nil, pathErr
	}
	path, pathStrErr := getFilePathFromList(pathList)
	if pathStrErr != nil {
		return nil, pathErr
	}
	length, err2 := readIntValueFromMap(mp, "length")
	if err2 != nil {
		return nil, err2
	}
	md5Sum, _ := readStringValueFromMap(mp, "md5sum")

	return model.NewFile(path, length, md5Sum), nil
}

func getFilePathFromList(list []interface{}) (string, error) {
	path := ""

	pathLength := len(list)
	for i := 0; i < pathLength; i++ {
		if i != 0 {
			path = path + "/"
		}

		value := list[i]
		strValue, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("Value %v cannot be represented as string", value)
		}

		path = path + strValue
	}

	return path, nil
}

func parseDirectoryNameFromDecodedInfoData(infoData *model.OrderedMap) (string, error) {
	if isMultiFileMode(infoData) {
		return parseNameFromDecodedInfoData(infoData)
	}
	return "", nil
}

func isMultiFileMode(infoData *model.OrderedMap) bool {
	_, err := readIntValueFromMap(infoData, "length")
	isMultiFile := err != nil
	return isMultiFile
}

func parseNameFromDecodedInfoData(infoData *model.OrderedMap) (string, error) {
	return readStringValueFromMap(infoData, "name")
}

// Helpers
func readStringValueFromMap(m *model.OrderedMap, key string) (string, error) {
	value, exists := m.GetExists(key)
	if !exists {
		return "", fmt.Errorf("String value not found in map for key %v", key)
	}

	castedValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("Unable to cast value %v to string", value)
	}

	return castedValue, nil
}

func readIntValueFromMap(m *model.OrderedMap, key string) (int, error) {
	value, exists := m.GetExists(key)
	if !exists {
		return 0, fmt.Errorf("Int value not found in map for key %v", key)
	}

	castedValue, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("Unable to cast value %v to int", value)
	}

	return castedValue, nil
}

func readBytesValueFromMap(m *model.OrderedMap, key string) ([]byte, error) {
	str, err := readStringValueFromMap(m, key)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func readListFromMap(m *model.OrderedMap, key string) ([]interface{}, error) {
	value, exists := m.GetExists(key)
	if !exists {
		return nil, fmt.Errorf("List of dictionaries value not found in map for key %v", key)
	}

	castedValue, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to cast value %v to []*model.OrderedMap (type is %v)", value, reflect.TypeOf(value))
	}

	return castedValue, nil
}

func readDictionaryValueFromMap(m *model.OrderedMap, key string) (*model.OrderedMap, error) {
	value, exists := m.GetExists(key)
	if !exists {
		return nil, fmt.Errorf("Dictionary value not found in map for key %v", key)
	}

	castedValue, ok := value.(*model.OrderedMap)
	if !ok {
		return nil, fmt.Errorf("Unable to cast value %v to *model.OrderedMap", value)
	}

	return castedValue, nil
}
