package config

import (
	"fmt"
	"strings"
)

func extractRunAddress(addrIn string) (string, int, error) {
	parts := strings.SplitN(addrIn, ":", 2)
	host := ""
	port := 0
	var err error
	if len(parts) == 2 {
		host = parts[0]
		_, err := fmt.Sscan(parts[1], &port)
		if err != nil {
			err = fmt.Errorf("PORT value is invalid: %v. Error: %s", parts[1], err.Error())
		}
	} else {
		host = addrIn
	}

	return host, port, err
}
