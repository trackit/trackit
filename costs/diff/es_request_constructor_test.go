//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package diff

import (
	"encoding/json"
	"testing"
	"time"
)

func TestQueryAccountFiltersMultipleAccounts(t *testing.T) {
	linkedAccountID := []int{
		123456,
		98765432,
	}
	expectedResult := `{"terms":{"usageAccountId":[123456,98765432]}}`
	res := createQueryAccountFilter(linkedAccountID)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestQueryAccountFiltersSingleAccount(t *testing.T) {
	linkedAccountID := []int{
		123456,
	}
	expectedResult := `{"terms":{"usageAccountId":[123456]}}`
	res := createQueryAccountFilter(linkedAccountID)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}

func TestQueryTimeRange(t *testing.T) {
	durationBegin, _ := time.Parse("2006-01-02", "2017-01-12")
	durationEnd, _ := time.Parse("2006-01-02", "2017-05-23")
	expectedResult := `{"range":{"usageStartDate":{"from":"2017-01-12T00:00:00Z","include_lower":true,"include_upper":true,"to":"2017-05-23T00:00:00Z"}}}`

	res := createQueryTimeRange(durationBegin, durationEnd)
	src, err := res.Source()
	if err != nil {
		t.Fatal(err)
	}
	jsonRes, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonRes) != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, string(jsonRes))
	}
}
