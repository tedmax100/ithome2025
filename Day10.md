---
title: "Day 10 10：DevSpace 入門與快速開發循環"
tags: 2025鐵人賽
date: 2025-07-20
---
# Day 10：DevSpace 入門與快速開發循環
重點
DevSpace 介紹與安裝
基本 devspace.yaml 配置
開發模式（devspace dev）與檔案監控
DevSpace UI 儀表板介紹
DevSpace 與 VS Code 整合
Lab
安裝 DevSpace
建立基本 devspace.yaml
使用 devspace dev 實現程式碼變更→自動部署的開發循環
體驗 DevSpace UI 監控應用狀態
設定 VS Code 與 DevSpace 整合

# DevSpace 是什麼？
DevSpace 是一個開源的雲原生開發工具，專門用於簡化 Kubernetes 應用程式的開發工作流程。它能夠自動化整個開發循環，包括：

監控程式碼變更：自動偵測原始碼的變更
建置映像檔：根據變更重新建置 Docker 映像檔
部署到 Kubernetes：將更新後的映像檔部署到 Kubernetes 叢集
狀態監控：透過內建 UI 儀表板監控部署狀態並提供即時反饋
檔案同步：直接將變更同步到容器，無需重建映像檔
DevSpace 的主要目標是讓開發人員能夠專注於程式碼開發，而不是花時間在重複性的部署工作上。

## 為什麼需要 DevSpace？
回顧昨天討論的傳統開發流程問題：

手動步驟多：修改程式碼→建置映像檔→載入映像檔→部署→測試
等待時間長：每個步驟都需要時間，影響開發效率
缺乏自動化：無法自動偵測變更並觸發重建
版本管理複雜：需要手動管理映像檔標籤
DevSpace 正是為了解決這些問題而設計的。它提供了一個完整的自動化開發循環，讓開發人員可以像在本地開發一樣，修改程式碼後立即看到變更效果。相較於其他類似工具，DevSpace 還提供了內建的 UI 儀表板，讓開發者可以更直觀地監控應用狀態。

## DevSpace 安裝
DevSpace 的安裝非常簡單，支援多種作業系統。

Linux
```bash
curl -L -o devspace "https://github.com/loft-sh/devspace/releases/latest/download/devspace-linux-amd64"
sudo install -c -m 0755 devspace /usr/local/bin
```

驗證安裝
安裝完成後，執行以下命令確認 DevSpace 已正確安裝：
```bash
> devspace version
v6.3.6
```

# DevSpace 基本概念
在深入了解 DevSpace 的使用方法前，我們先來了解幾個核心概念：
## 1. 開發工作流程
DevSpace 的工作流程包含以下階段：

Build：建置應用程式的容器映像檔
Deploy：將映像檔部署到 Kubernetes 叢集
Dev：啟動開發模式，監控檔案變更並同步到容器
Sync：將本地檔案變更直接同步到運行中的容器
Port-forwarding：自動設置連接埠轉發，方便本地訪問應用程式
UI：提供網頁界面監控應用狀態和日誌

## 2. 配置檔案 (devspace.yaml)
DevSpace 使用 devspace.yaml 檔案來定義整個工作流程。這個檔案指定了如何建置、部署和開發你的應用程式。

## 3. 運行模式
DevSpace 有幾種主要運行模式：

dev 模式：devspace dev - 持續監控檔案變更，自動觸發重建和重新部署
deploy 模式：devspace deploy - 執行一次完整的部署流程，適合 CI/CD 環境
UI 模式：devspace ui - 啟動 DevSpace UI 儀表板，監控應用狀態

## 建立基本的 devspace.yaml
現在，讓我們建立一個基本的 devspace.yaml 檔案。假設我們有一個簡單的 Node.js 應用程式，目錄結構如下：
```
my-node-app/
├── app.js
├── Dockerfile
├── package.json
└── kubernetes/
    └── deployment.yaml
```

以下是一個基本的 devspace.yaml 配置：

```yaml
version: v2beta1
name: my-node-app

# 定義要構建的映像檔
images:
  my-node-app:
    image: my-node-app
    dockerfile: ./Dockerfile

# 部署方式設定
deployments:
  my-deployment:
    kubectl:
      manifests:
        - kubernetes/deployment.yaml
        
# 開發模式設定
dev:
  # 自動同步本地檔案到容器
  sync:
    - imageSelector: my-node-app
      localSubPath: ./
      containerPath: /app
      excludePaths:
        - node_modules/
  # 自動轉發端口
  ports:
    - imageSelector: my-node-app
      forward:
        - port: 3000
          localPort: 8080
```
這個配置檔案告訴 DevSpace：

建置：使用 Dockerfile 建置名為 my-node-app 的映像檔
部署：使用 kubectl 部署 kubernetes/deployment.yaml 定義的資源
開發：監控本地檔案變更，並自動同步到容器的 /app 目錄
連接埠轉發：將容器的 3000 連接埠轉發到本地的 8080 連接埠

## 配置檔案詳解
讓我們更詳細地了解 devspace.yaml 的各個部分：

version
指定 DevSpace 配置檔案的版本。
```version: v2beta1```

name
指定專案名稱。
```name: my-node-app```

images
定義如何建置應用程式的映像檔。
```yaml
images:
  my-node-app:  # 映像檔的識別名稱
    image: my-node-app  # 實際的映像檔名稱
    dockerfile: ./Dockerfile  # Dockerfile 路徑
    context: ./  # 建置上下文
```

DevSpace 支援多種建置方式，包括：

Docker：使用 Docker 建置映像檔
Kaniko：在 Kubernetes 叢集中建置映像檔
Custom：使用自定義指令建置映像檔

deployments
定義如何將應用程式部署到 Kubernetes。
```yaml
deployments:
  my-deployment:  # 部署的識別名稱
    kubectl:
      manifests:
        - kubernetes/deployment.yaml  # Kubernetes 資源定義檔案
```
DevSpace 支援多種部署方式，包括：

kubectl：使用 kubectl 命令行工具
helm：使用 Helm 包管理器
kustomize：使用 Kustomize 客製化工具


dev
定義開發模式的行為。
```yaml
dev:
  # 自動同步本地檔案到容器
  sync:
    - imageSelector: my-node-app  # 要同步的容器
      localSubPath: ./  # 本地路徑
      containerPath: /app  # 容器內路徑
      excludePaths:
        - node_modules/  # 排除的路徑
  # 自動轉發端口
  ports:
    - imageSelector: my-node-app
      forward:
        - port: 3000  # 容器端口
          localPort: 8080  # 本地端口
```

## DevSpace UI 儀表板
DevSpace 的一個獨特功能是內建的 UI 儀表板，它提供了一個直觀的界面來監控應用狀態、查看日誌和管理開發環境。

要啟動 UI 儀表板，只需執行：```devspace ui```

UI 儀表板提供以下功能：

應用狀態監控：查看 Pod、Deployment 等資源的狀態
日誌查看：集中查看所有容器的日誌
容器終端：直接在瀏覽器中連接到容器終端
連接埠轉發管理：查看和管理所有轉發的連接埠
事件監控：查看 Kubernetes 事件

## 使用 DevSpace 開發模式
DevSpace 的開發模式是它最強大的功能之一。在這個模式下，DevSpace 會：

監控原始碼變更
自動重新建置映像檔（如需要）
自動重新部署到 Kubernetes（如需要）
直接同步檔案到容器（如已配置）
提供即時日誌輸出
設置連接埠轉發

要啟動開發模式，只需在專案目錄中執行：
```devspace dev```

### 開發模式工作流程
當你執行 devspace dev 時，以下是發生的事情：

DevSpace 讀取 devspace.yaml 配置檔案
建置所有定義的映像檔
部署應用程式到 Kubernetes 叢集
設置連接埠轉發
開始監控原始碼變更
當偵測到變更時，根據配置同步檔案或重新建置映像檔
顯示應用程式的日誌輸出

這個循環會一直持續，直到你按下 Ctrl+C 停止 DevSpace。當停止時，DevSpace 會自動清理它部署的資源（如有配置）。

## DevSpace 與 VS Code 整合
DevSpace 可以與 VS Code 無縫整合，提供更好的開發體驗。

### VS Code 擴充功能
首先，安裝以下 VS Code 擴充功能：

"Kubernetes" 擴充功能
"Remote - Containers" 擴充功能

### 使用 VS Code 的 Tasks 整合 DevSpace
在專案根目錄創建 .vscode/tasks.json 檔案：
```yaml
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "DevSpace: Start Development",
      "type": "shell",
      "command": "devspace dev",
      "problemMatcher": [],
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "DevSpace: Deploy",
      "type": "shell",
      "command": "devspace deploy",
      "problemMatcher": [],
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "DevSpace: Open UI",
      "type": "shell",
      "command": "devspace ui",
      "problemMatcher": [],
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    }
  ]
}
```

現在，你可以通過 VS Code 的命令面板（Ctrl+Shift+P）執行這些任務，選擇 "Tasks: Run Task" 然後選擇相應的 DevSpace 任務。

### 使用 VS Code 的 Launch 配置進行除錯
創建 .vscode/launch.json 檔案，設定遠端除錯：
```yaml
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Node.js: Remote Debug",
      "type": "node",
      "request": "attach",
      "port": 9229,
      "address": "localhost",
      "localRoot": "${workspaceFolder}",
      "remoteRoot": "/app",
      "preLaunchTask": "DevSpace: Start Development"
    }
  ]
}
```


這個配置允許你使用 VS Code 的除錯器連接到在 Kubernetes 中運行的 Node.js 應用程式。

# 實作：使用 DevSpace 開發 Node.js 應用程式
現在，讓我們通過一個實際的例子來演示如何使用 DevSpace 開發 Node.js 應用程式。

步驟 1：準備應用程式
首先，我們建立一個簡單的 Node.js 應用程式的dockerfile：
```dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 3000

CMD ["node", "app.js"]
```

步驟 2：建立 Kubernetes 部署檔案
建立 kubernetes 目錄並新增部署檔案：

```mkdir -p kubernetes```

建立 kubernetes/deployment.yaml：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: devspace-demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: devspace-demo
  template:
    metadata:
      labels:
        app: devspace-demo
    spec:
      containers:
      - name: devspace-demo
        image: devspace-demo
        ports:
        - containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: devspace-demo
spec:
  selector:
    app: devspace-demo
  ports:
  - port: 80
    targetPort: 3000
  type: ClusterIP
```

步驟 3：建立 devspace.yaml
在專案根目錄建立 devspace.yaml：
```yaml
version: v2beta1
name: devspace-demo

images:
  devspace-demo:
    image: devspace-demo
    dockerfile: ./Dockerfile

deployments:
  devspace-demo:
    kubectl:
      manifests:
        - kubernetes/deployment.yaml

dev:
  sync:
    - imageSelector: devspace-demo
      localSubPath: ./
      containerPath: /app
      excludePaths:
        - node_modules/
  ports:
    - imageSelector: devspace-demo
      forward:
        - port: 3000
          localPort: 8080
  open:
    - url: http://localhost:8080
```

注意，我們添加了 open 部分，這會在開發模式啟動後自動打開瀏覽器訪問應用程式。

步驟 4：啟動 DevSpace 開發模式
確保你的 KinD 叢集正在運行，然後執行：```devspace dev```

DevSpace 會執行以下操作：

建置 Docker 映像檔
將映像檔載入到 KinD 叢集
部署應用程式
設置連接埠轉發
開始監控檔案變更
打開瀏覽器訪問應用程式

步驟 5：測試開發循環
現在，讓我們修改 app.js 檔案，看看 DevSpace 如何自動更新應用程式：

保存檔案後，你應該會看到 DevSpace 自動：

偵測到檔案變更
將變更的檔案同步到容器
應用程式立即反映變更（無需重建映像檔）
打開瀏覽器訪問 http://localhost:8080，你應該會看到更新後的訊息。

步驟 6：使用 DevSpace UI
在另一個終端視窗執行：
```devspace ui```

這會啟動 DevSpace 的 UI 儀表板，你可以在瀏覽器中查看：

應用程式的部署狀態
容器日誌
事件和錯誤
資源使用情況
連接埠轉發狀態


## DevSpace 進階功能
DevSpace 還提供了許多進階功能，可以進一步優化開發工作流程：

1. 檔案同步
DevSpace 的檔案同步功能可以直接將變更的檔案同步到正在運行的容器中，無需重建映像檔：
```yaml
dev:
  sync:
    - imageSelector: devspace-demo
      localSubPath: ./
      containerPath: /app
      excludePaths:
        - node_modules/
        - .git/
      uploadExcludePaths:
        - dist/
      downloadExcludePaths:
        - logs/
```

2. 多環境配置
DevSpace 支援使用 profiles 來管理不同環境的配置：
```yaml
profiles:
  - name: production
    patches:
      - op: replace
        path: images.devspace-demo.image
        value: my-registry/devspace-demo:prod
  - name: staging
    patches:
      - op: replace
        path: deployments.devspace-demo.kubectl.manifests
        value:
          - kubernetes/staging-deployment.yaml
```

使用 devspace dev --profile=staging 可以啟用特定的環境配置。

3. 多服務開發
DevSpace 可以同時處理多個服務的開發：
```yaml
images:
  frontend:
    image: frontend
    dockerfile: ./frontend/Dockerfile
    context: ./frontend
  backend:
    image: backend
    dockerfile: ./backend/Dockerfile
    context: ./backend

deployments:
  app:
    kubectl:
      manifests:
        - kubernetes/frontend.yaml
        - kubernetes/backend.yaml

dev:
  sync:
    - imageSelector: frontend
      localSubPath: ./frontend
      containerPath: /app
    - imageSelector: backend
      localSubPath: ./backend
      containerPath: /app
```

4. 命令執行與鉤子
DevSpace 支援在開發過程中執行命令和設置鉤子：

```yaml
dev:
  terminal:
    imageSelector: devspace-demo
    command: ["bash"]
  hooks:
    - command: "echo 'Development started'"
      container: devspace-demo
      when: before:dev
    - command: "npm run test"
      container: devspace-demo
      when: after:sync
```

# 總結
DevSpace 是一個強大的工具，可以大幅簡化 Kubernetes 應用程式的開發流程。通過自動化建置、部署和監控過程，它讓開發人員可以專注於程式碼開發，而不是繁瑣的部署工作。

今天，我們學習了：

DevSpace 的基本概念和安裝方法
如何建立和配置 devspace.yaml
使用 devspace dev 實現快速開發循環
DevSpace UI 儀表板的功能和使用方法
如何將 DevSpace 與 VS Code 整合
如何通過實際例子使用 DevSpace 開發 Node.js 應用程式
DevSpace 的一些進階功能
在下一篇文章中，我們將深入探討 DevSpace 的更多進階功能，如多環境配置、調試支援和與 CI/CD 系統的整合。

# Lab 練習
練習 1：安裝 DevSpace 並設置開發環境
按照文章中的說明安裝 DevSpace
確認安裝成功：devspace version
確保你有一個運行中的 KinD 叢集：kind get clusters
練習 2：使用 DevSpace 開發一個簡單的應用程式
按照文章中的步驟建立 Node.js 應用程式
建立 Kubernetes 部署檔案
建立 devspace.yaml 配置檔案
啟動 DevSpace 開發模式：devspace dev
修改應用程式程式碼，觀察 DevSpace 自動同步
使用瀏覽器或 curl 測試應用程式
練習 3：探索 DevSpace UI 和 VS Code 整合
啟動 DevSpace UI：devspace ui
探索 UI 儀表板的各個功能
設定 VS Code 任務和啟動配置
使用 VS Code 任務啟動 DevSpace 開發環境
嘗試使用 VS Code 除錯器連接到容器中的應用程式
練習 4：使用 DevSpace 的進階功能
修改 devspace.yaml，添加更詳細的檔案同步配置
創建不同環境的 profiles
嘗試設置命令鉤子，在特定事件時執行命令
比較使用檔案同步和不使用檔案同步時的開發體驗
通過這些練習，你將能夠親身體驗 DevSpace 如何改善 Kubernetes 應用程式的開發體驗，並為後續學習更進階的功能打下基礎。