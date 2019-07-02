//   Copyright 2018 MSolution.IO
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

package reports

import (
	"errors"
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type conditionalFormat struct {
	value  string
	custom bool
	styles []string
}

type conditionalFormats []conditionalFormat

var conditionalStylesList = map[string]int{}

/* TODO: Store into JSON file ? */
var conditionsRaw = map[string]string{
	"above90percent":  `{"type": "cell", "criteria": ">", "value": "0.9"}`,
	"above85percent":  `{"type": "cell", "criteria": ">", "value": "0.85"}`,
	"above80percent":  `{"type": "cell", "criteria": ">", "value": "0.8"}`,
	"above75percent":  `{"type": "cell", "criteria": ">", "value": "0.75"}`,
	"above60percent":  `{"type": "cell", "criteria": ">", "value": "0.6"}`,
	"above30percent":  `{"type": "cell", "criteria": ">", "value": "0.3"}`,
	"validPercentage": `{"type": "cell", "criteria": ">", "value": "-1"}`,
	"positive":        `{"type": "cell", "criteria": ">", "value": "0"}`,
	"negative":        `{"type": "cell", "criteria": "<", "value": "0"}`,
	"zero":            `{"type": "cell", "criteria": "=", "value": "0"}`,
	"empty":           `{"type": "cell", "criteria": "=", "value": ""}`,
}

/* TODO: Update error handing (Errors should not interrupt spreadsheet generation since it is only styling issue) */
func (cs conditionalFormats) getConditionalFormatting(file *excelize.File) (string, error) {
	conditions := make([]string, 0)
	for _, condition := range cs {
		formattedCondition, err := condition.getConditionalFormatting(file)
		if err != nil {
			return "", err
		}
		conditions = append(conditions, formattedCondition)
	}
	return fmt.Sprintf("[%s]", strings.Join(conditions, ",")), nil
}

/* TODO: Update error handing (Errors should not interrupt spreadsheet generation since it is only styling issue) */
func (c conditionalFormat) getConditionalFormatting(file *excelize.File) (string, error) {
	condition := c.value
	if !c.custom {
		value, ok := conditionsRaw[c.value]
		if !ok {
			return "", errors.New(fmt.Sprintf("Condition %s not found", c.value))
		}
		condition = value
	}
	styleId, err := getConditionalStyleId(file, c.styles)
	if err != nil {
		return "", err
	}
	style := fmt.Sprintf(`{"format": %d}`, styleId)
	formattedCondition, err := mergeStringJson(condition, style)
	if err != nil {
		return "", err
	}
	return formattedCondition, nil
}

func getConditionalStyleId(file *excelize.File, styles []string) (int, error) {
	name := strings.Join(styles, "")
	if id, ok := conditionalStylesList[name]; ok {
		return id, nil
	} else {
		style, err := generateStyle(styles)
		if err != nil {
			return -1, err
		}
		return registerConditionalStyle(file, name, style)
	}
}

func registerConditionalStyle(file *excelize.File, name string, style string) (int, error) {
	styleId, err := file.NewConditionalStyle(style)
	if err != nil {
		return -1, err
	}
	conditionalStylesList[name] = styleId
	return styleId, nil
}
