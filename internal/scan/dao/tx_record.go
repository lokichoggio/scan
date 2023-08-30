package dao

import (
	"time"

	"scan-eth/pkg/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//type BigInt big.Int
//
//func (b *BigInt) Value() (driver.Value, error) {
//	if b != nil {
//		return (*big.Int)(b).String(), nil
//	}
//	return nil, nil
//}
//
//func (b *BigInt) Scan(value interface{}) error {
//	if value == nil {
//		b = nil
//	}
//
//	switch t := value.(type) {
//	case int64:
//		(*big.Int)(b).SetInt64(value.(int64))
//	default:
//		return fmt.Errorf("could not scan type %T into BigInt", t)
//	}
//
//	return nil
//}

type TxRecord struct {
	gorm.Model

	TxHash      string    `gorm:"column:tx_hash"`
	Method      string    `gorm:"column:method"`
	BlockNumber int64     `gorm:"column:block_number"`
	TxFrom      string    `gorm:"column:tx_from"`
	TxTo        string    `gorm:"column:tx_to"`
	TxValue     int64     `gorm:"column:tx_value"`
	TxFee       int64     `gorm:"column:tx_fee"`
	TxTime      time.Time `gorm:"column:tx_time"`
	// 1-普通交易，2-创建合约，3-代币交易，10-其他交易
	TxType int `gorm:"column:tx_type"`
	// 当交易类型是 创建合约
	ContractAddress string `gorm:"column:contract_address"`
	// 当交易类型是 代币交易
	TokenSymbol         string `gorm:"column:token_symbol"`
	TokenDecimals       uint8  `gorm:"column:token_decimals"`
	TokenTransferTo     string `gorm:"column:token_transfer_to"`
	TokenTransferAmount int64  `gorm:"column:token_transfer_amount"`
}

func (t TxRecord) TableName() string {
	return "tx_record"
}

func BatchCreateTxRecord(db *mysql.DB, records []TxRecord) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&records, 100).Error
}
