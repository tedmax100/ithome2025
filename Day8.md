---
title: "Day 8：Pod 日誌與 Debug"
tags: 2025鐵人賽
date: 2025-07-20
---

#### Day 8：Pod 日誌與 Debug
##### 重點
- kubectl logs、exec、describe 等指令使用
##### Lab
- 模擬應用錯誤，透過 logs/exec/describe 排查問題

在 Kubernetes 環境中，當應用程式出現問題時，能夠有效地診斷和排除故障是維護服務穩定性的關鍵。今天，我們將深入探討 Kubernetes 中的日誌收集和除錯技術，特別是 kubectl logs、kubectl exec 和 kubectl describe 等指令的使用方法。

![](https://static.learnkube.com/7bcf2c9e9dce01269c436a16b77b276f.png)
[k8s troubleshooting flowchart](https://static.learnkube.com/dac10c60ec5d2fe6bd3d3f8736cf0ce0.pdf)

# Kubernetes 除錯工具概述
Kubernetes 提供了多種工具來幫助開發者和運維人員診斷問題：

kubectl logs：查看容器的日誌輸出
kubectl exec：在容器中執行命令
kubectl describe：顯示資源的詳細信息
kubectl get：獲取資源的基本信息
kubectl port-forward：將本地端口轉發到 Pod
這些工具組合使用，可以幫助我們快速定位和解決問題。

## kubectl logs：容器日誌查看
kubectl logs 是查看容器日誌的主要工具，它可以顯示容器的標準輸出和標準錯誤。

基本用法
```bash
# 查看指定 Pod 的日誌
kubectl logs <pod-name>

# 查看特定命名空間中的 Pod 日誌
kubectl logs <pod-name> -n <namespace>

# 如果 Pod 有多個容器，指定容器名稱
kubectl logs <pod-name> -c <container-name>
```

進階用法
```bash

# 持續查看日誌（類似 tail -f）
kubectl logs <pod-name> -f

# 查看最近的 N 行日誌
kubectl logs <pod-name> --tail=100

# 查看特定時間段的日誌
kubectl logs <pod-name> --since=1h

# 查看日誌直到特定時間點
kubectl logs <pod-name> --since-time=2023-07-26T10:00:00Z

# 顯示時間戳
kubectl logs <pod-name> --timestamps
```

查看已終止容器的日誌
```bash
# 查看已終止的容器日誌
kubectl logs <pod-name> --previous
```

這對於分析容器崩潰原因非常有用。

查看由控制器管理的 Pod 日誌
對於由 Deployment、StatefulSet 等控制器管理的 Pod，可以使用標籤選擇器：
```bash
# 查看特定 Deployment 的所有 Pod 日誌
kubectl logs -l app=nginx

# 查看第一個匹配的 Pod 日誌
kubectl logs -l app=nginx --max-log-requests=1
```

kubectl exec：在容器中執行命令
kubectl exec 允許我們在容器中執行命令，就像 SSH 到遠程服務器一樣。這對於檢查容器內部狀態、檔案系統和網路連接非常有用。

基本用法
```bash
# 在 Pod 中執行單一命令
kubectl exec <pod-name> -- <command>

# 在特定命名空間的 Pod 中執行命令
kubectl exec <pod-name> -n <namespace> -- <command>

# 如果 Pod 有多個容器，指定容器
kubectl exec <pod-name> -c <container-name> -- <command>
```

常見用例
```bash
# 檢查容器環境變數
kubectl exec <pod-name> -- env

# 檢查容器檔案系統
kubectl exec <pod-name> -- ls -la /app

# 檢查網路連接
kubectl exec <pod-name> -- netstat -tuln

# 檢查進程
kubectl exec <pod-name> -- ps aux
```

交互式 Shell
```bash
# 啟動交互式 shell 會話
kubectl exec -it <pod-name> -- /bin/bash

# 如果容器沒有 bash，可以使用 sh
kubectl exec -it <pod-name> -- /bin/sh
```

交互式 shell 對於複雜的除錯非常有用，可以在容器內執行多個命令。

## kubectl describe：資源詳細信息
kubectl describe 提供了資源的詳細信息，包括事件、狀態和配置。這對於理解資源的當前狀態和診斷問題非常有用。

基本用法
```bash
# 描述特定 Pod
kubectl describe pod <pod-name>

# 描述特定命名空間中的 Pod
kubectl describe pod <pod-name> -n <namespace>

# 描述所有 Pod
kubectl describe pods

# 描述符合標籤選擇器的 Pod
kubectl describe pods -l app=nginx
```

其他資源類型
```bash
# 描述 Deployment
kubectl describe deployment <deployment-name>

# 描述 Service
kubectl describe service <service-name>

# 描述 Node
kubectl describe node <node-name>

# 描述 PersistentVolumeClaim
kubectl describe pvc <pvc-name>
```

重要信息
kubectl describe pod 輸出包含許多有用的信息：

基本信息：名稱、命名空間、標籤等
狀態：當前狀態、IP 地址、啟動時間
容器：容器狀態、映像、端口、環境變數
卷掛載：已掛載的卷和掛載點
條件：Pod 條件（如 PodScheduled、Initialized、Ready）
事件：與 Pod 相關的事件（最重要的除錯信息）
特別注意 Events 部分，它記錄了與 Pod 相關的重要事件，如調度、拉取映像、啟動容器等。當 Pod 出現問題時，這裡通常會有有用的錯誤信息。

# 常見問題診斷流程
讓我們看看如何結合使用這些工具來診斷常見問題：

1. Pod 無法啟動
當 Pod 無法啟動時，通常遵循以下步驟：

```bash
# 1. 檢查 Pod 狀態
kubectl get pod <pod-name>

# 2. 查看詳細信息和事件
kubectl describe pod <pod-name>

# 3. 如果容器曾經啟動過，查看之前的日誌
kubectl logs <pod-name> --previous
```

常見原因：

映像拉取失敗
資源不足
配置錯誤
健康檢查失敗


2. 應用程式崩潰
當應用程式啟動後崩潰：
```bash
# 1. 查看容器日誌
kubectl logs <pod-name>

# 2. 檢查重啟計數和最近事件
kubectl describe pod <pod-name>

# 3. 如果可能，進入容器檢查
kubectl exec -it <pod-name> -- /bin/sh
```

常見原因：

應用程式錯誤
配置錯誤
依賴服務不可用
資源限制（OOM）

3. 應用程式運行但行為異常
當應用程式運行但行為不符合預期：

```bash
# 1. 查看實時日誌
kubectl logs <pod-name> -f

# 2. 進入容器進行調查
kubectl exec -it <pod-name> -- /bin/sh

# 3. 檢查配置和環境變數
kubectl exec <pod-name> -- env
kubectl exec <pod-name> -- cat /path/to/config
```

常見原因：

配置錯誤
環境變數問題
網路連接問題
權限問題

4. 網路連接問題
當服務間無法通信：
```bash
# 1. 檢查服務定義
kubectl describe service <service-name>

# 2. 測試 DNS 解析
kubectl exec <pod-name> -- nslookup <service-name>

# 3. 測試網路連接
kubectl exec <pod-name> -- curl -v <service-name>:<port>

# 4. 檢查網路策略
kubectl get networkpolicy
```

常見原因：

服務選擇器配置錯誤
DNS 問題
網路策略阻止流量
端口配置錯誤

# 實作：模擬應用錯誤與排查
現在，讓我們通過實際操作來練習如何診斷和解決 Kubernetes 中的常見問題。

準備工作
首先，創建一個命名空間用於我們的實驗：
```bash
kubectl create namespace debug-demo
```
場景 1：映像拉取錯誤
讓我們部署一個使用不存在的映像的 Pod：
```bash
cat <<EOF | kubectl apply -f - -n debug-demo
apiVersion: v1
kind: Pod
metadata:
  name: image-error-pod
spec:
  containers:
  - name: nginx
    image: nginx:nonexistent-tag
EOF
```

現在，讓我們診斷問題：
```bash
# 檢查 Pod 狀態
kubectl get pod image-error-pod -n debug-demo

# 查看詳細信息和事件
kubectl describe pod image-error-pod -n debug-demo
```

在 Events 部分，你應該能看到類似這樣的錯誤：
```
Failed to pull image "nginx:nonexistent-tag": rpc error: code = NotFound desc = ...
```
解決方案：修改 Pod 定義，使用正確的映像標籤：
```bash
kubectl delete pod image-error-pod -n debug-demo

cat <<EOF | kubectl apply -f - -n debug-demo
apiVersion: v1
kind: Pod
metadata:
  name: image-error-pod
spec:
  containers:
  - name: nginx
    image: nginx:latest
EOF
```


# 總結
在本文中，我們深入探討了 Kubernetes 中的日誌收集和除錯技術，特別是 kubectl logs、kubectl exec 和 kubectl describe 等指令的使用方法。我們通過實際操作模擬了各種常見問題，並學習了如何診斷和解決這些問題。

有效的除錯需要系統性的方法：

觀察症狀：使用 kubectl get 確定資源狀態
收集信息：使用 kubectl describe 和 kubectl logs 獲取詳細信息
分析根本原因：根據日誌和事件確定問題根源
嘗試解決方案：修改配置、重新部署或其他修復措施
驗證修復：確認問題已解決
掌握這些除錯技術將幫助你更有效地管理 Kubernetes 應用程式，減少停機時間，提高服務可靠性。

