package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/model"
	sp "roobo.com/rosai-skills-kit-sdk-for-go/speech/speechlet"
	"roobo.com/sailor/glog"
	"roobo.com/sailor/util"
)

func init() {
	flag.Parse()
	if err := util.InitConf("./conf/app.json"); err != nil {
		log.Fatalf("load app.json error: %s", err)
	}
	ConfLog()
}

func main() {
	dm, err := getDialogModel()
	if err != nil {
		glog.Fatal(err)
	}
	rh := sp.RequestHandler{
		AppId:       "rosai1.ask.skill.tidepooler.12345",
		Speechlet:   &TidePooler{},
		DialogModel: dm,
	}
	http.Handle("/tidepooler", &rh)
	ip, port := getServerIP(), getServerPort()
	glog.Infof("start helloworld server on: %s:%d", ip, port)
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", ip, port),
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 3000 * time.Millisecond,
		IdleTimeout:  90 * time.Second,
	}
	glog.Fatal(srv.ListenAndServe())
}

func getDialogModel() (*model.DialogModel, error) {
	data, err := ioutil.ReadFile("./conf/dialog.json")
	if err != nil {
		return nil, err
	}
	dm := new(model.DialogModel)
	if err := json.Unmarshal(data, dm); err != nil {
		return nil, err
	}
	return dm, nil
}

func getServerIP() string {
	v, err := util.GetCfgVal("localhost", "server", "host")
	if err != nil {
		glog.Fatal("get server host error:", err)
	}
	if vv, ok := v.(string); ok {
		if net.ParseIP(vv) == nil {
			glog.Fatal("server ip set incorrect")
		}
		return vv
	}
	glog.Fatal("server ip configure error")
	return "0.0.0.0" // cannot reach here
}

func getServerPort() int {
	v, err := util.GetCfgVal(10000, "server", "port")
	if err != nil {
		glog.Fatal("get server prot error:", err)
	}
	if vv, ok := v.(int); ok {
		if vv <= 1023 || vv > 65535 {
			glog.Fatal("server port should greater than 1023 and less than 65535")
		}
		return vv
	}
	glog.Fatal("server port configure error")
	return 0 // cannot reach here
}

func ConfLog() {
	s, err1 := util.GetCfgVal("./log", "log", "log_dir")
	t, err2 := util.GetCfgVal("INFO", "log", "stderrthreshold")
	v, err3 := util.GetCfgVal(0, "log", "v")
	if err1 != nil || err2 != nil || err3 != nil {
		glog.Warningf("GetCfgVal log conf error: %v,%v,%v,%v", err1, err2, err3)
		s, t, v = "./log", "INFO", 0
	}
	flag.Lookup("log_dir").Value.Set(s.(string))
	flag.Lookup("stderrthreshold").Value.Set(t.(string))
	flag.Lookup("v").Value.Set(strconv.Itoa(v.(int)))
}
