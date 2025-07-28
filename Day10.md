---
title: "Day 10 ：Skaffold 入門與快速開發循環"
tags: 2025鐵人賽
date: 2025-07-20
---

#### Day 10：Skaffold 入門與快速開發循環
##### 重點
Skaffold 介紹與安裝
基本 skaffold.yaml 配置
開發模式（skaffold dev）與檔案監控
##### Lab
安裝 Skaffold
建立基本 skaffold.yaml
使用 skaffold dev 實現程式碼變更→自動部署的開發循環


# Skaffold 是什麼？
Skaffold 是由 Google 開發的開源工具，專門用於簡化 Kubernetes 應用程式的開發工作流程。它能夠自動化整個開發循環，包括：

監控程式碼變更：自動偵測原始碼的變更
建置映像檔：根據變更重新建置 Docker 映像檔
部署到 Kubernetes：將更新後的映像檔部署到 Kubernetes 叢集
狀態監控：監控部署狀態並提供即時反饋

Skaffold 的主要目標是讓開發人員能夠專注於程式碼開發，而不是花時間在重複性的部署工作上。

# 為什麼需要 Skaffold？
回顧昨天討論的傳統開發流程問題：

手動步驟多：修改程式碼→建置映像檔→載入映像檔→部署→測試
等待時間長：每個步驟都需要時間，影響開發效率
缺乏自動化：無法自動偵測變更並觸發重建
版本管理複雜：需要手動管理映像檔標籤

Skaffold 正是為了解決這些問題而設計的。它提供了一個完整的自動化開發循環，讓開發人員可以像在本地開發一樣，修改程式碼後立即看到變更效果。

# Skaffold 安裝
Skaffold 的安裝非常簡單，支援多種作業系統。

Linux
```bash
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
chmod +x skaffold
sudo mv skaffold /usr/local/bin
```

## 驗證安裝
安裝完成後，執行以下命令確認 Skaffold 已正確安裝：
```bash
> skaffold version
v2.9.0
```

# Skaffold 基本概念
在深入了解 Skaffold 的使用方法前，我們先來了解幾個核心概念：

## 1. Pipeline
Skaffold 的工作流程包含以下階段：

Build：建置應用程式的容器映像檔
Test：測試建置的映像檔
Deploy：將映像檔部署到 Kubernetes 叢集
Render：產生最終的 Kubernetes 資源配置
Port-forward：設置連接埠轉發，方便本地訪問應用程式
Log-tailing：顯示應用程式的日誌

## 2. 配置檔案 (skaffold.yaml)
Skaffold 使用 skaffold.yaml 檔案來定義整個工作流程。這個檔案指定了如何建置、測試和部署你的應用程式。

## 3. 運行模式
Skaffold 有兩種主要運行模式：

dev 模式：持續監控檔案變更，自動觸發重建和重新部署
run 模式：執行一次完整的工作流程，適合 CI/CD 環境


# 建立基本的 skaffold.yaml
現在，讓我們建立一個基本的 skaffold.yaml 檔案。假設我們有一個簡單的 Node.js 應用程式，目錄結構如下：
```
my-node-app/
├── app.js
├── Dockerfile
├── package.json
└── kubernetes/
    └── deployment.yaml
```

以下是一個基本的 skaffold.yaml 配置：
```yaml
apiVersion: skaffold/v4beta5
kind: Config
metadata:
  name: my-node-app
build:
  artifacts:
  - image: my-node-app
    docker:
      dockerfile: Dockerfile
deploy:
  kubectl:
    manifests:
    - kubernetes/deployment.yaml
```
這個配置檔案告訴 Skaffold：

建置：使用 Dockerfile 建置名為 my-node-app 的映像檔
部署：使用 kubectl 部署 kubernetes/deployment.yaml 定義的資源

## 配置檔案詳解
讓我們更詳細地了解 skaffold.yaml 的各個部分：

apiVersion
指定 Skaffold 配置檔案的版本。不同版本的 Skaffold 可能支援不同的配置選項。
```apiVersion: skaffold/v4beta5```


kind
指定配置檔案的類型，通常是 Config。
```kind: Config```

metadata
包含配置的元數據，如名稱。
```
metadata:
  name: my-node-app
```

build
定義如何建置應用程式的映像檔。

```
build:
  artifacts:
  - image: my-node-app  # 映像檔名稱
    docker:
      dockerfile: Dockerfile  # Dockerfile 路徑
```

Skaffold 支援多種建置器，包括：

Docker：使用 Docker 建置映像檔
Jib：無需 Dockerfile 建置 Java 應用程式
Buildpacks：使用 Cloud Native Buildpacks 建置映像檔
Bazel：使用 Bazel 建置系統
Custom：使用自定義指令建置映像檔


deploy
定義如何將應用程式部署到 Kubernetes。
```
deploy:
  kubectl:
    manifests:
    - kubernetes/deployment.yaml  # Kubernetes 資源定義檔案
```

Skaffold 支援多種部署器，包括：

kubectl：使用 kubectl 命令行工具
helm：使用 Helm 包管理器
kustomize：使用 Kustomize 客製化工具

# 使用 Skaffold 開發模式
Skaffold 的開發模式是它最強大的功能之一。在這個模式下，Skaffold 會：

監控原始碼變更
自動重新建置映像檔
自動重新部署到 Kubernetes
提供即時日誌輸出
設置連接埠轉發
要啟動開發模式，只需在專案目錄中執行：```skaffold dev```

## 開發模式工作流程
當你執行 skaffold dev 時，以下是發生的事情：

Skaffold 讀取 skaffold.yaml 配置檔案
建置所有定義的映像檔
部署應用程式到 Kubernetes 叢集
設置連接埠轉發
開始監控原始碼變更
當偵測到變更時，自動重複步驟 2-4
顯示應用程式的日誌輸出

這個循環會一直持續，直到你按下 Ctrl+C 停止 Skaffold。當停止時，Skaffold 會自動清理它部署的資源。

# 實作：使用 Skaffold 開發 Node.js 應用程式
現在，讓我們通過一個實際的例子來演示如何使用 Skaffold 開發 Node.js 應用程式。

步驟 1：準備應用程式
首先，我們建立一個簡單的 Node.js 應用程式的dockerfile：


```yaml
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 3000

CMD ["node", "app.js"]
```

步驟 2：建立 Kubernetes 部署檔案
建立 kubernetes 目錄並新增部署檔案： ```mkdir -p kubernetes```

建立 kubernetes/deployment.yaml：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: skaffold-demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: skaffold-demo
  template:
    metadata:
      labels:
        app: skaffold-demo
    spec:
      containers:
      - name: skaffold-demo
        image: skaffold-demo
        ports:
        - containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: skaffold-demo
spec:
  selector:
    app: skaffold-demo
  ports:
  - port: 80
    targetPort: 3000
  type: ClusterIP
```

步驟 3：建立 skaffold.yaml
在專案根目錄建立 skaffold.yaml：
```yaml
apiVersion: skaffold/v4beta5
kind: Config
metadata:
  name: skaffold-demo
build:
  artifacts:
  - image: skaffold-demo
    docker:
      dockerfile: Dockerfile
deploy:
  kubectl:
    manifests:
    - kubernetes/deployment.yaml
portForward:
- resourceType: service
  resourceName: skaffold-demo
  port: 80
  localPort: 8080
```

注意，我們添加了 portForward 部分，這會自動將服務的 80 連接埠映射到本地的 8080 連接埠。



步驟 4：啟動 Skaffold 開發模式
確保你的 KinD 叢集正在運行，然後執行：```skaffold dev```

Skaffold 會執行以下操作：

建置 Docker 映像檔
將映像檔載入到 KinD 叢集
部署應用程式
設置連接埠轉發
開始監控檔案變更


步驟 5：測試開發循環
現在，讓我們修改 app.js 檔案，看看 Skaffold 如何自動重新部署應用程式：

保存檔案後，你應該會看到 Skaffold 自動：

偵測到檔案變更
重新建置 Docker 映像檔
更新 Kubernetes 部署
保持連接埠轉發
打開瀏覽器訪問 http://localhost:8080，你應該會看到更新後的訊息。

# Skaffold 進階功能
Skaffold 還提供了許多進階功能，可以進一步優化開發工作流程：

1. 檔案同步
除了重建整個映像檔，Skaffold 還支援直接將變更的檔案同步到正在運行的容器中，大幅提高開發速度：
```yaml
build:
  artifacts:
  - image: skaffold-demo
    docker:
      dockerfile: Dockerfile
    sync:
      manual:
      - src: '*.js'
        dest: .
```
這個配置告訴 Skaffold，當 .js 檔案變更時，直接將它們複製到容器中，而不是重建整個映像檔。


2. 本地開發
Skaffold 支援將應用程式在本地運行，同時連接到 Kubernetes 叢集中的其他服務：
```yaml
build:
  artifacts:
  - image: skaffold-demo
    docker:
      dockerfile: Dockerfile
deploy:
  kubectl:
    manifests:
    - kubernetes/deployment.yaml
profiles:
- name: dev
  patches:
  - op: replace
    path: /deploy/kubectl/manifests
    value:
    - kubernetes/dev-deployment.yaml
```

使用 skaffold dev --profile=dev 可以啟用特定的開發配置。

3. 多服務開發
Skaffold 可以同時處理多個服務的開發：
```yaml
build:
  artifacts:
  - image: frontend
    context: ./frontend
    docker:
      dockerfile: Dockerfile
  - image: backend
    context: ./backend
    docker:
      dockerfile: Dockerfile
deploy:
  kubectl:
    manifests:
    - kubernetes/frontend.yaml
    - kubernetes/backend.yaml
```

4. 環境變數與模板
Skaffold 支援使用環境變數和模板來客製化配置：

```yaml
build:
  artifacts:
  - image: skaffold-demo
    docker:
      dockerfile: Dockerfile
deploy:
  kubectl:
    manifests:
    - kubernetes/deployment.yaml
  kustomize:
    paths: ["kustomize"]
```

# 總結
Skaffold 是一個強大的工具，可以大幅簡化 Kubernetes 應用程式的開發流程。通過自動化建置、部署和監控過程，它讓開發人員可以專注於程式碼開發，而不是繁瑣的部署工作。

今天，我們學習了：

Skaffold 的基本概念和安裝方法
如何建立和配置 skaffold.yaml
使用 skaffold dev 實現快速開發循環
如何通過實際例子使用 Skaffold 開發 Node.js 應用程式
Skaffold 的一些進階功能
在下一篇文章中，我們將深入探討 Skaffold 的更多進階功能，如多環境配置、調試支援和與 CI/CD 系統的整合。

# Lab 練習
練習 1：安裝 Skaffold 並設置開發環境
按照文章中的說明安裝 Skaffold
確認安裝成功：skaffold version
確保你有一個運行中的 KinD 叢集：kind get clusters
練習 2：使用 Skaffold 開發一個簡單的應用程式
按照文章中的步驟建立 Node.js 應用程式
建立 Kubernetes 部署檔案
建立 skaffold.yaml 配置檔案
啟動 Skaffold 開發模式：skaffold dev
修改應用程式程式碼，觀察 Skaffold 自動重新部署
使用瀏覽器或 curl 測試應用程式
練習 3：探索 Skaffold 的檔案同步功能
修改 skaffold.yaml，添加檔案同步配置
重新啟動 Skaffold 開發模式
修改 JavaScript 檔案，觀察 Skaffold 如何直接同步檔案而不重建映像檔
比較使用檔案同步和不使用檔案同步時的部署速度
通過這些練習，你將能夠親身體驗 Skaffold 如何改善 Kubernetes 應用程式的開發體驗，並為後續學習更進階的功能打下基礎。