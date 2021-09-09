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

package es

import (
	"encoding/json"
	"testing"
)

func TestToJsonableEmptySimplifiedCostsDocument(t *testing.T) {
	jsonRes := SimplifiedCostsDocument{}.ToJsonable()
	expectedResult := "{}"
	marshalled, _ := json.Marshal(jsonRes)
	if string(marshalled) != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, string(marshalled))
	}
}

func TestToJsonableSimpleSingleLevelResults(t *testing.T) {
	jsonRes := SimplifiedCostsDocument{
		Key:      "",
		HasValue: false,
		Value:    0,
		Children: []SimplifiedCostsDocument{
			{
				Key:          "AmazonRDS",
				HasValue:     true,
				Value:        12,
				Children:     []SimplifiedCostsDocument{},
				ChildrenKind: "",
			},
			{
				Key:          "AmazonSimpleDB",
				HasValue:     true,
				Value:        13,
				Children:     []SimplifiedCostsDocument{},
				ChildrenKind: "",
			},
		},
		ChildrenKind: "product",
	}.ToJsonable()
	expectedResult := `{
	"product": {
		"AmazonRDS": 12,
		"AmazonSimpleDB": 13
	}
}`
	marshalled, _ := json.MarshalIndent(jsonRes, "", "\t")
	if string(marshalled) != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, string(marshalled))
	}
}

func TestToJsonableSimpleTwoLevelsResults(t *testing.T) {
	jsonRes := SimplifiedCostsDocument{
		Key:      "",
		HasValue: false,
		Value:    0,
		Children: []SimplifiedCostsDocument{
			{
				Key:      "AmazonS3",
				HasValue: false,
				Value:    0,
				Children: []SimplifiedCostsDocument{
					{
						Key:          "123456789",
						HasValue:     true,
						Value:        42,
						Children:     []SimplifiedCostsDocument{},
						ChildrenKind: "",
					},
				},
				ChildrenKind: "account",
			},
			{
				Key:      "AmazonEC2",
				HasValue: false,
				Value:    0,
				Children: []SimplifiedCostsDocument{
					{
						Key:          "123456789",
						HasValue:     true,
						Value:        24,
						Children:     []SimplifiedCostsDocument{},
						ChildrenKind: "",
					},
					{
						Key:          "987654321",
						HasValue:     true,
						Value:        25,
						Children:     []SimplifiedCostsDocument{},
						ChildrenKind: "",
					},
				},
				ChildrenKind: "account",
			},
		},
		ChildrenKind: "product",
	}.ToJsonable()
	expectedResult := `{
	"product": {
		"AmazonEC2": {
			"account": {
				"123456789": 24,
				"987654321": 25
			}
		},
		"AmazonS3": {
			"account": {
				"123456789": 42
			}
		}
	}
}`
	marshalled, _ := json.MarshalIndent(jsonRes, "", "\t")
	if string(marshalled) != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, string(marshalled))
	}
}
