CREATE TABLE `dlts`
(
    `id`                 bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
    `draw_num`           varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '期号',
    `draw_time`          varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '开奖时间',
    `equipment_count`    int NULL DEFAULT 0 COMMENT '开奖设备号',
    `draw_pdf_url`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT 'pdf地址',
    `unsort_draw_result` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
    `f1`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区第一个号码',
    `f2`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区第二个号码',
    `f3`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区第三个号码',
    `f4`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区第四个号码',
    `f5`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区第五个号码',
    `b1`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '后区第一个号码',
    `b2`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '后区第二个号码',
    `oe`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '奇偶',
    `hz`                 int NULL DEFAULT 0 COMMENT '前后区和值',
    `ae_hz`              varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT 'hz的ABCDE',
    `qzh`                varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区的前中后的个数，前为1到17，中为18，后为19到35',
    `pool_balance`       decimal(15, 2) NULL DEFAULT 0.00 COMMENT '奖池奖金(元)',
    `total_sale_amount`  decimal(15, 2) NULL DEFAULT 0.00 COMMENT '当期总销售额(元)',
    `stake_count_101`    int NULL DEFAULT 0 COMMENT '一等奖注数',
    `stake_count_102`    int NULL DEFAULT 0 COMMENT '一等奖基本派奖注数',
    `stake_count_201`    int NULL DEFAULT 0 COMMENT '一等奖追加注数',
    `stake_count_202`    int NULL DEFAULT 0 COMMENT '一等奖追加派奖注数',
    `stake_count_301`    int NULL DEFAULT 0 COMMENT '二等奖注数',
    `stake_count_302`    int NULL DEFAULT 0 COMMENT '二等奖基本派奖注数',
    `stake_count_401`    int NULL DEFAULT 0 COMMENT '二等奖追加注数',
    `stake_count_402`    int NULL DEFAULT 0 COMMENT '二等奖追加派奖注数',
    `stake_count_501`    int NULL DEFAULT 0 COMMENT '三等奖注数',
    `stake_count_601`    int NULL DEFAULT 0 COMMENT '四等奖注数',
    `stake_count_701`    int NULL DEFAULT 0 COMMENT '五等奖注数',
    `stake_count_801`    int NULL DEFAULT 0 COMMENT '六等奖注数',
    `stake_count_901`    int NULL DEFAULT 0 COMMENT '七等奖注数',
    `stake_count_1001`   int NULL DEFAULT 0 COMMENT '八等奖注数',
    `stake_count_1101`   int NULL DEFAULT 0 COMMENT '九等奖注数',
    `stake_count_60`     int NULL DEFAULT 0 COMMENT '三等奖追加注数',
    `stake_count_80`     int NULL DEFAULT 0 COMMENT '四等奖追加注数',
    `stake_count_100`    int NULL DEFAULT 0 COMMENT '五等奖追加注数',
    `stake_amount_101`   int NULL DEFAULT 0 COMMENT '一等奖奖金(元)',
    `stake_amount_102`   int NULL DEFAULT 0 COMMENT '一等奖基本派奖奖金(元)',
    `stake_amount_201`   int NULL DEFAULT 0 COMMENT '一等奖追加奖金(元)',
    `stake_amount_202`   int NULL DEFAULT 0 COMMENT '一等奖追加派奖奖金(元)',
    `stake_amount_301`   int NULL DEFAULT 0 COMMENT '二等奖奖金(元)',
    `stake_amount_302`   int NULL DEFAULT 0 COMMENT '二等奖基本派奖奖金(元)',
    `stake_amount_401`   int NULL DEFAULT 0 COMMENT '二等奖追加奖金(元)',
    `stake_amount_402`   int NULL DEFAULT 0 COMMENT '二等奖追加派奖奖金(元)',
    `stake_amount_501`   int NULL DEFAULT 0 COMMENT '三等奖奖金(元)',
    `stake_amount_601`   int NULL DEFAULT 0 COMMENT '四等奖奖金(元)',
    `stake_amount_701`   int NULL DEFAULT 0 COMMENT '五等奖奖金(元)',
    `stake_amount_801`   int NULL DEFAULT 0 COMMENT '六等奖奖金(元)',
    `stake_amount_901`   int NULL DEFAULT 0 COMMENT '七等奖奖金(元)',
    `stake_amount_1001`  int NULL DEFAULT 0 COMMENT '八等奖奖金(元)',
    `stake_amount_1101`  int NULL DEFAULT 0 COMMENT '九等奖奖金(元)',
    `stake_amount_60`    int NULL DEFAULT 0 COMMENT '三等奖追加奖金(元)',
    `stake_amount_80`    int NULL DEFAULT 0 COMMENT '四等奖追加奖金(元)',
    `stake_amount_100`   int NULL DEFAULT 0 COMMENT '五等奖追加奖金(元)',
    `data_src`   int NULL DEFAULT 0 COMMENT '数据来源,0=官网,1=500.com',
    `created_at`         datetime NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic COMMENT='大乐透历史开奖信息';

CREATE TABLE `typs` (
                        `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                        `dlt_id` bigint DEFAULT NULL COMMENT 'dlts表id',
                        `prev_dlt_id` bigint DEFAULT NULL COMMENT 'dlts表id，历史上已经存在hm的id',
                        `typ_1` tinyint DEFAULT NULL COMMENT '类型1:7；6；5；4；',
                        `typ_2` tinyint DEFAULT NULL COMMENT '类型2:70；60；51；42；50；41；32；40；31；22；',
                        `hm` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '号码',
                        `created_at` datetime DEFAULT NULL COMMENT '创建时间',
                        PRIMARY KEY (`id` DESC) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `all_dlts` (
                        `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                        `hm` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '号码',
                        `oe`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '奇偶',
                        `hz`                 int NULL DEFAULT 0 COMMENT '前后区和值',
                        `ae_hz`              varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT 'hz的ABCDE',
                        `qzh`                varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区的前中后的个数，前为1到17，中为18，后为19到35',
                        `created_at` datetime DEFAULT NULL COMMENT '创建时间',
                        PRIMARY KEY (`id` DESC) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT = Dynamic COMMENT='大乐透所有可能的开奖号码';

CREATE TABLE `dlt_monis` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `a` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'A组号码用,隔开',
                            `b` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'B组号码用,隔开',
                            `c` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'C组号码用,隔开',
                            `d` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'D组号码用,隔开',
                            `e` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'E组号码用,隔开',
                            `typ`  varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT 'ABCDE中各自有几个号码',
                            `method`  varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '获取ABCDE中号码个数',
                            `cs`   int NULL DEFAULT 0 COMMENT '历史出现次数',
                            `comb`   int NULL DEFAULT 0 COMMENT '前区号码组合数(只有代码中使用,数据表中可以没有值)',
                            `created_at` datetime DEFAULT NULL COMMENT '创建时间',
                            `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
                            PRIMARY KEY (`id` DESC) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT = Dynamic COMMENT='大乐透模拟数据';


CREATE TABLE `ips` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `ip` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'IP',
                            `created_at` datetime DEFAULT NULL COMMENT '创建时间',
                            PRIMARY KEY (`id` DESC) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT = Dynamic COMMENT='IP表';


CREATE TABLE `ssqs` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `draw_num`           varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '期号',
                            `draw_time`          varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '开奖时间',
                            `week` varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '星期几',
                            `f1`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '红色球第一个号码',
                            `f2`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '红色球第二个号码',
                            `f3`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '红色球第三个号码',
                            `f4`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '红色球第四个号码',
                            `f5`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '红色球第五个号码',
                            `f6`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '红色球第六个号码',
                            `b1`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '蓝色球第一个号码',
                            `oe`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '奇偶',
                            `hz`                 int NULL DEFAULT 0 COMMENT '前后区和值',
                            `ae_hz`              varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT 'hz的ABCDE',
                            `qzh`                varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区的前中后的个数，前为1到16，中为17，后为18到33',
                            `pool_balance`       decimal(15, 2) NULL DEFAULT 0.00 COMMENT '奖池奖金(元)',
                            `total_sale_amount`  decimal(15, 2) NULL DEFAULT 0.00 COMMENT '当期总销售额(元)',
                            `content` varchar(9999) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '一等奖中奖情况',
                            `stake_count_1`    int NULL DEFAULT 0 COMMENT '一等奖注数',
                            `stake_amount_1`   int NULL DEFAULT 0 COMMENT '一等奖奖金(元)',
                            `stake_count_2`    int NULL DEFAULT 0 COMMENT '二等奖注数',
                            `stake_amount_2`   int NULL DEFAULT 0 COMMENT '二等奖奖金(元)',
                            `stake_count_3`    int NULL DEFAULT 0 COMMENT '三等奖注数',
                            `stake_amount_3`   int NULL DEFAULT 0 COMMENT '三等奖奖金(元)',
                            `stake_count_4`    int NULL DEFAULT 0 COMMENT '四等奖注数',
                            `stake_amount_4`   int NULL DEFAULT 0 COMMENT '四等奖奖金(元)',
                            `stake_count_5`    int NULL DEFAULT 0 COMMENT '五等奖注数',
                            `stake_amount_5`   int NULL DEFAULT 0 COMMENT '五等奖奖金(元)',
                            `stake_count_6`    int NULL DEFAULT 0 COMMENT '六等奖注数',
                            `stake_amount_6`   int NULL DEFAULT 0 COMMENT '六等奖奖金(元)',
                            `video_url`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '开奖视频地址',
                            `details_url`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '开奖详情地址',
                            `created_at` datetime DEFAULT NULL COMMENT '创建时间',
                            PRIMARY KEY (`id` DESC) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT = Dynamic COMMENT='双色球历史开奖号码(这里的开始期号为2013001,之前的数据获取不到了)';


CREATE TABLE `all_ssqs` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `hm` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '号码',
                            `oe`                 varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '奇偶',
                            `hz`                 int NULL DEFAULT 0 COMMENT '前后区和值',
                            `ae_hz`              varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT 'hz的ABCDE',
                            `qzh`                varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '前区的前中后的个数，前为1到16，中为17，后为18到33',
                            `created_at` datetime DEFAULT NULL COMMENT '创建时间',
                            PRIMARY KEY (`id` DESC) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT = Dynamic COMMENT='双色球所有可能的开奖号码';


CREATE TABLE `ssq_monis` (
                             `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                             `a` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'A组号码用,隔开',
                             `b` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'B组号码用,隔开',
                             `c` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'C组号码用,隔开',
                             `d` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'D组号码用,隔开',
                             `e` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'E组号码用,隔开',
                             `f` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'F组号码用,隔开',
                             `typ`  varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT 'ABCDEF中各自有几个号码',
                             `method`  varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '获取ABCDEF中号码个数',
                             `cs`   int NULL DEFAULT 0 COMMENT '历史出现次数',
                             `comb`   int NULL DEFAULT 0 COMMENT '前区号码组合数(只有代码中使用,数据表中可以没有值)',
                             `created_at` datetime DEFAULT NULL COMMENT '创建时间',
                             `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
                             PRIMARY KEY (`id` DESC) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT = Dynamic COMMENT='双色球模拟数据';


