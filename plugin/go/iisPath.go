package goplugin

import (
	"net/http"
	"strings"

	"github.com/opensec-cn/kunpeng/plugin"
	"github.com/opensec-cn/kunpeng/util"
)

type iisPath struct {
	info   plugin.Plugin
	result []plugin.Plugin
}

func init() {
	plugin.Regist("iis", &iisPath{})
}
func (d *iisPath) Init() plugin.Plugin {
	d.info = plugin.Plugin{
		Name:    "IIS 物理路径泄露",
		Remarks: "通过访问一个不存在的文件或者目录，得到web物理路径",
		Level:   4,
		Type:    "INFO",
		Author:  "wolf",
		References: plugin.References{
			KPID: "KP-0016",
		},
	}
	return d.info
}
func (d *iisPath) GetResult() []plugin.Plugin {
	var result = d.result
	d.result = []plugin.Plugin{}
	return result
}
func (d *iisPath) Check(URL string, meta plugin.TaskMeta) bool {
	request400, err := http.NewRequest("GET", URL+"/404-test.asp", nil)
	if err != nil {
		return false
	}
	resp, err := util.RequestDo(request400, false)
	if err != nil {
		return false
	}
	if resp.Other.StatusCode == 404 && strings.Contains(resp.ResponseRaw, "0x80070002") {
		result := d.info
		result.Response = resp.ResponseRaw
		result.Request = resp.RequestRaw
		d.result = append(d.result, result)
		return true
	}
	return false
}
