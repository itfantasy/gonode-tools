package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
)

func main() {

	conf, err := ini.Load(io.CurrentDir() + "conf.ini")
	if err != nil {
		fmt.Println(err)
	}

	srcPath := conf.Get("proto", "src")

	csPath := conf.Get("cs", "dst")
	goPath := conf.Get("go", "dst")
	jsPath := conf.Get("js", "dst")
	tsPath := conf.Get("js", "dst4ts")
	javaPath := conf.Get("java", "dst")

	files, err := io.ListDir(srcPath, ".proto")
	if err != nil {
		fmt.Println(err)
		return
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	for _, file := range files {
		if isOn(conf, "cs") {
			dstFile := strings.Replace(file, "proto", "pb", -1)
			dstFile += ".cs"
			cmd := exec.Command("bin/protobuf-net/ProtoGen/protogen", "-i:"+srcPath+file, "-o:"+csPath+dstFile, "-ns:proto")
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			} else {
				fmt.Println(cmd.Args)
			}
		}
		if isOn(conf, "go") {
			cmd := exec.Command("bin/protoc-gen-go/protoc", "--proto_path="+srcPath, "--go_out="+goPath, file)
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			} else {
				fmt.Println(cmd.Args)
			}
		}
		if isOn(conf, "java") {
			cmd := exec.Command("bin/protoc-gen-go/protoc", "--proto_path="+srcPath, "--java_out="+javaPath, file)
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			} else {
				fmt.Println(cmd.Args)
			}
		}
	}
	if isOn(conf, "js") {
		txt := "call pbjs -t static-module -w commonjs -o " + jsPath + "proto.js " + srcPath + "*.proto"
		if tsPath != "" {
			txt += "\r\n"
			txt += "call pbts -o " + tsPath + "proto.d.ts " + jsPath + "proto.js"
		}
		txt = strings.ReplaceAll(txt, "/", "\\")
		io.SaveFile(io.CurrentDir()+"do.cmd", txt)
		cmd := exec.Command(io.CurrentDir() + "do.cmd")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		} else {
			io.DeleteFile(io.CurrentDir() + "do.cmd")
			fmt.Println(cmd.Args)
		}

	}
}

func isOn(conf *ini.Conf, lang string) bool {
	return conf.GetInt(lang, "on", 0) > 0
}
