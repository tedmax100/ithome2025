---
title: "Day 4：Deployment 與 ReplicaSet"
tags: 2025鐵人賽
date: 2025-07-20
---

#### Day 4：Deployment 與 ReplicaSet
##### 重點
- Deployment、ReplicaSet 概念與用途
##### Lab
- 建立 Deployment，調整 replicas 數量，觀察自動擴縮容器

在昨天的課程中，我們學習了 Kubernetes 的基本資源類型和 YAML 配置，特別關注了 Pod 這個最小的部署單元。然而，在實際生產環境中，我們很少直接管理單個 Pod，而是使用更高級別的抽象，如 Deployment 和 ReplicaSet。今天，我們將深入了解這些關鍵資源，它們如何協同工作，以及如何利用它們實現應用程序的自動化管理。

# Deployment 概述
Deployment 是 Kubernetes 中最常用的資源類型之一，它提供了聲明式的應用程序更新能力，並管理 Pod 和 ReplicaSet 的創建和擴展。

**Deployment 的主要功能**
聲明式更新：只需定義期望狀態，Kubernetes 會自動調整實際狀態以匹配
滾動更新：在不停機的情況下更新應用程序
版本回滾：輕鬆回滾到先前的版本
擴展/縮減：調整應用程序的副本數量
暫停/恢復：暫停更新過程進行調試，然後恢復

**Deployment 如何工作**
當您創建 Deployment 時，Kubernetes 會：

創建一個 ReplicaSet
ReplicaSet 創建並管理指定數量的 Pod 副本
當您更新 Deployment 時，會創建一個新的 ReplicaSet，逐漸增加新 ReplicaSet 的副本數，同時減少舊 ReplicaSet 的副本數
最終，舊 ReplicaSet 的副本數變為 0，但它仍然保留在歷史記錄中，以便回滾

**ReplicaSet 概述**
ReplicaSet 的主要目的是維護一組 Pod 副本的穩定運行。它確保指定數量的 Pod 副本在任何時候都處於運行狀態。

ReplicaSet 的主要功能
維護 Pod 副本數量：確保指定數量的 Pod 副本運行
Pod 健康監控：自動替換失敗的 Pod
水平擴展：增加或減少 Pod 副本數量

**ReplicaSet vs. 舊版 ReplicationController**
ReplicaSet 是 ReplicationController 的下一代替代品，提供了更強大的選擇器支持：

ReplicationController 只支持等值選擇器（例如 app=nginx）
ReplicaSet 支持集合選擇器（例如 app in (nginx, apache)）
在實踐中，我們很少直接創建 ReplicaSet，而是通過 Deployment 間接創建和管理它們。

# Deployment YAML 結構
讓我們看看一個基本的 Deployment YAML 配置：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3                  # 副本數量
  selector:
    matchLabels:
      app: nginx               # 選擇器，用於識別 Pod
  template:                    # Pod 模板
    metadata:
      labels:
        app: nginx             # Pod 標籤
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
        ports:
        - containerPort: 80
        resources:
          limits:
            cpu: "0.5"
            memory: "512Mi"
          requests:
            cpu: "0.2"
            memory: "256Mi"
```

關鍵來欄位說明
replicas：指定需要運行的 Pod 副本數量
selector：定義 Deployment 如何找到要管理的 Pod
template：定義 Pod 的規格，包括容器、卷等
strategy（未顯示）：定義更新策略，如 RollingUpdate 或 Recreate

# 實作：創建和管理 Deployment
現在，讓我們通過實際操作來學習如何創建和管理 Deployment。

步驟 1：創建 Deployment YAML 文件
創建一個名為 nginx-deployment.yaml 的文件，內容如下：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
        ports:
        - containerPort: 80
```

步驟 2：使用 kubectl apply 部署
```bash
kubectl apply -f nginx-deployment.yaml
```

步驟 3：觀察 Deployment 狀態
```bash
# 查看 Deployment
kubectl get deployments

# 查看 ReplicaSet
kubectl get rs

# 查看 Pod
kubectl get pods

# 查看 Deployment 詳細信息
kubectl describe deployment nginx-deployment
```

步驟 4：擴展 Deployment
有兩種方式可以擴展 Deployment：

1. 使用 kubectl scale 命令：
```bash
kubectl scale deployment nginx-deployment --replicas=5
```

2. 修改 YAML 文件並重新應用：
```yaml
# 修改 nginx-deployment.yaml 中的 replicas 為 5
spec:
  replicas: 5
```

然後重新應用：
```bash
kubectl apply -f nginx-deployment.yaml
```

步驟 5：觀察擴展效果
```bash
kubectl get pods -w
```

這將實時顯示 Pod 的變化。您應該能看到新的 Pod 被創建，直到達到指定的副本數。

步驟 6：縮減 Deployment

```bash
kubectl scale deployment nginx-deployment --replicas=2
```

再次觀察 Pod 的變化：

```bash
kubectl get pods -w
```

您應該能看到多餘的 Pod 被終止，直到只剩下指定數量的副本。

# 更新 Deployment
Deployment 的一個強大功能是能夠無縫更新應用程序。

更新容器映像
```bash
kubectl set image deployment/nginx-deployment nginx=nginx:1.22
```

或者修改 YAML 文件並重新應用：
```yaml
# 修改 nginx-deployment.yaml 中的 image
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.22  # 更新版本
```

然後：
```bash
kubectl apply -f nginx-deployment.yaml
```

觀察更新過程
```kubectl rollout status deployment/nginx-deployment
```

您還可以查看 ReplicaSet 的變化：
```bash
kubectl get rs
```

您應該能看到一個新的 ReplicaSet 被創建，它的副本數逐漸增加，而舊的 ReplicaSet 的副本數逐漸減少。

回滾更新
如果新版本出現問題，您可以輕鬆回滾：
```bash
# 回滾到上一個版本
kubectl rollout undo deployment/nginx-deployment

# 回滾到特定版本
kubectl rollout undo deployment/nginx-deployment --to-revision=2
```

查看更新歷史
```bash
kubectl rollout history deployment/nginx-deployment
```

## Deployment 更新策略
Deployment 支持兩種更新策略：

RollingUpdate（默認）：逐步替換 Pod，確保服務不中斷
Recreate：先刪除所有舊 Pod，然後創建新 Pod（會導致短暫服務中斷）
您可以在 YAML 中指定更新策略：

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # 可以超出期望副本數的最大 Pod 數
      maxUnavailable: 1  # 更新過程中允許不可用的最大 Pod 數
```

## ReplicaSet 詳解
雖然我們通常通過 Deployment 間接使用 ReplicaSet，但了解 ReplicaSet 的工作原理對於理解 Kubernetes 的自動修復和擴展機制非常重要。

### ReplicaSet YAML 結構
```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: nginx-replicaset
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
```

### ReplicaSet 選擇器類型
ReplicaSet 支持多種選擇器類型：

1. 等值選擇器：
```yaml
selector:
  matchLabels:
    app: nginx
```

2. 集合選擇器：
```yaml
selector:
  matchExpressions:
    - key: app
      operator: In
      values:
        - nginx
        - web
    - key: tier
      operator: NotIn
      values:
        - dev
```

## Deployment 與 ReplicaSet 的關係
讓我們通過一個實例來理解 Deployment 和 ReplicaSet 之間的關係：

1. 創建 Deployment：
```bash
kubectl apply -f nginx-deployment.yaml
```

2. 檢查創建的資源：
```bash
kubectl get deploy,rs,pods
```

您會看到：

1 個 Deployment（nginx-deployment）
1 個 ReplicaSet（自動生成名稱，如 nginx-deployment-66b6c48dd5）
3 個 Pod（名稱以 ReplicaSet 名稱為前綴）

3. 更新 Deployment：

```bash
kubectl set image deployment/nginx-deployment nginx=nginx:1.22
```

再次檢查資源：
```bash
kubectl get deploy,rs,pods
```

現在您會看到：

1 個 Deployment（不變）
2 個 ReplicaSet（一個新的，一個舊的）
3 個 Pod（全部屬於新的 ReplicaSet）
這展示了 Deployment 如何通過創建新的 ReplicaSet 並逐步遷移 Pod 來實現無縫更新。

實驗室練習
基礎練習：創建和擴展 Deployment
創建一個運行 httpd:2.4 映像的 Deployment，初始副本數為 2
使用 kubectl scale 將副本數擴展到 4
觀察 Pod 的創建過程
將副本數縮減到 1
觀察 Pod 的終止過程
進階練習：Deployment 更新和回滾
創建一個運行 nginx:1.21 的 Deployment
更新 Deployment 使用 nginx:1.22
觀察滾動更新過程
回滾到原始版本
檢查更新歷史
挑戰練習：自定義更新策略
創建一個 Deployment，設置 maxSurge=2 和 maxUnavailable=0
更新 Deployment 的容器映像
觀察更新過程中的 Pod 數量變化
嘗試不同的 maxSurge 和 maxUnavailable 值，比較更新行為差異
故障排除
常見問題與解決方案
Deployment 卡在更新狀態：
Copy
kubectl rollout status deployment/nginx-deployment
# 如果卡住，可以檢查
kubectl describe deployment nginx-deployment
# 查看事件和條件
Pod 無法調度：
Copy
kubectl describe pod <pod-name>
# 查看事件，尋找調度失敗的原因
映像拉取失敗：
Copy
kubectl describe pod <pod-name>
# 查看事件，尋找 "Failed to pull image" 錯誤
重置卡住的更新：
Copy
kubectl rollout restart deployment/nginx-deployment
最佳實踐
始終使用資源限制：為 Pod 設置 CPU 和內存的請求和限制
使用適當的探針：添加存活探針和就緒探針
設置合理的更新策略：根據應用需求調整 maxSurge 和 maxUnavailable
使用標籤和註釋：合理組織和記錄資源
保持 YAML 文件版本控制：將配置文件納入版本控制系統
總結
今天，我們學習了 Kubernetes 中的 Deployment 和 ReplicaSet 資源：

Deployment 提供了聲明式的應用程序更新和管理能力
ReplicaSet 確保指定數量的 Pod 副本運行
Deployment 通過創建和管理 ReplicaSet 來實現滾動更新和回滾
我們可以輕鬆擴展和縮減應用程序，Kubernetes 會自動調整 Pod 數量
這些概念是 Kubernetes 應用程序管理的核心，掌握它們將使您能夠有效地部署和管理容器化應用程序。

