//   Copyright 2020 MSolution.IO
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

package indexes

import (
	"encoding/json"
)

type partialTemplate struct {
	Version int `json:"version"`
}

func extractVersionFromTemplate() error {
	for index, data := range versioningData {
		var obj partialTemplate
		err := json.Unmarshal([]byte(data.Template), &obj)
		if err != nil {
			return nil
		}
		versioningData[index].Version = obj.Version
	}
	return nil
}
