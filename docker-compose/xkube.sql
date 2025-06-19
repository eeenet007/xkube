/*
Navicat MySQL Data Transfer

Source Server         : 236.115
Source Server Version : 80029
Source Host           : 192.168.236.115:3306
Source Database       : db_xkube

Target Server Type    : MYSQL
Target Server Version : 80029
File Encoding         : 65001

Date: 2025-06-19 14:38:54
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `title` varchar(100) NOT NULL DEFAULT '',
  `status` int NOT NULL DEFAULT '2',
  `sort` int NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of group
-- ----------------------------
INSERT INTO `group` VALUES ('1', 'rbac', '权限管理', '2', '1');
INSERT INTO `group` VALUES ('3', 'xkube', 'k8s管理', '2', '2');
INSERT INTO `group` VALUES ('4', 'wiki', '文档中心', '2', '3');
INSERT INTO `group` VALUES ('5', 'cicd', 'CICD', '2', '4');
INSERT INTO `group` VALUES ('6', 'fav', '我的收藏', '2', '5');

-- ----------------------------
-- Table structure for node
-- ----------------------------
DROP TABLE IF EXISTS `node`;
CREATE TABLE `node` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '',
  `name` varchar(100) NOT NULL DEFAULT '',
  `level` int NOT NULL DEFAULT '1',
  `pid` bigint NOT NULL DEFAULT '0',
  `icons` varchar(200) DEFAULT NULL,
  `sorts` int DEFAULT '1',
  `remark` varchar(200) DEFAULT NULL,
  `status` int NOT NULL DEFAULT '2',
  `group_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=263 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of node
-- ----------------------------
INSERT INTO `node` VALUES ('1', '权限管理', 'rbac', '1', '0', 'layui-icon-auz', '7', '', '2', '1');
INSERT INTO `node` VALUES ('2', '目录结构', 'node', '2', '1', null, '4', '', '2', '1');
INSERT INTO `node` VALUES ('3', '列表', 'List', '3', '2', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('4', '增改', 'AddAndEdit', '3', '2', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('5', '删除', 'Delete', '3', '2', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('6', '管理员', 'user', '2', '1', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('7', '管理员列表', 'List', '3', '6', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('8', '添加管理员', 'Add', '3', '6', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('9', '更新管理员', 'Update', '3', '6', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('10', '删除管理员', 'Delete', '3', '6', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('11', '目录分组', 'group', '2', '1', null, '3', '', '2', '1');
INSERT INTO `node` VALUES ('12', '分组列表', 'List', '3', '11', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('13', '添加分组', 'Add', '3', '11', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('14', '修改分组', 'Update', '3', '11', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('15', '删除分组', 'Delete', '3', '11', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('16', '角色管理', 'role', '2', '1', null, '2', '', '2', '1');
INSERT INTO `node` VALUES ('17', '角色列表', 'List', '3', '16', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('18', '添加编辑角色', 'AddAndEdit', '3', '16', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('19', '删除角色', 'Delete', '3', '16', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('20', 'get roles', 'Getlist', '3', '16', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('21', 'show access', 'AccessToNode', '3', '16', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('22', 'add accsee', 'AddAccess', '3', '16', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('23', 'show role to userlist', 'RoleToUserList', '3', '16', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('24', 'add role to user', 'AddRoleToUser', '3', '16', null, '1', '', '2', '1');
INSERT INTO `node` VALUES ('27', '获取目录列表', 'Getlist', '3', '2', null, '1', 'add by kang', '2', '1');
INSERT INTO `node` VALUES ('28', '获取父节点id', 'GetPid', '3', '2', null, '1', 'add by kang', '2', '1');
INSERT INTO `node` VALUES ('29', 'RoleToNodeList', 'RoleToNodeList', '3', '16', null, '1', 'add by kang', '2', '1');
INSERT INTO `node` VALUES ('30', 'DelRoleToUser', 'DelRoleToUser', '3', '16', null, '1', 'add by kang', '2', '1');
INSERT INTO `node` VALUES ('31', 'DelRoleToNode', 'DelRoleToNode', '3', '16', null, '1', 'add by kang', '2', '1');
INSERT INTO `node` VALUES ('34', 'k8s管理', 'xkube', '1', '0', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('35', '无状态', 'deploy', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('36', '列表', 'v1/List', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('37', '集群管理', 'cluster', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('38', '集群列表', 'v1/List', '3', '37', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('39', '有状态', 'sts', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('40', '守护进程', 'ds', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('41', '任务', 'job', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('42', '定时任务', 'cronjob', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('43', '容器组', 'pod', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('44', '自定义资源', 'cdr', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('45', '自动伸缩', 'hpa', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('46', '事件中心', 'event', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('47', 'yam操作', 'apply', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('48', '命名空间', 'ns', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('49', '节点管理', 'node', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('50', '服务管理', 'svc', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('51', '路由', 'ing', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('52', 'configmap', 'cm', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('53', 'secrets', 'secret', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('54', '存储声明', 'pvc', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('55', '存储卷', 'pv', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('56', '存储类', 'storageclass', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('57', 'clusterrolebinding', 'clusterrolebinding', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('58', 'clusterroles', 'clusterroles', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('59', 'serviceaccounts', 'serviceaccounts', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('60', 'roles', 'roles', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('61', 'rolebinding', 'rolebinding', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('62', 'metrics', 'metrics', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('65', '文档中心', 'wiki', '1', '0', '', '0', '', '2', '4');
INSERT INTO `node` VALUES ('66', '列表', 'v1/List', '2', '65', '', '0', '', '2', '4');
INSERT INTO `node` VALUES ('67', '添加', 'v1/Add', '2', '65', '', '0', '', '2', '4');
INSERT INTO `node` VALUES ('68', '更新', 'v1/Update', '2', '65', '', '0', '', '2', '4');
INSERT INTO `node` VALUES ('69', '删除', 'v1/Del', '2', '65', '', '0', '', '2', '4');
INSERT INTO `node` VALUES ('70', '读取', 'v1/Read', '2', '65', '', '0', '', '2', '4');
INSERT INTO `node` VALUES ('71', '加密判断', 'v1/ReadEncry', '2', '65', '', '0', '', '2', '4');
INSERT INTO `node` VALUES ('72', 'GetridByuid', 'GetridByuid', '3', '16', '', '0', '', '2', '1');
INSERT INTO `node` VALUES ('73', '集群授权', 'cluster', '2', '1', '', '0', '', '2', '1');
INSERT INTO `node` VALUES ('74', '列表', 'List', '3', '73', '', '0', '', '2', '1');
INSERT INTO `node` VALUES ('75', '添加', 'Add', '3', '73', '', '0', '', '2', '1');
INSERT INTO `node` VALUES ('76', '删除', 'Delete', '3', '73', '', '0', '', '2', '1');
INSERT INTO `node` VALUES ('77', '添加', 'v1/Add', '3', '37', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('78', '更新', 'v1/Update', '3', '37', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('79', '编辑', 'v1/Edit', '3', '37', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('80', '删除', 'v1/Del', '3', '37', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('81', '详情', 'v1/Detail', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('82', '创建', 'v1/Create', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('83', '更新', 'v1/Modify', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('84', 'yaml更新', 'v1/ModifyByYaml', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('85', '删除', 'v1/Del', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('86', 'yaml查看', 'v1/Yaml', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('87', 'ReplicasetYaml', 'v1/ReplicasetYaml', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('88', 'Replicaset列表', 'v1/ReplicasetList', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('89', '回滚', 'v1/RollBack', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('90', '重启', 'v1/Restart', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('91', '标签', 'v1/Labels', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('92', '镜像', 'v1/Image', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('93', '列表', 'v1/List', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('94', '详情', 'v1/Detail', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('95', '创建', 'v1/Create', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('96', '更新', 'v1/Modify', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('97', '删除', 'v1/Del', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('98', 'yaml查看', 'v1/Yaml', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('99', '回滚', 'v1/RollBack', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('100', '重启', 'v1/Restart', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('101', '列表', 'v1/List', '3', '40', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('102', '详情', 'v1/Detail', '3', '40', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('103', '创建', 'v1/Create', '3', '40', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('104', '更新', 'v1/Modify', '3', '40', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('105', '删除', 'v1/Del', '3', '40', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('106', 'yaml查看', 'v1/Yaml', '3', '40', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('107', 'yaml更新', 'v1/ModifyByYaml', '3', '40', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('108', '列表', 'v1/List', '3', '41', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('109', '详情', 'v1/Detail', '3', '41', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('110', '创建', 'v1/Create', '3', '41', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('111', '更新', 'v1/Modify', '3', '41', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('112', '删除', 'v1/Del', '3', '41', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('113', 'yaml查看', 'v1/Yaml', '3', '41', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('114', '日志', 'v1/Log', '3', '41', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('115', '列表', 'v1/List', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('116', '详情', 'v1/Detail', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('117', '创建', 'v1/Create', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('118', '更新', 'v1/Modify', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('119', 'yaml更新', 'v1/ModifyByYaml', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('120', '删除', 'v1/Del', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('121', '查看yaml', 'v1/Yaml', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('122', '标签', 'v1/Labels', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('123', '列表', 'v1/List', '3', '50', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('124', '详情', 'v1/Detail', '3', '50', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('125', '创建', 'v1/Create', '3', '50', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('126', 'yaml更新', 'v1/ModifyByYaml', '3', '50', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('127', '删除', 'v1/Del', '3', '50', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('128', 'yaml查看', 'v1/Yaml', '3', '50', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('129', '列表', 'v1/List', '3', '52', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('130', '详情', 'v1/Detail', '3', '52', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('131', '创建', 'v1/Create', '3', '52', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('132', 'yaml更新', 'v1/ModifyByYaml', '3', '52', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('133', '删除', 'v1/Del', '3', '52', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('134', 'yaml查看', 'v1/Yaml', '3', '52', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('135', 'yaml查看', 'v1/Yaml', '3', '53', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('136', '列表', 'v1/List', '3', '53', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('137', '详情', 'v1/Detail', '3', '53', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('138', '创建', 'v1/Create', '3', '53', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('139', 'yaml更新', 'v1/ModifyByYaml', '3', '53', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('140', '删除', 'v1/Del', '3', '53', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('141', '列表', 'v1/List', '3', '51', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('142', '详情', 'v1/Detail', '3', '51', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('143', '创建', 'v1/Create', '3', '51', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('144', 'yaml更新', 'v1/ModifyByYaml', '3', '51', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('145', '删除', 'v1/Del', '3', '51', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('146', 'yaml查看', 'v1/Yaml', '3', '51', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('147', '列表', 'v1/List', '3', '43', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('148', 'container列表', 'v1/ContainerList', '3', '43', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('149', '详情', 'v1/Detail', '3', '43', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('150', '日志', 'v1/Log', '3', '43', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('151', '删除', 'v1/Del', '3', '43', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('152', 'yaml查看', 'v1/Yaml', '3', '43', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('153', '列表', 'v1/List', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('154', 'yaml查看', 'v1/Yaml', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('155', '创建', 'v1/Create', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('156', '删除', 'v1/Del', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('157', 'PodList', 'v1/PodList', '3', '62', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('158', 'PodUsage', 'PodUsage', '3', '62', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('159', 'NodeUsage', 'NodeUsage', '3', '62', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('160', 'NodeList', 'v1/NodeList', '3', '62', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('161', '列表', 'v1/List', '3', '49', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('162', '节点池列表', 'v1/PoolList', '3', '49', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('163', '详情', 'v1/Detail', '3', '49', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('164', 'yaml查看', 'v1/Yaml', '3', '49', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('165', '不调度', 'v1/Unschedulable', '3', '49', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('166', '排水', 'v1/Drain', '3', '49', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('167', '删除', 'v1/Del', '3', '49', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('168', '列表', 'v1/List', '3', '48', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('169', '详情', 'v1/Detail', '3', '48', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('170', '创建', 'v1/Create', '3', '48', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('171', '删除', 'v1/Del', '3', '48', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('172', 'yaml更新', 'v1/ModifyByYaml', '3', '48', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('173', 'yaml查看', 'v1/Yaml', '3', '48', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('174', '资源限制', 'v1/LimitRange', '3', '48', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('175', 'yaml创建', 'v1/CreateByYaml', '3', '47', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('176', 'ApplyYaml', 'v1/ApplyYaml', '3', '47', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('177', '列表', 'v1/List', '3', '44', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('178', 'yaml查看', 'v1/Yaml', '3', '44', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('179', '删除', 'v1/Del', '3', '44', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('180', '列表', 'v1/List', '3', '46', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('181', '应用集', 'appname', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('182', '列表', 'v1/List', '3', '56', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('183', 'yaml查看', 'v1/Yaml', '3', '56', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('184', '详情', 'v1/Detail', '3', '56', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('185', '列表', 'v1/List', '3', '55', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('186', '详情', 'v1/Detail', '3', '55', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('187', 'yaml查看', 'v1/Yaml', '3', '55', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('188', '列表', 'v1/List', '3', '54', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('189', 'yaml查看', 'v1/Yaml', '3', '54', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('190', '详情', 'v1/Detail', '3', '54', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('191', '列表', 'v1/List', '3', '57', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('192', '列表', 'v1/List', '3', '58', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('193', '列表', 'v1/List', '3', '59', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('194', '列表', 'v1/List', '3', '60', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('195', '列表', 'v1/List', '3', '61', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('196', 'yaml查看', 'v1/Yaml', '3', '57', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('197', 'yaml查看', 'v1/Yaml', '3', '58', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('198', 'yaml查看', 'v1/Yaml', '3', '59', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('199', 'yaml查看', 'v1/Yaml', '3', '60', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('200', 'yaml查看', 'v1/Yaml', '3', '61', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('201', '首页统计', 'v1/Count', '3', '37', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('202', 'beta列表', 'beta1/List', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('203', 'beta详情', 'beta1/Detail', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('204', 'beta创建', 'beta1/Create', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('205', 'beta更新', 'beta1/Modify', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('206', 'beta yaml更新', 'beta1/ModifyByYaml', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('207', 'beta删除', 'beta1/Del', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('208', 'beta查看yaml', 'beta1/Yaml', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('209', 'beta标签', 'beta1/Labels', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('210', '应用集列表', 'v1/List', '3', '181', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('211', '添加', 'v1/Add', '3', '181', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('212', '编辑', 'v1/Edit', '3', '181', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('213', '更新', 'v1/Update', '3', '181', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('214', '删除', 'v1/Del', '3', '181', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('215', 'CICD', 'cicd', '1', '0', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('216', '列表', 'v1/List', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('217', '添加', 'v1/Add', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('218', '更新', 'v1/Update', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('219', '编辑', 'v1/Edit', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('220', '删除', 'v1/Del', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('221', 'ak列表', 'v1/AkList', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('222', 'ak添加', 'v1/AkAdd', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('223', 'ak删除', 'v1/AkDel', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('224', 'cicd信息', 'v1/GetCicdInfo', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('225', '流水线信息', 'v1/GetPipelines', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('226', '应用名列表', 'v1/ListAppname', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('227', '更改cicd状态', 'v1/PostStatus', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('228', '阿里云id列表', 'v1/GetAliyunIdList', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('229', '运行流水线', 'pipeline/Start', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('230', '流水线运行列表', 'pipeline/ListRun', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('231', '获取流水线运行实例', 'pipeline/GetRun', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('232', '获取流水线日志', 'pipeline/GetJobLog', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('233', '通过AK获取组织ID', 'v1/GetOrganizationsByAk', '2', '215', '', '0', '', '2', '5');
INSERT INTO `node` VALUES ('234', '克隆', 'v1/Clone', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('235', '克隆', 'v1/Clone', '3', '39', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('236', '克隆', 'v1/Clone', '3', '50', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('237', '克隆', 'v1/Clone', '3', '52', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('238', '克隆', 'v1/Clone', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('239', '克隆', 'v1/Clone', '3', '53', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('240', 'beta1克隆', 'beta1/Clone', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('241', '副本', 'v1/Replicas', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('242', '升级策略', 'v1/Strategy', '3', '35', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('243', '运行一次', 'v1/Run', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('244', '运行一次', 'beta1/Run', '3', '42', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('245', 'websocket', 'terminal/ws', '3', '43', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('246', '克隆迁移', 'clone', '2', '34', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('247', '任务列表', 'v1/List', '3', '246', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('248', '集群详情', 'v1/Detail', '3', '37', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('249', '列表', 'v2beta2/List', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('250', 'yaml查看', 'v2beta2/Yaml', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('251', 'yaml更新', 'v2beta2/ModifyByYaml', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('252', '创建', 'v2beta2/Create', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('253', '删除', 'v2beta2/Del', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('254', 'yaml更新', 'v1/ModifyByYaml', '3', '45', '', '0', '', '2', '3');
INSERT INTO `node` VALUES ('255', '添加', 'Add', '3', '2', '', '0', '', '2', '1');
INSERT INTO `node` VALUES ('256', '获取角色id', 'GetridByuid', '3', '16', '', '0', '', '2', '1');
INSERT INTO `node` VALUES ('257', '上传', 'v1/Upload', '2', '65', '', '0', '', '2', '4');
INSERT INTO `node` VALUES ('258', '我的收藏', 'fav', '1', '0', '', '0', '', '2', '6');
INSERT INTO `node` VALUES ('259', '列表', 'v1/List', '2', '258', '', '0', '', '2', '6');
INSERT INTO `node` VALUES ('260', '添加', 'v1/Add', '2', '258', '', '0', '', '2', '6');
INSERT INTO `node` VALUES ('261', '删除', 'v1/Del', '2', '258', '', '0', '', '2', '6');
INSERT INTO `node` VALUES ('262', '我的集群', 'MyClusterList', '3', '73', '', '0', '', '2', '1');

-- ----------------------------
-- Table structure for node_roles
-- ----------------------------
DROP TABLE IF EXISTS `node_roles`;
CREATE TABLE `node_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `node_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=927 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of node_roles
-- ----------------------------
INSERT INTO `node_roles` VALUES ('128', '32', '4');
INSERT INTO `node_roles` VALUES ('146', '1', '10');
INSERT INTO `node_roles` VALUES ('147', '2', '10');
INSERT INTO `node_roles` VALUES ('148', '3', '10');
INSERT INTO `node_roles` VALUES ('149', '4', '10');
INSERT INTO `node_roles` VALUES ('150', '5', '10');
INSERT INTO `node_roles` VALUES ('151', '27', '10');
INSERT INTO `node_roles` VALUES ('152', '28', '10');
INSERT INTO `node_roles` VALUES ('153', '6', '10');
INSERT INTO `node_roles` VALUES ('154', '7', '10');
INSERT INTO `node_roles` VALUES ('155', '8', '10');
INSERT INTO `node_roles` VALUES ('156', '9', '10');
INSERT INTO `node_roles` VALUES ('157', '10', '10');
INSERT INTO `node_roles` VALUES ('158', '11', '10');
INSERT INTO `node_roles` VALUES ('159', '12', '10');
INSERT INTO `node_roles` VALUES ('160', '13', '10');
INSERT INTO `node_roles` VALUES ('161', '14', '10');
INSERT INTO `node_roles` VALUES ('162', '15', '10');
INSERT INTO `node_roles` VALUES ('163', '16', '10');
INSERT INTO `node_roles` VALUES ('164', '17', '10');
INSERT INTO `node_roles` VALUES ('165', '18', '10');
INSERT INTO `node_roles` VALUES ('166', '19', '10');
INSERT INTO `node_roles` VALUES ('167', '20', '10');
INSERT INTO `node_roles` VALUES ('168', '21', '10');
INSERT INTO `node_roles` VALUES ('169', '22', '10');
INSERT INTO `node_roles` VALUES ('170', '23', '10');
INSERT INTO `node_roles` VALUES ('171', '24', '10');
INSERT INTO `node_roles` VALUES ('172', '29', '10');
INSERT INTO `node_roles` VALUES ('173', '30', '10');
INSERT INTO `node_roles` VALUES ('174', '31', '10');
INSERT INTO `node_roles` VALUES ('175', '72', '10');
INSERT INTO `node_roles` VALUES ('176', '73', '10');
INSERT INTO `node_roles` VALUES ('177', '74', '10');
INSERT INTO `node_roles` VALUES ('178', '75', '10');
INSERT INTO `node_roles` VALUES ('179', '76', '10');
INSERT INTO `node_roles` VALUES ('197', '34', '9');
INSERT INTO `node_roles` VALUES ('198', '35', '9');
INSERT INTO `node_roles` VALUES ('199', '85', '9');
INSERT INTO `node_roles` VALUES ('200', '37', '9');
INSERT INTO `node_roles` VALUES ('201', '80', '9');
INSERT INTO `node_roles` VALUES ('202', '39', '9');
INSERT INTO `node_roles` VALUES ('203', '97', '9');
INSERT INTO `node_roles` VALUES ('204', '40', '9');
INSERT INTO `node_roles` VALUES ('205', '105', '9');
INSERT INTO `node_roles` VALUES ('206', '41', '9');
INSERT INTO `node_roles` VALUES ('207', '112', '9');
INSERT INTO `node_roles` VALUES ('208', '42', '9');
INSERT INTO `node_roles` VALUES ('209', '120', '9');
INSERT INTO `node_roles` VALUES ('210', '43', '9');
INSERT INTO `node_roles` VALUES ('211', '151', '9');
INSERT INTO `node_roles` VALUES ('212', '44', '9');
INSERT INTO `node_roles` VALUES ('213', '179', '9');
INSERT INTO `node_roles` VALUES ('214', '45', '9');
INSERT INTO `node_roles` VALUES ('215', '156', '9');
INSERT INTO `node_roles` VALUES ('216', '48', '9');
INSERT INTO `node_roles` VALUES ('217', '171', '9');
INSERT INTO `node_roles` VALUES ('218', '49', '9');
INSERT INTO `node_roles` VALUES ('219', '166', '9');
INSERT INTO `node_roles` VALUES ('220', '167', '9');
INSERT INTO `node_roles` VALUES ('221', '50', '9');
INSERT INTO `node_roles` VALUES ('222', '127', '9');
INSERT INTO `node_roles` VALUES ('223', '51', '9');
INSERT INTO `node_roles` VALUES ('224', '145', '9');
INSERT INTO `node_roles` VALUES ('225', '52', '9');
INSERT INTO `node_roles` VALUES ('226', '133', '9');
INSERT INTO `node_roles` VALUES ('227', '53', '9');
INSERT INTO `node_roles` VALUES ('228', '140', '9');
INSERT INTO `node_roles` VALUES ('229', '65', '9');
INSERT INTO `node_roles` VALUES ('230', '69', '9');
INSERT INTO `node_roles` VALUES ('231', '34', '8');
INSERT INTO `node_roles` VALUES ('232', '35', '8');
INSERT INTO `node_roles` VALUES ('233', '82', '8');
INSERT INTO `node_roles` VALUES ('234', '37', '8');
INSERT INTO `node_roles` VALUES ('235', '77', '8');
INSERT INTO `node_roles` VALUES ('236', '39', '8');
INSERT INTO `node_roles` VALUES ('237', '95', '8');
INSERT INTO `node_roles` VALUES ('238', '40', '8');
INSERT INTO `node_roles` VALUES ('239', '103', '8');
INSERT INTO `node_roles` VALUES ('240', '41', '8');
INSERT INTO `node_roles` VALUES ('241', '110', '8');
INSERT INTO `node_roles` VALUES ('242', '42', '8');
INSERT INTO `node_roles` VALUES ('243', '117', '8');
INSERT INTO `node_roles` VALUES ('244', '45', '8');
INSERT INTO `node_roles` VALUES ('245', '155', '8');
INSERT INTO `node_roles` VALUES ('246', '47', '8');
INSERT INTO `node_roles` VALUES ('247', '175', '8');
INSERT INTO `node_roles` VALUES ('248', '176', '8');
INSERT INTO `node_roles` VALUES ('249', '48', '8');
INSERT INTO `node_roles` VALUES ('250', '170', '8');
INSERT INTO `node_roles` VALUES ('251', '50', '8');
INSERT INTO `node_roles` VALUES ('252', '125', '8');
INSERT INTO `node_roles` VALUES ('253', '51', '8');
INSERT INTO `node_roles` VALUES ('254', '143', '8');
INSERT INTO `node_roles` VALUES ('255', '52', '8');
INSERT INTO `node_roles` VALUES ('256', '131', '8');
INSERT INTO `node_roles` VALUES ('257', '53', '8');
INSERT INTO `node_roles` VALUES ('258', '138', '8');
INSERT INTO `node_roles` VALUES ('259', '65', '8');
INSERT INTO `node_roles` VALUES ('260', '67', '8');
INSERT INTO `node_roles` VALUES ('261', '34', '7');
INSERT INTO `node_roles` VALUES ('262', '35', '7');
INSERT INTO `node_roles` VALUES ('263', '83', '7');
INSERT INTO `node_roles` VALUES ('264', '84', '7');
INSERT INTO `node_roles` VALUES ('265', '37', '7');
INSERT INTO `node_roles` VALUES ('266', '78', '7');
INSERT INTO `node_roles` VALUES ('267', '79', '7');
INSERT INTO `node_roles` VALUES ('268', '39', '7');
INSERT INTO `node_roles` VALUES ('269', '96', '7');
INSERT INTO `node_roles` VALUES ('270', '40', '7');
INSERT INTO `node_roles` VALUES ('271', '104', '7');
INSERT INTO `node_roles` VALUES ('272', '107', '7');
INSERT INTO `node_roles` VALUES ('273', '41', '7');
INSERT INTO `node_roles` VALUES ('274', '111', '7');
INSERT INTO `node_roles` VALUES ('275', '42', '7');
INSERT INTO `node_roles` VALUES ('276', '118', '7');
INSERT INTO `node_roles` VALUES ('277', '119', '7');
INSERT INTO `node_roles` VALUES ('278', '48', '7');
INSERT INTO `node_roles` VALUES ('279', '172', '7');
INSERT INTO `node_roles` VALUES ('280', '50', '7');
INSERT INTO `node_roles` VALUES ('281', '126', '7');
INSERT INTO `node_roles` VALUES ('282', '144', '7');
INSERT INTO `node_roles` VALUES ('283', '52', '7');
INSERT INTO `node_roles` VALUES ('284', '132', '7');
INSERT INTO `node_roles` VALUES ('285', '53', '7');
INSERT INTO `node_roles` VALUES ('286', '139', '7');
INSERT INTO `node_roles` VALUES ('287', '65', '7');
INSERT INTO `node_roles` VALUES ('288', '68', '7');
INSERT INTO `node_roles` VALUES ('289', '34', '6');
INSERT INTO `node_roles` VALUES ('290', '35', '6');
INSERT INTO `node_roles` VALUES ('291', '36', '6');
INSERT INTO `node_roles` VALUES ('292', '81', '6');
INSERT INTO `node_roles` VALUES ('293', '86', '6');
INSERT INTO `node_roles` VALUES ('294', '87', '6');
INSERT INTO `node_roles` VALUES ('295', '88', '6');
INSERT INTO `node_roles` VALUES ('296', '90', '6');
INSERT INTO `node_roles` VALUES ('297', '91', '6');
INSERT INTO `node_roles` VALUES ('298', '92', '6');
INSERT INTO `node_roles` VALUES ('299', '37', '6');
INSERT INTO `node_roles` VALUES ('300', '38', '6');
INSERT INTO `node_roles` VALUES ('302', '201', '6');
INSERT INTO `node_roles` VALUES ('303', '39', '6');
INSERT INTO `node_roles` VALUES ('304', '93', '6');
INSERT INTO `node_roles` VALUES ('305', '94', '6');
INSERT INTO `node_roles` VALUES ('306', '98', '6');
INSERT INTO `node_roles` VALUES ('307', '99', '6');
INSERT INTO `node_roles` VALUES ('308', '100', '6');
INSERT INTO `node_roles` VALUES ('309', '40', '6');
INSERT INTO `node_roles` VALUES ('310', '101', '6');
INSERT INTO `node_roles` VALUES ('311', '102', '6');
INSERT INTO `node_roles` VALUES ('312', '106', '6');
INSERT INTO `node_roles` VALUES ('313', '41', '6');
INSERT INTO `node_roles` VALUES ('314', '108', '6');
INSERT INTO `node_roles` VALUES ('315', '109', '6');
INSERT INTO `node_roles` VALUES ('316', '113', '6');
INSERT INTO `node_roles` VALUES ('317', '114', '6');
INSERT INTO `node_roles` VALUES ('318', '42', '6');
INSERT INTO `node_roles` VALUES ('319', '115', '6');
INSERT INTO `node_roles` VALUES ('320', '116', '6');
INSERT INTO `node_roles` VALUES ('321', '121', '6');
INSERT INTO `node_roles` VALUES ('322', '122', '6');
INSERT INTO `node_roles` VALUES ('323', '43', '6');
INSERT INTO `node_roles` VALUES ('324', '147', '6');
INSERT INTO `node_roles` VALUES ('325', '148', '6');
INSERT INTO `node_roles` VALUES ('326', '149', '6');
INSERT INTO `node_roles` VALUES ('327', '150', '6');
INSERT INTO `node_roles` VALUES ('328', '152', '6');
INSERT INTO `node_roles` VALUES ('329', '44', '6');
INSERT INTO `node_roles` VALUES ('330', '177', '6');
INSERT INTO `node_roles` VALUES ('331', '178', '6');
INSERT INTO `node_roles` VALUES ('332', '45', '6');
INSERT INTO `node_roles` VALUES ('333', '153', '6');
INSERT INTO `node_roles` VALUES ('334', '154', '6');
INSERT INTO `node_roles` VALUES ('335', '46', '6');
INSERT INTO `node_roles` VALUES ('336', '180', '6');
INSERT INTO `node_roles` VALUES ('337', '48', '6');
INSERT INTO `node_roles` VALUES ('338', '168', '6');
INSERT INTO `node_roles` VALUES ('339', '169', '6');
INSERT INTO `node_roles` VALUES ('340', '173', '6');
INSERT INTO `node_roles` VALUES ('341', '174', '6');
INSERT INTO `node_roles` VALUES ('342', '49', '6');
INSERT INTO `node_roles` VALUES ('343', '161', '6');
INSERT INTO `node_roles` VALUES ('344', '162', '6');
INSERT INTO `node_roles` VALUES ('345', '163', '6');
INSERT INTO `node_roles` VALUES ('346', '164', '6');
INSERT INTO `node_roles` VALUES ('347', '50', '6');
INSERT INTO `node_roles` VALUES ('348', '123', '6');
INSERT INTO `node_roles` VALUES ('349', '124', '6');
INSERT INTO `node_roles` VALUES ('350', '128', '6');
INSERT INTO `node_roles` VALUES ('351', '51', '6');
INSERT INTO `node_roles` VALUES ('352', '141', '6');
INSERT INTO `node_roles` VALUES ('353', '142', '6');
INSERT INTO `node_roles` VALUES ('354', '146', '6');
INSERT INTO `node_roles` VALUES ('355', '52', '6');
INSERT INTO `node_roles` VALUES ('356', '129', '6');
INSERT INTO `node_roles` VALUES ('357', '130', '6');
INSERT INTO `node_roles` VALUES ('358', '134', '6');
INSERT INTO `node_roles` VALUES ('359', '53', '6');
INSERT INTO `node_roles` VALUES ('360', '135', '6');
INSERT INTO `node_roles` VALUES ('361', '136', '6');
INSERT INTO `node_roles` VALUES ('362', '137', '6');
INSERT INTO `node_roles` VALUES ('363', '54', '6');
INSERT INTO `node_roles` VALUES ('364', '188', '6');
INSERT INTO `node_roles` VALUES ('365', '189', '6');
INSERT INTO `node_roles` VALUES ('366', '190', '6');
INSERT INTO `node_roles` VALUES ('367', '55', '6');
INSERT INTO `node_roles` VALUES ('368', '185', '6');
INSERT INTO `node_roles` VALUES ('369', '186', '6');
INSERT INTO `node_roles` VALUES ('370', '187', '6');
INSERT INTO `node_roles` VALUES ('371', '56', '6');
INSERT INTO `node_roles` VALUES ('372', '182', '6');
INSERT INTO `node_roles` VALUES ('373', '183', '6');
INSERT INTO `node_roles` VALUES ('374', '184', '6');
INSERT INTO `node_roles` VALUES ('375', '57', '6');
INSERT INTO `node_roles` VALUES ('376', '191', '6');
INSERT INTO `node_roles` VALUES ('377', '196', '6');
INSERT INTO `node_roles` VALUES ('378', '58', '6');
INSERT INTO `node_roles` VALUES ('379', '192', '6');
INSERT INTO `node_roles` VALUES ('380', '197', '6');
INSERT INTO `node_roles` VALUES ('381', '59', '6');
INSERT INTO `node_roles` VALUES ('382', '193', '6');
INSERT INTO `node_roles` VALUES ('383', '198', '6');
INSERT INTO `node_roles` VALUES ('384', '60', '6');
INSERT INTO `node_roles` VALUES ('385', '194', '6');
INSERT INTO `node_roles` VALUES ('386', '199', '6');
INSERT INTO `node_roles` VALUES ('387', '61', '6');
INSERT INTO `node_roles` VALUES ('388', '195', '6');
INSERT INTO `node_roles` VALUES ('389', '200', '6');
INSERT INTO `node_roles` VALUES ('390', '62', '6');
INSERT INTO `node_roles` VALUES ('391', '157', '6');
INSERT INTO `node_roles` VALUES ('392', '158', '6');
INSERT INTO `node_roles` VALUES ('393', '159', '6');
INSERT INTO `node_roles` VALUES ('394', '160', '6');
INSERT INTO `node_roles` VALUES ('395', '181', '6');
INSERT INTO `node_roles` VALUES ('396', '65', '6');
INSERT INTO `node_roles` VALUES ('397', '66', '6');
INSERT INTO `node_roles` VALUES ('398', '70', '6');
INSERT INTO `node_roles` VALUES ('399', '71', '6');
INSERT INTO `node_roles` VALUES ('400', '34', '5');
INSERT INTO `node_roles` VALUES ('401', '35', '5');
INSERT INTO `node_roles` VALUES ('402', '36', '5');
INSERT INTO `node_roles` VALUES ('403', '81', '5');
INSERT INTO `node_roles` VALUES ('404', '82', '5');
INSERT INTO `node_roles` VALUES ('405', '83', '5');
INSERT INTO `node_roles` VALUES ('406', '84', '5');
INSERT INTO `node_roles` VALUES ('407', '86', '5');
INSERT INTO `node_roles` VALUES ('408', '87', '5');
INSERT INTO `node_roles` VALUES ('409', '88', '5');
INSERT INTO `node_roles` VALUES ('410', '89', '5');
INSERT INTO `node_roles` VALUES ('411', '90', '5');
INSERT INTO `node_roles` VALUES ('412', '91', '5');
INSERT INTO `node_roles` VALUES ('413', '92', '5');
INSERT INTO `node_roles` VALUES ('414', '37', '5');
INSERT INTO `node_roles` VALUES ('415', '38', '5');
INSERT INTO `node_roles` VALUES ('416', '201', '5');
INSERT INTO `node_roles` VALUES ('417', '39', '5');
INSERT INTO `node_roles` VALUES ('418', '93', '5');
INSERT INTO `node_roles` VALUES ('419', '94', '5');
INSERT INTO `node_roles` VALUES ('420', '95', '5');
INSERT INTO `node_roles` VALUES ('421', '96', '5');
INSERT INTO `node_roles` VALUES ('422', '98', '5');
INSERT INTO `node_roles` VALUES ('423', '99', '5');
INSERT INTO `node_roles` VALUES ('424', '100', '5');
INSERT INTO `node_roles` VALUES ('425', '40', '5');
INSERT INTO `node_roles` VALUES ('426', '101', '5');
INSERT INTO `node_roles` VALUES ('427', '102', '5');
INSERT INTO `node_roles` VALUES ('428', '103', '5');
INSERT INTO `node_roles` VALUES ('429', '104', '5');
INSERT INTO `node_roles` VALUES ('430', '106', '5');
INSERT INTO `node_roles` VALUES ('431', '107', '5');
INSERT INTO `node_roles` VALUES ('432', '41', '5');
INSERT INTO `node_roles` VALUES ('433', '108', '5');
INSERT INTO `node_roles` VALUES ('434', '109', '5');
INSERT INTO `node_roles` VALUES ('435', '110', '5');
INSERT INTO `node_roles` VALUES ('436', '111', '5');
INSERT INTO `node_roles` VALUES ('437', '113', '5');
INSERT INTO `node_roles` VALUES ('438', '114', '5');
INSERT INTO `node_roles` VALUES ('439', '42', '5');
INSERT INTO `node_roles` VALUES ('440', '115', '5');
INSERT INTO `node_roles` VALUES ('441', '116', '5');
INSERT INTO `node_roles` VALUES ('442', '117', '5');
INSERT INTO `node_roles` VALUES ('443', '118', '5');
INSERT INTO `node_roles` VALUES ('444', '119', '5');
INSERT INTO `node_roles` VALUES ('445', '121', '5');
INSERT INTO `node_roles` VALUES ('446', '122', '5');
INSERT INTO `node_roles` VALUES ('447', '43', '5');
INSERT INTO `node_roles` VALUES ('448', '147', '5');
INSERT INTO `node_roles` VALUES ('449', '148', '5');
INSERT INTO `node_roles` VALUES ('450', '149', '5');
INSERT INTO `node_roles` VALUES ('451', '150', '5');
INSERT INTO `node_roles` VALUES ('452', '152', '5');
INSERT INTO `node_roles` VALUES ('453', '44', '5');
INSERT INTO `node_roles` VALUES ('454', '177', '5');
INSERT INTO `node_roles` VALUES ('455', '178', '5');
INSERT INTO `node_roles` VALUES ('456', '45', '5');
INSERT INTO `node_roles` VALUES ('457', '153', '5');
INSERT INTO `node_roles` VALUES ('458', '154', '5');
INSERT INTO `node_roles` VALUES ('459', '155', '5');
INSERT INTO `node_roles` VALUES ('460', '46', '5');
INSERT INTO `node_roles` VALUES ('461', '180', '5');
INSERT INTO `node_roles` VALUES ('462', '47', '5');
INSERT INTO `node_roles` VALUES ('463', '175', '5');
INSERT INTO `node_roles` VALUES ('464', '176', '5');
INSERT INTO `node_roles` VALUES ('465', '48', '5');
INSERT INTO `node_roles` VALUES ('466', '168', '5');
INSERT INTO `node_roles` VALUES ('467', '169', '5');
INSERT INTO `node_roles` VALUES ('468', '170', '5');
INSERT INTO `node_roles` VALUES ('469', '172', '5');
INSERT INTO `node_roles` VALUES ('470', '173', '5');
INSERT INTO `node_roles` VALUES ('471', '174', '5');
INSERT INTO `node_roles` VALUES ('472', '49', '5');
INSERT INTO `node_roles` VALUES ('473', '161', '5');
INSERT INTO `node_roles` VALUES ('474', '162', '5');
INSERT INTO `node_roles` VALUES ('475', '163', '5');
INSERT INTO `node_roles` VALUES ('476', '164', '5');
INSERT INTO `node_roles` VALUES ('477', '50', '5');
INSERT INTO `node_roles` VALUES ('478', '123', '5');
INSERT INTO `node_roles` VALUES ('479', '124', '5');
INSERT INTO `node_roles` VALUES ('480', '125', '5');
INSERT INTO `node_roles` VALUES ('481', '126', '5');
INSERT INTO `node_roles` VALUES ('482', '128', '5');
INSERT INTO `node_roles` VALUES ('483', '51', '5');
INSERT INTO `node_roles` VALUES ('484', '141', '5');
INSERT INTO `node_roles` VALUES ('485', '142', '5');
INSERT INTO `node_roles` VALUES ('486', '143', '5');
INSERT INTO `node_roles` VALUES ('487', '144', '5');
INSERT INTO `node_roles` VALUES ('488', '146', '5');
INSERT INTO `node_roles` VALUES ('489', '52', '5');
INSERT INTO `node_roles` VALUES ('490', '129', '5');
INSERT INTO `node_roles` VALUES ('491', '130', '5');
INSERT INTO `node_roles` VALUES ('492', '131', '5');
INSERT INTO `node_roles` VALUES ('493', '132', '5');
INSERT INTO `node_roles` VALUES ('494', '134', '5');
INSERT INTO `node_roles` VALUES ('495', '53', '5');
INSERT INTO `node_roles` VALUES ('496', '135', '5');
INSERT INTO `node_roles` VALUES ('497', '136', '5');
INSERT INTO `node_roles` VALUES ('498', '137', '5');
INSERT INTO `node_roles` VALUES ('499', '138', '5');
INSERT INTO `node_roles` VALUES ('500', '139', '5');
INSERT INTO `node_roles` VALUES ('501', '54', '5');
INSERT INTO `node_roles` VALUES ('502', '188', '5');
INSERT INTO `node_roles` VALUES ('503', '189', '5');
INSERT INTO `node_roles` VALUES ('504', '190', '5');
INSERT INTO `node_roles` VALUES ('505', '55', '5');
INSERT INTO `node_roles` VALUES ('506', '185', '5');
INSERT INTO `node_roles` VALUES ('507', '186', '5');
INSERT INTO `node_roles` VALUES ('508', '187', '5');
INSERT INTO `node_roles` VALUES ('509', '56', '5');
INSERT INTO `node_roles` VALUES ('510', '182', '5');
INSERT INTO `node_roles` VALUES ('511', '183', '5');
INSERT INTO `node_roles` VALUES ('512', '184', '5');
INSERT INTO `node_roles` VALUES ('513', '57', '5');
INSERT INTO `node_roles` VALUES ('514', '191', '5');
INSERT INTO `node_roles` VALUES ('515', '196', '5');
INSERT INTO `node_roles` VALUES ('516', '58', '5');
INSERT INTO `node_roles` VALUES ('517', '192', '5');
INSERT INTO `node_roles` VALUES ('518', '197', '5');
INSERT INTO `node_roles` VALUES ('519', '59', '5');
INSERT INTO `node_roles` VALUES ('520', '193', '5');
INSERT INTO `node_roles` VALUES ('521', '198', '5');
INSERT INTO `node_roles` VALUES ('522', '60', '5');
INSERT INTO `node_roles` VALUES ('523', '194', '5');
INSERT INTO `node_roles` VALUES ('524', '199', '5');
INSERT INTO `node_roles` VALUES ('525', '61', '5');
INSERT INTO `node_roles` VALUES ('526', '195', '5');
INSERT INTO `node_roles` VALUES ('527', '200', '5');
INSERT INTO `node_roles` VALUES ('528', '62', '5');
INSERT INTO `node_roles` VALUES ('529', '157', '5');
INSERT INTO `node_roles` VALUES ('530', '158', '5');
INSERT INTO `node_roles` VALUES ('531', '159', '5');
INSERT INTO `node_roles` VALUES ('532', '160', '5');
INSERT INTO `node_roles` VALUES ('533', '181', '5');
INSERT INTO `node_roles` VALUES ('534', '65', '5');
INSERT INTO `node_roles` VALUES ('535', '66', '5');
INSERT INTO `node_roles` VALUES ('536', '67', '5');
INSERT INTO `node_roles` VALUES ('537', '68', '5');
INSERT INTO `node_roles` VALUES ('538', '70', '5');
INSERT INTO `node_roles` VALUES ('539', '71', '5');
INSERT INTO `node_roles` VALUES ('540', '1', '4');
INSERT INTO `node_roles` VALUES ('541', '2', '4');
INSERT INTO `node_roles` VALUES ('542', '3', '4');
INSERT INTO `node_roles` VALUES ('543', '4', '4');
INSERT INTO `node_roles` VALUES ('544', '5', '4');
INSERT INTO `node_roles` VALUES ('545', '27', '4');
INSERT INTO `node_roles` VALUES ('546', '28', '4');
INSERT INTO `node_roles` VALUES ('547', '6', '4');
INSERT INTO `node_roles` VALUES ('548', '7', '4');
INSERT INTO `node_roles` VALUES ('549', '8', '4');
INSERT INTO `node_roles` VALUES ('550', '9', '4');
INSERT INTO `node_roles` VALUES ('551', '10', '4');
INSERT INTO `node_roles` VALUES ('552', '11', '4');
INSERT INTO `node_roles` VALUES ('553', '12', '4');
INSERT INTO `node_roles` VALUES ('554', '13', '4');
INSERT INTO `node_roles` VALUES ('555', '14', '4');
INSERT INTO `node_roles` VALUES ('556', '15', '4');
INSERT INTO `node_roles` VALUES ('557', '16', '4');
INSERT INTO `node_roles` VALUES ('558', '17', '4');
INSERT INTO `node_roles` VALUES ('559', '18', '4');
INSERT INTO `node_roles` VALUES ('560', '19', '4');
INSERT INTO `node_roles` VALUES ('561', '20', '4');
INSERT INTO `node_roles` VALUES ('562', '21', '4');
INSERT INTO `node_roles` VALUES ('563', '22', '4');
INSERT INTO `node_roles` VALUES ('564', '23', '4');
INSERT INTO `node_roles` VALUES ('565', '24', '4');
INSERT INTO `node_roles` VALUES ('566', '29', '4');
INSERT INTO `node_roles` VALUES ('567', '30', '4');
INSERT INTO `node_roles` VALUES ('568', '31', '4');
INSERT INTO `node_roles` VALUES ('569', '72', '4');
INSERT INTO `node_roles` VALUES ('570', '73', '4');
INSERT INTO `node_roles` VALUES ('571', '74', '4');
INSERT INTO `node_roles` VALUES ('572', '75', '4');
INSERT INTO `node_roles` VALUES ('573', '76', '4');
INSERT INTO `node_roles` VALUES ('574', '34', '4');
INSERT INTO `node_roles` VALUES ('575', '35', '4');
INSERT INTO `node_roles` VALUES ('576', '36', '4');
INSERT INTO `node_roles` VALUES ('577', '81', '4');
INSERT INTO `node_roles` VALUES ('578', '82', '4');
INSERT INTO `node_roles` VALUES ('579', '83', '4');
INSERT INTO `node_roles` VALUES ('580', '84', '4');
INSERT INTO `node_roles` VALUES ('581', '85', '4');
INSERT INTO `node_roles` VALUES ('582', '86', '4');
INSERT INTO `node_roles` VALUES ('583', '87', '4');
INSERT INTO `node_roles` VALUES ('584', '88', '4');
INSERT INTO `node_roles` VALUES ('585', '89', '4');
INSERT INTO `node_roles` VALUES ('586', '90', '4');
INSERT INTO `node_roles` VALUES ('587', '91', '4');
INSERT INTO `node_roles` VALUES ('588', '92', '4');
INSERT INTO `node_roles` VALUES ('589', '37', '4');
INSERT INTO `node_roles` VALUES ('590', '38', '4');
INSERT INTO `node_roles` VALUES ('591', '77', '4');
INSERT INTO `node_roles` VALUES ('592', '78', '4');
INSERT INTO `node_roles` VALUES ('593', '79', '4');
INSERT INTO `node_roles` VALUES ('594', '80', '4');
INSERT INTO `node_roles` VALUES ('595', '201', '4');
INSERT INTO `node_roles` VALUES ('596', '39', '4');
INSERT INTO `node_roles` VALUES ('597', '93', '4');
INSERT INTO `node_roles` VALUES ('598', '94', '4');
INSERT INTO `node_roles` VALUES ('599', '95', '4');
INSERT INTO `node_roles` VALUES ('600', '96', '4');
INSERT INTO `node_roles` VALUES ('601', '97', '4');
INSERT INTO `node_roles` VALUES ('602', '98', '4');
INSERT INTO `node_roles` VALUES ('603', '99', '4');
INSERT INTO `node_roles` VALUES ('604', '100', '4');
INSERT INTO `node_roles` VALUES ('605', '40', '4');
INSERT INTO `node_roles` VALUES ('606', '101', '4');
INSERT INTO `node_roles` VALUES ('607', '102', '4');
INSERT INTO `node_roles` VALUES ('608', '103', '4');
INSERT INTO `node_roles` VALUES ('609', '104', '4');
INSERT INTO `node_roles` VALUES ('610', '105', '4');
INSERT INTO `node_roles` VALUES ('611', '106', '4');
INSERT INTO `node_roles` VALUES ('612', '107', '4');
INSERT INTO `node_roles` VALUES ('613', '41', '4');
INSERT INTO `node_roles` VALUES ('614', '108', '4');
INSERT INTO `node_roles` VALUES ('615', '109', '4');
INSERT INTO `node_roles` VALUES ('616', '110', '4');
INSERT INTO `node_roles` VALUES ('617', '111', '4');
INSERT INTO `node_roles` VALUES ('618', '112', '4');
INSERT INTO `node_roles` VALUES ('619', '113', '4');
INSERT INTO `node_roles` VALUES ('620', '114', '4');
INSERT INTO `node_roles` VALUES ('621', '42', '4');
INSERT INTO `node_roles` VALUES ('622', '115', '4');
INSERT INTO `node_roles` VALUES ('623', '116', '4');
INSERT INTO `node_roles` VALUES ('624', '117', '4');
INSERT INTO `node_roles` VALUES ('625', '118', '4');
INSERT INTO `node_roles` VALUES ('626', '119', '4');
INSERT INTO `node_roles` VALUES ('627', '120', '4');
INSERT INTO `node_roles` VALUES ('628', '121', '4');
INSERT INTO `node_roles` VALUES ('629', '122', '4');
INSERT INTO `node_roles` VALUES ('630', '43', '4');
INSERT INTO `node_roles` VALUES ('631', '147', '4');
INSERT INTO `node_roles` VALUES ('632', '148', '4');
INSERT INTO `node_roles` VALUES ('633', '149', '4');
INSERT INTO `node_roles` VALUES ('634', '150', '4');
INSERT INTO `node_roles` VALUES ('635', '151', '4');
INSERT INTO `node_roles` VALUES ('636', '152', '4');
INSERT INTO `node_roles` VALUES ('637', '44', '4');
INSERT INTO `node_roles` VALUES ('638', '177', '4');
INSERT INTO `node_roles` VALUES ('639', '178', '4');
INSERT INTO `node_roles` VALUES ('640', '179', '4');
INSERT INTO `node_roles` VALUES ('641', '45', '4');
INSERT INTO `node_roles` VALUES ('642', '153', '4');
INSERT INTO `node_roles` VALUES ('643', '154', '4');
INSERT INTO `node_roles` VALUES ('644', '155', '4');
INSERT INTO `node_roles` VALUES ('645', '156', '4');
INSERT INTO `node_roles` VALUES ('646', '46', '4');
INSERT INTO `node_roles` VALUES ('647', '180', '4');
INSERT INTO `node_roles` VALUES ('648', '47', '4');
INSERT INTO `node_roles` VALUES ('649', '175', '4');
INSERT INTO `node_roles` VALUES ('650', '176', '4');
INSERT INTO `node_roles` VALUES ('651', '48', '4');
INSERT INTO `node_roles` VALUES ('652', '168', '4');
INSERT INTO `node_roles` VALUES ('653', '169', '4');
INSERT INTO `node_roles` VALUES ('654', '170', '4');
INSERT INTO `node_roles` VALUES ('655', '171', '4');
INSERT INTO `node_roles` VALUES ('656', '172', '4');
INSERT INTO `node_roles` VALUES ('657', '173', '4');
INSERT INTO `node_roles` VALUES ('658', '174', '4');
INSERT INTO `node_roles` VALUES ('659', '49', '4');
INSERT INTO `node_roles` VALUES ('660', '161', '4');
INSERT INTO `node_roles` VALUES ('661', '162', '4');
INSERT INTO `node_roles` VALUES ('662', '163', '4');
INSERT INTO `node_roles` VALUES ('663', '164', '4');
INSERT INTO `node_roles` VALUES ('664', '165', '4');
INSERT INTO `node_roles` VALUES ('665', '166', '4');
INSERT INTO `node_roles` VALUES ('666', '167', '4');
INSERT INTO `node_roles` VALUES ('667', '50', '4');
INSERT INTO `node_roles` VALUES ('668', '123', '4');
INSERT INTO `node_roles` VALUES ('669', '124', '4');
INSERT INTO `node_roles` VALUES ('670', '125', '4');
INSERT INTO `node_roles` VALUES ('671', '126', '4');
INSERT INTO `node_roles` VALUES ('672', '127', '4');
INSERT INTO `node_roles` VALUES ('673', '128', '4');
INSERT INTO `node_roles` VALUES ('674', '51', '4');
INSERT INTO `node_roles` VALUES ('675', '141', '4');
INSERT INTO `node_roles` VALUES ('676', '142', '4');
INSERT INTO `node_roles` VALUES ('677', '143', '4');
INSERT INTO `node_roles` VALUES ('678', '144', '4');
INSERT INTO `node_roles` VALUES ('679', '145', '4');
INSERT INTO `node_roles` VALUES ('680', '146', '4');
INSERT INTO `node_roles` VALUES ('681', '52', '4');
INSERT INTO `node_roles` VALUES ('682', '129', '4');
INSERT INTO `node_roles` VALUES ('683', '130', '4');
INSERT INTO `node_roles` VALUES ('684', '131', '4');
INSERT INTO `node_roles` VALUES ('685', '132', '4');
INSERT INTO `node_roles` VALUES ('686', '133', '4');
INSERT INTO `node_roles` VALUES ('687', '134', '4');
INSERT INTO `node_roles` VALUES ('688', '53', '4');
INSERT INTO `node_roles` VALUES ('689', '135', '4');
INSERT INTO `node_roles` VALUES ('690', '136', '4');
INSERT INTO `node_roles` VALUES ('691', '137', '4');
INSERT INTO `node_roles` VALUES ('692', '138', '4');
INSERT INTO `node_roles` VALUES ('693', '139', '4');
INSERT INTO `node_roles` VALUES ('694', '140', '4');
INSERT INTO `node_roles` VALUES ('695', '54', '4');
INSERT INTO `node_roles` VALUES ('696', '188', '4');
INSERT INTO `node_roles` VALUES ('697', '189', '4');
INSERT INTO `node_roles` VALUES ('698', '190', '4');
INSERT INTO `node_roles` VALUES ('699', '55', '4');
INSERT INTO `node_roles` VALUES ('700', '185', '4');
INSERT INTO `node_roles` VALUES ('701', '186', '4');
INSERT INTO `node_roles` VALUES ('702', '187', '4');
INSERT INTO `node_roles` VALUES ('703', '56', '4');
INSERT INTO `node_roles` VALUES ('704', '182', '4');
INSERT INTO `node_roles` VALUES ('705', '183', '4');
INSERT INTO `node_roles` VALUES ('706', '184', '4');
INSERT INTO `node_roles` VALUES ('707', '57', '4');
INSERT INTO `node_roles` VALUES ('708', '191', '4');
INSERT INTO `node_roles` VALUES ('709', '196', '4');
INSERT INTO `node_roles` VALUES ('710', '58', '4');
INSERT INTO `node_roles` VALUES ('711', '192', '4');
INSERT INTO `node_roles` VALUES ('712', '197', '4');
INSERT INTO `node_roles` VALUES ('713', '59', '4');
INSERT INTO `node_roles` VALUES ('714', '193', '4');
INSERT INTO `node_roles` VALUES ('715', '198', '4');
INSERT INTO `node_roles` VALUES ('716', '60', '4');
INSERT INTO `node_roles` VALUES ('717', '194', '4');
INSERT INTO `node_roles` VALUES ('718', '199', '4');
INSERT INTO `node_roles` VALUES ('719', '61', '4');
INSERT INTO `node_roles` VALUES ('720', '195', '4');
INSERT INTO `node_roles` VALUES ('721', '200', '4');
INSERT INTO `node_roles` VALUES ('722', '62', '4');
INSERT INTO `node_roles` VALUES ('723', '157', '4');
INSERT INTO `node_roles` VALUES ('724', '158', '4');
INSERT INTO `node_roles` VALUES ('725', '159', '4');
INSERT INTO `node_roles` VALUES ('726', '160', '4');
INSERT INTO `node_roles` VALUES ('727', '181', '4');
INSERT INTO `node_roles` VALUES ('728', '65', '4');
INSERT INTO `node_roles` VALUES ('729', '66', '4');
INSERT INTO `node_roles` VALUES ('730', '67', '4');
INSERT INTO `node_roles` VALUES ('731', '68', '4');
INSERT INTO `node_roles` VALUES ('732', '69', '4');
INSERT INTO `node_roles` VALUES ('733', '70', '4');
INSERT INTO `node_roles` VALUES ('734', '71', '4');
INSERT INTO `node_roles` VALUES ('735', '207', '9');
INSERT INTO `node_roles` VALUES ('736', '204', '8');
INSERT INTO `node_roles` VALUES ('737', '205', '7');
INSERT INTO `node_roles` VALUES ('738', '206', '7');
INSERT INTO `node_roles` VALUES ('739', '202', '6');
INSERT INTO `node_roles` VALUES ('740', '203', '6');
INSERT INTO `node_roles` VALUES ('741', '208', '6');
INSERT INTO `node_roles` VALUES ('742', '209', '6');
INSERT INTO `node_roles` VALUES ('743', '202', '5');
INSERT INTO `node_roles` VALUES ('744', '203', '5');
INSERT INTO `node_roles` VALUES ('745', '204', '5');
INSERT INTO `node_roles` VALUES ('746', '205', '5');
INSERT INTO `node_roles` VALUES ('747', '206', '5');
INSERT INTO `node_roles` VALUES ('749', '208', '5');
INSERT INTO `node_roles` VALUES ('750', '209', '5');
INSERT INTO `node_roles` VALUES ('751', '202', '4');
INSERT INTO `node_roles` VALUES ('752', '203', '4');
INSERT INTO `node_roles` VALUES ('753', '204', '4');
INSERT INTO `node_roles` VALUES ('754', '205', '4');
INSERT INTO `node_roles` VALUES ('755', '206', '4');
INSERT INTO `node_roles` VALUES ('756', '207', '4');
INSERT INTO `node_roles` VALUES ('757', '208', '4');
INSERT INTO `node_roles` VALUES ('758', '209', '4');
INSERT INTO `node_roles` VALUES ('759', '181', '9');
INSERT INTO `node_roles` VALUES ('760', '214', '9');
INSERT INTO `node_roles` VALUES ('761', '215', '9');
INSERT INTO `node_roles` VALUES ('762', '220', '9');
INSERT INTO `node_roles` VALUES ('763', '223', '9');
INSERT INTO `node_roles` VALUES ('764', '181', '8');
INSERT INTO `node_roles` VALUES ('765', '211', '8');
INSERT INTO `node_roles` VALUES ('766', '215', '8');
INSERT INTO `node_roles` VALUES ('767', '217', '8');
INSERT INTO `node_roles` VALUES ('768', '222', '8');
INSERT INTO `node_roles` VALUES ('769', '181', '7');
INSERT INTO `node_roles` VALUES ('770', '212', '7');
INSERT INTO `node_roles` VALUES ('771', '213', '7');
INSERT INTO `node_roles` VALUES ('772', '215', '7');
INSERT INTO `node_roles` VALUES ('773', '218', '7');
INSERT INTO `node_roles` VALUES ('774', '219', '7');
INSERT INTO `node_roles` VALUES ('775', '227', '7');
INSERT INTO `node_roles` VALUES ('776', '210', '6');
INSERT INTO `node_roles` VALUES ('777', '215', '6');
INSERT INTO `node_roles` VALUES ('778', '216', '6');
INSERT INTO `node_roles` VALUES ('780', '224', '6');
INSERT INTO `node_roles` VALUES ('781', '225', '6');
INSERT INTO `node_roles` VALUES ('782', '226', '6');
INSERT INTO `node_roles` VALUES ('783', '227', '6');
INSERT INTO `node_roles` VALUES ('784', '228', '6');
INSERT INTO `node_roles` VALUES ('785', '229', '6');
INSERT INTO `node_roles` VALUES ('786', '230', '6');
INSERT INTO `node_roles` VALUES ('787', '231', '6');
INSERT INTO `node_roles` VALUES ('788', '232', '6');
INSERT INTO `node_roles` VALUES ('789', '210', '5');
INSERT INTO `node_roles` VALUES ('790', '211', '5');
INSERT INTO `node_roles` VALUES ('791', '212', '5');
INSERT INTO `node_roles` VALUES ('792', '213', '5');
INSERT INTO `node_roles` VALUES ('793', '215', '5');
INSERT INTO `node_roles` VALUES ('794', '216', '5');
INSERT INTO `node_roles` VALUES ('795', '217', '5');
INSERT INTO `node_roles` VALUES ('796', '218', '5');
INSERT INTO `node_roles` VALUES ('797', '219', '5');
INSERT INTO `node_roles` VALUES ('799', '222', '5');
INSERT INTO `node_roles` VALUES ('800', '224', '5');
INSERT INTO `node_roles` VALUES ('801', '225', '5');
INSERT INTO `node_roles` VALUES ('802', '226', '5');
INSERT INTO `node_roles` VALUES ('803', '227', '5');
INSERT INTO `node_roles` VALUES ('804', '228', '5');
INSERT INTO `node_roles` VALUES ('805', '229', '5');
INSERT INTO `node_roles` VALUES ('806', '230', '5');
INSERT INTO `node_roles` VALUES ('807', '231', '5');
INSERT INTO `node_roles` VALUES ('808', '232', '5');
INSERT INTO `node_roles` VALUES ('809', '210', '4');
INSERT INTO `node_roles` VALUES ('810', '211', '4');
INSERT INTO `node_roles` VALUES ('811', '212', '4');
INSERT INTO `node_roles` VALUES ('812', '213', '4');
INSERT INTO `node_roles` VALUES ('813', '214', '4');
INSERT INTO `node_roles` VALUES ('814', '215', '4');
INSERT INTO `node_roles` VALUES ('815', '216', '4');
INSERT INTO `node_roles` VALUES ('816', '217', '4');
INSERT INTO `node_roles` VALUES ('817', '218', '4');
INSERT INTO `node_roles` VALUES ('818', '219', '4');
INSERT INTO `node_roles` VALUES ('819', '220', '4');
INSERT INTO `node_roles` VALUES ('820', '221', '4');
INSERT INTO `node_roles` VALUES ('821', '222', '4');
INSERT INTO `node_roles` VALUES ('822', '223', '4');
INSERT INTO `node_roles` VALUES ('823', '224', '4');
INSERT INTO `node_roles` VALUES ('824', '225', '4');
INSERT INTO `node_roles` VALUES ('825', '226', '4');
INSERT INTO `node_roles` VALUES ('826', '227', '4');
INSERT INTO `node_roles` VALUES ('827', '228', '4');
INSERT INTO `node_roles` VALUES ('828', '229', '4');
INSERT INTO `node_roles` VALUES ('829', '230', '4');
INSERT INTO `node_roles` VALUES ('830', '231', '4');
INSERT INTO `node_roles` VALUES ('831', '232', '4');
INSERT INTO `node_roles` VALUES ('832', '233', '5');
INSERT INTO `node_roles` VALUES ('833', '233', '6');
INSERT INTO `node_roles` VALUES ('834', '233', '4');
INSERT INTO `node_roles` VALUES ('835', '222', '7');
INSERT INTO `node_roles` VALUES ('836', '234', '8');
INSERT INTO `node_roles` VALUES ('837', '235', '8');
INSERT INTO `node_roles` VALUES ('838', '238', '8');
INSERT INTO `node_roles` VALUES ('839', '240', '8');
INSERT INTO `node_roles` VALUES ('840', '236', '8');
INSERT INTO `node_roles` VALUES ('841', '237', '8');
INSERT INTO `node_roles` VALUES ('842', '239', '8');
INSERT INTO `node_roles` VALUES ('843', '234', '5');
INSERT INTO `node_roles` VALUES ('844', '235', '5');
INSERT INTO `node_roles` VALUES ('845', '238', '5');
INSERT INTO `node_roles` VALUES ('846', '240', '5');
INSERT INTO `node_roles` VALUES ('847', '236', '5');
INSERT INTO `node_roles` VALUES ('848', '237', '5');
INSERT INTO `node_roles` VALUES ('849', '239', '5');
INSERT INTO `node_roles` VALUES ('850', '234', '4');
INSERT INTO `node_roles` VALUES ('851', '235', '4');
INSERT INTO `node_roles` VALUES ('852', '238', '4');
INSERT INTO `node_roles` VALUES ('853', '240', '4');
INSERT INTO `node_roles` VALUES ('854', '236', '4');
INSERT INTO `node_roles` VALUES ('855', '237', '4');
INSERT INTO `node_roles` VALUES ('856', '239', '4');
INSERT INTO `node_roles` VALUES ('857', '241', '6');
INSERT INTO `node_roles` VALUES ('858', '242', '6');
INSERT INTO `node_roles` VALUES ('859', '243', '6');
INSERT INTO `node_roles` VALUES ('860', '244', '6');
INSERT INTO `node_roles` VALUES ('861', '241', '5');
INSERT INTO `node_roles` VALUES ('862', '242', '5');
INSERT INTO `node_roles` VALUES ('863', '243', '5');
INSERT INTO `node_roles` VALUES ('864', '244', '5');
INSERT INTO `node_roles` VALUES ('865', '241', '4');
INSERT INTO `node_roles` VALUES ('866', '242', '4');
INSERT INTO `node_roles` VALUES ('867', '243', '4');
INSERT INTO `node_roles` VALUES ('868', '244', '4');
INSERT INTO `node_roles` VALUES ('869', '245', '5');
INSERT INTO `node_roles` VALUES ('870', '245', '6');
INSERT INTO `node_roles` VALUES ('871', '245', '4');
INSERT INTO `node_roles` VALUES ('872', '246', '5');
INSERT INTO `node_roles` VALUES ('873', '247', '5');
INSERT INTO `node_roles` VALUES ('874', '246', '4');
INSERT INTO `node_roles` VALUES ('875', '247', '4');
INSERT INTO `node_roles` VALUES ('876', '246', '6');
INSERT INTO `node_roles` VALUES ('877', '247', '6');
INSERT INTO `node_roles` VALUES ('878', '248', '7');
INSERT INTO `node_roles` VALUES ('879', '248', '4');
INSERT INTO `node_roles` VALUES ('880', '248', '8');
INSERT INTO `node_roles` VALUES ('881', '248', '5');
INSERT INTO `node_roles` VALUES ('882', '251', '7');
INSERT INTO `node_roles` VALUES ('883', '254', '7');
INSERT INTO `node_roles` VALUES ('884', '253', '9');
INSERT INTO `node_roles` VALUES ('885', '252', '8');
INSERT INTO `node_roles` VALUES ('886', '249', '6');
INSERT INTO `node_roles` VALUES ('887', '250', '6');
INSERT INTO `node_roles` VALUES ('888', '249', '4');
INSERT INTO `node_roles` VALUES ('889', '250', '4');
INSERT INTO `node_roles` VALUES ('890', '251', '4');
INSERT INTO `node_roles` VALUES ('891', '252', '4');
INSERT INTO `node_roles` VALUES ('892', '253', '4');
INSERT INTO `node_roles` VALUES ('893', '254', '4');
INSERT INTO `node_roles` VALUES ('894', '249', '5');
INSERT INTO `node_roles` VALUES ('895', '250', '5');
INSERT INTO `node_roles` VALUES ('896', '251', '5');
INSERT INTO `node_roles` VALUES ('897', '252', '5');
INSERT INTO `node_roles` VALUES ('899', '254', '5');
INSERT INTO `node_roles` VALUES ('900', '255', '4');
INSERT INTO `node_roles` VALUES ('901', '256', '4');
INSERT INTO `node_roles` VALUES ('902', '255', '10');
INSERT INTO `node_roles` VALUES ('903', '256', '10');
INSERT INTO `node_roles` VALUES ('904', '257', '4');
INSERT INTO `node_roles` VALUES ('905', '257', '5');
INSERT INTO `node_roles` VALUES ('906', '257', '7');
INSERT INTO `node_roles` VALUES ('907', '257', '8');
INSERT INTO `node_roles` VALUES ('908', '258', '4');
INSERT INTO `node_roles` VALUES ('909', '259', '4');
INSERT INTO `node_roles` VALUES ('910', '260', '4');
INSERT INTO `node_roles` VALUES ('911', '261', '4');
INSERT INTO `node_roles` VALUES ('912', '258', '5');
INSERT INTO `node_roles` VALUES ('913', '259', '5');
INSERT INTO `node_roles` VALUES ('914', '260', '5');
INSERT INTO `node_roles` VALUES ('915', '261', '5');
INSERT INTO `node_roles` VALUES ('916', '258', '6');
INSERT INTO `node_roles` VALUES ('917', '259', '6');
INSERT INTO `node_roles` VALUES ('918', '260', '6');
INSERT INTO `node_roles` VALUES ('919', '261', '6');
INSERT INTO `node_roles` VALUES ('920', '262', '4');
INSERT INTO `node_roles` VALUES ('921', '1', '5');
INSERT INTO `node_roles` VALUES ('922', '73', '5');
INSERT INTO `node_roles` VALUES ('923', '262', '5');
INSERT INTO `node_roles` VALUES ('924', '1', '6');
INSERT INTO `node_roles` VALUES ('925', '73', '6');
INSERT INTO `node_roles` VALUES ('926', '262', '6');

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '',
  `name` varchar(100) NOT NULL DEFAULT '',
  `remark` varchar(200) DEFAULT NULL,
  `status` int NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES ('1', 'Admin role', '超级管理员', '超级管理员不受权限控制', '2');
INSERT INTO `role` VALUES ('4', '', '普通管理员', '具备全部权限', '2');
INSERT INTO `role` VALUES ('5', '', '普通用户', '具备普通的增改重启等权限', '2');
INSERT INTO `role` VALUES ('6', '', '只读角色', '只能读取资源信息', '2');
INSERT INTO `role` VALUES ('7', '', '修改角色', '修改资源信息', '2');
INSERT INTO `role` VALUES ('8', '', '添加角色', '新增资源信息', '2');
INSERT INTO `role` VALUES ('9', '', '删除角色', '可以删除相关资源', '2');
INSERT INTO `role` VALUES ('10', '', '权限角色', '权限管理', '2');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `sessionid` varchar(32) DEFAULT '',
  `username` varchar(32) NOT NULL DEFAULT '',
  `password` varchar(32) NOT NULL DEFAULT '',
  `nickname` varchar(32) NOT NULL DEFAULT '',
  `telphone` varchar(12) DEFAULT '',
  `company` varchar(255) DEFAULT '',
  `department` varchar(255) DEFAULT '',
  `email` varchar(32) NOT NULL DEFAULT '',
  `remark` varchar(200) DEFAULT NULL,
  `status` int NOT NULL DEFAULT '1',
  `lastloginip` varchar(255) DEFAULT '',
  `lastlogintime` varchar(255) DEFAULT '',
  `createtime` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('1', 'lig20y6tpbpvkbkgcgn9rvkb9wlg0i8w', 'admin', '21232f297a57a5a743894a0e4a801fc3', '超级管理员', '10001', '广东省xxx公司', '研发中心', 'admin@qq.com', '超级管理员', '1', '10.4.34.15', '2025-06-19T06:34:31', '2015-07-24 10:06:49');

-- ----------------------------
-- Table structure for user_cluster
-- ----------------------------
DROP TABLE IF EXISTS `user_cluster`;
CREATE TABLE `user_cluster` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(32) DEFAULT NULL,
  `cluster_id` varchar(64) DEFAULT NULL,
  `createtime` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniqStr` (`username`,`cluster_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of user_cluster
-- ----------------------------

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of user_roles
-- ----------------------------

-- ----------------------------
-- Table structure for wiki
-- ----------------------------
DROP TABLE IF EXISTS `wiki`;
CREATE TABLE `wiki` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '',
  `xcolumn` varchar(64) DEFAULT 'default' COMMENT '文章分组',
  `sketch` varchar(255) DEFAULT '',
  `content` mediumtext,
  `author` varchar(64) DEFAULT '',
  `status` int DEFAULT '1',
  `authkey` varchar(255) DEFAULT '',
  `createtime` varchar(32) NOT NULL DEFAULT '',
  `updatetime` varchar(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of wiki
-- ----------------------------
INSERT INTO `wiki` VALUES ('1', '安装说明文档', '帮助中心', '', '### 功能特性：\n------------\n>   该软件是一款基于client-go、layui、layuimini、beego开发的kubernetes多集群管理系统，该系统具备将多个IDC、公有云的k8s进行统一进行管理。比kubernetes dashboard的功能更丰富，界面更友好更直观。\n\n#### 具备功能\n                \n- 跨公有云、跨IDC的多k8s集群统一管理平台\n- 具备节点、节电池、命名空间、clusterRoleBinding、clusterRoles、RoleBinding、Roles、serviceAccounts的创建、详情、yaml文件查看、删除等功能。\n- workload方面支持对deployment、statefuleset、dameset、cronjob、job、pod容器组、cdr自定义资源、hpa伸缩的功能创建、yaml查看修改、删除功能。\n- 无状态【deployment]:功能具备yaml在线编辑、yaml下载,在线修改升级策略,在线标签修改,在线重启,镜像更新，查看关联的pod对象、查看关联的service、ingress，创建hpa对deployment进行自动伸缩容，可以在线操作回滚到指定的镜像版本，查看该deployment的相关事件，在线查看日志，ssh终端登录关联的pod。\n- 有状态【statefulset】：功能和deployment的类似，除有少量差异之外，相关功能基本一致。\n- 守护进程【daemonset】的功能主要是对deamonset的相关信息进行分类查看。\n- 任务【job】：具备日志、信息、事件、状态的查询功能。\n- 定时任务【cronjob】：在线通过图形库界面进行创建、yaml文件进行创建、对计划任务的在线修改，在线更改状态等功能。\n- 容器组【pod】：具备ssh登录、日志查看、实时查看pod的内存、cpu使用情况【需k8s环境安装metric-beat】等功能。\n- 扩缩容【hpa】：在线图形化操作，根据pod的cpu、内存使用情况、定义pod的扩缩容。\n- 自定义资源【cdr】:自定义资源的信息查看\n- yaml操作：可以通过在线的各种deployment、service、ingress、cronjob等yaml文件模板来进行资源的创建。\n- 服务【service】:支持通过yaml来进行创建，对service的yaml配置查看和修改、支持图形操作创建。\n- 路由【ingres】：支持通过yaml来进行创建，对service的yaml配置查看和修改、支持图形操作创建，目前只支持nginx-ingress。\n- 配置【configmap】：支持图形化、yaml配置的查看、创建、修改和删除。\n- 保密字典【secret】:支持图形化、yaml配置的查看、创建、修改和删除。\n- 存储声明【pvc】:pvc的yaml查看、创建、修改、删除。\n- 存储卷【pv】:存储卷信息的yaml查看、创建、修改、删除。\n- 存储类【storageclass】:存储类信息的yaml查看、创建、修改、删除。\n- 事件信息：查看当前集群中发生的事件信息。\n- 应用集：按照资源的标签appname=myapp进行划分，将该项目所涉及的资源整合到统一界面便于管理。\n- 权限管理：按照角色进行权限划分：超级管理员、普通管理员、只读等角色，并按照集群进行授权，只有授权的了对应集群权限的用户才能访问该集群的资源。\n- 文档中心：markdown格式的文档编辑器、用于运维文档记录。\n     \n### 功能截图\n---\n\n\n#### 下一阶段计划\n---\n* 新增CICD持续集成、持续部署的功能。\n* 根据用户反馈新增功能、或修复bug。\n* 增加镜像、应用的管理。\n\n### 安装说明：\n---\n#### 第一步：安装mysql\n	安装完mysql以后，将xkube.sql导入到数据库中，然后修改conf/app.conf中的如下配置：\n	db_host = 127.0.0.1  #mysql的IP\n	db_port = 3306	#MySQL的端口\n	db_user = root	#mysql的用户名\n	db_pass = root#123	#mysqld 密码\n	db_name = xkube	#数据库名\n	\n#### 第二步：安装redis\n	安装完redis以后：然后修改conf/app.conf中的如下配置：\n	redisDb = \"192.168.10.171:6379\"	#redis的IP和端口\n	redisPasswd = \"redis#123\"	#redis配置设置的密码\n	SessionProviderConfig = \"192.168.10.171:6379,100,redis#123\"	#将redis的IP、端口、密码进行更改，100这个数字保留即可。\n	\n#### 第三步：启动服务\n	修改完以上配置后：进入xkube解压的目录直接执行：nohub ./xkube & \n	就可以通过http://ip:8001/index 进行访问了。注意:直接访问http://ip:8001/ 会出现404，需要附上/index 这个路径。\n', 'admin', '0', '', '2024-01-15 14:24:00', '2024-01-19 15:09:43');
INSERT INTO `wiki` VALUES ('2', 'test123', 'default', '', '## test12', 'admin', '0', '', '2024-01-16 10:29:55', '2024-04-07 10:47:06');
INSERT INTO `wiki` VALUES ('3', 'xkube使用手册', '帮助中心', '', '[TOC]\n\n\n# 功能特性\n- 跨公有云、跨IDC的多k8s集群统一管理平台\n- 具备节点、节电池、命名空间、clusterRoleBinding、clusterRoles、RoleBinding、Roles、serviceAccounts的创建、详情、yaml文件查看、删除等功能。\n- workload方面支持对deployment、statefuleset、dameset、cronjob、job、pod容器组、cdr自定义资源、hpa伸缩的功能创建、yaml查看修改、删除功能。\n- 无状态【deployment]:功能具备yaml在线编辑、yaml下载,在线修改升级策略,在线标签修改,在线重启,镜像更新，查看关联的pod对象、查看关联的service、ingress，创建hpa对deployment进行自动伸缩容，可以在线操作回滚到指定的镜像版本，查看该deployment的相关事件，在线查看日志，ssh终端登录关联的pod。\n- 有状态【statefulset】：功能和deployment的类似，除有少量差异之外，相关功能基本一致。\n- 守护进程【daemonset】的功能主要是对deamonset的相关信息进行分类查看。\n- 任务【job】：具备日志、信息、事件、状态的查询功能。\n- 定时任务【cronjob】：在线通过图形库界面进行创建、yaml文件进行创建、对计划任务的在线修改，在线更改状态等功能。\n- 容器组【pod】：具备ssh登录、日志查看、实时查看pod的内存、cpu使用情况【需k8s环境安装metric-beat】等功能。\n- 扩缩容【hpa】：在线图形化操作，根据pod的cpu、内存使用情况、定义pod的扩缩容。\n- 自定义资源【cdr】:自定义资源的信息查看\n- yaml操作：可以通过在线的各种deployment、service、ingress、cronjob等yaml文件模板来进行资源的创建。\n- 服务【service】:支持通过yaml来进行创建，对service的yaml配置查看和修改、支持图形操作创建。\n- 路由【ingres】：支持通过yaml来进行创建，对service的yaml配置查看和修改、支持图形操作创建，目前只支持nginx-ingress。\n- 配置【configmap】：支持图形化、yaml配置的查看、创建、修改和删除。\n- 保密字典【secret】:支持图形化、yaml配置的查看、创建、修改和删除。\n- 存储声明【pvc】:pvc的yaml查看、创建、修改、删除。\n- 存储卷【pv】:存储卷信息的yaml查看、创建、修改、删除。\n- 存储类【storageclass】:存储类信息的yaml查看、创建、修改、删除。\n- 事件信息：查看当前集群中发生的事件信息。\n- 应用集：按照资源的标签appname=myapp进行划分，将该项目所涉及的资源整合到统一界面便于管理。\n- 权限管理：按照角色进行权限划分：超级管理员、普通管理员、只读等角色，并按照集群进行授权，只有授权的了对应集群权限的用户才能访问该集群的资源。\n- 文档中心：markdown格式的文档编辑器、用于运维文档记录。\n- CICD:支持对接阿里云流水线，可以试下流水线的运行、信息查看、流程查看、日志查看等功能。\n- aws的eks集群管理：eks的管理需要进行通过aws进行认证才能进行管理，可以通过bearerToken实现对接管理。\n# 快速入手\n### 第一步：安装mysql\n	安装完mysql以后，将xkube.sql导入到数据库中，然后修改conf/app.conf中的如下配置：\n	db_host = 127.0.0.1  #mysql的IP\n	db_port = 3306	#MySQL的端口\n	db_user = root	#mysql的用户名\n	db_pass = root#123	#mysqld 密码\n	db_name = xkube	#数据库名\n	\n### 第二步：安装redis\n	安装完redis以后：然后修改conf/app.conf中的如下配置：\n	redisDb = \"192.168.100.17:6379\"	#redis的IP和端口\n	redisPasswd = \"redis#123\"	#redis配置设置的密码\n	SessionProviderConfig = \"192.168.100.171:6379,100,redis#123\"	#将redis的IP、端口、密码进行更改，100这个数字保留即可。\n	\n### 第三步：启动服务\n	修改完以上配置后：Linux 环境下进入目录直接执行：nohub ./xkube & ，windows环境下执行./xkube.exe即可启动\n	就可以通过http://ip:8001/index 进行访问了。注意:直接访问http://ip:8001/ 会出现404，需要附上/index 这个路径。\n\n### 第四步：前端若有nginx反向代理需增加如下配置：\n    location ~^/xkube/pod/terminal/ws {\n        proxy_pass http://127.0.0.1:8080;\n        proxy_http_version 1.1;\n        proxy_set_header Upgrade $http_upgrade;\n        proxy_set_header Connection \"upgrade\";\n        proxy_set_header Origin \"\";\n    }\n\n### 第五步：后台设置\n    1. 登录后台后，默认用户名和密码都是admin，在k8s列表里进行添加需要管理的集群。\n    2. 当有多个集群时，可以为自己设置一个常用集群，该设置会记录到cookie，下一次登录会继续管理该集群。\n    3. 完成以上两步以后就可以在线管理k8s 了。\n\n### 第六步：后台主面板功能说明\n\n![基本信息](/images/0-1.png)\n\n# 1.k8s管理\n## 1.1.登录后首页板块说明\n1.快捷入口：略\n2.资源统计：该部分统计是当有人访问过某个集群对应的资源列表的同时，就会将对应资源的数量记录到redis，然后在通过接口读取出来，该处资源数量非实时更新，当有资源变动以后，需要访问一次对应的资源列表，此处的数值才会更新。\n3.最新文章：略\n4.集群事件：该处的数据是当用户设置了常用集群以后，读取当前k8s集群中的事件。\n\n## 1.2.添加k8s集群配置\n\n| 字段 | 描述 | 是否必须 |\n| --- | --- | --- |\n| 集群ID | 集群的简短描述，支持的英文数字下划线或横线，例如:xkube-a1 | 是 |\n| k8s版本 | xkube会通过该版本号采用不同的api版本来读取资源，务必正确 | 是 |\n| 集群名称 | 中文描述 | 是 |\n| 内网IP入口 | http服务的内网IP | 否 |\n| 公网入口IP | http服务的公网入口IP | 否 |\n| IDC | 集群所属IDC或公有云的描述，例如：aliyun | 否 |\n| 备注 | 该集群的简短描述 | 否 |\n| 状态 | 记录该集群的状态 | 否 |\n| BearerToken | 此功能主要是针对需要进行额外的认证才能读取api的情况,例如aws的eks，eks的bearerToken创建参考“1.5.aws的eks怎样创建不过期token” | 否 |\n| kubeconfig | 调用k8s api必须的配置文件，通常在/root/.kube/config获取 | 是 |\n\n\n## 1.3.设为常用集群\n> 功能描述：在多个k8s集群的情况下，每个人可能管理的k8s集群会不同，可以通过设置该功能，将信息记录到cookie,下一次登录以后默认就可以管理设置好的集群。如果出现每个k8s都需要经常管理的情况，可以直接通过在资源列表右上角的集群进行无缝切换。\n## 1.4.wiki\n> 功能描述：wiki功能模块采用markedown格式来进行记录，可以将一些资源的维护或信息记录到wiki文档中。\n## 1.5.aws的eks怎样创建不过期token\n\n> 第一步：kubectl -n kube-system create serviceaccount kubeconfig-sa\n\n> 第二步：kubectl create clusterrolebinding kubeconfig-sa-binding --clusterrole=cluster-admin --serviceaccount=kube-system:kubeconfig-sa\n\n> 第三步：vim编写secret.yaml\n\n```\napiVersion: v1\nkind: Secret\nmetadata:\n  name: kubeconfig-sa-token\n  namespace: kube-system\n  annotations:\n    kubernetes.io/service-account.name: kubeconfig-sa\ntype: kubernetes.io/service-account-token\n```\n> 第四步：kubectl apply -f secret.yaml -n kube-system\n\n> 第五步：TOKEN=$(kubectl -n kube-system get secret kubeconfig-sa-token -o jsonpath=\\\\\\\\\\\\\\\'{.data.token}\\\\\\\\\\\\\\\'| base64 --decode)\n\n> 第六步：输出token的值：echo $TOKEN \n\n# 2.集群信息\n## 2.1.节点池\n> 功能描述：该功能主要针对阿里云的 容器服务ACK的功能，ack里有一个节电池，主要是将ack的节电池列表显示出来。\n## 2.2.节点\n### 2.2.1.节点详情\n> 节点详情中包含几个主要功能：1.节点的基本信息 2.内存和cpu的当前使用情况[需要k8s安装metric-server] 3.该节点的容器列表 4.yaml查看、节点排水、调度设置、将该节点移除k8s集群 详情参考如下图：\n![基本信息](/images/2.2.1-1.png)\n![内存cpu监控](/images/2.2.1-2.png)\n![节点容器列表](/images/2.2.1-3.png)\n\n## 2.3.命名空间\n### 2.3.1.添加命名空间\n> 图形界面创建命名空间：\n\n![创建命名空间](/images/2.3.1-1.png)\n> yaml创建命名空间：\n\n```\napiVersion: v1\nkind: Namespace\nmetadata:\n  labels:\n    appname: my-app\n  name: my-namespace\n```\n### 2.3.2.命名空间资源限制\n> 资源限制：该限制是将该命名空间下的资源总的cpu和内存资源限制，设置时注意格式:cpu:1核=1000m,可以设置为1、100m、500m，内存格式：200Mi,1Gi\n\n![资源限制](/images/2.3.2-1.png)\n\n### 2.3.3.命名空间yaml编辑\n> 命名空间yaml文件的修改，在建立以后，一般是修改标签。\n\n![yaml编辑](/images/2.3.3-1.png)\n\n## 2.4.clusterRoleBinding\n> 该功能主要是将集群资源clusterRoleBinding列出来,主要是作为信息的查看。平时运维很少用到\n\n![clusterRoleBinding](/images/2.4-1.png)\n\n## 2.5.clusterRoles\n> 该功能主要是将集群资源clusterRoles列出来,主要是作为信息的查看。平时运维很少用到\n\n![clusterRoles](/images/2.5-1.png)\n\n## 2.6.RoleBinding\n> 该功能主要是将集群资源RoleBinding列出来,主要是作为信息的查看。平时运维很少用到\n\n![RoleBinding](/images/2.6-1.png)\n\n## 2.7.Roles\n> 该功能主要是将集群资源clusterRoleBinding列出来,主要是作为信息的查看。平时在操作权限时会用到\n\n![Roles](/images/2.7-1.png)\n\n## 2.8.ServiceAccount\n> 该功能主要是将集群资源ServiceAccount列出来,主要是作为信息的查看。平时在操作权限时会用到\n\n![ServiceAccount](/images/2.8-1.png)\n\n\n# 3.工作负载\n## 3.1.无状态[deployment]\n### 3.1.1.deployment列表\n> 该列表页具备在顶部搜索信息框内实现关键字、命名空间对deployment的搜索，同时可以在右上角进行k8s集群的切换。表格右上方支持对表格的列筛选、表格导出、打印功能。表格的最后一列有对deployment进行操作的快捷键，详情页/cicd/yaml编辑可以参考下图，其中CICD是采用本地流水线或调用阿里云流水线API来实现持续集成和持续部署[前提需要在本地或阿里云进行配置好流水线]，重启按钮是将该deploy下的pod重启，删除是将该deployment删除掉相当于命令：kubectl delete\n\n![列表页](/images/3.1.1-1.png)\n![详情页](/images/3.1.1-2.png)\n![CICD](/images/3.1.1-3.png)\n![YAML编辑](/images/3.1.1-4.png)\n\n### 3.1.2.界面创建deployment\n> deployment的创建必须填的参数：集群选择、名称、镜像地址默认读取当前设置的常用k8s集群\n1.命名空间：可以为空，为空时命名空间为default\n2.标签：可以为空，可以填写多个key-value值，多个时点击后边的加号按钮增加输入框。\n附加配置部分：\n1.资源限制：当勾选资源限制复选框时，会显示资源限制的输入框,按照格式输入即可,1核=1000m,1GiB=1024Mi,1Gi即为1GB的内存,200Mi=200MB内存\n2.健康检查：当勾选健康检查时显示健康检查的输入框,就绪检查readinessProbe,检测类型：URL检测、端口存活检测、执行命令。路径[path]:检查的URL，端口[port]：探测的端口，延迟探测[initialDelaySeconds]:启动pod后延迟多少秒开始检测,探测频率[periodSeconds]:每隔多少秒检测一次,正常阀值[successThreshold]:检测成功多少次算成功,失败阀值[failureThreshold]:失败多少次算失败，超时[timeoutSeconds]:检测的超时时间\n3.生命周期:启动后执行[postStart]:在pod启动前先执行命令，停止前执行[preStop]:在pod停止前执行命令\n4.环境变量[env]：以key-value形式添加一个或多个变量。\n5.service:在创建deployment的同时创建同名的service服务\n\n![deploy创建](/images/3.1.2-1.png)\n\n### 3.1.3.yaml创建deployment\n> 通过yaml配置创建deployment，自己提供配置或选择已有模板进行修改，参考下图：\n\n![yaml创建](/images/3.1.3-1.png)\n\n### 3.1.4.副本伸缩\n> 通过更改副本的数量，来更改pod的数量,参考如下配置,当deployment有设置hpa时，可能导致pod数量冲突，更改无法生效。\n\n```\nspec:\n  replicas: 2\n```\n![pod副本修改](/images/3.1.4-1.png)\n### 3.1.5.重新部署\n>该功能主要是将服务进行重启,将原有pod销毁，重新创建新的pod.和列表页的重启是同一功能\n命令参考：kubectl rollout restart -n test-namespace deployment test-app\n\n![重新部署](/images/3.1.5-1.png)\n\n### 3.1.6.升级策略\n> 升级策略的调整主要调整如下几个参数：\n\n```\n  minReadySeconds: 10  #最小准备时间\n  strategy:\n    rollingUpdate:\n      maxSurge: 25%%%%%%%%     #超过pod的数量或百分比\n      maxUnavailable: 25%%%%%%%%       #最大不可用pod的数量或百分比\n    type: RollingUpdate         #策略类型：滚动更新RollingUpdate 和替换升级：Recreate\n```\n![策略调整](/images/3.1.6-1.png)\n\n### 3.1.7.标签修改\n> 增删改标签labels。点击加号增加key-value输入框，点击减号删除一对标签。key和value只支持英文字母和数字及/,-,_。其他特殊字符不支持。\n```\nmetadata:\n  annotations:\n    deployment.kubernetes.io/revision: \"1\"\n  creationTimestamp: \"2023-12-19T09:52:56Z\"\n  generation: 2\n  labels:\n    app: test-kk\n    appname: test-kk\n```\n![标签修改](/images/3.1.7-1.png)\n\n### 3.1.8.镜像更新\n> 镜像更新主要是更新参数中的image和name。id不用更改\n```\n    spec:\n      containers:\n      - image: nginx:latest\n        imagePullPolicy: Always\n        name: test-kk\n```\n![镜像升级](/images/3.1.8-1.png)\n\n### 3.1.9.创建HPA\n>hpa配置参考如下,主要作用在于自动对资源进行扩缩容\n\n```\nspec:\n  maxReplicas: 5    #最大pod数量\n  metrics:\n  - resource:\n      name: memory\n      target:\n        averageUtilization: 90  #当内存资源90%%%%%%%%时触发扩容事件\n        type: Utilization\n    type: Resource\n  - resource:\n      name: cpu\n      target:\n        averageUtilization: 90  #当cpu资源达到90%%%%%%%%时触发扩容事件\n        type: Utilization\n    type: Resource\n  minReplicas: 2    #最小pod数量，\n  scaleTargetRef:\n    apiVersion: apps/v1\n    kind: Deployment\n    name: test-kk\n```\n![HPA列表](/images/3.1.9-1.png)\n![创建HPA](/images/3.1.9-2.png)\n\n### 3.1.10.回滚配置\n>通过该点击下图中按钮\"回滚到此版本\"将deployment回滚到之前的版本,yaml详情按钮可以查看历史版本的配置。\n\n```\nspec:\n  progressDeadlineSeconds: 600\n  replicas: 2\n  revisionHistoryLimit: 10  #此参数控制保留配置的历史版本数量\n```\n![历史版本](/images/3.1.10-1.png)\n\n\n## 3.2.有状态[statefulset]\n>有状态的所有管理功能和deployment的管理功能一致，可参考3.1.无状态的帮助文档\n\n![sts列表](/images/3.2.1-1.png)\n\n## 3.3.守护进程集[daemonset]\n>daemonset部分功能会少很多，各个按钮的含义和3.1.无状态的功能点一致。具体功能参考3.1.无状态\n\n![ds列表](/images/3.3.1-1.png)\n\n## 3.4.任务[job]\n> job主要是由cronjob定时任务触发创建的任务，可以通过job详情查看job具体的执行情况，并通过job关联的pod查看执行结果日志以及yaml的配置信息。\n\n![job列表](/images/3.4.1-1.png)\n![job详情](/images/3.4.1-2.png)\n\n## 3.5.定时任务[cronjob]\n### 3.5.1.cronjob列表\n> 列表中的 暂停：true时表示该任务已停止运行，false时会定期执行，计划：表示执行运行任务的时间点，参考Linux的crontab定时任务语法，活跃：当任务在运行中是值为1表示，否则值为0.列表中最后一列的启动和暂停按钮用于启动或停止任务，对应修改“suspend:”的值，修改后的值对应列表的暂停这一列的值。CICD：主要用于通过调用流水线来更新镜像。\n\n![cronjob列表](/images/3.5.1-1.png)\n\n### 3.5.2.界面创建cronjob\n> 必填的几个参数：当前集群、命名空间、名称、镜像地址、定时规则。各项参数对应的yaml配置,参考3.5.4.cronjob的yaml语法解释\n\n![cronjob创建](/images/3.5.2-1.png)\n\n### 3.5.3.定时任务时间设置语法\n> 和Linux crontab的语法相同：* * * * * 分别对应分钟（0-59）、小时（0-23）、日期（1-31）、月份（1-12）、星期（0-7）,\n*：星号代表所有可能的值；\n,：逗号用于指定多个值；\n-：横线用于表示一个范围值；\n/：斜杠用于表示重复的频率。\n示例：\n\\\\\\\\\\\\\\\\* * * * * :每分钟执行一次；\n0 * * * * :每小时执行一次；\n0 0 * * * :每天午夜执行一次；\n0 0 * * 1 :每周一午夜执行一次；\n0 0 1 * * :每月1号午夜执行一次。\n3,15 * * * * :每小时的第3和第15分钟执行一次\n3,15 8-11 * * * :每天的8点到11点的第3和第15分钟执行一次\n3,15 8-11 */2 * * :每隔两天的上午8点到11点的第3和第15分钟执行一次\n3,15 8-11 * * 1 :每周一上午8点到11点的第3和第15分钟执行一次\n30 21 * * * :每晚的21:30执行一次\n45 4 1,10,22 * * :每月1、10、22日的4 : 45执行一次\n10 1 * * 6,0 :每周六、周日的1 : 10执行一次\n0,30 18-23 * * * :每天18 : 00至23 : 00之间每隔30分钟执行一次\n0 23 * * 6 :每星期六的晚上23: 00 pm执行一次\n* */1 * * * :每一小时执行一次\n* 23-7/1 * * * :每天晚上23点到第二天7点之间，每隔一小时执行一次\n\n![cron语法建](/images/3.5.3-1.png)\n\n### 3.5.4.cronjob的yaml语法解释\n```yaml\n#curl用法\napiVersion: batch/v1\nkind: CronJob\nmetadata:\n  name: testapp-task-test   #cronjob名称\n  labels:\n    app: testapp-task-test  #标签\n    appname: myapp          #标签\n  namespace: test-app  #命名空间\nspec:\n  schedule: \"*/1 * * * *\"\n  suspend: true #是否暂停\n  jobTemplate:\n    spec:\n      template:\n        spec:\n          containers:\n          - name: testapp-task-test #containers名称\n            image: xxx-registry-vpc.cn-shenzhen.cr.aliyuncs.com/xxx/curl:latest #镜像地址\n            imagePullPolicy: IfNotPresent   #镜像策略\n            args:   #启动时执行的命令及参数\n            - \"/bin/sh\"\n            - \"-ec\"\n            - \"curl -o /dev/stdout --max-time 300 --connect-timeout 10 -H \\\\\\\\\\\\\\\'host:baidu.com\\\\\\\\\\\\\\\' http://svcName.test-app.svc.cluster.local:8001/test.jsp\"\n          restartPolicy: OnFailure  #重启策略,仅支持Never或OnFailure\n```\n\n## 3.6.容器组[pod]\n### 3.6.1.pod列表\n> pod列表默认读取所有命名空间的pod。pod列表中的值非实时刷新，最新的值需要手动刷新页面。针对pod的监控信息、pod日志查看、pod终端可以在操作这一列点击进入\n\n![pod列表](/images/3.6.1-1.png)\n\n### 3.6.2.查看pod详情\n> pod详情页将pod 的基本信息、状态、存储挂载、事件等信息展示出来，参考下图：\n\n![pod详情](/images/3.6.2-1.png)\n![pod详情](/images/3.6.2-2.png)\n\n### 3.6.3.查看pod日志\n> pod的日志查看命令等同：kubectl log pod xxxx -n namespace查看pod的日志。\n\n![pod日志](/images/3.6.3-1.png)\n\n### 3.6.4.登录容器终端\n> 在pod列表中的终端按钮、或者详情中的终端按钮进入pod终端。pod终端是基于websocket实现，在安装配置时可以参考：快速入手的第四步。\n\n![pod终端](/images/3.6.4-1.png)\n\n### 3.6.5.查看pod的资源使用情况\n> pod的资源查看是通过点击列表，操作这一列中的第二个图标按钮进行查看。每个两秒读取一次pod的内存和cpu使用信息并画曲线图。该功能需要k8s安装了metric-server才行。可以通过kubectl top pod来验证k8s是否有安装metric-server。\n\n![pod资源查看](/images/3.6.5-1.png)\n\n### 3.6.6.pod yaml创建\n> 例子如下：\n```yaml\napiVersion: v1\nkind: Pod\nmetadata:\n  namespace: zx-app #命名空间\n  name: ops-tools   #pod名称\nspec:\n  containers:\n  - name: ops-tools #容器名称\n    image: xxx-registry.cn-shenzhen.cr.aliyuncs.com/public/ops-tools:202401012203   #镜像地址\n    command: [\"sleep\",\"36000\"] #执行命令\n```\n\n\n## 3.7.自定义资源[cdr]\n> 自定义资源列表，主要针对一些非k8s自身资源，一些自主开发、或插件的云原生开发的资源信息打印出来。\n\n![cdr列表](/images/3.7.1-1.png)\n\n## 3.8.自动伸缩[hpa]\n> hpa是用来配置deployment的自动扩缩容的功能、根据cpu和内存的使用阀值自动扩缩pod的数量，具体创建可以参考3.1.9部分\n\n![hpa列表](/images/3.7.1-1.png)\n\n## 3.9.yaml创建资源\n> 该功能对应的命令：kubectl appy -f xxx.yaml的功能，主要是将yaml配置复制到输入框，选择对应的集群进行创建。\n\n![yaml创建](/images/3.9.1-1.png)\n\n# 4.网络\n## 4.1.服务[service]\n### 4.1.1.sevice列表\n> 列表中的ip端口：是将service的IP和映射的端口组合，多个端口的逗号隔开了。内部端点：是selector 的app或k8s-app的值+targetport+protocol来显示的，有多个端口的逗号隔开。外部端点：当type: LoadBalancer 时才会显示的对外提供访问的IP和端口。外部端点这一列默认隐藏了\n\n![svc列表](/images/4.1.1-1.png)\n\n### 4.1.2.service详情\n> service详情 除了将service的基本信息显示以外，还会将关联的deployment、pod、statefulset、关联的事件均显示在这里。参考下图\n\n![svc详情](/images/4.1.2-1.png)\n![svc详情](/images/4.1.2-2.png)\n\n### 4.1.3.界面创建service\n> svc名称：service的名称，建议和关联的deployment或statefulset名称一致，关联:deployment或statefulset的名称，\nsvc类型：通常为ClusterIp,当有负载均衡时为LoadBlancer。NodePort时，访问入口是节点IP和端口。当创建无头服务时需要勾选“实例间服务发现”的复选框。\nsvc端口部分：\n名称：分别是端口名称\n服务端口：对外提供服务的端口\n容器端口：deployment暴露出来的端口，以及选择协议\n如果有多个端口时可以点击后边的+号。建议给service加上appname的标签，便于应用集将信息关联起来。例如：appname:test-app,当应用集是test-app时，可以关联所有标签是：appname:test-app的资源。\n\n![svc创建](/images/4.1.2-1.png)\n\n### 4.1.3.yaml创建service\n\n```yaml\napiVersion: v1\nkind: Service\nmetadata:\n  labels:\n    app: testapp-readonly   #标签\n  name: testapp-readonly    #service名称\n  namespace: test-app       #命名空间\nspec:\n  ports:\n  - name: http          #端口名称\n    port: 8001          #对外提供服务的端口\n    protocol: TCP       #协议\n    targetPort: 8001    #deployment暴漏的端口\n  selector:\n    app: testapp-readonly   #关联的deployment的名称\n  sessionAffinity: None\n  type: ClusterIP           #服务类型\n```\n\n## 4.2.路由[ingress]\n### 4.2.1.ingress列表\n\n![ingress列表](/images/4.2.1-1.png)\n> test1\n\n### 4.2.2.界面创ingress\n> 通过界面创建ingress,目前该功能只支持配置单一域名，支持多个后端backend,支持配置证书，标签。当需要配置多个路径不同后端时，店家+号增加输入框。\n选择正确的集群、命名空间以后，输入ingress的名称和 域名。\n路径：默认是根路径，\n匹配规则：一般选择Prefix前缀匹配。exact:准确匹配，ImplementationSpecific：匹配方法取决于 IngressClass。\n服务：service的名称\n端口：service的端口\n\n![ingress创建](/images/4.2.2-1.png)\n\n### 4.2.3.ingress详情\n> 详情主要展示了ingress的基本信息，路由规则、以及事件信息，参考下图：\n\n![ingress详情](/images/4.2.3-1.png)\n\n### 4.2.4.yaml创建ingress\n> yaml创建可以参考如下配置\n\n```yaml\napiVersion: extensions/v1beta1\nkind: Ingress\nmetadata:\n  annotations:\n    kubernetes.io/ingress.class: nginx\n  labels:\n    app: test-app-ing\n  name: test-app-ing\nspec:\n  rules:\n    - host: test.xxx.com.cn #域名\n      http:\n        paths:\n          - backend:\n              serviceName: test-app-nginx   #service的名称\n              servicePort: 80\n            path: /\n    - host: test2.xxx.com.cn    #域名\n      http:\n        paths:\n          - backend:\n              serviceName: test1-app        #service的名称\n              servicePort: 80\n            path: /test1    #路径，不同路径走不同的后端时，可以添加backend\n          - backend:\n              serviceName: test-all-app     #service的名称\n              servicePort: 80\n            path: /             #路径，一般为根路径\n  tls:\n    - hosts:\n        - test1.xxx.com.cn      #域名1\n        - test2.xxx.com.cn      #域名2\n      secretName: xxx-ssl-cert  #将证书和私钥存到secret,这是secret的名称,需要域名1和域名2都可以使用该证书\n\n```\n\n\n# 5.配置管理\n## 5.1.配置项[configmap]\n### 5.1.1.configmap列表\n>configmap列表参考下图\n\n![configmap列表](/images/5.1.1-1.png)\n\n### 5.1.2.configmap详情\n> 详情页将configmap的基本信息展示出来，此处的键值信息不可修改，修改需要通过yaml编辑修改。\n\n![configmap详情](/images/5.1.2-1.png)\n\n### 5.1.3.界面创configmap\n> 选择集群、命名空间、输入configmap名称，输入键值对即可，多个键值对，可以点击加号增加输入框。也支持多个标签添加。\n\n![configmap创建](/images/5.1.3-1.png)\n\n### 5.1.4.yaml创建configmap\n>示例一 以单一键值对，key-value的形式存储\n\n```yaml\napiVersion: v1\nkind: ConfigMap\ndata:\n  kubernetes-pods: mini\n  sync_podMonitor: \"false\"\n  sync_serviceMonitor: \"false\"\nmetadata:\n  name: coredns\n  namespace: kube-system\n  resourceVersion: \"7141240\"\n  uid: cec6dc70-a28c-4e36-83bb-8fecb2f0a7c4\n```\n> 示例二 以整体配置文件的形式\n\n```yaml\napiVersion: v1\nkind: ConfigMap\ndata:\n  Corefile: |\n    .:53 {\n        errors\n        health {\n           lameduck 5s\n        }\n        ready\n\n        kubernetes cluster.local in-addr.arpa ip6.arpa {\n\n          pods verified\n          fallthrough in-addr.arpa ip6.arpa\n          ttl 30\n        }\n        prometheus :9153\n        forward . /etc/resolv.conf\n        cache 30\n        loop\n        reload\n        loadbalance\n    }\nmetadata:\n  name: coredns\n  namespace: kube-system\n  resourceVersion: \"7141240\"\n  uid: cec6dc70-a28c-4e36-83bb-8fecb2f0a7c4\n```\n\n## 5.2.保密字典[secret]\n### 5.2.1.secret列表\n>secret列表列表参考下图\n\n![secret列表](/images/5.2.1-1.png)\n\n### 5.2.2.secret详情\n> 详情页将secret的基本信息展示出来，此处的键值信息不可修改，修改需要通过yaml编辑修改。\n\n![secret详情](/images/5.2.2-1.png)\n\n### 5.2.3.界面创secret\n> 选择集群、命名空间、输入secret名称,\n类型：Opaquue 通常为此类型，在此类型下可以输入多个键值对即可，也可以输入文件类型的，文件名或文件内容。镜像仓库登录密钥：改类型下输入的信息参考secret详情中的值，直接将认证信息粘贴进去即可。TLS证书，将crt或pem格式的证书文件粘贴到cert内容，私钥key文件粘贴到Key内。\n标签可以通过点击加号添加多个。\n\n![界面创secret](/images/5.2.3-1.png)\n\n### 5.2.4.yaml创建secret\n> Opaque类型的yaml参考\n\n```yaml\napiVersion: v1\nkind: Secret\ndata:\n  ALICLOUD_AK: xxxTFRBdddlNck\n  ALICLOUD_AS: xxxxR2N0SccclBWMG\n  GIT_TOKEN: xxxxaHbbbbbZ\nmetadata:\n  labels:\n    app.kubernetes.io/instance: alertmanager2\n  name: alertmanager2\n  namespace: ops-app\ntype: Opaque\n\n```\n\n> 镜像仓库登录密钥类型的yaml参考\n\n```yaml\napiVersion: v1\nkind: Secret\ndata:\n  .dockerconfigjson: ewoJImF1dGhzIjogewoJCSJ3h2WTJGc1FIQmpZWFYwYnpJd01qQTZaRlZaUUVwTldtNXBRalJ6U21oYSIKCQl9Cgl9Cn0K\nmetadata:\n  labels:\n    app: test-app\n  name: pull-secrets\n  namespace: test-app\ntype: kubernetes.io/dockerconfigjson\n\n```\n\n> TLS证书类型的yaml参考\n\n```yaml\napiVersion: v1\nkind: Secret\ndata:\n  tls.crt: LS0tLS1CRUdJTiBDRycW5rY2tTemlZU1F0amlwSWNKREVIc1hvCjRIQT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==\n  tls.key: LS0tLS1Wk92L1Jxc0tnSzNBQ2lzS1FML3JwSUptcXl5Vys0TWxCa2oxVUE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=\nmetadata:\n  labels:\n    app: xxx-tls-secret\n  name: xxx-tls-secret\n  namespace: argocd\ntype: kubernetes.io/tls\n\n```\n\n# 6.存储\n## 6.1.存储申明[Persistent Volume claim]\n### 6.1.1.pvc列表\n> pvc列表功能参考下图:\n\n![pvc列表](/images/6.1.1-1.png)\n### 6.1.2.pvc详情\n> pvc详情功能参考下图:\n\n![pvc详情](/images/6.1.2-1.png)\n### 6.1.3.pvc的yaml讲解\n\n```yaml\napiVersion: v1\nkind: PersistentVolumeClaim\nmetadata:\n  name: my-pvc      #pvc名称\n  namespace: dev   # 命名空间\nspec:\n  accessModes:  #访问模式\n    - ReadWriteMany #以read-write挂载到多个节点\n  resources:\n    requests:\n      storage: 5Gi  # PVC允许申请的大小\n\n```\n\n\n## 6.2.存储卷[Persistent Volume]\n### 6.2.1.pv列表\n> pv列表功能参考下图:\n![pv列表](/images/6.2.1-1.png)\n\n### 6.2.2.pv详情\n> pv详情功能参考下图:\n\n![pv详情](/images/6.2.2-1.png)\n### 6.2.2.pv的yaml讲解\n\n```yaml\napiVersion: v1\nkind: PersistentVolume\nmetadata:\n  name: nginx-pv    # pv的名称\nspec:\n  accessModes:      # 访问模式\n    - ReadWriteMany # PV以read-write挂载到多个节点\n  capacity:  # 容量\n    storage: 2Gi    # pv可用的大小   \n  nfs:\n    path: /nfs/data/     # NFS的挂载路径\n    server: 10.0.0.16    # NFS服务器地址\n```\n\n\n## 6.3.存储类[storageclass]\n### 6.3.1.存储类列表\n> 存储类列表功能参考下图:\n\n![存储类列表](/images/6.3.1-1.png)\n\n### 6.3.2.存储类详情\n> 存储类详情功能参考下图:\n\n![存储类详情](/images/6.3.2-1.png)\n### 6.3.2.存储类的yaml讲解\n\n\n# 7.运维管理\n## 7.1.事件中心\n### 7.1.1.事件查询\n> 将当前集群的所有事件查询统一到此处，便于对于集群的事件的整体掌握。\n\n![事件查询列表](/images/7.1.1-1.png)\n## 7.2.应用集\n### 7.2.1.应用集添加\n> 应用集是指k8s中资源的集合，将同一个域名下的各种deployment/service/configmap/等资源关联一起便于信息查询。该关联是通过在资源中的标签:appname关联的，相同的appname关联到一个集合下。\n参考如下配置：\n\n![appname标签](/images/7.2.1-1.png)\n\n### 7.2.2.应用维护wiki建立\n> 可以将每个应用名创建一个独立的栏目，该栏目下存放该应用所有的维护文档。\n\n## 7.3.CICD\n### 7.3.1.CICD列表\n>CICD 目前引用阿里云的流水线，将流水线的功能关联到k8s集群，便于进行管理，在创建cicd之前需要在阿里云有开通并配置流水线。然后在阿里云帐号设置 这里添加上阿里云帐号的ak信息，并将流水线的组织ID填写上。\n\n![CICD列表](/images/7.3.1-1.png)\n\n### 7.3.2.阿里云帐号设置\n> 添加界面各字段说明如下：\n> - aliyunId:建议用阿里云的主账号名称，或者用一串能表示阿里云帐号意思的英文或数字的字符串。例如：aliyun001\n> - accesskeyId: 登录阿里云后-点击右上角的头像--》权限与安全--》AccessKey管理--》创建AccessKey,就会生成AccessKey Id 和 AccessKey Secret了。\n> - accesskeySecret: 同上\n\n![添加阿里云帐号](/images/7.3.2-1.png)\n![添加阿里云帐号](/images/7.3.2-2.png)\n\n### 7.3.3.添加阿里云流水线\n> 各添加字段说明如下：\n> - 名称：流水线名称,建议和deployment的名称一致,\n> - 应用名：类似为流水线定义一个分组,通过该分组可以查询该组下所有的流水线。\n> - 当前集群:该流水线当前发布的集群，如过名称和deployment一致，集群和命名空间也选择一致。当存在同一个流水线发布到多个k8s集群时，通过集群名称和命名空间区分。\n> - 命名空间：当前流水线发布的命名空间，建议和deployment的命名空间一致\n> - 阿里云帐号：选择该流水线所在的阿里云帐号。需先在运维-CICD-列表的左上方-阿里云帐号设置这里添加。\n> - 组织ID：流水线里会存在一个企业，选择正确的企业下的流水线。\n> -流水线ID：例如：https://flow.aliyun.com/pipelines/2679891/current中的2679891\n\n![添加流水线](/images/7.3.3-1.png)\n\n### 7.3.4.CICD详情\n\n![流水线详情](/images/7.3.4-1.png)\n\n### 7.3.5.怎样查看阿里云流水线组织ID\n> 组织ID不需要手动获取，在添加阿里云帐号时，可以通过组织ID输入框旁边的按钮点击获取，一方面可以验证ak信息是否正确，另一方面能获取到该帐号下的所有组织名称和组织ID。\n\n![添加阿里云帐号](/images/7.3.5-1.png)\n\n### 7.3.6.怎样查看阿里云流水线ID\n> 在流水线的运行界面和编辑界面时，可以复制浏览器中的URL，例如：https://flow.aliyun.com/pipelines/2908037/current ，其中的2908307就是流水线的ID。\n\n### 7.3.7.怎样在阿里云创建AK\n> 参考下图\n\n![阿里云创建AK](/images/7.3.7-1.png)\n\n### 7.3.8.资源克隆功能说明\n> 资源克隆是将资源中的deployment/statefulset/cronjob/configmap/secret/service资源拷贝一份到统一集群 或不同集群的不同命名空间下。使用场景主要是：应用上线、服务迁移集群、跨集群容灾部署\n\n![资源克隆](/images/7.3.8-1.png)\n\n# 8.权限管理\n## 8.1.管理员\n### 8.1.1.管理员列表及修改\n> 管理员的添加、删除、修改、角色授权、均在此页面内完成。\n![管理员列表](/images/8.1.1-1.png)\n![管理员修改](/images/8.1.1-2.png)\n\n## 8.2.角色列表\n\n![角色列表](/images/8.2.1-1.png)\n\n### 8.2.1.角色添加功能授权\n> 在角色列表点击授权列表、弹出框会弹出该角色赋予的授权的URL路径。点击添加会弹出目录结构框，然后勾选路径并选择授权给角色。取消授权：先复选框勾选路径，然后点击取消授权。默认已对不同角色设置好权限，如果不是太熟悉，建议不要在角色中删除授权。\n参考下图：\n\n![角色添加授权](/images/8.2.1-2.png)\n![角色添加授权2](/images/8.2.1-3.png)\n\n### 8.2.2.角色赋予用户\n> 在角色列表、点击用户列表在弹出框中操作，参考下图：\n\n![用户赋予角色](/images/8.2.2-1.png)\n\n## 8.3.目录分组\n### 8.3.1.目录分组列表\n> 目录分组的功能主要是将不同功能类型的URL进行分组，便于查看。\n\n![目录分组列表](/images/8.3.1-1.png)\n\n## 8.4.目录结构\n### 8.4.1.目录结构的添加\n> 该页面主要是针对有新开发的功能或路径时，需要添加到此处，并给予角色授权。自定义开时才用到。默认已配置好路径，建议不要删除里面的路径。\n\n![目录列表](/images/8.4.1-1.png)\n\n### 8.4.2.目录授权给角色\n> 该功能和8.2.1的功能一致，只不过可以单独在此页面进行角色的功能授权。\n\n## 8.5.集群授权\n### 8.5.1.集群授权列表\n> 默认新建用户没有任何集群的权限，如果需要给用户授予某个集群的访问权限，在此处添加，如果选哟授予所有集群的权限，选择全部集群来授权。某个用户需要授权多个集群时，添加多条记录即可。\n\n![集群授权](/images/8.5.1-1.png)\n\n# 9.文档中心\n> 文档中心主要作用是用来针对服务、应用等编写维护文档。文档中心采用markdown的语法来进行记录。对于比较重要的信息可以采用加密的方式。\n\n## 9.1.栏目列表\n> 栏目是由添加或修改文章时定义，不支持单独添加，当两篇文章有相同的栏目名称时，同属一个栏目。一般栏目主要起到一个分组作用，就是将相同类型的文章划分到一个分组，便于查找。例如：应用集--wiki，一般认为一个应用则为一个栏目，然后该应用所有的维护文档可以放到这个应用名的栏目下。\n\n![栏目列表](/images/9.1.1-1.png)\n\n## 9.2.文章列表\n### 9.2.1文章列表\n![文章列表](/images/9.2.1-1.png)\n\n### 9.2.2.文章添加\n> 添加文章注意事项：\n> 1. 栏目：如果需要将文章放到已存在的栏目下，输入的栏目名称必须和已有的名称一致。否则会产生一个新的栏目。\n> 2. 标题：标题格式为不包含特殊字符的中英文即可。\n> 3. 密码：如果需要给阅读该文章加个密码，可以在此处设置一个阅读密码。如果需要取消该密码，在编辑时将该输入框置空后提交。\n> 4. 文章内容：编辑器采用markedown的格式进行编辑，具体的语法可以参考文章列表-顶部的markdown语法示例。或者参考9.2.5的URL。\n\n![文章添加](/images/9.2.2-1.png)\n\n### 9.2.3.文章修改\n> 1. 一种是没有加密的文章，直接编辑修改,栏目、标题、文章内容都可以修改。如果需要加上米，直接在阅读密码输入框输入密码后点击提交。\n> 2. 另一种是修改加密文章，修改前需要先输入密码，如果密码失败会刷新页面重新提示输入密码。输入正确密码以后，可以修改栏目、标题、文章，密码是以\"******\"的形式显示，表示有密码，如果需要修改密码，直接在输入框输入新的密码以后点击提交。\n\n### 9.2.4.文章加密\n> 添加文章时在阅读密码处设置一个密码即可。由于默认是所有用户都能阅读文章，而某些文章中的内容不希望一部分用户看到，这时候可以加一个密码，只有自己或者部门内的人知道这个密码才能查看和修改。如果不希望用户查看文章，还可以通过给该用户设定特定的角色，并在角色权限赋予中不添加文章查看权限即可。\n> #### 如果不知道密码怎么办？\n> 答：可以在数据库中将该文章的authkey的这个字段置空，然后再编辑重新设置一个密码。\n### 9.2.5.markdown语法\n> 参考URL：https://pandao.github.io/editor.md/', 'admin', '0', '', '2024-04-18 12:32:03', '2024-04-24 10:56:17');

-- ----------------------------
-- Table structure for xkb_appname
-- ----------------------------
DROP TABLE IF EXISTS `xkb_appname`;
CREATE TABLE `xkb_appname` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `appname` varchar(64) NOT NULL DEFAULT '' COMMENT '应用名',
  `remarks` varchar(255) DEFAULT '',
  `createtime` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniqapp` (`appname`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of xkb_appname
-- ----------------------------

-- ----------------------------
-- Table structure for xkb_cicd
-- ----------------------------
DROP TABLE IF EXISTS `xkb_cicd`;
CREATE TABLE `xkb_cicd` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `cicd_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `appname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `cluster_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `namespace` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `cicd_type` int NOT NULL DEFAULT '0',
  `status` int DEFAULT '1',
  `remarks` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `last_runtime` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `createtime` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniqCicdName` (`cicd_name`,`cluster_id`,`namespace`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of xkb_cicd
-- ----------------------------

-- ----------------------------
-- Table structure for xkb_cicd_ak
-- ----------------------------
DROP TABLE IF EXISTS `xkb_cicd_ak`;
CREATE TABLE `xkb_cicd_ak` (
  `id` int NOT NULL AUTO_INCREMENT,
  `aliyun_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `accesskey_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `accesskey_secret` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `organization_id` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `remarks` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `createtime` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of xkb_cicd_ak
-- ----------------------------

-- ----------------------------
-- Table structure for xkb_cicd_pipelines
-- ----------------------------
DROP TABLE IF EXISTS `xkb_cicd_pipelines`;
CREATE TABLE `xkb_cicd_pipelines` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `cicd_id` int NOT NULL DEFAULT '0',
  `aliyun_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `organization_id` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `pipeline_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of xkb_cicd_pipelines
-- ----------------------------

-- ----------------------------
-- Table structure for xkb_clone_log
-- ----------------------------
DROP TABLE IF EXISTS `xkb_clone_log`;
CREATE TABLE `xkb_clone_log` (
  `id` int NOT NULL AUTO_INCREMENT,
  `clusterid` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `namespace` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `restype` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `objname` varchar(512) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `target_clusterid` varchar(512) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `target_namespace` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `target_objname` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `status` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `result` varchar(512) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `remarks` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `user` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `createtime` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of xkb_clone_log
-- ----------------------------

-- ----------------------------
-- Table structure for xkb_cluster
-- ----------------------------
DROP TABLE IF EXISTS `xkb_cluster`;
CREATE TABLE `xkb_cluster` (
  `id` int NOT NULL AUTO_INCREMENT,
  `cluster_id` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `cluster_name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `idc_name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `kube_version` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `kube_config` text CHARACTER SET utf8mb3 COLLATE utf8_general_ci,
  `bearer_token` text CHARACTER SET utf8mb3 COLLATE utf8_general_ci COMMENT '用于二次认证例如aws的eks',
  `lan_slbip` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `wan_slbip` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `status` int DEFAULT '1',
  `remarks` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  `createtime` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of xkb_cluster
-- ----------------------------

-- ----------------------------
-- Table structure for xkb_favorite
-- ----------------------------
DROP TABLE IF EXISTS `xkb_favorite`;
CREATE TABLE `xkb_favorite` (
  `id` int NOT NULL AUTO_INCREMENT,
  `fav_type` varchar(64) DEFAULT '',
  `cluster_id` varchar(64) DEFAULT '',
  `name_space` varchar(64) DEFAULT '',
  `fav_name` varchar(64) DEFAULT '',
  `username` varchar(64) DEFAULT '',
  `createtime` varchar(32) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniqFav` (`fav_type`,`cluster_id`,`name_space`,`fav_name`,`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of xkb_favorite
-- ----------------------------
