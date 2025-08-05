## 第一週：Kubernetes 基礎組件與 YAML 解析

### Day 1: 為什麼後端工程師必須學習 KinD？

* KinD 的定位與價值
* 容器化開發環境的重要性
* 本地 Kubernetes 環境對後端開發的影響

### Day 2: Kubernetes 組件架構與 YAML 基礎

* Kubernetes 核心組件與架構
* YAML 語法與結構基礎
* Kubernetes 資源定義的共同模式
  * apiVersion, kind, metadata, spec 等共同結構
* 實作：分析不同類型資源的 YAML 結構差異

### Day 3: Pod YAML 詳解 - 容器的基本單位

* Pod 定義與用途
* Pod YAML 結構深入解析
  * 容器定義
  * 環境變數
  * 資源限制
  * 重啟策略
* 實作：創建並測試各種 Pod 設定

### Day 4: Deployment YAML 詳解 - 應用管理與更新

* Deployment 與 ReplicaSet 關係
* Deployment YAML 結構深入解析
  * replicas 設定
  * 選擇器機制
  * 更新策略
  * 滾動更新設定
* 實作：部署應用並實驗不同更新策略

### Day 5: Service YAML 詳解 - 服務發現與負載均衡

* Service 類型與網路模型
* Service YAML 結構深入解析
  * 選擇器與標籤
  * 端口映射
  * 不同類型 Service 的配置差異
* 實作：為應用創建不同類型的 Service 並測試

### Day 6: ConfigMap 與 Secret YAML - 配置管理

* 配置與敏感資訊管理
* ConfigMap YAML 結構與生成方式
  * 內聯數據
  * 檔案引用
* Secret YAML 結構與安全考量
* 實作：從檔案與命令生成 ConfigMap 和 Secret

### Day 7: PersistentVolume 與 PVC YAML - 儲存管理

* Kubernetes 儲存模型
* PersistentVolume YAML 結構
* PersistentVolumeClaim YAML 結構
* 儲存類 (StorageClass) 設定
* 實作：設定持久化儲存並驗證資料保留

## 第二週：高級工作負載與 DevOps 整合

### Day 8: StatefulSet YAML - 有狀態應用部署

* StatefulSet 使用場景與特性
* StatefulSet YAML 結構深入解析
  * 穩定網路標識
  * 有序部署與更新
  * 持久化設定
* 實作：部署 PostgreSQL 數據庫

### Day 9: Job 與 CronJob YAML - 批處理工作

* 批處理工作負載類型
* Job YAML 結構詳解
  * 完成與重試策略
  * 並行處理設定
* CronJob YAML 結構詳解
  * 排程語法
  * 歷史限制與並發策略
* 實作：**使用 Flyway 進行數據庫遷移 Job**

### Day 10: Ingress YAML - 流量路由與 TLS

* Ingress 資源與控制器
* Ingress YAML 結構深入解析
  * 路由規則
  * TLS 設定
  * 註釋 (annotations) 與控制器特定設定
* 實作：配置 Ingress 並設定 HTTPS

### Day 11: NetworkPolicy YAML - 網路安全策略

* Kubernetes 網絡策略模型
* NetworkPolicy YAML 結構深入解析
  * 選擇器與標籤匹配
  * 入站與出站規則
  * 協議與端口設定
* 實作：實施網絡隔離規則

### Day 12: HorizontalPodAutoscaler YAML - 自動擴展

* Kubernetes 自動擴展機制
* HPA YAML 結構深入解析
  * 指標類型
  * 擴展策略
  * 冷卻期設定
* 實作：實施 CPU 與自定義指標擴展

### Day 13: 資源管理 YAML - 請求與限制

* Kubernetes 資源模型
* ResourceQuota YAML 結構
* LimitRange YAML 結構
* 實作：設定命名空間資源限制

### Day 14: Helm Chart 結構與模板

* Helm Chart 目錄結構
* YAML 模板化技術
* 值覆蓋機制
* 實作：創建可配置的應用 Chart

## 第三週：多環境部署與 GitOps 工作流

### Day 15: ConfigMap 與環境變數注入深入解析

* 環境變數注入方式比較
* ConfigMap 作為環境變數來源
* ConfigMap 作為檔案掛載
* 實作：**將 SQL 遷移腳本存儲在 ConfigMap 中**

### Day 16: 多環境配置策略 - YAML 環境特定變化

* 環境分離設計模式
* 使用 Kustomize 管理環境差異
* Helm 值文件環境管理
* 實作：建立開發/測試/生產環境配置

### Day 17: CI/CD 中的 YAML 生成與管理

* YAML 檔案的動態生成
* CI/CD 環境中的密鑰管理
* 版本控制策略
* 實作：**CI 管道中動態生成 ConfigMap**

### Day 18: GitOps 工作流與 YAML 版本控制

* GitOps 原則
* ArgoCD/Flux 配置結構
* YAML 宣告式部署
* 實作：設定基本 GitOps 管道

### Day 19: 從 Git 倉庫自動載入配置

* Kubernetes 與外部配置來源整合
* Git 倉庫與 Kubernetes 的橋接方式
* 使用 Init Container 從 Git 拉取配置
* 實作：**從 Git 倉庫載入遷移腳本**

### Day 20: Kubernetes Operator 與 CRD YAML

* 自定義資源定義 (CRD)
* 控制器模式與 Operator
* CRD YAML 結構分析
* 實作：部署資料庫 Operator

### Day 21: YAML 生成工具與最佳實踐

* YAML 生成工具比較
* 模板引擎與程式化生成
* YAML 質量與維護性
* 實作：使用程式化方法生成複雜 YAML

## 第四週：進階場景與整合

### Day 22: 多容器 Pod 設計模式與 YAML 配置

* Sidecar 模式
* Ambassador 模式
* Adapter 模式
* 實作：實現各種多容器設計模式

### Day 23: InitContainer 深入應用場景

* InitContainer 工作原理與生命週期
* 依賴檢查與準備工作
* 資源下載與初始化
* 實作：**使用 InitContainer 從 Git 提取遷移檔案**

### Day 24: 容器生命週期 Hook 與 YAML 設定

* postStart 與 preStop 鉤子
* 健康檢查與就緒探針
* 啟動與結束序列設計
* 實作：實現優雅啟動與關閉應用

### Day 25: 進階服務配置 - ExternalName 與 Headless Service

* Headless Service 用途與配置
* ExternalName Service 與外部整合
* 服務發現策略
* 實作：實現不同類型服務發現

### Day 26: 資料庫部署與管理 YAML 最佳實踐

* 有狀態資料庫部署考量
* 備份與恢復策略
* 連接與認證配置
* 實作：部署生產級 PostgreSQL 實例


### Day 27: 完整微服務架構搭建

* 目標：**建立一個完整的微服務示範環境**
* 內容：
  * 部署簡單前端應用 (React/Vue)
  * 部署多個後端服務 (REST API)
  * 設置 PostgreSQL 數據庫
  * 使用 Ingress 設定路由
  * 服務間通信配置
* 實作：
  * 創建完整的應用拓撲圖
  * 部署所有組件並確保互相通信
  * 建立端到端操作流程
  * YAML 清單含前端、API、數據庫的完整定義


### Day 28: OpenTelemetry Collector 與可觀測性基礎

* 目標：**在 KinD 上建立核心可觀測性基礎設施**
* 內容：
  * OpenTelemetry 架構概述
  * OpenTelemetry Collector 部署模式
  * Collector YAML 配置詳解：
    * 接收器 (receivers)
    * 處理器 (processors)
    * 導出器 (exporters)
    * 管道 (pipelines)
  * 接收各種格式的遙測數據
* 實作：
  * 部署 OpenTelemetry Operator
  * 配置 Collector DaemonSet
  * 設定基本數據收集管道
  * 實現自動檢測 (auto-instrumentation) 側車注入



### Day 29: Grafana 堆疊部署 - Mimir, Loki, Tempo

* 目標：**部署完整的 Grafana 可觀測性套件**
* 內容：
  * Grafana 作為統一觀測平台
  * Mimir 指標存儲系統部署與配置
    * 與 Prometheus 兼容的高可用指標系統
    * 長期存儲與查詢優化
  * Loki 日誌聚合系統部署與配置
    * 標籤索引與日誌查詢
    * 日誌流管道設定
  * Tempo 分布式追蹤系統部署與配置
    * 追蹤數據存儲與查詢
    * 與 OpenTelemetry 集成
* 實作：
  * 使用 Helm 部署 Grafana 堆疊
  * 配置 OTel Collector 向各系統發送數據
  * 設定跨組件關聯 (exemplars)
  * 建立統一的可觀測性儀表板



### Day 30: 可觀測性融合與應用檢測

* 目標：**將微服務應用與可觀測性堆疊集成**
* 內容：
  * 應用程式檢測策略：
    * 自動檢測配置
    * 手動檢測最佳實踐
  * 各語言 SDK 使用方法 (Node.js, Java, Go 等)
  * 使用 OpenTelemetry 注釋與標籤增強數據
  * 關聯 ID 傳播和上下文管理
  * RED 與 USE 監控方法論
* 實作：
  * 為所有服務啟用 OpenTelemetry 檢測
  * 配置自定義度量導出
  * 實現跨服務追蹤
  * 使用 Grafana 建立完整的監控儀表板
  * 設置多維度告警



### Day 31: xk6-disruptor 混沌測試與可觀測性驗證

* 目標：**實施混沌工程並利用可觀測性驗證系統韌性**
* 內容：
  * 混沌工程原則與實踐
  * xk6-disruptor 功能與用例：
    * Pod 故障注入
    * 網絡延遲與中斷
    * 資源壓力測試
  * k6 負載測試與 xk6-disruptor 整合
  * 使用可觀測性數據評估系統行為
  * 基於 SLO 的韌性測試
* 實作：
  * 部署 xk6-disruptor
  * 設計混沌實驗計劃
  * 執行網絡與資源故障注入
  * 使用 Grafana 堆疊分析影響
  * 制定改進建議與防護措施


```
                                  [Ingress Controller]
                                          |
                                          v
                [前端應用] <-----> [API Gateway/後端服務群] <-----> [PostgreSQL]
                     ^                     ^
                     |                     |
                     v                     v
             [OpenTelemetry Collector 集群]
                     |
          +----------+----------+
          |          |          |
          v          v          v
    [Grafana]    [Mimir]     [Loki]     [Tempo]
          ^          ^          ^          ^
          |          |          |          |
          +----------+----------+----------+
                     |
                     v
             [xk6-disruptor]

```
