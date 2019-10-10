package web

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"golang.org/x/net/websocket"
	"html/template"
	"net/http"
)

var distBox *rice.Box

func StartServer(port uint, access bool) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port %d", port)
	}

	distBox = rice.MustFindBox("dist") // go.rice 文件盒子
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(distBox.HTTPBox())))

	http.HandleFunc("/", RootHandle)
	http.HandleFunc("/index.html", middleware(indexPage))
	http.HandleFunc("/containers", middleware(ContainerHandle))

	http.Handle("/ws", websocket.Handler(WSHandler))
	if access {
		return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
	fmt.Println("现在只监听localhost，请注意")
	return http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	tmpl := boxTmplParse("index", "index.html")
	_ = tmpl.Execute(w, nil)
}

// boxTmplParse ricebox 载入文件内容, 并进行模板解析
func boxTmplParse(name string, fileNames ...string) (tmpl *template.Template) {
	tmpl = template.New(name)
	for k := range fileNames {
		_, _ = tmpl.Parse(distBox.MustString(fileNames[k]))
	}
	return
}
