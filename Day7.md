---
title: "Day 7：資源管理（Label、Selector、Namespace）"
tags: 2025鐵人賽
date: 2025-07-20
---

#### Day 7：資源管理（Label、Selector、Namespace）
##### 重點
- Label 與 Selector 用法、Namespace 管理
##### Lab
- 建立多個 Namespace，利用 Label 管理與查詢資源

在前幾篇文章中，我們探討了 Kubernetes 的基本資源和網路模型。今天，我們將聚焦於 Kubernetes 中的資源管理機制，特別是 Label、Selector 和 Namespace 這三個重要概念。這些機制是有效組織和管理 Kubernetes 資源的關鍵工具。

# Label 概念與用法
什麼是 Label？
Label 是附加到 Kubernetes 對象（如 Pod、Node、Service 等）上的鍵值對，用於組織和選擇資源子集。Label 不提供唯一性，多個對象可以有相同的 Label。

## Label 的特性
鍵值對格式：key=value
可自定義：用戶可以根據需求自由定義
可多個：一個資源可以有多個 Label
可修改：可以隨時添加、修改或刪除

## Label 的命名規則
Label 的命名有特定規則：

鍵（Key）：

可選前綴：如 example.com/，不超過 253 字符
名稱：必須以字母或數字開頭和結尾，可包含字母、數字、連字符、點和下劃線
不超過 63 字符
值（Value）：

必須以字母或數字開頭和結尾
可包含字母、數字、連字符、點和下劃線
不超過 63 字符

## Label 的常見用途
環境區分：environment=dev/staging/production
應用識別：app=frontend/backend/database
版本標記：version=v1/v2/latest
團隊歸屬：team=team-a/team-b
發布追踪：release=stable/canary

## 如何為資源添加 Label
可以在創建資源時添加 Label，或者為現有資源添加 Label。

創建時添加 Label：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    app: nginx
    environment: production
    tier: frontend
spec:
  containers:
  - name: nginx
    image: nginx:1.21
```

為現有資源添加 Label：
```bash
# 使用 kubectl label 命令
kubectl label pods nginx version=v1

# 添加多個 Label
kubectl label pods nginx team=frontend release=stable
```

修改和刪除 Label
修改 Label：
```bash
# 覆蓋現有 Label
kubectl label --overwrite pods nginx environment=staging
```

刪除 Label：
```bash
# 刪除 Label（注意 key 後面的減號）
kubectl label pods nginx environment-
```

# Selector 概念與用法
什麼是 Selector？
Selector 是一種機制，用於根據 Label 選擇 Kubernetes 資源。它是許多 Kubernetes 資源（如 Service、Deployment）的核心功能，用於確定哪些 Pod 是它們的目標。

## Selector 的類型
Kubernetes 支持兩種類型的 Selector：

1. 等值選擇器（Equality-based）：使用 =、== 或 != 運算符
2. 集合選擇器（Set-based）：使用 in、notin 和 exists 運算符

## 等值選擇器示例
```bash
# 選擇 environment=production 的所有 Pod
kubectl get pods -l environment=production

# 選擇 tier!=frontend 的所有 Pod
kubectl get pods -l tier!=frontend

# 選擇同時滿足多個條件的 Pod
kubectl get pods -l environment=production,tier=frontend
```

集合選擇器示例
```bash
# 選擇 environment 為 production 或 staging 的 Pod
kubectl get pods -l 'environment in (production,staging)'

# 選擇 tier 不是 frontend 或 backend 的 Pod
kubectl get pods -l 'tier notin (frontend,backend)'

# 選擇有 version 標籤的 Pod
kubectl get pods -l 'version'

# 選擇沒有 version 標籤的 Pod
kubectl get pods -l '!version'
```

## Selector 在資源定義中的應用
Selector 在多種 Kubernetes 資源中使用，如 Deployment、Service、ReplicaSet 等。

在 Deployment 中使用 Selector：
```bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
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

在 Service 中使用 Selector：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
```

高級 Selector 表達式
Deployment 和其他資源支持更複雜的 matchExpressions 選擇器：

```yaml
selector:
  matchExpressions:
    - {key: environment, operator: In, values: [production, staging]}
    - {key: tier, operator: NotIn, values: [frontend]}
```

支持的運算符包括：
- In：值必須匹配列表中的一個
- NotIn：值不能匹配列表中的任何一個
- Exists：資源必須包含指定的標籤鍵（不檢查值）
- DoesNotExist：資源不能包含指定的標籤鍵

# Namespace 概念與管理
什麼是 Namespace？
Namespace 是 Kubernetes 中的虛擬集群，它們在同一物理集群上提供了一種資源隔離機制。Namespace 可以將集群資源劃分為多個邏輯單元，適用於多租戶環境或將不同環境（開發、測試、生產）隔離。

默認 Namespace
Kubernetes 初始化時會創建幾個默認 Namespace：

default：默認 Namespace，未指定時使用
kube-system：Kubernetes 系統組件使用
kube-public：自動創建且所有用戶可讀，用於集群使用
kube-node-lease：用於節點心跳檢測


查看 Namespace
```
# 列出所有 Namespace
kubectl get namespaces

# 或簡寫
kubectl get ns
```

創建 Namespace
可以通過 YAML 文件或直接使用 kubectl 命令創建 Namespace。

使用 YAML 文件：
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: development
```

保存為 development-ns.yaml 並應用：
```bash
kubectl apply -f development-ns.yaml
```

使用 kubectl 命令：
```bash
kubectl create namespace production
```

在特定 Namespace 中操作
在使用 kubectl 時，可以通過 --namespace 或 -n 選項指定 Namespace：

```bash
# 在 production Namespace 中創建資源
kubectl apply -f nginx.yaml --namespace=production

# 查看特定 Namespace 中的資源
kubectl get pods -n production
```

也可以在資源定義中指定 Namespace：
```bash
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  namespace: production
spec:
  containers:
  - name: nginx
    image: nginx:1.21
```

設置默認 Namespace
可以更改 kubectl 的默認 Namespace，避免每次都要指定：

```bash
# 設置當前上下文的默認 Namespace
kubectl config set-context --current --namespace=production

# 驗證設置
kubectl config view --minify | grep namespace:
```

刪除 Namespace
```bash
kubectl delete namespace development
```

刪除 Namespace 會刪除其中的所有資源，請謹慎操作。


Namespace 的限制
並非所有 Kubernetes 資源都受 Namespace 約束。一些資源（如 Node、PersistentVolume）是集群級別的，不屬於任何 Namespace。

查看哪些資源受 Namespace 約束：

```bash
# 查看受 Namespace 約束的資源
kubectl api-resources --namespaced=true

# 查看不受 Namespace 約束的資源
kubectl api-resources --namespaced=false
```

# 實作：多 Namespace 環境與 Label 管理
現在，讓我們通過實際操作來理解 Label、Selector 和 Namespace 的使用。

步驟 1：創建多個 Namespace
我們將創建三個 Namespace，代表不同的環境：

```bash
# 創建開發環境 Namespace
kubectl create namespace dev

# 創建測試環境 Namespace
kubectl create namespace staging

# 創建生產環境 Namespace
kubectl create namespace prod
```

確認 Namespace 創建成功：

```bash
kubectl get namespaces
```

步驟 2：在不同 Namespace 中部署應用
我們將在每個 Namespace 中部署一個簡單的 Nginx 應用，並添加不同的 Label。

開發環境：

```yaml
cat <<EOF | kubectl apply -f - -n dev
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
    environment: development
    version: v1
    team: frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
        environment: development
        version: v1
        team: frontend
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
    environment: development
spec:
  selector:
    app: nginx
    environment: development
  ports:
  - port: 80
    targetPort: 80
EOF
```

測試環境：
```bash
cat <<EOF | kubectl apply -f - -n staging
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
    environment: staging
    version: v1
    team: frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
        environment: staging
        version: v1
        team: frontend
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
    environment: staging
spec:
  selector:
    app: nginx
    environment: staging
  ports:
  - port: 80
    targetPort: 80
EOF
```

生產環境：
```bash
cat <<EOF | kubectl apply -f - -n prod
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
    environment: production
    version: stable
    team: frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
        environment: production
        version: stable
        team: frontend
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
    environment: production
spec:
  selector:
    app: nginx
    environment: production
  ports:
  - port: 80
    targetPort: 80
EOF
```

步驟 3：使用 Label 查詢和管理資源
現在，讓我們使用 Label 來查詢和管理這些資源。

查看所有 Namespace 中的 Nginx Deployment：
```bash
kubectl get deployments --all-namespaces -l app=nginx
```
查看特定環境的 Pod：
```bash
# 查看開發環境的 Pod
kubectl get pods --all-namespaces -l environment=development

# 查看生產環境的 Pod
kubectl get pods --all-namespaces -l environment=production
```

查看特定版本的 Deployment：

```bash
# 查看 v1 版本的 Deployment
kubectl get deployments --all-namespaces -l version=v1

# 查看 stable 版本的 Deployment
kubectl get deployments --all-namespaces -l version=stable
```

使用多個條件查詢：
```bash
# 查看開發環境中的 frontend 團隊資源
kubectl get all --all-namespaces -l environment=development,team=frontend

# 查看非生產環境的所有服務
kubectl get services --all-namespaces -l 'environment in (development,staging)'
```

步驟 4：修改資源的 Label
讓我們修改一些資源的 Label，以模擬版本升級或環境變更。

將開發環境的應用升級到 v2：
```bash
# 更新 Deployment 的 Label
kubectl label deployment nginx version=v2 --overwrite -n dev

# 更新 Pod 模板的 Label（需要編輯 Deployment）
kubectl edit deployment nginx -n dev
```


在編輯器中，找到 spec.template.metadata.labels 部分，將 version: v1 改為 version: v2。

添加新的 Label：
```bash
# 為測試環境添加 release 標籤
kubectl label deployment nginx release=beta -n staging

# 為生產環境添加 release 標籤
kubectl label deployment nginx release=stable -n prod
```

使用新 Label 查詢：
```bash
# 查看 beta 發布的資源
kubectl get all --all-namespaces -l release=beta

# 查看 v2 版本的資源
kubectl get all --all-namespaces -l version=v2
```

步驟 5：使用 Namespace 隔離資源
現在，讓我們嘗試在不同的 Namespace 中創建同名資源，以展示 Namespace 的隔離作用。

在每個 Namespace 中創建一個 ConfigMap：

```bash
# 在開發環境創建 ConfigMap
kubectl create configmap app-config -n dev --from-literal=ENV=development

# 在測試環境創建同名 ConfigMap
kubectl create configmap app-config -n staging --from-literal=ENV=staging

# 在生產環境創建同名 ConfigMap
kubectl create configmap app-config -n prod --from-literal=ENV=production
```

查看不同 Namespace 中的 ConfigMap：

```bash
# 列出所有 Namespace 中的 ConfigMap
kubectl get configmaps --all-namespaces

# 查看開發環境的 ConfigMap 內容
kubectl describe configmap app-config -n dev

# 查看生產環境的 ConfigMap 內容
kubectl describe configmap app-config -n prod
```
可以看到，雖然 ConfigMap 名稱相同，但它們在不同的 Namespace 中是完全隔離的，內容也不同。

# Label 和 Namespace 的最佳實踐
Label 最佳實踐
制定標準化的標籤方案：

為組織定義一套標準的標籤鍵和值
記錄並共享這些標準，確保團隊一致使用
使用有意義的標籤：

選擇描述性強的標籤名稱
避免過於籠統或過於具體的標籤
避免過度使用：

只使用真正需要的標籤
太多標籤會增加管理複雜性
考慮使用前綴：

對於組織特定的標籤，考慮使用域名前綴
例如：example.com/environment 而不是 environment
定期審查和更新：

定期檢查標籤使用情況
移除過時或不再使用的標籤
Namespace 最佳實踐
按環境或團隊劃分：

為不同環境（開發、測試、生產）創建獨立的 Namespace
或按團隊/項目劃分 Namespace
設置資源配額：

為每個 Namespace 設置 ResourceQuota，限制資源使用
防止一個 Namespace 消耗過多集群資源
使用網路策略：

實施 NetworkPolicy 控制 Namespace 間的通信
增強安全性和隔離性
設置默認限制：

使用 LimitRange 為 Namespace 中的容器設置默認資源限制
確保所有工作負載都有合理的資源限制
權限管理：

使用 RBAC 控制對 Namespace 的訪問
為不同團隊分配適當的權限

# 總結
在本文中，我們深入探討了 Kubernetes 中的資源管理機制：Label、Selector 和 Namespace。這些機制是有效組織和管理 Kubernetes 資源的關鍵工具。

Label 和 Selector 提供了一種靈活的方式來組織和選擇資源，使我們能夠按環境、應用、版本等維度管理資源。Namespace 則提供了資源隔離的機制，使不同團隊或環境的資源可以共存於同一集群中。

通過實際操作，我們學習了如何創建和管理多個 Namespace，如何使用 Label 來標記和選擇資源，以及如何利用這些機制來實現更有效的資源管理。

在下一篇文章中，我們將探討 Kubernetes 的配置管理，包括 ConfigMap 和 Secret 資源，它們是管理應用配置和敏感數據的重要工具。