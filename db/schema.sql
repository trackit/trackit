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

CREATE TABLE user (
	id INTEGER NOT NULL AUTO_INCREMENT,
	created       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	email         VARCHAR(254) NOT NULL,
	auth          VARCHAR(255) NOT NULL,
	next_external VARCHAR(96)      NULL,
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT unique_email UNIQUE KEY (email)
);

CREATE TABLE aws_account (
	id       INTEGER      NOT NULL AUTO_INCREMENT,
	created  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	user_id  INTEGER      NOT NULL,
	pretty   VARCHAR(255) NOT NULL,
	role_arn VARCHAR(255) NOT NULL,
	external VARCHAR(255) NOT NULL,
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_user   FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE aws_bill_repository (
	id                     INTEGER       NOT NULL AUTO_INCREMENT,
	created                TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified               TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	aws_account_id         INTEGER       NOT NULL,
	bucket                 VARCHAR(63)   NOT NULL,
	prefix                 VARCHAR(1024) NOT NULL,
	last_imported_manifest DATETIME      NOT NULL DEFAULT "1970-01-01 00:00:00",
	next_update            DATETIME      NOT NULL DEFAULT "1970-01-01 00:00:00",
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_aws_account    FOREIGN KEY (aws_account_id) REFERENCES aws_account(id) ON DELETE CASCADE
--	CONSTRAINT unique_per_account     UNIQUE  KEY (aws_account_id, bucket, prefix)
);

CREATE VIEW aws_bill_repository_due_update AS
	SELECT * FROM aws_bill_repository WHERE next_update <= NOW()
;

CREATE TABLE aws_product_pricing_update (
	id INTEGER NOT NULL AUTO_INCREMENT,
	product VARCHAR(255) NOT NULL,
	etag VARCHAR(255) NOT NULL,
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT UNIQUE KEY (product)
);

CREATE TABLE aws_product_pricing_ec2 (
	sku VARCHAR(255) NOT NULL,
	etag VARCHAR(255) NOT NULL,
	region VARCHAR(255) NOT NULL,
	instance_type VARCHAR(255) NOT NULL,
	current_generation BOOLEAN NOT NULL,
	vcpu INTEGER NOT NULL,
	memory VARCHAR(255) NOT NULL,
	storage VARCHAR(255) NOT NULL,
	network_performance VARCHAR(255) NOT NULL,
	tenancy VARCHAR(255) NOT NULL,
	operating_system VARCHAR(255) NOT NULL,
	ecu VARCHAR(255) NOT NULL,
	CONSTRAINT PRIMARY KEY (etag, sku)
);
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

ALTER TABLE user ADD parent_user_id INTEGER NULL;
ALTER TABLE user ADD CONSTRAINT parent_user FOREIGN KEY (parent_user_id) REFERENCES user(id) ON DELETE CASCADE;
--   Copyright 2018 MSolution.IO
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

ALTER TABLE aws_bill_repository ADD status VARCHAR(255) NULL;
--   Copyright 2018 MSolution.IO
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

ALTER TABLE aws_bill_repository ADD error VARCHAR(255) NOT NULL DEFAULT "";
ALTER TABLE aws_bill_repository DROP COLUMN status;

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

CREATE TABLE forgotten_password (
	id       INTEGER      NOT NULL AUTO_INCREMENT,
	created  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id  INTEGER      NOT NULL,
	token    VARCHAR(255) NOT NULL,
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_user   FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

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

--   Copyright 2018 MSolution.IO
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

ALTER TABLE aws_account_update_job ADD ec2Error VARCHAR(255) NOT NULL DEFAULT "";

--   Copyright 2018 MSolution.IO
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

ALTER TABLE aws_bill_repository ADD grace_update DATETIME NOT NULL DEFAULT "1970-01-01 00:00:00";

--   Copyright 2018 MSolution.IO
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

ALTER TABLE user ADD aws_customer_identifier varchar(255) NOT NULL DEFAULT "";

--   Copyright 2018 MSolution.IO
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

CREATE TABLE emailed_anomaly (
	id             INTEGER      NOT NULL AUTO_INCREMENT,
	aws_account_id INTEGER      NOT NULL,
	product        VARCHAR(255) NOT NULL,
	recipient      VARCHAR(255) NOT NULL,
	date           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_aws_account FOREIGN KEY (aws_account_id) REFERENCES aws_account(id) ON DELETE CASCADE
);

--   Copyright 2018 MSolution.IO
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

ALTER TABLE aws_account ADD payer BOOL NOT NULL DEFAULT "1";

--   Copyright 2018 MSolution.IO
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

CREATE TABLE shared_account (
  id                     INTEGER   NOT NULL AUTO_INCREMENT,
  account_id             INTEGER   NOT NULL,
  user_id                INTEGER   NOT NULL,
  user_permission        INTEGER   NOT NULL DEFAULT 0,
  sharing_accepted       BOOL      NOT NULL DEFAULT 0,
  CONSTRAINT PRIMARY KEY (id),
  CONSTRAINT foreign_aws_account FOREIGN KEY (account_id) REFERENCES aws_account(id) ON DELETE CASCADE,
  CONSTRAINT foreign_user_id FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

--   Copyright 2018 MSolution.IO
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

ALTER TABLE user ADD aws_customer_entitlement bool NOT NULL DEFAULT 1;

--   Copyright 2018 MSolution.IO
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

ALTER TABLE aws_account_update_job ADD historyError VARCHAR(255) NOT NULL DEFAULT "";

--   Copyright 2018 MSolution.IO
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

CREATE TABLE aws_account_plugins_job (
	id                     INTEGER      NOT NULL AUTO_INCREMENT,
	created                TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	aws_account_id         INTEGER      NOT NULL,
	completed              TIMESTAMP    NOT NULL DEFAULT 0,
	worker_id              VARCHAR(255) NOT NULL,
	jobError               VARCHAR(255) NOT NULL DEFAULT "",
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_aws_account FOREIGN KEY (aws_account_id) REFERENCES aws_account(id) ON DELETE CASCADE
);

ALTER TABLE aws_account ADD next_update_plugins DATETIME NOT NULL DEFAULT "1970-01-01 00:00:00";
ALTER TABLE aws_account ADD grace_update_plugins DATETIME NOT NULL DEFAULT "1970-01-01 00:00:00";

--   Copyright 2018 MSolution.IO
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

CREATE VIEW aws_account_plugins_due_update AS
	SELECT * FROM aws_account WHERE next_update_plugins <= NOW() AND grace_update_plugins <= NOW()
;

--   Copyright 2018 MSolution.IO
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

CREATE OR REPLACE VIEW aws_bill_repository_due_update AS
	SELECT * FROM aws_bill_repository WHERE next_update <= NOW() AND grace_update <= NOW()
;

--   Copyright 2018 MSolution.IO
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

ALTER TABLE aws_bill_update_job MODIFY COLUMN error VARCHAR(2000) NOT NULL;

--   Copyright 2018 MSolution.IO
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

ALTER TABLE user ADD next_update_entitlement DATETIME NOT NULL DEFAULT "1970-01-01 00:00:00";

CREATE VIEW user_entitlement_due_update AS
	SELECT * FROM user WHERE next_update_entitlement <= NOW()
;

--   Copyright 2018 MSolution.IO
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

CREATE OR REPLACE VIEW aws_bill_repository_due_update AS
	SELECT * FROM aws_bill_repository WHERE next_update <= NOW()
;

CREATE OR REPLACE VIEW aws_account_due_update AS
	SELECT * FROM aws_account WHERE next_update <= NOW()
;

CREATE OR REPLACE VIEW aws_account_plugins_due_update AS
	SELECT * FROM aws_account WHERE next_update_plugins <= NOW()
;

ALTER TABLE aws_bill_repository DROP COLUMN grace_update;
ALTER TABLE aws_account DROP COLUMN grace_update, DROP COLUMN grace_update_plugins;

--   Copyright 2018 MSolution.IO
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

CREATE TABLE aws_account_reports_job (
	id                     INTEGER      NOT NULL AUTO_INCREMENT,
	created                TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	aws_account_id         INTEGER      NOT NULL,
	completed              TIMESTAMP    NOT NULL DEFAULT 0,
	worker_id              VARCHAR(255) NOT NULL,
	jobError               VARCHAR(255) NOT NULL DEFAULT "",
	spreadsheetError       VARCHAR(255) NOT NULL DEFAULT "",
	costDiffError          VARCHAR(255) NOT NULL DEFAULT "",
	ec2UsageReportError    VARCHAR(255) NOT NULL DEFAULT "",
	rdsUsageReportError    VARCHAR(255) NOT NULL DEFAULT "",
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_aws_account FOREIGN KEY (aws_account_id) REFERENCES aws_account(id) ON DELETE CASCADE
);

--   Copyright 2018 MSolution.IO
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

ALTER TABLE aws_account ADD last_spreadsheet_report_generation DATETIME NOT NULL DEFAULT "1970-01-01 00:00:00";
ALTER TABLE aws_account_update_job ADD monthly_reports_generated bool NOT NULL DEFAULT 0;
