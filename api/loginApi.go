package api

import (
	"bvtracker/dao"
	"bvtracker/g"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strings"

	"io/ioutil"
	"net/http"
	"net/url"
)

func Get_phone_code(c *gin.Context) {

	var jsMap map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &jsMap)
	phone_num := jsMap["phone_num"].(string)

	code := 0
	msg := "success"
	data := make([]map[string]interface{}, 0)

	//获取手机验证码
	err, phone_code := g.SmsSdkByALiYUN(phone_num)
	if err.Error() == "OK" {
		entry := make(map[string]interface{}, 0)
		entry["phone_num"] = phone_num
		entry["phone_code"] = phone_code

		data = append(data, entry)
	} else {
		code = 1
		msg = err.Error()
		if strings.Contains(msg, "触发分钟级流控Permits:1") {
			msg = "操作太频繁，一个手机号一分钟只能登入一次"
		}
		if strings.Contains(msg, "触发小时级流控Permits:5") {
			msg = "操作太频繁，一个手机号一小时只能登入五次"
		}
		if strings.Contains(msg, "触发天级流控Permits:10") {
			msg = "操作太频繁，一个手机号一天只能登入十次"
		}
	}

	g.GinSuccessResponse(c, data, code, msg)

}

func Login(c *gin.Context) {

	var jsMap map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &jsMap)

	phone_num := jsMap["phone_num"].(string)
	phone_code := jsMap["phone_code"].(string)

	true_phone_code := g.Phone_codes()[phone_num]

	code := 0
	msg := "success"
	data := make([]map[string]interface{}, 0)

	if len(phone_code) == 3 {
		phone_code = "0" + phone_code
	}
	//fmt.Println("phone_code==", phone_code, "  true_phone_code==", true_phone_code)

	if true_phone_code != phone_code {
		code = 1
		msg = "手机验证码不正确，请重新输入"
	} else {

		dao.Login_user(phone_num, "")

		entry := make(map[string]interface{}, 0)
		entry["token"] = "bvtracker"
		entry["phone_num"] = phone_num
		data = append(data, entry)

	}

	g.GinSuccessResponse(c, data, code, msg)

}

func Wetchat_code(c *gin.Context) {

	W := WeCharClient{
		Appid:       g.Config().Appid,
		Secret:      g.Config().Appsecret,
		RedirectUri: g.Config().Callback,
		Scope:       "snsapi_login",
	}
	result := W.AuthCodeUrl("state")
	fmt.Println("result==", result)

}

func Exchange(c *gin.Context) {

	code := c.Query("code") //查询请求URL后面的参数
	//fmt.Println("code_get==", code)

	W := WeCharClient{
		Appid:       g.Config().Appid,
		Secret:      g.Config().Appsecret,
		RedirectUri: g.Config().Callback,
		Scope:       "snsapi_login",
	}

	access_token, err := W.Exchanges(code)
	if err != nil {
		fmt.Println("err==", err)
	}
	log.Println("access_token==", access_token)
	log.Println("Unionid==", access_token.Unionid)
	log.Println("openid==", access_token.OpenId)
	user_info, err := W.GetUserInfo(access_token.AccessToken, access_token.OpenId, g.ChineseSimplified)
	fmt.Println("user_info==", user_info)

}

// 获取认证二维码url
func (w *WeCharClient) AuthCodeUrl(state string) string {
	return fmt.Sprintf(
		"https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		w.Appid, url.QueryEscape(w.RedirectUri), w.Scope, state)
}

// 获取access token
func (w *WeCharClient) Exchanges(code string) (AccessToken, error) {

	reUrl := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		w.Appid, w.Secret, code)
	log.Println("reUrl==", reUrl)
	if response, err := g.Requests("GET", reUrl, nil); err == nil && response.StatusCode == http.StatusOK {

		body := response.Body
		defer body.Close()

		if bodyByte, err := ioutil.ReadAll(body); err == nil {

			var result AccessToken

			log.Println("bodyByte==", string(bodyByte))

			if err := json.Unmarshal(bodyByte, &result); err == nil {
				fmt.Println("result==", result)
				return result, nil

			}

		}
	}
	return AccessToken{}, errors.New("get access token fail")
}

// 刷新access token
func (w *WeCharClient) ReGetAccessToken(refreshToken string) (AccessToken, error) {
	reUrl := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s#wechat_redirect",
		w.Appid, refreshToken)
	if response, err := g.Requests("GET", reUrl, nil); err == nil && response.StatusCode == http.StatusOK {

		body := response.Body
		defer body.Close()
		if bodyByte, err := ioutil.ReadAll(body); err == nil {
			var result AccessToken
			if err := json.Unmarshal(bodyByte, &result); err == nil {
				return result, nil
			}
		}
	}
	return AccessToken{}, errors.New("get access token fail")
}

// 获取用户信息
func (w *WeCharClient) GetUserInfo(accessToken string, openId string, lang ...string) (UserInfo, error) {

	reUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s#wechat_redirect", accessToken, openId)
	if len(lang) > 0 {
		reUrl += fmt.Sprintf("&%s", lang[0])
	}
	if response, err := g.Requests("GET", reUrl, nil); err == nil && response.StatusCode == http.StatusOK {

		body := response.Body
		defer body.Close()
		if bodyByte, err := ioutil.ReadAll(body); err == nil {
			var result UserInfo
			if err := json.Unmarshal(bodyByte, &result); err == nil {
				return result, nil
			}
		}
	}
	return UserInfo{}, errors.New("get user info fail")
}
