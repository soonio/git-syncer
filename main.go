package main

import (
	"fmt"
	"github.com/gookit/color"
	"os"
	"os/exec"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		color.Redln("请输入有效的配置文件.")
		return
	} else {
		initializeViper(os.Args[1])
	}

	home := Config.Storage

	for _, name := range Config.Repo {
		dir := fmt.Sprintf("%s/%s", home, name[0:len(name)-4])
		fi, _ := os.Stat(dir)
		if fi == nil { // 判断是否有项目目录
			_ = command("/bin/rm", "-rf", dir).Run()         // 进行尝试性清空
			_ = command(home, "/bin/mkdir", "-p", dir).Run() // 重新创建对应目录

			c := command(dir, "/usr/bin/git", "clone", fmt.Sprintf(Config.Remote.Origin, name), dir) // 重新clone代码
			err := c.Run()
			if err != nil {
				panic(err)
			}

			_ = command(dir, "/usr/bin/git", "remote", "add", "new", fmt.Sprintf(Config.Remote.New, name)).Run() // 添加新的远程地址
		}

		_ = command(dir, "/usr/bin/git", "remote", "prune", "origin") // 删除git上不存在的分支

		remoteBranches := remoteBranch(dir)
		localBranches := localBranch(dir)

		// 对比分支，下拉不存在的分支
		for _, o := range remoteBranches {
			if !contain(localBranches, o) {
				c := command(dir, "/usr/bin/git", "checkout", "-b", o, "origin/"+o)
				_ = c.Run()
			}
		}
		// 推送到新的远程地址
		localBranches = localBranch(dir)
		for _, branch := range localBranches {
			c := command(dir, "/usr/bin/git", "push", "new", branch)
			_ = c.Run()
		}
	}
}

// 获取远程分支
func remoteBranch(dir string) []string {
	c := command(dir, "/usr/bin/git", "branch", "-r")
	output, err := c.Output()
	if err != nil {
		panic(err)
	}
	var remoteBranches []string
	branches := strings.Split(string(output), "\n")
	for _, branch := range branches {
		if strings.Contains(branch, "->") {
			continue
		}
		branch = strings.Replace(branch, " ", "", -1)
		if branch == "" {
			continue
		}
		branch = branch[7:]
		remoteBranches = append(remoteBranches, branch)
	}
	return remoteBranches
}

// 获取本地分支
func localBranch(dir string) []string {
	c := command(dir, "/usr/bin/git", "branch")
	output, err := c.Output()
	if err != nil {
		panic(err)
	}
	var localBranches []string
	branches := strings.Split(string(output), "\n")
	for _, branch := range branches {
		branch = strings.Replace(branch, " ", "", -1)
		if branch == "" {
			continue
		}
		localBranches = append(localBranches, strings.Replace(branch, "*", "", -1))
	}
	return localBranches
}

func contain(s []string, d string) bool {
	for _, i := range s {
		if i == d {
			return true
		}
	}
	return false
}

func command(dir, name string, arg ...string) *exec.Cmd {
	c := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, arg...),
		Dir:  dir,
	}
	return c
}
