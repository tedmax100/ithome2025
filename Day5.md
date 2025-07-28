---
title: "Day 5：Service 與 Kubernetes 網路模型"
tags: 2025鐵人賽
date: 2025-07-20
---


#### Day 5：Service 與 Kubernetes 網路模型
##### 重點
Service 類型（ClusterIP、NodePort、LoadBalancer）
Kubernetes 網路模型基礎
KinD 內建的 CNI 介紹

##### Lab
部署 NodePort Service，從本機 curl 訪問服務
測試不同 Service 類型的行為
觀察網路封包流動


在前面的文章中，我們已經了解了 Pod 和 Deployment 等基本資源。今天，我們將深入探討 Kubernetes 的網路模型和 Service 資源，這是讓應用能夠被訪問和發現的關鍵組件。

# Kubernetes 網路模型基礎
Kubernetes 網路模型基於以下四個基本原則：

1. Pod 間通信：所有 Pod 可以不通過 NAT 直接通信
2. Node 與 Pod 通信：Node 上的代理程序（如 kubelet）可以與該 Node 上的所有 Pod 通信
3. Pod 自身網路：Pod 內的容器共享網路命名空間，可以通過 localhost 互相訪問
4. 外部與 Pod 通信：需要通過 Service 或 Ingress 等資源實現

這種網路模型與傳統的 Docker 網路有很大不同。在 Docker 中，容器默認使用橋接網路，容器間通信需要通過端口映射或自定義網路。而 Kubernetes 為每個 Pod 分配一個集群內唯一的 IP 地址，使 Pod 之間可以直接通信，就像傳統的虛擬機一樣。


# Pod 網路實現的挑戰
實現 Kubernetes 網路模型面臨幾個主要挑戰：

跨節點通信：Pod 可能分布在不同節點上，需要解決跨節點通信問題
IP 地址管理：需要為每個 Pod 分配唯一的 IP 地址
路由與轉發：確保數據包能正確路由到目標 Pod
網路策略：實現網路隔離和安全控制

# Container Network Interface (CNI)
Kubernetes 使用 CNI（Container Network Interface）插件來實現其網路模型。CNI 是一個規範和一組庫，用於配置 Linux 容器的網路接口。

在 KinD（Kubernetes in Docker）中，默認使用的 CNI 插件是 kindnet，這是一個基於 ptp (point-to-point) 和 host-local IPAM 插件的簡單 CNI 實現。它為每個節點分配一個子網，並使用 iptables 規則實現跨節點的 Pod 通信。

其他常見的 CNI 插件包括：

Calico：提供高性能、可擴展的網路解決方案，支持網路策略
Flannel：簡單的覆蓋網路，專注於 Kubernetes 的網路連接
Cilium：基於 eBPF 的網路和安全解決方案，提供 L7 策略支持
Weave Net：創建跨節點的虛擬網路，支持加密通信

# Service 資源介紹
雖然 Pod 有自己的 IP 地址，但這些 IP 是不穩定的——當 Pod 重啟或被替換時，IP 地址會改變。此外，當使用 Deployment 時，可能同時運行多個 Pod 副本。這就需要一種機制來抽象這些變化，提供穩定的網路端點，這就是 Service 的作用。

Service 是 Kubernetes 中的一個抽象概念，它定義了一組 Pod 的邏輯集合和這些 Pod 的策略。Service 通過標籤選擇器（Label Selector）來確定哪些 Pod 屬於它，並為這些 Pod 提供統一的訪問入口。

## Service 的基本 YAML 結構
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
  - port: 80        # Service 暴露的端口
    targetPort: 9376 # Pod 上的目標端口
  type: ClusterIP   # Service 類型
```

## Service 的主要類型
Kubernetes 提供了多種類型的 Service，以滿足不同的網路需求：

1. ClusterIP（默認類型）
ClusterIP 是最基本的 Service 類型，它在集群內部提供一個穩定的 IP 地址，僅在集群內部可訪問。
分配一個集群內部的虛擬 IP
只能在集群內部訪問
適用於內部微服務通信

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-clusterip-service
spec:
  selector:
    app: MyApp
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

用途：集群內部服務間的通信，如前端訪問後端 API。

2. NodePort
NodePort 在 ClusterIP 的基礎上，在每個節點上開放一個靜態端口（默認範圍：30000-32767），通過 <NodeIP>:<NodePort> 可以從集群外部訪問服務。

在 ClusterIP 基礎上，在每個節點上開放一個靜態端口
可以通過 <NodeIP>:<NodePort> 從集群外部訪問
端口範圍通常為 30000-32767
適用於開發和測試環境

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-nodeport-service
spec:
  selector:
    app: MyApp
  ports:
  - port: 80        # 集群內部訪問端口
    targetPort: 8080 # Pod 端口
    nodePort: 30007  # 節點端口（可選，不指定則自動分配）
  type: NodePort
```
用途：開發環境或小型生產環境中暴露服務，便於直接訪問。


3. LoadBalancer
LoadBalancer 在 NodePort 的基礎上，自動創建外部負載均衡器（需要雲提供商支持），並將流量轉發到 Service。

在 NodePort 基礎上，創建一個外部負載均衡器
自動分配一個外部 IP，將流量轉發到 Service
需要雲提供商支持或使用 MetalLB 等工具
適用於生產環境的外部訪問

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-loadbalancer-service
spec:
  selector:
    app: MyApp
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```
用途：在雲環境中暴露服務，如公開的 Web 應用。

4. ExternalName
ExternalName 將服務映射到一個外部 DNS 名稱，而不是選擇 Pod。
   
將 Service 映射到一個 DNS 名稱
不使用選擇器和集群 IP
通過返回 CNAME 記錄實現

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-externalname-service
spec:
  type: ExternalName
  externalName: my.database.example.com

```

用途：將外部服務引入集群內部，如外部數據庫。


## Service 發現機制
Kubernetes 提供兩種服務發現機制：

環境變量：當 Pod 啟動時，Kubernetes 會為每個活動的 Service 設置環境變量
DNS：Kubernetes 集群中的 DNS 服務（如 CoreDNS）為每個 Service 創建 DNS 記錄
使用 DNS 是推薦的方式，Pod 可以通過 <service-name>.<namespace>.svc.cluster.local 訪問服務。在同一命名空間中，可以直接使用 <service-name>。


# Service 實戰：部署和測試不同類型的 Service
現在讓我們通過實際操作來理解不同類型的 Service。我們將部署一個簡單的 Web 應用，並使用不同類型的 Service 來訪問它。

步驟 1：部署一個簡單的 Web 應用
首先，讓我們創建一個名為 web-deployment.yaml 的文件：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
  labels:
    app: web
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
        ports:
        - containerPort: 80
```

使用 kubectl apply 部署這個應用：


kubectl apply -f web-deployment.yaml
確認 Pod 已經運行：


kubectl get pods -l app=web
輸出應該顯示 3 個運行中的 Pod：


NAME                       READY   STATUS    RESTARTS   AGE
web-app-6d8d8c4f7d-5bjvx   1/1     Running   0          30s
web-app-6d8d8c4f7d-8xjqp   1/1     Running   0          30s
web-app-6d8d8c4f7d-qpvgz   1/1     Running   0          30s

步驟 2：創建 ClusterIP Service
現在，讓我們創建一個 ClusterIP 類型的 Service，這是默認類型：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-clusterip
spec:
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP  # 可以省略，因為這是默認類型
```

保存為 web-clusterip.yaml 並應用：


kubectl apply -f web-clusterip.yaml
查看 Service：


kubectl get svc web-clusterip
輸出應該類似於：


NAME           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
web-clusterip  ClusterIP   10.96.123.456   <none>        80/TCP    30s
測試 ClusterIP Service（在集群內部）：


# 創建一個臨時 Pod 來測試
kubectl run temp --rm -it --image=busybox -- sh

# 在臨時 Pod 內部執行
wget -O- web-clusterip
# 或者使用 IP 地址
wget -O- 10.96.123.456
你應該能看到 Nginx 的歡迎頁面 HTML。

步驟 3：創建 NodePort Service
接下來，讓我們創建一個 NodePort 類型的 Service：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-nodeport
spec:
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30080  # 指定節點端口，可以省略讓 Kubernetes 自動分配
  type: NodePort
```

保存為 web-nodeport.yaml 並應用：


kubectl apply -f web-nodeport.yaml
查看 Service：


kubectl get svc web-nodeport
輸出應該類似於：


NAME           TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
web-nodeport   NodePort   10.96.234.567   <none>        80:30080/TCP   30s


在 KinD 環境中測試 NodePort Service：

由於 KinD 是在 Docker 容器中運行的 Kubernetes 集群，我們需要先獲取節點容器的 IP：
```bash
> docker container ls
```

找到 KinD 控制平面節點（通常名為 kind-control-plane）：


CONTAINER ID   IMAGE                  COMMAND                  CREATED       STATUS       PORTS                       NAMES
abcdef123456   kindest/node:v1.24.0   "/usr/local/bin/entr…"   2 hours ago   Up 2 hours   127.0.0.1:50443->6443/tcp   kind-control-plane

現在我們可以使用 curl 訪問 NodePort Service：

```bash
> curl http://localhost:30080
```

如果你使用的是 Mac 或 Windows，KinD 會自動將節點端口映射到主機，所以可以直接訪問 localhost。如果你使用的是 Linux，可能需要使用節點容器的 IP 地址：


docker container inspect kind-control-plane -f '{{.NetworkSettings.Networks.kind.IPAddress}}'

然後使用這個 IP 地址訪問：

curl http://<node-ip>:30080
你應該能看到 Nginx 的歡迎頁面 HTML。

步驟 4：模擬 LoadBalancer Service
在雲環境中，LoadBalancer 類型的 Service 會自動創建一個外部負載均衡器。但在本地 KinD 環境中，我們需要使用 MetalLB 等工具來模擬這個行為。由於設置 MetalLB 相對複雜，我們可以使用 port-forward 來模擬 LoadBalancer 的行為：

```bash
> kubectl port-forward svc/web-clusterip 8080:80
```

現在你可以在瀏覽器中訪問 http://localhost:8080 或使用 curl http://localhost:8080


觀察網路封包流動
要深入了解 Kubernetes 網路模型，我們可以觀察網路封包的流動。在 KinD 環境中，我們可以進入節點容器並使用 tcpdump 工具來捕獲網路流量。

首先，讓我們進入控制平面節點：
```bash
docker exec -it kind-control-plane bash
```

安裝 tcpdump（如果尚未安裝）：

apt-get update && apt-get install -y tcpdump


現在我們可以捕獲特定接口上的流量。例如，捕獲 eth0 接口上的 HTTP 流量：

tcpdump -i eth0 port 80 -n

在另一個終端中，使用 curl 訪問我們的 NodePort Service：


curl http://localhost:30080


你應該能在 tcpdump 輸出中看到相關的網路封包。

Service 網路流量分析
讓我們分析一下當我們訪問 Service 時，網路流量是如何流動的：

1. ClusterIP Service:
   - 客戶端 -> kube-proxy（iptables/IPVS 規則）-> 目標 Pod

2. NodePort Service:
   - 外部客戶端 -> 節點 IP:NodePort -> kube-proxy -> 目標 Pod

3. LoadBalancer Service:
   - 外部客戶端 -> 負載均衡器 IP -> 節點 IP:NodePort -> kube-proxy -> 目標 Pod

在 Kubernetes 中，kube-proxy 組件負責實現 Service 的網路功能。它有三種模式：
- userspace 模式（最早的實現，性能較差）
- iptables 模式（默認模式，使用 Linux iptables 規則）
- IPVS 模式（更高效的實現，適用於大規模集群）

在 KinD 中，默認使用 iptables 模式。我們可以查看節點上的 iptables 規則來了解 Service 是如何實現的：
```bash
# 在節點容器內執行
iptables-save | grep web-clusterip
```

你應該能看到一系列與我們的 ClusterIP Service 相關的 iptables 規則。


## KinD 內建的 CNI 介紹
KinD 默認使用 kindnet 作為 CNI 插件，這是一個簡單的 CNI 實現，專為 KinD 設計。它使用 ptp (point-to-point) 插件為每個 Pod 創建一個虛擬網絡接口，並使用 host-local IPAM 插件分配 IP 地址。它通過在 Docker 橋接網路上實現 Kubernetes 網路模型，使 Pod 可以跨節點通信。

kindnet 的主要特點：

- 為每個節點分配一個子網
- 使用 iptables 規則實現跨節點的 Pod 通信
- 簡單輕量，適合測試和開發環境
- 不支持網路策略（Network Policy）

我們可以查看 KinD 集群中的 CNI 配置：

```bash
# 在節點容器內執行
cat /etc/cni/net.d/10-kindnet.conflist
```

輸出應該類似於：
```json
{
  "cniVersion": "0.3.1",
  "name": "kindnet",
  "plugins": [
    {
      "type": "ptp",
      "ipMasq": false,
      "ipam": {
        "type": "host-local",
        "dataDir": "/run/cni-ipam-state",
        "routes": [
          { "dst": "10.244.0.0/16" }
        ],
        "ranges": [
          [
            {
              "subnet": "10.244.1.0/24"
            }
          ]
        ]
      }
    },
    {
      "type": "portmap",
      "capabilities": {
        "portMappings": true
      }
    }
  ]
}
```

這個配置顯示 kindnet 使用 ptp 插件和 host-local IPAM 插件，為 Pod 分配 10.244.1.0/24 子網中的 IP 地址。


# 總結
今天我們學習了 Kubernetes 的網路模型和 Service 資源。我們了解了不同類型的 Service（ClusterIP、NodePort、LoadBalancer）及其用途，並通過實際操作部署和測試了這些 Service。我們還深入探討了 Kubernetes 網路模型的原理，以及 KinD 中使用的 CNI 插件。

Kubernetes 的網路模型設計得非常靈活和強大，能夠滿足各種複雜的網路需求。Service 資源則提供了一種抽象機制，使應用能夠被穩定地訪問和發現，無論底層 Pod 如何變化。

在下一篇文章中，我們將探討 Ingress 資源，它提供了更高級的 HTTP 路由功能，是構建現代微服務架構的重要組件。

