package main

import (
	"fmt"
	"github.com/gookit/color"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	Git = "/usr/bin/git"
)

func main() {

	color.Blueln(fmt.Sprintf("[%s] ✨✨✨ Git仓库同步计划开始[start-sync]", time.Now().Format("2006-01-02 15:04:05")))

	if len(os.Args) < 2 {
		color.Redln("请输入有效的配置文件.")
		return
	} else {
		initializeViper(os.Args[1])
	}

	home := Config.Storage

	for _, current := range Config.Repo {
		var backup = current
		for _, d := range Config.Maps {
			if d.Key == current {
				backup = d.Value
				break
			}
		}

		dir := fmt.Sprintf("%s/%s", home, current[0:len(current)-4])

		color.Blueln(fmt.Sprintf("[%s] 🚀 开始同步 %s", time.Now().Format("2006-01-02 15:04:05"), current))

		fi, err := os.Stat(dir)
		if err != nil || !fi.IsDir() { // 判断是否有项目目录
			color.Greenln(fmt.Sprintf("\t不存在目录%s", dir))
			_ = command("/bin/rm", "-rf", dir).Run()           // 进行尝试性清空
			err = command(home, "/bin/mkdir", "-p", dir).Run() // 重新创建对应目录
			if err != nil {
				color.Redln(fmt.Sprintf("\t创建%s失败", dir))
				color.Grayln(err.Error())
			}

			repo := fmt.Sprintf(Config.Remote.Origin, current)
			c := command(dir, Git, "clone", repo, dir) // 重新clone代码
			err = c.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("\tgit clone %s 到 %s失败", repo, dir))
				color.Grayln("\t" + err.Error())
			}

			newRemote := fmt.Sprintf(Config.Remote.New, backup)
			err = command(dir, Git, "remote", "add", "new", newRemote).Run() // 添加新的远程地址
			if err != nil {
				color.Redln(fmt.Sprintf("\tgit添加新 %s 到 %s失败", newRemote, dir))
				color.Grayln("\t" + err.Error())
			}
		}

		_ = command(dir, Git, "pull", "origin").Run()              // 刷新原始分支
		err = command(dir, Git, "remote", "prune", "origin").Run() // 删除git上不存在的分支
		if err != nil {
			color.Redln(fmt.Sprintf("\tgit清除不存在的分支失败%s", dir))
			color.Grayln("\t" + err.Error())
		}

		remoteBranches := remoteBranch(dir)
		localBranches := localBranch(dir)

		// 对比分支，下拉不存在的分支
		for _, o := range remoteBranches {
			if !contain(localBranches, o) {
				c := command(dir, Git, "checkout", "-b", o, "origin/"+o)
				err = c.Run()
				if err != nil {
					color.Redln(fmt.Sprintf("\t检出远程分支失败%s", dir))
					color.Grayln(c.String())
					color.Grayln("\t" + err.Error())
				}
			}
		}
		// 推送到新的远程地址
		localBranches = localBranch(dir)
		for _, branch := range localBranches {
			a := command(dir, Git, "reset", "--hard")
			err = a.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("\t重置代码失败%s:%s", current, branch))
				color.Grayln(a.String())
				color.Grayln("\t" + err.Error())
			}

			b := command(dir, Git, "checkout", branch)
			err = b.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("\t切换分支失败%s %s", current, branch))
				color.Grayln(b.String())
				color.Grayln("\t" + err.Error())
			}

			c := command(dir, Git, "pull", "origin", branch)
			err = c.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("\t分支下拉失败%s", current))
				color.Grayln("\t" + c.String())
				color.Grayln("\t" + err.Error())

				_ = command(dir, Git, "branch", "-D", branch).Run()
			}

			d := command(dir, Git, "push", "new", branch)
			err = d.Run()
			if err != nil {
				color.Redln(fmt.Sprintf("\t推送到新的远程分支失败%s", current))
				color.Grayln("\t" + d.String())
				color.Grayln("\t" + err.Error())
			} else {
				color.Greenln(fmt.Sprintf("\t推到新仓库成功%s:%s", current, branch))
			}
			time.Sleep(time.Second)
		}
	}

	color.Blueln(fmt.Sprintf("[%s] ✨✨✨ Git仓库同步计划结束[end-sync]", time.Now().Format("2006-01-02 15:04:05")))
}

// 获取远程分支
func remoteBranch(dir string) []string {
	c := command(dir, Git, "branch", "-r")
	output, err := c.Output()
	if err != nil {
		color.Redln(fmt.Sprintf("\tgit获取远程分支失败%s", dir))
		color.Grayln("\t" + err.Error())
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
		color.Redln(fmt.Sprintf("\tgit获取远程分支失败%s", dir))
		color.Grayln("\t" + err.Error())
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
