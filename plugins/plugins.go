//   Copyright 2019 MSolution.IO
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

// Package plugins implements loading all the plugins at init-time
package plugins

import (
	// This is solely for the side effects of importation, i.e. registration of the plugins
	_ "github.com/trackit/trackit/plugins/account/networkEc2"
	_ "github.com/trackit/trackit/plugins/account/s3Traffic"
	_ "github.com/trackit/trackit/plugins/account/unattachedEIP"
	_ "github.com/trackit/trackit/plugins/account/unusedEBS"
)
