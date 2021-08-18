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

CREATE TABLE tagbot_discount_code (
       id                         INTEGER       NOT NULL AUTO_INCREMENT,
       code                       VARCHAR(10)   NOT NULL,
       description                VARCHAR(200)  NOT NULL DEFAULT "",
       created                    TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
       modified                   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
       CONSTRAINT PRIMARY KEY (id),
       CONSTRAINT unique_code UNIQUE KEY (code)
);

ALTER TABLE tagbot_user ADD free_tier_end_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE tagbot_user ADD discount_code_id INTEGER NULL DEFAULT NULL;
ALTER TABLE tagbot_user ADD CONSTRAINT foreign_discount_code FOREIGN KEY (discount_code_id) REFERENCES tagbot_discount_code(id) ON DELETE SET NULL;

-- Set all the existing rows to have a usable value for free_tier_end_at (i.e. the default of 14 days)
UPDATE tagbot_user INNER JOIN user ON user.id = tagbot_user.user_id SET tagbot_user.free_tier_end_at = DATE_ADD(user.created, INTERVAL 14 DAY);
