# 测试API接口
$baseUrl = "http://localhost:8888"

Write-Host "=== 测试用户注册 ===" -ForegroundColor Green
$registerData = @{
    username = "testuser"
    password = "123456"
    email = "test@example.com"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/v1/user/register" -Method POST -Body $registerData -ContentType "application/json"
    Write-Host "注册响应:" -ForegroundColor Cyan
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "注册失败: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

Write-Host "`n=== 测试用户登录 ===" -ForegroundColor Green
$loginData = @{
    username = "testuser"
    password = "123456"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/v1/user/login" -Method POST -Body $loginData -ContentType "application/json"
    Write-Host "登录响应:" -ForegroundColor Cyan
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "登录失败: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

Write-Host "`n=== 测试mDNS扫描（小范围） ===" -ForegroundColor Green
$scanData = @{
    cidr = "127.0.0.1/32"
    portRanges = @("5353")
    timeout = 2
    concurrency = 10
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/v1/scan" -Method POST -Body $scanData -ContentType "application/json"
    Write-Host "扫描响应:" -ForegroundColor Cyan
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "扫描失败: $_" -ForegroundColor Red
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Green
