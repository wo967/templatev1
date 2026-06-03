# 清理并重新构建项目
Write-Host "=== 清理旧的构建文件 ===" -ForegroundColor Green
Remove-Item -Path "E:\code\template\templatev1\api\api.exe" -ErrorAction SilentlyContinue
Remove-Item -Path "E:\code\template\templatev1\rpc\rpc.exe" -ErrorAction SilentlyContinue

Write-Host "`n=== 同步Go模块依赖 ===" -ForegroundColor Green
Set-Location "E:\code\template\templatev1"
go work sync

Write-Host "`n=== 构建RPC服务 ===" -ForegroundColor Green
Set-Location "E:\code\template\templatev1\rpc"
go mod tidy
if ($LASTEXITCODE -eq 0) {
    go build -o rpc.exe .
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ RPC服务构建成功" -ForegroundColor Green
    } else {
        Write-Host "✗ RPC服务构建失败" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "✗ RPC模块依赖同步失败" -ForegroundColor Red
    exit 1
}

Write-Host "`n=== 构建API服务 ===" -ForegroundColor Green
Set-Location "E:\code\template\templatev1\api"
go mod tidy
if ($LASTEXITCODE -eq 0) {
    go build -o api.exe .
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ API服务构建成功" -ForegroundColor Green
    } else {
        Write-Host "✗ API服务构建失败" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "✗ API模块依赖同步失败" -ForegroundColor Red
    exit 1
}

Write-Host "`n=== 所有服务构建成功！ ===" -ForegroundColor Green
Write-Host "`n启动说明：" -ForegroundColor Cyan
Write-Host "1. 先启动RPC服务: .\rpc\rpc.exe" -ForegroundColor White
Write-Host "2. 再启动API服务: .\api\api.exe" -ForegroundColor White
Write-Host "3. 测试注册接口: curl -X POST http://localhost:8888/v1/user/register -H 'Content-Type: application/json' -d '{`"username`":`"test`",`"password`":`"123456`",`"email`":`"test@test.com`"}'" -ForegroundColor White
Write-Host "4. 测试登录接口: curl -X POST http://localhost:8888/v1/user/login -H 'Content-Type: application/json' -d '{`"username`":`"test`",`"password`":`"123456`"}'" -ForegroundColor White
