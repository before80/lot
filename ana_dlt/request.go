package ana_dlt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/before80/lot/lg"
)

// GenerateRequest 前端发送的请求结构体
type GenerateRequest struct {
	RandomCount         int      `json:"randomCount"`
	Combinations        []string `json:"combinations"`
	Parity              []string `json:"parity"`
	HzMin               int      `json:"hzMin"`
	HzMax               int      `json:"hzMax"`
	RemoveHistory       int      `json:"removeHistory"`
	Danma               []string `json:"danma"`
	FrontInclude        []string `json:"frontInclude"`
	BackInclude         []string `json:"backInclude"`
	CombinationsInclude []string `json:"combinationsInclude"`
	CombinationsExclude []string `json:"combinationsExclude"`
	FrontExclude        []string `json:"frontExclude"`
	BackExclude         []string `json:"backExclude"`
	FrontPosition       []string `json:"frontPosition"`
	NewFourRepeat       int      `json:"newFourRepeat"`
}

// GenerateResponse 返回给前端的响应结构体
type GenerateResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Error   string   `json:"error,omitempty"`
	Numbers []string `json:"numbers,omitempty"`
	Count   int      `json:"count,omitempty"`
}

// 验证组合类型
var validCombinations = map[string]bool{
	"OtherT": true, "T1": true, "T2": true, "T3": true, "T4": true,
	"T5": true, "T6": true, "T7": true, "T8": true, "T9": true,
	"T10": true, "T11": true, "T12": true, "T13": true, "T14": true,
	"T15": true, "T16": true, "T17": true,
}

// 验证奇偶比例
var validParity = map[string]bool{
	"07": true, "16": true, "25": true, "34": true,
	"43": true, "52": true, "61": true, "70": true,
}

// validateRequest 验证请求参数
func validateRequest(req *GenerateRequest) error {
	// 验证随机注数
	if req.RandomCount < 1 || req.RandomCount > 100 {
		return fmt.Errorf("随机注数必须在1-100之间")
	}

	// 验证组合类型
	if len(req.Combinations) == 0 {
		return fmt.Errorf("至少选择一个组合类型")
	}
	for _, combo := range req.Combinations {
		if !validCombinations[combo] {
			return fmt.Errorf("无效的组合类型: %s", combo)
		}
	}

	// 验证奇偶比例
	if len(req.Parity) == 0 {
		return fmt.Errorf("至少选择一个奇偶比例")
	}
	for _, parity := range req.Parity {
		if !validParity[parity] {
			return fmt.Errorf("无效的奇偶比例: %s", parity)
		}
	}

	// 验证和值范围
	if req.HzMin < 18 || req.HzMin > 188 {
		return fmt.Errorf("和值最小值必须在18-188之间")
	}
	if req.HzMax < 18 || req.HzMax > 188 {
		return fmt.Errorf("和值最大值必须在18-188之间")
	}
	if req.HzMin > req.HzMax {
		return fmt.Errorf("和值最小值不能大于最大值")
	}

	// 验证胆码数量
	if len(req.Danma) > 4 {
		return fmt.Errorf("胆码最多只能选择4个")
	}
	for _, danma := range req.Danma {
		if !isValidDanma(danma) {
			return fmt.Errorf("无效的胆码: %s", danma)
		}
	}

	// 验证前区必须包括的号码
	for _, num := range req.FrontInclude {
		if !isValidFrontNumber(num) {
			return fmt.Errorf("无效的前区号码: %s", num)
		}
	}

	// 验证后区必须包括的号码
	for _, num := range req.BackInclude {
		if !isValidBackNumber(num) {
			return fmt.Errorf("无效的后区号码: %s", num)
		}
	}

	// 验证后续必须包含的组合
	for _, pair := range req.CombinationsInclude {
		if !isValidCombinationPair(pair) {
			return fmt.Errorf("无效的组合对: %s", pair)
		}
	}

	// 验证前区必须排除的号码
	for _, num := range req.FrontExclude {
		if !isValidFrontNumber(num) {
			return fmt.Errorf("无效的前区排除号码: %s", num)
		}
	}

	// 验证后区必须排除的号码
	for _, num := range req.BackExclude {
		if !isValidBackNumber(num) {
			return fmt.Errorf("无效的后区排除号码: %s", num)
		}
	}

	// 验证新4重号数
	if req.NewFourRepeat < 0 || req.NewFourRepeat > 300 {
		return fmt.Errorf("新4重号数必须在0-300之间")
	}

	return nil
}

// 验证胆码（01-35）
func isValidDanma(danma string) bool {
	if len(danma) != 2 {
		return false
	}

	// 检查是否在01-35范围内
	validNumbers := map[string]bool{
		"01": true, "02": true, "03": true, "04": true, "05": true,
		"06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true,
		"16": true, "17": true, "18": true, "19": true, "20": true,
		"21": true, "22": true, "23": true, "24": true, "25": true,
		"26": true, "27": true, "28": true, "29": true, "30": true,
		"31": true, "32": true, "33": true, "34": true, "35": true,
	}

	return validNumbers[danma]
}

// 验证前区号码（01-35）
func isValidFrontNumber(num string) bool {
	return isValidDanma(num) // 前区号码范围和胆码一样
}

// 验证后区号码（01-12）
func isValidBackNumber(num string) bool {
	if len(num) != 2 {
		return false
	}

	validNumbers := map[string]bool{
		"01": true, "02": true, "03": true, "04": true, "05": true,
		"06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true,
	}

	return validNumbers[num]
}

// 验证组合对（01,02到11,12）
func isValidCombinationPair(pair string) bool {
	// 验证格式：01,02
	parts := strings.Split(pair, ",")
	if len(parts) != 2 {
		return false
	}

	// 验证两个都是有效的后区号码
	if !isValidBackNumber(parts[0]) || !isValidBackNumber(parts[1]) {
		return false
	}

	// 验证第一个数字小于第二个数字
	return parts[0] < parts[1]
}

// GenerateNumbersHandler 生成号码处理器
func GenerateNumbersHandler(w http.ResponseWriter, r *http.Request) {
	// 记录请求开始时间
	startTime := time.Now()
	lg.InfoToFileAndStdOut(fmt.Sprintf("run 1\n"))
	// 只允许POST方法
	if r.Method != "POST" {
		response := GenerateResponse{
			Success: false,
			Error:   "只支持POST方法",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// 解析JSON请求体
	var req GenerateRequest
	decoder := json.NewDecoder(r.Body)
	//decoder.DisallowUnknownFields() // 不允许未知字段
	lg.InfoToFileAndStdOut(fmt.Sprintf("run 2\n"))
	if err := decoder.Decode(&req); err != nil {
		response := GenerateResponse{
			Success: false,
			Error:   fmt.Sprintf("JSON解析错误: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("run 3\n"))
	// 验证请求数据
	if err := validateRequest(&req); err != nil {
		response := GenerateResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("run 4\n"))
	// 记录请求参数（用于调试）
	lg.InfoToFileAndStdOut(fmt.Sprintf("收到请求: req=%#v", req))
	defer func() {
		if r := recover(); r != nil {
			lg.InfoToFileAndStdOut(fmt.Sprintf("发生错误 %v \n", r))
		}
	}()

	xuHaoSt := &XuHaoSt{
		Tx:               req.Combinations,
		Oes:              req.Parity,
		HzMin:            req.HzMin,
		HzMax:            req.HzMax,
		RemoveHis:        req.RemoveHistory,
		FrontDanMaHms:    req.Danma,
		FrontIncludeHms:  req.FrontInclude,
		BackIncludeHms:   req.BackInclude,
		BackIncludeCombs: req.CombinationsInclude,
		FrontExcludeHms:  req.FrontExclude,
		BackExcludeHms:   req.BackExclude,
		BackExcludeCombs: req.CombinationsExclude,
		QzhSlice:         req.FrontPosition,
		Ch4MustGetCount:  req.NewFourRepeat,
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("run 5\n"))
	// 生成随机号码
	numbers := DltXHaoForWeb1(xuHaoSt, req.RandomCount)
	lg.InfoToFileAndStdOut(fmt.Sprintf("run 6\n"))
	// 计算处理时间
	processingTime := time.Since(startTime).Milliseconds()

	// 返回成功响应
	response := GenerateResponse{
		Success: true,
		Message: fmt.Sprintf("成功生成 %d 注号码，耗时 %dms", len(numbers), processingTime),
		Numbers: numbers,
		Count:   len(numbers),
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("run 7\n"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
