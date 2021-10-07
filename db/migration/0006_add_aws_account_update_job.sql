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

CREATE TABLE aws_account_update_job (
	id                     INTEGER      NOT NULL AUTO_INCREMENT,
	created                TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	aws_account_id         INTEGER      NOT NULL,
	completed              TIMESTAMP    NOT NULL DEFAULT 0,
	worker_id              VARCHAR(255) NOT NULL,
	jobError               VARCHAR(255) NOT NULL DEFAULT "",
	rdsError               VARCHAR(255) NOT NULL DEFAULT "",
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_aws_account FOREIGN KEY (aws_account_id) REFERENCES aws_account(id) ON DELETE CASCADE
);

ALTER TABLE aws_account ADD next_update DATETIME NOT NULL DEFAULT "1970-01-01 00:00:00";
ALTER TABLE aws_account ADD grace_update DATETIME NOT NULL DEFAULT "1970-01-01 00:00:00";

CREATE VIEW aws_account_due_update AS
	SELECT * FROM aws_account WHERE next_update <= NOW() AND grace_update <= NOW()
;
