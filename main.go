package main

import (
    "net/http"
    "fmt"
    "io"
    "os"
    "encoding/json"
    "path"
    "time"
    "strconv"
    "file-uploader/util"
    "file-uploader/db"
)

const RootPath = "/opt/fshome/"
const Domain = "https://www.yinjk.cn/"

type Response struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Err  string `json:"err"`
}

func main() {
    http.HandleFunc("/upload", upload)
    //http.ListenAndServeTLS(":8082","ssl/1_www.yinjk.cn_bundle.crt", "ssl/2_www.yinjk.cn.key", nil)
    http.ListenAndServe(":8082", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域（跨域请求）
    r.ParseMultipartForm(32 << 20)
    file, handler, err := r.FormFile("file")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()
    // 获取上传文件的md5值
    fMd5 := util.Md5File(file)
    b, info := db.ContainsFile(fMd5)
    if b { //如果服务器中已经有这张图片，直接返回该图片地址
        data := Response{Code: 1, Msg: info.Url}
        result, _ := json.Marshal(data)
        w.Write(result)
        return
    }
    //构造目录
    fileName := handler.Filename
    fileType := path.Ext(fileName)
    now := time.Now()
    year := now.Format("2006")
    month := now.Format("01")
    day := now.Format("02")
    root := "images/" + year + "/" + month + "/" + day + "/"

    if !util.ExistsFile(RootPath + root) {
        os.Mkdir(RootPath + root, 0777)
    }
    filePath := root + strconv.Itoa(now.Nanosecond()) + fileType //将文件名转化成时间戳
    f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()

    //从头开始读文件，保存文件到服务器
    file.Seek(0,0) //offset偏移位置，whence为0时表示从文件开始偏移，为1时表示从当前位置偏移，为2时表示从文件结尾偏移
    io.Copy(f, file)

    //save to DB...
    fi, _ := f.Stat()
    fileInfo := db.FileInfo{Name:fileName, Path: RootPath + filePath, Ip:r.Host, Url: Domain + filePath, FileType:fileType,
        Md5: fMd5, Size: fi.Size(), CreateTime:time.Now()}
    url, _ := db.SaveFile(fileInfo) //将文件信息保存如数据库
    data := Response{Code: 1, Msg: url}
    result, _ := json.Marshal(data)
    w.Write(result)
}






