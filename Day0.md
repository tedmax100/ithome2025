---
title: "Day01_為什麼後端工程師必須學習 KinD？"
tags: 2025鐵人賽
date: 2025-07-20
---

今年是 AI 生態圈快速發展的一年，許多應用在快速萌芽，但身為後端開發者的我來說，打好基礎還是比較重要的。過去幾年我專注於語言、架構或可觀測性工程的學習，近期則發現自己對 Kubernetes（k8s）還不夠熟悉。因此，這個月我決定從頭學習 Kubernetes 的基礎。不過，身為開發者而非 Infra 工程師，單機的 k8s 環境就很夠用了。這次，我選擇了 KinD（Kubernetes in Docker）作為學習主題，因為它更貼近現代 CI/CD、測試自動化等情境，也更容易在不同平台上快速建立多節點叢集。

## 什麼是 KinD？

[KinD（Kubernetes in Docker）](https://kind.sigs.k8s.io/)是一個可以在本地端利用 Docker 容器快速建立 Kubernetes 叢集的工具。它特別適合開發、測試、CI/CD pipeline 等場景，讓你不用依賴雲端或實體機器，就能輕鬆模擬多節點的 k8s 環境。

## 為什麼後端工程師應該學習 KinD？

1. **現代雲端架構趨勢**  
   現在許多後端應用都會部署在 Kubernetes 上，學會 KinD 就等於學會如何在本地模擬這種雲端環境。你可以先在本地驗證應用的部署、設定與運作，降低在雲端出錯的機率。

2. **開發與測試自動化利器**

    KinD 天生適合整合到 CI/CD pipeline，常被用於自動化測試環境。你可以在本機或 CI server 上快速建立、銷毀叢集，大幅提升測試與開發效率。


3. **熟悉 k8s 多節點與網路行為**

    KinD 支援多節點部署，幫助你了解 k8s 真實運作時的網路、服務發現、Pod 調度等細節，這是 Minikube 較難模擬的。


4. **容器化思維與除錯**  
   KinD 完全基於 Docker，讓你更熟悉容器生態，也方便直接測試本地自建 Image 的行為。

5. **團隊協作與跨平台一致性**  
   KinD 的叢集定義完全程式化（YAML），可納入版本控制，讓團隊每個人都能用相同設定快速建立一致的 k8s 環境。

---

# 預計的學習路線

## 第一週：Kubernetes 與 KinD 基礎

#### Day 1：雲原生與 Kubernetes 概念
##### 重點
- 雲原生、Kubernetes 架構、核心元件（Pod、Node、Cluster、Service、Deployment）
##### Lab
- 整理 Kubernetes 架構圖，並用文字描述每個元件功能

#### Day 2：KinD 安裝與啟動
##### 重點
- KinD、kubectl 安裝與設定
##### Lab
- 安裝 KinD 並啟動本地叢集，使用 kubectl get nodes 驗證

#### Day 3：Kubernetes 資源與 YAML 基礎
##### 重點
- Kubernetes 資源種類、YAML 結構說明
##### Lab
- 撰寫一個簡單的 Pod YAML，kubectl apply 部署並觀察狀態

#### Day 4：Deployment 與 ReplicaSet
##### 重點
- Deployment、ReplicaSet 概念與用途
##### Lab
- 建立 Deployment，調整 replicas 數量，觀察自動擴縮容器

#### Day 5：Service 與 Kubernetes 網路模型
##### 重點
Service 類型（ClusterIP、NodePort、LoadBalancer）
Kubernetes 網路模型基礎
KinD 內建的 CNI 介紹

##### Lab
部署 NodePort Service，從本機 curl 訪問服務
測試不同 Service 類型的行為
觀察網路封包流動

#### Day 6: Flannel CNI 深入解析
重點
Kubernetes CNI 概念與架構
Flannel CNI 工作原理與組件
網路模型與封包流向
各種後端模式比較
Lab
建立使用 Flannel CNI 的自訂 KinD 叢集
測試 Pod 間通訊
分析 Flannel 網路行為

#### Day 7：資源管理（Label、Selector、Namespace）
##### 重點
- Label 與 Selector 用法、Namespace 管理
##### Lab
- 建立多個 Namespace，利用 Label 管理與查詢資源

#### Day 8：Pod 日誌與 Debug
##### 重點
- kubectl logs、exec、describe 等指令使用
##### Lab
- 模擬應用錯誤，透過 logs/exec/describe 排查問題

---

## 第二週：開發流程與 Skaffold 整合
#### Day 9：自訂映像檔部署與開發工作流程
##### 重點
如何讓 KinD 使用本地 Docker 映像檔（kind load docker-image）
開發工作流程問題：手動部署的繁瑣步驟
##### Lab
將本地映像檔直接用於 Deployment，驗證服務可用性
記錄開發過程中需要執行的重複步驟（建置→載入→部署→測試）

#### Day 10：DevSpace 入門與快速開發循環
##### 重點
DevSpace 介紹與安裝
基本 devspace.yaml 配置
開發模式（devspace dev）與檔案監控
DevSpace UI 儀表板介紹
DevSpace 與 VS Code 整合

##### Lab
安裝 DevSpace
建立基本 devspace.yaml
分析 devspace.yaml 的各個部分
使用 devspace dev 實現程式碼變更→自動部署的開發循環
體驗 DevSpace UI 監控應用狀態
設定 VS Code 與 DevSpace 整合

#### Day 11：KinD 進階配置管理與 DevSpace 多應用開發與環境變數
##### 重點
KinD 進階配置管理

多節點叢集設計模式與最佳實踐

為不同應用場景定制 KinD 配置

配置版本控制與團隊共享策略

使用腳本自動化 KinD 環境管理

DevSpace 與 KinD 配置聯動

在 DevSpace 工作流中創建/銷毀 KinD 叢集

使用 DevSpace hooks 管理 KinD 生命週期

配置 DevSpace 使用特定 KinD 叢集

DevSpace 多應用（前後端）設定

結合 ConfigMap 與環境變數

使用 DevSpace profiles 區分環境

##### Lab
建立包含前後端的 devspace.yaml

使用 ConfigMap 注入環境變數

設定 dev/staging profiles

實作練習

設計三種不同用途的 KinD 配置（開發、測試、模擬生產）

創建自動化腳本管理 KinD 環境生命週期

配置 DevSpace 與不同 KinD 環境整合

使用 ConfigMap 為不同環境注入配置

設計 DevSpace profiles 對應不同環境需求

實現完整的開發-測試工作流程

#### Day 12：Volume 與持久化儲存（整合 DevSpace
##### 重點
Volume、PersistentVolume、PersistentVolumeClaim 概念
DevSpace 與持久化資料的處理
##### Lab
建立 PVC 並掛載於 Pod，測試資料持久性
在 DevSpace 開發循環中處理持久化資料
配置 DevSpace 的持久化資料同步功能

#### Day 13：健康檢查與 DevSpace 狀態監控
##### 重點
Liveness、Readiness Probe 概念與設定
DevSpace 的狀態監控功能與 UI 儀表板
##### Lab
為 Deployment 加入健康檢查
設定 DevSpace 等待部署健康後才繼續
使用 DevSpace UI 監控應用健康狀態


## 第三週：Helm 與 CI/CD 整合

- Day 14-16：從基礎到進階，讓學習者先掌握 Helm 核心概念
  - 先需要部署簡單應用（Day 14）
  - 然後需要處理複雜配置（Day 15-16）
- Day 17-18：將 Helm 知識應用到實際的 Kubernetes 工作負載上
  - 最後處理特殊工作負載（Day 17-18）

```
基礎 Helm 概念
    ↓
Chart 開發技能
    ↓
應用到多容器場景
    ↓
應用到批次作業場景
```

為 CI/CD 做準備
Day 14-18 建立的 Chart 和配置，正好是 Day 19-20 CI/CD 要用到的
學習者會有完整的 Chart 可以用來實作自動化部署

Day 14：DevSpace 與 Helm 基礎整合
重點
Helm 基礎概念介紹
Chart 結構與最佳實踐
DevSpace 支援 Helm 部署
Lab
設定 DevSpace 使用 Helm 部署應用
建立簡單的 Helm chart
使用 DevSpace 的值覆寫功能

Day 15：Helm Chart 開發深入（上）
重點
Chart.yaml 與 values.yaml 設計
模板函數與管道
條件判斷與流程控制
Lab
開發完整的應用 Chart
使用條件判斷處理不同環境
實作常用 Helper 模板

Day 16：Helm Chart 開發深入（下）
重點
依賴管理與子 Chart
Hook 機制與生命週期
Chart 測試與驗證
Lab
建立具有依賴關係的 Chart
實作部署前後的 Hook
編寫 Chart 測試

Day 17：多容器 Pod 與 DevSpace 同步更新
重點
多容器設計模式
InitContainer 使用場景
DevSpace 的檔案同步機制
Lab
設計多容器 Pod 的 Helm Chart
配置 DevSpace 檔案同步
測試熱重載功能

Day 18：Job 與 CronJob 的 Helm 實作
重點
Job 與 CronJob 的 Helm 模板
批次作業最佳實踐
DevSpace 的 Job 測試策略
Lab
建立 Job Chart
實作 CronJob 模板
使用 DevSpace 測試 Job

Day 19：GitHub Actions 與 CI/CD（上）
重點
GitHub Actions 基礎
Docker 映像檔建置流程
自動化測試設定
Lab
設定基本 CI pipeline
配置映像檔建置與推送
實作自動化測試
Day 20：GitHub Actions 與 CI/CD（下）
重點
多環境部署策略
Helm Chart 發布流程
版本控制與標籤管理
Lab
實作完整 CD 流程
配置 Chart 發布
自動化版本管理
Day 21：Service Account 與 RBAC 安全性
重點
Service Account 設計
RBAC 權限管理
環境隔離策略
Lab
建立 Service Account
配置 RBAC 規則
實作環境隔離

---
## 第四週：CI/CD 與團隊協作
#### Day DevSpace 與 CI/CD 整合
##### 重點
DevSpace CI/CD 模式（devspace deploy、devspace build）
整合 GitHub Actions 或 GitLab CI
DevSpace 的無頭模式（headless mode）

##### Lab
設定 CI pipeline 使用 DevSpace 自動測試與部署
模擬 PR 流程中的自動測試
配置 DevSpace 的無頭模式運行

#### Day 23：團隊開發工作流程設計
##### 重點
團隊共享 KinD 與 DevSpace 配置
Git 工作流程與環境分離
DevSpace 的團隊協作功能
##### Lab
設計團隊開發工作流程
建立共享配置與腳本
配置 DevSpace 的團隊協作設定

#### Day 24：微服務開發與測試策略
##### 重點
微服務架構下的本地開發策略
服務依賴與模擬（Service Mocking）
DevSpace 的多組件開發
##### Lab
設計微服務測試策略
使用 DevSpace 管理多服務開發
配置服務間依賴關係

---
# 第五週：可觀測性與進階開發技術
#### Day 25：Prometheus、Grafana 與開發指標
##### 重點
使用 Helm 與 DevSpace 部署監控堆疊
應用程式指標導出與監控
DevSpace 與監控工具的整合

##### Lab
部署 Prometheus & Grafana
設定應用程式導出指標
建立開發指標 Dashboard
配置 DevSpace 自動部署監控堆疊

#### Day 26：日誌收集與開發除錯
##### 重點
Loki 部署與配置
開發過程中的結構化日誌
DevSpace 的日誌查看功能
##### Lab
部署 Loki
配置應用程式產生結構化日誌
使用 DevSpace 查看集中式日誌
配置 Grafana 查詢開發日誌


#### Day 27：分散式追蹤與微服務除錯
##### 重點
Tempo 部署與 OpenTelemetry 整合
微服務間的追蹤與效能分析
DevSpace 與追蹤工具的整合
##### Lab
部署 Tempo 與 OpenTelemetry Collector
在應用程式中加入追蹤
分析微服務調用鏈
使用 DevSpace 管理追蹤配置

#### Day 28：OpenTelemetry Collector 進階配置
##### 重點
OpenTelemetry Collector 架構深入解析
收集器部署模式（Agent vs Gateway）
自定義處理管道（pipeline）與處理器（processors）
##### Lab
部署 OTel Collector 為 DaemonSet（每節點 Agent 模式）
配置多種資料來源（logs, metrics, traces）收集
設定資料轉換與過濾處理器
使用 DevSpace 管理 OTel 配置更新

#### Day 29：OpenTelemetry 自動檢測與 Kubernetes 整合
##### 重點
OpenTelemetry 自動檢測（Auto-instrumentation）技術
Kubernetes 運算子（Operator）與側車注入（Sidecar Injection）
服務網格（Service Mesh）與 OTel 整合
##### Lab
部署 OpenTelemetry Operator
配置自動檢測注入（Java, Node.js, Python 等）
使用註解（annotations）控制檢測行為
實現零程式碼變更的可觀測性
使用 DevSpace 簡化可觀測性堆疊部署

#### Day 30：混沌測試與可觀測性驗證
##### 重點
xk6-disruptor 混沌測試
可觀測性驗證與 SLO 監控
使用 OpenTelemetry 分析系統行為
DevSpace 在測試環境中的應用

##### Lab
部署測試微服務架構
使用 xk6-disruptor 注入網路延遲與中斷
透過 OpenTelemetry 收集的指標、日誌與追蹤分析影響
建立可觀測性驅動的韌性改進
使用 DevSpace 管理測試環境生命週期


