package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/skip2/go-qrcode"
	"net/http"
	"os"
	"strings"
	"time"
	"tiny_tool/layout"
	"tiny_tool/log"
	"tiny_tool/tool"
)

const (
	PORT = "1323"
)

var UpdateKeys func([]interface{})

func ReadMsg(id string, ws *websocket.Conn) {
	for {
		//msg := ""
		_, msg, err := ws.ReadMessage()
		//err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			offline(ws)
			break
		}
		m := make(map[string]interface{})
		if err := json.Unmarshal(msg, &m); err == nil {
			if m["type"] == "log" {
				log.Console(m["level"].(float64), time.Now().Format("01-02 15:04:05"), "["+id+"]", m["tag"].(string)+":"+m["msg"].(string))
			} else if m["type"] == "func" {
				if UpdateKeys != nil {
					UpdateKeys(m["keys"].([]interface{}))
				}
			} else if m["type"] == "NativeLayout" {
				layout.PrintLayout(string(msg))
			}else if m["type"]=="heartBeat" {
				updateHeartBeat(id)
			}
		}
	}
}


func GetApkDownloadUrl() string {
	ip := tool.GetIP()
	return "http://" + ip + ":1323/apk/build/outputs/apk/debug/app-debug.apk"
}

func GetWsPath() string {
	return "ws://" + tool.GetIP() + ":" + PORT + "/ws"
}

func StartServer() {
	e := echo.New()
	//e.Logger.SetOutput(log.EchoLogger{})
	e.Logger.SetOutput(e.Logger.Output())
	e.Use(middleware.Recover())
	e.GET("/ws", Connect)
	e.GET("/", func(context echo.Context) error {
		html = strings.ReplaceAll(html, "$SRC", "http://127.0.0.1:"+PORT+"/qrCode")
		return context.HTML(200, html)
	})
	e.GET("/qrCode", getQrCode)

	androidDir := tool.GetCurrentPath() + "/android"
	if len(os.Args) > 2 {
		androidDir = os.Args[2]
	}
	e.Static("/apk", androidDir)
	go heartBeat()
	err := e.Start(":" + PORT)
	panic(err)
}
func getQrCode(context echo.Context) error {
	png, _ := qrcode.Encode(GetApkDownloadUrl(), qrcode.Medium, 256)
	return context.Stream(200, "image/png", bytes.NewBuffer(png))
}
func heartBeat() {
	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		unix := time.Now().Unix()
		data := make(map[string]interface{})
		data["type"] = "heartBeat"
		data["time"] = unix
		PublishMsg(data)
	}
}

func Connect(context echo.Context) error {
	request := context.Request()
	response := context.Response()
	AndroidId := request.RemoteAddr
	upgrader := websocket.Upgrader{}

	upgrader.HandshakeTimeout=time.Millisecond*100
	upgrader.Error= func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.V(status," Error")
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {

		return true
	}
	if ws, err := upgrader.Upgrade(response, request, nil);err ==nil{
		online(AndroidId, ws)
		ReadMsg(AndroidId, ws)
	}else{
		log.V(err)
	}
	return nil
}

func PublishMsg(data interface{}) {
	if marshal, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		publishMsg(marshal)
	}
}

var html = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="author" content="mengkun">
    <meta name="renderer" content="webkit">
    <meta http-equiv="Cache-Control" content="no-siteapp">
    <title>扫码安装调试包</title>

    <script>
        /**
         * 本多合一收款码作者 @megnkun(https://mkblog.cn)，转载请注明！
         * Github： https://github.com/mengkunsoft/OneQRCode
         *
         * 以下内容请用专门的 Html 编辑器（如 VS Code）进行编辑，切勿直接用记事本编辑，否则可能出现中文乱码！
         */
    </script>

    <style>
        * {
            margin: 0;
            padding: 0;
            font-family: Microsoft yahei;
        }

        body {
            background-color: #fff;
        }

        .code-item {
            width: 100%;
            max-width: 400px;
            margin: 0 auto;
            padding-bottom: 1px;
            display: none;
            background-color: darkorange;
        }

        .code-title {
            text-align: center;
            color: red;
            line-height: 50px;
            font-size: 20px;
            height: 50px;
             background-color: lightgray;
            background-position: center;
            background-repeat: no-repeat;
        }

        .code-area {
            text-align: center;
        }

        .code-area img {
            margin: 50px auto;
            width: 60%;
            min-width: 100px;
            background: #fff;
            padding: 10px;
            border-radius: 5px;
        }

        .code-footer {
            height: 80px;
            text-align: center;
            background-color: lightgray;
            color: #666;
            line-height: 80px;
            font-size: 20px;
        }

        #code-all > .code-title {
            background-image: url("https://bkimg.cdn.bcebos.com/pic/a5c27d1ed21b0ef4d58253f3d5c451da81cb3e31?x-bce-process=image/watermark,image_d2F0ZXIvYmFpa2UxMTY=,g_7,xp_5,yp_5/format,f_auto");
        }
    </style>
</head>
<body>
<!-- 万能收款码展示区域 -->
<div class="code-item" id="code-all">
    <div class="code-title">
        <span>扫码设备必须与本机在同一局域网</span>
    </div>
    <div class="code-area">
        <img id="page-url"
             src="https://bkimg.cdn.bcebos.com/pic/a5c27d1ed21b0ef4d58253f3d5c451da81cb3e31?x-bce-process=image/watermark,image_d2F0ZXIvYmFpa2UxMTY=,g_7,xp_5,yp_5/format,f_auto">
    </div>
    <div class="code-footer">扫码下载调试包</div>
</div>


<script>
    /* 其它，显示“万能码” */
    document.getElementById("page-url").src = "$SRC";
    document.getElementById("code-all").style.display = "block";

</script>

</body>
</html>
`
