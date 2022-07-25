package src

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var (
	searchUrl = "https://www.youdao.com/result?lang=en&word="
	email     = os.Getenv("email163")
	password  = os.Getenv("password163")
	word      = "name"

	cacheDir       string
	cookieFileName string
	cookieValid    bool
)

func init() {
	if email == "" || password == "" {
		errCheck(errors.New("请设置email或password"))
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("获取用户home目录失败: %s", err)
		errCheck(err)
	}

	cacheDir = path.Join(dir, ".youda-translate-shell")
	if _, err = os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, os.ModePerm)
		errCheck(err)
	}

	cookieFileName = path.Join(cacheDir, "cookies.tmp")
}

func Login(keyword string) {
	word = keyword
	Run(myTasks)
	select {}
}

//参考 https://blog.csdn.net/yes169yes123/article/details/109562220
func Run(tasks func() chromedp.Tasks) {
	// chromdp依赖context上限传递参数
	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),

		// 以默认配置的数组为基础，覆写headless参数
		// 当然也可以根据自己的需要进行修改，这个flag是浏览器的设置
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)...,
	)

	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	ctx, _ = chromedp.NewContext(
		ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	// 执行我们自定义的任务 - myTasks函数在第4步
	if err := chromedp.Run(ctx, tasks()); err != nil {
		log.Fatal(err)
		return
	}

	select {}
}

func myTasks() chromedp.Tasks {
	return chromedp.Tasks{login()}
}

// 保存Cookies
func saveCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// cookies的获取对应是在devTools的network面板中
		// 1. 获取cookies
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return
		}

		// 2. 序列化
		cookiesData, _ := json.Marshal(cookies)

		// 3. 存储到临时文件
		if err = ioutil.WriteFile(cookieFileName, cookiesData, 0755); err != nil {
			errCheck(err)
			return
		}

		fmt.Println("cookie saved")
		return
	}
}

// 加载Cookies
func loadCookies() chromedp.ActionFunc {

	return func(ctx context.Context) error {
		// 如果cookies临时文件不存在则直接跳过
		if _, _err := os.Stat(cookieFileName); os.IsNotExist(_err) {
			log.Println("cookie不存在")
			errCheck(_err)
			return nil
		}

		cookieValid = true

		log.Println("cookie存在")
		// 如果存在则读取cookies的数据
		cookiesData, err := ioutil.ReadFile(cookieFileName)
		if err != nil {
			errCheck(err)
			return nil
		}

		// 反序列化
		cookiesParams := network.SetCookiesParams{}
		if err = json.Unmarshal(cookiesData, &(cookiesParams.Cookies)); err != nil {
			errCheck(err)
			return nil
		}

		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

func login() chromedp.Tasks {
	//return chromedp.Tasks{
	//	chromedp.Navigate(searchUrl + word),
	//	//chromedp.Sleep(time.Second * 3),
	//	chromedp.Click("word-operate add"),
	//}
	if cookieValid {
		log.Println("login in")

		return chromedp.Tasks{
			chromedp.Navigate(searchUrl)}
	}

	log.Println("try login")
	return chromedp.Tasks{
		chromedp.Navigate(searchUrl + word),
		//chromedp.Navigate(loginURL),
		//click login button
		chromedp.Click("a[class=login]"),

		// email login
		chromedp.Click(`.j-headimg >div:nth-child(1)`),
		//click and type in email
		chromedp.Click(`input[name=email]`),
		chromedp.SendKeys(`input[name=email]`, email),
		//click and type in passwd
		chromedp.Click(`input[name=password]`),
		chromedp.SendKeys(`input[name=password]`, password),

		//no need login for days
		chromedp.Clear("#un-login"),

		//chromedp.SetAttributes("", map[string]string{""}),
		//do login
		chromedp.Click("#dologin"),
		//chromedp.Sleep(time.Second * 3),
		chromedp.Click("#searchLayout > div > div.search_result.center_container > div > div > section > div.simple-explain > div > div > div > div > div.title > div > a.word-operate.add.added"),
		saveCookies(),
	}
}
