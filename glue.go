package glue

import (
    "github.com/guyannanfei25/go_common"
    sj "github.com/guyannanfei25/go-simplejson"
)

func Init(conf *sj.Json) error {
    // init logger
    logDir     := conf.Get("log_info").Get("dir").MustString("/tmp")
    logName    := conf.Get("log_info").Get("name").MustString("log.info")
    logLevel   := conf.Get("log_info").Get("level").MustInt(1)

    if err := common.InitLogger(logDir, logName, logLevel); err != nil {
        return err
    }

    // init gc
    maxMem     := conf.Get("gc_info").Get("max_mem_m").MustInt(2048)
    checkInter := conf.Get("gc_info").Get("check_interval_s").MustInt(20)
    
    go common.IntervalGC(maxMem, checkInter)

    // init pid
    pidName    := conf.Get("pid_info").Get("file").MustString("/tmp/glue.pid")

    if err := common.InitPidFile(pidName); err != nil {
        return err
    }

    maxProc    := conf.Get("proc_info").Get("max_proc").MustInt(-1)

    common.InitRunProcs(maxProc)

    return nil
}
