/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80016
 Source Host           : localhost:3306
 Source Schema         : bbs

 Target Server Type    : MySQL
 Target Server Version : 80016
 File Encoding         : 65001

 Date: 19/05/2022 22:11:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS `article`;
CREATE TABLE `article`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '论贴的标题，主题',
  `category_id` int(11) UNSIGNED NOT NULL COMMENT '分类id，每篇论贴隶属一个分类',
  `user_id` int(10) UNSIGNED NOT NULL COMMENT '创建者用户id',
  `status` int(11) NULL DEFAULT NULL COMMENT '状态，-1 已删除 1-未发布，2-已经发布',
  `support_count` int(11) NULL DEFAULT NULL COMMENT '点赞数量',
  `read_count` int(11) NULL DEFAULT NULL COMMENT '阅读数量',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '修改时间',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '论贴内容',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `title`(`title`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `foreign_category`(`category_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 45 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '论贴表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for article_follower
-- ----------------------------
DROP TABLE IF EXISTS `article_follower`;
CREATE TABLE `article_follower`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NULL DEFAULT NULL,
  `article_id` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `user_id`(`user_id`, `article_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 83 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(22) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '分类/主题名称',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '分类表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `article_id` int(11) NOT NULL,
  `user_id` int(11) NULL DEFAULT NULL COMMENT '评论者用户id',
  `replied_user_id` int(11) NULL DEFAULT -1 COMMENT '被回复者id, 默认0表示直接评论',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  `support_count` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '点赞数量',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '评论内容',
  `status` int(11) NULL DEFAULT 1 COMMENT '状态1正常，-1已删除',
  `comment_id` int(11) NULL DEFAULT NULL COMMENT '父级评论id',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `article_id`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 132 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '评论表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for follower
-- ----------------------------
DROP TABLE IF EXISTS `follower`;
CREATE TABLE `follower`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `followed_id` int(11) UNSIGNED NOT NULL COMMENT '被关注者id',
  `follower_id` int(11) UNSIGNED NOT NULL COMMENT '关注者id',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `followed_id`(`followed_id`, `follower_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '关注关系表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for label
-- ----------------------------
DROP TABLE IF EXISTS `label`;
CREATE TABLE `label`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '标签名',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  `status` int(11) NULL DEFAULT NULL COMMENT '状态 1 正常 -1已删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '标签表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for label_article
-- ----------------------------
DROP TABLE IF EXISTS `label_article`;
CREATE TABLE `label_article`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `label_id` int(11) NOT NULL,
  `article_id` int(10) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `label_id`(`label_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 232 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `status` int(8) NOT NULL DEFAULT 1 COMMENT '状态， -1已删除，1-正常使用',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  `creator_id` int(11) UNSIGNED NOT NULL COMMENT '消息创建者',
  `target_id` int(11) UNSIGNED NOT NULL COMMENT '目标用户id',
  `type` int(8) NULL DEFAULT NULL COMMENT '消息类型 1-未读， 2-已读',
  `kind` int(11) NULL DEFAULT NULL COMMENT '消息种类，1-评论 2-回复',
  `comment_id` int(11) NULL DEFAULT NULL COMMENT '评论id',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `target_id`(`target_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 76 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名',
  `account` varchar(22) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '账号',
  `password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '邮箱',
  `telephone_number` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '手机号',
  `status` int(11) NULL DEFAULT 1 COMMENT '状态， 1-正常使用 2-已注销 3-黑名单',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  `role` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '身份角色，3 -超级管理员，2-管理员，1-普通用户',
  `gender` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '性别，男女',
  `image_url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '头像路径',
  `birthday` datetime(0) NULL DEFAULT NULL COMMENT '生日',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `account`(`account`) USING BTREE,
  INDEX `account_password`(`account`, `password`) USING BTREE COMMENT '账号密码联合索引',
  INDEX `name`(`name`) USING BTREE,
  INDEX `email`(`email`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 70 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO `user`(`id`, `name`, `account`, `password`, `email`, `telephone_number`, `status`, `create_time`, `update_time`, `role`, `gender`,`birthday`) VALUES (1, 'admin', 'admin', '299f39bf4e1d6d2328717ae49f016172', '123@qq.com', '18212345721', 1, '2022-04-03 23:31:40', '2022-02-27 15:16:56', '3', '男', '2004-05-19 01:21:37');


INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (1, 'java', '2021-10-15 15:45:23', '2021-11-20 18:43:57', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (2, 'golang', '2021-10-15 15:45:56', '2021-11-20 19:23:26', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (3, 'c', '2021-10-15 15:46:02', '2021-10-15 18:29:38', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (4, 'c++', '2021-10-15 15:46:05', '2021-10-15 18:29:41', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (5, 'python', '2021-10-15 15:46:12', '2021-10-15 18:29:44', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (6, '其他标签', '2021-10-25 16:40:30', '2021-10-25 16:40:34', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (8, 'JavaScript', '2021-11-20 19:24:01', '2022-05-17 23:44:05', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (9, 'sql', '2021-11-20 19:28:54', '2022-05-14 15:57:54', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (10, 'linux', '2022-05-14 15:58:02', '2022-05-14 15:58:02', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (11, '旅游', '2022-05-14 16:08:02', '2022-05-14 16:08:02', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (12, '生活', '2022-05-17 21:23:32', '2022-05-17 21:23:32', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (13, '工具', '2022-05-17 23:42:11', '2022-05-17 23:42:37', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (14, '分享', '2022-05-18 00:48:56', '2022-05-18 00:48:56', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (15, '笔记', '2022-05-18 00:49:03', '2022-05-18 00:49:03', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (16, '前端', '2022-05-18 20:52:01', '2022-05-18 20:52:01', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (17, '后端', '2022-05-18 20:52:06', '2022-05-18 20:52:06', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (18, 'redis', '2022-05-18 20:54:54', '2022-05-18 20:54:54', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (19, '框架', '2022-05-18 21:40:52', '2022-05-18 21:40:52', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (20, '操作系统', '2022-05-18 22:11:05', '2022-05-18 22:12:08', 1);
INSERT INTO `label`(`id`, `name`, `create_time`, `update_time`, `status`) VALUES (21, '计网', '2022-05-19 00:27:52', '2022-05-19 00:27:52', 1);


INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (1, '生活', '2021-10-14 02:59:13', '2021-11-20 17:18:07');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (2, '工作', '2021-10-14 03:03:02', '2021-10-14 03:03:05');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (4, '大学', '2021-10-14 03:13:59', '2021-10-14 03:13:59');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (5, '程序员', '2021-10-14 03:37:08', '2021-10-14 03:38:57');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (6, '其他分类', '2021-10-25 16:40:08', '2021-11-20 17:25:14');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (8, '新分类', '2021-11-20 17:25:36', '2021-11-20 17:29:46');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (9, '学术', '2021-11-20 17:30:09', '2022-05-14 15:57:33');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (10, '家庭', '2022-05-14 15:57:20', '2022-05-14 15:57:20');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (11, '娱乐', '2022-05-14 15:58:53', '2022-05-14 15:58:53');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (12, '数据库', '2022-05-14 16:27:59', '2022-05-14 16:27:59');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (13, '影视', '2022-05-17 23:43:03', '2022-05-17 23:43:03');
INSERT INTO `category`(`id`, `name`, `create_time`, `update_time`) VALUES (14, '问题', '2022-05-18 01:50:08', '2022-05-18 01:50:08');
