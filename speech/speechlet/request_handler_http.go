package speechlet

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (rh *RequestHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ERROR] RequestHandler(AppId:%s) read request body content error: %s",
			rh.AppId, err)
		respStr := "internal error"
		fmt.Fprint(rw, respStr)
		return
	}
	log.Println("INFO] request >> ", string(reqBytes))
	respBytes, err := rh.HandleCall(reqBytes)
	if err != nil {
		log.Printf("ERROR] RequestHandler(AppId:%s) HandleCall error: %s", rh.AppId, err)
		fmt.Fprint(rw, err.Error())
		return
	}
	_, err = rw.Write(respBytes)
	if err != nil {
		log.Printf("ERROR] RequestHandler(AppId:%s) write http.ResponseWrite error: %s",
			rh.AppId, err)
	}
}
