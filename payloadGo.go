package  payloadGo

import (
	"fmt"
	"os"
  "strconv"
)
func getSelfPath() (path string) {
  path, err := os.Executable()
  if err != nil { panic (err) }
  return path
}

func getFile(path string) (file *os.File){
	file, err := os.Open(path)
	if err != nil { panic(err) }
  return file
}

func getFileInfo(path string) (fileInfo os.FileInfo){
  fileInfo, err := os.Stat(path)
	if err != nil { panic(err) }
  return fileInfo
}

func getLengthPayload(file *os.File, fileInfo os.FileInfo ) (lengthPayload int)  {
//get length payload and bytes->string->int from file
  bufLenPayload := make([]byte, 8) //creation d'un array de 8 octets pour stocker nos 8 caracteres 00000000 qui est la taille de  notre payload
  file.ReadAt(bufLenPayload, fileInfo.Size() - int64(9))
  lengthPayload, err := strconv.Atoi(string(bufLenPayload)) // convertit nos bytes en string puis en int int(string(bytes)))
  if err != nil { return 0 }
  return   lengthPayload
}

func deletePayload(file *os.File, fileInfo os.FileInfo, path string  ) (int,string){
  lengthPayload := getLengthPayload(file, fileInfo)
  if lengthPayload == 0 { return 201,"yet_delete" }
    newSize:= fileInfo.Size() - int64(9) - int64(lengthPayload)
    os.Truncate(path, newSize)
    return 200,"delete"
}

func putPayload(path string, payload string, file *os.File, fileInfo os.FileInfo) (int,string){
  lengthPayload  := getLengthPayload(file, fileInfo)
  if lengthPayload != 0 { deletePayload(file, fileInfo,path)  }

  f, err := os.OpenFile(path,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil { fmt.Println(err) }

  defer f.Close()
  if _, err := f.WriteString(payload+fmt.Sprintf("%08d", len(payload))+"\n"); err != nil {
    fmt.Println(err)
  }
  return 200,"okey"
}

func getPayload(path string,  file *os.File, fileInfo os.FileInfo, lengthPayload int ) (int,string,string){
  if lengthPayload == 0 { 
    return 404,"empty or corrupt",""
    } else {
    bufPayload := make([]byte, lengthPayload) //on met la taille du payload en int
    file.ReadAt(bufPayload, fileInfo.Size() - int64(9) - int64(lengthPayload) )
    return 200, "okey", string(bufPayload[:len(bufPayload)-1])
   }

}

type response struct {
    REQUEST_METHOD string
    PARAMETERS string
    STATUS_CODE int
    STATUS_MESSAGE string
    DATA string
}

func DELETE(path string) response{
  if path == "" {
    path = getSelfPath()
  }
  status_code, message := deletePayload(getFile(path),getFileInfo(path), path)
  return response{"DELETE","path: "+path,status_code,message,""}
}

func PUT(path string, payload string) response{
  if path == "" {
    path = getSelfPath()
  }
  if payload == "" { payload = "\ndefaultPayload:\n\t- key1: false\n"}

  status_code, message :=  putPayload(path,payload,getFile(path),getFileInfo(path))
  return response{"PUT","path: "+path+" Payload: "+payload,status_code,message,""}
}

func GET(path string) response{
  if path == "" {
    path = getSelfPath()
  }

  status_code, message, data := getPayload(path,getFile(path),getFileInfo(path),getLengthPayload(getFile(path),getFileInfo(path)))
  return response{"GET","path: "+path,status_code,message,data}
}



//func main(){
//  log.Printf("%+v\n", GET(""))
//  log.Printf("%+v\n", PUT("","ddddddd"))
//  log.Printf("%+v\n", GET(""))
//  log.Printf("%+v\n", DELETE(""))
//  log.Printf("%+v\n", GET(""))
//}
