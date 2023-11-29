package etc

import "github.com/ethereum/go-ethereum/common"

type (
	Server struct {
		HttpPort      int    `yaml:"http_port"`
		RunMode       string `yaml:"run_mode"`
		FileDir       string `yaml:"file_dir"`
		NodeUrl       string `yaml:"node_url"`
		NodeMode      string `yaml:"node_mode"`
		NodeUser      string `yaml:"node_user"`
		NodePwd       string `yaml:"node_pwd"`
		StartHeight   int64  `yaml:"start_height"`
		NftHeight     int64  `yaml:"nft_height"`
		ForceRollback int    `yaml:"force_rollback"`
		ChainMode     int    `yaml:"chain_mode"` // 0-EVM 1-TRX
	}

	Mysql struct {
		Name     string `yaml:"name"`
		IP       string `yaml:"ip"`
		Port     int    `yaml:"port"`
		Db       string `yaml:"db"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Debug    bool   `yaml:"debug"`
	}

	Redis struct {
		IP       string `yaml:"ip"`
		Port     int    `yaml:"port"`
		DB       []int  `yaml:"db"`
		Password string `yaml:"password"`
	}

	Level struct {
		ChainId uint64 `yaml:"chain_id"`
	}

	Contract struct {
		ArcTokenAddr string `yaml:"arc_token"`
	}

	EthContract struct {
		ArcTokenAddr common.Address
	}

	ServerUrl struct {
		RedPacketShareUrl string `yaml:"red_packet_share_url"`
		JavaHttpDevUrl    string `yaml:"java_http_dev_url"`
		JavaHttpUrl       string `yaml:"java_http_url"`
		JavaWebSocketUrl  string `yaml:"java_web_socket_url"`
		GetTxsListUrl     string `yaml:"get_txs_list_url"`
	}

	config struct {
		Server      *Server      `yaml:"server"`
		Mysql       []*Mysql     `yaml:"mysql"`
		Redis       *Redis       `yaml:"redis"`
		Level       *Level       `yaml:"level"`
		Contract    *Contract    `yaml:"contract"`
		EthContract *EthContract `yaml:"ethContract"`
		ServerUrl   *ServerUrl   `yaml:"server_url"`
	}
)
