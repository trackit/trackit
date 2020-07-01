package utils

import "strconv"

// GetRegionForURL removes the last character of an AWS region if it's not a number
func GetRegionForURL(region string) string {
	if len(region) <= 0 {
		return region
	}

	_, err := strconv.ParseInt(region[len(region)-1:], 10, 32)
	if err != nil {
		region = region[:len(region)-1]
	}

	return region
}
