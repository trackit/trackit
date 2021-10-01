package reports

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/trackit/trackit/aws"
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

func getAwsAccount(account string, aas []aws.AwsAccount) *aws.AwsAccount {
	for _, aa := range aas {
		if aa.AwsIdentity == account {
			return &aa
		}
	}
	return nil
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

func downloadFile(url string) (data []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := res.Body.Close(); err == nil {
			err = closeErr
		}
	}()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, res.Body)
	if err != nil {
		return
	}
	data = buffer.Bytes()
	return
}
