package domain

const (
	RunModeProd = "prod" // 生产模式
	RunModeDev  = "dev"  // 开发模式
)

type AppConfig struct {
	RunMode                  string                   `default:"dev"`
	ProfitDistributionConfig ProfitDistributionConfig `default:"{TaskHour: 2, TaskMinute: 0}"`
}
