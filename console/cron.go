package console

import (
	"fmt"
	"project/logic/controll"

	"github.com/robfig/cron/v3"
)

var Conrs *cron.Cron

func init() {
	Conrs := cron.New()
	_, err := Conrs.AddFunc("0 3 * * *", virtual)
	if err != nil {
		fmt.Errorf("定时任务出错: %v", err)
	}
	Conrs.Start()
	defer Conrs.Stop()
	select {}
}
func virtual() {
	fmt.Println("123")
	_, err := controll.RataApiUpdate()
	if err != nil {
		fmt.Println("Error: %v", err)
	}
}
