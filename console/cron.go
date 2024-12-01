package console

import (
	"fmt"
	"project/logic/service"

	"github.com/robfig/cron/v3"
)

var Conrs *cron.Cron

func init() {
	Conrs = cron.New()
	_, err := Conrs.AddFunc("00 03 * * *", virtual)
	if err != nil {
		fmt.Printf("定时任务出错: %v\n", err)
		return
	}
	Conrs.Start()
	// defer Conrs.Stop()
	// select {}
}

func virtual() {
	fmt.Println("正在更新货币信息")
	_, err := service.ServiceGroupApp.RateServiceGroup.RataApiUpdate()
	if err != nil {
		fmt.Println("Error: %v", err)
	}
	fmt.Println("更新已完成")
}
