package services

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"scan-eth/internal/common/const"
	"scan-eth/internal/scan/config"
	"scan-eth/internal/scan/dao"
	"scan-eth/pkg/log"
	"scan-eth/pkg/mysql"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type Service struct {
	c        *config.Config
	erc20ABI abi.ABI
	db       *mysql.DB
	client   *ethclient.Client
	wsClient *ethclient.Client
}

func New(c *config.Config) (*Service, error) {
	// erc20 token abi
	f, err := os.Open(c.ERC20ABIFilePath)
	if err != nil {
		return nil, fmt.Errorf("open ERC20ABI file error: %s", err)
	}
	erc20ABI, err := abi.JSON(f)
	if err != nil {
		return nil, fmt.Errorf("abi json error: %s", err)
	}
	_ = f.Close()

	db, err := mysql.InitConn(&c.Db)
	if err != nil {
		return nil, fmt.Errorf("mysql InitConn error: %s", err)
	}

	client, err := ethclient.Dial(c.EthAddr)
	if err != nil {
		return nil, fmt.Errorf("new client error: %s", err)
	}

	wsClient, err := ethclient.Dial(c.EthWsAddr)
	if err != nil {
		return nil, fmt.Errorf("new ws client error: %s", err)
	}

	return &Service{
		c:        c,
		erc20ABI: erc20ABI,
		db:       db,
		client:   client,
		wsClient: wsClient,
	}, nil
}

func (s *Service) Run() error {
	recentBlockNumber, err := s.client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("get recent BlockNumber error: %s", err)
	}

	log.Infof("recent blockNumber: %d", recentBlockNumber)

	go func() {
		if err = s.Scan(int64(recentBlockNumber)); err != nil {
			log.Errorf("Scan error: %s", err)
		}
	}()

	return s.Subscribe()
}

func (s *Service) Stop() error {
	s.client.Close()
	s.wsClient.Close()
	return s.db.Close()
}

func (s *Service) Scan(recentBlockNumber int64) error {
	logger := log.WithFields(zap.String("func", "Scan"))

	for blockNumber := s.c.StartBlockNumber; blockNumber <= recentBlockNumber; blockNumber++ {
		if blockNumber%s.c.ScanIntervalBlock == 0 {
			time.Sleep(time.Second * time.Duration(s.c.ScanIntervalTimeSeconds))
		}

		go func(bn int64) {
			s.Process(bn, logger)
		}(blockNumber)
	}

	return nil
}

func (s *Service) Subscribe() error {
	logger := log.WithFields(zap.String("func", "Subscribe"))

	headers := make(chan *types.Header)
	sub, err := s.wsClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return fmt.Errorf("SubscribeNewHead error: %s", err)
	}

	logger.Info("SubscribeNewHead ok")

	for {
		select {
		case e := <-sub.Err():
			logger.Errorf("Subscribe error: %s", e)
		case header := <-headers:
			go s.Process(header.Number.Int64(), logger)
		}
	}
}

func (s *Service) Process(blockNumber int64, logger *zap.SugaredLogger) {
	logger.Infof("process block number: %d", blockNumber)

	block, err := s.client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		logger.Errorf("blockNumber [%d] BlockByNumber error: %s", blockNumber, err)
		return
	}

	logger.Infof("blockNumber [%d] has [%d] txs", blockNumber, len(block.Transactions()))

	var txRecords []dao.TxRecord

	for _, tx := range block.Transactions() {
		sender, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			logger.Errorf("blockNumber [%d] tx [%s] get sender error: %s", blockNumber, tx.Hash().Hex(), err)
			continue
		}

		from := sender.Hex()
		to := func() string {
			if tx.To() == nil {
				return ""
			}
			return tx.To().Hex()
		}()
		value := tx.Value().Int64()

		receipt, err := s.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			logger.Errorf("blockNumber [%d] tx [%s] TransactionReceipt error: %s", blockNumber, tx.Hash().Hex(), err)
			continue
		}

		txRecord := dao.TxRecord{
			TxHash: tx.Hash().Hex(),
			// todo
			Method:      "",
			BlockNumber: block.Number().Int64(),
			TxFrom:      from,
			TxTo:        to,
			TxValue:     value,
			TxFee:       new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed))).Int64(),
			// todo tx time
			TxTime: tx.Time(),
		}

		switch {
		// ordinary tx
		case to != "" && value != 0:
			// filter target address
			if !s.c.ContainAddress([]string{from, to}) {
				continue
			}

			txRecord.TxType = _const.TxTypeOrdinary
		// create contract
		case to == "" && value == 0:
			// filter target address
			if !s.c.ContainAddress([]string{from, receipt.ContractAddress.Hex()}) {
				continue
			}

			txRecord.TxType = _const.TxTypeCreateContract
			txRecord.ContractAddress = receipt.ContractAddress.Hex()
		// contract tx
		case to != "" && value == 0:
			// erc20 token transfer tx
			if len(tx.Data()) < 4 {
				continue
			}
			if !bytes.Equal(tx.Data()[:4], s.erc20ABI.Methods["transfer"].ID) {
				continue
			}

			// erc20 token info
			token, ok := getToken(tx.To().Hex())
			if !ok {
				token, err = isERC20Token(*tx.To(), s.client)
				if err != nil {
					logger.Debugf("blockNumber [%d] tx [%s] isERC20Token error: %s", blockNumber, tx.Hash().Hex(), err)
					continue
				}
				setToken(tx.To().Hex(), token)
			}

			// token transfer input
			input, err := s.erc20ABI.Methods["transfer"].Inputs.Unpack(tx.Data()[4:])
			if err != nil {
				logger.Errorf("blockNumber [%d] tx [%s] erc20ABI Unpack error: %s", blockNumber, tx.Hash().Hex(), err)
				continue
			}
			if len(input) != 2 {
				logger.Errorf("blockNumber [%d] tx [%s] erc20 token transfer input: %+v", blockNumber, tx.Hash().Hex(), input)
				continue
			}

			transferTo := input[0].(common.Address)
			transferAmount := input[1].(*big.Int)

			// filter target address
			if !s.c.ContainAddress([]string{from, to, transferTo.Hex()}) {
				continue
			}

			txRecord.TxType = _const.TxTypeTokenTransfer
			txRecord.TokenSymbol = token.Symbol
			txRecord.TokenDecimals = token.Decimals
			txRecord.TokenTransferTo = transferTo.Hex()
			txRecord.TokenTransferAmount = transferAmount.Int64()
		default:
			// filter target address
			if !s.c.ContainAddress([]string{from}) {
				continue
			}

			txRecord.TxType = _const.TxTypeOthers
		}

		txRecords = append(txRecords, txRecord)
	}

	if len(txRecords) > 0 {
		if e := dao.BatchCreateTxRecord(s.db, txRecords); e != nil {
			logger.Errorf("blockNumber [%d] BatchCreateTxRecord error: %s, txRecords: %+v", blockNumber, e, txRecords)
			return
		}

		logger.Infof("blockNumber [%d] BatchCreateTxRecord success, txRecords: %+v", blockNumber, txRecords)
	}
}
