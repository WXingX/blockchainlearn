package config

import (
	logging "SyncChainEvent/logger"
	"SyncChainEvent/store/gdb"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Log      *logging.LogConf `toml:"log" mapstructure:"log" json:"log"`
	DBCfg    *gdb.Config      `toml:"db" mapstructure:"db" json:"db"`
	ChainCfg *ChainCfg        `toml:"chain_cfg" mapstructure:"chain_cfg" json:"chain_cfg"`
}

type ChainCfg struct {
	Name            string `toml:"name" mapstructure:"name" json:"name"`
	ID              int64  `toml:"id" mapstructure:"id" json:"id"`
	TokenAddress    string `toml:"token_address" mapstructure:"token_address" json:"token_address"`
	DeployerAddress string `toml:"deployer_address" mapstructure:"deployer_address" json:"deployer_address"`
	DelayBlockNum   uint64 `toml:"delay_block_num" mapstructure:"delay_block_num" json:"delay_block_num"`
	RPCUrl          string `toml:"rpc_url" mapstructure:"rpc_url" json:"rpc_url"`
}

// UnmarshalConfig unmarshal conifg file
// @params path: the path of config dir
func UnmarshalConfig(configFilePath string) (*Config, error) {
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("toml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CNFT")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

// UnmarshalCmdConfig unmarshal conifg file
// @params path: the path of config dir
func UnmarshalCmdConfig() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
