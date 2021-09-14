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
