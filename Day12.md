---
title: "Day 12：Volume 與持久化儲存（整合 DevSpace"
tags: 2025鐵人賽
date: 2025-07-20
---


#### Day 12：Volume 與持久化儲存（整合 DevSpace

##### 重點

Volume、PersistentVolume、PersistentVolumeClaim 概念
DevSpace 與持久化資料的處理

##### Lab

建立 PVC 並掛載於 Pod，測試資料持久性
在 DevSpace 開發循環中處理持久化資料
配置 DevSpace 的持久化資料同步功能

# Kubernetes 持久化儲存基礎
在 Kubernetes 中，容器是臨時性的，當容器重啟或 Pod 重新調度時，容器內的資料會丟失。為了解決這個問題，Kubernetes 提供了多種持久化儲存機制。

## Volume 類型
Kubernetes 支援多種 Volume 類型：

emptyDir: 臨時存儲，Pod 刪除時資料消失
hostPath: 掛載節點上的目錄，適用於單節點開發環境
configMap/secret: 掛載配置資料或敏感資訊
persistentVolumeClaim (PVC): 請求持久化儲存資源

## PersistentVolume (PV) 與 PersistentVolumeClaim (PVC)
Kubernetes 將儲存資源的供應與消費分離：

PersistentVolume (PV): 由管理員建立的儲存資源，與具體儲存實現相關
PersistentVolumeClaim (PVC): 使用者對儲存的請求，應用程式通過 PVC 使用 PV
PV 範例
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-example
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: standard
  hostPath:
    path: /data/pv-example
```

PVC 範例
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-example
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: standard
```

在 Pod 中使用 PVC
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-with-pvc
spec:
  containers:
  - name: app
    image: nginx
    volumeMounts:
    - mountPath: /data
      name: data-volume
  volumes:
  - name: data-volume
    persistentVolumeClaim:
      claimName: pvc-example
```

StorageClass 與動態佈建
StorageClass 允許動態佈建 PV，無需管理員手動建立：
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-ssd
```

# DevSpace 與持久化儲存整合
DevSpace 提供了多種方式來處理開發過程中的持久化資料。

在 DevSpace 中定義 PVC
在 devspace.yaml 中，我們可以定義部署時需要的 PVC：
```yaml
version: v2beta1
name: my-app

deployments:
  app:
    kubectl:
      manifests:
        - kubernetes/pvc.yaml
        - kubernetes/deployment.yaml
```

其中 pvc.yaml 定義了持久化卷宣告：


```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: app-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

而 deployment.yaml 使用了這個 PVC：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: app
        image: my-app
        volumeMounts:
        - name: app-data
          mountPath: /app/data
      volumes:
      - name: app-data
        persistentVolumeClaim:
          claimName: app-data
```

DevSpace 文件同步與持久化資料
DevSpace 的文件同步功能可以與容器中的持久化卷協同工作，確保本地開發的變更即時反映在容器中：
```yaml
version: v2beta1
name: my-app

images:
  app:
    image: my-app
    dockerfile: ./Dockerfile

deployments:
  app:
    kubectl:
      manifests:
        - kubernetes/pvc.yaml
        - kubernetes/deployment.yaml

dev:
  sync:
    - imageSelector: app
      localSubPath: ./src
      containerPath: /app/src
    - imageSelector: app
      localSubPath: ./data
      containerPath: /app/data  # 這裡對應 PVC 掛載點
      excludePaths:
        - "*.log"
        - "*.tmp"
```

在這個配置中，我們將本地的 ./data 目錄同步到容器中掛載 PVC 的 /app/data 目錄。這樣，即使 Pod 重啟，資料也會保留在 PVC 中。

## 持久化資料的開發模式
在開發過程中，我們可能需要不同的持久化資料處理策略：

開發模式：使用文件同步，本地資料更改直接反映到容器
測試模式：使用獨立的 PVC，模擬生產環境
生產模式：使用生產級別的儲存解決方案

DevSpace 的 profiles 功能可以幫助我們管理這些不同的模式：
```yaml
version: v2beta1
name: my-app

# 基本配置...

profiles:
  - name: dev
    patches:
      # 開發模式配置...
  
  - name: test
    patches:
      - op: replace
        path: deployments.app.kubectl.manifests
        value:
          - kubernetes/pvc-test.yaml
          - kubernetes/deployment-test.yaml
  
  - name: prod
    patches:
      - op: replace
        path: deployments.app.kubectl.manifests
        value:
          - kubernetes/pvc-prod.yaml
          - kubernetes/deployment-prod.yaml
```

# 實例：使用 DevSpace 開發資料庫應用
讓我們通過一個具體的例子來說明如何在 DevSpace 中處理持久化資料。假設我們正在開發一個使用 MongoDB 的應用。

項目結構
```
my-app/
├── src/
│   └── ...
├── kubernetes/
│   ├── mongodb-pvc.yaml
│   ├── mongodb.yaml
│   ├── app-deployment.yaml
│   └── app-service.yaml
├── Dockerfile
└── devspace.yaml
```

MongoDB PVC 定義
```yaml
# kubernetes/mongodb-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mongodb-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

MongoDB 部署

```yaml
# kubernetes/mongodb.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mongodb
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mongodb
  template:
    metadata:
      labels:
        app: mongodb
    spec:
      containers:
      - name: mongodb
        image: mongo:5.0
        ports:
        - containerPort: 27017
        volumeMounts:
        - name: mongodb-data
          mountPath: /data/db
      volumes:
      - name: mongodb-data
        persistentVolumeClaim:
          claimName: mongodb-data
---
apiVersion: v1
kind: Service
metadata:
  name: mongodb
spec:
  selector:
    app: mongodb
  ports:
  - port: 27017
    targetPort: 27017
```

應用部署
```yaml
# kubernetes/app-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: app
        image: my-app
        env:
        - name: MONGODB_URI
          value: "mongodb://mongodb:27017/mydb"
```

DevSpace 配置
```yaml
# devspace.yaml
version: v2beta1
name: my-app

images:
  app:
    image: my-app
    dockerfile: ./Dockerfile

deployments:
  mongodb:
    kubectl:
      manifests:
        - kubernetes/mongodb-pvc.yaml
        - kubernetes/mongodb.yaml
  app:
    kubectl:
      manifests:
        - kubernetes/app-deployment.yaml
        - kubernetes/app-service.yaml
    dependsOn:
      - mongodb

dev:
  sync:
    - imageSelector: app
      localSubPath: ./src
      containerPath: /app/src
  ports:
    - imageSelector: app
      forward:
        - port: 3000
          localPort: 3000
    - imageSelector: mongodb
      forward:
        - port: 27017
          localPort: 27017
```

這個配置做了以下幾件事：

部署 MongoDB 與對應的 PVC
部署我們的應用，並設置依賴關係（確保 MongoDB 先啟動）
設置文件同步，將本地 src 目錄同步到容器
設置端口轉發，允許從本地訪問應用和數據庫

# 使用 DevSpace 開發
啟動開發環境：`devspace dev`

DevSpace 會：

構建應用映像檔
部署 MongoDB 和 PVC
部署應用
設置文件同步和端口轉發
監控文件變更

由於 MongoDB 使用 PVC，即使 Pod 重啟，資料也會保留。同時，由於設置了端口轉發，我們可以從本地直接連接到 MongoDB：
`mongo localhost:27017/mydb`


## 進階：備份與恢復持久化資料
在開發過程中，我們可能需要備份和恢復數據庫資料。DevSpace 的 hooks 功能可以幫助我們實現這一點：
```yaml
version: v2beta1
name: my-app

# 基本配置...

hooks:
  - command: ["kubectl", "exec", "deployment/mongodb", "--", "mongodump", "--archive=/tmp/backup.gz", "--gzip"]
    name: backup-db
    when: before:deploy:exit
    where: local
  - command: ["kubectl", "cp", "mongodb-pod:/tmp/backup.gz", "./backups/backup-$(date +%Y%m%d%H%M%S).gz"]
    name: save-backup
    when: after:backup-db
    where: local
```

這個配置在 DevSpace 退出前執行數據庫備份，並將備份文件保存到本地。

# Lab 實作：使用 DevSpace 開發有狀態應用
實作 1：建立 PVC 並測試資料持久性
創建一個簡單的 PVC 定義：
```yaml
# pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: data-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

創建一個使用 PVC 的 Pod：

```yaml
# pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: data-pod
spec:
  containers:
  - name: busybox
    image: busybox
    command: ["sh", "-c", "while true; do sleep 3600; done"]
    volumeMounts:
    - name: data-volume
      mountPath: /data
  volumes:
  - name: data-volume
    persistentVolumeClaim:
      claimName: data-pvc
```

部署並測試：
```bash
kubectl apply -f pvc.yaml
kubectl apply -f pod.yaml

# 創建一個文件
kubectl exec data-pod -- sh -c "echo 'Hello PVC' > /data/test.txt"

# 驗證文件存在
kubectl exec data-pod -- cat /data/test.txt

# 刪除 Pod
kubectl delete pod data-pod

# 重新創建 Pod
kubectl apply -f pod.yaml

# 驗證文件是否還存在
kubectl exec data-pod -- cat /data/test.txt
```

實作 2：在 DevSpace 開發循環中處理持久化資料
創建一個簡單的 Node.js 應用，使用文件系統存儲資料：
```javascript
// app.js
const express = require('express');
const fs = require('fs');
const path = require('path');

const app = express();
const dataDir = process.env.DATA_DIR || '/app/data';

// 確保資料目錄存在
if (!fs.existsSync(dataDir)) {
  fs.mkdirSync(dataDir, { recursive: true });
}

app.use(express.json());

// 獲取所有筆記
app.get('/notes', (req, res) => {
  const notesFile = path.join(dataDir, 'notes.json');
  
  if (!fs.existsSync(notesFile)) {
    return res.json([]);
  }
  
  const notes = JSON.parse(fs.readFileSync(notesFile, 'utf8'));
  res.json(notes);
});

// 添加筆記
app.post('/notes', (req, res) => {
  const { title, content } = req.body;
  const notesFile = path.join(dataDir, 'notes.json');
  
  let notes = [];
  if (fs.existsSync(notesFile)) {
    notes = JSON.parse(fs.readFileSync(notesFile, 'utf8'));
  }
  
  const newNote = {
    id: Date.now(),
    title,
    content,
    createdAt: new Date().toISOString()
  };
  
  notes.push(newNote);
  fs.writeFileSync(notesFile, JSON.stringify(notes, null, 2));
  
  res.status(201).json(newNote);
});

const port = process.env.PORT || 3000;
app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});
```

Dockerfile
```
FROM node:16-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 3000

CMD ["node", "app.js"]
```

創建 Kubernetes 配置：
```yaml
# kubernetes/pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: notes-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notes-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: notes-app
  template:
    metadata:
      labels:
        app: notes-app
    spec:
      containers:
      - name: notes-app
        image: notes-app
        ports:
        - containerPort: 3000
        env:
        - name: DATA_DIR
          value: /app/data
        volumeMounts:
        - name: notes-data
          mountPath: /app/data
      volumes:
      - name: notes-data
        persistentVolumeClaim:
          claimName: notes-data
```

```yaml
# kubernetes/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: notes-app
spec:
  selector:
    app: notes-app
  ports:
  - port: 80
    targetPort: 3000
  type: ClusterIP
```

創建 DevSpace 配置：
```yaml
# devspace.yaml
version: v2beta1
name: notes-app

images:
  app:
    image: notes-app
    dockerfile: ./Dockerfile

deployments:
  app:
    kubectl:
      manifests:
        - kubernetes/pvc.yaml
        - kubernetes/deployment.yaml
        - kubernetes/service.yaml

dev:
  sync:
    - imageSelector: app
      localSubPath: ./
      containerPath: /app
      excludePaths:
        - node_modules/
        - .git/
  ports:
    - imageSelector: app
      forward:
        - port: 3000
          localPort: 3000
  open:
    - url: http://localhost:3000/notes
```

啟動 DevSpace：`devspace dev`

測試應用：
```
# 添加筆記
curl -X POST http://localhost:3000/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"測試筆記","content":"這是一個測試筆記"}'

# 獲取所有筆記
curl http://localhost:3000/notes
```

重啟 Pod 並驗證資料持久性：
```bash
kubectl delete pod -l app=notes-app
# 等待 Pod 重新啟動

# 再次獲取所有筆記
curl http://localhost:3000/notes
```


實作 3：配置 DevSpace 的持久化資料同步功能
在某些情況下，我們可能希望將本地資料目錄與容器中的持久化卷同步，以便在本地查看和修改持久化資料。

修改 DevSpace 配置：
```yaml
# devspace.yaml
version: v2beta1
name: notes-app

images:
  app:
    image: notes-app
    dockerfile: ./Dockerfile

deployments:
  app:
    kubectl:
      manifests:
        - kubernetes/pvc.yaml
        - kubernetes/deployment.yaml
        - kubernetes/service.yaml

dev:
  sync:
    - imageSelector: app
      localSubPath: ./
      containerPath: /app
      excludePaths:
        - node_modules/
        - .git/
        - data/
    - imageSelector: app
      localSubPath: ./data
      containerPath: /app/data
      onUpload:
        restartContainer: false
  ports:
    - imageSelector: app
      forward:
        - port: 3000
          localPort: 3000
  open:
    - url: http://localhost:3000/notes
```

創建本地資料目錄：`mkdir -p data`
啟動 DevSpace：`devspace dev`

添加一些筆記，然後檢查本地資料目錄：
```bash
# 編輯本地文件
# 修改 data/notes.json 中的某個筆記標題

# 查看變更是否反映在應用中
curl http://localhost:3000/notes

```


# 總結
在本文中，我們深入探討了 Kubernetes 中的持久化儲存概念，以及如何在 DevSpace 開發流程中處理持久化資料。我們學習了：

Kubernetes 中的 Volume、PersistentVolume 和 PersistentVolumeClaim 概念
如何在 DevSpace 中配置和使用持久化儲存
如何通過文件同步功能與持久化資料交互
如何使用 profiles 為不同環境配置不同的持久化策略
如何備份和恢復持久化資料
通過這些知識和實踐，我們可以在 Kubernetes 環境中有效地開發和管理有狀態應用，並確保資料的持久性和一致性。

在實際開發中，持久化儲存的選擇應該根據應用的需求和環境的特性來決定。在本地開發環境中，簡單的 hostPath 或 emptyDir 可能已經足夠；而在生產環境中，我們可能需要使用更可靠的儲存解決方案，如雲服務提供商的持久化磁碟或分佈式儲存系統。

無論選擇哪種儲存解決方案，DevSpace 都能幫助我們簡化開發流程，並確保開發環境與生產環境的一致性。