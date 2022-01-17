#!/usr/bin/env bash
#cd $(dirname "$0") || exit
#
#echo "🚀 $(date "+%Y-%m-%d %H:%M:%S") 开始同步github和coding"
#
#git pull --all
#git branch -a
#git tag --points-at
#
#git push coding
#git push coding --tag

Origin=git@gitlab.10in.com
New=git@e.coding.net

# 同步仓库
Sync(){
  echo "$Origin:$1.git $New:sh-10in/$1.git"

  dir="repo/$1"

  #  判断本地仓库是否存在，不存在则进行初始化
  if [ ! -d "$dir" ]; then
#    mkdir -p "$dir"
    echo ${$dir/\//\//}
    echo "文件夹不存在"
  fi

  # 判断本地仓库是否已经添加了新的仓库地址，没有的话，则添加新的远程仓库地址

  # 拉取最新的git远程代码，推送到coding远程仓库中
  # 拉取最新的coding远程代码，推送到git远程仓库中

}
Sync pp,api
#Sync pp/api
#Sync pp/php-middleware
#Sync pp/front
