package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrettyDisplay(message string, v interface{}) {
	fmt.Printf("====%s====\n", message)
	empJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(empJSON))
	fmt.Println("==================================\n")
}
