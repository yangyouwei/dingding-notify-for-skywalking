package main

import (
        "bytes"
        "encoding/json"
        "text/template"
        "flag"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "strings"
)

var (
        //服务器端口号
        serverPort string
        //钉钉告警url
        dingdingUrl string
)

//钉钉告警模板
const templatealarm = `{
     "msgtype": "text",
     "text": {
         "content": "%s"
     },
     "at": {
         "isAtAll": false
     }
 }`

type Mgs struct {
        ScopeId      int    `json:"scopeId"`
        Name         string `json:"name"`
        Id0          int    `json:"id0"`
        Id1          int    `json:"id1"`
        AlarmMessage string `json:"alarmMessage"`
        StartTime    int64    `json:"startTime"`
}

const tmpl  = `ScopeId = {{.ScopeId}}
Name = {{.Name}}
Id0 = {{.Id0}}
Id1 = {{.Id1}}
StartTime = {{.StartTime}}
AlarmMessage = {{.AlarmMessage}}`

func init() {

        flag.StringVar(&serverPort, "p", "9201", "server port")
        flag.StringVar(&dingdingUrl, "u", "", "alarm webhook")
}

func sendMsg(w http.ResponseWriter, r *http.Request) {
defer fmt.Fprintf(w, "ok\n")
        bodys, _ := ioutil.ReadAll(r.Body)
        messages := string(bodys)
        nochline := strings.Replace(messages, "\n", "", -1)
        nochline = strings.Replace(messages, "\r", "", -1)
        jsonms := []byte(nochline)
        msslice := []Mgs{}

        err:=json.Unmarshal(jsonms,&msslice)
        if err != nil {
                fmt.Println(err)
        }
        a := printtmpl(msslice)
        fmt.Println(a)

        msg := strings.NewReader(fmt.Sprintf(templatealarm, a))
        //发送钉钉消息
        res, err := http.Post(dingdingUrl, "application/json", msg)
        fmt.Println(msg)
        if err != nil || res.StatusCode != 200 {
                log.Print("send alram msg to dingding fatal.", err, res.StatusCode)
        }
}

func main() {
        flag.Parse()
        if dingdingUrl == "" {
                log.Fatal("dingdingUrl cannot be empty usage -h get help. ")
        }

        http.HandleFunc("/alarm", sendMsg)

        //启动web服务器
        if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil); err != nil {
                log.Fatal("server start fatal.", err)
        }

}

func printtmpl(s []Mgs) string {
        l := len(s)
        var alarms string
        for i := 0;i<l ;i++ {
                a := s[i]
                tmpl, err := template.New("alarm").Parse(tmpl)  //建立一个模板
                if err != nil {
                        panic(err)
                }
                buf := new(bytes.Buffer)
                err = tmpl.Execute(buf, a)
                if err != nil {
                        panic(err)
                }
                alarms = alarms + fmt.Sprint( buf) + "\n\n"
        }
        return alarms
}
