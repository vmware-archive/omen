package userio

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetConfirmation(confMsg string) bool {
	fmt.Println(confMsg)
	reader := bufio.NewReader(os.Stdin)

	for {
		applyBytes, _, _ := reader.ReadLine()
		applyChanges := string(applyBytes)

		if strings.EqualFold(applyChanges, "Y") {
			return true
		}

		if strings.EqualFold(applyChanges, "N") {
			return false
		}
	}
}
