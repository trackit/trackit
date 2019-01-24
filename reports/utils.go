package reports

import (
	"fmt"

	"github.com/trackit/trackit-server/aws"
)

func formatMetric(value float64) interface{} {
	if value == -1 {
		return "N/A"
	}
	return value
}

func formatMetricPercentage(value float64) interface{} {
	if value == -1 {
		return "N/A"
	}
	return value / 100
}

func getTotal(values map[string]float64) float64 {
	var total float64
	for _, value := range values {
		total += value
	}
	return total
}

func formatTags(tags map[string]string) []string {
	formattedTags := make([]string, 0, len(tags))
	for key, value := range tags {
		formattedTags = append(formattedTags, fmt.Sprintf("%s:%s", key, value))
	}
	return formattedTags
}

func getIdentities(aas []aws.AwsAccount) []string {
	identities := make([]string, len(aas))
	for i, account := range aas {
		identities[i] = account.AwsIdentity
	}
	return identities
}
