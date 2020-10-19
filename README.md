# Convertible bond news
打新债信息推送机器人

# 推送方式
## [Server酱（微信）](http://sc.ftqq.com/3.version)
登录到Server酱，绑定微信之后，获取到sckey即可，样式如下  
![ft](https://raw.githubusercontent.com/Cyronlee/cbnew-go/master/imgs/ft.png)
## [BARK（iOS）](https://github.com/Finb/Bark)
在AppStore下载BARK，获取到barkkey（机器码）即可，样式如下  
![bark](https://raw.githubusercontent.com/Cyronlee/cbnew-go/master/imgs/bark.png)

# 使用
1.github actions部署(推荐)
- fork本项目后，在setting里配置推送的秘钥，即可。
- 可以通过修改.github/workflows下面的yaml文件来修改定时推送时间，注意，时间为UTC时间。

2.下载源码自己编译、部署
```bash
# sckey 和 barkkey 至少设置其中一个
cbnew.exe -sckey=XXXX -barkkey=XXXX

# 手动设置推送时间为 8:50，默认每天 9:00 分推送一次
cbnew.exe -sckey=XXXX -h=8 -m=50
```

# 数据来源
- [集思录](https://www.jisilu.cn/data/cbnew/#pre)
- [东方财富](http://data.eastmoney.com/kzz/default.html)
- [免费节假日API](http://tool.bitefu.net/jiari/)

# TODO
- [x] ~~优化推送消息的格式~~
- [x] ~~加入工作日的判断~~
- [x] ~~增加github actions部署~~

# 参考
[V2EX：cbnew-python](https://github.com/crazygit/cbnew)
