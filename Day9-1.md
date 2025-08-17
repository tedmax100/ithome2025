# Day 9: DevSpace 入門與開發工作流

## 🎯 學習目標

今天我們將學習 DevSpace，這是一個專為 Kubernetes 開發者設計的工具，能夠大幅提升開發效率。作為新手，我們會從最基礎的概念開始，逐步建立完整的開發工作流。

### 📚 學習重點
- 理解 DevSpace 的核心概念與優勢
- 安裝和配置 DevSpace 環境
- 創建第一個 DevSpace 項目
- 實現代碼熱重載功能
- 掌握本地調試技巧
- 整合現有開發工具

---

## 🤔 為什麼需要 DevSpace？

### 傳統 Kubernetes 開發的痛點

想像一下，沒有 DevSpace 時的開發流程：

```mermaid
graph TD
    A[修改代碼] --> B[重新構建 Docker 映像]
    B --> C[推送到映像倉庫]
    C --> D[更新 Kubernetes 配置]
    D --> E[重新部署到集群]
    E --> F[測試應用]
    F --> G{測試通過？}
    G -->|否| A
    G -->|是| H[完成]
    
    style A fill:#ffcdd2
    style B fill:#ffcdd2
    style C fill:#ffcdd2
    style D fill:#ffcdd2
    style E fill:#ffcdd2
```

**問題**：
- ⏰ 每次修改需要 5-10 分鐘的構建部署時間
- 🔄 重複性工作太多
- 🐛 調試困難
- 💸 資源浪費

### DevSpace 的解決方案

```mermaid
graph TD
    A[修改代碼] --> B[DevSpace 自動同步]
    B --> C[容器內熱重載]
    C --> D[立即測試]
    D --> E{測試通過？}
    E -->|否| A
    E -->|是| F[完成]
    
    style A fill:#c8e6c9
    style B fill:#c8e6c9
    style C fill:#c8e6c9
    style D fill:#c8e6c9
```

**優勢**：
- ⚡ 秒級的代碼同步
- 🔥 熱重載支持
- 🎯 直接在 Kubernetes 環境中開發
- 🛠️ 豐富的調試功能

---

## 📦 DevSpace 核心概念

### 架構概覽

```mermaid
graph TB
    subgraph "本地開發環境"
        DEV[開發者]
        CODE[源代碼]
        DEVSPACE[DevSpace CLI]
    end
    
    subgraph "Kubernetes 集群"
        subgraph "開發命名空間"
            POD[開發 Pod]
            SVC[Service]
            CONFIG[ConfigMap/Secret]
        end
    end
    
    DEV --> CODE
    CODE --> DEVSPACE
    DEVSPACE -->|文件同步| POD
    DEVSPACE -->|端口轉發| POD
    DEVSPACE -->|日誌流| POD
    
    style DEVSPACE fill:#e3f2fd
    style POD fill:#e8f5e8
```

### 核心功能

1. **文件同步 (File Sync)**
   - 本地代碼變更自動同步到容器
   - 支持雙向同步
   - 智能過濾不需要的文件

2. **端口轉發 (Port Forwarding)**
   - 自動設置本地到 Pod 的端口映射
   - 支持多端口轉發
   - 動態端口分配

3. **日誌流 (Log Streaming)**
   - 實時查看容器日誌
   - 多容器日誌聚合
   - 彩色輸出和過濾

4. **開發容器 (Dev Container)**
   - 為開發優化的容器配置
   - 包含開發工具和依賴
   - 支持遠程調試

---

## 🛠️ 實作 1：安裝和配置 DevSpace

### 安裝 DevSpace

```bash
# macOS (使用 Homebrew)
brew install devspace

# Linux/macOS (使用 curl)
curl -s -L "https://github.com/loft-sh/devspace/releases/latest" | sed -nE 's!.*"([^"]*devspace-linux-amd64)".*!https://github.com\1!p' | xargs -n 1 curl -L -o devspace && chmod +x devspace && sudo mv devspace /usr/local/bin

# Windows (使用 PowerShell)
md -Force "$Env:APPDATA\devspace"; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]'Tls,Tls11,Tls12'; Invoke-WebRequest -UseBasicParsing ((Invoke-WebRequest -URI "https://github.com/loft-sh/devspace/releases/latest" -UseBasicParsing).Content -replace ".*`"([^`"]*devspace-windows-amd64.exe)`".*","https://github.com/`$1") -o $Env:APPDATA\devspace\devspace.exe; $env:Path += ";" + $Env:APPDATA + "\devspace"; [Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::User)
```

### 驗證安裝

```bash
# 檢查版本
devspace --version

# 檢查 Kubernetes 連接
devspace list spaces
```

### 初始化配置

```bash
# 登錄 DevSpace (可選，用於雲端功能)
devspace login

# 設置默認配置
devspace use context  # 選擇 Kubernetes 上下文
```

---

## 🚀 實作 2：創建第一個 DevSpace 項目

### 項目結構

讓我們創建一個簡單的 Node.js Web 應用：

```bash
# 創建項目目錄
mkdir my-devspace-app
cd my-devspace-app

# 初始化 Node.js 項目
npm init -y
```

### 創建應用代碼

```javascript
// app.js
const express = require('express');
const app = express();
const port = process.env.PORT || 3000;

// 靜態文件服務
app.use(express.static('public'));

// API 路由
app.get('/api/hello', (req, res) => {
  res.json({ 
    message: 'Hello from DevSpace!', 
    timestamp: new Date().toISOString(),
    version: '1.0.0'
  });
});

// 健康檢查
app.get('/health', (req, res) => {
  res.json({ status: 'healthy', uptime: process.uptime() });
});

// 主頁面
app.get('/', (req, res) => {
  res.send(`
    <!DOCTYPE html>
    <html>
    <head>
        <title>DevSpace Demo</title>
        <style>
            body { font-family: Arial, sans-serif; margin: 40px; }
            .container { max-width: 600px; margin: 0 auto; }
            .button { 
                background: #007cba; 
                color: white; 
                padding: 10px 20px; 
                border: none; 
                border-radius: 5px; 
                cursor: pointer; 
            }
            .response { 
                margin-top: 20px; 
                padding: 15px; 
                background: #f5f5f5; 
                border-radius: 5px; 
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>🚀 DevSpace Demo App</h1>
            <p>這是一個使用 DevSpace 開發的應用示例</p>
            <button class="button" onclick="callAPI()">調用 API</button>
            <div id="response" class="response" style="display:none;"></div>
        </div>
        
        <script>
            async function callAPI() {
                try {
                    const response = await fetch('/api/hello');
                    const data = await response.json();
                    document.getElementById('response').style.display = 'block';
                    document.getElementById('response').innerHTML = 
                        '<h3>API 響應:</h3><pre>' + JSON.stringify(data, null, 2) + '</pre>';
                } catch (error) {
                    console.error('Error:', error);
                }
            }
        </script>
    </body>
    </html>
  `);
});

app.listen(port, () => {
  console.log(\`🚀 應用運行在端口 \${port}\`);
  console.log(\`📝 環境: \${process.env.NODE_ENV || 'development'}\`);
});
```

### 更新 package.json

```json
{
  "name": "my-devspace-app",
  "version": "1.0.0",
  "description": "DevSpace demo application",
  "main": "app.js",
  "scripts": {
    "start": "node app.js",
    "dev": "nodemon app.js",
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "dependencies": {
    "express": "^4.18.2"
  },
  "devDependencies": {
    "nodemon": "^3.0.1"
  }
}
```

### 創建 Dockerfile

```dockerfile
# Dockerfile
FROM node:18-alpine

WORKDIR /app

# 複製 package 文件
COPY package*.json ./

# 安裝依賴 (包括開發依賴，用於開發環境)
RUN npm install

# 複製應用代碼
COPY . .

# 暴露端口
EXPOSE 3000

# 開發模式啟動命令
CMD ["npm", "run", "dev"]
```

---

## ⚙️ 實作 3：創建 DevSpace 配置

### 初始化 DevSpace

```bash
# 在項目根目錄執行
devspace init
```

DevSpace 會引導你完成初始化：

```
? How do you want to deploy this project?
> Use existing Kubernetes manifests or kustomizations
  Use Helm to deploy this project  
  Use existing docker-compose.yaml
  Deploy using kubectl apply
```

選擇 "Deploy using kubectl apply"

### DevSpace 配置文件

初始化後會生成 `devspace.yaml`：

```yaml
# devspace.yaml
version: v2beta1
name: my-devspace-app

# 映像配置
images:
  app:
    image: my-devspace-app
    dockerfile: ./Dockerfile
    context: ./

# 部署配置
deployments:
  app:
    kubectl:
      manifests:
        - k8s/

# 開發配置
dev:
  app:
    # 選擇要開發的 Pod
    labelSelector:
      app: my-devspace-app
    
    # 文件同步
    sync:
      - path: ./
        container: /app
        excludePaths:
          - node_modules/
          - .git/
          - .devspace/
    
    # 端口轉發
    ports:
      - port: "3000"
    
    # 開啟終端
    terminal:
      enabled: true
    
    # 自動重啟
    restartHelper:
      inject: true

# 變量配置
vars:
  REGISTRY: ""
  IMAGE: "my-devspace-app"

# 配置文件
profiles:
  - name: production
    patches:
      - op: replace
        path: dev
        value: {}
      - op: replace
        path: images.app.dockerfile
        value: ./Dockerfile.prod
```

### 創建 Kubernetes 清單

```bash
# 創建 k8s 目錄
mkdir k8s
```

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-devspace-app
  labels:
    app: my-devspace-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-devspace-app
  template:
    metadata:
      labels:
        app: my-devspace-app
    spec:
      containers:
      - name: app
        image: my-devspace-app
        ports:
        - containerPort: 3000
        env:
        - name: NODE_ENV
          value: "development"
        - name: PORT
          value: "3000"
        # 開發環境配置
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        # 健康檢查
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: my-devspace-app-service
  labels:
    app: my-devspace-app
spec:
  selector:
    app: my-devspace-app
  ports:
  - port: 80
    targetPort: 3000
    protocol: TCP
  type: ClusterIP
```

---

## 🔥 實作 4：熱重載開發

### 啟動開發模式

```bash
# 啟動 DevSpace 開發模式
devspace dev
```

啟動過程：

```mermaid
graph TD
    A[devspace dev] --> B[構建映像]
    B --> C[部署到 Kubernetes]
    C --> D[啟動文件同步]
    D --> E[設置端口轉發]
    E --> F[開啟日誌流]
    F --> G[準備就緒]
    
    style G fill:#c8e6c9
```

### 測試熱重載

1. **修改應用代碼**：

```javascript
// 修改 app.js 中的 API 響應
app.get('/api/hello', (req, res) => {
  res.json({ 
    message: '🎉 Hello from DevSpace! (已更新)', 
    timestamp: new Date().toISOString(),
    version: '1.1.0',  // 更新版本號
    developer: 'Your Name'  // 新增字段
  });
});
```

2. **觀察自動同步**：

DevSpace 會自動檢測文件變更並同步到容器：

```
[info]   Sync: Upload 1 change(s) (Exclude: node_modules/, .git/, .devspace/)
[info]   Sync: Successfully uploaded 1 change(s)
[info]   Container restarted because of file change
```

3. **測試應用**：

在瀏覽器中訪問 `http://localhost:3000`，點擊 "調用 API" 按鈕，應該看到更新後的響應。

### 文件同步配置詳解

```yaml
# devspace.yaml 中的同步配置
dev:
  app:
    sync:
      - path: ./              # 本地路徑
        container: /app       # 容器內路徑
        excludePaths:         # 排除的路徑
          - node_modules/
          - .git/
          - .devspace/
          - "*.log"
        downloadExcludePaths: # 僅下載時排除
          - tmp/
        uploadExcludePaths:   # 僅上傳時排除
          - logs/
        onUpload:            # 上傳後執行的命令
          restartContainer: true
```

---

## 🐛 實作 5：本地調試

### 配置調試環境

1. **更新 package.json**：

```json
{
  "scripts": {
    "start": "node app.js",
    "dev": "nodemon app.js",
    "debug": "nodemon --inspect=0.0.0.0:9229 app.js"
  }
}
```

2. **更新 DevSpace 配置**：

```yaml
# devspace.yaml
dev:
  app:
    # ... 其他配置
    
    # 端口轉發 (添加調試端口)
    ports:
      - port: "3000"
      - port: "9229"  # Node.js 調試端口
    
    # 容器啟動命令覆蓋
    command: ["npm", "run", "debug"]
```

### VS Code 調試配置

創建 `.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "DevSpace Debug",
      "type": "node",
      "request": "attach",
      "port": 9229,
      "address": "localhost",
      "localRoot": "${workspaceFolder}",
      "remoteRoot": "/app",
      "protocol": "inspector",
      "restart": true,
      "skipFiles": [
        "<node_internals>/**"
      ]
    }
  ]
}
```

### 調試流程

```mermaid
graph TD
    A[啟動 DevSpace Dev] --> B[容器以調試模式運行]
    B --> C[端口轉發 9229]
    C --> D[VS Code 連接調試器]
    D --> E[設置斷點]
    E --> F[觸發請求]
    F --> G[斷點命中]
    G --> H[調試變量和執行流程]
    
    style D fill:#e3f2fd
    style G fill:#fff3e0
```

### 調試示例

1. **在 VS Code 中設置斷點**：
   - 在 `app.js` 的 API 處理函數中設置斷點

2. **啟動調試**：
   ```bash
   devspace dev
   ```

3. **連接調試器**：
   - 在 VS Code 中按 F5 或點擊調試按鈕

4. **觸發斷點**：
   - 在瀏覽器中訪問 API 端點

---

## 🔧 實作 6：整合開發工具

### 整合 Docker Compose

有時你可能需要本地服務（如數據庫）：

```yaml
# docker-compose.dev.yaml
version: '3.8'
services:
  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: developer
      POSTGRES_PASSWORD: devpass
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

### 環境變量管理

```yaml
# devspace.yaml
dev:
  app:
    # ... 其他配置
    
    # 環境變量
    env:
      - name: DATABASE_URL
        value: "postgresql://developer:devpass@host.docker.internal:5432/myapp"
      - name: REDIS_URL
        value: "redis://host.docker.internal:6379"
      - name: NODE_ENV
        value: "development"
      - name: DEBUG
        value: "app:*"
```

### 整合 Git Hooks

創建 `.devspace/hooks/before_deploy.sh`：

```bash
#!/bin/bash
# 部署前檢查
echo "🔍 執行部署前檢查..."

# 檢查代碼格式
npm run lint || {
  echo "❌ 代碼格式檢查失敗"
  exit 1
}

# 運行測試
npm test || {
  echo "❌ 測試失敗"
  exit 1
}

echo "✅ 部署前檢查通過"
```

在 `devspace.yaml` 中配置：

```yaml
hooks:
  - when:
      before:
        deployments: app
    execute:
      - .devspace/hooks/before_deploy.sh
```

---

## 📊 DevSpace 工作流程圖

### 完整開發流程

```mermaid
graph TD
    subgraph "開發者本地環境"
        A[編寫代碼]
        B[Git 提交]
        C[DevSpace Dev]
    end
    
    subgraph "DevSpace 處理"
        D[構建映像]
        E[部署到 K8s]
        F[文件同步]
        G[端口轉發]
        H[日誌流]
    end
    
    subgraph "Kubernetes 集群"
        I[開發 Pod]
        J[Service]
        K[ConfigMap]
    end
    
    subgraph "開發工具"
        L[VS Code 調試]
        M[瀏覽器測試]
        N[API 測試工具]
    end
    
    A --> C
    C --> D
    D --> E
    E --> F
    E --> G
    E --> H
    F --> I
    G --> L
    G --> M
    H --> A
    
    B --> A
    I --> J
    I --> K
    L --> I
    M --> J
    N --> J
    
    style C fill:#e3f2fd
    style I fill:#e8f5e8
    style L fill:#fff3e0
```

---

## 🔍 故障排除指南

### 調試命令參考

```bash
# 基本狀態檢查
devspace status
devspace list pods
devspace list services

# 日誌查看
devspace logs
devspace logs --follow
devspace logs --container <container-name>

# 進入容器
devspace enter
devspace enter --container <container-name>

# 同步管理
devspace sync
devspace reset sync

# 清理資源
devspace purge
devspace reset pods
```

### 實用調試技巧

#### 1. 詳細狀態檢查

```bash
# 查看完整狀態信息
devspace status --verbose

# 檢查特定組件狀態
devspace status --component sync
devspace status --component portforwarding
devspace status --component logs
```

#### 2. 高級日誌查看

```bash
# 查看最近 100 行日誌
devspace logs --tail 100

# 只查看錯誤日誌
devspace logs --level error

# 查看特定時間範圍的日誌
devspace logs --since 1h

# 多容器日誌並行查看
devspace logs --all-containers
```

#### 3. 文件同步調試

```bash
# 查看同步狀態詳情
devspace status sync --verbose

# 手動觸發同步
devspace sync --upload-only
devspace sync --download-only

# 重置並重新開始同步
devspace reset sync
devspace dev --sync-only
```

#### 4. 網路連接調試

```bash
# 測試端口轉發
devspace list ports

# 手動設置端口轉發
devspace port-forward 3000:3000

# 檢查 Service 連接
kubectl get svc -l app=my-devspace-app
kubectl describe svc my-devspace-app-service
```

---

## 🎯 最佳實踐與技巧

### 1. 開發環境優化

#### 資源限制配置

```yaml
# k8s/deployment.yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "200m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

#### 快速重啟配置

```yaml
# devspace.yaml
dev:
  app:
    restartHelper:
      inject: true
      path: /tmp/devspace-restart
    command: ["sh", "-c", "while true; do npm run dev; sleep 1; done"]
```

### 2. 多人協作配置

#### 命名空間隔離

```yaml
# devspace.yaml
vars:
  NAMESPACE: "dev-${DEVSPACE_USER}"

deployments:
  app:
    namespace: "${NAMESPACE}"
    kubectl:
      createNamespace: true
```

#### 共享配置

```yaml
# devspace.yaml
imports:
  - git: https://github.com/company/devspace-configs
    path: base-config.yaml
```

### 3. 性能優化

#### 構建緩存

```yaml
# devspace.yaml
images:
  app:
    buildKit:
      enabled: true
    cache:
      disabled: false
```

#### 選擇性同步

```yaml
# devspace.yaml
dev:
  app:
    sync:
      - path: ./src
        container: /app/src
      - path: ./package.json
        container: /app/package.json
        file: true  # 單文件同步
```

### 4. 安全考慮

#### 開發環境隔離

```yaml
# devspace.yaml
deployments:
  app:
    kubectl:
      manifests:
        - k8s/
    namespace: "dev-sandbox"
```

#### 敏感信息管理

```bash
# 使用 Secret 管理敏感信息
kubectl create secret generic app-secrets \
  --from-literal=db-password=secret123 \
  --from-literal=api-key=abc123
```

```yaml
# k8s/deployment.yaml
env:
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: app-secrets
      key: db-password
```

---

## 🚀 進階應用場景

### 1. 微服務開發

```yaml
# devspace.yaml - 多服務配置
images:
  frontend:
    image: myapp/frontend
    dockerfile: ./frontend/Dockerfile
  backend:
    image: myapp/backend
    dockerfile: ./backend/Dockerfile
  worker:
    image: myapp/worker
    dockerfile: ./worker/Dockerfile

deployments:
  frontend:
    kubectl:
      manifests:
        - k8s/frontend/
  backend:
    kubectl:
      manifests:
        - k8s/backend/
  worker:
    kubectl:
      manifests:
        - k8s/worker/

dev:
  frontend:
    labelSelector:
      app: frontend
    ports:
      - port: "3000"
    sync:
      - path: ./frontend/src
        container: /app/src
  
  backend:
    labelSelector:
      app: backend
    ports:
      - port: "8080"
      - port: "9229"  # 調試端口
    sync:
      - path: ./backend/src
        container: /app/src
```

### 2. 數據庫開發

```yaml
# devspace.yaml - 包含數據庫
dependencies:
  - name: postgres
    git: https://github.com/company/postgres-k8s
    
dev:
  app:
    # 等待數據庫就緒
    wait:
      - name: postgres-service
        timeout: 300
    env:
      - name: DATABASE_URL
        value: "postgresql://user:pass@postgres-service:5432/mydb"
```

### 3. 前端開發

```yaml
# devspace.yaml - React 應用
images:
  frontend:
    image: react-app
    dockerfile: ./Dockerfile.dev

dev:
  frontend:
    sync:
      - path: ./src
        container: /app/src
      - path: ./public
        container: /app/public
    ports:
      - port: "3000"
    env:
      - name: FAST_REFRESH
        value: "true"
      - name: WDS_SOCKET_HOST
        value: "localhost"
      - name: CHOKIDAR_USEPOLLING
        value: "true"
```

#### React 開發 Dockerfile

```dockerfile
# Dockerfile.dev
FROM node:18-alpine

WORKDIR /app

# 安裝依賴
COPY package*.json ./
RUN npm install

# 複製源代碼
COPY . .

# 暴露端口
EXPOSE 3000

# 開發模式啟動
CMD ["npm", "start"]
```

### 4. 全棧應用開發

```yaml
# devspace.yaml - 全棧配置
images:
  frontend:
    image: myapp/frontend
    dockerfile: ./frontend/Dockerfile.dev
  backend:
    image: myapp/backend
    dockerfile: ./backend/Dockerfile.dev

deployments:
  database:
    helm:
      chart:
        name: postgresql
        repo: https://charts.bitnami.com/bitnami
      values:
        auth:
          postgresPassword: devpass
          database: myapp
  
  backend:
    kubectl:
      manifests:
        - k8s/backend/
  
  frontend:
    kubectl:
      manifests:
        - k8s/frontend/

dev:
  backend:
    labelSelector:
      app: backend
    ports:
      - port: "8080"
      - port: "9229"
    sync:
      - path: ./backend/src
        container: /app/src
    env:
      - name: DATABASE_URL
        value: "postgresql://postgres:devpass@postgresql:5432/myapp"
  
  frontend:
    labelSelector:
      app: frontend
    ports:
      - port: "3000"
    sync:
      - path: ./frontend/src
        container: /app/src
    env:
      - name: REACT_APP_API_URL
        value: "http://localhost:8080"
```

---

## 📋 實用命令清單

### 基本操作

```bash
# 初始化項目
devspace init

# 開發模式
devspace dev

# 部署應用
devspace deploy

# 清理資源
devspace purge
```

### 調試命令

```bash
# 查看狀態
devspace status
devspace status --verbose

# 查看日誌
devspace logs
devspace logs --follow --tail 100

# 進入容器
devspace enter
devspace enter --container app

# 端口轉發
devspace port-forward 3000:3000
```

### 配置管理

```bash
# 使用配置文件
devspace use profile production
devspace deploy --profile production

# 變量管理
devspace set var IMAGE_TAG=v1.2.3
devspace print --var IMAGE_TAG

# 上下文切換
devspace use context minikube
devspace use namespace dev-team
```

### 高級操作

```bash
# 構建映像
devspace build
devspace build --tag latest

# 同步管理
devspace sync --upload-only
devspace sync --download-only
devspace reset sync

# 依賴管理
devspace update dependencies
devspace build dependencies
```

---

## 🔧 配置文件完整示例

### 生產級 devspace.yaml

```yaml
# devspace.yaml
version: v2beta1
name: my-production-app

# 變量定義
vars:
  REGISTRY: "registry.company.com"
  IMAGE_TAG: "dev-${DEVSPACE_RANDOM}"
  NAMESPACE: "dev-${DEVSPACE_USER}"

# 映像配置
images:
  app:
    image: "${REGISTRY}/my-app:${IMAGE_TAG}"
    dockerfile: ./Dockerfile
    context: ./
    buildKit:
      enabled: true
    cache:
      disabled: false

# 部署配置
deployments:
  database:
    helm:
      chart:
        name: postgresql
        repo: https://charts.bitnami.com/bitnami
        version: "12.1.9"
      values:
        auth:
          postgresPassword: "${DATABASE_PASSWORD}"
          database: myapp
        primary:
          persistence:
            size: 1Gi
  
  app:
    kubectl:
      manifests:
        - k8s/
      kustomize: true
    namespace: "${NAMESPACE}"

# 開發配置
dev:
  app:
    labelSelector:
      app: my-app
    
    # 文件同步
    sync:
      - path: ./src
        container: /app/src
        excludePaths:
          - "*.tmp"
          - "*.log"
      - path: ./config
        container: /app/config
        file: true
    
    # 端口轉發
    ports:
      - port: "3000"
      - port: "9229"
        bind: "127.0.0.1"
    
    # 環境變量
    env:
      - name: NODE_ENV
        value: "development"
      - name: DATABASE_URL
        value: "postgresql://postgres:${DATABASE_PASSWORD}@postgresql:5432/myapp"
      - name: LOG_LEVEL
        value: "debug"
    
    # 容器配置
    container: app
    command: ["npm", "run", "dev"]
    workingDir: /app
    
    # 終端訪問
    terminal:
      enabled: true
    
    # 自動重啟
    restartHelper:
      inject: true
      path: /tmp/devspace-restart

# 依賴配置
dependencies:
  - name: shared-config
    git: https://github.com/company/k8s-configs
    subPath: base/
    
# 鉤子配置
hooks:
  - when:
      before:
        deployments: app
    execute:
      - echo "🔍 執行部署前檢查..."
      - npm run lint
      - npm test
  - when:
      after:
        deployments: app
    execute:
      - echo "✅ 部署完成"

# 配置文件
profiles:
  - name: production
    patches:
      - op: replace
        path: dev
        value: {}
      - op: replace
        path: images.app.image
        value: "${REGISTRY}/my-app:${IMAGE_TAG}"
      - op: add
        path: images.app.build.buildArgs
        value:
          NODE_ENV: production
  
  - name: staging
    patches:
      - op: replace
        path: vars.NAMESPACE
        value: "staging"
      - op: replace
        path: deployments.database.helm.values.primary.persistence.size
        value: "5Gi"

# 本地開發配置
localRegistry:
  enabled: true
  name: localhost:5000
```

---

## 🎯 常見問題與解決方案

### 1. 文件同步問題

#### 問題：文件同步不工作

**症狀**：修改代碼後容器內沒有更新

**解決方案**：
```bash
# 檢查同步狀態
devspace status sync

# 重置同步
devspace reset pods
devspace dev

# 檢查排除規則
devspace print --var sync
```

#### 問題：同步速度慢

**症狀**：文件同步需要很長時間

**解決方案**：
```yaml
# 優化同步配置
dev:
  app:
    sync:
      - path: ./src
        container: /app/src
        excludePaths:
          - node_modules/
          - "*.log"
          - ".git/"
        # 使用批量上傳
        uploadExcludePaths:
          - "*.tmp"
```

### 2. 端口轉發問題

#### 問題：端口轉發失敗

**症狀**：無法訪問應用端口

**解決方案**：
```bash
# 檢查端口狀態
devspace status

# 手動設置端口轉發
devspace port-forward 3000:3000

# 檢查端口衝突
lsof -i :3000
```

### 3. 容器啟動問題

#### 問題：容器啟動失敗

**症狀**：Pod 一直處於 CrashLoopBackOff 狀態

**解決方案**：
```bash
# 查看詳細日誌
devspace logs --tail 100

# 檢查 Pod 狀態
kubectl describe pod $(devspace list pods -q)

# 進入容器調試
devspace enter --container app

# 檢查資源限制
kubectl top pods
```

### 4. 映像構建問題

#### 問題：映像構建失敗

**症狀**：devspace dev 時構建錯誤

**解決方案**：
```bash
# 清理構建緩存
devspace reset

# 強制重新構建
devspace dev --force-build

# 檢查 Dockerfile
docker build -t test .
```

---

## 📊 性能監控與優化

### 監控開發環境

```yaml
# k8s/monitoring.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: monitoring-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
    scrape_configs:
      - job_name: 'my-app'
        static_configs:
          - targets: ['my-app-service:3000']

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
      - name: prometheus
        image: prom/prometheus:latest
        ports:
        - containerPort: 9090
        volumeMounts:
        - name: config
          mountPath: /etc/prometheus
      volumes:
      - name: config
        configMap:
          name: monitoring-config
```

### 性能優化技巧

#### 1. 構建優化

```dockerfile
# 多階段構建
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:18-alpine AS runtime
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
```

#### 2. 同步優化

```yaml
# 選擇性同步
dev:
  app:
    sync:
      - path: ./src
        container: /app/src
        # 只同步必要的文件類型
        includePaths:
          - "*.js"
          - "*.jsx"
          - "*.ts"
          - "*.tsx"
          - "*.css"
          - "*.scss"
```

#### 3. 資源優化

```yaml
# k8s/deployment.yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "200m"
# 啟用水平擴展
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: my-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  minReplicas: 1
  maxReplicas: 3
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## 🚀 CI/CD 整合

### GitHub Actions 整合

```yaml
# .github/workflows/devspace.yml
name: DevSpace CI/CD

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run tests
      run: npm test
    
    - name: Run linting
      run: npm run lint

  deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Install DevSpace
      run: |
        curl -s -L "https://github.com/loft-sh/devspace/releases/latest" | sed -nE 's!.*"([^"]*devspace-linux-amd64)".*!https://github.com\1!p' | xargs -n 1 curl -L -o devspace
        chmod +x devspace
        sudo mv devspace /usr/local/bin
    
    - name: Setup Kubernetes
      uses: azure/k8s-set-context@v1
      with:
        method: kubeconfig
        kubeconfig: ${{ secrets.KUBE_CONFIG }}
    
    - name: Deploy with DevSpace
      run: |
        devspace use profile production
        devspace deploy --force-build
      env:
        DEVSPACE_ENV: production
```

### GitLab CI 整合

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"

test:
  stage: test
  image: node:18-alpine
  script:
    - npm ci
    - npm run test
    - npm run lint

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

deploy:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache curl
    - curl -s -L "https://github.com/loft-sh/devspace/releases/latest" | sed -nE 's!.*"([^"]*devspace-linux-amd64)".*!https://github.com\1!p' | xargs -n 1 curl -L -o devspace
    - chmod +x devspace && mv devspace /usr/local/bin
  script:
    - devspace use profile production
    - devspace set var IMAGE_TAG=$CI_COMMIT_SHA
    - devspace deploy
  only:
    - main
```

---

## 📋 總結

### 🎉 今天學到了什麼

通過今天的學習，我們完成了以下內容：

1. **DevSpace 基礎概念**
   - 理解了 DevSpace 解決的問題
   - 掌握了核心架構和功能
   - 了解了與傳統開發流程的差異

2. **實際操作技能**
   - 安裝和配置 DevSpace 環境
   - 創建完整的 DevSpace 項目
   - 實現代碼熱重載功能
   - 配置本地調試環境

3. **高級功能應用**
   - 整合外部開發工具
   - 多服務開發配置
   - 性能優化技巧
   - 故障排除方法

### 🔑 關鍵收穫

```mermaid
mindmap
  root((DevSpace 核心價值))
    開發效率
      秒級代碼同步
      熱重載支持
      即時反饋
    調試能力
      遠程調試
      實時日誌
      容器內調試
    團隊協作
      環境一致性
      配置共享
      命名空間隔離
    生產準備
      多環境配置
      CI/CD 整合
      部署自動化
```

### 💡 最佳實踐要點

1. **配置管理**
   - 使用 Profile 管理不同環境
   - 合理配置文件同步規則
   - 善用環境變量

2. **性能優化**
   - 啟用構建緩存
   - 選擇性文件同步
   - 合理設置資源限制

3. **團隊協作**
   - 命名空間隔離
   - 共享基礎配置
   - 統一開發流程

4. **故障排除**
   - 熟悉調試命令
   - 查看詳細日誌
   - 理解常見問題

### 🚀 下一步學習建議

1. **深入實踐**
   - 在實際項目中使用 DevSpace
   - 嘗試不同的配置選項
   - 探索高級功能

2. **擴展學習**
   - 學習 Helm 整合
   - 了解 CI/CD 流水線整合
   - 探索雲端開發環境

3. **社區參與**
   - 關注 DevSpace 官方文檔
   - 參與社區討論
   - 分享使用經驗

### 📈 學習進度檢查

完成今天學習後，你應該能夠：

- [ ] 解釋 DevSpace 的核心概念和優勢
- [ ] 獨立安裝和配置 DevSpace 環境
- [ ] 創建和配置 DevSpace 項目
- [ ] 實現代碼熱重載和本地調試
- [ ] 解決常見的開發問題
- [ ] 應用最佳實踐優化開發流程

### 🎯 實際應用場景

DevSpace 特別適合以下場景：

1. **微服務開發** - 多服務協同開發
2. **雲原生應用** - 直接在 K8s 環境開發
3. **團隊協作** - 統一開發環境
4. **CI/CD 整合** - 自動化部署流程
5. **遠程開發** - 雲端開發環境

通過今天的學習，你已經掌握了 DevSpace 的核心技能，可以開始在實際項目中應用這些知識，大幅提升 Kubernetes 應用的開發效率！

