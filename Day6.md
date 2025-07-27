---
title: "Day 6: Flannel CNI 深入解析"
tags: 2025鐵人賽
date: 2025-07-20
---

#### Day 6: Flannel CNI 深入解析
##### 重點
Kubernetes CNI 概念與架構
Flannel CNI 工作原理與組件
網路模型與封包流向
各種後端模式比較

##### Lab
建立使用 Flannel CNI 的自訂 KinD 叢集
測試 Pod 間通訊
分析 Flannel 網路行為


在前一篇文章中，我們介紹了 Kubernetes 的網路模型和 Service 資源。今天，我們將深入探討 Kubernetes 的網路實現，特別是 Flannel 這個流行的容器網路介面 (CNI) 插件。

# Kubernetes CNI 概念與架構
## CNI 是什麼？
CNI (Container Network Interface) 是一個規範和一組庫，用於配置 Linux 容器的網路介面。Kubernetes 使用 CNI 來設置 Pod 網路，允許不同的網路解決方案以插件形式提供。

## CNI 的核心概念
1. 網路命名空間 (Network Namespace): Linux 內核功能，提供網路隔離
2. CNI 插件: 實現 CNI 規範的可執行文件，負責配置容器網路
3. CNI 配置: 定義網路如何設置的 JSON 文件

### CNI 在 Kubernetes 中的角色
在 Kubernetes 中，當 Pod 被調度到節點上時，kubelet 會：

1. 為 Pod 創建網路命名空間
2. 調用 CNI 插件配置 Pod 的網路
3. CNI 插件分配 IP 地址並設置路由規則

這個過程確保了 Kubernetes 網路模型的四個基本原則：

- 所有 Pod 可以不經 NAT 互相通信
- 所有節點可以不經 NAT 與所有 Pod 通信
- Pod 自己的 IP 與其他系統看到的 IP 相同
- Pod 內的容器共享網路命名空間

# Flannel CNI 工作原理與組件
## Flannel 簡介
Flannel 是一個為 Kubernetes 設計的簡單、輕量級的 CNI 插件，它通過創建覆蓋網路 (overlay network) 或使用現有的底層網路來實現 Pod 間的通信。

## Flannel 的主要組件
1. flanneld: 核心守護進程，在每個節點上運行
2. flannel-cni-plugin: 實現 CNI 規範的插件
3. kube-flannel-ds: Kubernetes DaemonSet，在每個節點上部署 Flannel


## Flannel 的工作流程
當 Flannel 部署到 Kubernetes 集群時：

1. 初始化: flanneld 從 Kubernetes API 或配置文件獲取網路配置
2. 子網分配: 為每個節點分配一個 CIDR 子網
3. 路由設置: 設置必要的路由規則，使 Pod 可以跨節點通信
4. 封包處理: 根據後端類型處理跨節點的 Pod 流量

Flannel 的配置存儲
Flannel 將網路配置存儲在 Kubernetes ConfigMap 中，包括：
- 網路 CIDR
- 後端類型
- 子網長度
- 其他特定後端的配置

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-flannel-cfg
  namespace: kube-system
data:
  net-conf.json: |
    {
      "Network": "10.244.0.0/16",
      "Backend": {
        "Type": "vxlan"
      }
    }
```

## Flannel 網路模型與封包流向
### Flannel 網路模型
Flannel 為每個節點分配一個子網，例如：

- 節點 1: 10.244.1.0/24
- 節點 2: 10.244.2.0/24
- 節點 3: 10.244.3.0/24

每個 Pod 從所在節點的子網中獲取 IP 地址。


### 封包流向分析
讓我們分析當 Pod A (10.244.1.2) 在節點 1 上，與 Pod B (10.244.2.3) 在節點 2 上通信時，封包的流向：

1. Pod A 發送封包:
    - Pod A 發送封包到 Pod B (10.244.2.3)
    - 封包離開 Pod A 的網路命名空間，進入節點 1 的網路命名空間

2. 節點 1 上的路由決策:
   - 節點 1 查找路由表，發現 10.244.2.0/24 子網在節點 2 上
   - 根據後端類型，封包被適當處理（封裝或路由）

3. 封包傳輸到節點 2:
   - 對於 VXLAN 後端，封包被封裝在 UDP 中
   - 封包通過物理網路傳輸到節點 2

4. 節點 2 上的處理:
   - 節點 2 接收封包，解封裝（如果需要）
   - 封包根據路由表轉發到 Pod B

5. Pod B 接收封包:
   - 封包進入 Pod B 的網路命名空間
   - Pod B 處理封包並回應

這個過程對應用是完全透明的，應用只看到普通的 IP 通信。

# 圖解 Flannel VXLAN 模式下的封包流向
```
Pod A (10.244.1.2) -> Pod B (10.244.2.3)

+------------------+                         +------------------+
| 節點 1            |                         | 節點 2            |
|                  |                         |                  |
| +-------------+  |                         |  +-------------+ |
| | Pod A       |  |                         |  | Pod B       | |
| | 10.244.1.2  |  |                         |  | 10.244.2.3  | |
| +-------------+  |                         |  +-------------+ |
|        |         |                         |         ^        |
|        v         |                         |         |        |
| +-------------+  |                         |  +-------------+ |
| | flannel0    |  |                         |  | flannel0    | |
| | (VXLAN介面)  |  |                         |  | (VXLAN介面)  | |
| +-------------+  |                         |  +-------------+ |
|        |         |                         |         ^        |
|        v         |                         |         |        |
| +-------------+  |    VXLAN封裝的UDP封包    |  +-------------+ |
| | eth0        | -|------------------------->| eth0        | |
| | 192.168.1.2 |  |                         |  | 192.168.1.3| |
| +-------------+  |                         |  +-------------+ |
+------------------+                         +------------------+
```

# Flannel 的各種後端模式比較
Flannel 支持多種後端模式，每種模式有不同的特點和適用場景：

## 1. VXLAN 後端 (默認)
工作原理：使用 Linux 內核的 VXLAN 功能創建覆蓋網路。

**優點：**
   - 良好的跨節點性能
   - 在大多數環境中工作良好
   - 完全在內核空間中運行

**缺點：**
  - 相比直接路由有輕微的性能開銷
  - 需要內核支持 VXLAN
適用場景：大多數生產環境

## 2. Host-GW (Host Gateway) 後端
工作原理：使用節點作為網關，通過直接路由實現 Pod 間通信。

**優點：**
   - 性能接近原生網路
   - 設置簡單，無需封裝/解封裝

**缺點：**
   - 要求所有節點在同一個二層網路
   - 不適用於跨子網或雲環境
適用場景：所有節點在同一數據中心的私有雲

## 3. UDP 後端
工作原理：使用 UDP 封裝在用戶空間中轉發封包。

**優點：**
- 兼容性最好，幾乎在所有環境中都能工作
- 不需要特殊的內核功能

**缺點：*
- 性能較差，用戶空間和內核空間之間有多次上下文切換
- CPU 使用率高

適用場景：測試環境或不支持其他後端的特殊環境

## 4. AWS VPC 後端
工作原理：利用 AWS VPC 路由表實現 Pod 網路。

**優點：**
  - 性能好，無需封裝
  - 與 AWS 基礎設施集成

**缺點：**
  - 僅適用於 AWS
  - 受 AWS 路由表條目限制

適用場景：AWS 上的 Kubernetes 集群

### 後端性能比較
根據測試數據，不同後端的性能排序（從高到低）：
  - Host-GW: 接近原生網路性能
  - AWS VPC: 在 AWS 上接近原生性能
  - VXLAN: 輕微性能損失（約 10-20%）
  - UDP: 顯著性能損失（約 50%以上）

選擇後端時應考慮環境限制和性能需求。

# 實作：使用 Flannel CNI 的 KinD 叢集
現在讓我們動手創建一個使用 Flannel CNI 的 KinD 叢集，並測試 Pod 間的通信。

## 步驟 1：創建 KinD 配置文件
KinD 默認使用的是 kindnet CNI，我們需要創建一個自定義配置來使用 Flannel。

創建一個名為 kind-flannel-config.yaml 的文件：

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: flannel-cluster
networking:
  disableDefaultCNI: true  # 禁用默認的 kindnet CNI
  podSubnet: "10.244.0.0/16"  # Flannel 默認使用的子網
nodes:
- role: control-plane
- role: worker
- role: worker
```

## 步驟 2：創建 KinD 叢集
```bash
> kind create cluster --config kind-flannel-config.yaml
```

這將創建一個包含 1 個控制平面節點和 2 個工作節點的叢集，但沒有 CNI。


## 步驟 3：部署 Flannel
下載 Flannel 清單文件：
```bash
curl -LO https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
```

應用 Flannel 清單：
```bash
kubectl apply -f kube-flannel.yml
```

等待 Flannel 部署完成：

```bash
kubectl -n kube-system wait --for=condition=ready pod -l app=flannel --timeout=300s
```

## 步驟 4：驗證網路功能
創建測試 Pod：

```bash
# 創建第一個 Pod
kubectl run pod1 --image=busybox --command -- sleep 3600

# 創建第二個 Pod
kubectl run pod2 --image=busybox --command -- sleep 3600
```

等待 Pod 運行：
```bash
kubectl wait --for=condition=ready pod pod1 pod2 --timeout=60s
```

獲取 Pod IP 地址：
```bash
POD1_IP=$(kubectl get pod pod1 -o jsonpath='{.status.podIP}')
POD2_IP=$(kubectl get pod pod2 -o jsonpath='{.status.podIP}')

echo "Pod 1 IP: $POD1_IP"
echo "Pod 2 IP: $POD2_IP"
```

測試 Pod 間通信：
```bash
# 從 Pod 1 ping Pod 2
kubectl exec pod1 -- ping -c 4 $POD2_IP

# 從 Pod 2 ping Pod 1
kubectl exec pod2 -- ping -c 4 $POD1_IP
```

如果 ping 成功，說明 Flannel 網路工作正常。

## 步驟 5：分析 Flannel 網路行為
查看 Flannel 配置：
```bash
kubectl -n kube-system get configmap kube-flannel-cfg -o yaml
```

查看節點上的 Flannel 介面：
```bash
# 獲取一個工作節點的名稱
WORKER_NODE=$(kubectl get nodes -o jsonpath='{.items[?(@.metadata.labels.kubernetes\.io/hostname!="flannel-cluster-control-plane")].metadata.name}' | awk '{print $1}')

# 在節點上檢查 Flannel 介面
docker exec $WORKER_NODE ip addr show flannel.1
docker exec $WORKER_NODE ip route | grep flannel
```

查看 VXLAN 封裝：
```bash
# 安裝 tcpdump
docker exec $WORKER_NODE apt-get update && docker exec $WORKER_NODE apt-get install -y tcpdump

# 捕獲 VXLAN 流量 (UDP 8472 端口)
docker exec $WORKER_NODE tcpdump -i any port 8472 -n
```

在另一個終端執行 Pod 間通信，觀察 VXLAN 封包。


## 步驟 6：探索 Flannel 的路由表
檢查節點上的路由表：
```bash
# 控制平面節點
docker exec flannel-cluster-control-plane route -n

# 工作節點
docker exec $WORKER_NODE route -n
```

觀察 Flannel 添加的路由，特別是指向其他節點子網的路由。

# Flannel 的優缺點與適用場景
## 優點
  - 簡單易用：配置簡單，易於部署和維護
  - 輕量級：資源消耗低
  - 穩定可靠：經過廣泛測試和使用
  - 多後端支持：適應不同環境需求

## 缺點
  - 功能有限：缺乏高級網路功能，如網路策略
  - 性能開銷：尤其是在 VXLAN 和 UDP 模式下
  - 故障診斷複雜：跨節點網路問題診斷困難
  - 缺乏細粒度控制：網路控制能力有限

## 適用場景
Flannel 最適合：
  - 小型到中型集群：資源需求適中
  - 開發和測試環境：簡單易用
  - 對網路功能要求不高的生產環境：基本連接性足夠
  - 初學者和教學環境：概念簡單，易於理解

對於需要高級網路功能（如網路策略、流量加密、深度監控）的環境，可能需要考慮 Calico、Cilium 等更強大的 CNI 插件。

# 總結
在本文中，我們深入探討了 Flannel CNI 的工作原理、組件和網路模型。我們了解了 Flannel 如何實現 Kubernetes 的網路模型，以及不同後端模式的特點和適用場景。

通過實際操作，我們創建了一個使用 Flannel 的 KinD 叢集，並測試了 Pod 間的通信。我們還分析了 Flannel 的網路行為，包括 VXLAN 封裝和路由表設置。

Flannel 作為一個簡單、輕量級的 CNI 插件，非常適合初學者和小型環境。但對於更複雜的網路需求，可能需要考慮其他 CNI 插件。

在下一篇文章中，我們將探討 Kubernetes 的儲存系統，包括 PersistentVolume 和 PersistentVolumeClaim。