CREATE TABLE `user`(
    `id` bigint(20) not null auto_increment,
    `user_id` bigint(20) not null ,
    `username` varchar(64) collate utf8mb4_general_ci not null,
    `password` varchar(64) collate utf8mb4_general_ci not null,
    `email` varchar(64) collate utf8mb4_general_ci ,
    `gender` tinyint(4) not null default '0',
    `create_time` timestamp null default current_timestamp,
    `update_time` timestamp null default current_timestamp on update current_timestamp,
    primary key (`id`),
    UNIQUE KEY `idx_username` (`username`) using btree ,
    UNIQUE KEY `idx_user_id` (`user_id`) using btree
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

DROP TABLE IF EXISTS `community`;
CREATE TABLE `community`(
    `id` int(11) not null auto_increment,
    `community_id` int(10) unsigned not null ,
    `community_name` varchar(128) collate utf8mb4_general_ci not null ,
    `introduction` varchar(256) collate utf8mb4_general_ci not null ,
    `create_time` timestamp not null default current_timestamp,
    `update_time` timestamp not null default current_timestamp on update current_timestamp,
    primary key (`id`),
    unique key `idx_community_id`(`community_id`),
    unique key `idx_community_name`(`community_name`)
)engine =InnoDB default charset =utf8mb4 collate=utf8mb4_general_ci;

INSERT INTO `community` values('1','1','go','goland','2025-01-01 08:10:10','2025-01-01 08:10:10'),
                              ('2','2','leetcode','刷题','2021-02-02 08:10:10','2021-02-02 08:10:10'),
                              ('3','3','java','dead','2022-02-02 08:10:10','2022-02-02 08:10:10'),
                              ('4','4','python','666','2023-02-02 08:10:10','2023-02-02 08:10:10');

drop table if exists `post`;
create table `post`(
    `id` bigint(20) not null auto_increment,
    `post_id` bigint(20) not null comment '帖子id',
    `title` varchar(128) collate utf8mb4_general_ci not null comment '标题',
    `content` varchar(8192) collate utf8mb4_general_ci not null comment '内容',
    `author_id` bigint(20) not null comment '作者的用户id',
    `community_id` bigint(20) not null comment '所属社区',
    `status` tinyint(4) not null default '1' comment '帖子状态',
    `create_time` timestamp null default current_timestamp  comment '创建时间',
    `update_time` timestamp null default current_timestamp on update current_timestamp comment '更新时间',
    primary key ('id'),
    unique key `idx_post_id` (`post_id`),
    key `idx_author_id` (`author_id`),
    key `idx_community_id` (`community_id`)
)engine=InnoDB default charset =utf8mb4 collate =utf8mb4_general_ci;