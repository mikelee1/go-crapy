package message

import (
	"fmt"
	"github.com/qinxin0720/QcloudSms-go/QcloudSms"
	"go-crapy/config"
	"net/http"
	"encoding/json"
)

func SendMsg(msgContent string) error {
	var (
		err error
	)
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	qcloudsms, err := QcloudSms.NewQcloudSms(int(conf.AppID), conf.AppKey)
	if err != nil {
		return err
	}

	err = qcloudsms.SmsSingleSender.SendWithParam(86, conf.Phone, int(conf.RegisteTemp),
		[]string{msgContent, conf.ExpireMinute},
		"", "", "", Callback)
	if err != nil {
		return err
	}
	fmt.Println("send!")
	return nil
}

func Callback(err error, resp *http.Response, resData string) {
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		resp := Response{}
		err = json.Unmarshal([]byte(resData), &resp)
		if err != nil {
			panic(err)
		}
		if resp.Result != 0 || resp.Errmsg != "OK" {
			panic(fmt.Errorf("result is not 0"))
		}
		fmt.Println("response data: ", resData)
	}
}

type Response struct {
	Result int
	Errmsg string
	Ext    string
	Sid    string
	Fee    int
}
