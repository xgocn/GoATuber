package control

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func initControlRouter() {
	r := gin.Default()
	e.R = r

	//加载静态控制器页面
	r.Use(static.Serve("/control", static.LocalFile("./control/frontend", true)))
	r.NoRoute(func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			content, err := os.ReadFile("./control/frontend/index.html")
			if err != nil {
				c.Writer.WriteHeader(404)
				c.Writer.WriteString("Not Found")
				return
			}
			c.Writer.WriteHeader(200)
			c.Writer.Header().Add("Accept", "text/html")
			c.Writer.Write(content)
			c.Writer.Flush()
		}
	})

	//获取配置文件信息
	r.GET("/init", getConfig)

	//配置文件改动
	config := r.Group("/config")
	{
		config.POST("/dict", modifyDict)
		config.POST("/filter", modifyFilter)
		config.POST("/llm", modifyLLM)
		config.POST("/monitor", modifyMonitor)
		config.POST("/mood", modifyMood)
		config.POST("/proxy", modifyProxy)
		config.POST("/speech", modifySpeech)
		config.POST("/voice", modifyVoice)
		//监听包
		config.POST("/bilibili", modifyBiliBili)
		//工具包
		config.POST("/memory", modifyMemory)
		//应用包
		config.POST("/azure", modifyAzure) //TODO:azure需要再加上对gpt_role的设置接口
		config.POST("/baidu", modifyBaidu)
		config.POST("/openai", modifyOpenai) //TODO:openai需要再加上对gpt_role的设置接口
		config.POST("/pinecone", modifyPinecone)
		config.POST("/xunfei", modifyXfyun)
	}

	//命令相关接口
	command := r.Group("/command")
	{
		command.GET("/read", readText)
		command.GET("/chat", chat)
	}

	//启动进程
	r.GET("/run", run)
	//终止进程（目前没用——和项目结构有关系，得研究一下怎么切开
	r.GET("/stop", stop)

	//调用浏览器打开前端页面
	cmd := exec.Command("cmd", "/c", "start", "http://127.0.0.1:9000/control")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()

	r.Run(":9000")
}
