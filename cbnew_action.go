package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"strconv"
)

var (
	sckey   string
	barkkey string
	hour    int
	minute  int
)

func main() {
	fmt.Printf("运行中...每天%v时%v分推送\r\n", hour, minute)
	doJob()
	//startScheduler(hour, minute, time.Hour*24)
}

func init() {
	flag.StringVar(&sckey, "sckey", "", "请设置SKEY")
	flag.StringVar(&sckey, "barkkey", "", "请设置barkkey")
	flag.IntVar(&hour, "h", 9, "请设置开始小时")
	flag.IntVar(&minute, "m", 0, "请设置开始分钟")
	flag.Parse()
	if sckey == "" && barkkey == "" {
		panic("请至少设置sckey或barkkey其中一种")
	}
}
func isWorkDay(date string) (isworkday bool) {
	isworkday = true //默认是工作日
	resp, err:= http.Get("http://tool.bitefu.net/jiari/?" + "d=" + date + "&back=json")
	if err != nil {
		fmt.Println("节假日api请求出错：")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("节假日api返回结果出错：")
	}
	//var info Info
	//json.Unmarshal(body, &info)
	/*
	工作日对应结果为 0, 休息日对应结果为 1, 节假日对应的结果为 2
	*/
	var mapResult map[string]interface{}
	if err := json.Unmarshal([]byte(body), &mapResult); err != nil {
		fmt.Println("节假日api返回结果解析出错：" + string(body))
		return isworkday
	}
	if _, err := mapResult[date]; !err {
		fmt.Println("节假日api返回结果解析出错：" + string(body))
		return isworkday
	}

	if mapResult[date] == 1 || mapResult[date] == 2{
		isworkday = false
		fmt.Println("节假日api返回结果成功：" + string(body))
		fmt.Println("isworkday：false")
	}
	return isworkday

}
func pushInfo(title string, text string) {
	fmt.Println(title)
	fmt.Println(text)
	//fmt.Println(url.QueryEscape(text))
	// 使用Server酱推送
	if sckey != "" {
		resp, _ := http.Get("https://sc.ftqq.com/" + sckey + ".send?text=" + url.QueryEscape(title) + "&desp=" + url.QueryEscape(text))
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Server酱推送结果" + string(body))
	}

	// 使用Bark推送
	if barkkey != "" {
		resp, _ := http.Get("https://api.day.app/" + barkkey + "/" + title + "/" + text)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("BARK推送结果" + string(body))
	}
}

func doJob() {
	timeStr:=time.Now().Format("20060102")
	isworkday := isWorkDay(timeStr)
	applyList, listList := getTodayCbInfo()
	fmt.Println("=====isworkday?" + strconv.FormatBool(isworkday))
	//fmt.Sprintf("%t", isworkday)
	
	if !isworkday {
		return
	}
	
	if len(applyList) == 0 {
		pushInfo("今日无可打新债", "")
	} else {
		var text string
		for _, apply := range applyList {
			apply = "- " + apply //markdown
			text += apply + "\r\n"
		}
		pushInfo("今日可打新债", text)
	}

	if len(listList) == 0 {
		pushInfo("今日无上市债券", "")
	} else {
		var text string
		for _, list := range listList {
			list = "- " + list //markdown
			text += list + "\r\n"
		}
		pushInfo("今日上市债券", text)

	}
}

func startScheduler(hour int, minute int, duration time.Duration) {
	now := time.Now()
	first := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 00, 0, now.Location())
	if first.Before(now) {
		first = first.Add(time.Hour * 24)
	}
	starter := time.NewTicker(first.Sub(now))
	<-starter.C
	starter.Stop()

	scheduler := time.NewTicker(duration)
	for {
		doJob()
		<-scheduler.C
	}
	defer scheduler.Stop()
}

func getTodayCbInfo() (applyList []string, listList []string) {
	uri := "https://www.jisilu.cn/data/cbnew/pre_list/"

	client := &http.Client{}

	req, _ := http.NewRequest("POST", uri, strings.NewReader("rp=22"))

	req.Header.Set("Origin", "https://www.jisilu.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://www.jisilu.cn/data/cbnew/")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	var info Info
	json.Unmarshal(body, &info)

	todayStr := time.Now().Format("2006-01-02")
	for _, row := range info.Rows {
		if row.Cell.ApplyDate == todayStr {
			applyList = append(applyList, conv(row.Cell))
		}
		if row.Cell.ListDate == todayStr {
			listList = append(listList, conv(row.Cell))
		}
	}

	return applyList, listList
}

func conv(cell Cell) string {
	return "名称:" + cell.Name + ",申购日期:" + cell.ApplyDate + ",上市日期:" + cell.ListDate
}

type Info struct {
	Page int   `json:"page"`
	Rows []Row `json:"rows"`
}

type Row struct {
	Id   string `json:"id"`
	Cell Cell   `json:"cell"`
}

type Cell struct {
	Name      string `json:"bond_nm"`
	ApplyDate string `json:"apply_date"`
	ListDate  string `json:"list_date"`
	Advise    string `json:"jsl_advise_text"`
	BondId    string `json:"bond_id"`
	DrawRate  string `json:"lucky_draw_rt"`
	Rating    string `json:"rating_cd"`
}
