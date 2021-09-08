CREATE TABLE `tab_wx_kf_log`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL,
    `input`      text         DEFAULT NULL,
    `response`   varchar(255) DEFAULT NULL,
    `open_kfid`   varchar(255) DEFAULT NULL,
    `to_user_id`   int(10) unsigned DEFAULT NULL,
    `msg_id`      varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY          `idx_users_open_id` (`to_user_id`),
    KEY          `idx_kf_open_id` (`open_kfid`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
