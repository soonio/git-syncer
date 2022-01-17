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
	}

	initializeViper(os.Args[1])

	//判断本地仓库是否存在，不存在则进行初始化
	//判断本地仓库是否已经添加了新的仓库地址，没有的话，则添加新的远程仓库地址
	//拉取最新的git远程代码，推送到coding远程仓库中
	//拉取最新的coding远程代码，推送到git远程仓库中
	//	(\/)(?!.*\1).*\.git
	for _, name := range Config.Repo {
		dir := fmt.Sprintf("repo/%s", name[0:len(name)-4])
		fi, _ := os.Stat(dir)
		if fi == nil { // 判断是否有项目目录
			_ = command("/bin/rm", "-rf", dir).Run()   // 进行尝试性清空
			_ = command(dir, "/bin/mkdir", "-p").Run() // 重新创建对应目录

			cmd := command(dir, "/usr/bin/git", "clone", fmt.Sprintf(Config.Remote.Origin, name), dir) // 重新clone代码
			err := cmd.Run()
			if err != nil {
				panic(err)
			}

			_ = command(dir, "/usr/bin/git", "remote", "add", "new", fmt.Sprintf(Config.Remote.New, name)) // 添加新的远程地址
		}

		_ = command(dir, "/usr/bin/git", "remote", "prune", "origin") // 删除git上不存在的分支

		// 获取远程分支
		cmd := command(dir, "/usr/bin/git", "branch", "-r")
		output, err := cmd.Output()
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

		// 获取本地分支
		cmd = command(dir, "/usr/bin/git", "branch")
		output, err = cmd.Output()
		if err != nil {
			panic(err)
		}
		var localBranches []string
		branches = strings.Split(string(output), "\n")
		for _, branch := range branches {
			branch = strings.Replace(branch, " ", "", -1)
			if branch == "" {
				continue
			}
			localBranches = append(localBranches, strings.Replace(branch, "*", "", 1))
		}
		// 对比分支，下拉不存在的分支
		for _, o := range remoteBranches {
			if !contain(localBranches, o) {
				cmd = command(dir, "/usr/bin/git", "checkout", "-b", o, "origin/"+o)
				err = cmd.Run()
			}
		}

	}
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
	cmd := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, arg...),
		Dir:  dir,
	}
	return cmd
}
