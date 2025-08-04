# RiskBird Go SDK

风鸟企业信息查询 Go 语言 SDK。

## 功能特性

- 🔍 企业信息搜索
- 🏢 企业基本信息查询  
- 🌐 ICP备案信息查询
- 📱 APP信息查询
- 📱 微信小程序信息查询
- 💼 招聘信息查询
- 📄 软件著作权信息查询
- 💰 投资信息查询（注意：部分企业可能无投资数据）
- 🏪 分支机构信息查询
- 👥 股东信息查询

## 安装

```bash
go get github.com/dean2021/riskbird
```

## 快速开始

```go
import (
    "fmt"
    "log"
    riskbird "github.com/dean2021/riskbird"
)

// 创建配置
config := &riskbird.Config{
    Cookie:    "your_cookie_here", // 风鸟网站的Cookie
    UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    Timeout:   5, // 5分钟超时
    Delay:     2, // 2秒延时
}

// 创建客户端
client := riskbird.NewClient(config)

// 搜索企业
results, err := client.Search("腾讯")
if err != nil {
    log.Fatalf("搜索失败: %v", err)
}

// 获取企业信息
for _, result := range results {
    fmt.Printf("企业: %s, ID: %s\n", result.EntName, result.Entid)
    
    // 获取详细信息
    info, _ := client.GetCompanyInfo(result.Entid)
    icpInfo, _ := client.GetICPInfo(result.Entid, 1)
    appInfo, _ := client.GetAppInfo(result.Entid, 1)
    
    fmt.Printf("法人: %s, ICP: %d条, APP: %d个\n", 
        info.PersonName, icpInfo.Total, appInfo.Total)
}
```

## 主要方法

### 搜索相关
- `Search(name string) ([]SearchResult, error)` - 搜索企业（获取所有页面结果）

### 企业信息
- `GetCompanyInfo(entid string) (*CompanyInfo, error)` - 获取企业基本信息

### 分页查询方法
- `GetICPInfo(entid string, page int) (*PageInfo, error)` - 获取ICP备案信息
- `GetAppInfo(entid string, page int) (*PageInfo, error)` - 获取APP信息
- `GetWxAppInfo(entid string, page int) (*PageInfo, error)` - 获取微信小程序信息
- `GetJobInfo(entid string, page int) (*PageInfo, error)` - 获取招聘信息
- `GetCopyrightInfo(entid string, page int) (*PageInfo, error)` - 获取软件著作权信息
- `GetInvestInfo(entid string, page int) (*PageInfo, error)` - 获取投资信息
- `GetBranchInfo(entid string, page int) (*PageInfo, error)` - 获取分支机构信息
- `GetPartnerInfo(entid string, page int) (*PageInfo, error)` - 获取股东信息

## Cookie 配置

⚠️ **重要**：如果遇到 `{"code":9999,"msg":"数据错误"}` 错误，说明Cookie无效或未配置！

### 获取Cookie

1. 访问 [风鸟网站](https://www.riskbird.com) 并登录
2. 打开浏览器开发者工具（F12）→ Network 标签页
3. 在网站上搜索企业，查看请求头中的 Cookie
4. 复制完整的 Cookie 字符串

**Cookie必须包含关键字段**：`app-uuid`、`token`、`userinfo`

### 配置选项

```go
type Config struct {
    Cookie    string // 风鸟网站的Cookie（必需）
    UserAgent string // 自定义User-Agent
    Proxy     string // 代理地址
    Timeout   int    // 超时时间（分钟），默认5分钟
    Delay     int    // 请求延时（秒），默认1秒
}
```

### 注意事项

- 建议设置适当的请求延时（1-3秒）
- Cookie可能会过期，需要定期更新
- 某些数据可能需要付费账号才能访问
- 支持代理设置，适用于网络受限环境

## 数据结构

### SearchResult - 搜索结果
```go
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
```

### CompanyInfo - 企业基本信息
```go
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
```

### PageInfo - 分页信息
```go
type PageInfo struct {
    Total   int64          `json:"total"`   // 总数
    Page    int64          `json:"page"`    // 当前页
    Size    int64          `json:"size"`    // 每页大小
    HasNext bool           `json:"hasNext"` // 是否有下一页
    Data    []gjson.Result `json:"data"`    // 数据
}
```

## 运行示例

```bash
cd example
go run main.go
```

## 使用技巧

### 🔍 查询企业信息

#### 基本信息查询
```go
// 获取企业基本信息
companyInfo, err := client.GetCompanyInfo(entid)
```

#### 分页查询（单页数据）
```go
// 获取ICP备案信息（分页）
icpInfo, err := client.GetICPInfo(entid, 1) // 第1页

// 获取APP信息（分页）
appInfo, err := client.GetAppInfo(entid, 1)

// 获取微信小程序信息（分页）
wxAppInfo, err := client.GetWxAppInfo(entid, 1)

// 获取招聘信息（分页）
jobInfo, err := client.GetJobInfo(entid, 1)

// 获取软件著作权信息（分页）
copyrightInfo, err := client.GetCopyrightInfo(entid, 1)

// 获取投资信息（分页）
investInfo, err := client.GetInvestInfo(entid, 1)

// 获取分支机构信息（分页）
branchInfo, err := client.GetBranchInfo(entid, 1)

// 获取股东信息（分页）
partnerInfo, err := client.GetPartnerInfo(entid, 1)
```

#### 获取全部数据（自动分页）
```go
// 获取所有ICP备案信息
allICPInfo, err := client.GetAllICPInfo(entid)

// 获取所有APP信息
allAppInfo, err := client.GetAllAppInfo(entid)

// 获取所有微信小程序信息
allWxAppInfo, err := client.GetAllWxAppInfo(entid)

// 获取所有招聘信息
allJobInfo, err := client.GetAllJobInfo(entid)

// 获取所有软件著作权信息
allCopyrightInfo, err := client.GetAllCopyrightInfo(entid)

// 获取所有投资信息
allInvestInfo, err := client.GetAllInvestInfo(entid)

// 获取所有分支机构信息
allBranchInfo, err := client.GetAllBranchInfo(entid)

// 获取所有股东信息
allPartnerInfo, err := client.GetAllPartnerInfo(entid)
```

### 错误处理

```go
// 检查Cookie是否有效
results, err := client.Search("测试企业")
if err != nil {
    if strings.Contains(err.Error(), "数据错误") {
        log.Println("Cookie无效或已过期，请更新Cookie")
    }
    return err
}
```

### 批量查询

```go
// 搜索多个企业
companies := []string{"腾讯", "阿里巴巴", "百度"}
for _, company := range companies {
    results, err := client.Search(company)
    if err != nil {
        log.Printf("搜索 %s 失败: %v", company, err)
        continue
    }
    
    // 处理结果...
    time.Sleep(time.Duration(config.Delay) * time.Second) // 延时
}
```

### 选择查询方式

```go
// 方式1: 获取单页数据（适合预览或分页显示）
icpInfo, err := client.GetICPInfo(entid, 1)
if err == nil {
    fmt.Printf("总数: %d, 当前页: %d条\n", icpInfo.Total, len(icpInfo.Data))
}

// 方式2: 获取所有数据（推荐，自动处理分页）
allICPInfo, err := client.GetAllICPInfo(entid)
if err == nil {
    fmt.Printf("获取到所有ICP信息: %d条\n", len(allICPInfo))
    // 直接处理所有数据
    for _, icp := range allICPInfo {
        // 处理每条ICP信息
    }
}
```

### 特别说明

#### 投资信息查询
投资信息API可能因为以下原因返回空数据：
- 企业确实没有对外投资
- API参数配置问题
- 数据源暂时不可用

如果遇到投资信息查询问题，建议：
1. 先验证其他API（如ICP、APP信息）是否正常
2. 检查Cookie是否有效
3. 尝试查询其他已知有投资信息的企业

```go
// 投资信息查询示例
allInvestInfo, err := client.GetAllInvestInfo(entid)
if err != nil {
    log.Printf("投资信息查询失败: %v", err)
} else if len(allInvestInfo) == 0 {
    fmt.Println("该企业暂无投资信息")
} else {
    fmt.Printf("投资企业总数: %d\n", len(allInvestInfo))
}
```

## 依赖项

- `github.com/imroc/req/v3` - HTTP客户端
- `github.com/tidwall/gjson` - JSON解析
- `github.com/tidwall/sjson` - JSON构建

## 版本要求

- Go 1.21+

## 免责声明

本工具仅用于合法的信息收集目的，请勿用于非法用途。使用者需要自行承担使用风险。