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

package mail

import (
	"testing"
	"bytes"
)

func TestSendMail(t *testing.T) {
	m := Mail{
		"",
		"",
		"",
		"",
		"team@msolution.io",
		"thibaut@trackit.io",
		"test subject!",
		"test body!",
	}
	msg := m.buildMessage()
	template := []byte(
		"From: " + m.Sender + "\r\n" +
			"To: " + m.Recipient + "\r\n" +
			"Subject: " + m.Subject + "\r\n" +
			"\r\n" +
			"test body!",
		)
	if !bytes.Equal(msg, template) {
		t.Fatalf("Unexcepted message: (%s) instead of (%s)", msg, template)
	}
}
