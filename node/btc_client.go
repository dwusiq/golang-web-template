package node

import (
	"walletSynV2/etc"

	"github.com/btcsuite/btcd/rpcclient"
)

type BtcNodeClient struct {
	// client    *ethclient.Client
	client *rpcclient.Client
	// chainId   uint64
	host      string
	user      string
	pwd       string
	isTestnet bool
	Result    interface{}
}

func NewBtcClient(host string, user string, pwd string, isDev bool) (*BtcNodeClient, error) {

	// 连接到比特币节点
	connConfig := &rpcclient.ConnConfig{
		Host:         etc.Conf.Server.NodeUrl,  // 比特币节点的RPC地址
		User:         etc.Conf.Server.NodeUser, // 比特币节点的RPC用户名
		Pass:         etc.Conf.Server.NodePwd,  // 比特币节点的RPC密码
		HTTPPostMode: true,                     // 是否使用HTTP POST模式进行通信
		DisableTLS:   true,                     // 是否禁用TLS加密

	}

	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		return nil, err
	}

	c := new(BtcNodeClient)
	c.isTestnet = isDev
	c.host = host
	c.user = user
	c.pwd = pwd
	c.client = client
	return c, nil

	// rpcDial, err := rpc.Dial(host)
	// if err != nil {
	// 	return nil, err
	// }

	// c := new(BtcNodeClient)
	// c.client = ethclient.NewClient(rpcDial)
	// // id, err := c.client.ChainID(context.Background())
	// if err != nil {
	// 	return nil, err
	// }
	// // c.chainId = id.Uint64()
	// c.host = host
	// c.isTestnet = isDev
	// c.rpcClient = rpcDial
	// return c, nil
}

func (c *BtcNodeClient) GetClient() *rpcclient.Client {
	return c.client
}

// func (c *EthNodeClient) GetRpcClient() *rpc.Client {
// 	return c.rpcClient
// }

// func (c *EthNodeClient) GetChainId() uint64 {
// 	return c.chainId
// }

// func (c *EthNodeClient) GetHost() string {
// 	return c.host
// }
