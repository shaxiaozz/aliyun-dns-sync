// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type DomainRecords struct {
	DomainName string `json:"DomainName"`
	Line       string `json:"Line"`
	Locked     bool   `json:"Locked"`
	Rr         string `json:"RR"`
	RecordID   string `json:"RecordId"`
	Status     string `json:"Status"`
	TTL        int    `json:"TTL"`
	Type       string `json:"Type"`
	Value      string `json:"Value"`
	Weight     int    `json:"Weight"`
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Alidns
	config.Endpoint = tea.String("alidns.cn-shenzhen.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

/**
* 使用STS鉴权方式初始化账号Client，推荐此方式。
* @param accessKeyId
* @param accessKeySecret
* @param securityToken
* @return Client
* @throws Exception
 */
func CreateClientWithSTS(accessKeyId *string, accessKeySecret *string, securityToken *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
		// 必填，您的 Security Token
		SecurityToken: securityToken,
		// 必填，表明使用 STS 方式
		Type: tea.String("sts"),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Alidns
	config.Endpoint = tea.String("alidns.cn-shenzhen.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

func _main(args []*string, domainName, accessKeyID, accessKeySecret string) (domainRecordsdata []*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord, _err error) {
	// 请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html
	client, _err := CreateClient(tea.String(accessKeyID), tea.String(accessKeySecret))
	if _err != nil {
		return nil, _err
	}

	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(domainName),
		PageNumber: tea.Int64(1),
		PageSize:   tea.Int64(500),
	}

	runtime := &util.RuntimeOptions{}
	resp, _err := client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, runtime)
	if _err != nil {
		return nil, _err
	}

	// 由于aliyun-sdk PageSize最大只能支持500，因此需要查询第二页数据
	describeDomainRecordsRequestPage2 := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(domainName),
		PageNumber: tea.Int64(2),
		PageSize:   tea.Int64(500),
	}

	runtimePage2 := &util.RuntimeOptions{}
	respPage2, _err := client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequestPage2, runtimePage2)
	if _err != nil {
		return nil, _err
	}

	// 数据合并
	resp.Body.DomainRecords.Record = append(resp.Body.DomainRecords.Record, respPage2.Body.DomainRecords.Record...)

	return resp.Body.DomainRecords.Record, _err
}

// 处理dns解析数据
func DomainRecordsHandler(domainRecordsdata []*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord, fileName string) (err error) {
	// 实例化结构体
	var dnsdata DomainRecords
	// 循环遍历
	for _, domainRecords := range domainRecordsdata {
		domainRecordsUnit, err := json.Marshal(domainRecords)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		domainRecordsStr := string(domainRecordsUnit)
		domainRecordsByte := []byte(domainRecordsStr)
		err = json.Unmarshal([]byte(domainRecordsByte), &dnsdata)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		//fmt.Println("DomainName:", dnsdata.DomainName, "Type:", dnsdata.Type, "RR:", dnsdata.Rr, "Value:", dnsdata.Value)

		// 数据写入
		fileData := "\n" + dnsdata.Rr + "    IN    " + dnsdata.Type + "    " + dnsdata.Value
		err = FileHandler(fileName, fileData)
	}
	fmt.Println("同步阿里云DNS解析完成")
	return nil
}

// 初始化文件
func InitFile(fileName string) (err error) {
	// 初始化文件
	if _, err = os.Stat(fileName); err == nil {
		// 清空文件内容
		err = ioutil.WriteFile(fileName, []byte{}, 0644)
		if err != nil {
			fmt.Println("清空文件内容时出错：", err)
			return err
		}
		fmt.Println(fileName + "文件内容已清空")
	} else if os.IsNotExist(err) {
		// 创建文件
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("创建文件时出错：", err)
			return err
		}
		defer file.Close()
		fmt.Println(fileName, "文件已创建")
		return err
	} else {
		fmt.Println("判断文件是否存在时出错：", err)
		return err
	}

	// 数据写入
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return err
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	initFileText := `$TTL 60
@    IN SOA    @ rname.invalid. (
                    0    ; serial
                    1D    ; refresh
                    1H    ; retry
                    1W    ; expire
                    3H )    ; minimum
    NS    @`
	write.WriteString(initFileText)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return nil
}

// 文件操作
func FileHandler(fileName string, fileData string) (err error) {
	// 数据写入
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return err
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(fileData)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return nil
}

// 执行Bind 配置文件热加载命令
func ExecCommand() (err error) {
	cmd := exec.Command("/bin/bash", "-c", `rndc reload`)

	// 创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return err
	}

	// 执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return err
	}

	// 读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return err
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return err
	}
	fmt.Printf("rndc reload command stdout: %s", bytes)
	return nil
}
func main() {
	// 定义命令行参数
	var filePath string
	var domainName string
	var accessKeyID string
	var accessKeySecret string

	flag.StringVar(&filePath, "file", "", "文件路径")
	flag.StringVar(&domainName, "domainname", "", "阿里云域名")
	flag.StringVar(&accessKeyID, "accesskey-id", "", "Aliyun AccessKey ID")
	flag.StringVar(&accessKeySecret, "accesskey-secret", "", "Aliyun AccessKey Secret")

	// 自定义提示信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	// 解析命令行参数
	flag.Parse()

	// 检查是否有参数
	if len(os.Args) == 1 {
		flag.Usage()
		return
	}
	// 初始化文件
	err := InitFile(filePath)
	if err != nil {
		panic(err)
	}

	// 获取阿里云dns数据
	domainRecordsdata, err := _main(tea.StringSlice(os.Args[1:]), domainName, accessKeyID, accessKeySecret)
	if err != nil {
		panic(err)
	}

	// 处理阿里云dns数据
	err = DomainRecordsHandler(domainRecordsdata, filePath)
	if err != nil {
		panic(err)
	}

	// 重载Bind配置文件
	err = ExecCommand()
	if err != nil {
		panic(err)
	}
}
