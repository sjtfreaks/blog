package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// 发起 HTTP GET 请求
	res, err := http.Get("https://www.github.com/trending")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("请求失败，状态码：%d", res.StatusCode)
	}

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 当前日期
	today := time.Now().Format("2006-01-02")
	md_name := "github_trending_" + today + ".md"

	//判断文件具体目录是否存在
	_, err = os.Stat("content/posts/github" + md_name)
	if err == nil {
		//文件存在
		os.Remove(md_name)
	}

	// 创建 Markdown 文件
	file, err := os.Create("content/posts/github" + md_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 写入 Markdown 文件头部
	_, err = file.WriteString("---\ntitle: " + today + " GitHub 热门榜单\ndate: " + today + "\ndraft: true\nauthor: 'jobcher'\ntags: ['github']\ncategories: ['github']\nseries: ['github']\n---\n\n")
	if err != nil {
		log.Fatal(err)
	}

	// 查找所有的 trending repository
	doc.Find(".Box .Box-row").Each(func(i int, s *goquery.Selection) {
		// 提取标题和作者
		title := strings.TrimSpace(s.Find("h2.h3 a").AttrOr("href", ""))
		author := strings.TrimSpace(s.Find("span.text-normal").First().Text())
		url := strings.TrimSpace(s.Find("h2.h3 a").AttrOr("href", ""))
		desc := strings.TrimSpace(s.Find("p.col-9").Text())

		// 输出标题和作者
		fmt.Printf("Repository %d:\n", i+1)
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Description: %s\n", desc)
		fmt.Printf("URL: https://github.com%s\n", url)
		fmt.Printf("Author: %s\n\n", author)

		// 将信息以 Markdown 格式写入文件
		content := fmt.Sprintf("### Repository %d:\n", i+1)
		content += fmt.Sprintf("- Title: %s\n", title)
		content += fmt.Sprintf("- Description: %s\n", desc)
		content += fmt.Sprintf("- URL: https://github.com%s\n", url)
		content += fmt.Sprintf("- Author: %s\n\n", author)

		_, err = file.WriteString(content)
		if err != nil {
			log.Fatal(err)
		}

	})
	fmt.Println("Done.")

	//使用go 提交本地代码
	cmd := exec.Command("git add . && git commit -m 'update' && git push origin master")
	cmd.Output()
	cmd.Run()
	fmt.Println("提交成功")

}
