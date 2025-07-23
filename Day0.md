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

#### Day 10：Skaffold 入門與快速開發循環
##### 重點
Skaffold 介紹與安裝
基本 skaffold.yaml 配置
開發模式（skaffold dev）與檔案監控
##### Lab
安裝 Skaffold
建立基本 skaffold.yaml
使用 skaffold dev 實現程式碼變更→自動部署的開發循環

#### Day 11：KinD 進階配置管理與 Skaffold 多應用開發與環境變數
##### 重點
KinD 進階配置管理
多節點叢集設計模式與最佳實踐
為不同應用場景定制 KinD 配置
配置版本控制與團隊共享策略
使用腳本自動化 KinD 環境管理

Skaffold 與 KinD 配置聯動
在 Skaffold 工作流中創建/銷毀 KinD 叢集
使用 Skaffold hooks 管理 KinD 生命週期
配置 Skaffold 使用特定 KinD 叢集

Skaffold 多應用（前後端）設定
結合 ConfigMap 與環境變數
使用 Skaffold profiles 區分環境

##### Lab
建立包含前後端的 skaffold.yaml
使用 ConfigMap 注入環境變數
設定 dev/staging profiles

實作練習
設計三種不同用途的 KinD 配置（開發、測試、模擬生產）
創建自動化腳本管理 KinD 環境生命週期
配置 Skaffold 與不同 KinD 環境整合
使用 ConfigMap 為不同環境注入配置
設計 Skaffold profiles 對應不同環境需求
實現完整的開發-測試工作流程

#### Day 12：Volume 與持久化儲存（整合 Skaffold）
##### 重點
Volume、PersistentVolume、PersistentVolumeClaim 概念
Skaffold 與持久化資料的處理
##### Lab
建立 PVC 並掛載於 Pod，測試資料持久性
在 Skaffold 開發循環中處理持久化資料

#### Day 12：健康檢查與 Skaffold 狀態監控
##### 重點
Liveness、Readiness Probe 概念與設定
Skaffold 的 statusCheck 功能
##### Lab
為 Deployment 加入健康檢查
設定 Skaffold 等待部署健康後才繼續

#### Day 13：Skaffold 與 Helm 整合
##### 重點
Skaffold 支援 Helm 部署
Helm 基礎與 Chart 結構
##### Lab
設定 Skaffold 使用 Helm 部署應用
建立簡單的 Helm chart 並透過 Skaffold 部署

#### Day 14：Skaffold 除錯模式與資源監控
##### 重點
Skaffold debug 模式
遠端除錯設定
資源限制與監控
##### Lab
使用 skaffold debug 進行遠端除錯
設定 Pod 資源限制，觀察 Skaffold 監控輸出


---

## 第三週：Kubernetes 進階功能與開發整合
#### Day 15：多容器 Pod 與 Skaffold 同步更新
##### 重點
多容器設計模式、InitContainer 用法
Skaffold 的檔案同步（file sync）功能
##### Lab
設計多容器 Pod，並在 Skaffold 中設定檔案同步
測試程式碼變更不重建容器的快速更新

#### Day 16：Job、CronJob 與 Skaffold 整合
##### 重點
批次作業（Job）、定時任務（CronJob）
如何在開發流程中測試 Job
##### Lab
建立 Job 與 CronJob
使用 Skaffold 部署並測試 Job 執行
#### ay 17：Service Account、RBAC 與開發安全性
##### 重點
Service Account、RBAC 權限控管
開發環境的安全性考量
##### Lab
建立 Service Account，設定簡單 RBAC 規則
在 Skaffold 開發流程中驗證權限設定

#### Day 18：自訂 Helm Chart 與 Skaffold 值覆寫
##### 重點
撰寫自訂 Helm Chart
Skaffold 的 Helm 值覆寫功能
##### Lab
將應用包成 Helm Chart
使用 Skaffold 的 setValues 或 valuesFiles 功能

#### Day 19：Kustomize 與 Skaffold 整合
##### 重點
Kustomize 基礎與用法
Skaffold 支援 Kustomize 部署
##### Lab
建立 Kustomize 配置
設定 Skaffold 使用 Kustomize 部署

#### Day 20：K8s Dashboard 與開發可視化
##### 重點
在 KinD 叢集啟用 Dashboard
結合 Skaffold port-forward 功能
##### Lab
啟用 Dashboard，設定 Skaffold 自動 port-forward
使用 Dashboard 監控開發中的應用

#### Day 21：開發環境常見問題排解
##### 重點
Skaffold 與 KinD 常見錯誤
排錯流程與最佳實踐
##### Lab
模擬常見問題並排解（映像檔問題、網路問題等）

---
## 第四週：CI/CD 與團隊協作
#### Day 22：Skaffold 與 CI/CD 整合
##### 重點
Skaffold CI/CD 模式（skaffold run/build/deploy）
整合 GitHub Actions 或 GitLab CI
##### Lab
設定 CI pipeline 使用 Skaffold 自動測試與部署
模擬 PR 流程中的自動測試

#### Day 23：團隊開發工作流程設計
##### 重點
團隊共享 KinD 與 Skaffold 配置
Git 工作流程與環境分離
##### Lab
設計團隊開發工作流程
建立共享配置與腳本

#### Day 24：微服務開發與測試策略
##### 重點
微服務架構下的本地開發策略
服務依賴與模擬（Service Mocking）
##### Lab
設計微服務測試策略
使用 Skaffold 管理多服務開發


---
# 第五週：可觀測性與進階開發技術
#### Day 25：Prometheus、Grafana 與開發指標
##### 重點
使用 Helm 與 Skaffold 部署監控堆疊
應用程式指標導出與監控
##### Lab
部署 Prometheus & Grafana
設定應用程式導出指標
建立開發指標 Dashboard

#### Day 26：日誌收集與開發除錯
##### 重點
Loki 部署與配置
開發過程中的結構化日誌
Skaffold 日誌集中顯示
##### Lab
部署 Loki
配置應用程式產生結構化日誌
使用 Grafana 查詢開發日誌

#### Day 27：分散式追蹤與微服務除錯
##### 重點
Tempo 部署與 OpenTelemetry 整合
微服務間的追蹤與效能分析
##### Lab
部署 Tempo 與 OpenTelemetry Collector
在應用程式中加入追蹤
分析微服務調用鏈

#### Day 28：OpenTelemetry Collector 進階配置
##### 重點
OpenTelemetry Collector 架構深入解析
收集器部署模式（Agent vs Gateway）
自定義處理管道（pipeline）與處理器（processors）
##### Lab
部署 OTel Collector 為 DaemonSet（每節點 Agent 模式）
配置多種資料來源（logs, metrics, traces）收集
設定資料轉換與過濾處理器
使用 Skaffold 管理 OTel 配置更新

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

#### Day 30：混沌測試與可觀測性驗證
##### 重點
xk6-disruptor 混沌測試
可觀測性驗證與 SLO 監控
使用 OpenTelemetry 分析系統行為
##### Lab
部署測試微服務架構
使用 xk6-disruptor 注入網路延遲與中斷
透過 OpenTelemetry 收集的指標、日誌與追蹤分析影響
建立可觀測性驅動的韌性改進

