package main

import (
	"fmt"
	//"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var start, end int
	fmt.Printf("请输入起始页：")
	fmt.Scan(&start)
	fmt.Printf("请输入终止页：")
	fmt.Scan(&end)

	DoWork(start, end) //工作函数
}
func SpiderOneJoy(url string) (title, content string, err error) {
	rst, err4 := HttpGet(url)
	if err4 != nil {
		err = err4
		return
	}
	//取关键信息
	//取标题<h1> 只取一个
	rel := regexp.MustCompile(`<h1>(?s:(.*?))</h1>`)
	if rel == nil {
		err = fmt.Errorf("%s", "regexp.MustCompile(`<h1>(?s:(.*?))<h1>`)")
		return
	}
	//取标题内容
	tmpTitle := rel.FindAllStringSubmatch(rst, 1) //最后一个参数为1，只过滤第一个
	for _, data := range tmpTitle {

		title = data[1]
		title = strings.Replace(title, "\t", "", -1)
		break
	}

	rel2 := regexp.MustCompile(`<div class="content-txt pt10">(?s:(.*?))<a id="prev" href="`)
	if rel2 == nil {
		err = fmt.Errorf("regexp.MustCompile err")
		return
	}
	//取标题内容
	tmpContent := rel2.FindAllStringSubmatch(rst, -1) //最后一个参数为1，只过滤第一个
	for _, data := range tmpContent {
		content = data[1]
		content = strings.Replace(content, "\t", "", -1)
		content = strings.Replace(content, "<br />", "", -1)
		content = strings.Replace(content, "\n", "", -1)
		content = strings.Replace(content, "&nbsp;", "", -1)
		break
	}
	return

}

func DoWork(start int, end int) {
	fmt.Printf("正在爬取第%d页到第%d页的网址\n", start, end)
	page := make(chan int)

	for i := start; i <= end; i++ {
		//定义一个函数、爬取主页面
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第 %d 个页面爬取完成", <-page)
	}
}

func StoreJoyToFile(i int, fileTitle, fileContent []string) {
	file, err5 := os.Create(strconv.Itoa(i) + ".txt")
	if err5 != nil {
		fmt.Println("os.Create", err5)
		return
	}
	defer file.Close()
	n := len(fileContent)
	for i := 0; i < n; i++ {
		file.Write([]byte(fileTitle[i] + " :\r\n"))
		file.Write([]byte(fileContent[i] + "\r\n"))
		file.Write([]byte("*************************************************************\r\n"))

	}

}
func SpiderPage(i int, page chan int) {
	//爬取的url
	url := "https://www.pengfu.com/xiaohua_" + strconv.Itoa(i) + ".html"
	fmt.Printf("正在爬取第%d个页面：%s\n", i, url)

	rst, err := HttpGet(url) //rst=获取网页内容
	if err != nil {
		fmt.Println("main-38:HttpGet(url) err=", err)
	}

	re := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))"`)
	if re == nil {
		fmt.Println("匹配正则表达式出错")
		return
	}

	//取关键信息
	joyUrls := re.FindAllStringSubmatch(rst, -1)
	fileTitle := make([]string, 0)
	fileContent := make([]string, 0)

	//从关键信息中过滤出子url
	for _, data := range joyUrls {
		//fmt.Println("url=", data[0])
		//根据url爬取每一个笑话
		title, content, err2 := SpiderOneJoy(data[1])
		if err2 != nil {
			fmt.Println("SpiderOneJoy(data[1]) err=", err2)
			return
		}
		//fmt.Println("title= ", title)
		//fmt.Println("content= ", content)
		fileTitle = append(fileTitle, title)
		fileContent = append(fileContent, content)

	}
	fmt.Printf("本页共搜集到 %d 条\n", len(fileTitle)) //追加内容

	StoreJoyToFile(i, fileTitle, fileContent)
	page <- i

}
func HttpGet(url string) (rst string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	//读取网页内容
	buf := make([]byte, 4*1024)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		rst += string(buf[:n])
	}
	return rst, err
}
