package util

import (
    "strconv"
    "fmt"
    "os"
    "crypto/md5"
    "io"
    "encoding/hex"
    "mime/multipart"
)
/**
 * 将文件大小转换成人类认识的展示方式
 */
func FormatSize(size int) string{
    if size < 2048 {
        return strconv.Itoa(size) + "kb"
    }
    //to Mb
    fSize, _ := strconv.ParseFloat(strconv.Itoa(size), 32/64)
    mSize := fSize/1024
    sSize := fmt.Sprintf("%.2f", mSize) + "Mb"
    return sSize
}

/**
 * 查询文件（文件件）是否存在
 */
func ExistsFile(path string) bool {
    _, err := os.Stat(path) //os.Stat获取文件信息
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}

/**
 * 返回上传文件的MD5校验码
 */
func Md5File(file multipart.File) string {
    m := md5.New()
    _, err := io.Copy(m, file)
    if err != nil {
        fmt.Println(err)
        return ""
    }
    md5s := hex.EncodeToString(m.Sum(nil))
    return md5s
}

/**
 * 异常处理
 */
func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}