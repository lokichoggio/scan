CREATE TABLE `tx_record` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `tx_hash` char(66) NOT NULL COMMENT '交易哈希',
    `method` varchar(64) NOT NULL DEFAULT '' COMMENT '方法',
    `block_number` bigint unsigned NOT NULL COMMENT '区块高度',
    `tx_from` char(42) NOT NULL COMMENT '交易来源地址',
    `tx_to` char(42) NOT NULL DEFAULT '' COMMENT '交易接收地址',
    `tx_value` bigint NOT NULL DEFAULT '0' COMMENT '交易值',
    `tx_fee` bigint NOT NULL DEFAULT '0' COMMENT '交易费用',
    `tx_time` timestamp NOT NULL COMMENT '交易时间',

    /* 1-普通交易，2-创建合约，3-代币交易，10-其他交易 */
    `tx_type` tinyint NOT NULL COMMENT '交易类型',
    /* 当交易类型是 创建合约 */
    `contract_address` char(42) NOT NULL DEFAULT '' COMMENT '合约地址',
    /* 当交易类型是 代币交易 */
    `token_symbol` varchar(64) DEFAULT NULL COMMENT '代币标志',
    `token_decimals` int NOT NULL DEFAULT '0' COMMENT '代币小数位',
    `token_transfer_to` char(42) NOT NULL DEFAULT '' COMMENT '代币接收地址',
    `token_transfer_amount` bigint NOT NULL DEFAULT '0' COMMENT '代币交易值',

    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_tx_hash` (`tx_hash`),
    KEY `idx_tx_from_address` (`tx_from`),
    KEY `idx_tx_to_address` (`tx_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易记录表';
