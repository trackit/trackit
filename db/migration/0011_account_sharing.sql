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

CREATE TABLE account_sharing (
	id                     INTEGER      NOT NULL AUTO_INCREMENT,
	account_id             INTEGER      NOT NULL COMMENT 'Invite account id',
	owner_id               INTEGER      NOT NULL COMMENT 'user id of the Owner of the account',
	user_id                INTEGER      NOT NULL COMMENT 'id of the invited user',
	user_permission        TINYINT(2)   NOT NULL DEFAULT 0 COMMENT '0 : admin, 1 : standard, 2 : read-only',
	account_status         TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '0 : user has never loged in (pending status), 1 : user has logged once',
	CONSTRAINT PRIMARY KEY (id)
);
