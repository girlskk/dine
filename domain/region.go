package domain

type Country string

const (
	CountryMY Country = "MY"
	CountrySG Country = "SG"
)

func (Country) Values() []string {
	return []string{
		string(CountryMY),
		string(CountrySG),
	}
}

type Province string

const (
	// 马来西亚地区
	ProvinceMY01 Province = "MY-01" // 柔佛州
	ProvinceMY02 Province = "MY-02" // 吉打州
	ProvinceMY03 Province = "MY-03" // 吉兰丹州
	ProvinceMY04 Province = "MY-04" // 马六甲州
	ProvinceMY05 Province = "MY-05" // 森美兰州
	ProvinceMY06 Province = "MY-06" // 彭亨州
	ProvinceMY07 Province = "MY-07" // 槟城州
	ProvinceMY08 Province = "MY-08" // 霹雳州
	ProvinceMY09 Province = "MY-09" // 玻璃市州
	ProvinceMY10 Province = "MY-10" // 雪兰莪州
	ProvinceMY11 Province = "MY-11" // 登嘉楼州
	ProvinceMY12 Province = "MY-12" // 沙巴州
	ProvinceMY13 Province = "MY-13" // 砂拉越州
	ProvinceMY14 Province = "MY-14" // 吉隆坡联邦直辖区
	ProvinceMY15 Province = "MY-15" // 纳闽联邦直辖区
	ProvinceMY16 Province = "MY-16" // 布城联邦直辖区

	// 新加坡地区
	ProvinceSG01 Province = "SG-01" // 中区社区发展理事会
	ProvinceSG02 Province = "SG-02" // 东北社区发展理事会
	ProvinceSG03 Province = "SG-03" // 西北社区发展理事会
	ProvinceSG04 Province = "SG-04" // 东南社区发展理事会
	ProvinceSG05 Province = "SG-05" // 西南社区发展理事会
)

func (Province) Values() []string {
	return []string{
		string(ProvinceMY01),
		string(ProvinceMY02),
		string(ProvinceMY03),
		string(ProvinceMY04),
		string(ProvinceMY05),
		string(ProvinceMY06),
		string(ProvinceMY07),
		string(ProvinceMY08),
		string(ProvinceMY09),
		string(ProvinceMY10),
		string(ProvinceMY11),
		string(ProvinceMY12),
		string(ProvinceMY13),
		string(ProvinceMY14),
		string(ProvinceMY15),
		string(ProvinceMY16),
		string(ProvinceSG01),
		string(ProvinceSG02),
		string(ProvinceSG03),
		string(ProvinceSG04),
		string(ProvinceSG05),
	}
}
