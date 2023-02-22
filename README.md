# LiveTV
将 Youtube 直播作为 IPTV 电视源

## 安装方法

下载对应livetv二进制文件、assert目录、view目录与youtube文件，赋予权限后执行以下命令

install youtube /usr/local/bin/youtube

将assert目录与view目录与livetv放置在同一路径下。

## 使用方法

./livetv -DIR /etc -PORT 9000 -URL 0.0.0.0

参数：

数据库与log路径，默认：/etc

-DIR /etc

监听IP，默认：0.0.0.0

-URL 0.0.0.0

监听端口，默认：9000

-PORT 9000

默认的登入密码是 "password",爲了你的安全请及时修改。

首先你要知道如何在外界访问到你的主机，如果你使用 VPS 或者独立伺服器，可以访问 `http://你的主机ip:9000`

首先你需要在设定区域点击“自动填充”，设定正确的URL。然后点击“储存设定”。

然后就可以添加频道需要代理的勾选代理直播数据流，频道添加成功后就能M3U8档案列的地址进行播放了。

当你使用Kodi之类的播放器，可以考虑使用第一行的M3U档案URL进行播放，会自动生成包含所有频道信息的播放列表。
