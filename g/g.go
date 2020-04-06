package g

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	ChineseSimplified  = "zh_CN" // 中文简体
	ChineseTraditional = "zh_TW" // 中文繁体
	English            = "en"    // 英文
)

var (
	phone_codes = make(map[string]string, 0)
	phoneLock   = new(sync.RWMutex)
)

func Phone_codes() map[string]string {
	phoneLock.RLock()
	defer phoneLock.RUnlock()
	return phone_codes
}

func CreateCaptcha() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}

func SmsSdkByALiYUN(phoneNum string) (error, string) {
	//生成4位数的随机码
	randomNum := CreateCaptcha()

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", Config().AccessKeyId, Config().AccessKeySecret)
	if err != nil {
		return err, randomNum
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phoneNum
	request.SignName = "必维"
	request.TemplateCode = "SMS_138066377"
	request.TemplateParam = "{\"code\":\"" + randomNum + "\"}"

	response, err := client.SendSms(request)
	if err == nil {
		phone_codes[phoneNum] = randomNum
	}
	msg := response.Message

	return ArgError(msg), randomNum

}

type ArgError string

//重写error方法
func (e ArgError) Error() string {
	return fmt.Sprintf("%s", string(e))
}

func GinSuccessResponse(c *gin.Context, data interface{}, code int, message ...string) {
	res := gin.H{"code": code, "data": data}
	if len(message) == 0 {
		res["msg"] = "ok"
	} else {
		res["msg"] = message[0]
	}
	c.JSON(http.StatusOK, res)
}

func GinErrorResponse(c *gin.Context, code int, message string) {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Println(runtime.FuncForPC(pc).Name(), file, line)
	}
	data := make([]map[string]interface{}, 0)
	res := gin.H{"code": code, "data": data, "message": message}
	c.AbortWithStatusJSON(http.StatusBadRequest, res)
}

func Requests(method string, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest(method, url, body)
	request.Header.Set("Content-Type", "application/json")
	return client.Do(request)
}
