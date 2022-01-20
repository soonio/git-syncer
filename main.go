package main

import (
	"fmt"
	"github.com/gookit/color"
	"os"
	"os/exec"
	"strings"
)

const (
	Git = "/usr/bin/git"
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
			color.Greenln(fmt.Sprintf("不存在目录%s", dir))
			err := command("/bin/rm", "-rf", dir).Run() // 进行尝试性清空
			if err != nil {
				color.Redln(fmt.Sprintf("尝试删除%s失败", dir))
				color.Grayln(err.Error())
			}

			err = command(home, "/bin/mkdir", "-p", dir).Run() // 重新创建对应目录
			if err != nil {
				color.Redln(fmt.Sprintf("创建%s失败", dir))
				color.Grayln(err.Error())
			}

			repo := fmt.Sprintf(Config.Remote.Origin, name)
			c := command(dir, Git, "clone", repo, dir) // 重新clone代码
			err = c.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("git clone %s 到 %s失败", repo, dir))
				color.Grayln(err.Error())
			}

			newRemote := fmt.Sprintf(Config.Remote.New, name)
			err = command(dir, Git, "remote", "add", "new", newRemote).Run() // 添加新的远程地址
			if err != nil {
				color.Redln(fmt.Sprintf("git添加新 %s 到 %s失败", newRemote, dir))
				color.Grayln(err.Error())
			}
		}

		err := command(dir, Git, "remote", "prune", "origin").Run() // 删除git上不存在的分支
		if err != nil {
			color.Redln(fmt.Sprintf("git清除不存在的分支失败%s", dir))
			color.Grayln(err.Error())
		}

		remoteBranches := remoteBranch(dir)
		localBranches := localBranch(dir)

		// 对比分支，下拉不存在的分支
		for _, o := range remoteBranches {
			if !contain(localBranches, o) {
				c := command(dir, Git, "checkout", "-b", o, "origin/"+o)
				err = c.Run()
				if err != nil {
					color.Redln(fmt.Sprintf("检出远程分支失败%s", dir))
					color.Grayln(c.String())
					color.Grayln(err.Error())
				}
			}
		}
		// 推送到新的远程地址
		localBranches = localBranch(dir)
		for _, branch := range localBranches {
			a := command(dir, Git, "reset", "--hard")
			err = a.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("重置代码失败%s:%s", name, branch))
				color.Grayln(a.String())
				color.Grayln(err.Error())
			}

			b := command(dir, Git, "pull", "origin", branch)
			err = b.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("下拉代码失败%s", name))
				color.Grayln(b.String())
				color.Grayln(err.Error())
			}

			c := command(dir, Git, "push", "new", branch)
			err = c.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("推送到新的远程分支失败%s", name))
				color.Grayln(c.String())
				color.Grayln(err.Error())
			} else {
				color.Greenln(fmt.Sprintf("推到新仓库成功%s:%s", name, branch))
			}
		}
	}
}

// 获取远程分支
func remoteBranch(dir string) []string {
	c := command(dir, Git, "branch", "-r")
	output, err := c.Output()
	if err != nil {
		color.Redln(fmt.Sprintf("git获取远程分支失败%s", dir))
		color.Grayln(err.Error())
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
		if branch[0:4] == "new/" {
			continue
		}
		branch = branch[7:]
		remoteBranches = append(remoteBranches, branch)
	}
	return remoteBranches
}

// 获取本地分支
func localBranch(dir string) []string {
	c := command(dir, Git, "branch")
	output, err := c.Output()
	if err != nil {
		color.Redln(fmt.Sprintf("git获取远程分支失败%s", dir))
		color.Grayln(err.Error())
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
