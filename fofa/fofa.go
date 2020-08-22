package fofa

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strconv"
	"strings"
)

var (
	fofaSessionTag = "_fofapro_ars_session"
)

func (f *Fofa) Run() {
	//初始化用来存储结果
	f.ipNodes = make([]ipNode, 0)

	if !f.FofaComma {
		f.runIP()
	} else {
		f.runRaw()
	}

	saveFile := common.GenFileName("fofa")
	logNumber := 0
	for k, n := range f.ipNodes {
		if k == 0 {
			_ = common.SaveFile("端口存活,IP,PORT,网址,标题,中间件,指纹", saveFile)
		}
		line := fmt.Sprintf("%v,%s,%s,%s,%s,%s,%s",
			n.Alive,
			n.Ip,
			n.Port,
			n.Domain,
			n.Title,
			n.Server,
			n.Finger)
		if err := common.SaveFile(line, saveFile); err != nil {
			f.ErrChannel <- common.LogBuild("Fofa",
				fmt.Sprintf("收集信息任务失败，无法保存结果:%s", saveFile), common.FAULT)
			return
		}
		logNumber++
	}
	if logNumber == 0 {
		f.ErrChannel <- common.LogBuild("Fofa",
			fmt.Sprintf("收集信息任务完成，未有结果"), common.INFO)
	} else {
		f.ErrChannel <- common.LogBuild("Fofa",
			fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", logNumber, saveFile), common.INFO)
	}
	return
}

func (f *Fofa) runRaw() {
	f.ErrChannel <- common.LogBuild("Fofa",
		fmt.Sprintf("开始收集信息:查询语法%s,间隔%d秒", f.Ip, f.Interval),
		common.INFO)
}

func (f *Fofa) runIP() {
	f.ErrChannel <- common.LogBuild("Fofa",
		fmt.Sprintf("开始收集信息:IP范围%s,间隔%d秒", f.Ip, f.Interval),
		common.INFO)
	ips := strings.Split(f.Ip, ",")
	for _, ip := range ips {
		base, start, end, err := parseIP(ip)
		if err != nil {
			f.ErrChannel <- err.Error()
			return
		}
		for {
			if start > end {
				break
			}
			f.get(fmt.Sprintf("%s.%d", base, start))
			start++
		}
	}
}

func parseIP(ip string) (base string, start, end int, err error) {
	fromTo := strings.Split(ip, "-")
	ipStart := fromTo[0]
	err = fmt.Errorf(common.LogBuild("Fofa", "起始IP格式错误(eg. 1.1.1.1)", common.FAULT))

	tIp := strings.Split(ipStart, ".")
	if len(tIp) != 4 {
		return
	}
	start, _ = strconv.Atoi(tIp[3])
	end = start
	if len(fromTo) == 2 {
		end, _ = strconv.Atoi(fromTo[1])
	}
	if end == 0 {
		return
	}
	base = fmt.Sprintf("%s.%s.%s", tIp[0], tIp[1], tIp[2])
	err = nil
	return
}
