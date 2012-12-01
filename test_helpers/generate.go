package testHelpers

import (
	"fmt"
	"time"
)

func GenerateEmailAddress() string {
	return fmt.Sprintf("email-%d@mydomain.com", time.Now().UnixNano())
}

func GenerateScreenName() string {
	return fmt.Sprintf("screenName%d", time.Now().UnixNano())
}
