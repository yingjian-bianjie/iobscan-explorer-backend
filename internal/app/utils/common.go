package utils

import "fmt"

func FormatPhoneNumber(phonePrefix, phoneNumber string) string {
	if phonePrefix[0] == '+' {
		phonePrefix = phonePrefix[1:]
	}
	return fmt.Sprintf("%s%s", phonePrefix, phoneNumber)
}
