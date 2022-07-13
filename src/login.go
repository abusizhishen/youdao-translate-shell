package src

import (
	"context"
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
	searchUrl = "https://www.youdao.com/"

	loginURL = "https://c.youdao.com/common-login-web/index.html"
	email    = os.Getenv("email163")
	password = os.Getenv("password163")

	cacheDir       string
	cookieFileName string
	cookieValid    bool
)

func init() {
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

//参考 https://blog.csdn.net/yes169yes123/article/details/109562220
func Login() {
	// chromdp依赖context上限传递参数
	ctx, _ := chromedp.NewExecAllocator(
		context.Background(),

		// 以默认配置的数组为基础，覆写headless参数
		// 当然也可以根据自己的需要进行修改，这个flag是浏览器的设置
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)...,
	)

	ctx, _ = context.WithTimeout(ctx, 30*time.Second)
	ctx, _ = chromedp.NewContext(
		ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),
	)

	// 执行我们自定义的任务 - myTasks函数在第4步
	if err := chromedp.Run(ctx, myTasks()); err != nil {
		log.Fatal(err)
		return
	}

	select {}
}

func myTasks() chromedp.Tasks {
	loadCookies()
	if cookieValid {
		return chromedp.Tasks{
			chromedp.Navigate(searchUrl),
			doSomething(),
		}
	}

	return login()
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
		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return
		}

		// 3. 存储到临时文件
		if err = ioutil.WriteFile(cookieFileName, cookiesData, 0755); err != nil {
			errCheck(err)
			return
		}
		return
	}
}

// 加载Cookies
func loadCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 如果cookies临时文件不存在则直接跳过
		if _, _err := os.Stat(cookieFileName); os.IsNotExist(_err) {
			errCheck(err)
			return
		}

		// 如果存在则读取cookies的数据
		cookiesData, err := ioutil.ReadFile(cookieFileName)
		if err != nil {
			errCheck(err)
			return
		}

		// 反序列化
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			errCheck(err)
			return
		}

		cookieValid = true
		// 设置cookies
		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

func login() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(searchUrl),
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
		chromedp.Click("#un-login"),

		//do login
		chromedp.Click("#dologin"),
		saveCookies(),
	}
}

func doSomething() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		select {}
	}
}
