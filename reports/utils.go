package reports

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trackit/trackit-server/aws"
)

func mergeStringJson(style1 string, style2 string) (string, error) {
	merged := make(map[string]interface{})
	err := json.NewDecoder(strings.NewReader(style1)).Decode(&merged)
	if err != nil {
		return "", nil
	}
	err = json.NewDecoder(strings.NewReader(style2)).Decode(&merged)
	if err != nil {
		return "", nil
	}
	output, err := json.Marshal(merged)
	if err != nil {
		return "", nil
	}
	return string(output), nil
}

func formatAwsAccount(aa aws.AwsAccount) string {
	return fmt.Sprintf("%s (%s)", aa.Pretty, aa.AwsIdentity)
}

func getAwsIdentities(aas []aws.AwsAccount) []string {
	identities := make([]string, len(aas))
	for index, account := range aas {
		identities[index] = account.AwsIdentity
	}
	return identities
}

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
