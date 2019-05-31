package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
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

func (h ReplaceHelper) walkCallback(path string, f os.FileInfo, err error) error {

	if err != nil {
		return err
	}
	if f == nil {
		return nil
	}
	if f.IsDir() {
		//fmt.Pringln("DIR:",path)
		return nil
	}

	//文件类型需要进行过滤

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		//err
		return err
	}
	content := string(buf)

	//替换
	newContent := strings.Replace(content, h.OldText, h.NewText, -1)

	//重新写入
	ioutil.WriteFile(path, []byte(newContent), 0)

	return err
}

func DoAction(ctx *cli.Context) error{

	if len(ctx.Args()) != 2 {
		return errors.New("useage:pkgreplace src obj")
	}
	srcname := ctx.Args()[0]
	dstname := ctx.Args()[1]

	cpstr:="cp -r "+srcname+" "+dstname

	exec_shell(cpstr)

	oldText:="\""+srcname+"/"
	newText:="\""+dstname+"/"

	helper := ReplaceHelper{
		Root:    dstname,
		OldText: oldText,
		NewText: newText,
	}
	err := helper.DoWrok()
	if err == nil {
		fmt.Println("done!")
	} else {
		fmt.Println("error:", err.Error())
	}
	return nil
}
func exec_shell(s string) (string, error){
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	return out.String(), err
}