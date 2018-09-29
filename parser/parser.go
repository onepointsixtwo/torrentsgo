package parser

import (
	"fmt"
	"github.com/onepointsixtwo/torrentsgo/bencoding"
	"github.com/onepointsixtwo/torrentsgo/model"
	"io"
	"net/url"
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

func parseMetaInfoFromDecodedData(data map[string]interface{}) (*model.MetaInfo, error) {
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

func parseAnnounceUrlsFromDecodedData(data map[string]interface{}) ([]*url.URL, error) {
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

func parseCreationDateFromDecodedData(data map[string]interface{}) time.Time {
	timestamp, _ := readIntValueFromMap(data, "creation date")
	return time.Unix(int64(timestamp), 0)
}

func parseCommentFromDecodedData(data map[string]interface{}) string {
	str, _ := readStringValueFromMap(data, "comment")
	return str
}

func parseCreatedByFromDecodedData(data map[string]interface{}) string {
	str, _ := readStringValueFromMap(data, "created by")
	return str
}

func parseEncodingFromDecodedData(data map[string]interface{}) string {
	str, _ := readStringValueFromMap(data, "encoding")
	return str
}

// Info parsing

func parseInfoFromDecodedData(data map[string]interface{}) (*model.Info, error) {
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

func parsePieceLengthFromDecodedInfoData(infoData map[string]interface{}) (int, error) {
	return readIntValueFromMap(infoData, "piece length")
}

func parsePiecesDataFromDecodedInfoData(infoData map[string]interface{}) ([]byte, error) {
	return readBytesValueFromMap(infoData, "pieces")
}

func parsePrivateFromDecodedInfoData(infoData map[string]interface{}) int {
	value, _ := readIntValueFromMap(infoData, "private")
	return value
}

func parseFilesFromDecodedInfoData(infoData map[string]interface{}) ([]*model.File, error) {
	if isMultiFileMode(infoData) {
		return parseMultiFileModeFilesFromDecodedInfoData(infoData)
	}
	return parseSingleFileModeFilesFromDecodedInfoData(infoData)
}

func parseSingleFileModeFilesFromDecodedInfoData(infoData map[string]interface{}) ([]*model.File, error) {
	file, err := parseFileFromMap(infoData, "name")
	if err != nil {
		return nil, err
	}

	files := make([]*model.File, 1)
	files[0] = file
	return files, nil
}

func parseMultiFileModeFilesFromDecodedInfoData(infoData map[string]interface{}) ([]*model.File, error) {
	filesList, err := readListOfDictionariesValueFromMap(infoData, "files")
	if err != nil {
		return nil, err
	}

	filesLength := len(filesList)
	files := make([]*model.File, 0)
	for i := 0; i < filesLength; i++ {
		dict := filesList[i]
		file, err := parseFileFromMap(dict, "path")
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

func parseFileFromMap(mp map[string]interface{}, pathName string) (*model.File, error) {
	fileName, err := readStringValueFromMap(mp, pathName)
	if err != nil {
		return nil, err
	}
	length, err2 := readIntValueFromMap(mp, "length")
	if err != nil {
		return nil, err2
	}
	md5Sum, _ := readStringValueFromMap(mp, "md5sum")

	return model.NewFile(fileName, length, md5Sum), nil
}

func parseDirectoryNameFromDecodedInfoData(infoData map[string]interface{}) (string, error) {
	if isMultiFileMode(infoData) {
		return parseNameFromDecodedInfoData(infoData)
	}
	return "", nil
}

func isMultiFileMode(infoData map[string]interface{}) bool {
	_, err := readIntValueFromMap(infoData, "length")
	isMultiFile := err != nil
	return isMultiFile
}

func parseNameFromDecodedInfoData(infoData map[string]interface{}) (string, error) {
	return readStringValueFromMap(infoData, "name")
}

// Helpers
func readStringValueFromMap(m map[string]interface{}, key string) (string, error) {
	value, exists := m[key]
	if !exists {
		return "", fmt.Errorf("String value not found in map for key %v", key)
	}

	castedValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("Unable to cast value %v to string", value)
	}

	return castedValue, nil
}

func readIntValueFromMap(m map[string]interface{}, key string) (int, error) {
	value, exists := m[key]
	if !exists {
		return 0, fmt.Errorf("Int value not found in map for key %v", key)
	}

	castedValue, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("Unable to cast value %v to int", value)
	}

	return castedValue, nil
}

func readBytesValueFromMap(m map[string]interface{}, key string) ([]byte, error) {
	str, err := readStringValueFromMap(m, key)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func readListOfDictionariesValueFromMap(m map[string]interface{}, key string) ([]map[string]interface{}, error) {
	value, exists := m[key]
	if !exists {
		return nil, fmt.Errorf("List of dictionaries value not found in map for key %v", key)
	}

	castedValue, ok := value.([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to cast value %v to []map[string]interface{}", value)
	}

	return castedValue, nil
}

func readDictionaryValueFromMap(m map[string]interface{}, key string) (map[string]interface{}, error) {
	value, exists := m[key]
	if !exists {
		return nil, fmt.Errorf("Dictionary value not found in map for key %v", key)
	}

	castedValue, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to cast value %v to map[string]interface{}", value)
	}

	return castedValue, nil
}
