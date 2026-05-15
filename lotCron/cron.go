package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	// 嵌入时区数据到二进制文件中
	_ "time/tzdata"

	"github.com/before80/lot/ana_dlt"
	"github.com/before80/lot/lg"
	"github.com/pelletier/go-toml/v2"
)

// ScheduleConfig 执行时间配置
type ScheduleConfig struct {
	Timezone   string   `toml:"timezone"`    // 时区，默认为 Asia/Shanghai
	TimePoints []string `toml:"time_points"` // 执行时间点列表，格式："星期几:时:分"

	// 高级配置
	CheckInterval int    `toml:"check_interval"` // 检查间隔（秒），默认60
	ConfigFile    string `toml:"-"`              // 配置文件路径
}

// TimePoint 时间点解析结果
type TimePoint struct {
	Weekday time.Weekday // 星期几
	Hour    int          // 时
	Minute  int          // 分
}

type TaskScheduler struct {
	lastRunTime   time.Time
	mu            sync.RWMutex
	location      *time.Location
	timePoints    []TimePoint // 执行时间点配置
	config        *ScheduleConfig
	checkInterval time.Duration
}

func main() {
	// 加载配置
	config, err := loadConfig()
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("加载配置失败: %v", err))
		return
	}

	// 初始化调度器
	scheduler, err := NewTaskScheduler(config)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("初始化调度器失败: %v", err))
		return
	}

	lg.InfoToFile("🚀 定时任务调度器启动")
	lg.InfoToFile(fmt.Sprintf("📅 执行时间（%s）：%s", config.Timezone, formatTimePoints(scheduler.timePoints)))
	lg.InfoToFile(fmt.Sprintf("⏱️ 检查间隔：%d秒", config.CheckInterval))
	lg.InfoToFile(fmt.Sprintf("下次预计执行时间: %v", scheduler.calculateNextRunTime(time.Now().In(scheduler.location))))

	// 显示时区信息
	//scheduler.showTimezoneInfo()

	// 立即检查一次
	scheduler.checkAndRun()

	// 按配置的间隔检查
	ticker := time.NewTicker(scheduler.checkInterval)
	defer ticker.Stop()

	lg.InfoToFile(fmt.Sprintf("⏰ 开始监控: %v，每%d秒检查一次...", time.Now().In(scheduler.location), config.CheckInterval))

	tempNum := 0

	for now := range ticker.C {
		//lg.InfoToFile(fmt.Sprintf("now = %v now.In(scheduler.location)= %v\n", now, now.In(scheduler.location)))
		scheduler.checkAndRunWithTime(now.In(scheduler.location))
		if tempNum == 0 {
			scheduler.printStatus()
		}
		if tempNum < 360 {
			tempNum++
		} else {
			tempNum = 0
		}

		//// 每天打印一次状态
		//if now.In(scheduler.location).Second() == 0 {
		//
		//}
	}
}

// 加载配置文件
func loadConfig() (*ScheduleConfig, error) {
	configFile := "schedule_config.toml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	// 默认配置
	config := &ScheduleConfig{
		ConfigFile:    configFile,
		Timezone:      "Asia/Shanghai",
		TimePoints:    []string{"1:23:58", "3:23:58", "6:23:58"}, // 默认周一、三、六 23:58
		CheckInterval: 60,                                        // 默认60秒检查一次
	}

	// 尝试读取配置文件
	if file, err := os.ReadFile(configFile); err == nil {
		var loadedConfig ScheduleConfig
		if err1 := toml.Unmarshal(file, &loadedConfig); err1 != nil {
			lg.InfoToFile(fmt.Sprintf("TOML配置文件解析失败，使用默认配置: %v", err1))
		} else {
			// 合并配置
			if loadedConfig.Timezone != "" {
				config.Timezone = loadedConfig.Timezone
			}
			if len(loadedConfig.TimePoints) > 0 {
				config.TimePoints = loadedConfig.TimePoints
			}
			if loadedConfig.CheckInterval > 0 {
				config.CheckInterval = loadedConfig.CheckInterval
			}
			lg.InfoToFile(fmt.Sprintf("从配置文件 %s 加载配置成功", configFile))
		}
	} else {
		lg.InfoToFile(fmt.Sprintf("配置文件 %s 不存在，使用默认配置", configFile))
		// 创建默认配置文件
		if err1 := createDefaultConfig(configFile); err1 != nil {
			lg.InfoToFile(fmt.Sprintf("创建默认配置文件失败: %v", err1))
		}
	}

	return config, nil
}

// 创建默认配置文件
func createDefaultConfig(configFile string) error {
	// 直接使用带注释的字符串，更清晰
	configContent := `# 定时任务调度器配置
# 时区设置（支持所有标准时区，如 Asia/Shanghai, America/New_York 等）
timezone = "Asia/Shanghai"

# 检查间隔（秒），建议设置为60（每分钟检查一次）
check_interval = 60

# 执行时间点配置
# 格式：星期几:时:分
# 星期几：0=周日，1=周一，2=周二，3=周三，4=周四，5=周五，6=周六
# 示例：
#   "1:23:58"  # 周一 23:58
#   "3:23:58"  # 周三 23:58
#   "6:23:58"  # 周六 23:58
time_points = [
    "1:23:58",
    "3:23:58",
    "6:23:58",
]

# 更多示例：
# time_points = [
#     "1:09:30",   # 周一 09:30
#     "1:14:00",   # 周一 14:00
#     "3:12:00",   # 周三 12:00
#     "6:08:00",   # 周六 08:00
#     "6:16:30",   # 周六 16:30
# ]
`

	return os.WriteFile(configFile, []byte(configContent), 0644)
}

// 保存配置到文件
func saveConfig(config *ScheduleConfig) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(config.ConfigFile, data, 0644)
}

func NewTaskScheduler(config *ScheduleConfig) (*TaskScheduler, error) {
	// 解析时区
	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		lg.InfoToFile(fmt.Sprintf("⚠️ 无法加载时区 %s，尝试其他中国时区: %v", config.Timezone, err))

		// 尝试其他中国时区名称
		locationsToTry := []string{
			"Asia/Shanghai",
			"Asia/Chongqing",
			"Asia/Harbin",
			"Asia/Urumqi",
			"Asia/Taipei",
			"Asia/Hong_Kong",
			"Asia/Macau",
		}

		for _, locName := range locationsToTry {
			location, err = time.LoadLocation(locName)
			if err == nil {
				lg.InfoToFile(fmt.Sprintf("✅ 成功加载时区: %s", locName))
				config.Timezone = locName // 更新配置中的时区
				break
			}
		}

		if err != nil {
			lg.InfoToFile("⚠️ 无法加载任何中国时区，使用 UTC+8")
			// 使用固定偏移（UTC+8）
			location = time.FixedZone("CST", 8*60*60)
		}
	} else {
		lg.InfoToFile(fmt.Sprintf("✅ 成功加载时区: %s", config.Timezone))
	}

	// 解析时间点配置
	timePoints, err := parseTimePoints(config.TimePoints)
	if err != nil {
		return nil, fmt.Errorf("解析时间点配置失败: %v", err)
	}

	// 设置检查间隔
	checkInterval := time.Duration(config.CheckInterval) * time.Second
	if checkInterval < time.Second {
		checkInterval = time.Second // 最小1秒
	}

	scheduler := &TaskScheduler{
		location:      location,
		timePoints:    timePoints,
		config:        config,
		checkInterval: checkInterval,
	}

	// 显示解析的时间点
	lg.InfoToFile(fmt.Sprintf("✅ 解析到 %d 个执行时间点", len(timePoints)))
	for i, tp := range timePoints {
		lg.InfoToFile(fmt.Sprintf("  时间点 %d: 星期%s %02d:%02d",
			i+1,
			formatWeekday(tp.Weekday),
			tp.Hour,
			tp.Minute))
	}

	return scheduler, nil
}

// 解析时间点配置
func parseTimePoints(timePoints []string) ([]TimePoint, error) {
	result := make([]TimePoint, 0, len(timePoints))

	for _, tpStr := range timePoints {
		parts := strings.Split(tpStr, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("时间点格式错误: %s，应为 星期几:时:分", tpStr)
		}

		// 解析星期几
		weekday, err := strconv.Atoi(parts[0])
		if err != nil || weekday < 0 || weekday > 6 {
			return nil, fmt.Errorf("星期几必须为0-6的数字: %s", parts[0])
		}

		// 解析小时
		hour, err := strconv.Atoi(parts[1])
		if err != nil || hour < 0 || hour > 23 {
			return nil, fmt.Errorf("小时必须为0-23的数字: %s", parts[1])
		}

		// 解析分钟
		minute, err := strconv.Atoi(parts[2])
		if err != nil || minute < 0 || minute > 59 {
			return nil, fmt.Errorf("分钟必须为0-59的数字: %s", parts[2])
		}

		result = append(result, TimePoint{
			Weekday: time.Weekday(weekday), // 转换为 time.Weekday 类型
			Hour:    hour,
			Minute:  minute,
		})
	}

	return result, nil
}

func (s *TaskScheduler) checkAndRun() {
	s.checkAndRunWithTime(time.Now().In(s.location))
}

func (s *TaskScheduler) checkAndRunWithTime(beijingTime time.Time) {
	if s.shouldRunTask(beijingTime) && s.canRun(beijingTime) {
		lg.InfoToFile(fmt.Sprintf("🎯 触发执行（%s）: %s", s.config.Timezone, beijingTime.Format("2006-01-02 15:04:05")))
		lg.InfoToFile(fmt.Sprintf("   🖥️  系统时间: %s", time.Now().Format("2006-01-02 15:04:05")))
		s.runTask()
		s.recordRun(beijingTime)
	}
}

func (s *TaskScheduler) shouldRunTask(t time.Time) bool {
	// 检查所有配置的时间点
	for _, tp := range s.timePoints {
		if s.matchesTimePoint(t, tp) {
			return true
		}
	}
	return false
}

// 检查时间是否匹配时间点
func (s *TaskScheduler) matchesTimePoint(t time.Time, tp TimePoint) bool {
	return t.Weekday() == tp.Weekday &&
		t.Hour() == tp.Hour &&
		t.Minute() == tp.Minute
}

func (s *TaskScheduler) canRun(t time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 如果上次执行是同一分钟，则不再执行（防止重复执行）
	if !s.lastRunTime.IsZero() &&
		s.lastRunTime.Year() == t.Year() &&
		s.lastRunTime.YearDay() == t.YearDay() &&
		s.lastRunTime.Hour() == t.Hour() &&
		s.lastRunTime.Minute() == t.Minute() {
		lg.InfoToFile(fmt.Sprintf("⏭️  当前时间点已执行过，跳过: %s", t.Format("2006-01-02 15:04")))
		return false
	}

	return true
}

func (s *TaskScheduler) recordRun(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastRunTime = t
	lg.InfoToFile(fmt.Sprintf("✅ 记录执行时间: %s", t.Format("2006-01-02 15:04:05")))
}

func (s *TaskScheduler) runTask() {
	lg.InfoToFile(">>> 开始执行定时任务")

	// 获取执行时的北京时间
	beijingTime := time.Now().In(s.location)

	// 执行业务逻辑
	s.executeBusinessLogic(beijingTime)

	lg.InfoToFile(">>> 定时任务执行完成")
}

// executeBusinessLogic 执行任务
func (s *TaskScheduler) executeBusinessLogic(execTime time.Time) {
	// 这里替换为你的实际业务代码
	executeRes1 := ana_dlt.GetDltKj()
	if executeRes1 {
		//// 1. 生成模拟数据
		//ana_dlt.BatchMoni(3)
		//
		//// 2. 更新monis表
		//ana_dlt.UpdateMonisTable()
		//
		//// 3. 生成 Excel文件
		//ana_dlt.DltDataToExcel(false)

		// 或者直接
		ana_dlt.DltDataToExcel(true)

		lg.InfoToFile(">>> 定时任务执行成功")
	} else {
		lg.ErrorToFile(">>> 定时任务执行Fail")
	}
}

func (s *TaskScheduler) printStatus() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	currentTime := time.Now().In(s.location)

	// 显示状态信息
	printSeparator("-", 50)
	lg.InfoToFile(fmt.Sprintf("📈 调度器状态报告（%s）:", s.config.Timezone))
	lg.InfoToFile(fmt.Sprintf("   当前时间: %s", currentTime.Format("2006-01-02 15:04:05")))

	if !s.lastRunTime.IsZero() {
		lg.InfoToFile(fmt.Sprintf("   上次执行时间: %s", s.lastRunTime.Format("2006-01-02 15:04:05")))
		lg.InfoToFile(fmt.Sprintf("   距今: %v", currentTime.Sub(s.lastRunTime)))
	} else {
		lg.InfoToFile("   上次执行时间: 从未执行")
	}

	// 显示下次预计执行时间
	nextRun := s.calculateNextRunTime(currentTime)
	lg.InfoToFile(fmt.Sprintf("   下次预计执行: %s", nextRun.Format("2006-01-02 15:04:05")))
	lg.InfoToFile(fmt.Sprintf("   剩余时间: %v", nextRun.Sub(currentTime)))

	// 显示配置的时间点
	lg.InfoToFile("   配置的时间点:")
	for i, tp := range s.timePoints {
		lg.InfoToFile(fmt.Sprintf("     %d. 星期%s %02d:%02d",
			i+1, formatWeekday(tp.Weekday), tp.Hour, tp.Minute))
	}

	printSeparator("-", 50)
}

func (s *TaskScheduler) calculateNextRunTime(current time.Time) time.Time {
	// 如果没有时间点，返回一天后
	if len(s.timePoints) == 0 {
		return current.AddDate(0, 0, 1)
	}

	// 记录最小时间差和对应的下一个执行时间
	minDuration := time.Duration(365*24) * time.Hour // 一年的最大时长
	var nextRunTime time.Time
	found := false

	// 检查未来7周（足够覆盖所有每周重复的时间点）
	for weekOffset := 0; weekOffset < 7; weekOffset++ {
		for dayOffset := 0; dayOffset < 7; dayOffset++ {
			// 计算要检查的日期
			checkDate := current.AddDate(0, 0, weekOffset*7+dayOffset)

			for _, tp := range s.timePoints {
				// 构造这个时间点
				candidate := time.Date(checkDate.Year(), checkDate.Month(), checkDate.Day(),
					tp.Hour, tp.Minute, 0, 0, s.location)

				// 调整到正确的星期几
				daysToAdjust := (int(tp.Weekday) - int(candidate.Weekday()) + 7) % 7
				candidate = candidate.AddDate(0, 0, daysToAdjust)

				// 检查是否在当前时间之后
				if candidate.After(current) {
					// 计算时间差
					duration := candidate.Sub(current)

					// 如果这是找到的更早的时间
					if !found || duration < minDuration {
						minDuration = duration
						nextRunTime = candidate
						found = true
					}
				}
			}
		}
	}

	if found {
		return nextRunTime
	}

	// 理论上不会执行到这里
	// 返回一周后的第一个时间点作为后备
	firstTp := s.timePoints[0]
	nextWeek := current.AddDate(0, 0, 7)
	nextRun := time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(),
		firstTp.Hour, firstTp.Minute, 0, 0, s.location)

	daysToAdjust := (int(firstTp.Weekday) - int(nextRun.Weekday()) + 7) % 7
	nextRun = nextRun.AddDate(0, 0, daysToAdjust)

	return nextRun
}

//func (s *TaskScheduler) calculateNextRunTime(current time.Time) time.Time {
//	// 找到下一个符合条件的日期和时间
//	for i := 0; i < 365; i++ { // 最多查找一年
//		for _, tp := range s.timePoints {
//			// 构造该时间点
//			candidate := time.Date(current.Year(), current.Month(), current.Day(),
//				tp.Hour, tp.Minute, 0, 0, s.location)
//
//			// 调整星期
//			daysToAdd := (int(tp.Weekday) - int(current.Weekday()) + 7) % 7
//			candidate = candidate.AddDate(0, 0, daysToAdd)
//
//			// 如果候选时间在当前时间之后，返回
//			if candidate.After(current) {
//				return candidate
//			}
//		}
//		// 检查下一天
//		current = current.AddDate(0, 0, 1)
//	}
//
//	// 如果没找到（理论上不会发生），返回一天后
//	return current.AddDate(0, 0, 1)
//}

func (s *TaskScheduler) showTimezoneInfo() {
	// 显示时区信息
	now := time.Now()
	localTime := now.In(s.location)
	utcTime := now.UTC()

	printSeparator("=", 50)
	lg.InfoToFile("🌍 时区信息:")
	lg.InfoToFile(fmt.Sprintf("   配置时区: %s", s.config.Timezone))
	lg.InfoToFile(fmt.Sprintf("   实际时区: %s", s.location))
	lg.InfoToFile(fmt.Sprintf("   系统时间: %s", now.Format("2006-01-02 15:04:05 MST")))
	lg.InfoToFile(fmt.Sprintf("   目标时区时间: %s", localTime.Format("2006-01-02 15:04:05 MST")))
	lg.InfoToFile(fmt.Sprintf("   UTC时间: %s", utcTime.Format("2006-01-02 15:04:05 MST")))

	// 计算时区偏移
	_, offset := localTime.Zone()
	hours := offset / 3600
	minutes := (offset % 3600) / 60

	if minutes > 0 {
		lg.InfoToFile(fmt.Sprintf("   时区偏移: UTC%+d:%02d", hours, minutes))
	} else {
		lg.InfoToFile(fmt.Sprintf("   时区偏移: UTC%+d", hours))
	}
	printSeparator("=", 50)
}

// ReloadConfig 重新加载配置（热重载）
func (s *TaskScheduler) ReloadConfig() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 重新解析时间点
	timePoints, err := parseTimePoints(config.TimePoints)
	if err != nil {
		return err
	}

	// 更新时区
	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.location = location
	s.timePoints = timePoints
	s.config = config
	s.checkInterval = time.Duration(config.CheckInterval) * time.Second

	lg.InfoToFile("✅ 配置重新加载成功")
	return nil
}

// AddTimePoint 添加时间点
func (s *TaskScheduler) AddTimePoint(weekday time.Weekday, hour, minute int) error {
	tp := TimePoint{Weekday: weekday, Hour: hour, Minute: minute}

	// 检查是否已存在
	for _, existing := range s.timePoints {
		if existing.Weekday == tp.Weekday && existing.Hour == tp.Hour && existing.Minute == tp.Minute {
			return fmt.Errorf("时间点已存在: 星期%s %02d:%02d", formatWeekday(weekday), hour, minute)
		}
	}

	s.timePoints = append(s.timePoints, tp)

	// 更新配置文件
	s.config.TimePoints = append(s.config.TimePoints,
		fmt.Sprintf("%d:%02d:%02d", weekday, hour, minute))

	return saveConfig(s.config)
}

// RemoveTimePoint 删除时间点
func (s *TaskScheduler) RemoveTimePoint(weekday time.Weekday, hour, minute int) error {
	newTimePoints := make([]TimePoint, 0, len(s.timePoints))
	newConfigPoints := make([]string, 0, len(s.config.TimePoints))

	targetStr := fmt.Sprintf("%d:%02d:%02d", weekday, hour, minute)
	found := false

	// 从时间点列表删除
	for _, tp := range s.timePoints {
		if !(tp.Weekday == weekday && tp.Hour == hour && tp.Minute == minute) {
			newTimePoints = append(newTimePoints, tp)
		} else {
			found = true
		}
	}

	// 从配置字符串列表删除
	for _, tpStr := range s.config.TimePoints {
		if tpStr != targetStr {
			newConfigPoints = append(newConfigPoints, tpStr)
		}
	}

	if !found {
		return fmt.Errorf("时间点不存在: 星期%s %02d:%02d", formatWeekday(weekday), hour, minute)
	}

	s.timePoints = newTimePoints
	s.config.TimePoints = newConfigPoints

	return saveConfig(s.config)
}

// ListTimePoints 列出所有时间点
func (s *TaskScheduler) ListTimePoints() []string {
	result := make([]string, len(s.timePoints))
	for i, tp := range s.timePoints {
		result[i] = fmt.Sprintf("星期%s %02d:%02d", formatWeekday(tp.Weekday), tp.Hour, tp.Minute)
	}
	return result
}

// GetConfigSummary 获取配置摘要
func (s *TaskScheduler) GetConfigSummary() map[string]interface{} {
	return map[string]interface{}{
		"timezone":       s.config.Timezone,
		"time_points":    s.ListTimePoints(),
		"check_interval": s.config.CheckInterval,
		"total_points":   len(s.timePoints),
	}
}

// 辅助函数：格式化时间点显示
func formatTimePoints(timePoints []TimePoint) string {
	if len(timePoints) == 0 {
		return "无执行时间点"
	}

	var builder strings.Builder
	for i, tp := range timePoints {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("星期%s %02d:%02d",
			formatWeekday(tp.Weekday), tp.Hour, tp.Minute))
	}
	return builder.String()
}

// 辅助函数：格式化星期几
func formatWeekday(weekday time.Weekday) string {
	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
	if int(weekday) < len(weekdays) {
		return weekdays[weekday]
	}
	return fmt.Sprintf("%d", weekday)
}

// 辅助函数：打印分隔线
func printSeparator(char string, length int) {
	lg.InfoToFile(strings.Repeat(char, length))
}
