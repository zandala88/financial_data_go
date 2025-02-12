# financial_data_go
# 金融数据可视化与预测系统

![Language](https://img.shields.io/badge/language-golang-brightgreen)

## 一、系统介绍

本项目大致分为以下6个模块，分别为

1. 上市公司信息
2. 股票市场信息
3. 公募基金信息
4. 期货数据信息
5. 宏观经济信息
6. 用户个人界面

## 二、功能展示

### 1. 上市公司信息

该模块提供上司公司的信息展示，可通过省份进行筛选或者进行搜索

![company-main.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/company-main.png?raw=true)

![company-detail.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/company-detail.png?raw=true)

### 2. 股票市场

股票市场模块中，分为主界面和单支股票的详情界面

#### 2.1 主界面

拥有数据展示图表以及每支股票的具体信息，并且可以点击查看详情内容

![stock-main-1.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/stock-main-1.png?raw=true)

![stock-main-2.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/stock-main-2.png?raw=true)

#### 2.2 详情界面

详情界面中有 **K线图**、**移动平均线图**、**MACD图**、**RSI图**，此外还有其他多个功能，重点为预测结果功能，会通过训练的LSTM模型对下一次交易日的收盘价进行预测

![stock-detail.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/stock-detail.png?raw=true)

![stock-pre.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/stock-pre.png?raw=true)

### 3. 公募基金

公募基金模块和股票市场模块类似，此处不过多阐述

![fund-main.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/fund-main.png?raw=true)

![fund-main-2.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/fund-main-2.png?raw=true)

### 4. 期货数据

期货数据模块拥有交易日历的展示以及期货数据的展示、期货数据的展示包括**成交量**、**成交金额**、**持仓量**、**收盘价**、**同比数据**

![futures-main-1.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/futures-main-1.png?raw=true)

### 5. 宏观经济

宏观经济模块展示了**Shibor利率**，**历年GDP**、**居民消费指数**

![economics.png](https://github.com/zandala88/financial_data_go/blob/main/.github/img/economics.png?raw=true)

## 三、技术栈

### 1. 数据库

数据库使用`MySQL`以及`Redis`

### 2. 后端

#### 2.1 预测相关-Python

预测模型：`pytorch`

预测接口：`flask`

#### 2.2 其他-Golang

`Gin`+`Gorm`

### 3.前端

前端程序基于`Vue3.x`开发，UI框架基于`Element Plus`，开发语言为`javascript`，图表插件主要使用到了`echarts`

## 鸣谢

### 数据来源 - tushare

感谢tushare提供免费的数据

官网：https://tushare.pro



