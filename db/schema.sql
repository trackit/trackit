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

CREATE TABLE aws_update_job (
	id                     INTEGER      NOT NULL AUTO_INCREMENT,
	created                TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified               TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	aws_bill_repository_id INTEGER      NOT NULL,
	expired                TIMESTAMP    NOT NULL DEFAULT ADDTIME(CURRENT_TIMESTAMP, '02:00:00'),
	status                 VARCHAR(255) NOT NULL,
	worker_id              VARCHAR(32)  NOT NULL,
	error                  VARCHAR(255) NOT NULL,
	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT foreign_bill_repository (aws_bill_repository_id) REFERENCES aws_bill_repository(id) ON DELETE CASCADE
);
