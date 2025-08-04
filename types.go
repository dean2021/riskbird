package riskbird

import (
	"github.com/tidwall/gjson"
)

// Config 风鸟查询配置
type Config struct {
	Cookie    string // 风鸟网站的Cookie
	UserAgent string // 自定义User-Agent
	Proxy     string // 代理地址
	Timeout   int    // 超时时间（分钟）
	Delay     int    // 请求延时（秒）
}

// Client 风鸟查询客户端
type Client struct {
	config *Config
}

// CompanyInfo 企业基本信息
type CompanyInfo struct {
	EntName    string   `json:"entName"`    // 企业名称
	PersonName string   `json:"personName"` // 法人代表
	EntStatus  string   `json:"entStatus"`  // 经营状态
	TelList    []string `json:"telList"`    // 电话
	EmailList  []string `json:"emailList"`  // 邮箱
	RecConcat  string   `json:"recConcat"`  // 注册资本
	EsDate     string   `json:"esDate"`     // 成立日期
	YrAddress  string   `json:"yrAddress"`  // 注册地址
	OpScope    string   `json:"opScope"`    // 经营范围
	Uniscid    string   `json:"uniscid"`    // 统一社会信用代码
	Entid      string   `json:"entid"`      // 企业ID
}

// SearchResult 搜索结果
type SearchResult struct {
	EntName   string   `json:"ENTNAME"`   // 企业名称
	Faren     string   `json:"faren"`     // 法人代表
	EntStatus string   `json:"ENTSTATUS"` // 经营状态
	Tels      []string `json:"tels"`      // 电话
	Emails    []string `json:"emails"`    // 邮箱
	RegConcat string   `json:"regConcat"` // 注册资本
	EsDate    string   `json:"esDate"`    // 成立日期
	Dom       string   `json:"dom"`       // 注册地址
	Uniscid   string   `json:"UNISCID"`   // 统一社会信用代码
	Entid     string   `json:"entid"`     // 企业ID
}

// ICPInfo ICP备案信息
type ICPInfo struct {
	Webname  string `json:"webname"`  // 网站名称
	Hostname string `json:"hostname"` // 域名
	Icpnum   string `json:"icpnum"`   // 网站备案/许可证号
}

// AppInfo APP信息
type AppInfo struct {
	Appname            string `json:"appname"`            // 名称
	UpdateDateAndroid  string `json:"updateDateAndroid"`  // 更新时间
	Brief              string `json:"brief"`              // 简介
	IconUrl            string `json:"iconUrl"`            // logo
	DownloadCountLevel string `json:"downloadCountLevel"` // market
}

// WxAppInfo 小程序信息
type WxAppInfo struct {
	Name   string `json:"name"`   // 名称
	Cate   string `json:"cate"`   // 分类
	Logo   string `json:"logo"`   // 头像
	Qrcode string `json:"qrcode"` // 二维码
}

// JobInfo 招聘信息
type JobInfo struct {
	Position  string `json:"position"`  // 招聘职位
	Education string `json:"education"` // 学历要求
	Region    string `json:"region"`    // 工作地点
	Pdate     string `json:"pdate"`     // 发布日期
}

// CopyrightInfo 软件著作权信息
type CopyrightInfo struct {
	Sname string `json:"sname"` // 软件名称
	Snum  string `json:"snum"`  // 登记号
}

// InvestInfo 投资信息
type InvestInfo struct {
	EntName     string `json:"entName"`     // 企业名称
	PersonName  string `json:"personName"`  // 法人
	EntStatus   string `json:"entStatus"`   // 状态
	FunderRatio string `json:"funderRatio"` // 投资比例
	Entid       string `json:"entid"`       // PID
}

// BranchInfo 分支机构信息
type BranchInfo struct {
	BrName      string `json:"brName"`      // 企业名称
	BrPrincipal string `json:"brPrincipal"` // 法人
	EntStatus   string `json:"entStatus"`   // 状态
	Entid       string `json:"entid"`       // PID
}

// PartnerInfo 股东信息
type PartnerInfo struct {
	ShaName     string `json:"shaName"`     // 股东名称
	FundedRatio string `json:"fundedRatio"` // 持股比例
	SubConAm    string `json:"subConAm"`    // 认缴出资金额
	ShaId       string `json:"shaId"`       // PID
}

// PageInfo 分页信息
type PageInfo struct {
	Total   int64          `json:"total"`   // 总数
	Page    int64          `json:"page"`    // 当前页
	Size    int64          `json:"size"`    // 每页大小
	HasNext bool           `json:"hasNext"` // 是否有下一页
	Data    []gjson.Result `json:"data"`    // 数据
}

// EntityMap 实体映射配置
type EntityMap struct {
	Name    string            // 接口名字
	Api     string            // API 地址
	Field   []string          // 获取的字段名称
	KeyWord []string          // 关键词
	Total   int64             // 统计数量
	GNum    string            // 获取数量的json关键词
	SData   map[string]string // 接口请求POST参数
}
