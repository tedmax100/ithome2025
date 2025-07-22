---
title: "Day 1 KinD 安裝與啟動"
tags: 2025鐵人賽
date: 2025-07-20
---

#### KinD 安裝與啟動
##### 重點
- KinD、kubectl 安裝與設定
##### Lab
- 安裝 KinD 並啟動本地叢集
- 使用 kubectl get nodes 驗證
- 執行 kubectl get pods、kubectl get svc 等指令，初步體驗



Kubernetes（K8s）是一個分散式容器編排平台，本身沒有「本地單機」的官方安裝方式。

Minikube、k3d、KinD 都是「本地端模擬 K8s」的工具，讓你不用真的建一整套生產環境，也能體驗/測試 K8s 的部署流程。

你在這些環境寫的 YAML、用的 kubectl 指令，幾乎都跟生產環境一樣。

- 本地快速搭建測試環境：不用每次 push 到遠端，自己本機就能部署、測試、debug。
- CI/CD 腳本驗證：可以模擬 GitHub Actions、GitLab CI 的 deploy 步驟，確保腳本沒問題再上線。
- 熟悉 K8s 資源操作：練習 kubectl、YAML、Helm，與生產環境無縫接軌。
- 排查容器化問題：本地就能重現容器啟動、網路、設定等問題，減少「雲端才爆炸」的風險。
- 跨團隊溝通：理解 K8s 架構與資源，有助於與 DevOps/SRE 協作。


# 為什麼選 KinD？Minikube、k3d、KinD 差異比較
1. 這些工具是什麼？
Minikube、k3d、kind 都是用來在本地端快速建立 Kubernetes 叢集的工具。
它們本質上都在幫你「模擬一個小型的 K8s 環境」，方便開發、測試、學習，不需像生產環境那樣複雜。Minikube、k3d、KinD 都是用來在本地端快速建立 Kubernetes 叢集的工具。

它們本質上都在幫你「模擬一個小型的 K8s 環境」，方便開發、測試、學習，不需像生產環境那樣複雜。

各工具簡介與差異
| 工具 | 底層運作 | 特色與優勢 | 適合情境 |
| --- | ------- | -------- | ------- |
| Minikube | VM/Container | 最接近原生 K8s，安裝簡單，支援多種驅動（Docker、VM）、有 Dashboard、外掛多 | 學習、開發、初學者、單機測試 |
| k3d |	Docker Container | 基於輕量級 k3s（非官方 K8s），超快啟動，佔用資源少，適合多節點模擬 | 需要多節點、資源有限、CI 測試 |
| KinD | Docker Container | 使用官方 K8s，全部跑在 Docker 裡，超快啟動，非常適合 CI/CD pipeline | 自動化測試、CI/CD、開發環境 |

# 為什麼 KinD 適合後端開發工程師？
- 官方原生 K8s：KinD 使用的是官方 Kubernetes 發行版，行為最接近生產環境。
- 超快啟動、易於重建：所有節點都以 Docker container 執行，建立、銷毀環境非常快速，適合開發與 CI/CD。
- 多節點支援：輕鬆模擬多節點叢集，方便測試服務發現、Pod 調度等情境。
- 腳本化/自動化友善：KinD 支援 YAML 定義叢集拓撲，適合納入版本控制、團隊協作與自動化流程。
- 本地映像直接使用：可直接將本地 Docker image 載入 KinD 叢集，測試流程順暢。
- 資源消耗低：相比 VM 型方案（如 Minikube），KinD 僅需 Docker，對硬體資源需求較低。


## kind 指令有哪些功能？常用子指令總覽

| 指令 | 功能說明 | 描述 |
| --- | ------- | --- |
| kind | create cluster | 建立一個新的 Kubernetes 叢集（可加參數自訂名稱等） |
| kind |  delete cluster | 刪除一個叢集（預設為 kind，也可指定名稱） |
| kind |  get clusters | 列出目前所有由 kind 建立的叢集 |
| kind |  get kubeconfig | 輸出 kubeconfig 設定（可用於多叢集/多帳號場景） |
| kind |  load docker-image | 將本地 Docker image 載入 kind 叢集 |
| kind |  export logs | 匯出叢集的詳細 log（除錯用） |
| kind |  build node-image | 從 Kubernetes 原始碼自製 node image（進階用） |
| kind |  version | 顯示 kind 版本 |
| kind |  completion | 產生 shell 自動補全腳本 |
| kind |  export kubeconfig | 匯出 kubeconfig 檔案 |

#### 常用情境舉例
- 重建測試環境：kind delete cluster && kind create cluster
- 測試本地鏡像：docker build -t myapp:latest . → kind load docker-image myapp:latest
- 多叢集管理：kind create cluster --name demo、kind get clusters
- 除錯：kind export logs --name demo

## Kubectl 與 kind 指令
kind 和 kubectl 都是 Kubernetes 相關工具，但他們的用途和功能不同，彼此之間有明確的分工和關聯：

**kind** 是一個用來建立、管理本地 Kubernetes 叢集的工具。
它主要負責：
- 建立/刪除/管理本地的 K8s 叢集（都跑在 Docker container 裡）
- 幫你準備好一個可以練習或測試的 K8s 環境

**kubectl** 是Kubernetes 的官方命令列管理工具。
它主要負責：
- 跟現有的 Kubernetes 叢集溝通
- 部署、查詢、管理 K8s 上的各種資源（Pod、Service、Deployment 等）


```
kind ──> 建立本地 K8s 叢集 ──>
                              \
                                → kubectl ──> 管理/查詢/部署 K8s 資源
```

# Lab
## 安裝 KinD

> 必須先安裝好 Docker。沒有 Docker，KinD 無法運作。

參考 [KinD 官方網站 Installation](https://kind.sigs.k8s.io/docs/user/quick-start#installation) 的教學。

小弟我是 Linux Mint 的系統，首先我需要確認自己的 CPU 架構是 x86 還是 arm 指令集。
```bash=
> uname -m
x86_64
```

根據指令下載對應版本的 kind 並給予可執行權限，鐵人賽當下 KinD 最新是 `v0.29.0` 版本，之後應該當下最新的版本來換版本號進行安裝。
```bash=
# For AMD64 / x86_64
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.29.0/kind-linux-amd64

chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

輸入指令確定 kind 安裝完成，且版本如預期。
```bash=
> kind version
kind v0.29.0 go1.24.2 linux/amd64
```

### 安裝 kubectl
參考[k8s Install kubectl on Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
因為剛剛已經確認過是x86指令集了。
```bash=
# 下載 kubectl 最新版本
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

# 安裝到系統路徑
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# 確認安裝完成
kubectl version
```

## 建立本地 Kubernetes Cluster
我們能透過 `get  Gets one of [clusters, nodes, kubeconfig]` 這指令，來取得 cluster、nodes、kubeconfig 的現有清單，確定我們目前是沒有建立 cluster 的。
```bash=
> kind get clusters
No kind clusters found.
```

然後我們就能建立 k8s cluster 了。透過 `create` 指令，以及 `--name` 就能定義 cluster 名稱，若是沒給名稱，則預設是`kind`。所以可以同時建立多個不同名稱的 kind 叢集，彼此互不影響，這就很像 Docker-compose。

```bash=
> kind create cluster --name demo

Creating cluster "demo" ...
 # 確認並下載（如果本地沒有）Kubernetes 節點所需的 Docker image。
 ✓ Ensuring node image (kindest/node:v1.33.1) 🖼 
 # 準備要啟動的節點（這裡通常就是一個 control-plane node，也可以設定多節點）。
 ✓ Preparing nodes 📦  
 # 產生並寫入 Kubernetes cluster 的設定檔。
 ✓ Writing configuration 📜 
 # 啟動 Kubernetes 的 control-plane node（即主要管理節點）。
 ✓ Starting control-plane 🕹️ 
 # 安裝 CNI (Container Network Interface)，讓 Pod 之間可以互相通訊。
 ✓ Installing CNI 🔌 
 # 安裝預設的 StorageClass，讓你可以在叢集裡建立 PersistentVolumeClaim 等儲存資源。
 ✓ Installing StorageClass 💾 
# kind 已經自動幫你把 kubectl 的預設 context 切換到這個新建立的 demo 叢集。
# 之後你用 kubectl 指令時，會直接操作這個叢集。
Set kubectl context to "kind-demo"
You can now use your cluster with:

kubectl cluster-info --context kind-demo

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community
```



### 驗證 Cluster 狀態
```bash=
kubectl cluster-info
# 這是你的 Kubernetes 叢集的 API Server 位址（本機連接埠 35123）
Kubernetes control plane is running at https://127.0.0.1:35123
# 這是叢集內 DNS 服務的位置。
CoreDNS is running at https://127.0.0.1:35123/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

# 如果你要進一步除錯，可以用這個指令取得更詳細的診斷資訊。
To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

當我們透過瀏覽器打開`https://127.0.0.1:35123`，會看到以下錯誤 403 Forbidden 訊息。
```jsonld=
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "forbidden: User \"system:anonymous\" cannot get path \"/\"",
  "reason": "Forbidden",
  "details": {},
  "code": 403
}
```

因為 K8s API Server 需要權限驗證，你沒有帶上 kubeconfig 憑證資訊，所以被視為「匿名使用者（system:anonymous）」。匿名使用者沒有權限存取 API Server 的根路徑，因此回傳 403 Forbidden。通常我們會透過`kubectl` 來存取這的資訊，因為 kubectl 會自動讀取 kubeconfig 檔案，帶上認證資訊（token、憑證等）跟 API Server 溝通。
除非我們有 token。

#### 取得憑證
我們能透過 `kubectl config view` 查詢。
找到目前 context 對應的 user 設定。
- 如果是 token authentication，會有 token: 欄位。
- 如果是 client-certificate authentication，會有 certificate-authority-data、client-certificate-data: 和 client-key-data:。這 3 個資訊用 base64 encode 分別存成 `ca.crt`、`cert.crt` 和 `key.key`，3個檔案後。就能透過 curl 指令來看件資訊了。

```bash=
> kubectl config view --raw

apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZ
    server: https://127.0.0.1:35123
  name: kind-demo
contexts:
- context:
    cluster: kind-demo
    user: kind-demo
  name: kind-demo
current-context: kind-demo
kind: Config
preferences: {}
users:
- name: kind-demo
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk
    client-key-data: LS0tLS1CRU==
```

這下我們就能成功的獲取該 K8s API server 跟目錄下所有可操作的 API 路徑了，這裡只列舉一部分。
```bash=
> curl --cert cert.crt --key key.key --cacert ca.crt https://127.0.0.1:35123

{
  "paths": [
    "/.well-known/openid-configuration",
    "/api",
    "/api/v1",
    "/apis",
    "/healthz",
    "/livez",
    "/metrics",
    "/openapi/v2",
    "/openapi/v3",
    "/openapi/v3/",
    "/readyz",
    "/version"
  ]
}%                                
```

所以若要用 API 來存取 k8s 除了用 kubectl，其實也能依靠`client certificate authentication`來存取，只要憑證沒過期沒被撤銷就能繼續用。


### 刪除 Cluster (清理資源)
當我們做完實驗或測試時，就能把資源給刪除乾淨了，畢竟 KinD 就是來給我們測試跟實驗用的，不適合跑在環境中。
我們能先查看現有 KinD cluster 有哪些，這時應該只有剛剛建立出來的 demo cluster。
```bash
> kind get clusters
demo
```

刪除跟建立Cluster的指令非常相似。
```bash=
# 刪除單一 Cluster
kind delete cluster --name demo

# 刪除所有 Cluster
for c in $(kind get clusters); do kind delete cluster --name $c; done

# 清理 Docker 殘留資源
docker system prune -a

# 刪除 kubeconfig 設定
# KinD 產生的 kubeconfig 檔案預設寫在 ~/.kube/config，如果你只用 KinD，可以直接清空或移除
rm -rf ~/.kube/config
```

# 總結
KinD 是開發者「本地端玩 Kubernetes」的超好用利器，能讓你無痛練習、測試、Debug，跟生產環境無縫接軌，是每個現代後端工程師都該會的技能！

如果你還沒玩過 KinD，真的推薦馬上動手試試，你會發現 K8s 沒那麼可怕，而且本地端也能很快上手、驗證想法！

https://www.blueshoe.io/blog/minikube-vs-k3d-vs-kind-vs-getdeck-beiboot/#performance-evaluation