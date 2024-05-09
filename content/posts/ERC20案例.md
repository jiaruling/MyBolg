---
title: "ERC20案例"
date: 2024-05-09T16:03:11+08:00
lastmod: 2024-05-09T16:03:11+08:00
author: ["me"]
categories: ["web3"] # 没有分类界面可以不填写
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

## 上下文变量

- 合约函数的背后是transaction，上下文变量访问的是transaction中的信息。
- 两个上下文变量: tx和msg

![](/image/context.png)

## ganache-cli

### 什么是 Ganache CLI？

Ganache CLI（前身为 TestRPC）是一个基于 Node.js 的命令行工具，用于快速搭建本地的以太坊区块链网络。它能够模拟完整的区块链环境，包括部署合约、交易确认等功能，方便开发者在本地环境中进行智能合约的开发和测试。

### 安装 Ganache CLI

使用 npm 来安装 Ganache CLI，只需要在命令行中运行 `npm install -g ganache-cli`。安装完成后，就可以在命令行中使用 `ganache-cli` 命令来启动 Ganache CLI。

### 使用 Ganache CLI 进行开发和测试

启动 Ganache CLI 后，它会默认在本地创建一个开发用的区块链网络。可以通过访问 http://localhost:8545 来查看区块链网络的状态和交易信息。

在启动的区块链网络中，你可以：

1. 部署智能合约：使用 `web3.js` 或其他以太坊开发库，通过 Ganache CLI 部署和测试智能合约。
2. 发送交易：在本地环境中测试交易的发送和确认速度，以及合约交互的效果。
3. 账户管理：Ganache CLI 提供了一些默认的测试账户，你可以使用这些账户进行测试，也可以自行导入现有账户。
4. 虚拟时间控制：你可以手动增加区块、调整区块时间，以模拟各种场景。

### 总结

Ganache CLI 是一个强大的工具，为以太坊开发者提供了一个本地化的开发和测试环境。它的简便安装和易于使用，使得开发者能够更高效地测试和调试智能合约和 DApp。在区块链开发过程中，Ganache CLI 是不可或缺的利器。

## MetaMask钱包的安装和使用

1. 导入本地的测试网络`Ganache-cli`
   - 网络名称: `localhost:8545`
   - RPC url: `http://127.0.0.1:8545`
   - 链ID: 1337
   - 货币符号: ETH
2. 通过的私钥导入账户
3. 通过浏览器Remix(Injected Provider - MetaMask)与MetaMask钱包相链接

## ERC20

> 同质化资产：例如货币，股票，没有唯一性。
> 地址: [ERC20 - OpenZeppelin Docs](https://docs.openzeppelin.com/contracts/4.x/erc20)

```solidity
// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.0 <0.9.0;

contract BalanceManager {
    mapping(address=>uint256) public balanceOf;

    string public name = "MYDOLLAR";
    string public symbol = "$";
    uint8 public decimals = 4;

    constructor(uint256 total){
        balanceOf[msg.sender] = total;
    }

    function transfer(address to, uint256 amount) public  {
        address from = msg.sender;
        uint256 fb = balanceOf[from];
        uint256 tb = balanceOf[to];

        require(amount <= fb, "from account do not have enough money!");

        fb -= amount;
        tb += amount;
        balanceOf[from] = fb;
        balanceOf[to] = tb;
    }
}
```

- 注意：在Ganache的测试环境中，使用0.8.20以下的版本进行编译，要不然部署的时候会编译环境和部署环境不匹配的错误。
