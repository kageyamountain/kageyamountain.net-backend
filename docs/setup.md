## 開発環境セットアップ
Windows環境を例に書いているが、Macもhomebrewで `mise` をインストールすれば同様の手順でセットアップ可能なはず。

### 前提

- Git
- WSL や Git Bash などのUnix系シェル実行環境
- Docker Desktop
- Pathは適宜通すこと

### インストール手順
パッケージ管理ツールの `Scoop` をインストールする。  
`Windows PowerShell` で以下のコマンドを実行する。
```shell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
```
これ以降のコマンドはUnix系シェルでOK

環境管理ツールの `mise` をインストールする。
```shell
scoop install mise
mise install
```

go mod tidyで依存関係を整理する。
```shell
go mod tidy
```