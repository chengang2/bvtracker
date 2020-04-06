package dao

import (
	"bvtracker/g"
	"fmt"
	"strconv"
)

func Get_report_travel(reportNo string) ([]map[string]interface{}, error) {

	var sql string
	data := make([]map[string]interface{}, 0)
	if len(reportNo) == 11 {
		tab := reportNo[0:2]
		if tab == "59" {
			sql = "SELECT b1.ReportNO, b1.ReportLocation, b1.ReceiveDate, b1.ServiceLevel, b1.Dueday, b2.chargeteam, b2.logintime," +
				" b2.PhotoTime, b2.TAtime, b2.CuttingTime, b2.TestFinishedTime, b2.DCDtime, b2.Drafttime, b2.Checktime, b2.logouttime" +
				" FROM ( select * from bok_if_reportinfo where ReportNO = '" + reportNo + "' ) AS b1 LEFT JOIN reportstate b2 ON b1.Id = b2.ParentID"
		} else if tab == "66" {
			sql = "SELECT b1.ReportNO, b1.ReportLocation, b1.ReceiveDate, b1.ServiceLevel, b1.Dueday, b2.chargeteam, b2.logintime," +
				" b2.PhotoTime, b2.TAtime, b2.CuttingTime, b2.TestFinishedTime, b2.DCDtime, b2.Drafttime, b2.Checktime, b2.logouttime" +
				" FROM ( select * from sh_lab_data.bok_if_reportinfo where ReportNO = '" + reportNo + "' ) AS b1 LEFT JOIN sh_lab_data.reportstate b2 ON b1.Id = b2.ParentID"
		} else if tab == "91" {
			sql = "SELECT b1.ReportNO, b1.ReportLocation, b1.ReceiveDate, b1.ServiceLevel, b1.Dueday, b2.chargeteam, b2.logintime," +
				" b2.PhotoTime, b2.TAtime, b2.CuttingTime, b2.TestFinishedTime, b2.DCDtime, b2.Drafttime, b2.Checktime, b2.logouttime" +
				" FROM ( select * from nj_database.bok_if_reportinfo where ReportNO = '" + reportNo + "' ) AS b1 LEFT JOIN nj_database.reportstate b2 ON b1.Id = b2.ParentID"
		} else {
			return data, g.ArgError("实验室代码不正确，请重新输入")
		}

	} else {
		return data, g.ArgError("订单号不正确，请重新输入")
		//sql += " and PONO='"+reportNo+"'"
	}

	//select * from reportstate a left join (select * from bok_if_reportinfo where ReportNO= '111111') b on  a.ParentID = b.Id

	//fmt.Println("sql==", sql)
	data, err := SelectMysqlData(sql)

	//if len(records) == 0 {
	//	sql = "select * from reportinfor_view where StyleNO='"+reportNo+"'"
	//	fmt.Println("sql==",sql)
	//	records,err = SelectMysqlData(sql)
	//	if len(records) == 0 {
	//		sql = "select * from reportinfor_view where FPUNo='"+reportNo+"'"
	//		fmt.Println("sql==",sql)
	//		records,err = SelectMysqlData(sql)
	//		if len(records) == 0 {
	//			sql = "select * from reportinfor_view where GPUNo='"+reportNo+"'"
	//			fmt.Println("sql==",sql)
	//			records,err = SelectMysqlData(sql)
	//		}
	//	}
	//}

	return data, err
}

func Get_report_status(reportNo, status, begin_time, end_time, phone_num string, pageNo, pageSize int) ([]map[string]interface{}, error) {

	user_sql := "select role_id,email from bv_user where phone_num='" + phone_num + "'"
	records, err := SelectMysqlData(user_sql)
	if err != nil {
		return records, err
	} else if len(records) > 0 {
		var sql string
		role_id := records[0]["role_id"].(string)
		email := ""
		emais := records[0]["email"]
		if emais != nil {
			email = emais.(string)
		}
		if role_id == "0" {
			//IF((b3.`status` = 'Logout'),'已完成','测试中')
			sql = "SELECT b2.ReportNO,b2.ReportLocation, case b3.`status` when 'Logout' then '已完成' when 'Cs Follow' then '客户确认' else '未完成' end AS status," +
				" b2.ReceiveDate,b2.ServiceLevel,b2.Dueday,b2.PONO,b2.StyleNO,b2.FPUNo,b2.GPUNo,b2.Submitter,b2.SubmittingFor,b2.InvoiceRecipient" +
				" FROM ( SELECT * FROM  nj_database.bok_if_report_excontactor WHERE 1=1 "
			//ContactPhoneNo = '" + phone_num + "' OR ContactEmail = '" + email + "'
				if email != ""{
					sql += " and ContactEmail = '" + email + "'"
				}else{
					sql += " and ContactPhoneNo = '" + phone_num + "'"
				}
				sql += " ) AS b1  JOIN ( select * from nj_database.bok_if_reportinfo where  CREATEON >= DATE_SUB(CURDATE(), INTERVAL 3 MONTH) "
			if reportNo != "" {
				sql += " AND ReportNO='" + reportNo + "'"
			}
			if begin_time != "" {
				sql += " AND unix_timestamp(CREATEON) >= unix_timestamp('" + begin_time + "')"
			}
			if end_time != "" {
				sql += " AND unix_timestamp(CREATEON) <= unix_timestamp('" + end_time + "')"
			}
			sql += " ) AS b2 ON b1.ParentID = b2.Id "
			sql += "  JOIN ( select * from nj_database.reportstate where 1=1 "
			if reportNo != "" {
				sql += " AND report='" + reportNo + "'"
			}
			if status != "" {
				if status == "0" {
					sql += " AND status='Logout'"
				}
				if status == "1" {
					sql += " AND status='Cs Follow'"
				}
				if status == "2" {
					sql += " AND status <> 'Cs Follow' AND status <> 'Logout'"
				}
			}
			if begin_time != "" {
				sql += " AND unix_timestamp(datain) >= unix_timestamp('" + begin_time + "')"
			}
			if end_time != "" {
				sql += " AND unix_timestamp(datain) <= unix_timestamp('" + end_time + "')"
			}
			sql += " ) AS b3 ON b3.ParentID = b1.ParentID"
			limitfrom := (pageNo - 1) * pageSize
			sql += " order by b2.id desc limit " + strconv.Itoa(limitfrom) + "," + strconv.Itoa(pageSize)
			//sql = fmt.Sprintf("select * from reportinfor_view  where ReportNO = '%s' AND ReceiveDate >= '%s' AND ReceiveDate <= '%s' AND status = '%s' " +
			//	" and FIND_IN_SET(%s, ContartIDs)",
			//	reportNo,begin_time,end_time,status,phone_num)
		}
		if role_id == "1" {
			sql = ""
		}
		//fmt.Println("sql==", sql)
		records, err = SelectMysqlData(sql)

		return records, err
	} else {
		records := make([]map[string]interface{}, 0)
		return records, nil
	}

}

func Get_report_status_count(reportNo, status, begin_time, end_time, phone_num string) int {

	user_sql := "select role_id,email from bv_user where phone_num=" + phone_num
	records, err := SelectMysqlData(user_sql)
	if err != nil {
		return 0
	} else if len(records) > 0 {
		var sql string
		role_id := records[0]["role_id"].(string)
		email := ""
		emais := records[0]["email"]
		if emais != nil {
			email = emais.(string)
		}
		if role_id == "0" {
			//SELECT b2.ReportNO,b2.ReportLocation, IF((b3.`status` = 'Logout'),'已完成','测试中') AS `status`, b2.ReceiveDate,b2.ServiceLevel,b2.Dueday,b2.PONO,b2.StyleNO,b2.FPUNo,b2.GPUNo,b2.Submitter,b2.SubmittingFor,b2.InvoiceRecipient
			//FROM ( SELECT * FROM  bok_if_report_excontactor WHERE bok_if_report_excontactor.ContactPhoneNo = '15161621716' OR bok_if_report_excontactor.ContactEmail = 'xxx' ) AS b1
			//LEFT JOIN (select * from bok_if_reportinfo where bok_if_reportinfo.ReportNO = '59191120151' and bok_if_reportinfo.ReceiveDate between '2019-04-20' and '2019-04-30') AS b2 ON b1.ParentID = b2.Id
			//LEFT JOIN (select * from reportstate where report = '59191120151' and status <> "Logout" and datain between '2019-04-20' and '2019-04-30') AS b3 ON b3.ParentID = b1.ParentID
			sql = "SELECT b2.ReportNO,b2.ReportLocation, case b3.`status` when 'Logout' then '已完成' when 'Cs Follow' then '客户确认' else '未完成' end AS status," +
				" b2.ReceiveDate,b2.ServiceLevel,b2.Dueday,b2.PONO,b2.StyleNO,b2.FPUNo,b2.GPUNo,b2.Submitter,b2.SubmittingFor,b2.InvoiceRecipient" +
				" FROM ( SELECT * FROM  nj_database.bok_if_report_excontactor WHERE 1=1 "
			//ContactPhoneNo = '" + phone_num + "' OR ContactEmail = '" + email + "'
			if email != ""{
				sql += " and ContactEmail = '" + email + "'"
			}else{
				sql += " and ContactPhoneNo = '" + phone_num + "'"
			}
			sql += " ) AS b1  JOIN ( select * from nj_database.bok_if_reportinfo where CREATEON >= DATE_SUB(CURDATE(), INTERVAL 3 MONTH) "
			if reportNo != "" {
				sql += " AND ReportNO='" + reportNo + "'"
			}
			if begin_time != "" {
				sql += " AND unix_timestamp(CREATEON) >= unix_timestamp('" + begin_time + "')"
			}
			if end_time != "" {
				sql += " AND unix_timestamp(CREATEON) <= unix_timestamp('" + end_time + "')"
			}
			sql += " ) AS b2 ON b1.ParentID = b2.Id "
			sql += "  JOIN ( select * from nj_database.reportstate where 1=1 "
			if reportNo != "" {
				sql += " AND report='" + reportNo + "'"
			}
			if status != "" {
				if status == "0" {
					sql += " AND status='Logout'"
				}
				if status == "1" {
					sql += " AND status='Cs Follow'"
				}
				if status == "2" {
					sql += " AND status <> 'Cs Follow' AND status <> 'Logout'"
				}
			}
			if begin_time != "" {
				sql += " AND unix_timestamp(datain) >= unix_timestamp('" + begin_time + "')"
			}
			if end_time != "" {
				sql += " AND unix_timestamp(datain) <= unix_timestamp('" + end_time + "')"
			}
			sql += " ) AS b3 ON b3.ParentID = b1.ParentID"

		}
		if role_id == "1" {
			sql = ""
		}
		count := CountRecord(sql)
		return count
	} else {
		return 0
	}

}

func Login_user(phone_num, wx_code string) {

	sql := "select login_num from bv_user where 1=1"
	if phone_num != "" {
		sql += " and phone_num=" + phone_num
	}
	if wx_code != "" {
		sql += " and wx_code='" + wx_code + "'"
	}
	records, _ := SelectMysqlData(sql)
	rcord := 0
	if len(records) > 0 {
		rcord, _ = strconv.Atoi(records[0]["login_num"].(string))
	}

	if rcord == 0 {

		insert_sql := fmt.Sprintf("insert into bv_user (phone_num,wx_code,role_id) values ('%s','%s',%d)", phone_num, wx_code, 0)

		InsertMysqlData(insert_sql)
		//if result {
		//	fmt.Println("用户首次登入，插入成功!!!")
		//}
	} else {
		update_sql := fmt.Sprintf("update bv_user set login_num=%d where phone_num='%s'", rcord+1, phone_num)

		UpdateMysqlData(update_sql)
	}

}
