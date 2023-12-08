package custom_validators

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
	"time"
)

func IsValidRUT(rut string) bool {
	rutRegex := regexp.MustCompile(`^(\d{1,8})-([\dKk]{1})$`)
	if !rutRegex.MatchString(rut) {
		return false
	}

	matches := rutRegex.FindStringSubmatch(rut)
	num, dv := matches[1], strings.ToUpper(matches[2])

	var total int
	factors := []int{2, 3, 4, 5, 6, 7}
	for i, j := len(num)-1, 0; i >= 0; i, j = i-1, j+1 {
		digit := int(num[i] - '0')
		total += digit * factors[j%6]
	}

	remainder := total % 11
	calculatedDV := 11 - remainder

	var strCalculatedDV string
	if calculatedDV == 11 {
		strCalculatedDV = "0"
	} else if calculatedDV == 10 {
		strCalculatedDV = "K"
	} else {
		strCalculatedDV = string(rune(calculatedDV + '0'))
	}

	return dv == strCalculatedDV
}

func IsValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regex.MatchString(email)
}

func IsValidPhone(phoneNumber string) bool {
	regex := regexp.MustCompile(`^\(\+569\)\d{8}$`)
	return regex.MatchString(phoneNumber)
}

func IsValidUrl(url string) bool {
	regex := regexp.MustCompile(`^https?://[a-z0-9\.-]+\.[a-z]{2,4}/?.*$`)
	return regex.MatchString(url)
}

func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

func IsMinLengthChar(value string, minLength int) bool {
	return len(value) >= minLength
}

func IsValidDateTime(stringDate string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", stringDate)
	return err == nil
}

func IsInValues(value string, allowedValues []string) bool {

	for _, allowedValue := range allowedValues {
		if value == allowedValue {
			return true
		}
	}
	return false
}
func GreaterThanOrEqualZero(value int64) bool {
	return value >= 0

}
