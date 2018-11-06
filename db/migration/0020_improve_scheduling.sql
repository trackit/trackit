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
