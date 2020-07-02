--   Copyright 2020 MSolution.IO
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

CREATE TABLE aws_account_update_most_used_tags_job (
	id                     INTEGER      NOT NULL AUTO_INCREMENT,
	created                TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	aws_account_id         INTEGER      NOT NULL,
	completed              TIMESTAMP    NOT NULL DEFAULT 0,
	worker_id              VARCHAR(255) NOT NULL,
	job_error              VARCHAR(255) NOT NULL DEFAULT "",
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_aws_account FOREIGN KEY (aws_account_id) REFERENCES aws_account(id) ON DELETE CASCADE
);

CREATE TABLE most_used_tags (
	id                     INTEGER      NOT NULL AUTO_INCREMENT,
	report_date            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	aws_account_id         INTEGER      NOT NULL,
	tags                   VARCHAR(255) NOT NULL,
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_aws_account FOREIGN KEY (aws_account_id) REFERENCES aws_account(id) ON DELETE CASCADE
);