---
title: "Day 9 ：自訂映像檔部署與開發工作流程"
tags: 2025鐵人賽
date: 2025-07-20
---

## 第二週：開發流程與 Skaffold 整合
#### Day 9：自訂映像檔部署與開發工作流程
##### 重點
如何讓 KinD 使用本地 Docker 映像檔（kind load docker-image）
開發工作流程問題：手動部署的繁瑣步驟
##### Lab
將本地映像檔直接用於 Deployment，驗證服務可用性
記錄開發過程中需要執行的重複步驟（建置→載入→部署→測試）


在前面的分享中，我們已經學習了如何使用 Kubernetes 部署應用程式，以及如何診斷和解決 Pod 相關問題。今天，我們將進一步深入探討開發流程，特別是如何在本地 Kubernetes 環境（KinD）中使用團隊自定義開發的 Docker 映像檔，以及開發過程中常見的工作流程挑戰。

# 為什麼需要在本地使用自訂映像檔？
在實際開發過程中，我們經常需要：

測試尚未發布到公共或私有倉庫的映像檔
快速迭代開發，避免每次都推送映像檔到遠端倉庫
在離線環境中工作
減少開發過程中的網絡依賴

這時，能夠直接在本地 Kubernetes 集群中使用本地建置的 Docker 映像檔就變得非常重要。


# KinD 與本地 Docker 映像檔
## KinD 的映像檔管理機制
KinD 是一個在 Docker 容器中運行 Kubernetes 集群的工具。由於 KinD 節點本身是 Docker 容器，它們有自己獨立的容器運行時環境，無法直接訪問主機上的 Docker 映像檔。

這意味著，當你在本地建置了一個 Docker 映像檔，KinD 集群內的容器運行時無法自動看到這個映像檔。你需要明確地將映像檔載入到 KinD 集群中。

## 使用 kind load 載入本地映像檔
KinD 提供了 kind load 命令，專門用於將本地 Docker 映像檔載入到 KinD 集群中：
```bash
# 語法
kind load docker-image <image-name>:<tag> --name <cluster-name>

# 範例：載入名為 my-app:v1 的映像檔到名為 kind 的集群
kind load docker-image my-app:v1 --name kind
```

如果你沒有指定 --name 參數，KinD 會使用默認的集群名稱 "kind"。



### 映像檔載入過程
當執行 kind load docker-image 命令時，KinD 會：

從本地 Docker daemon process 中獲取指定的映像檔
將映像檔打包並傳輸到 KinD 集群的節點容器中
在節點容器內解包映像檔，使其可用於 Kubernetes Pod
這個過程完成後，你就可以在 Kubernetes 部署文件中引用這個映像檔，而不需要從遠端倉庫拉取。

# 實作：自訂映像檔部署流程
讓我們通過一個實際的例子來演示如何建置自訂映像檔並部署到 KinD 集群中。

## 步驟 1：建立一個簡單的應用程式
首先，我們建立一個簡單的 Node.js 應用程式：

建立專案目錄：
```bash
mkdir my-node-app
cd my-node-app
```

建立 app.js 文件：
```javaascript
const http = require('http');

const server = http.createServer((req, res) => {
  res.statusCode = 200;
  res.setHeader('Content-Type', 'text/plain');
  res.end('Hello from my custom Docker image in KinD!\n');
});

const PORT = process.env.PORT || 3000;
server.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
```

建立 Dockerfile：
```dockerfile
FROM node:18-alpine

WORKDIR /app

COPY app.js .

EXPOSE 3000

CMD ["node", "app.js"]
```

## 步驟 2：建置 Docker 映像檔
使用 Docker 建置映像檔：
```bash
docker build -t my-node-app:v1 .
```

確認映像檔已成功建置：
```bash
docker images | grep my-node-app
```

## 步驟 3：載入映像檔到 KinD 集群
假設我們已經有一個運行中的 KinD 集群（如果沒有，可以使用 kind create cluster 建立一個）：
```bash
kind load docker-image my-node-app:v1
```

你應該會看到類似以下的輸出：

```
Image: "my-node-app:v1" with ID "sha256:1a2b3c4d..." not yet present on node "kind-control-plane", loading...
```

## 步驟 4：建立 Kubernetes Deploayment 文件
建立 deployment.yaml 文件：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-node-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-node-app
  template:
    metadata:
      labels:
        app: my-node-app
    spec:
      containers:
      - name: my-node-app
        image: my-node-app:v1
        ports:
        - containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: my-node-app
spec:
  selector:
    app: my-node-app
  ports:
  - port: 80
    targetPort: 3000
  type: ClusterIP
```

注意，這裡我們直接使用 my-node-app:v1 作為映像檔名稱，而沒有指定倉庫地址。這是因為我們已經將映像檔載入到 KinD 集群中，Kubernetes 可以直接使用它。

## 步驟 5：部署應用程式
```bash
kubectl apply -f deployment.yaml
```

確認部署狀態：
```bash
kubectl get pods
```

## 步驟 6：測試應用程式
使用 port-forward 來訪問應用程式：
```bash
kubectl port-forward svc/my-node-app 8080:80
```

現在，你可以在瀏覽器中訪問 http://localhost:8080

如果一切正常，你應該會看到 "Hello from my custom Docker image in KinD!" 的回應。


# 開發工作流程問題
在上面的例子中，我們經歷了從建置映像檔到部署應用程式的完整流程。但在實際開發過程中，這個流程會變得非常繁瑣，特別是當你需要頻繁修改代碼並測試時。

讓我們來記錄一下典型的開發迭代流程：

## 傳統開發迭代流程
修改代碼：更新應用程式代碼
建置映像檔：docker build -t my-app:v1 .
載入映像檔：kind load docker-image my-app:v1
重新部署：kubectl apply -f deployment.yaml
等待 Pod 就緒：kubectl get pods -w
測試應用程式：通過 port-forward 或其他方式測試
查看日誌：kubectl logs <pod-name>
重複以上步驟

這個流程有幾個明顯的問題：

重複性工作：每次修改代碼都需要執行相同的命令序列
等待時間長：建置映像檔、載入映像檔和重新部署都需要時間
手動操作多：容易出錯，特別是在處理多個服務時
版本管理複雜：需要手動管理映像檔標籤
缺乏自動化：無法自動檢測代碼變更並觸發重建


## 一個簡單的自動化腳本
為了減輕一些手動工作，我們可以建立一個簡單的腳本來自動執行這些步驟：

建立 deploy.sh 文件：
```shell
#!/bin/bash

# 設定變數
APP_NAME="my-node-app"
IMAGE_TAG="v1"
DEPLOYMENT_FILE="deployment.yaml"

# 建置映像檔
echo "Building Docker image..."
docker build -t ${APP_NAME}:${IMAGE_TAG} .

# 載入映像檔到 KinD
echo "Loading image to KinD..."
kind load docker-image ${APP_NAME}:${IMAGE_TAG}

# 部署或更新應用程式
echo "Deploying to Kubernetes..."
kubectl apply -f ${DEPLOYMENT_FILE}

# 等待 Pod 就緒
echo "Waiting for Pod to be ready..."
kubectl wait --for=condition=ready pod -l app=${APP_NAME} --timeout=60s

# 顯示 Pod 狀態
echo "Pod status:"
kubectl get pods -l app=${APP_NAME}

# 設定 port-forward（背景執行）
echo "Setting up port-forward..."
pkill -f "kubectl port-forward svc/${APP_NAME}" || true
kubectl port-forward svc/${APP_NAME} 8080:80 &

echo "Application is accessible at http://localhost:8080"
```

使腳本可執行並運行：
```bash
chmod +x deploy.sh
./deploy.sh
```

這個腳本雖然簡化了流程，但仍然有局限性：

不能自動檢測代碼變更
每次都會重建整個映像檔，即使只有小部分代碼變更
不支持多服務開發
缺乏開發環境與生產環境的配置差異管理

更好的解決方案：開發工具
實際上，有多種專門的工具可以解決這些開發工作流程問題：

Skaffold：Google 開發的工具，可以自動化 Kubernetes 應用程式的開發工作流程
Tilt：專注於多服務應用程式的開發工具
DevSpace：為 Kubernetes 設計的開發環境工具
Telepresence：允許本地代碼與遠端 Kubernetes 集群交互
在明天的分享中，我們將深入探討 Skaffold，並了解它如何徹底改變 Kubernetes 應用程式的開發體驗。


# 總結
今天，我們學習了：

為什麼在本地開發過程中需要使用自訂 Docker 映像檔
如何使用 kind load docker-image 將本地映像檔載入到 KinD 集群
完整的自訂映像檔部署流程
傳統開發工作流程中的問題和挑戰
如何使用簡單的腳本來自動化部分流程
這些知識為我們下一步學習更先進的開發工具奠定了基礎。在實際開發中，理解底層的工作原理和手動流程是非常重要的，即使最終我們會使用工具來自動化這些流程。

# Lab 練習
練習 1：自訂映像檔部署
建立一個簡單的 Web 應用程式（可以使用任何你熟悉的語言）
為應用程式建立 Dockerfile
建置映像檔並載入到 KinD 集群
建立 Kubernetes 部署文件並部署應用程式
驗證服務可用性
練習 2：記錄開發工作流程
修改你在練習 1 中建立的應用程式（例如，更改顯示的文字）
記錄從修改代碼到驗證變更所需的所有步驟
計算完成一次迭代所需的時間
思考如何改進這個流程
嘗試建立一個自動化腳本來簡化流程
在下一節課中，我們將介紹 DevSpace
