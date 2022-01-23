CREATE TABLE IF NOT EXISTS `logs` (`uid`, `application`, `error`, `time`);
CREATE TABLE IF NOT EXISTS `requests` (`uid`, `application`, `client_ip`, `method`, `path`, `url`, `host`, `headers`, `body`, `referrer`, `user_agent`, `time`);
CREATE TABLE IF NOT EXISTS `responses` (`uid`, `request_uid`, `application`, `client_ip`, `status`, `body`, `path`, `headers`, `size`, `time`);