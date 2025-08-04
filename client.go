package riskbird

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// NewClient 创建新的风鸟查询客户端
func NewClient(config *Config) *Client {
	if config == nil {
		config = &Config{
			Timeout: 5,
			Delay:   1,
		}
	}
	return &Client{
		config: config,
	}
}

// Search 根据关键词搜索企业（获取所有页面的结果）
func (c *Client) Search(name string) ([]SearchResult, error) {
	url := "https://www.riskbird.com/riskbird-api/newSearch"
	var allResults []SearchResult
	pageNo := 1
	pageSize := 10

	for {
		searchData := map[string]string{
			"searchKey":           name,
			"pageNo":              fmt.Sprintf("%d", pageNo),
			"range":               fmt.Sprintf("%d", pageSize),
			"referer":             "search",
			"queryType":           "1",
			"selectConditionData": "{\"status\":\"\",\"sort_field\":\"\"}",
		}

		marshal, err := json.Marshal(searchData)
		if err != nil {
			return nil, fmt.Errorf("关键词处理失败: %s", err.Error())
		}

		content, err := c.request(url, string(marshal))
		if err != nil {
			return nil, err
		}

		// 检查API是否成功
		if !gjson.Get(content, "success").Bool() {
			return nil, fmt.Errorf("搜索失败: %s", gjson.Get(content, "msg").String())
		}

		enList := gjson.Get(content, "data.list").Array()
		if len(enList) == 0 {
			if pageNo == 1 {
				return nil, fmt.Errorf("没有查询到关键词: %s", name)
			}
			break // 没有更多数据了
		}

		// 解析当前页的结果
		for _, item := range enList {
			var result SearchResult
			if err := json.Unmarshal([]byte(item.Raw), &result); err == nil {
				allResults = append(allResults, result)
			} else {
				log.Printf("解析企业信息失败: %s", err.Error())
			}
		}

		// 检查是否还有下一页
		totalCount := gjson.Get(content, "data.total").Int()
		if int64(len(allResults)) >= totalCount {
			break // 已获取所有数据
		}

		// 如果当前页数据少于页面大小，说明是最后一页
		if len(enList) < pageSize {
			break
		}

		pageNo++

		// 添加延时避免请求过快
		if c.config.Delay > 0 {
			time.Sleep(time.Duration(c.config.Delay) * time.Second)
		}
	}

	return allResults, nil
}

// GetCompanyInfo 根据企业ID获取企业基本信息
func (c *Client) GetCompanyInfo(entid string) (*CompanyInfo, error) {
	detailRes, _, err := c.getBaseInfo(entid)
	if err != nil {
		return nil, err
	}
	var info CompanyInfo
	if err := json.Unmarshal([]byte(detailRes.Raw), &info); err != nil {
		return nil, fmt.Errorf("解析企业信息失败: %s", err.Error())
	}
	return &info, nil
}

// GetICPInfo 获取ICP备案信息
func (c *Client) GetICPInfo(entid string, page int) (*PageInfo, error) {
	return c.getInfoByPage(entid, page, "propertyIcp")
}

// GetAllICPInfo 获取所有ICP备案信息
func (c *Client) GetAllICPInfo(entid string) ([]gjson.Result, error) {
	return c.getAllInfoByType(entid, "propertyIcp")
}

// GetAppInfo 获取APP信息
func (c *Client) GetAppInfo(entid string, page int) (*PageInfo, error) {
	return c.getInfoByPage(entid, page, "propertyApp")
}

// GetAllAppInfo 获取所有APP信息
func (c *Client) GetAllAppInfo(entid string) ([]gjson.Result, error) {
	return c.getAllInfoByType(entid, "propertyApp")
}

// GetWxAppInfo 获取微信小程序信息
func (c *Client) GetWxAppInfo(entid string, page int) (*PageInfo, error) {
	return c.getInfoByPage(entid, page, "propertyMiniprogram")
}

// GetAllWxAppInfo 获取所有微信小程序信息
func (c *Client) GetAllWxAppInfo(entid string) ([]gjson.Result, error) {
	return c.getAllInfoByType(entid, "propertyMiniprogram")
}

// GetJobInfo 获取招聘信息
func (c *Client) GetJobInfo(entid string, page int) (*PageInfo, error) {
	return c.getInfoByPage(entid, page, "propertyJob")
}

// GetAllJobInfo 获取所有招聘信息
func (c *Client) GetAllJobInfo(entid string) ([]gjson.Result, error) {
	return c.getAllInfoByType(entid, "propertyJob")
}

// GetCopyrightInfo 获取软件著作权信息
func (c *Client) GetCopyrightInfo(entid string, page int) (*PageInfo, error) {
	return c.getInfoByPage(entid, page, "propertyCopyright")
}

// GetAllCopyrightInfo 获取所有软件著作权信息
func (c *Client) GetAllCopyrightInfo(entid string) ([]gjson.Result, error) {
	return c.getAllInfoByType(entid, "propertyCopyright")
}

// GetInvestInfo 获取投资信息
func (c *Client) GetInvestInfo(entid string, page int) (*PageInfo, error) {
	return c.getInfoByPage(entid, page, "companyInvest")
}

// GetAllInvestInfo 获取所有投资信息
func (c *Client) GetAllInvestInfo(entid string) ([]gjson.Result, error) {
	return c.getAllInfoByType(entid, "companyInvest")
}

// GetBranchInfo 获取分支机构信息
func (c *Client) GetBranchInfo(entid string, page int) (*PageInfo, error) {
	return c.getInfoByPage(entid, page, "propertyBranch")
}

// GetAllBranchInfo 获取所有分支机构信息
func (c *Client) GetAllBranchInfo(entid string) ([]gjson.Result, error) {
	return c.getAllInfoByType(entid, "propertyBranch")
}

// GetPartnerInfo 获取股东信息
func (c *Client) GetPartnerInfo(entid string, page int) (*PageInfo, error) {
	return c.getInfoByPage(entid, page, "propertyPartner")
}

// GetAllPartnerInfo 获取所有股东信息
func (c *Client) GetAllPartnerInfo(entid string) ([]gjson.Result, error) {
	return c.getAllInfoByType(entid, "propertyPartner")
}

// getAllInfoByType 获取指定类型的所有数据（自动分页）
func (c *Client) getAllInfoByType(entid string, extractType string) ([]gjson.Result, error) {
	var allData []gjson.Result
	page := 1

	for {
		pageInfo, err := c.getInfoByPage(entid, page, extractType)
		if err != nil {
			// 如果是500错误且是第一页，可能是该企业没有此类数据，返回空数组而不是错误
			if page == 1 && (strings.Contains(err.Error(), "500") || strings.Contains(err.Error(), "未知错误")) {
				return []gjson.Result{}, nil
			}
			return nil, err
		}

		// 如果没有数据，直接返回
		if pageInfo.Total == 0 {
			return []gjson.Result{}, nil
		}

		// 添加当前页的数据
		allData = append(allData, pageInfo.Data...)

		// 检查是否还有更多数据
		if int64(len(allData)) >= pageInfo.Total {
			break
		}

		page++
		// 添加延时避免请求过快
		if c.config.Delay > 0 {
			time.Sleep(time.Duration(c.config.Delay) * time.Second)
		}
	}

	return allData, nil
}

// getInfoByPage 通用的分页信息获取方法
func (c *Client) getInfoByPage(entid string, page int, extractType string) (*PageInfo, error) {
	// 先获取企业基本信息以获得正确的orderNo
	baseResponse, err := c.request("https://www.riskbird.com/api/ent/query?entId="+entid, "")
	if err != nil {
		return nil, err
	}
	r := gjson.Parse(baseResponse)
	orderNo := r.Get("orderNo").String()
	if orderNo == "" {
		return nil, fmt.Errorf("无法获取orderNo")
	}

	url := "https://www.riskbird.com/riskbird-api/companyInfo/list"
	li := map[string]interface{}{
		"filterCnd":   0,
		"page":        page,
		"size":        100,
		"orderNo":     orderNo,
		"extractType": extractType,
		"sortField":   "",
		"filterMap":   map[string]interface{}{},
	}
	marshal, err := json.Marshal(li)
	if err != nil {
		return nil, err
	}
	response, err := c.request(url, string(marshal))
	if err != nil {
		return nil, err
	}
	content := gjson.Parse(response)
	if content.Get("code").String() != "20000" {
		// 调试输出
		log.Printf("API响应: %s", response)
		return nil, fmt.Errorf("获取数据失败: code=%s, msg=%s", content.Get("code").String(), content.Get("msg").String())
	}

	info := &PageInfo{
		Size:  100,
		Total: content.Get("data.totalCount").Int(),
		Data:  content.Get("data.apiData").Array(),
	}
	return info, nil
}

// getBaseInfo 获取企业基本信息
func (c *Client) getBaseInfo(entid string) (gjson.Result, gjson.Result, error) {
	response, err := c.request("https://www.riskbird.com/api/ent/query?entId="+entid, "")
	if err != nil {
		return gjson.Result{}, gjson.Result{}, err
	}
	r := gjson.Parse(response)
	result := r.Get("basicResult.apiData.list.jbxxInfo")
	enJsonTMP, _ := sjson.Set(result.Raw, "orderNo", r.Get("orderNo").String())
	enBaseInfo := r.Get("basicResult.apiData.count")
	return gjson.Parse(enJsonTMP), enBaseInfo, nil
}

// request 发送HTTP请求
func (c *Client) request(url string, data string) (string, error) {
	client := req.C()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTLSFingerprintChrome()

	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.60 Safari/537.36",
		"Accept":       "text/html,application/json,application/xhtml+xml, image/jxr, */*",
		"App-Device":   "WEB",
		"Content-Type": "application/json",
		"Origin":       "https://www.riskbird.com",
		"Referer":      "https://www.riskbird.com/ent/",
	}

	if c.config.Cookie != "" {
		headers["Cookie"] = c.config.Cookie
	}

	if c.config.UserAgent != "" {
		headers["User-Agent"] = c.config.UserAgent
	}

	client.SetCommonHeaders(headers)

	if c.config.Proxy != "" {
		client.SetProxyURL(c.config.Proxy)
	}

	if c.config.Timeout > 0 {
		client.SetTimeout(time.Duration(c.config.Timeout) * time.Minute)
	}

	if c.config.Delay > 0 {
		time.Sleep(time.Duration(c.config.Delay) * time.Second)
	}

	req := client.R()

	method := "GET"
	if data != "" {
		method = "POST"
		req.SetBody(data)
		if !strings.Contains(url, "newSearch") {
			req.SetHeader("Xs-Content-Type", "application/json")
		}
	}

	resp, err := req.Send(method, url)
	if err != nil {
		log.Printf("请求错误 %s: %s\n", url, err)
		time.Sleep(5 * time.Second)
		return c.request(url, data)
	}

	if resp.StatusCode == 200 {
		rs := gjson.Parse(resp.String())
		if rs.Get("state").String() == "limit:auth" {
			return "", errors.New("您今日的查询次数已达到上限, 请更换cookie")
		}
		return resp.String(), nil
	} else {
		switch resp.StatusCode {
		case 403:
			return "", errors.New("IP被禁止访问网站，请更换IP")
		case 401:
			return "", errors.New("Cookie有问题或过期，请重新获取")
		case 302:
			return "", errors.New("需要更新Cookie")
		case 404:
			return "", fmt.Errorf("请求错误 404: %s", url)
		default:
			return "", fmt.Errorf("未知错误: %d", resp.StatusCode)
		}
	}
}
