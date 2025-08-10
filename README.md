# kageyamountain.net-backend

## 開発環境セットアップ（Windows）

### 前提

- Git
- WSL や Git Bash などのUnix系シェル実行環境
- Docker Desktop
- Pathは適宜通すこと

### Scoopのインストール
パッケージ管理ツールの `Scoop` をインストールする。  
`Windows PowerShell` で以下のコマンドを実行する。
```shell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
```
これ以降のコマンドはUnix系シェルでOK

### miseのインストール
環境管理ツールの `mise` をインストールする。
```shell
scoop install mise
mise install
```
