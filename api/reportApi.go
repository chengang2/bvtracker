package api

import (
	"bvtracker/dao"
	"bvtracker/g"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func Report_data(c *gin.Context) {

	var jsMap map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &jsMap)

	reportNo := jsMap["reportNo"].(string)

	code := 0
	msg := "success"
	data := make([]map[string]interface{}, 0)
	var err error
	data, err = dao.Get_report_travel(reportNo)
	if err != nil {
		code = 1
		msg = err.Error()
	}
	g.GinSuccessResponse(c, data, code, msg)

}

func Report_manager(c *gin.Context) {

	var jsMap map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &jsMap)
	begin_time := jsMap["begin_time"]
	end_time := jsMap["end_time"]
	status := jsMap["status"]
	reportNo := jsMap["reportNo"]
	phone_num := jsMap["phone_num"].(string)
	pageNo := int(jsMap["indexin"].(float64))
	pageSize := int(jsMap["sizein"].(float64))
	code := 0
	msg := "success"
	data := make([]map[string]interface{}, 0)
	var err error
	var btime, etime, stat, reportno string
	btime, etime, stat, reportno = "", "", "", ""
	if begin_time != nil {
		btime = begin_time.(string)
	}
	if end_time != nil {
		etime = end_time.(string)
	}
	if status != nil {
		stat = status.(string)
	}
	if reportNo != nil {
		reportno = reportNo.(string)
	}

	data, err = dao.Get_report_status(reportno, stat, btime, etime, phone_num, pageNo, pageSize)
	entry := make(map[string]interface{})
	entry["mysql_total_num"] = dao.Get_report_status_count(reportno, stat, btime, etime, phone_num)
	entry["items"] = data

	if err != nil {
		code = 1
		msg = err.Error()
	}
	g.GinSuccessResponse(c, entry, code, msg)
}
