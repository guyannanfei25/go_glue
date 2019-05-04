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

	if conf.Get("db_conf").Get("enable").MustBool(false) {
		dsn := conf.Get("db_conf").Get("dsn").MustString("")
		if dsn == "" {
			return fmt.Errorf("db conf dsn is empty")
		}

		user := conf.Get("db_conf").Get("user").MustString("")
		if user == "" {
			return fmt.Errorf("db conf user is empty")
		}

		password := conf.Get("db_conf").Get("password").MustString("")
		if password == "" {
			return fmt.Errorf("db conf password is empty")
		}

		host := conf.Get("db_conf").Get("host").MustString("")
		if host == "" {
			return fmt.Errorf("db conf host is empty")
		}

		port := conf.Get("db_conf").Get("port").MustInt()

		dbName := conf.Get("db_conf").Get("db_name").MustString("")
		if dbName == "" {
			return fmt.Errorf("db conf db_name is empty")
		}

		fDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", user, password, host, port, dbName, dsn)
		if err := common.InitGorm(fDsn); err != nil {
			return err
		}
	}

	if conf.Get("redis_conf").Get("enable").MustBool(false) {
		addr := conf.Get("redis_conf").Get("addr").MustString("")
		if addr == "" {
			return fmt.Errorf("redis_conf addr is empty")
		}

		passwd := conf.Get("redis_conf").Get("passwd").MustString("")

		if err := common.InitRedis(addr, passwd); err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	return common.Close()
}
