--   Copyright 2017 MSolution.IO
--
--   Licensed under the Apache License, Version 2.0 (the "License");
--   you may not use this file except in compliance with the License.
--   You may obtain a copy of the License at
--
--       http://www.apache.org/licenses/LICENSE-2.0
--
--   Unless required by applicable law or agreed to in writing, software
--   distributed under the License is distributed on an "AS IS" BASIS,
--   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
--   See the License for the specific language governing permissions and
--   limitations under the License.

CREATE TABLE aws_bill_update_job (
	id                     INTEGER      NOT NULL AUTO_INCREMENT,
	created                TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified               TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	aws_bill_repository_id INTEGER      NOT NULL,
	expired                TIMESTAMP    NOT NULL DEFAULT ADDTIME(CURRENT_TIMESTAMP, '02:00:00'),
	completed              TIMESTAMP    NOT NULL DEFAULT 0,
	worker_id              VARCHAR(255) NOT NULL,
	error                  VARCHAR(255) NOT NULL,
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_bill_repository FOREIGN KEY (aws_bill_repository_id) REFERENCES aws_bill_repository(id) ON DELETE CASCADE
);
