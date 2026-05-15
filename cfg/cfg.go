package cfg

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/before80/lot/lg"
	"github.com/spf13/viper"
)

// DefaultConfig 定义整体 JSON 文件的结构
type DefaultConfig struct {
	AeMax77777     [][]string `mapstructure:"ae_max_77777"`
	AeHisMax77777  int        `mapstructure:"ae_his_max_77777"`
	AeMax116666    [][]string `mapstructure:"ae_max_116666"`
	AeHisMax116666 int        `mapstructure:"ae_his_max_116666"`
	AeMax155555    [][]string `mapstructure:"ae_max_155555"`
	AeHisMax155555 int        `mapstructure:"ae_his_max_155555"`
	AeMax194444    [][]string `mapstructure:"ae_max_194444"`
	AeHisMax194444 int        `mapstructure:"ae_his_max_194444"`
	AeMax215432    [][]string `mapstructure:"ae_max_215432"`
	AeHisMax215432 int        `mapstructure:"ae_his_max_215432"`
	AeMax224441    [][]string `mapstructure:"ae_max_224441"`
	AeHisMax224441 int        `mapstructure:"ae_his_max_224441"`
	AeMax224432    [][]string `mapstructure:"ae_max_224432"`
	AeHisMax224432 int        `mapstructure:"ae_his_max_224432"`
	AeMax233333    [][]string `mapstructure:"ae_max_233333"`
	AeHisMax233333 int        `mapstructure:"ae_his_max_233333"`
	AeMax253322    [][]string `mapstructure:"ae_max_253322"`
	AeHisMax253322 int        `mapstructure:"ae_his_max_253322"`
	AeMax272222    [][]string `mapstructure:"ae_max_272222"`
	AeHisMax272222 int        `mapstructure:"ae_his_max_272222"`

	DltTxStr              string `mapstructure:"dlt_tx_str"`
	DltOeStr              string `mapstructure:"dlt_oe_str"`
	DltHzMin              int    `mapstructure:"dlt_hz_min"`
	DltHzMax              int    `mapstructure:"dlt_hz_max"`
	DltRemoveHis          int    `mapstructure:"dlt_remove_his"`
	DltFrontDanMaHmStr    string `mapstructure:"dlt_front_danma_hm_str"`
	DltFrontIncludeHmStr  string `mapstructure:"dlt_front_include_hm_str"`
	DltBackIncludeHmStr   string `mapstructure:"dlt_back_include_hm_str"`
	DltBackIncludeCombStr string `mapstructure:"dlt_back_include_comb_str"`
	DltFrontExcludeHmStr  string `mapstructure:"dlt_front_exclude_hm_str"`
	DltBackExcludeHmStr   string `mapstructure:"dlt_back_exclude_hm_str"`
	DltBackExcludeCombStr string `mapstructure:"dlt_back_exclude_comb_str"`
	DltQzhStr             string `mapstructure:"dlt_qzh_str"`
	DltCh4MustGETCount    int    `mapstructure:"dlt_ch4_must_GET_count"`

	SsqTxStr             string `mapstructure:"ssq_tx_str"`
	SsqOeStr             string `mapstructure:"ssq_oe_str"`
	SsqHzMin             int    `mapstructure:"ssq_hz_min"`
	SsqHzMax             int    `mapstructure:"ssq_hz_max"`
	SsqRemoveHis         int    `mapstructure:"ssq_remove_his"`
	SsqFrontIncludeHmStr string `mapstructure:"ssq_front_include_hm_str"`
	SsqBackIncludeHmStr  string `mapstructure:"ssq_back_include_hm_str"`
	SsqFrontExcludeHmStr string `mapstructure:"ssq_front_exclude_hm_str"`
	SsqBackExcludeHmStr  string `mapstructure:"ssq_back_exclude_hm_str"`
	SsqQzhStr            string `mapstructure:"ssq_qzh_str"`
	SsqCh4MustGETCount   int    `mapstructure:"ssq_ch4_must_GET_count"`

	ChromePath           string   `mapstructure:"chrome_path"`
	Db                   DBConfig `mapstructure:"db"`
	UseGoRodToGetLottery int      `mapstructure:"use_go_rod_to_get_lottery"`
	CloseBrowser         int      `mapstructure:"close_browser"`
	EnvIsLocal           int      `mapstructure:"env_is_local"`
	UseProxy             int      `mapstructure:"use_proxy"`
	ProxyScheme          string   `mapstructure:"proxy_scheme"`
	ProxyHost            string   `mapstructure:"proxy_host"`
	ProxyPort            int      `mapstructure:"proxy_port"`
	ProxyUsername        string   `mapstructure:"proxy_username"`
	ProxyPassword        string   `mapstructure:"proxy_password"`
	UniqueMdFilepath     string   `mapstructure:"unique_md_filepath"`
	BrowserWidth         int      `mapstructure:"browser_width"`
	BrowserHeight        int      `mapstructure:"browser_height"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

var Default DefaultConfig

func init() {
	var err error
	Default, err = getDefaultConfigInfo()
	if err != nil {
		panic(fmt.Sprintf("获取默认配置信息出现错误：%v", err))
	}
}

// GetDefaultConfigInfo 获取默认配置信息
func getDefaultConfigInfo1() (defaultConfig DefaultConfig, err error) {
	viper.SetConfigName("Default") // 配置文件名称（不包含扩展名）
	viper.SetConfigType("toml")    // 配置文件类型
	configPath := "./config"
	_, err1 := os.Stat(configPath)

	if os.IsNotExist(err1) {
		configPath = "../config"
	}

	viper.AddConfigPath(configPath) // 配置文件所在目录

	// 读取配置文件
	if err = viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			lg.ErrorToFile(fmt.Sprintf("配置文件未找到: %v\n", err))
		}
		return
	}
	if err = viper.Unmarshal(&defaultConfig); err != nil {
		lg.ErrorToFile(fmt.Sprintf("解析配置信息到对象时出错: %v\n", err))
		return
	}

	return defaultConfig, err
}

func getDefaultConfigInfo() (defaultConfig DefaultConfig, err error) {
	v := viper.New() // 强烈建议使用新实例，避免全局污染

	configPath := "./config"
	if _, err1 := os.Stat(configPath); os.IsNotExist(err1) {
		configPath = "../config"
	}

	v.AddConfigPath(configPath)
	v.SetConfigType("toml")

	// 1. 读取 Default.toml
	v.SetConfigName("Default")
	if err = v.ReadInConfig(); err != nil {
		return defaultConfig, err
	}

	// 2. 合并 Second.toml（可选）
	v.SetConfigName("xuhao")
	if err = v.MergeInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return defaultConfig, err
		}
		// 不存在可忽略
	}

	// 3. 反序列化
	if err = v.Unmarshal(&defaultConfig); err != nil {
		return defaultConfig, err
	}

	return defaultConfig, nil
}

func Str2StrSliceWithSeparator(str, separator string) (slice []string) {
	stringSlice := strings.Split(str, separator)
	for _, v := range stringSlice {
		str1 := strings.TrimSpace(v)
		if str1 != "" {
			slice = append(slice, str1)
		}
	}
	return
}
