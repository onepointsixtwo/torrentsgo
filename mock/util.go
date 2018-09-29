package mock

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(anything interface{}) string {
	b, err := json.MarshalIndent(anything, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return fmt.Sprintf(string(b))
}
