package goplugin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/opensec-cn/kunpeng/plugin"
	"github.com/opensec-cn/kunpeng/util"
)

type tomcatWeakPass struct {
	info   plugin.Plugin
	result []plugin.Plugin
}

func init() {
	plugin.Regist("tomcat", &tomcatWeakPass{})
}
func (d *tomcatWeakPass) Init() plugin.Plugin {
	d.info = plugin.Plugin{
		Name:    "Apache Tomcat 弱口令",
		Remarks: "攻击者通过此漏洞可以登陆管理控制台，通过部署功能可直接获取服务器权限。",
		Level:   0,
		Type:    "WEAKPWD",
		Author:  "wolf",
		References: plugin.References{
			URL:  "https://github.com/vulhub/vulhub/tree/master/tomcat/tomcat8",
			KPID: "KP-0020",
		},
	}
	return d.info
}
func (d *tomcatWeakPass) GetResult() []plugin.Plugin {
	var result = d.result
	d.result = []plugin.Plugin{}
	return result
}
func (d *tomcatWeakPass) Check(URL string, meta plugin.TaskMeta) bool {
	userList := []string{
		"admin", "tomcat", "apache", "root", "manager",
	}
	for _, user := range userList {
		for _, pass := range meta.PassList {
			pass = strings.Replace(pass, "{user}", user, -1)
			request, err := http.NewRequest("GET", URL+"/manager/html", nil)
			request.SetBasicAuth(user, pass)
			resp, err := util.RequestDo(request, true)
			if err != nil {
				return false
			}
			if resp.Other.StatusCode == 404 {
				return false
			}
			if resp.Other.StatusCode == 200 {
				if strings.Contains(resp.ResponseRaw, "/manager/html/reload") || strings.Contains(resp.ResponseRaw, "Tomcat Web Application Manager") {
					result := d.info
					result.Response = resp.ResponseRaw
					result.Request = resp.RequestRaw
					result.Remarks = fmt.Sprintf("弱口令：%s,%s,%s", user, pass, result.Remarks)
					d.result = append(d.result, result)
					return true
				}
				//200 又没关键字的可能不是tomcat
				return false
			}
		}
	}
	return false
}
