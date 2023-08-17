package main

import (
	"log"
	"os"
	"text/template"
)

const gptTmpl = `gpt:
  token: {{.TokenGPT}}
`

const yandexCloudTmpl = `yandexCloud:
  oAuthToken: {{.OAuth}}
  folderId: {{.FolderID}}
`

const serverTmpl = `server:
  port: {{.Port}}
  host: "{{.Host}}"
`
const (
	path            = "configs"
	cloud           = path + "/api"
	server          = path + "/server"
	gptFileName     = cloud + "/gpt.yml"
	yaCloudFileName = cloud + "/yandex_cloud.yml"
	servFileName    = server + "/server.yml"
)

func main() {
	serverT := template.Must(template.New("server").Parse(serverTmpl))
	gptT := template.Must(template.New("gpt").Parse(gptTmpl))
	ycT := template.Must(template.New("cloud").Parse(yandexCloudTmpl))

	err := os.Mkdir(path, 0777)
	if err != nil {
		log.Println(err)
	}

	err = os.Mkdir(cloud, 0777)
	if err != nil {
		log.Println(err)
	}
	err = os.Mkdir(server, 0777)
	if err != nil {
		log.Println(err)
	}

	gptConf, err := os.Create(gptFileName)
	if err != nil {
		log.Fatal(err)
	}
	err = gptT.Execute(gptConf, map[string]string{
		"TokenGPT": os.Getenv("TOKEN_GPT"),
	})
	if err != nil {
		log.Fatal(err)
	}

	yaCloudConf, err := os.Create(yaCloudFileName)
	if err != nil {
		log.Fatal(err)
	}
	err = ycT.Execute(yaCloudConf, map[string]string{
		"OAuth":    os.Getenv("YANDEX_CLOUD_OAuth"),
		"FolderID": os.Getenv("YANDEX_CLOUD_FOLDER_ID"),
	})
	if err != nil {
		log.Fatal(err)
	}

	serverConf, err := os.Create(servFileName)
	if err != nil {
		log.Fatal(err)
	}
	err = serverT.Execute(serverConf, map[string]string{
		"Host": os.Getenv("HOST"),
		"Port": os.Getenv("PORT"),
	})
	if err != nil {
		log.Fatal(err)
	}

}
