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
- Deployment: 管理 Pod 的複製和更新
- ReplicaSet: 確保指定數量的 Pod 副本運行
- StatefulSet: 管理有狀態應用的部署和擴展
- DaemonSet: 確保所有（或部分）節點運行一個 Pod 副本
- Job: 運行一次性任務
- CronJob: 按照時間表運行 Job

## 服務與網路資源 (Service & Networking)
- Service: 為一組 Pod 提供穩定的網路端點
- Ingress: 管理外部訪問服務的規則
- NetworkPolicy: 定義 Pod 之間的網路隔離策略

## 配置與儲存資源 (Config & Storage)
- ConfigMap: 存儲非機密配置數據
- Secret: 存儲機密數據
- Volume: 提供 Pod 的持久化存儲
- PersistentVolume: 集群級別的存儲資源
- PersistentVolumeClaim: 請求存儲資源

## 集群資源 (Cluster Resources)
- Node: 集群中的工作機器
- Namespace: 提供資源隔離機制
- ResourceQuota: 限制命名空間資源使用
- Role/ClusterRole: 定義權限規則
- RoleBinding/ClusterRoleBinding: 將角色綁定到用戶

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