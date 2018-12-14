// SPDX-License-Identifier: Apache-2.0

package share

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseBool(str string) (bool, error) {
	b, err := strconv.ParseBool(str)
	if err == nil {
		return b, err
	}

	if strings.EqualFold(str, "yes") || strings.EqualFold(str, "y") || strings.EqualFold(str, "on") {
		return true, nil
	} else if strings.EqualFold(str, "no") || strings.EqualFold(str, "n") || strings.EqualFold(str, "off") {
		return false, nil
	}

	return false, fmt.Errorf("ParseBool")
}
