---
title: "Gas与转账收款"
date: 2024-05-13T14:11:26+08:00
lastmod: 2024-05-13T14:11:26+08:00
author: ["me"]
categories: [""] # 没有分类界面可以不填写
tags: [""] # 标签
description: ""
weight:
slug: ""
draft: false # 是否为草稿
comments: true # 本页面是否显示评论
reward: true # 打赏
mermaid: true #是否开启mermaid
showToc: true # 显示目录
TocOpen: true # 自动展开目录
hidemeta: false # 是否隐藏文章的元信息，如发布日期、作者等
disableShare: true # 底部不显示分享栏
showbreadcrumbs: true #顶部显示路径
#   cover:
#     image: "" #图片路径例如：posts/tech/123/123.png
#     zoom: # 图片大小，例如填写 50% 表示原图像的一半大小
#     caption: "" #图片底部描述
#     alt: ""
#     relative: false
---

## 区块链是一个经济系统

- 计算与存储资源都是稀缺的，区块链的工作需要消耗资源
- 共识、trustless需要矿工的工作，而矿工需要激励
- Transaction的执行有成本(gas)，gas费成为矿工的激励
- ether(Native token)是这个经济生态系统的通行货币

![](/image/gas.png)

### 关心的问题

- 合约执行中的经济成本，即gas问题
- 智能合约实现货币的流通，即转账收款功能

### 货币单位

| Uint       | wei Value      | wei                       | ether value    |
| ---------- | -------------- | ------------------------- | -------------- |
| wei        | 1 wei          | 1                         | $10^{-18}$ ETH |
| kwei       | $10^3$ wei     | 1,000                     | $10^{-15}$ ETH |
| mwei       | $10^6$ wei     | 1,000,000                 | $10^{-12}$ ETH |
| gwei       | $10^9$ wei     | 1,000,000,000             | $10^{-9}$ ETH  |
| microether | $10^{12}$ wei  | 1,000,000,000,000         | $10^{-6}$ ETH  |
| milliether | $10^{15}$ wei  | 1,000,000,000,000,000     | $10^{-3}$ ETH  |
| ether      | $10^{18}$  wei | 1,000,000,000,000,000,000 | 1 ETH          |

```solidity
// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

contract EtherUnits {
    uint public oneWei = 1 wei;
    bool public isOneWei = 1 wei == 1;

    uint public oneEther = 1 ether;
    bool public isOneEther = 1 ether == 1e18;
}
```

### 合约持有ether

- **address.balance**: 合约可以有钱
- 合约与其它外部账号或者EOA之间可以转账
- multisig钱包

### gas、gas fee、gas price

- 实际的gas是完全由执行逻辑决定的，一个固定逻辑的合约函数执行gas没有变换
- gas price是由市场定价
- `gasfee = gas * gasprice`

### Gaslimit与Gasleft()函数

- 交易发起者设定最多消耗多少: gaslimit
- 合约之间调用，调用者可以设定gaslimit
- 区块本身有一个gaslimit(社区决定的)
- Gasleft()是以上因素作为限定，与当前gas消耗一起计算的结果

### 退款规则

- 剩余没有用完的gas会”退款“
- 如果可用gas耗尽，会终止交易执行
- 交易失败，已经用了的gas不退

### gas编码实战

### 转账

#### 设计思路

- 理解转账的关键是理解合约收款的**设计安排**
- 设计安排拆分为两步:
  - 被调用函数解析逻辑
    1. 自定义函数匹配selector
    2. Receive函数匹配为空的calldata
    3. Fallback函数兜底
    4. 除了fallback必须在尾部，其它元素顺序无关
  - 检查逻辑
    - 解析逻辑如果成功则会输出一个函数，检查逻辑就是检查一个情况: `value > 0`并且这个函数没有被payable修饰。如果出现这个情况，调用失败终止，否则执行函数逻辑。

#### receive()函数

> 单纯转账calldata为空，为了使得fallback的职责清晰，solidity安排了一个特殊函数receive()来处理它

#### 转账系统的遗留方式

- Solidity中的转账函数send和transfer是旧的转账设计，有各种缺陷
- 新的转账设计没有专门的转账函数，而是普通函数调用的伴生物
- send和transfer就是gaslimit为2300的calldata为空的call，区别在于transfer处理了call的返回值
- 建议使用新的转账设计
