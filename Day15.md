# Day 15：Helm Chart 開發深入（上）

## 重點
- Chart.yaml 與 values.yaml 設計
- 模板函數與管道
- 條件判斷與流程控制

## 學習目標
- 深入了解 Chart.yaml 與 values.yaml 的設計原則
- 掌握 Helm 模板函數與管道的使用
- 學會條件判斷與流程控制
- 開發完整的應用 Chart

---

## 1. Chart.yaml 設計原則

### 為什麼 Chart.yaml 很重要？

Chart.yaml 就像是你的 Helm Chart 的「身分證」，它告訴 Helm 和其他人這個 Chart 是什麼、版本多少、依賴什麼。想像你在開發一個應用程式，Chart.yaml 就是你的 package.json 或 requirements.txt。

### 基本結構解析

```yaml
apiVersion: v2              # 使用 Helm 3 的 API 版本
name: my-app               # Chart 名稱，必須是小寫字母和連字符
description: 完整的後端應用 Helm Chart  # 簡短描述
type: application          # Chart 類型：application 或 library
version: 0.1.0            # Chart 版本（遵循語義化版本）
appVersion: "1.0.0"       # 應用程式版本
dependencies:             # 依賴的其他 Charts
  - name: postgresql
    version: "~12.1.0"    # 波浪號表示接受相容版本
    repository: "https://charts.bitnami.com/bitnami"
    condition: postgresql.enabled  # 條件式依賴
```

**重要概念說明：**

- **Chart version vs appVersion**：
  - Chart version 是你這個部署模板的版本
  - appVersion 是你實際應用程式的版本
  - 例如：你的應用是 v2.0.0，但你的部署模板可能還是 v0.1.0

- **依賴管理**：
  - `condition` 讓你可以選擇性安裝依賴
  - `~12.1.0` 表示接受 12.1.x 的任何版本，但不接受 12.2.0

### 版本管理策略

這是很多人容易搞混的地方。讓我用實際例子說明：

```yaml
# 情境 1：應用程式更新了，但部署方式沒變
version: 0.1.0        # Chart 版本不變
appVersion: "1.2.0"   # 應用版本更新

# 情境 2：新增了 Ingress 支援，這是部署模板的變更
version: 0.2.0        # Chart 版本更新（新功能）
appVersion: "1.2.0"   # 應用版本可能不變

# 情境 3：修復了模板中的 bug
version: 0.1.1        # Chart 版本修補更新
appVersion: "1.2.0"   # 應用版本不變
```

---

## 2. values.yaml 設計架構

### 為什麼需要好的 values.yaml 設計？

values.yaml 是你的 Chart 的「控制面板」。一個好的設計讓使用者可以：
- 輕鬆理解有哪些可配置選項
- 快速找到需要修改的設定
- 安全地覆蓋預設值而不破壞其他設定

### 結構化配置的思維

想像你在設計一個汽車的控制面板，你會把相關的控制項放在一起：

```yaml
# 全域設定 - 影響整個 Chart 的設定
global:
  imageRegistry: ""        # 統一的映像倉庫
  imagePullSecrets: []     # 統一的拉取密鑰
  
# 應用設定 - 關於你的應用本身
app:
  name: my-app
  version: "1.0.0"
  
# 映像設定 - 容器映像相關
image:
  repository: myorg/my-app  # 映像倉庫
  tag: ""                   # 空字串表示使用 Chart.appVersion
  pullPolicy: IfNotPresent  # 拉取策略
```

**設計原則：**
1. **分組邏輯**：相關設定放在同一個區塊
2. **預設值策略**：提供合理的預設值
3. **註釋說明**：重要設定要有註釋
4. **一致性**：命名風格要一致

### 環境特定配置的重要性

在真實世界中，開發、測試、生產環境的需求完全不同：

```yaml
# values-dev.yaml - 開發環境
image:
  tag: "dev-latest"      # 開發環境用最新版本
  pullPolicy: Always     # 總是拉取最新映像

resources:
  limits:
    cpu: 200m            # 開發環境資源需求較小
    memory: 256Mi

env:
  normal:
    NODE_ENV: development
    LOG_LEVEL: debug     # 開發環境需要詳細日誌

# values-prod.yaml - 生產環境
image:
  tag: "v1.2.3"         # 生產環境用固定版本

resources:
  limits:
    cpu: 1000m          # 生產環境需要更多資源
    memory: 1Gi

autoscaling:
  enabled: true         # 生產環境啟用自動擴展
  minReplicas: 3
  maxReplicas: 20

env:
  normal:
    NODE_ENV: production
    LOG_LEVEL: warn     # 生產環境只記錄警告以上
```

**為什麼這樣設計？**
- **開發環境**：重視開發效率，需要詳細日誌，資源可以較少
- **生產環境**：重視穩定性和效能，需要固定版本、自動擴展

---

## 3. 模板函數與管道

### 什麼是模板函數？

Helm 使用 Go 的模板引擎，模板函數就像是小工具，幫你處理資料。管道（pipe）讓你可以串連多個函數。

### 常用函數實例解析

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "my-app.fullname" . }}    # 使用 helper 函數產生完整名稱
  labels:
    {{- include "my-app.labels" . | nindent 4 }}  # 使用管道縮排
spec:
  replicas: {{ .Values.replicaCount | default 1 }}  # 預設值函數
```

**函數說明：**
- `include`：呼叫其他模板函數
- `nindent 4`：縮排 4 個空格（YAML 格式要求）
- `default 1`：如果 replicaCount 沒設定，使用 1

### 字串處理的實際應用

```yaml
metadata:
  name: {{ .Values.app.name | lower | replace "_" "-" }}
  labels:
    version: {{ .Chart.AppVersion | quote }}
    environment: {{ .Values.environment | upper }}
    region: {{ .Values.region | default "us-west-2" }}
```

**為什麼需要這些處理？**
- `lower`：Kubernetes 資源名稱必須小寫
- `replace "_" "-"`：Kubernetes 不允許底線，要轉成連字符
- `quote`：確保值被正確引用，避免 YAML 解析錯誤
- `upper`：標籤值可能需要大寫格式

### 管道的威力

管道讓你可以串連多個操作：

```yaml
# 原始值：My_App_Name
# 處理後：my-app-name
name: {{ .Values.rawName | lower | replace "_" "-" | trunc 63 | trimSuffix "-" }}
```

這個管道做了什麼？
1. `lower`：轉小寫 → my_app_name
2. `replace "_" "-"`：替換底線 → my-app-name
3. `trunc 63`：截斷到 63 字符（Kubernetes 限制）
4. `trimSuffix "-"`：移除結尾的連字符

---

## 4. 條件判斷與流程控制

### 為什麼需要條件判斷？

不是每個部署都需要相同的資源。例如：
- 開發環境可能不需要 Ingress
- 小型部署可能不需要自動擴展
- 某些環境可能使用外部資料庫

### 基本條件判斷

```yaml
# templates/ingress.yaml
{{- if .Values.ingress.enabled -}}    # 只有啟用時才建立
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ include "my-app.fullname" . }}
spec:
  {{- if .Values.ingress.className }}   # 可選的 className
  ingressClassName: {{ .Values.ingress.className }}
  {{- end }}
  rules:
    {{- range .Values.ingress.hosts }}  # 迴圈處理多個主機
    - host: {{ .host | quote }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: {{ include "my-app.fullname" $ }}  # 注意這裡用 $
                port:
                  number: {{ $.Values.service.port }}
    {{- end }}
{{- end }}
```

**重要概念：**
- `{{- if ... -}}`：條件判斷，`-` 移除多餘空白
- `{{- range ... }}`：迴圈處理陣列
- `$`：在迴圈中引用根上下文

### 複雜邏輯處理

```yaml
# 環境變數的智慧處理
env:
  # 基本環境變數
  {{- range $key, $value := .Values.env.normal }}
  - name: {{ $key | quote }}
    value: {{ $value | quote }}
  {{- end }}
  
  # 資料庫連線（根據配置決定）
  {{- if .Values.database }}
  - name: DB_HOST
    value: {{ .Values.database.host | quote }}
  {{- if .Values.database.existingSecret }}
  # 使用現有的 Secret
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: {{ .Values.database.existingSecret }}
        key: password
  {{- else }}
  # 直接使用配置的密碼（不推薦用於生產環境）
  - name: DB_PASSWORD
    value: {{ .Values.database.password | quote }}
  {{- end }}
  {{- end }}
```

**這段程式碼的邏輯：**
1. 先處理所有一般環境變數
2. 如果有資料庫配置，加入資料庫相關環境變數
3. 密碼優先使用 Secret，其次使用直接配置

---

## 5. Helper 模板實作

### 為什麼需要 Helper 模板？

想像你在寫程式，如果同樣的邏輯要重複寫很多次，你會把它抽取成函數。Helper 模板就是 Helm Chart 的「函數」。

### 核心 Helper 模板解析

```yaml
{{/*
應用名稱 - 基礎的名稱生成
*/}}
{{- define "my-app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}
```

**這個 Helper 做什麼？**
1. 如果用戶設定了 `nameOverride`，使用它
2. 否則使用 Chart 名稱
3. 截斷到 63 字符（Kubernetes 限制）
4. 移除結尾可能的連字符

```yaml
{{/*
完整應用名稱 - 包含 Release 名稱的完整名稱
*/}}
{{- define "my-app.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}
```

**這個邏輯處理什麼情況？**
- 如果用戶設定了完整名稱覆蓋，直接使用
- 如果 Release 名稱已經包含應用名稱，避免重複
- 否則組合 Release 名稱和應用名稱

### 標籤 Helper 的重要性

```yaml
{{/*
通用標籤 - 符合 Kubernetes 建議的標籤
*/}}
{{- define "my-app.labels" -}}
helm.sh/chart: {{ include "my-app.chart" . }}
{{ include "my-app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
```

**為什麼這些標籤重要？**
- `helm.sh/chart`：標識這是由哪個 Chart 建立的
- `app.kubernetes.io/version`：應用版本，方便追蹤
- `app.kubernetes.io/managed-by`：標識管理工具（Helm）

---

## Lab：開發完整的應用 Chart

### 實作步驟詳解

#### 步驟 1：建立 Chart 結構

```bash
# 建立新 Chart
helm create my-backend-app
cd my-backend-app

# 清理預設檔案，我們要從頭開始
rm -rf templates/*
rm values.yaml
```

**為什麼要清理？**
因為 `helm create` 產生的預設模板比較簡單，我們要建立更完整的版本。

#### 步驟 2：理解 Deployment 模板

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "my-backend-app.fullname" . }}
  labels:
    {{- include "my-backend-app.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount | default 1 }}
  {{- end }}
```

**這段邏輯的思考：**
- 如果啟用了自動擴展，就不要設定固定的 replicas
- 否則使用配置的副本數，預設為 1

#### 步驟 3：測試與驗證的重要性

```bash
# 驗證語法 - 檢查模板語法是否正確
helm lint .

# 測試模板渲染 - 看看實際產生的 YAML
helm template my-backend-app .

# 測試條件判斷 - 確保條件邏輯正確
helm template my-backend-app . --set ingress.enabled=true
```

**為什麼要這樣測試？**
- `helm lint`：抓出語法錯誤
- `helm template`：看到實際的 Kubernetes YAML
- 條件測試：確保不同配置下都能正確工作

---

## 總結

### 關鍵學習重點

1. **Chart.yaml 設計**
   - 版本管理不只是數字，要有策略
   - 依賴管理要考慮條件式載入
   - 元數據要完整，方便其他人理解

2. **values.yaml 架構**
   - 結構化設計讓配置更清晰
   - 環境特定配置避免一個檔案管所有環境
   - 合理的預設值減少配置負擔

3. **模板函數與管道**
   - 函數是工具，管道是流水線
   - 字串處理要考慮 Kubernetes 的命名規則
   - 預設值處理讓 Chart 更健壯

4. **條件判斷與流程控制**
   - 不是所有資源都需要在所有環境中建立
   - 複雜邏輯要分解成清晰的條件判斷
   - 迴圈處理讓配置更靈活

### 實際應用建議

- **開發階段**：多用 `helm template` 測試
- **測試階段**：在不同環境用不同的 values 檔案
- **生產階段**：確保版本管理策略清晰

通過今天的學習，你已經掌握了開發專業級 Helm Chart 的核心技能。明天我們會學習更進階的主題，包括 Hook 機制和安全性考量。
