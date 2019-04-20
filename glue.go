package glue

import (
	"fmt"

	sj "github.com/guyannanfei25/go-simplejson"
	"github.com/guyannanfei25/go_common"
)

func Init(conf *sj.Json) error {
	// init logger
	logDir := conf.Get("log_info").Get("dir").MustString("/tmp")
	logName := conf.Get("log_info").Get("name").MustString("log.info")
	logLevel := conf.Get("log_info").Get("level").MustInt(1)

	if err := common.InitLogger(logDir, logName, logLevel); err != nil {
		return err
	}

	// init gc
	maxMem := conf.Get("gc_info").Get("max_mem_m").MustInt(2048)
	checkInter := conf.Get("gc_info").Get("check_interval_s").MustInt(20)

	go common.IntervalGC(maxMem, checkInter)

	// init pid
	pidName := conf.Get("pid_info").Get("file").MustString("/tmp/glue.pid")

	if err := common.InitPidFile(pidName); err != nil {
		return err
	}

	maxProc := conf.Get("proc_info").Get("max_proc").MustInt(-1)

	common.InitRunProcs(maxProc)

	dsn := conf.Get("db_conf").Get("dsn").MustString("")
	if dsn == "" {
		return fmt.Errorf("db conf dsn is empty")
	}
	if err := common.InitGorm(dsn); err != nil {
		return err
	}

	return nil
}

func Close() error {
	return common.Close()
}
