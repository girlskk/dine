package util

import (
	"encoding/json"
	"fmt"
	"log"
)

// PrettyJson 打印json格式的结构体数据
func PrettyJson(req any) {
	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
	}
	prettyOutput := string(data)
	fmt.Println(prettyOutput)
}
