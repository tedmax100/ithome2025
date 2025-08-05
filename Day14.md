---
title: "Day 14 DevSpace 與 Helm 整合"
tags: 2025鐵人賽
date: 2025-07-20
---


#### Day 14 DevSpace 與 Helm 整合
##### 重點
DevSpace 支援 Helm 部署
Helm 基礎與 Chart 結構
DevSpace 的 Helm 值覆寫機制
##### Lab
設定 DevSpace 使用 Helm 部署應用
建立簡單的 Helm chart 並透過 DevSpace 部署
使用 DevSpace 的值覆寫功能適應不同環境

在現代 Kubernetes 應用部署中，Helm 已成為一個重要的包管理工具。今天我們將探討如何在 DevSpace 中整合 Helm，實現更靈活的應用部署管理。

# 為什麼要在 DevSpace 中使用 Helm？
## 簡化部署流程

DevSpace 可直接使用 Helm chart 進行部署
無需切換工具，統一開發體驗

## 環境配置管理

透過 DevSpace 的值覆寫機制管理不同環境
開發、測試、生產環境配置集中管理

## 開發效率提升

快速切換不同版本的應用配置
即時同步本地更改到集群


# Helm 基礎概念
什麼是 Helm？
Helm 是 Kubernetes 的包管理器，它可以：

打包應用及其依賴
版本管理
配置管理
簡化部署流程

## Helm Chart 結構
```
mychart/
  ├── Chart.yaml          # Chart 的基本信息
  ├── values.yaml         # 默認配置值
  ├── charts/            # 依賴的子 charts
  └── templates/         # 模板文件
      ├── deployment.yaml
      ├── service.yaml
      └── _helpers.tpl
```

# DevSpace 與 Helm 整合
基本配置
```
# devspace.yaml
version: v2beta1
name: my-app

deployments:
  my-release:
    helm:
      chart:
        name: ./charts/my-app
      values:
        image:
          repository: my-app
          tag: latest
        service:
          type: ClusterIP
          port: 80
```

值覆寫機制
```
# devspace.yaml
deployments:
  my-release:
    helm:
      chart:
        name: ./charts/my-app
      values:
        # 基礎值
        replicaCount: 1
      valuesFiles:
        - values.yaml
        - values.dev.yaml
```

# DevSpace Helm 實用功能
1. 自動重新部署
當 Helm 配置發生變更時，DevSpace 會自動重新部署：
```yaml
deployments:
  backend:
    helm:
      chart:
        name: ./charts/my-app
      # 啟用自動重新部署
      replaceImageTags: true
      # 等待部署就緒
      wait: true
      timeout: 180
```

條件部署
根據不同環境或條件部署不同配置：
```yaml
deployments:
  backend:
    helm:
      chart:
        name: ./charts/my-app
      # 根據環境變數決定配置
      values:
        environment: ${DEVSPACE_PROFILE}
        debug: ${ENABLE_DEBUG}
```

3. 依賴管理
處理多個 Helm chart 之間的依賴關係：
```yaml
deployments:
  backend:
    helm:
      chart:
        name: ./charts/my-app
      # 安裝依賴
      dependencies:
        - name: mongodb
          repository: https://charts.bitnami.com/bitnami
          version: "10.0.0"
```


# 實作範例：建立完整的 Helm Chart
1. 創建 Helm Chart 結構
```helm create my-app```

2. 配置 Chart.yaml
```yaml
# charts/my-app/Chart.yaml
apiVersion: v2
name: my-app
description: A Helm chart for my application
type: application
version: 0.1.0
appVersion: "1.0.0"
```

3. 配置 values.yaml
```yaml
# charts/my-app/values.yaml
replicaCount: 1

image:
  repository: my-app
  pullPolicy: IfNotPresent
  tag: "latest"

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 50m
    memory: 64Mi
```

4. 創建部署模板
```yaml
# charts/my-app/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "my-app.fullname" . }}
  labels:
    {{- include "my-app.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "my-app.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "my-app.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - containerPort: 80
              protocol: TCP
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
```

5. 配置 DevSpace
```yaml
# devspace.yaml
version: v2beta1
name: my-app

images:
  app:
    image: my-app
    dockerfile: ./Dockerfile
    tags:
      - latest
      - ${DEVSPACE_GIT_COMMIT}

deployments:
  app:
    helm:
      chart:
        name: ./charts/my-app
      values:
        image:
          repository: my-app
          tag: ${DEVSPACE_GIT_COMMIT}
        ingress:
          enabled: true
          hosts:
            - host: dev.example.com
              paths:
                - path: /
                  pathType: Prefix

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
        - port: 80
          localPort: 3000
```

環境特定配置
開發環境配置
```yaml
# values.dev.yaml
replicaCount: 1
resources:
  limits:
    cpu: 200m
    memory: 256Mi
ingress:
  enabled: true
  hosts:
    - host: dev.example.com
```

生產環境配置
```
# values.prod.yaml
replicaCount: 3
resources:
  limits:
    cpu: 500m
    memory: 512Mi
ingress:
  enabled: true
  hosts:
    - host: prod.example.com
```


DevSpace 配置檔案結構化
```
# devspace.yaml
version: v2beta1
name: my-app

# 定義映像檔
images:
  app:
    image: my-app
    dockerfile: ./Dockerfile

# Helm 部署配置
deployments:
  app:
    helm:
      chart:
        name: ./charts/my-app
      values:
        image:
          repository: my-app
          tag: ${DEVSPACE_IMAGE_TAG}
      valuesFiles:
        - values.yaml
        - values.${DEVSPACE_PROFILE}.yaml

# 開發環境配置
dev:
  sync:
    - imageSelector: app
      localSubPath: ./
      containerPath: /app
  ports:
    - imageSelector: app
      forward:
        - port: 80
          localPort: 3000

# 環境配置檔案
profiles:
  - name: dev
    patches:
      - op: merge
        path: deployments.app.helm.values
        value:
          ingress:
            enabled: true
            hosts:
              - host: dev.example.com
  
  - name: prod
    patches:
      - op: merge
        path: deployments.app.helm.values
        value:
          ingress:
            enabled: true
            hosts:
              - host: prod.example.com
```


Lab 練習
Lab 1：建立基本 Helm Chart
創建新的 Helm Chart：`helm create my-app`


修改 values.yaml：
```yaml
replicaCount: 1
image:
  repository: my-app
  tag: latest
service:
  type: ClusterIP
  port: 80
```

配置 DevSpace：
```yaml
version: v2beta1
name: my-app

deployments:
  app:
    helm:
      chart:
        name: ./charts/my-app
```


Lab 2：環境特定配置
創建環境特定的值文件：
```yaml
# values.dev.yaml
environment: development
replicaCount: 1
```
```
# values.prod.yaml
environment: production
replicaCount: 3
```


在 DevSpace 中使用不同環境：

```yaml
version: v2beta1
name: my-app

deployments:
  app:
    helm:
      chart:
        name: ./charts/my-app
      valuesFiles:
        - values.yaml
        - values.${DEVSPACE_PROFILE}.yaml

profiles:
  - name: dev
  - name: prod
```


Lab 3：自動化部署流程
創建部署腳本：
```bash
#!/bin/bash
# deploy.sh

ENVIRONMENT=$1

if [ -z "$ENVIRONMENT" ]; then
  echo "Usage: ./deploy.sh <environment>"
  exit 1
fi

devspace use profile $ENVIRONMENT
devspace deploy
```

執行部署：
```bash
chmod +x deploy.sh
./deploy.sh dev  # 部署到開發環境
./deploy.sh prod # 部署到生產環境
```

總結
通過整合 Helm 和 DevSpace，我們可以：

標準化部署流程

使用 Helm Chart 管理應用配置
統一不同環境的部署方式
靈活的配置管理

環境特定的值覆寫
簡化配置變更
改進的開發體驗

快速切換環境
簡化測試和除錯流程
下一步學習