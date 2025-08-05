
# Day 16：Helm Chart 開發深入（下）

## 重點

- 依賴管理與子 Chart
- Hook 機制與生命週期
- Chart 測試與驗證

## 學習目標

- 掌握複雜應用的依賴關係管理
- 理解並實作 Helm Hook 機制
- 建立完整的 Chart 測試策略
- 開發企業級的 Helm Chart

---

## 1. 依賴管理與子 Chart

### 為什麼需要依賴管理？

想像你在開發一個電商應用，你需要：

- 前端服務（React 應用）
- 後端 API（Node.js 服務）
- 資料庫（PostgreSQL）
- 快取服務（Redis）
- 訊息佇列（RabbitMQ）

如果每個服務都要單獨管理，你會面臨：

- 配置重複且容易出錯
- 版本不一致的問題
- 部署順序難以控制
- 環境間的差異難以管理

### 依賴關係的設計思維

依賴管理有兩種主要方式：

#### 1. 外部依賴（推薦用於基礎設施）

```yaml
# Chart.yaml
apiVersion: v2
name: my-ecommerce-app
description: 完整的電商應用
version: 0.1.0
appVersion: "1.0.0"

dependencies:
  # 資料庫 - 使用官方 Chart
  - name: postgresql
    version: "~12.1.0"
    repository: "https://charts.bitnami.com/bitnami"
    condition: postgresql.enabled
    tags:
      - database
  
  # Redis 快取
  - name: redis
    version: "~17.3.0" 
    repository: "https://charts.bitnami.com/bitnami"
    condition: redis.enabled
    tags:
      - cache
  
  # 訊息佇列
  - name: rabbitmq
    version: "~11.1.0"
    repository: "https://charts.bitnami.com/bitnami"
    condition: rabbitmq.enabled
    tags:
      - messaging
```

**為什麼這樣設計？**

- **condition**：讓你可以選擇性安裝（例如開發環境可能使用外部資料庫）
- **tags**：可以按類別批次啟用/停用
- **波浪號版本**：接受相容的補丁版本，但避免破壞性更新

#### 2. 子 Chart（推薦用於應用元件）

```yaml
# 主 Chart 結構
my-ecommerce-app/
├── Chart.yaml
├── values.yaml
├── charts/
│   ├── frontend/          # 前端子 Chart
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   ├── backend-api/       # 後端 API 子 Chart
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   └── admin-panel/       # 管理後台子 Chart
└── templates/
    └── ingress.yaml       # 統一的 Ingress 配置
```

**子 Chart 的優勢：**

- 每個服務有獨立的配置和模板
- 可以單獨開發和測試
- 重用性高（其他專案也可以使用）
- 版本控制更精確

### 實際的依賴配置策略

```yaml
# values.yaml - 主 Chart 的配置
global:
  imageRegistry: "my-registry.com"
  environment: "production"
  
# 資料庫配置
postgresql:
  enabled: true
  auth:
    postgresPassword: "secure-password"
    database: "ecommerce_db"
  primary:
    persistence:
      enabled: true
      size: 10Gi
      storageClass: "fast-ssd"

# Redis 配置  
redis:
  enabled: true
  auth:
    enabled: true
    password: "redis-password"
  master:
    persistence:
      enabled: true
      size: 2Gi

# 前端子 Chart 配置
frontend:
  image:
    repository: my-registry.com/frontend
    tag: "v1.2.0"
  service:
    type: ClusterIP
    port: 80
  env:
    API_URL: "https://api.mystore.com"
  
# 後端 API 子 Chart 配置
backend-api:
  image:
    repository: my-registry.com/backend-api
    tag: "v1.2.0"
  service:
    type: ClusterIP
    port: 3000
  env:
    DATABASE_URL: "postgresql://user:pass@postgresql:5432/ecommerce_db"
    REDIS_URL: "redis://redis:6379"
```

**配置的層次思維：**

- **Global**：影響所有子 Chart 的設定
- **依賴層級**：每個依賴服務的專屬配置
- **子 Chart 層級**：應用元件的個別設定

### 依賴管理的最佳實踐

```bash
# 更新依賴
helm dependency update

# 檢查依賴狀態
helm dependency list

# 建構包含依賴的 Chart
helm dependency build
```

**重要觀念：**

- 資料庫、快取等基礎設施用外部依賴
- 應用元件用子 Chart
- 開發環境可以停用某些依賴（使用外部服務）

---

## 2. Hook 機制與生命週期

### 什麼是 Helm Hook？

Hook 就像是「事件觸發器」，讓你可以在特定時機執行特定任務。想像你在部署一個應用：

1. **部署前**：需要檢查資料庫是否可用
2. **部署中**：正常建立所有資源
3. **部署後**：需要執行資料庫遷移
4. **升級前**：需要備份現有資料
5. **刪除前**：需要清理外部資源

### Hook 的生命週期階段

Helm 提供了完整的生命週期 Hook：

```yaml
# templates/pre-install-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ include "my-app.fullname" . }}-pre-install"
  labels:
    {{- include "my-app.labels" . | nindent 4 }}
  annotations:
    # 這是關鍵：定義 Hook 類型和權重
    "helm.sh/hook": pre-install
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": before-hook-creation,hook-succeeded
spec:
  template:
    metadata:
      name: "{{ include "my-app.fullname" . }}-pre-install"
    spec:
      restartPolicy: Never
      containers:
      - name: pre-install-job
        image: "{{ .Values.hooks.image.repository }}:{{ .Values.hooks.image.tag }}"
        command:
          - /bin/sh
          - -c
          - |
            echo "檢查資料庫連線..."
            until nc -z {{ .Values.database.host }} {{ .Values.database.port }}; do
              echo "等待資料庫啟動..."
              sleep 2
            done
            echo "資料庫連線成功！"
```

**Hook 註釋說明：**

- `helm.sh/hook: pre-install`：在安裝前執行
- `helm.sh/hook-weight: "-5"`：執行順序（數字越小越早執行）
- `helm.sh/hook-delete-policy`：Hook 完成後的清理策略

### 實際應用場景的 Hook 實作

#### 1. 資料庫遷移 Hook

```yaml
# templates/post-install-migration.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ include "my-app.fullname" . }}-db-migration"
  annotations:
    "helm.sh/hook": post-install,post-upgrade
    "helm.sh/hook-weight": "1"
    "helm.sh/hook-delete-policy": before-hook-creation
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: db-migration
        image: "{{ .Values.app.image.repository }}:{{ .Values.app.image.tag }}"
        command:
          - /bin/sh
          - -c
          - |
            echo "開始資料庫遷移..."
            npm run migrate
            if [ $? -eq 0 ]; then
              echo "資料庫遷移成功"
            else
              echo "資料庫遷移失敗"
              exit 1
            fi
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: {{ include "my-app.fullname" . }}-db-secret
              key: database-url
```

**為什麼在 post-install 和 post-upgrade？**

- 新安裝時需要建立資料庫結構
- 升級時需要更新資料庫結構
- 在主應用啟動前完成，避免應用啟動失敗

#### 2. 備份 Hook

```yaml
# templates/pre-upgrade-backup.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ include "my-app.fullname" . }}-backup-{{ .Release.Revision }}"
  annotations:
    "helm.sh/hook": pre-upgrade
    "helm.sh/hook-weight": "-10"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: backup
        image: postgres:13
        command:
          - /bin/sh
          - -c
          - |
            echo "開始備份資料庫..."
            TIMESTAMP=$(date +%Y%m%d_%H%M%S)
            pg_dump $DATABASE_URL > /backup/backup_${TIMESTAMP}_rev_{{ .Release.Revision }}.sql
            echo "備份完成：backup_${TIMESTAMP}_rev_{{ .Release.Revision }}.sql"
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: {{ include "my-app.fullname" . }}-db-secret
              key: database-url
        volumeMounts:
        - name: backup-volume
          mountPath: /backup
      volumes:
      - name: backup-volume
        persistentVolumeClaim:
          claimName: {{ include "my-app.fullname" . }}-backup-pvc
```

**備份策略的考量：**

- 使用 Release.Revision 確保每次升級都有唯一備份
- Hook 成功後保留備份檔案（不自動刪除）
- 使用 PVC 確保備份持久化

#### 3. 清理 Hook

```yaml
# templates/pre-delete-cleanup.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ include "my-app.fullname" . }}-cleanup"
  annotations:
    "helm.sh/hook": pre-delete
    "helm.sh/hook-weight": "1"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: cleanup
        image: "{{ .Values.hooks.image.repository }}:{{ .Values.hooks.image.tag }}"
        command:
          - /bin/sh
          - -c
          - |
            echo "開始清理外部資源..."
          
            # 清理 S3 儲存桶中的檔案
            aws s3 rm s3://{{ .Values.storage.bucket }}/{{ .Release.Name }}/ --recursive
          
            # 清理外部 DNS 記錄
            curl -X DELETE "https://api.cloudflare.com/client/v4/zones/{{ .Values.dns.zoneId }}/dns_records/{{ .Values.dns.recordId }}" \
              -H "Authorization: Bearer {{ .Values.dns.apiToken }}"
          
            echo "外部資源清理完成"
        env:
        - name: AWS_ACCESS_KEY_ID
          valueFrom:
            secretKeyRef:
              name: {{ include "my-app.fullname" . }}-aws-secret
              key: access-key-id
        - name: AWS_SECRET_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: {{ include "my-app.fullname" . }}-aws-secret
              key: secret-access-key
```

### Hook 的執行順序與錯誤處理

```yaml
# values.yaml 中的 Hook 配置
hooks:
  image:
    repository: alpine/curl
    tag: latest
  
  # Hook 執行的超時設定
  timeout: 300
  
  # 錯誤處理策略
  failurePolicy: "fail"  # fail 或 ignore
```

**執行順序範例：**

1. `pre-install` (weight: -10) → 環境檢查
2. `pre-install` (weight: -5) → 依賴檢查
3. 主要資源建立
4. `post-install` (weight: 1) → 資料庫遷移
5. `post-install` (weight: 5) → 初始化資料

---

## 3. Chart 測試與驗證

### 為什麼 Chart 測試很重要？

想像你開發了一個複雜的 Chart，包含多個服務、資料庫、快取等。沒有測試的情況下：

- 部署到生產環境才發現配置錯誤
- 不同環境的行為不一致
- 升級時破壞現有功能
- 團隊成員修改 Chart 時引入 bug

### Helm 內建測試機制

Helm 提供了內建的測試框架，讓你可以驗證部署是否成功：

```yaml
# templates/tests/test-connection.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "my-app.fullname" . }}-test-connection"
  labels:
    {{- include "my-app.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test
    "helm.sh/hook-delete-policy": before-hook-creation,hook-succeeded
spec:
  restartPolicy: Never
  containers:
  - name: test-connection
    image: busybox:1.35
    command:
      - /bin/sh
      - -c
      - |
        echo "測試應用程式連線..."
      
        # 測試主服務是否可達
        if wget --spider --timeout=10 http://{{ include "my-app.fullname" . }}:{{ .Values.service.port }}/health; then
          echo "✅ 主服務連線成功"
        else
          echo "❌ 主服務連線失敗"
          exit 1
        fi
      
        # 測試資料庫連線
        if nc -z {{ .Values.database.host }} {{ .Values.database.port }}; then
          echo "✅ 資料庫連線成功"
        else
          echo "❌ 資料庫連線失敗"
          exit 1
        fi
      
        echo "🎉 所有連線測試通過！"
```

**執行測試：**

```bash
# 安裝 Chart
helm install my-app ./my-app

# 執行測試
helm test my-app

# 查看測試結果
kubectl logs -l "app.kubernetes.io/name=my-app,helm.sh/hook=test"
```

### 功能測試的完整實作

```yaml
# templates/tests/test-functionality.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "my-app.fullname" . }}-test-functionality"
  annotations:
    "helm.sh/hook": test
    "helm.sh/hook-weight": "2"
spec:
  restartPolicy: Never
  containers:
  - name: test-functionality
    image: curlimages/curl:7.85.0
    command:
      - /bin/sh
      - -c
      - |
        set -e
        echo "開始功能測試..."
      
        BASE_URL="http://{{ include "my-app.fullname" . }}:{{ .Values.service.port }}"
      
        # 測試健康檢查端點
        echo "測試健康檢查..."
        HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/health)
        if [ "$HEALTH_RESPONSE" = "200" ]; then
          echo "✅ 健康檢查通過"
        else
          echo "❌ 健康檢查失敗，HTTP 狀態碼：$HEALTH_RESPONSE"
          exit 1
        fi
      
        # 測試 API 端點
        echo "測試 API 端點..."
        API_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/api/v1/status)
        if [ "$API_RESPONSE" = "200" ]; then
          echo "✅ API 端點測試通過"
        else
          echo "❌ API 端點測試失敗，HTTP 狀態碼：$API_RESPONSE"
          exit 1
        fi
      
        # 測試資料庫連線（透過 API）
        echo "測試資料庫連線..."
        DB_RESPONSE=$(curl -s $BASE_URL/api/v1/db-status | grep -o '"status":"ok"')
        if [ "$DB_RESPONSE" = '"status":"ok"' ]; then
          echo "✅ 資料庫連線測試通過"
        else
          echo "❌ 資料庫連線測試失敗"
          exit 1
        fi
      
        echo "🎉 所有功能測試通過！"
```

### 效能測試

```yaml
# templates/tests/test-performance.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "my-app.fullname" . }}-test-performance"
  annotations:
    "helm.sh/hook": test
    "helm.sh/hook-weight": "3"
spec:
  restartPolicy: Never
  containers:
  - name: test-performance
    image: jordi/ab:latest
    command:
      - /bin/sh
      - -c
      - |
        echo "開始效能測試..."
      
        BASE_URL="http://{{ include "my-app.fullname" . }}:{{ .Values.service.port }}"
      
        # 使用 Apache Bench 進行負載測試
        echo "執行負載測試（100 個請求，10 個併發）..."
        ab -n 100 -c 10 -g /tmp/results.tsv $BASE_URL/health
      
        # 分析結果
        FAILED_REQUESTS=$(grep "Failed requests" /tmp/ab_output.txt | awk '{print $3}')
        REQUESTS_PER_SECOND=$(grep "Requests per second" /tmp/ab_output.txt | awk '{print $4}')
      
        echo "效能測試結果："
        echo "- 失敗請求數：$FAILED_REQUESTS"
        echo "- 每秒請求數：$REQUESTS_PER_SECOND"
      
        # 設定效能標準
        if [ "$FAILED_REQUESTS" -gt "5" ]; then
          echo "❌ 失敗請求過多"
          exit 1
        fi
      
        echo "✅ 效能測試通過"
```

### 安全性測試

```yaml
# templates/tests/test-security.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "my-app.fullname" . }}-test-security"
  annotations:
    "helm.sh/hook": test
    "helm.sh/hook-weight": "4"
spec:
  restartPolicy: Never
  containers:
  - name: test-security
    image: alpine:3.16
    command:
      - /bin/sh
      - -c
      - |
        apk add --no-cache curl
      
        echo "開始安全性測試..."
      
        BASE_URL="http://{{ include "my-app.fullname" . }}:{{ .Values.service.port }}"
      
        # 測試是否暴露敏感資訊
        echo "檢查敏感資訊暴露..."
        SENSITIVE_RESPONSE=$(curl -s $BASE_URL/debug || echo "not_found")
        if [ "$SENSITIVE_RESPONSE" = "not_found" ]; then
          echo "✅ 未暴露除錯端點"
        else
          echo "⚠️  警告：除錯端點可能暴露敏感資訊"
        fi
      
        # 檢查 HTTP 安全標頭
        echo "檢查安全標頭..."
        SECURITY_HEADERS=$(curl -s -I $BASE_URL/health | grep -i "x-frame-options\|x-content-type-options\|x-xss-protection")
        if [ -n "$SECURITY_HEADERS" ]; then
          echo "✅ 安全標頭檢查通過"
        else
          echo "⚠️  警告：缺少安全標頭"
        fi
      
        echo "安全性測試完成"
```

### 測試的自動化整合

```yaml
# .github/workflows/helm-test.yml
name: Helm Chart 測試

on:
  push:
    paths:
      - 'charts/**'
  pull_request:
    paths:
      - 'charts/**'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
  
    - name: 設定 Kubernetes
      uses: helm/kind-action@v1.4.0
      with:
        cluster_name: test-cluster
  
    - name: 安裝 Helm
      uses: azure/setup-helm@v3
      with:
        version: '3.10.0'
  
    - name: Lint Chart
      run: |
        helm lint charts/my-app
  
    - name: 測試模板渲染
      run: |
        helm template test-release charts/my-app
  
    - name: 安裝 Chart
      run: |
        helm install test-release charts/my-app \
          --wait --timeout=300s
  
    - name: 執行測試
      run: |
        helm test test-release
  
    - name: 清理
      if: always()
      run: |
        helm uninstall test-release
```

---

## Lab：建立具有依賴關係的 Chart

### 實作目標

建立一個完整的電商應用 Chart，包含：

- 前端 React 應用
- 後端 Node.js API
- PostgreSQL 資料庫
- Redis 快取
- 完整的 Hook 和測試

### 步驟 1：建立主 Chart 結構

```bash
# 建立主 Chart
helm create ecommerce-platform
cd ecommerce-platform

# 清理預設檔案
rm -rf templates/*
rm values.yaml
```

### 步驟 2：設定依賴關係

```yaml
# Chart.yaml
apiVersion: v2
name: ecommerce-platform
description: 完整的電商平台 Helm Chart
type: application
version: 0.1.0
appVersion: "1.0.0"

dependencies:
  # PostgreSQL 資料庫
  - name: postgresql
    version: "12.1.9"
    repository: "https://charts.bitnami.com/bitnami"
    condition: postgresql.enabled
    tags:
      - database
  
  # Redis 快取
  - name: redis
    version: "17.3.7"
    repository: "https://charts.bitnami.com/bitnami" 
    condition: redis.enabled
    tags:
      - cache

maintainers:
  - name: "Your Name"
    email: "your.email@company.com"

keywords:
  - ecommerce
  - nodejs
  - react
  - postgresql
  - redis
```

### 步驟 3：配置 values.yaml

```yaml
# values.yaml
global:
  imageRegistry: ""
  imagePullSecrets: []

# 應用配置
app:
  name: ecommerce-platform
  
# 前端配置
frontend:
  enabled: true
  replicaCount: 2
  image:
    repository: my-registry/ecommerce-frontend
    tag: "v1.0.0"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 80
  env:
    REACT_APP_API_URL: "/api"

# 後端 API 配置
backend:
  enabled: true
  replicaCount: 3
  image:
    repository: my-registry/ecommerce-backend
    tag: "v1.0.0"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 3000
  env:
    NODE_ENV: production
    PORT: "3000"

# PostgreSQL 配置
postgresql:
  enabled: true
  auth:
    postgresPassword: "secure-postgres-password"
    username: "ecommerce_user"
    password: "secure-user-password"
    database: "ecommerce_db"
  primary:
    persistence:
      enabled: true
      size: 20Gi
      storageClass: "fast-ssd"

# Redis 配置
redis:
  enabled: true
  auth:
    enabled: true
    password: "secure-redis-password"
  master:
    persistence:
      enabled: true
      size: 2Gi

# Ingress 配置
ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: ecommerce.example.com
      paths:
        - path: /
          pathType: Prefix
          service: frontend
        - path: /api
          pathType: Prefix
          service: backend
  tls:
    - secretName: ecommerce-tls
      hosts:
        - ecommerce.example.com

# Hook 配置
hooks:
  image:
    repository: alpine/curl
    tag: "latest"
```

### 步驟 4：實作 Hook

```yaml
# templates/hooks/pre-install-check.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ include "ecommerce-platform.fullname" . }}-pre-install-check"
  labels:
    {{- include "ecommerce-platform.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": pre-install
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": before-hook-creation,hook-succeeded
spec:
  template:
    metadata:
      name: "{{ include "ecommerce-platform.fullname" . }}-pre-install-check"
    spec:
      restartPolicy: Never
      containers:
      - name: pre-install-check
        image: "{{ .Values.hooks.image.repository }}:{{ .Values.hooks.image.tag }}"
        command:
          - /bin/sh
          - -c
          - |
            echo "🔍 執行安裝前檢查..."
          
            # 檢查必要的 StorageClass
            {{- if .Values.postgresql.primary.persistence.enabled }}
            echo "檢查 StorageClass: {{ .Values.postgresql.primary.persistence.storageClass }}"
            # 這裡可以加入實際的檢查邏輯
            {{- end }}
          
            # 檢查 Ingress Controller
            {{- if .Values.ingress.enabled }}
            echo "檢查 Ingress Controller: {{ .Values.ingress.className }}"
            {{- end }}
          
            echo "✅ 安裝前檢查完成"
```

```yaml
# templates/hooks/post-install-migration.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ include "ecommerce-platform.fullname" . }}-db-migration"
  labels:
    {{- include "ecommerce-platform.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": post-install,post-upgrade
    "helm.sh/hook-weight": "1"
    "helm.sh/hook-delete-policy": before-hook-creation
spec:
  template:
    metadata:
      name: "{{ include "ecommerce-platform.fullname" . }}-db-migration"
    spec:
      restartPolicy: Never
      containers:
      - name: db-migration
        image: "{{ .Values.backend.image.repository }}:{{ .Values.backend.image.tag }}"
        command:
          - /bin/sh
          - -c
          - |
            echo "🚀 開始資料庫遷移..."
            
            # 等待資料庫準備就緒
            echo "等待 PostgreSQL 準備就緒..."
            until pg_isready -h {{ include "postgresql.primary.fullname" .Subcharts.postgresql }} -p 5432; do
              echo "等待資料庫..."
              sleep 2
            done
            
            # 執行資料庫遷移
            echo "執行資料庫遷移腳本..."
            npm run migrate
            
            # 檢查遷移結果
            if [ $? -eq 0 ]; then
              echo "✅ 資料庫遷移成功完成"
            else
              echo "❌ 資料庫遷移失敗"
              exit 1
            fi
            
            # 初始化基礎資料
            echo "初始化基礎資料..."
            npm run seed:basic
            
            echo "🎉 資料庫設定完成"
        env:
        - name: DATABASE_URL
          value: "postgresql://{{ .Values.postgresql.auth.username }}:{{ .Values.postgresql.auth.password }}@{{ include "postgresql.primary.fullname" .Subcharts.postgresql }}:5432/{{ .Values.postgresql.auth.database }}"
        - name: NODE_ENV
          value: "migration"

```

步驟 5：建立應用模板
```yaml
# templates/backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: "{{ include "ecommerce-platform.fullname" . }}-backend"
  labels:
    {{- include "ecommerce-platform.labels" . | nindent 4 }}
    app.kubernetes.io/component: backend
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.backend.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "ecommerce-platform.selectorLabels" . | nindent 6 }}
      app.kubernetes.io/component: backend
  template:
    metadata:
      labels:
        {{- include "ecommerce-platform.selectorLabels" . | nindent 8 }}
        app.kubernetes.io/component: backend
    spec:
      containers:
      - name: backend
        image: "{{ .Values.backend.image.repository }}:{{ .Values.backend.image.tag }}"
        imagePullPolicy: {{ .Values.backend.image.pullPolicy }}
        ports:
        - name: http
          containerPort: {{ .Values.backend.service.port }}
          protocol: TCP
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: http
          initialDelaySeconds: 5
          periodSeconds: 5
        env:
        {{- range $key, $value := .Values.backend.env }}
        - name: {{ $key }}
          value: {{ $value | quote }}
        {{- end }}
        - name: DATABASE_URL
          value: "postgresql://{{ .Values.postgresql.auth.username }}:{{ .Values.postgresql.auth.password }}@{{ include "postgresql.primary.fullname" .Subcharts.postgresql }}:5432/{{ .Values.postgresql.auth.database }}"
        - name: REDIS_URL
          value: "redis://:{{ .Values.redis.auth.password }}@{{ include "redis.fullname" .Subcharts.redis }}-master:6379"
        resources:
          {{- toYaml .Values.backend.resources | nindent 12 }}
```

```yaml
# templates/frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: "{{ include "ecommerce-platform.fullname" . }}-frontend"
  labels:
    {{- include "ecommerce-platform.labels" . | nindent 4 }}
    app.kubernetes.io/component: frontend
spec:
  replicas: {{ .Values.frontend.replicaCount }}
  selector:
    matchLabels:
      {{- include "ecommerce-platform.selectorLabels" . | nindent 6 }}
      app.kubernetes.io/component: frontend
  template:
    metadata:
      labels:
        {{- include "ecommerce-platform.selectorLabels" . | nindent 8 }}
        app.kubernetes.io/component: frontend
    spec:
      containers:
      - name: frontend
        image: "{{ .Values.frontend.image.repository }}:{{ .Values.frontend.image.tag }}"
        imagePullPolicy: {{ .Values.frontend.image.pullPolicy }}
        ports:
        - name: http
          containerPort: 80
          protocol: TCP
        livenessProbe:
          httpGet:
            path: /
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: http
          initialDelaySeconds: 5
          periodSeconds: 5
        env:
        {{- range $key, $value := .Values.frontend.env }}
        - name: {{ $key }}
          value: {{ $value | quote }}
        {{- end }}
        resources:
          {{- toYaml .Values.frontend.resources | nindent 12 }}
```

步驟 6：建立完整的測試套件

```yaml
# templates/tests/test-connectivity.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "ecommerce-platform.fullname" . }}-test-connectivity"
  labels:
    {{- include "ecommerce-platform.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test
    "helm.sh/hook-weight": "1"
    "helm.sh/hook-delete-policy": before-hook-creation,hook-succeeded
spec:
  restartPolicy: Never
  containers:
  - name: test-connectivity
    image: busybox:1.35
    command:
      - /bin/sh
      - -c
      - |
        echo "🔍 開始連線測試..."
        
        # 測試前端服務
        echo "測試前端服務連線..."
        if nc -z {{ include "ecommerce-platform.fullname" . }}-frontend {{ .Values.frontend.service.port }}; then
          echo "✅ 前端服務連線成功"
        else
          echo "❌ 前端服務連線失敗"
          exit 1
        fi
        
        # 測試後端 API 服務
        echo "測試後端 API 服務連線..."
        if nc -z {{ include "ecommerce-platform.fullname" . }}-backend {{ .Values.backend.service.port }}; then
          echo "✅ 後端 API 服務連線成功"
        else
          echo "❌ 後端 API 服務連線失敗"
          exit 1
        fi
        
        {{- if .Values.postgresql.enabled }}
        # 測試 PostgreSQL 連線
        echo "測試 PostgreSQL 資料庫連線..."
        if nc -z {{ include "postgresql.primary.fullname" .Subcharts.postgresql }} 5432; then
          echo "✅ PostgreSQL 資料庫連線成功"
        else
          echo "❌ PostgreSQL 資料庫連線失敗"
          exit 1
        fi
        {{- end }}
        
        {{- if .Values.redis.enabled }}
        # 測試 Redis 連線
        echo "測試 Redis 快取連線..."
        if nc -z {{ include "redis.fullname" .Subcharts.redis }}-master 6379; then
          echo "✅ Redis 快取連線成功"
        else
          echo "❌ Redis 快取連線失敗"
          exit 1
        fi
        {{- end }}
        
        echo "🎉 所有連線測試通過！"
```

```yaml
# templates/tests/test-api-functionality.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "ecommerce-platform.fullname" . }}-test-api"
  labels:
    {{- include "ecommerce-platform.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test
    "helm.sh/hook-weight": "2"
    "helm.sh/hook-delete-policy": before-hook-creation,hook-succeeded
spec:
  restartPolicy: Never
  containers:
  - name: test-api
    image: curlimages/curl:7.85.0
    command:
      - /bin/sh
      - -c
      - |
        echo "🧪 開始 API 功能測試..."
        
        BACKEND_URL="http://{{ include "ecommerce-platform.fullname" . }}-backend:{{ .Values.backend.service.port }}"
        
        # 測試健康檢查端點
        echo "測試健康檢查端點..."
        HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" $BACKEND_URL/health)
        if [ "$HEALTH_STATUS" = "200" ]; then
          echo "✅ 健康檢查端點正常"
        else
          echo "❌ 健康檢查端點異常，狀態碼：$HEALTH_STATUS"
          exit 1
        fi
        
        # 測試 API 版本端點
        echo "測試 API 版本端點..."
        VERSION_RESPONSE=$(curl -s $BACKEND_URL/api/v1/version)
        if echo "$VERSION_RESPONSE" | grep -q "version"; then
          echo "✅ API 版本端點正常"
          echo "API 版本資訊：$VERSION_RESPONSE"
        else
          echo "❌ API 版本端點異常"
          exit 1
        fi
        
        # 測試資料庫連線狀態
        echo "測試資料庫連線狀態..."
        DB_STATUS=$(curl -s $BACKEND_URL/api/v1/db-status)
        if echo "$DB_STATUS" | grep -q '"connected":true'; then
          echo "✅ 資料庫連線狀態正常"
        else
          echo "❌ 資料庫連線狀態異常"
          echo "回應內容：$DB_STATUS"
          exit 1
        fi
        
        # 測試快取連線狀態
        echo "測試快取連線狀態..."
        CACHE_STATUS=$(curl -s $BACKEND_URL/api/v1/cache-status)
        if echo "$CACHE_STATUS" | grep -q '"connected":true'; then
          echo "✅ 快取連線狀態正常"
        else
          echo "❌ 快取連線狀態異常"
          echo "回應內容：$CACHE_STATUS"
          exit 1
        fi
        
        echo "🎉 所有 API 功能測試通過！"
```

```yaml
# templates/tests/test-end-to-end.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "ecommerce-platform.fullname" . }}-test-e2e"
  labels:
    {{- include "ecommerce-platform.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test
    "helm.sh/hook-weight": "3"
    "helm.sh/hook-delete-policy": before-hook-creation,hook-succeeded
spec:
  restartPolicy: Never
  containers:
  - name: test-e2e
    image: curlimages/curl:7.85.0
    command:
      - /bin/sh
      - -c
      - |
        echo "🔄 開始端對端測試..."
        
        BACKEND_URL="http://{{ include "ecommerce-platform.fullname" . }}-backend:{{ .Values.backend.service.port }}"
        
        # 模擬使用者註冊流程
        echo "測試使用者註冊..."
        REGISTER_RESPONSE=$(curl -s -X POST $BACKEND_URL/api/v1/auth/register \
          -H "Content-Type: application/json" \
          -d '{"email":"test@example.com","password":"testpass123","name":"Test User"}')
        
        if echo "$REGISTER_RESPONSE" | grep -q "success\|token\|user"; then
          echo "✅ 使用者註冊測試通過"
        else
          echo "⚠️  使用者註冊測試跳過（可能已存在）"
        fi
        
        # 測試使用者登入
        echo "測試使用者登入..."
        LOGIN_RESPONSE=$(curl -s -X POST $BACKEND_URL/api/v1/auth/login \
          -H "Content-Type: application/json" \
          -d '{"email":"test@example.com","password":"testpass123"}')
        
        if echo "$LOGIN_RESPONSE" | grep -q "token"; then
          echo "✅ 使用者登入測試通過"
          TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        else
          echo "❌ 使用者登入測試失敗"
          echo "回應內容：$LOGIN_RESPONSE"
          exit 1
        fi
        
        # 測試受保護的 API 端點
        echo "測試受保護的 API 端點..."
        PROFILE_RESPONSE=$(curl -s $BACKEND_URL/api/v1/user/profile \
          -H "Authorization: Bearer $TOKEN")
        
        if echo "$PROFILE_RESPONSE" | grep -q "email\|name"; then
          echo "✅ 受保護 API 端點測試通過"
        else
          echo "❌ 受保護 API 端點測試失敗"
          exit 1
        fi
        
        # 測試商品列表 API
        echo "測試商品列表 API..."
        PRODUCTS_RESPONSE=$(curl -s $BACKEND_URL/api/v1/products)
        
        if echo "$PRODUCTS_RESPONSE" | grep -q "products\|\[\]"; then
          echo "✅ 商品列表 API 測試通過"
        else
          echo "❌ 商品列表 API 測試失敗"
          exit 1
        fi
        
        echo "🎉 端對端測試完成！"
```

步驟 7：建立部署腳本
```yaml
#!/bin/bash
# scripts/deploy.sh

set -e

CHART_NAME="ecommerce-platform"
NAMESPACE="ecommerce"
RELEASE_NAME="my-ecommerce"

echo "🚀 開始部署 $CHART_NAME..."

# 建立 namespace（如果不存在）
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# 更新依賴
echo "📦 更新 Chart 依賴..."
helm dependency update

# 驗證 Chart
echo "🔍 驗證 Chart..."
helm lint .

# 測試模板渲染
echo "🧪 測試模板渲染..."
helm template $RELEASE_NAME . --namespace $NAMESPACE > /tmp/rendered-templates.yaml
echo "渲染的模板已儲存到 /tmp/rendered-templates.yaml"

# 部署或升級
if helm list -n $NAMESPACE | grep -q $RELEASE_NAME; then
    echo "⬆️  升級現有部署..."
    helm upgrade $RELEASE_NAME . \
        --namespace $NAMESPACE \
        --wait \
        --timeout=600s \
        --history-max=5
else
    echo "🆕 執行新部署..."
    helm install $RELEASE_NAME . \
        --namespace $NAMESPACE \
        --wait \
        --timeout=600s \
        --create-namespace
fi

echo "✅ 部署完成！"

# 執行測試
echo "🧪 執行 Helm 測試..."
helm test $RELEASE_NAME --namespace $NAMESPACE

# 顯示部署狀態
echo "📊 部署狀態："
helm status $RELEASE_NAME --namespace $NAMESPACE

echo "🎉 所有操作完成！"
```

步驟 8：建立不同環境的配置
```yaml
# values-dev.yaml
postgresql:
  primary:
    persistence:
      size: 5Gi  # 開發環境較小的儲存空間

redis:
  master:
    persistence:
      size: 1Gi

backend:
  replicaCount: 1  # 開發環境單一副本
  env:
    NODE_ENV: development
    LOG_LEVEL: debug

frontend:
  replicaCount: 1
  env:
    REACT_APP_DEBUG: "true"

ingress:
  hosts:
    - host: ecommerce-dev.example.com
      paths:
        - path: /
          pathType: Prefix
          service: frontend
        - path: /api
          pathType: Prefix
          service: backend
```

```yaml
# values-prod.yaml
postgresql:
  primary:
    persistence:
      size: 100Gi
      storageClass: "premium-ssd"
    resources:
      requests:
        memory: 2Gi
        cpu: 1000m
      limits:
        memory: 4Gi
        cpu: 2000m

redis:
  master:
    persistence:
      size: 10Gi
    resources:
      requests:
        memory: 512Mi
        cpu: 250m

backend:
  replicaCount: 5
  env:
    NODE_ENV: production
    LOG_LEVEL: warn
  resources:
    requests:
      memory: 512Mi
      cpu: 250m
    limits:
      memory: 1Gi
      cpu: 500m

frontend:
  replicaCount: 3
  resources:
    requests:
      memory: 128Mi
      cpu: 100m
    limits:
      memory: 256Mi
      cpu: 200m

# 生產環境啟用自動擴展
autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80
```

步驟 9：測試與驗證
```bash
# 測試不同環境的部署
echo "測試開發環境配置..."
helm template test-dev . -f values-dev.yaml > /tmp/dev-templates.yaml

echo "測試生產環境配置..."
helm template test-prod . -f values-prod.yaml > /tmp/prod-templates.yaml

# 比較不同環境的差異
echo "比較環境差異..."
diff /tmp/dev-templates.yaml /tmp/prod-templates.yaml || true

# 執行完整的測試流程
echo "執行完整測試..."
./scripts/deploy.sh
```

# 
總結與最佳實踐
關鍵學習重點
依賴管理策略

基礎設施（資料庫、快取）使用外部依賴
應用元件使用子 Chart
條件式依賴讓不同環境有彈性
Hook 機制的威力

pre-install：環境檢查和準備
post-install：資料庫遷移和初始化
pre-upgrade：備份和準備
pre-delete：清理外部資源
完整的測試策略

連線測試：確保服務間通訊正常
功能測試：驗證 API 端點正常運作
端對端測試：模擬真實使用者流程
效能測試：確保系統效能符合要求
企業級部署考量

環境特定配置分離
自動化部署腳本
完整的錯誤處理
安全性檢查


# 實際應用建議
開發階段：

使用 helm template 驗證模板正確性
頻繁執行 helm lint 檢查語法
建立完整的測試套件
測試階段：

在不同環境測試不同的 values 檔案
驗證 Hook 的執行順序和錯誤處理
測試升級和回滾流程
生產階段：

使用 GitOps 工作流程
實施完整的監控和告警
定期備份和災難恢復測試
進階主題預告
通過今天的深入學習，你已經掌握了：

複雜應用的依賴管理
Hook 機制的實際應用
企業級的測試策略
這些技能讓你能夠開發和維護生產級別的 Helm Chart。在實際工作中，這些知識將幫助你：

管理複雜的微服務架構
實現可靠的部署流程
確保應用的高可用性和穩定性
