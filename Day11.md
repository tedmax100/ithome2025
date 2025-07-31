---
title: "Day 11：KinD 進階配置管理與 DevSpace 多應用開發與環境變數"
tags: 2025鐵人賽
date: 2025-07-20
---

#### Day 11：KinD 進階配置管理與 DevSpace 多應用開發與環境變數

##### 重點
KinD 進階配置管理
多節點叢集設計模式與最佳實踐
為不同應用場景定制 KinD 配置
配置版本控制與團隊共享策略
使用腳本自動化 KinD 環境管理
DevSpace 與 KinD 配置聯動
在 DevSpace 工作流中創建/銷毀 KinD 叢集
使用 DevSpace hooks 管理 KinD 生命週期
配置 DevSpace 使用特定 KinD 叢集
DevSpace 多應用（前後端）設定
結合 ConfigMap 與環境變數
使用 DevSpace profiles 區分環境

##### Lab

建立包含前後端的 devspace.yaml
使用 ConfigMap 注入環境變數
設定 dev/staging profiles
實作練習
  設計三種不同用途的 KinD 配置（開發、測試、模擬生產）
  創建自動化腳本管理 KinD 環境生命週期
  配置 DevSpace 與不同 KinD 環境整合
  使用 ConfigMap 為不同環境注入配置
  設計 DevSpace profiles 對應不同環境需求
  實現完整的開發-測試工作流程


在前一篇文章中，我們介紹了 DevSpace 的基本使用方法，了解了如何使用它來簡化 Kubernetes 應用程式的開發流程。今天，我們將深入探討 KinD 的進階配置管理，以及如何結合 DevSpace 進行多應用開發和環境變數管理。

# KinD 進階配置管理
KinD 配置檔案基礎
KinD (Kubernetes in Docker) 使用 YAML 配置檔案來定義叢集的結構和行為。基本的配置檔案如下：

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: my-cluster
```
這個最小配置會創建一個名為 my-cluster 的單節點叢集。但 KinD 的真正強大之處在於它的進階配置選項。

## 多節點叢集設計
在生產環境中，Kubernetes 叢集通常包含多個節點，分為控制平面節點和工作節點。使用 KinD，我們可以模擬這種多節點架構：
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: multi-node-cluster
nodes:
- role: control-plane
- role: worker
- role: worker
```

這個配置創建了一個有 1 個 Control plane 節點和 2 個 worker 節點的叢集。

多控制平面高可用配置
對於高可用性設計，我們可以配置多個控制平面節點：

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: ha-cluster
nodes:
- role: control-plane
- role: control-plane
- role: control-plane
- role: worker
- role: worker
- role: worker
```

這個配置創建了一個有 3 個控制平面節點和 3 個工作節點的高可用叢集。

節點特定配置
KinD 允許為每個節點配置特定的設置，例如暴露端口、掛載目錄等：
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
  extraMounts:
  - hostPath: /path/on/host
    containerPath: /path/in/node
- role: worker
```

這個配置將主機的 80 和 443 端口映射到 control plane 節點，並掛載了一個目錄。

網路配置
KinD 支持自定義 Pod 子網和服務子網：
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  podSubnet: "10.244.0.0/16"
  serviceSubnet: "10.96.0.0/12"
```
這對於模擬特定網路環境或避免與現有網路衝突非常有用。

## 為不同應用場景定制 KinD 配置
根據不同的開發需求，我們可以創建專門的 KinD 配置：

開發環境配置
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: dev-cluster
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80
    hostPort: 8080
  - containerPort: 443
    hostPort: 8443
```

測試環境配置
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: prod-sim-cluster
nodes:
- role: control-plane
- role: control-plane
- role: control-plane
- role: worker
- role: worker
- role: worker
```

## 配置版本控制與團隊共享策略
在團隊環境中，保持 KinD 配置的一致性至關重要。以下是一些最佳實踐：

將配置文件納入版本控制
使用目錄結構組織配置：
```
kind-configs/
├── dev/
│   └── kind-dev.yaml
├── test/
│   └── kind-test.yaml
└── prod-sim/
    └── kind-prod-sim.yaml
```

文檔化配置目的和使用方法：在每個配置文件中添加註釋，說明其用途和特殊設置。
使用環境變數參數化配置：
```
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: ${CLUSTER_NAME:-dev-cluster}
```

然後使用腳本處理變數替換：```envsubst < kind-template.yaml > kind-config.yaml```

## 使用腳本自動化 KinD 環境管理
為了簡化 KinD 叢集的創建和管理，我們可以創建自動化腳本：

創建叢集腳本 (create-cluster.sh)
```shell
#!/bin/bash
set -e

CLUSTER_TYPE=${1:-dev}
CLUSTER_NAME=${2:-kind-cluster}

echo "Creating $CLUSTER_TYPE cluster named $CLUSTER_NAME"

# 選擇配置文件
CONFIG_FILE="./kind-configs/$CLUSTER_TYPE/kind-$CLUSTER_TYPE.yaml"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Configuration file $CONFIG_FILE not found!"
    exit 1
fi

# 環境變數替換
export CLUSTER_NAME=$CLUSTER_NAME
envsubst < $CONFIG_FILE > /tmp/kind-config.yaml

# 創建叢集
kind create cluster --config /tmp/kind-config.yaml

echo "Cluster $CLUSTER_NAME created successfully!"
```

刪除叢集腳本 (delete-cluster.sh)

```shell
#!/bin/bash
set -e

CLUSTER_NAME=${1:-kind-cluster}

echo "Deleting cluster $CLUSTER_NAME"
kind delete cluster --name $CLUSTER_NAME

echo "Cluster $CLUSTER_NAME deleted successfully!"
```

叢集狀態檢查腳本 (check-cluster.sh)

```shell
#!/bin/bash
set -e

CLUSTER_NAME=${1:-kind-cluster}

echo "Checking status of cluster $CLUSTER_NAME"

if kind get clusters | grep -q $CLUSTER_NAME; then
    echo "Cluster $CLUSTER_NAME exists"
    kubectl --context kind-$CLUSTER_NAME get nodes
    kubectl --context kind-$CLUSTER_NAME get pods -A
else
    echo "Cluster $CLUSTER_NAME does not exist"
    exit 1
fi
```

# DevSpace 與 KinD 配置聯動
DevSpace 可以與 KinD 緊密集成，實現完整的開發循環管理。

在 DevSpace 工作流中創建/銷毀 KinD 叢集
DevSpace 提供了 hooks 功能，可以在工作流的不同階段執行命令。我們可以利用這個功能來管理 KinD 叢集的生命週期。

```yaml
version: v2beta1
name: my-app

images:
  my-app:
    image: my-app
    dockerfile: Dockerfile

deployments:
  my-app:
    kubectl:
      manifests:
      - kubernetes/*.yaml

hooks:
  - command: ["./scripts/create-cluster.sh", "dev", "devspace-cluster"]
    when: before:deploy
    where: local
  - command: ["./scripts/delete-cluster.sh", "devspace-cluster"]
    when: after:deploy:exit
    where: local
```

這個配置會在 DevSpace 部署前創建一個 KinD 叢集，並在 DevSpace 結束後刪除它。

## 使用 DevSpace hooks 管理 KinD 生命週期
DevSpace 的 hooks 可以更精細地控制 KinD 叢集的生命週期：
```yaml
version: v2beta1
name: my-app

images:
  my-app:
    image: my-app
    dockerfile: Dockerfile

deployments:
  my-app:
    kubectl:
      manifests:
      - kubernetes/*.yaml

profiles:
  - name: dev
    hooks:
      - command: ["./scripts/create-cluster.sh", "dev", "devspace-dev"]
        when: before:deploy
        where: local
      - command: ["./scripts/delete-cluster.sh", "devspace-dev"]
        when: after:deploy:exit
        where: local
  - name: test
    hooks:
      - command: ["./scripts/create-cluster.sh", "test", "devspace-test"]
        when: before:deploy
        where: local
      - command: ["./scripts/delete-cluster.sh", "devspace-test"]
        when: after:deploy:exit
        where: local
```

使用 devspace dev --profile=dev 或 devspace dev --profile=test 可以啟動不同的環境。

## 配置 DevSpace 使用特定 KinD 叢集
如果你已經有一個運行中的 KinD 叢集，可以配置 DevSpace 使用它：

```yaml
version: v2beta1
name: my-app

images:
  my-app:
    image: my-app
    dockerfile: Dockerfile

deployments:
  my-app:
    kubectl:
      manifests:
      - kubernetes/*.yaml
      kubeContext: kind-my-existing-cluster
```

kubeContext 參數指定了 DevSpace 應該使用的 Kubernetes 上下文。對於 KinD 叢集，上下文名稱通常是 kind-<cluster-name>。

## DevSpace 多應用（前後端）設定
現代應用程式通常包含前端和後端組件。DevSpace 可以輕鬆管理這種多應用架構。

基本多應用配置
```yaml
version: v2beta1
name: fullstack-app

images:
  frontend:
    image: frontend
    context: ./frontend
    dockerfile: Dockerfile
  backend:
    image: backend
    context: ./backend
    dockerfile: Dockerfile

deployments:
  fullstack:
    kubectl:
      manifests:
      - kubernetes/frontend.yaml
      - kubernetes/backend.yaml
```

這個配置定義了前端和後端兩個組件，每個組件都有自己的 Dockerfile 和部署清單。

文件同步優化
對於前端應用，我們通常希望在開發過程中快速同步文件變更，而不是重建整個映像檔：
```yaml
version: v2beta1
name: fullstack-app

images:
  frontend:
    image: frontend
    context: ./frontend
    dockerfile: Dockerfile
  backend:
    image: backend
    context: ./backend
    dockerfile: Dockerfile

deployments:
  fullstack:
    kubectl:
      manifests:
      - kubernetes/frontend.yaml
      - kubernetes/backend.yaml

dev:
  sync:
    - imageSelector: frontend
      localSubPath: ./frontend/src
      containerPath: /app/src
      excludePaths:
        - node_modules/
    - imageSelector: backend
      localSubPath: ./backend/src
      containerPath: /app/src
      excludePaths:
        - node_modules/
```

結合 ConfigMap 與環境變數
在 Kubernetes 中，ConfigMap 是管理應用配置的標準方式。DevSpace 可以幫助我們管理這些配置：

創建 ConfigMap 模板 (kubernetes/configmap.yaml)
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  API_URL: "http://backend:8080"
  DEBUG: "true"
  ENVIRONMENT: "development"
```

在部署中使用 ConfigMap

前端部署 (kubernetes/frontend.yaml):
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend
        image: frontend
        ports:
        - containerPort: 3000
        envFrom:
        - configMapRef:
            name: app-config
```

後端部署 (kubernetes/backend.yaml):
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: backend
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: app-config
```

使用 DevSpace profiles 區分環境
DevSpace 的 profiles 功能允許我們為不同環境定義不同的配置：

```yaml
version: v2beta1
name: fullstack-app

images:
  frontend:
    image: frontend
    context: ./frontend
    dockerfile: Dockerfile
  backend:
    image: backend
    context: ./backend
    dockerfile: Dockerfile

deployments:
  fullstack:
    kubectl:
      manifests:
      - kubernetes/base/*.yaml

profiles:
  - name: dev
    patches:
      - op: add
        path: deployments.fullstack.kubectl.manifests
        value:
          - kubernetes/base/*.yaml
          - kubernetes/overlays/dev/*.yaml
  - name: staging
    patches:
      - op: add
        path: deployments.fullstack.kubectl.manifests
        value:
          - kubernetes/base/*.yaml
          - kubernetes/overlays/staging/*.yaml
  - name: production
    patches:
      - op: replace
        path: images.frontend.dockerfile
        value: Dockerfile.prod
      - op: replace
        path: images.backend.dockerfile
        value: Dockerfile.prod
      - op: add
        path: deployments.fullstack.kubectl.manifests
        value:
          - kubernetes/base/*.yaml
          - kubernetes/overlays/production/*.yaml
```

這個配置定義了三個環境：dev、staging 和 production。每個環境使用不同的配置文件，production 環境還使用了不同的 Dockerfile。


# 實際案例：全棧應用開發環境
現在，讓我們通過一個實際案例來演示如何結合 KinD 和 DevSpace 創建一個完整的全棧應用開發環境。

項目結構
```
fullstack-app/
├── frontend/
│   ├── src/
│   ├── Dockerfile
│   └── Dockerfile.prod
├── backend/
│   ├── src/
│   ├── Dockerfile
│   └── Dockerfile.prod
├── kubernetes/
│   ├── base/
│   │   ├── frontend.yaml
│   │   ├── backend.yaml
│   │   └── services.yaml
│   └── overlays/
│       ├── dev/
│       │   └── configmap.yaml
│       ├── staging/
│       │   └── configmap.yaml
│       └── production/
│           └── configmap.yaml
├── kind-configs/
│   ├── dev/
│   │   └── kind-dev.yaml
│   ├── staging/
│   │   └── kind-staging.yaml
│   └── production/
│       └── kind-production.yaml
├── scripts/
│   ├── create-cluster.sh
│   └── delete-cluster.sh
└── devspace.yaml
```

KinD 開發環境配置 (kind-configs/dev/kind-dev.yaml)
```
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: ${CLUSTER_NAME:-dev-cluster}
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80
    hostPort: 8080
  - containerPort: 3000
    hostPort: 3000
  - containerPort: 8000
    hostPort: 8000
```

DevSpace 配置 (devspace.yaml)
```yaml
version: v2beta1
name: fullstack-app

images:
  frontend:
    image: frontend
    context: ./frontend
    dockerfile: Dockerfile
  backend:
    image: backend
    context: ./backend
    dockerfile: Dockerfile

deployments:
  fullstack:
    kubectl:
      manifests:
      - kubernetes/base/*.yaml
      - kubernetes/overlays/dev/*.yaml

dev:
  sync:
    - imageSelector: frontend
      localSubPath: ./frontend/src
      containerPath: /app/src
      excludePaths:
        - node_modules/
    - imageSelector: backend
      localSubPath: ./backend/src
      containerPath: /app/src
      excludePaths:
        - node_modules/
  ports:
    - imageSelector: frontend
      forward:
        - port: 3000
          localPort: 3000
    - imageSelector: backend
      forward:
        - port: 8000
          localPort: 8000
  open:
    - url: http://localhost:3000

hooks:
  - command: ["./scripts/create-cluster.sh", "dev", "fullstack-dev"]
    when: before:deploy
    where: local
  - command: ["./scripts/delete-cluster.sh", "fullstack-dev"]
    when: after:deploy:exit
    where: local

profiles:
  - name: staging
    patches:
      - op: replace
        path: deployments.fullstack.kubectl.manifests
        value:
          - kubernetes/base/*.yaml
          - kubernetes/overlays/staging/*.yaml
    hooks:
      - command: ["./scripts/create-cluster.sh", "staging", "fullstack-staging"]
        when: before:deploy
        where: local
      - command: ["./scripts/delete-cluster.sh", "fullstack-staging"]
        when: after:deploy:exit
        where: local
  - name: production
    patches:
      - op: replace
        path: images.frontend.dockerfile
        value: Dockerfile.prod
      - op: replace
        path: images.backend.dockerfile
        value: Dockerfile.prod
      - op: replace
        path: deployments.fullstack.kubectl.manifests
        value:
          - kubernetes/base/*.yaml
          - kubernetes/overlays/production/*.yaml
    hooks:
      - command: ["./scripts/create-cluster.sh", "production", "fullstack-production"]
        when: before:deploy
        where: local
      - command: ["./scripts/delete-cluster.sh", "fullstack-production"]
        when: after:deploy:exit
        where: local
```


基本 Kubernetes 配置
前端服務 (kubernetes/base/frontend.yaml):
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend
        image: frontend
        ports:
        - containerPort: 3000
        envFrom:
        - configMapRef:
            name: app-config
---
apiVersion: v1
kind: Service
metadata:
  name: frontend
spec:
  selector:
    app: frontend
  ports:
  - port: 3000
    targetPort: 3000
  type: ClusterIP
```

後端服務 (kubernetes/base/backend.yaml):

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: backend
        ports:
        - containerPort: 8000
        envFrom:
        - configMapRef:
            name: app-config
---
apiVersion: v1
kind: Service
metadata:
  name: backend
spec:
  selector:
    app: backend
  ports:
  - port: 8000
    targetPort: 8000
  type: ClusterIP
```

開發環境配置 (kubernetes/overlays/dev/configmap.yaml):
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  API_URL: "http://backend:8000"
  FRONTEND_URL: "http://frontend:3000"
  DEBUG: "true"
  ENVIRONMENT: "development"
```

啟動開發環境
使用以下命令啟動開發環境：
```devspace dev```

DevSpace 會：

創建 KinD 叢集
建置前端和後端映像檔
部署應用到叢集
設置端口轉發
監控文件變更並自動同步或重建
在瀏覽器中打開應用

切換到其他環境
要切換到 staging 環境：`devspace dev --profile=staging`
要切換到 production 環境：`devspace dev --profile=production`

## 最佳實踐與總結
KinD 最佳實踐
根據需求調整資源：為 KinD 節點分配適當的 CPU 和內存
使用持久卷：對於需要持久存儲的應用，配置持久卷
網路規劃：確保 KinD 叢集的網路不與現有網路衝突
版本控制配置：將 KinD 配置文件納入版本控制
自動化腳本：使用腳本簡化叢集管理

DevSpace 最佳實踐
使用文件同步：對於前端開發，優先使用文件同步而不是重建映像檔
合理組織 profiles：使用 profiles 分離不同環境的配置
利用 hooks：使用 hooks 自動化環境設置和清理
監控資源使用：定期檢查資源使用情況，避免開發環境資源不足
整合 CI/CD：將 DevSpace 配置與 CI/CD 流程整合 

# 總結
在本文中，我們深入探討了 KinD 的進階配置管理和 DevSpace 的多應用開發能力。通過結合這兩個強大的工具，我們可以創建一個高效、靈活的 Kubernetes 開發環境，大大提高開發效率。

KinD 提供了一個輕量級但功能強大的 Kubernetes 環境，可以根據不同的需求進行定制。DevSpace 則簡化了應用的構建、部署和監控過程，特別是對於多組件應用。

通過使用環境變數、ConfigMap 和 profiles，我們可以輕鬆管理不同環境的配置，確保開發、測試和生產環境的一致性。自動化腳本和 hooks 進一步簡化了環境管理，使開發人員可以專注於代碼而不是基礎設施。

在下一篇文章中，我們將探討如何將這些工具與 CI/CD 流程集成，實現從開發到部署的完整自動化。


# Lab 練習
練習 1：創建多節點 KinD 叢集
創建一個包含 1 個控制平面節點和 2 個工作節點的 KinD 配置文件
使用該配置文件創建叢集
驗證節點狀態
練習 2：使用 DevSpace 開發全棧應用
創建一個簡單的前端應用（使用 React 或 Vue）
創建一個簡單的後端 API（使用 Node.js 或其他語言）
為兩個應用創建 Dockerfile
創建 Kubernetes 部署文件
配置 DevSpace 並啟動開發環境
修改前端和後端代碼，觀察變更如何自動同步
練習 3：使用 profiles 管理不同環境
為開發、測試和生產環境創建不同的 ConfigMap
在 DevSpace 配置中添加對應的 profiles
測試在不同 profile 下啟動應用
觀察不同環境中的配置差異
通過這些練習，你將能夠掌握 KinD 和 DevSpace 的進階功能，為團隊建立高效的 Kubernetes 開發環境。