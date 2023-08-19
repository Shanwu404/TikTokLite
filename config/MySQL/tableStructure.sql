SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id`       bigint          NOT NULL AUTO_INCREMENT COMMENT '用户id，自增主键',
  `username` varchar(32)     NOT NULL                COMMENT '用户名',
  # 长度60是为了后续采用bcrypt存储密码考虑。明文应当是合法的ASCII字符。
  `password` varchar(60)     NOT NULL                COMMENT '用户密码',
  PRIMARY KEY (`id`),
  CONSTRAINT `username_unique` UNIQUE (`username`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  COMMENT = '用户表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos` (
  `id`           bigint          NOT NULL AUTO_INCREMENT COMMENT '视频唯一标识，自增主键',
  `author_id`    bigint          NOT NULL                COMMENT '视频作者id',
  `play_url`     varchar(255)    NOT NULL                COMMENT '视频播放地址',
  `cover_url`    varchar(255)    NOT NULL                COMMENT '视频封面地址',
  `publish_time` datetime        NOT NULL                COMMENT '发布日期',
  `title`        varchar(255)    DEFAULT NULL            COMMENT '视频名称',
  PRIMARY KEY (`id`),
  INDEX (`author_id`),
  INDEX (`publish_time`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  COMMENT = '视频表';

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `id`           bigint          NOT NULL AUTO_INCREMENT COMMENT '评论id，自增主键',
    `user_id`      bigint          NOT NULL                COMMENT '评论发布用户id',
    `video_id`     bigint          NOT NULL                COMMENT '评论视频id',
    `content`      varchar(500)    NOT NULL                COMMENT '评论内容',
    `create_date`  datetime        NOT NULL                COMMENT '评论发布时间',
    `canceled`     tinyint         NOT NULL DEFAULT '0'    COMMENT '评论发布默认置0，取消后置1',
    PRIMARY KEY (`id`),
    INDEX `video_id` (`video_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  COMMENT ='评论表';

-- ----------------------------
-- Table structure for likes
-- ----------------------------
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes` (
  `id`       bigint          NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id`  bigint          NOT NULL                COMMENT '执行点赞用户id',
  `video_id` bigint          NOT NULL                COMMENT '被点赞的视频id',
  `cancel`   tinyint         NOT NULL DEFAULT '0'    COMMENT '默认点赞为0，取消赞为1',
  PRIMARY KEY (`id`),
  # 这个联合唯一索引必要性需要澄清
  UNIQUE KEY `userIdtoVideoIdIdx` (`user_id`,`video_id`) USING BTREE,
  INDEX `user_id_index`  (`user_id`),
  INDEX `video_id_index` (`video_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  COMMENT = '点赞表';

-- ----------------------------
-- Table structure for follows
-- ----------------------------
DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows`
(
  `id`          bigint         NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `to_user_id`  bigint         NOT NULL                COMMENT '被关注用户id',
  `follower_id` bigint         NOT NULL                COMMENT '执行关注的用户id',
  PRIMARY KEY (`id`),
  INDEX (`to_user_id`),
  INDEX (`follower_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  COMMENT ='关注表';

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages`
(
    `id`           bigint          NOT NULL AUTO_INCREMENT COMMENT '消息id，自增主键',
    `to_user_id`   bigint          NOT NULL COMMENT '接收用户id',
    `from_user_id` bigint          NOT NULL COMMENT '发送用户id',
    `content`      varchar(500)    NOT NULL COMMENT '消息内容',
    `create_time`  datetime(4)        NOT NULL COMMENT '消息发送时间',
    PRIMARY KEY (`id`),
    INDEX (`from_user_id`, `to_user_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  COMMENT ='消息表';

