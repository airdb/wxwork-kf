CREATE TABLE `tab_talk`
(
    `id`         int(10)   NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL     DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL     DEFAULT CURRENT_TIMESTAMP,
    `open_kfid`  varchar(100)       DEFAULT NULL,
    `to_userid`  varchar(100)       DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 0
  DEFAULT CHARSET = utf8;


CREATE TABLE `tab_message`
(
    `id`             int(10)   NOT NULL AUTO_INCREMENT,
    `created_at`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     timestamp NULL     DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`     timestamp NULL     DEFAULT CURRENT_TIMESTAMP,
    `talk_id`        INT(10)            DEFAULT NULL,
    `msg_from`       varchar(100)       DEFAULT NULL,
    `origin`         tinyint(4)         DEFAULT NULL,
    `msg_id`         varchar(100)       DEFAULT NULL,
    `msg_type`       varchar(16)        DEFAULT NULL,
    `send_time`      timestamp NULL     DEFAULT CURRENT_TIMESTAMP,
    `service_userid` varchar(100)       DEFAULT NULL,
    `content`        text               DEFAULT NULL,
    `raw`            text               DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 0
  DEFAULT CHARSET = utf8;