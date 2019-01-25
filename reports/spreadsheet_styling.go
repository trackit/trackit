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
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var stylesList = map[string]int{}

/* TODO: Store into JSON file ? */
var stylesRaw = map[string]string{
	"borders": `{"border": [{"type": "left", "color": "#000000", "style": 1},
							{"type": "right", "color": "#000000", "style": 1},
							{"type": "top", "color": "#000000", "style": 1},
							{"type": "bottom", "color": "#000000", "style": 1}]}`,
	"bold": `{"font": {"bold": true}}`,
	"centerText": `{"alignment": {"horizontal": "center", "vertical": "center"}}`,
	"price": `{"number_format": 176, "decimal_places": 3}`,
	"percentage": `{"number_format": 10}`,
	"green": `{"font": {"color": "#006600"}, "fill": {"type": "pattern", "pattern": 1, "color": ["#CCFFCC"]}}`,
	"orange": `{"font": {"color": "#C65911"}, "fill": {"type": "pattern", "pattern": 1, "color": ["#F8CBAD"]}}`,
	"red": `{"font": {"color": "#CC0000"}, "fill": {"type": "pattern", "pattern": 1, "color": ["#FFCCCC"]}}`,
}

func getStyleId(file *excelize.File, styles []string) (int, error){
	name := strings.Join(styles, "")
	if id, ok := stylesList[name]; ok {
		return id, nil
	} else {
		style, err := generateStyle(styles)
		if err != nil {
			return -1, err
		}
		return registerStyle(file, name, style)
	}
}

func generateStyle(styles []string) (string, error) {
	output := "{}"
	for _, styleName := range styles {
		if style, ok := stylesRaw[styleName]; ok {
			merged, err := mergeStringJson(output, style)
			if err != nil {
				return "", err
			}
			output = merged
		}
	}
	return output, nil
}

func registerStyle(file *excelize.File, name string, style string) (int, error) {
	styleId, err := file.NewStyle(style)
	if err != nil {
		return -1, err
	}
	stylesList[name] = styleId
	return styleId, nil
}
