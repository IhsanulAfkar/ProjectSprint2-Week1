# Project Sprint Batch 2 Week 1

## Installation

ofc install dulu go-nya

```
go mod init FOLDER_PROJECT # klo dah ada file go.mod gausah
go run .\tes.go
```

### Setup migrations

run di powershell
```
mkdir db/migrations
# pake scoop
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
scoop install migrate
```
