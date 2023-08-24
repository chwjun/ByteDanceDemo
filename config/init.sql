CREATE TABLE `users` (
                         `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                         `created_at` datetime(3) DEFAULT NULL COMMENT '记录创建时间',
                         `updated_at` datetime(3) DEFAULT NULL COMMENT '记录更新时间',
                         `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
                         `name` varchar(191) NOT NULL COMMENT '用户名',
                         `password` varchar(191) NOT NULL COMMENT '用户密码',
                         `role` varchar(191) NOT NULL COMMENT '角色',
                         `avatar` varchar(191) DEFAULT 'http://yourserver.com/default_avatar.jpg' COMMENT '用户头像链接',
                         `background_image` varchar(191) DEFAULT 'http://yourserver.com/default_background.jpg' COMMENT '用户个人页顶部大图链接',
                         `signature` text COMMENT '个人简介',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `name` (`name`),
                         KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `casbin_rule` (
                               `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                               `p_type` varchar(100) DEFAULT NULL,
                               `v0` varchar(100) DEFAULT NULL,
                               `v1` varchar(100) DEFAULT NULL,
                               `v2` varchar(100) DEFAULT NULL,
                               `v3` varchar(100) DEFAULT NULL,
                               `v4` varchar(100) DEFAULT NULL,
                               `v5` varchar(100) DEFAULT NULL,
                               PRIMARY KEY (`id`),
                               UNIQUE INDEX idx_unique_casbin (p_type,v0,v1,v2)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `likes` (
                         `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                         `created_at` datetime(3) DEFAULT NULL COMMENT '记录创建时间',
                         `updated_at` datetime(3) DEFAULT NULL COMMENT '记录更新时间',
                         `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
                         `user_id` bigint(20) unsigned NOT NULL COMMENT '点赞用户id',
                         `video_id` bigint(20) unsigned NOT NULL COMMENT '点赞视频id',
                         `liked` bigint(20) NOT NULL DEFAULT '1' COMMENT '默认1表示已点赞，0表示未点赞',
                         PRIMARY KEY (`id`),
                         KEY `idx_likes_deleted_at` (`deleted_at`),
                         KEY `idx_user_video_liked` (`user_id`,`video_id`,`liked`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `comments` (
                            `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                            `created_at` datetime(3) DEFAULT NULL COMMENT '记录创建时间',
                            `updated_at` datetime(3) DEFAULT NULL COMMENT '记录更新时间',
                            `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
                            `user_id` bigint(20) unsigned NOT NULL COMMENT '发布评论的用户id',
                            `video_id` bigint(20) unsigned NOT NULL COMMENT '评论视频的id',
                            `content` varchar(191) NOT NULL COMMENT '评论的内容',
                            `action_type` enum('1','2') NOT NULL COMMENT '评论行为，1表示已发布评论，2表示删除评论',
                            PRIMARY KEY (`id`),
                            KEY `idx_comments_deleted_at` (`deleted_at`),
                            KEY `idx_comments_user_id` (`user_id`),
                            KEY `idx_video_action` (`video_id`,`action_type`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `messages` (
                            `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                            `created_at` datetime(3) DEFAULT NULL COMMENT '记录创建时间',
                            `updated_at` datetime(3) DEFAULT NULL COMMENT '记录更新时间',
                            `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
                            `sender_id` bigint(20) unsigned NOT NULL COMMENT '发送message的user id',
                            `receiver_id` bigint(20) unsigned NOT NULL COMMENT '接收message的user id',
                            `content` varchar(191) NOT NULL COMMENT '消息内容',
                            `action_type` int(11) DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            KEY `idx_messages_deleted_at` (`deleted_at`),
                            KEY `idx_sender_receiver` (`sender_id`,`receiver_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `videos` (
                          `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                          `created_at` datetime(3) DEFAULT NULL COMMENT '记录创建时间',
                          `updated_at` datetime(3) DEFAULT NULL COMMENT '记录更新时间',
                          `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
                          `author_id` bigint(20) unsigned NOT NULL COMMENT '视频作者id',
                          `title` varchar(191) NOT NULL COMMENT '视频标题',
                          `play_url` varchar(191) NOT NULL COMMENT '视频播放地址',
                          `cover_url` varchar(191) NOT NULL COMMENT '视频封面地址',
                          PRIMARY KEY (`id`),
                          KEY `idx_videos_deleted_at` (`deleted_at`),
                          KEY `idx_videos_author_id` (`author_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `relations` (
                             `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                             `created_at` datetime(3) DEFAULT NULL COMMENT '记录创建时间',
                             `updated_at` datetime(3) DEFAULT NULL COMMENT '记录更新时间',
                             `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
                             `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
                             `following_id` bigint(20) unsigned NOT NULL COMMENT 'user id关注的用户id',
                             `followed` bigint(20) NOT NULL DEFAULT '0' COMMENT '默认0表示未关注，1表示已关注',
                             PRIMARY KEY (`id`),
                             KEY `idx_relations_deleted_at` (`deleted_at`),
                             KEY `idx_user_following_followed` (`user_id`,`following_id`,`followed`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;