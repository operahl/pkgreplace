package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ReplaceHelper struct {
	Root    string //根目录
	OldText string //需要替换的文本
	NewText string //新的文本
}

func (h *ReplaceHelper) DoWrok() error {

	return filepath.Walk(h.Root, h.walkCallback)

}
func CreateDateDir(Path string) string {

	if _, err := os.Stat(Path); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(Path, 0777) //0777也可以os.ModePerm
		os.Chmod(Path, 0777)
	}
	return Path
}
func (h ReplaceHelper) walkCallback(path string, f os.FileInfo, err error) error {

	if err != nil {
		return err
	}
	if f == nil {
		return nil
	}
	 if strings.Index(path, "/.git")>0{
	 	return nil
	 }
	if f.IsDir() {
		//fmt.Println("DIR:",path)
		newDirPath := strings.Replace(path, h.OldText, h.NewText, -1)
		//fmt.Println("NEWDIR:",newDirPath)
		CreateDateDir(newDirPath)
		return nil
	}
	//fmt.Println("FILE:",path)
	newFilePath := strings.Replace(path, h.OldText, h.NewText, -1)
	//fmt.Println("NEWFILE:",newFilePath)

	//文件类型需要进行过滤

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		//err
		return err
	}
	content := string(buf)

	//替换
	var newContent string
	if path == "goframe/go.mod"{
		newContent = strings.Replace(content, "module "+h.OldText, "module "+h.NewText, -1)
	}else if strings.Index(path, "goserver.sh")>0 {
		newContent = strings.Replace(content, "PRG=\""+h.OldText, "PRG=\""+h.NewText, -1)
	}else{
		newContent = strings.Replace(content, "\""+h.OldText+"/", "\""+h.NewText+"/", -1)
	}
	//重新写入
	ioutil.WriteFile(newFilePath, []byte(newContent), 0777)

	return err
}

func DoAction(ctx *cli.Context) error{

	if len(ctx.Args()) != 2 {
		return errors.New("useage:pkgreplace src obj")
	}
	srcname := ctx.Args()[0]
	dstname := ctx.Args()[1]

	if srcname == dstname{
		return errors.New("src != obj")
	}
	helper := ReplaceHelper{
		Root:    srcname,
		OldText: srcname,
		NewText: dstname,
	}
	err := helper.DoWrok()

	if err == nil {
		fmt.Println("done!")
	} else {
		fmt.Println("error:", err.Error())
	}
	return nil
}