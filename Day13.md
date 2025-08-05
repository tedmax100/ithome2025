---
title: "Day 13：健康檢查與 DevSpace 狀態監控"
tags: 2025鐵人賽
date: 2025-07-20
---

#### Day 13：健康檢查與 DevSpace 狀態監控
##### 重點
Liveness、Readiness Probe 概念與設定
DevSpace 的狀態監控功能與 UI 儀表板
##### Lab
為 Deployment 加入健康檢查
設定 DevSpace 等待部署健康後才繼續
使用 DevSpace UI 監控應用健康狀態


在 Kubernetes 環境中，應用程式的健康狀態監控是確保服務穩定性的關鍵。今天我們將深入探討 Kubernetes 的健康檢查機制，以及如何使用 DevSpace 進行狀態監控。

# 健康檢查基礎概念
## Liveness Probe
Liveness Probe 用於確定容器是否處於運行狀態。如果探測失敗，Kubernetes 會重啟容器。
```
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 3
  periodSeconds: 3
```

Readiness Probe
Readiness Probe 用於確定容器是否準備好接收流量。失敗時，Pod 會從服務負載均衡中移除。
```
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

Startup Probe
Startup Probe 用於處理需要較長啟動時間的應用，確保其他探針在應用完全啟動前不會干擾。
```
startupProbe:
  httpGet:
    path: /startup
    port: 8080
  failureThreshold: 30
  periodSeconds: 10
```

# DevSpace 健康檢查整合
DevSpace 提供了與 Kubernetes 健康檢查機制的緊密整合。以下是一個完整的配置示例：

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
        - kubernetes/*.yaml
    wait:
      enabled: true
      timeout: 300
      conditions:
        - type: PodReady
          status: "True"
        - type: Available
          status: "True"

dev:
  sync:
    - imageSelector: app
      localSubPath: ./
      containerPath: /app
  ports:
    - imageSelector: app
      forward:
        - port: 8080
          localPort: 8080
  ui:
    enabled: true
```

# 實作範例：健康檢查應用
1. 建立基礎應用程式

```
// app.js
const express = require('express');
const app = express();

// 應用狀態管理
let isHealthy = true;
let isReady = false;

// 模擬應用啟動過程
setTimeout(() => {
  isReady = true;
  console.log('Application is ready to serve traffic');
}, 10000);

// 健康檢查端點
app.get('/health', (req, res) => {
  if (isHealthy) {
    res.status(200).json({ status: 'healthy' });
  } else {
    res.status(500).json({ status: 'unhealthy' });
  }
});

// 就緒檢查端點
app.get('/ready', (req, res) => {
  if (isReady) {
    res.status(200).json({ status: 'ready' });
  } else {
    res.status(503).json({ status: 'not ready' });
  }
});

// 啟動檢查端點
app.get('/startup', (req, res) => {
  if (isReady) {
    res.status(200).json({ status: 'started' });
  } else {
    res.status(503).json({ status: 'starting' });
  }
});

// 監控指標端點
app.get('/metrics', (req, res) => {
  res.json({
    uptime: process.uptime(),
    memoryUsage: process.memoryUsage(),
    cpuUsage: process.cpuUsage()
  });
});

const port = process.env.PORT || 8080;
app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});
```

Dockerfile 配置
```
FROM node:16-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 8080

CMD ["node", "app.js"]
```

 Kubernetes 部署配置
 ```
 # kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: health-check-demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: health-check-demo
  template:
    metadata:
      labels:
        app: health-check-demo
    spec:
      containers:
      - name: app
        image: health-check-demo
        ports:
        - containerPort: 8080
        startupProbe:
          httpGet:
            path: /startup
            port: 8080
          failureThreshold: 30
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
```

DevSpace 配置
```
# devspace.yaml
version: v2beta1
name: health-check-demo

images:
  app:
    image: health-check-demo
    dockerfile: ./Dockerfile

deployments:
  app:
    kubectl:
      manifests:
        - kubernetes/*.yaml
    wait:
      enabled: true
      timeout: 300
      conditions:
        - type: PodReady
          status: "True"
        - type: Available
          status: "True"

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
        - port: 8080
          localPort: 8080
  ui:
    enabled: true
    metrics:
      - name: "Application Uptime"
        query: "process_uptime_seconds"
      - name: "Memory Usage"
        query: "process_resident_memory_bytes"
```

# DevSpace UI 監控功能
DevSpace 提供了豐富的 UI 監控功能，可以幫助開發者即時了解應用狀態：

Pod 狀態監控
    容器運行狀態
    資源使用情況
    健康檢查結果

日誌查看
    即時日誌流
    日誌過濾
    多容器日誌整合

資源使用監控
    CPU 使用率
    記憶體使用率
    網路流量

自定義指標
    應用特定指標
    性能指標
    業務指標


## 最佳實踐
健康檢查設計
    確保檢查輕量且快速
    避免過度檢查影響性能
    合理設置超時時間

監控指標選擇
    選擇關鍵業務指標
    監控系統資源使用
    追蹤錯誤率和延遲

告警配置
    設置合適的閾值
    配置適當的告警頻率
    定義清晰的告警規則



# 總結
健康檢查和狀態監控是確保應用可靠運行的關鍵。通過合理配置 Kubernetes 的健康檢查機制，結合 DevSpace 的監控功能，我們可以：

及時發現並處理應用問題
確保服務的高可用性
優化資源使用
提升開發效率
在實際開發中，應根據應用特性選擇合適的健康檢查策略，並利用 DevSpace 提供的工具進行有效的狀態監控。

# Lab 練習
練習 1：實現完整的健康檢查
實現上述範例應用
添加自定義健康檢查邏輯
測試不同的故障場景
練習 2：配置 DevSpace 監控
設置 DevSpace UI
配置自定義監控指標
實現告警機制
練習 3：性能測試
使用 Apache Bench 進行負載測試
觀察健康檢查行為
分析監控數據