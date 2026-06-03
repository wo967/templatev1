# 启动所有服务
Write-Host "=== TemplateV1 服务启动器 ===" -ForegroundColor Cyan

# 检查是否已构建
if (-not (Test-Path "E:\code\template\templatev1\rpc\rpc.exe")) {
    Write-Host "RPC服务未构建，正在构建..." -ForegroundColor Yellow
    Set-Location "E:\code\template\templatev1\rpc"
    go build -o rpc.exe .
}

if (-not (Test-Path "E:\code\template\templatev1\api\api.exe")) {
    Write-Host "API服务未构建，正在构建..." -ForegroundColor Yellow
    Set-Location "E:\code\template\templatev1\api"
    go build -o api.exe .
}

Write-Host "`n请选择启动模式：" -ForegroundColor Cyan
Write-Host "1. 启动所有服务（RPC + API）" -ForegroundColor White
Write-Host "2. 仅启动RPC服务" -ForegroundColor White
Write-Host "3. 仅启动API服务" -ForegroundColor White
Write-Host "4. 重新构建并启动" -ForegroundColor White
$choice = Read-Host "请输入选项 (1-4)"

switch ($choice) {
    "1" {
        Write-Host "`n=== 启动RPC服务 ===" -ForegroundColor Green
        Start-Process -FilePath "E:\code\template\templatev1\rpc\rpc.exe" -WorkingDirectory "E:\code\template\templatev1\rpc"
        Write-Host "RPC服务正在启动..." -ForegroundColor Yellow
        Start-Sleep -Seconds 2
        
        Write-Host "`n=== 启动API服务 ===" -ForegroundColor Green
        Start-Process -FilePath "E:\code\template\templatev1\api\api.exe" -WorkingDirectory "E:\code\template\templatev1\api"
        Write-Host "API服务正在启动..." -ForegroundColor Yellow
        
        Write-Host "`n✓ 所有服务已启动！" -ForegroundColor Green
        Write-Host "API服务地址: http://localhost:8888" -ForegroundColor White
        Write-Host "RPC服务地址: 127.0.0.1:8083" -ForegroundColor White
        Write-Host "`n按任意键打开测试脚本..." -ForegroundColor Cyan
        $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
        .\test.ps1
    }
    "2" {
        Write-Host "`n=== 启动RPC服务 ===" -ForegroundColor Green
        Set-Location "E:\code\template\templatev1\rpc"
        .\rpc.exe
    }
    "3" {
        Write-Host "`n=== 启动API服务 ===" -ForegroundColor Green
        Set-Location "E:\code\template\templatev1\api"
        .\api.exe
    }
    "4" {
        Write-Host "`n=== 重新构建项目 ===" -ForegroundColor Green
        Set-Location "E:\code\template\templatev1"
        .\build.ps1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "`n构建成功！是否启动服务？(Y/N)" -ForegroundColor Cyan
            $startChoice = Read-Host
            if ($startChoice -eq "Y" -or $startChoice -eq "y") {
                & $MyInvocation.MyCommand.Path
            }
        }
    }
    default {
        Write-Host "无效的选项" -ForegroundColor Red
    }
}
