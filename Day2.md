---
title: "Day 2 Kubernetes 資源與 YAML 基礎"
tags: 2025鐵人賽
date: 2025-07-20
---
#### Kubernetes 資源與 YAML 基礎

##### 重點

- Kubernetes 資源種類、YAML 結構說明

##### Lab

- 撰寫一個簡單的 Pod YAML，kubectl apply 部署並觀察狀態

----
在上一篇文章中，我們成功安裝了 KinD 並啟動了一個本地 Kubernetes 叢集。今天，我們將深入了解 Kubernetes 資源的種類以及如何使用 YAML 來定義這些資源。

# Kubernetes 資源概述
Kubernetes 資源(Resources) 是指系統中可以建立、修改和管理的對象。每種資源都有特定的用途和行為。以下是一些核心資源類型：

## 工作負載資源 (Workload Resources)
- Pod: 最小的可部署單位，包含一個或多個容器
  - Docker 對比: 類似於 Docker 容器，但 Pod 可以包含多個容器，這些容器共享網絡和存儲命名空間，類似於 Docker 中的 --net=container 和 --volumes-from 功能。
- Deployment: 管理 Pod 的複製和更新
  - Docker 對比: 類似於 Docker Compose 中的服務定義，但增加了自動擴展、滾動更新等功能。
- ReplicaSet: 確保指定數量的 Pod 副本運行
  - Docker 對比: 相當於 Docker Compose 中的 scale 功能，但更加自動化和健壯。
- StatefulSet: 管理有狀態應用的部署和擴展
  - Docker 對比: 在 Docker Compose 中沒有直接對應，但類似於為有狀態服務(如數據庫)定義的服務配置，增加了穩定網絡標識和持久存儲。
- DaemonSet: 確保所有（或部分）節點運行一個 Pod 副本
- Job: 運行一次性任務
- CronJob: 按照時間表運行 Job

## 服務與網路資源 (Service & Networking)
- Service: 為一組 Pod 提供穩定的網路端點
  - Docker 對比: 類似於 Docker Compose 中的服務發現和端口映射，但更加強大和靈活。
- Ingress: 管理外部訪問服務的規則
- NetworkPolicy: 定義 Pod 之間的網路隔離策略

## 配置與儲存資源 (Config & Storage)
- ConfigMap: 存儲非機密配置數據
  - Docker 對比: 類似於 Docker 中的環境變量或 Docker Compose 中的 environment 和 env_file 配置。
- Secret: 存儲機密數據
  - Docker 對比: 類似於 Docker 的環境變量，但專為敏感數據設計，有加密功能。
- Volume: 提供 Pod 的持久化存儲
  - Docker 對比: 直接對應 Docker 的 volumes 和 bind mounts。
- PersistentVolume: 集群級別的存儲資源
- PersistentVolumeClaim: 請求存儲資源

## 集群資源 (Cluster Resources)
- Node: 集群中的工作機器
- Namespace: 提供資源隔離機制
- ResourceQuota: 限制命名空間資源使用
- Role/ClusterRole: 定義權限規則
- RoleBinding/ClusterRoleBinding: 將角色綁定到用戶

```meramid
graph TD
  subgraph "集群資源"
    A[Node] -->|"運行於"| B[Namespace]
    B -->|"限制資源使用"| C[ResourceQuota]
    B -->|"定義權限"| D[Role/ClusterRole]
    B -->|"關聯用戶與權限"| E[RoleBinding/ClusterRoleBinding]
  end
  
  subgraph "工作負載資源"
    F[Pod] -->|"包含"| G[Container]
    H[Deployment] -->|"創建和管理"| I[ReplicaSet]
    I -->|"創建和管理"| F
    J[StatefulSet] -->|"直接管理"| F
    K[DaemonSet] -->|"在每個節點創建"| F
    L[Job] -->|"創建一次性"| F
    M[CronJob] -->|"定時創建"| L
  end
  
  subgraph "服務與網路"
    N[Service] -->|"暴露和負載均衡"| F
    O[Ingress] -->|"提供外部訪問"| N
    P[NetworkPolicy] -->|"控制網路流量"| F
  end
  
  subgraph "配置與儲存"
    Q[ConfigMap] -->|"提供配置"| F
    R[Secret] -->|"提供敏感配置"| F
    S[Volume] -->|"提供存儲"| F
    T[PersistentVolume] -->|"分配給"| U[PersistentVolumeClaim]
    U -->|"掛載到"| F
  end
  
  A ---|"調度和運行"| F
```

集群資源關係
Node → Namespace：「運行於」- Node 是實體機器，而 Namespace 是邏輯隔離單位，Node 上運行著屬於各個 Namespace 的資源。

Namespace → ResourceQuota：「限制資源使用」- ResourceQuota 定義了 Namespace 中可以使用的資源總量限制。

Namespace → Role/ClusterRole：「定義權限」- Role 定義了在特定 Namespace 中的權限，ClusterRole 則定義了集群範圍的權限。

Namespace → RoleBinding/ClusterRoleBinding：「關聯用戶與權限」- 這些資源將用戶或服務帳戶與角色綁定，授予相應權限。

工作負載資源關係
Pod → Container：「包含」- Pod 是最小部署單位，包含一個或多個容器。

Deployment → ReplicaSet：「創建和管理」- Deployment 創建 ReplicaSet 並管理其生命週期，實現滾動更新等功能。

ReplicaSet → Pod：「創建和管理」- ReplicaSet 確保指定數量的 Pod 副本運行。

StatefulSet → Pod：「直接管理」- StatefulSet 直接管理 Pod，提供穩定的網絡標識和持久存儲。

DaemonSet → Pod：「在每個節點創建」- DaemonSet 確保每個節點上運行一個 Pod 副本，適合節點監控、日誌收集等場景。

Job → Pod：「創建一次性」- Job 創建 Pod 執行一次性任務，任務完成後 Pod 終止。

CronJob → Job：「定時創建」- CronJob 按照時間表定期創建 Job。

服務與網路關係
Service → Pod：「暴露和負載均衡」- Service 為一組 Pod 提供穩定的網絡端點和負載均衡。

Ingress → Service：「提供外部訪問」- Ingress 定義從集群外部訪問 Service 的規則。

NetworkPolicy → Pod：「控制網路流量」- NetworkPolicy 定義 Pod 間的網絡隔離策略。

配置與儲存關係
ConfigMap → Pod：「提供配置」- ConfigMap 為 Pod 提供非敏感的配置數據。

Secret → Pod：「提供敏感配置」- Secret 為 Pod 提供敏感的配置數據，如密碼、令牌等。

Volume → Pod：「提供存儲」- Volume 為 Pod 提供臨時或持久的存儲空間。

PersistentVolume → PersistentVolumeClaim：「分配給」- PersistentVolume 是集群中的存儲資源，通過 PersistentVolumeClaim 分配給 Pod。

PersistentVolumeClaim → Pod：「掛載到」- PersistentVolumeClaim 請求存儲資源並掛載到 Pod 中。

核心關係
Node ↔ Pod：「調度和運行」- Kubernetes 調度器將 Pod 分配到 Node 上運行，Node 提供計算資源。

# YAML 在 Kubernetes 中的應用
Kubernetes 使用 YAML (或 JSON) 格式的檔案來定義資源，並使用 kubectl apply 命令部署到叢集中。

## Kubernetes 資源 YAML 結構
所有 Kubernetes 資源 YAML 都有以下基本結構：
```yaml
apiVersion: <API 版本>
kind: <資源類型>
metadata:
  name: <資源名稱>
  namespace: <命名空間>
  labels:
    <標籤鍵>: <標籤值>
spec:
  # 資源特定配置
```

apiVersion: 使用的 Kubernetes API 版本
kind: 要創建的資源類型
metadata: 幫助識別對象的數據，如名稱、UID、命名空間等標籤
spec: 資源的期望狀態

```yaml
apiVersion: v1        # Kubernetes API 版本
kind: Pod             # 資源類型
metadata:             # 元數據
  name: nginx-pod     # 資源名稱
  labels:             # 標籤
    app: nginx
spec:                 # 規格/期望狀態
  containers:         # 容器定義
  - name: nginx       # 容器名稱
    image: nginx:1.14 # 容器映像
    ports:            # 容器端口
    - containerPort: 80
```

### Docker Compose YAML vs Kubernetes YAML
```yaml
# Docker Compose YAML
version: '3'
services:
  nginx:
    image: nginx:1.21
    ports:
      - "80:80"
    volumes:
      - ./html:/usr/share/nginx/html
```

```yaml
# Kubernetes Pod YAML
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.21
    ports:
    - containerPort: 80
    volumeMounts:
    - name: html-volume
      mountPath: /usr/share/nginx/html
  volumes:
  - name: html-volume
    hostPath:
      path: ./html
```
主要區別：

Docker Compose 使用 services 作為頂層鍵，而 Kubernetes 使用 kind 指定資源類型
Kubernetes 需要 apiVersion 指定 API 版本
Kubernetes 的 metadata 包含更多識別和分類信息
Kubernetes 將容器定義放在 spec.containers 下，而不是頂層

#### 多容器設計模式
在 Docker Compose 中，通常每個服務都是獨立的容器。而在 Kubernetes 中，Pod 可以包含多個容器，這些容器共享網路和儲存空間，類似於 Docker 中的 network 和 volume 共享功能：

```yaml
# Kubernetes 多容器 Pod
apiVersion: v1
kind: Pod
metadata:
  name: multi-container-pod
spec:
  containers:
  - name: nginx
    image: nginx:1.21
  - name: log-collector
    image: busybox
    command: ["/bin/sh", "-c", "tail -f /var/log/nginx/access.log"]
    volumeMounts:
    - name: logs
      mountPath: /var/log/nginx
  volumes:
  - name: logs
    emptyDir: {}
```


# 實作：撰寫和部署一個簡單的 Pod YAML
現在，讓我們創建一個簡單的 Pod YAML 文件，並使用 kubectl apply 命令部署它。

## 步驟 1: 創建 Pod YAML 文件
創建一個名為 nginx-pod.yaml 的文件，內容如下：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.21
    ports:
    - containerPort: 80
```
這個 YAML 文件定義了一個運行 Nginx 1.21 版本的 Pod，並暴露 80 端口。


## 步驟 2: 部署 Pod 到 KinD 叢集
首先確保 KinD 叢集正在運行：
```bash
> kind get clusters
# 如果沒有運行中的叢集，創建一個
# kind create cluster --name demo
```

```bash
> kubectl apply -f nginx-pod.yaml
pod/nginx-pod created
```

## 步驟 3: 驗證 Pod 狀態
```bash
> kubectl get pods
NAME        READY   STATUS    RESTARTS   AGE
nginx-pod   1/1     Running   0          30s
```
如果 Pod 狀態為 Running，表示部署成功。



## 步驟 4: 查看 Pod 詳細資訊
這會顯示 Pod 的詳細信息，包括：

Pod 的基本信息（名稱、命名空間、節點等）
容器的狀態和詳情
事件日誌（可以幫助排查問題）

```bash
> kubectl describe pod nginx-pod
Name:         nginx-pod
Namespace:    default
Priority:     0
Node:         demo-control-plane/172.18.0.2
Start Time:   Wed, 23 Jul 2025 10:15:30 +0800
Labels:       app=nginx
Status:       Running
IP:           10.244.0.5
IPs:
  IP:  10.244.0.5
Containers:
  nginx:
    Container ID:   containerd://...
    Image:          nginx:1.21
    Image ID:       docker.io/library/nginx@sha256:...
    Port:           80/TCP
    Host Port:      0/TCP
    State:          Running
      Started:      Wed, 23 Jul 2025 10:15:35 +0800
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-... (ro)
...
```

## 步驟 5: 訪問 Pod 中的 Nginx 服務
我們可以使用 port-forward 功能將 Pod 的端口轉發到本地：
```bash
> kubectl port-forward nginx-pod 8080:80
Forwarding from 127.0.0.1:8080 -> 80
```

現在，我們可以在瀏覽器中訪問 http://localhost:8080 或使用 curl 命令：
```bash
> curl http://localhost:8080
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
...
```

## 步驟 6: 查看 Pod 日誌
```bash
> kubectl logs nginx-pod
```
這將顯示 Nginx 容器的日誌輸出。


## 步驟 7: 進入 Pod 容器
```bash
> kubectl exec -it nginx-pod -- /bin/bash
root@nginx-pod:/# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
root@nginx-pod:/# exit
```

## 步驟 8: 刪除 Pod
完成實驗後，我們可以刪除 Pod：

```bash
> kubectl delete pod nginx-pod
pod "nginx-pod" deleted
```

或者使用 YAML 文件刪除：

```bash
> kubectl delete -f nginx-pod.yaml
pod "nginx-pod" deleted
```

## 多容器 Pod 示例
Pod 可以包含多個容器，這些容器共享網絡和存儲空間。以下是一個多容器 Pod 的示例：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: multi-container-pod
  labels:
    app: web
spec:
  containers:
  - name: nginx
    image: nginx:1.21
    ports:
    - containerPort: 80
  - name: log-collector
    image: busybox
    command: ["/bin/sh", "-c", "while true; do echo 'Collecting logs...' >> /var/log/collect.log; sleep 10; done"]
    volumeMounts:
    - name: shared-logs
      mountPath: /var/log
  volumes:
  - name: shared-logs
    emptyDir: {}

```

在這個例子中：

第一個容器運行 Nginx 服務器
第二個容器運行一個簡單的日誌收集器
兩個容器共享一個名為 shared-logs 的 volume

# 總結
今天我們學習了 Kubernetes 的基本資源類型和 YAML 配置文件的結構。我們通過實際操作，創建了一個簡單的 Pod，並使用 kubectl 命令進行了部署、觀察和管理。

YAML 文件是與 Kubernetes 交互的主要方式，理解其結構和語法對於有效管理 Kubernetes 資源至關重要。在接下來的課程中，我們將探索更複雜的資源類型和配置，如 Deployment、Service 等，以及如何組合使用它們來部署完整的應用。

記住，Kubernetes 的強大之處在於其聲明式 API——你只需要描述你想要的狀態，Kubernetes 會負責實現並維護這個狀態。

下一篇文章中，我們將探討 Deployment 和 ReplicaSet，了解如何管理 Pod 的複製和更新。