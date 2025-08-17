
# Kubernetes 學習路徑 (使用 KinD)

非常好的整理！基於您的重點，我將幫您調整學習路徑，著重在 KinD、DevSpace、Helm 以及搭建前後端叢集、PostgreSQL，同時保留後半段可觀測性與混沌工程的內容。這樣的計劃更加聚焦於實用技能。

## 第一週：Kubernetes 基礎與 KinD

### Day 1: KinD 入門與本地開發環境設置

 **學習目的** ：掌握 KinD 的安裝、配置和基本操作，建立高效的本地開發環境。

 **實作內容** ：

1. 安裝 KinD、kubectl 和其他必要工具
2. 創建多節點 KinD 集群
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">kind: Cluster
   apiVersion: kind.x-k8s.io/v1alpha4
   nodes:
   - role: control-plane
     extraPortMappings:
     - containerPort: 80
       hostPort: 80
     - containerPort: 443
       hostPort: 443
   - role: worker
   - role: worker
   </code></div></div></pre>
3. 熟悉 KinD 特有命令和選項
4. 了解 KinD 的局限性與真實集群差異
5. 設置開發工具與 KinD 的集成 (VS Code, kubectl 插件等)

### Day 2: Kubernetes 組件架構與 YAML 基礎

 **學習目的** ：理解 Kubernetes 架構和核心組件，掌握 YAML 語法基礎。

 **實作內容** ：

1. 探索 KinD 集群中的系統組件 (`kubectl get pods -A`)
2. 解析 YAML 語法和結構
3. 理解 K8s 資源的共同模式 (apiVersion, kind, metadata, spec)
4. 練習 Docker Compose 到 Kubernetes YAML 的轉換
5. 使用 `kubectl explain` 學習資源定義

### Day 3: Pod YAML 詳解與實作

 **學習目的** ：掌握 Pod 配置和生命週期管理。

 **實作內容** ：

1. 創建和管理基本 Pod
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: v1
   kind: Pod
   metadata:
     name: nginx-pod
     labels:
       app: nginx
   spec:
     containers:
     - name: nginx
       image: nginx:alpine
       ports:
       - containerPort: 80
       resources:
         requests:
           memory: "64Mi"
           cpu: "100m"
         limits:
           memory: "128Mi"
           cpu: "200m"
   </code></div></div></pre>
2. 設定容器環境變數、資源限制和存活探針
3. 實驗多容器 Pod 設計
4. 使用 `kubectl port-forward` 訪問 Pod
5. 練習各種故障診斷命令

### Day 4: Deployment YAML 實作與管理

 **學習目的** ：理解 Deployment 管理 Pod 的機制，掌握應用更新策略。

 **實作內容** ：

1. 創建管理多副本的 Deployment
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: web-app
   spec:
     replicas: 3
     selector:
       matchLabels:
         app: web-app
     strategy:
       type: RollingUpdate
       rollingUpdate:
         maxSurge: 1
         maxUnavailable: 0
     template:
       metadata:
         labels:
           app: web-app
       spec:
         containers:
         - name: web-app
           image: nginx:1.19
   </code></div></div></pre>
2. 實施滾動更新和版本回滾
3. 比較不同更新策略 (RollingUpdate vs Recreate)
4. 使用標籤選擇器管理 Pod 集
5. 監控 Deployment 狀態和事件

### Day 5: Service YAML 與網路模型

 **學習目的** ：理解 Kubernetes 服務發現和內部網路架構。

 **實作內容** ：

1. 創建 ClusterIP Service 連接到 Deployment
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: v1
   kind: Service
   metadata:
     name: web-service
   spec:
     selector:
       app: web-app
     ports:
     - port: 80
       targetPort: 80
     type: ClusterIP
   </code></div></div></pre>
2. 設置 NodePort Service 訪問應用
3. 測試 KinD 中的 LoadBalancer 服務 (使用 metallb)
4. 理解服務選擇器和標籤關係
5. 實驗服務的會話黏性和負載均衡行為

### Day 6: ConfigMap 與 Secret 配置管理

 **學習目的** ：掌握應用配置和敏感信息管理的最佳實踐。

 **實作內容** ：

1. 創建和使用不同類型的 ConfigMap
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: v1
   kind: ConfigMap
   metadata:
     name: app-config
   data:
     app.properties: |
       environment=development
       log.level=debug
       feature.x=true
     config.json: |
       {
         "apiUrl": "http://api-service",
         "maxConnections": 100
       }
   </code></div></div></pre>
2. 將 ConfigMap 掛載為文件和環境變數
3. 創建並使用 Secret 存儲敏感數據
4. 配置 Pod 使用 ConfigMap 和 Secret
5. 測試配置更新如何影響應用行為

### Day 7: PersistentVolume 與 PVC 存儲管理

 **學習目的** ：學習 Kubernetes 持久化存儲機制。

 **實作內容** ：

1. 在 KinD 中配置本地存儲
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: v1
   kind: PersistentVolume
   metadata:
     name: local-pv
   spec:
     capacity:
       storage: 1Gi
     accessModes:
     - ReadWriteOnce
     persistentVolumeReclaimPolicy: Retain
     storageClassName: standard
     hostPath:
       path: /tmp/data
   </code></div></div></pre>
2. 創建 PVC 並連接到 Pod
3. 配置動態存儲供應 (使用 local-path-provisioner)
4. 測試數據持久性 (刪除 Pod 後驗證數據保留)
5. 實驗不同的訪問模式和存儲類

## 第二週：DevSpace 與高級工作負載

### Day 8: DevSpace 入門與開發工作流

 **學習目的** ：使用 DevSpace 建立高效的 Kubernetes 開發工作流。

 **實作內容** ：

1. 安裝和配置 DevSpace
2. 創建基本 DevSpace 配置文件
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">version: v2beta1
   name: my-project

   vars:
     IMAGE: my-registry.com/my-image

   pipelines:
     dev:
       run: |-
         run_dependencies --all
         create_deployments --all
         start_dev --all

   deployments:
     app:
       helm:
         chart:
           name: ./charts/app
         values:
           image: ${IMAGE}
   </code></div></div></pre>
3. 實現代碼與容器的熱重載
4. 使用 DevSpace 進行本地調試
5. 整合 DevSpace 與現有開發工具

### Day 9: StatefulSet 與有狀態應用部署

 **學習目的** ：理解 StatefulSet 的特性，學習部署有狀態應用。

 **實作內容** ：

1. 部署 PostgreSQL StatefulSet
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: apps/v1
   kind: StatefulSet
   metadata:
     name: postgres
   spec:
     serviceName: postgres
     replicas: 1
     selector:
       matchLabels:
         app: postgres
     template:
       metadata:
         labels:
           app: postgres
       spec:
         containers:
         - name: postgres
           image: postgres:13
           env:
           - name: POSTGRES_PASSWORD
             valueFrom:
               secretKeyRef:
                 name: postgres-secret
                 key: password
           ports:
           - containerPort: 5432
           volumeMounts:
           - name: data
             mountPath: /var/lib/postgresql/data
     volumeClaimTemplates:
     - metadata:
         name: data
       spec:
         accessModes: ["ReadWriteOnce"]
         resources:
           requests:
             storage: 1Gi
   </code></div></div></pre>
2. 設置 PostgreSQL 的持久化存儲
3. 配置 Headless Service 提供穩定網絡標識
4. 實施基本的數據庫備份和恢復策略
5. 測試 Pod 重建後的數據持久性

### Day 10: Job 與 CronJob 實作數據庫遷移

 **學習目的** ：掌握批處理任務執行和數據庫遷移自動化。

 **實作內容** ：

1. 創建 Flyway 數據庫遷移 Job
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: batch/v1
   kind: Job
   metadata:
     name: db-migration
   spec:
     template:
       spec:
         containers:
         - name: flyway
           image: flyway/flyway:7
           args:
           - -url=jdbc:postgresql://postgres:5432/mydb
           - -user=postgres
           - -password=$(POSTGRES_PASSWORD)
           - -connectRetries=10
           - migrate
           env:
           - name: POSTGRES_PASSWORD
             valueFrom:
               secretKeyRef:
                 name: postgres-secret
                 key: password
           volumeMounts:
           - name: migrations
             mountPath: /flyway/sql
         volumes:
         - name: migrations
           configMap:
             name: db-migrations
         restartPolicy: OnFailure
   </code></div></div></pre>
2. 使用 ConfigMap 存儲 SQL 遷移腳本
3. 配置 CronJob 實施定期數據處理
4. 實現 Job 的重試和並行執行策略
5. 設計 Job 依賴和順序執行機制

### Day 11: Ingress 部署與路由配置

 **學習目的** ：學習在 KinD 環境中設置 Ingress 進行流量路由。

 **實作內容** ：

1. 在 KinD 中部署 Nginx Ingress Controller
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
   </code></div></div></pre>
2. 配置基本 Ingress 路由
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: networking.k8s.io/v1
   kind: Ingress
   metadata:
     name: app-ingress
     annotations:
       nginx.ingress.kubernetes.io/rewrite-target: /
   spec:
     rules:
     - host: app.local
       http:
         paths:
         - path: /
           pathType: Prefix
           backend:
             service:
               name: web-service
               port:
                 number: 80
         - path: /api
           pathType: Prefix
           backend:
             service:
               name: api-service
               port:
                 number: 8080
   </code></div></div></pre>
3. 實現路徑重寫和重定向
4. 配置 TLS 終止 (自簽名證書)
5. 測試多服務路由和負載均衡

### Day 12: HorizontalPodAutoscaler 自動擴展

 **學習目的** ：設置基於負載的自動擴展機制。

 **實作內容** ：

1. 部署 Metrics Server 到 KinD 集群
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
   </code></div></div></pre>
2. 配置基於 CPU 的 HPA
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: autoscaling/v2
   kind: HorizontalPodAutoscaler
   metadata:
     name: api-hpa
   spec:
     scaleTargetRef:
       apiVersion: apps/v1
       kind: Deployment
       name: api-service
     minReplicas: 2
     maxReplicas: 10
     metrics:
     - type: Resource
       resource:
         name: cpu
         target:
           type: Utilization
           averageUtilization: 50
   </code></div></div></pre>
3. 創建負載測試工具生成流量
4. 監控自動擴展行為和指標
5. 調整 HPA 參數優化擴展性能

### Day 13: NetworkPolicy 網絡安全策略

 **學習目的** ：理解和實施 Kubernetes 網絡隔離。

 **實作內容** ：

1. 確認 KinD 使用支持 NetworkPolicy 的 CNI
2. 創建默認拒絕策略
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: default-deny
     namespace: production
   spec:
     podSelector: {}
     policyTypes:
     - Ingress
     - Egress
   </code></div></div></pre>
3. 實施微服務間精確通信控制
4. 配置命名空間隔離策略
5. 使用網絡工具測試和驗證隔離效果

### Day 14: Helm Chart 入門與應用打包

 **學習目的** ：使用 Helm 簡化應用部署和管理。

 **實作內容** ：

1. 安裝 Helm 並了解基本概念
2. 使用現有 Chart 部署應用
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">helm repo add bitnami https://charts.bitnami.com/bitnami
   helm install my-release bitnami/nginx
   </code></div></div></pre>
3. 創建自定義 Helm Chart 結構
4. 編寫模板和使用內置函數
5. 設計不同環境的值文件 (values.yaml)

## 第三週：應用架構與多環境配置

### Day 15: Helm 高級模板與依賴管理

 **學習目的** ：掌握 Helm 高級功能，建立複雜的應用包。

 **實作內容** ：

1. 編寫高級 Helm 模板 (條件、循環、函數)
2. 管理 Chart 依賴關係
3. 創建 Umbrella Chart 組織多服務應用
4. 設計可重用的 Helm 庫 (Library Charts)
5. 使用 Helm Hook 管理部署生命週期

### Day 16: ConfigMap 與環境變數注入深入解析

 **學習目的** ：深入理解配置管理和注入策略。

 **實作內容** ：

1. 比較不同環境變數注入方法
2. 實現 SQL 遷移腳本的 ConfigMap 存儲
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: v1
   kind: ConfigMap
   metadata:
     name: db-migrations
   data:
     V1__Create_tables.sql: |
       CREATE TABLE users (
         id SERIAL PRIMARY KEY,
         username VARCHAR(50) NOT NULL UNIQUE,
         email VARCHAR(100) NOT NULL
       );
     V2__Add_timestamps.sql: |
       ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT NOW();
       ALTER TABLE users ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();
   </code></div></div></pre>
3. 設計多層次配置管理策略
4. 實現配置熱更新和應用重載
5. 為不同環境配置不同 ConfigMap

### Day 17: 前後端分離架構設計

 **學習目的** ：設計和實現前後端分離的微服務架構。

 **實作內容** ：

1. 設計前後端分離應用架構
2. 部署 React 前端應用
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: frontend
   spec:
     replicas: 2
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
           image: my-frontend:latest
           ports:
           - containerPort: 80
   </code></div></div></pre>
3. 部署 Node.js/Java/Go 後端 API 服務
4. 配置服務間通信和 API 訪問
5. 實現前端與後端的整合測試

### Day 18: 數據庫高可用配置與連接管理

 **學習目的** ：深入理解數據庫部署和連接池配置。

 **實作內容** ：

1. 設計 PostgreSQL 高可用架構
2. 配置連接池和資源限制
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: v1
   kind: ConfigMap
   metadata:
     name: db-config
   data:
     db-config.properties: |
       # Connection Pool Settings
       maximumPoolSize=20
       minimumIdle=5
       connectionTimeout=30000
       idleTimeout=600000
       maxLifetime=1800000
   </code></div></div></pre>
3. 實現數據庫健康檢查和監控
4. 設計數據庫備份和恢復流程
5. 優化數據庫性能和資源使用

### Day 19: 從 Git 倉庫自動載入配置

 **學習目的** ：實現配置與代碼分離，從外部源動態加載配置。

 **實作內容** ：

1. 設置配置 Git 倉庫結構
2. 實現 InitContainer 從 Git 拉取配置
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: v1
   kind: Pod
   metadata:
     name: app-with-config
   spec:
     initContainers:
     - name: git-clone
       image: alpine/git
       command:
       - git
       - clone
       - --single-branch
       - --branch=main
       - https://github.com/user/config-repo.git
       - /config
       volumeMounts:
       - name: config-volume
         mountPath: /config
     containers:
     - name: app
       image: my-app:latest
       volumeMounts:
       - name: config-volume
         mountPath: /app/config
     volumes:
     - name: config-volume
       emptyDir: {}
   </code></div></div></pre>
3. 實現數據庫遷移腳本從 Git 加載
4. 設計配置變更偵測機制
5. 建立配置版本控制與回滾策略

### Day 20: DevSpace 高級開發工作流

 **學習目的** ：建立高效的雲原生開發工作流。

 **實作內容** ：

1. 設置完整的 DevSpace 開發環境
2. 配置開發、測試和生產環境設定檔
3. 實現代碼變更熱重載
4. 設置遠程調試和日誌流
5. 優化構建流程和鏡像管理

### Day 21: 完整微服務架構部署

 **學習目的** ：整合所學知識，部署完整的微服務系統。

 **實作內容** ：

1. 設計完整的微服務拓撲
2. 使用 Helm 部署整個應用堆疊
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">helm install myapp ./myapp-chart \
     --set frontend.replicas=2 \
     --set api.replicas=3 \
     --set database.storage=2Gi
   </code></div></div></pre>
3. 配置服務間依賴和啟動順序
4. 實現端到端的功能測試
5. 文檔化部署架構和流程

## 第四週：可觀測性與混沌工程

### Day 22: 多容器 Pod 設計模式實踐

 **學習目的** ：掌握多容器 Pod 的高級設計模式。

 **實作內容** ：

1. 實現日誌收集 Sidecar
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: v1
   kind: Pod
   metadata:
     name: app-with-sidecar
   spec:
     containers:
     - name: app
       image: my-app:latest
       volumeMounts:
       - name: logs
         mountPath: /app/logs
     - name: log-collector
       image: fluent/fluent-bit
       volumeMounts:
       - name: logs
         mountPath: /logs
       - name: config
         mountPath: /fluent-bit/etc/
     volumes:
     - name: logs
       emptyDir: {}
     - name: config
       configMap:
         name: fluent-bit-config
   </code></div></div></pre>
2. 部署 Ambassador 代理容器
3. 創建 Adapter 適配器容器
4. 設計容器間共享和通信
5. 實施資源管理和優先級控制

### Day 23: Prometheus 與基本監控

 **學習目的** ：設置基本的應用監控和指標收集。

 **實作內容** ：

1. 在 KinD 中部署 Prometheus
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
   helm install prometheus prometheus-community/prometheus
   </code></div></div></pre>
2. 配置服務發現和目標選擇
3. 為應用添加 Prometheus 指標導出
4. 設計基本的監控儀表板
5. 配置簡單的告警規則

### Day 24: Grafana 與可視化

 **學習目的** ：建立強大的可視化和儀表板。

 **實作內容** ：

1. 部署 Grafana 並連接到 Prometheus
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">helm install grafana grafana/grafana
   </code></div></div></pre>
2. 創建應用性能儀表板
3. 設計數據庫監控面板
4. 配置用戶和團隊權限
5. 設置告警通知渠道

### Day 25: Loki 與日誌管理

 **學習目的** ：實現集中式日誌收集和分析。

 **實作內容** ：

1. 部署 Loki 日誌聚合系統
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">helm install loki grafana/loki-stack \
     --set grafana.enabled=false
   </code></div></div></pre>
2. 配置 Promtail 收集容器日誌
3. 實現日誌過濾和標籤策略
4. 創建日誌查詢和儀表板
5. 設置日誌告警和異常檢測

### Day 26: Tempo 與分布式追蹤

 **學習目的** ：實現分布式系統的追蹤和性能分析。

 **實作內容** ：

1. 部署 Tempo 追蹤系統
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">helm install tempo grafana/tempo
   </code></div></div></pre>
2. 為應用添加 OpenTelemetry 檢測
3. 實現服務間追蹤上下文傳播
4. 創建追蹤視圖和服務圖
5. 使用追蹤分析性能瓶頸

### Day 27: OpenTelemetry Collector 部署

 **學習目的** ：建立統一的遙測數據收集管道。

 **實作內容** ：

1. 部署 OpenTelemetry Operator
   <pre><div class="code-enhance--_fiUF hljs language-bash"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-bash" data-code-tools="">helm install opentelemetry-operator open-telemetry/opentelemetry-operator
   </code></div></div></pre>
2. 配置 Collector 收集管道
   <pre><div class="code-enhance--_fiUF hljs language-yaml"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-yaml" data-code-tools="">apiVersion: opentelemetry.io/v1alpha1
   kind: OpenTelemetryCollector
   metadata:
     name: otel-collector
   spec:
     config: |
       receivers:
         otlp:
           protocols:
             grpc:
             http:
       processors:
         batch:
       exporters:
         prometheus:
           endpoint: "0.0.0.0:8889"
         loki:
           endpoint: "http://loki:3100/loki/api/v1/push"
         otlp:
           endpoint: "tempo:4317"
           tls:
             insecure: true
       service:
         pipelines:
           metrics:
             receivers: [otlp]
             processors: [batch]
             exporters: [prometheus]
           logs:
             receivers: [otlp]
             processors: [batch]
             exporters: [loki]
           traces:
             receivers: [otlp]
             processors: [batch]
             exporters: [otlp]
   </code></div></div></pre>
3. 實現自動檢測注入
4. 配置多源數據收集
5. 優化數據處理和導出

### Day 28: Grafana 堆疊整合

 **學習目的** ：整合所有可觀測性系統，建立統一觀測平台。

 **實作內容** ：

1. 整合 Prometheus、Loki 和 Tempo
2. 創建跨數據源儀表板
3. 設置數據源關聯 (exemplars)
4. 實施統一的告警系統
5. 優化數據存儲和查詢性能

### Day 29: 應用檢測與可觀測性最佳實踐

 **學習目的** ：將應用與可觀測性堆疊深度集成。

 **實作內容** ：

1. 為前端應用添加 RUM (Real User Monitoring)
2. 為後端服務實施 OpenTelemetry 檢測
3. 實現自定義業務指標收集
4. 配置 RED (Rate, Error, Duration) 監控
5. 設計服務水平目標 (SLO) 監控

### Day 30: xk6-disruptor 混沌測試入門

 **學習目的** ：理解混沌工程原則，學習基本的故障注入。

 **實作內容** ：

1. 安裝 k6 和 xk6-disruptor
2. 創建基本的負載測試腳本
3. 設計簡單的故障注入實驗
4. 配置 Pod 和網絡故障場景
5. 監控故障影響和系統行為

### Day 31: 高級混沌測試與韌性驗證

 **學習目的** ：實施複雜的混沌實驗，驗證系統韌性。

 **實作內容** ：

1. 設計混沌實驗計劃和假設
2. 實施複合故障場景
   <pre><div class="code-enhance--_fiUF hljs language-javascript"><div class="code-enhance-header--NVMaE"><div class="code-enhance-header-right--R62yQ"><span class="code-enhance-copy--XipCY"><span>Copy</span></span></div></div><div class="code-enhance-content--fGI3Q"><code class="hljs language-javascript" data-code-tools="">import { check } from 'k6';
   import http from 'k6/http';
   import { PodDisruptor } from 'k6/x/disruptor';

   export default function () {
     const podDisruptor = new PodDisruptor({
       namespace: 'default',
       labelSelector: 'app=api-service',
       gracePeriod: '10s',
       duration: '30s',
     });

     podDisruptor.terminate();

     // Test if system remains responsive
     const res = http.get('http://frontend-service');
     check(res, {
       'is status 200': (r) => r.status === 200,
       'response time < 500ms': (r) => r.timings.duration < 500,
     });
   }
   </code></div></div></pre>
3. 使用可觀測性堆疊監控實驗
4. 分析系統行為和故障模式
5. 實施改進和加固措施

## 項目總結與成果展示

 **學習目的** ：整合所學知識，展示完整的雲原生應用部署。

 **實作內容** ：

1. 完整架構文檔和拓撲圖
2. 所有 YAML 配置的 Git 倉庫
3. Helm Chart 包含前端、後端、數據庫和可觀測性
4. 可觀測性儀表板和告警配置
5. 混沌測試報告和韌性評估

這個學習路徑專注於使用 KinD 作為本地環境，著重在 DevSpace、Helm 以及完整的前後端架構，並保留可觀測性和混沌工程的重要內容。同時忽略了 ArgoCD 和 Kustomize 等工具，使學習更加聚焦和高效。
