# software_engineering_project


## 分支介紹 RedirectByPermission
這是專門處理不同權限渲染的分支，為避免master的節點過多，因此創建此分支

關於遠端分支操作詳請參見，裡面講得很清楚
https://blog.longwin.com.tw/2013/11/git-create-remote-branch-2013/

另外，你們可能注意到github上多了「Compare & pull request」，這是是在說將分支合併到主幹上的請求，當功能都做完時再按就好了


## github 協作教學
這邊不會講解原理，網路上的講解很多。只告訴大家如何快速地操作
### 從遠端版本庫首次下載至指定目錄

    git clone https://github.com/hsinnnlu/software_engineering_project.git
記得下載完後要更新專案的套件與大家的同步哦！
### 將自己的更新推送至遠端版本庫
三個步驟上傳更新
```
git add (欲上傳的檔案)
```
```
git commit -m "(更新標題)"
```
```
git push
```
### 更新本地端至遠端版本庫
如果自己的版本落後於遠端版本庫則需要做這件事。需要注意「conflict 衝突」  
這邊先介紹簡單的更新（拉取）方法：
```
git pull
```
如遇衝突先緩一下，記得在群組（開會）提出，不然怕合併掉重要的程式

### git本地端的操作
git有很多版本管理的指令可以在網路上學，以下指令熟記的話可以大幅幫助開發管理
1. git checkout 還原  
2. git branch 分支
3. git status 追蹤狀態
4. git log 版本紀錄

另外其他好用的指令就給大家去探索了