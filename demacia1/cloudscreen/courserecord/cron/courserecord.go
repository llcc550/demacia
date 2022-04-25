package main

import (
	"context"
	"demacia/cloudscreen/courserecord/cron/internal/config"
	"demacia/cloudscreen/courserecord/cron/internal/logic"
	"flag"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/courserecord-cron.yaml", "the config file")

func main() {

	flag.Parse()

	var c config.Config

	conf.MustLoad(*configFile, &c)

	generateLogic := logic.NewRecordGenerateLogic(context.Background(), &c)

	generateLogic.GenerateCourseRecord()
}
