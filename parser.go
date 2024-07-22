package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL = "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"
	address   = "0xYourEthereumAddress"
)

var transferEventSig = []byte("Transfer(address,address,uint256)")

var (
	once   sync.Once
	parser *Parser
)

type ParserImpl interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}

type Transaction struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

type Parser struct {
	logs         chan types.Log
	transactions []Transaction
	sub          ethereum.Subscription
}

func getParserInstance() *Parser {
	once.Do(func() {
		parser = &Parser{
			transactions: make([]Transaction, 0),
			logs:         make(chan types.Log),
		}
	})
	return parser
}

func (p *Parser) GetCurrentBlock() int {
	length := len(p.transactions)
	if length == 0 {
		return -1
	}
	return int(p.transactions[length-1].Value)
}

func (p *Parser) Subscribe(address string) bool {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return false
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(address)},
	}

	p.sub, err = client.SubscribeFilterLogs(context.Background(), query, p.logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
		return false
	}

	go func() {
		for {
			select {
			case err := <-p.sub.Err():
				log.Fatalf("Error: %v", err)
			case vLog := <-p.logs:
				transferEvent := common.BytesToHash(transferEventSig)
				if vLog.Topics[0] == transferEvent {
					fmt.Println("Transfer Event Detected")
					transferEvent := Transaction{}
					err := json.Unmarshal(vLog.Data, &transferEvent)
					if err != nil {
						log.Fatalf("Failed to unmarshal log data: %v", err)
					}

					fmt.Printf("From: %s\n", transferEvent.From.Hex())
					fmt.Printf("To: %s\n", transferEvent.To.Hex())
					fmt.Printf("Value: %s\n", transferEvent.Value.String())
					p.transactions = append(p.transactions, transferEvent)
				}
			}
		}
	}()
	return true
}

func (p *Parser) GetTransactions(address string) []Transaction {
	return p.transactions
}
