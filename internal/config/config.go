package config

import (
	"os"

	"ai-knowledge/internal/common"

	"github.com/naoina/toml"
)

// Config 配置文件
type Config struct {
	Debug     bool             `toml:"debug"`
	Address   string           `toml:"address"`
	DB        *DbConfig        `toml:"db"`
	Milvus    *MilvusConfig    `toml:"milvus"`
	LLM       *LLMConfig       `toml:"llm"`
	Embedding *EmbeddingConfig `toml:"embedding"`
}

// 关系型数据库配置
type DbConfig struct {
	Address      string `toml:"address"`        // 数据库连接地址
	Port         int    `toml:"port"`           // 数据库端口
	MaxIdleConns int    `toml:"max_idle_conns"` // 连接池最大连接数
	MaxOpenConns int    `toml:"max_open_conns"` // 默认打开连接数
	User         string `toml:"user"`           // 数据库用户名
	Password     string `toml:"password"`       // 数据库密码
	DbName       string `toml:"db_name"`        // 数据库名
}

// milvus 向量数据库
type MilvusConfig struct {
	Address string `toml:"address"` // 数据库连接地址
	Port    int    `toml:"port"`    // 数据库端口
}

// 模型配置
type LLMConfig struct {
	BaseUrl     string  `toml:"base_url"`
	ApiKey      string  `toml:"api_key"`
	Model       string  `toml:"model"`
	Temperature float64 `toml:"temperature"`
	MaxTokens   int     `toml:"max_tokens"`
	TopP        float64 `toml:"top_p"`
	// TopK        int     `toml:"top_k"`
	FrequencyPenalty float64 `toml:"frequency_penalty"`
	PresencePenalty  float64 `toml:"presence_penalty"`
}

// 向量配置
type EmbeddingConfig struct {
	BaseUrl string `toml:"base_url"`
	ApiKey  string `toml:"api_key"`
	Model   string `toml:"model"`
}

// NewConfig 初始化一个server配置文件对象
func NewConfig(path string) (cfgChan chan *Config, err error) {
	if path == "" {
		path = common.GetRootDir() + "config/cfg.toml"
	}
	cfgChan = make(chan *Config, 0)
	// 读取配置文件
	cfg, err := readConfFile(path)
	if err != nil {
		return
	}
	go watcher(cfgChan, path)
	go func() {
		cfgChan <- cfg
	}()
	return
}

// ReadConfFile 读取配置文件
func readConfFile(path string) (cfg *Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg = new(Config)
	if err := toml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return
}
