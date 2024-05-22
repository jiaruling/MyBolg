---
title: "Delegatecall与代理模式"
date: 2024-05-15T10:29:15+08:00
lastmod: 2024-05-15T10:29:15+08:00
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

## 存储布局-Storage Layout

> 成员变量按照他们出现的顺序在storage中按照一种规则依次堆放，这使得每个成员变量具有固定的位置。在slot中一个变量从低位开始存

![](/image/slot.png)

### 值类型的堆叠规则

- 值类型需要的空间就是这个数据类型的数据块大小(字节)
- 值类型数据如果不能在当前slot剩余空间存的下，则新起一个slot
- 结构和数组不会与其他数据分享slot

### 动态数组和映射

- 动态数组和mapping参与成员变量的堆叠占有一个32位slot,但是数据通过哈希运算存在别的存储slot中，不参与堆叠
- 动态数组的slot存放数组大小，mapping空着

### layout中的递归和继承

- struct和array的数据的堆叠跟上述成员变量的堆叠规则一样
- 对于合约继承，按照3c线性化的结果以此堆叠，并且允许父子数据共享于一个slot
- 如官方文档所言，layout规则并不是内部技术，而是与语言的外部行为有关

### 编码实战

```solidity
// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

contract StorageLayout {
    uint256 public num; // slot0
    address public sender; // slot1
    Person person; // slot2 3 4
    bool[12] successs; //slot5
    mapping(address=>uint) balances; // slot6
    uint256 public value; //slot7

    struct Person {
        uint256 num; // slot0
        address sender; // slot1
        bool[12] success; // slot1
        uint256 value; // slot2
    }

    function setVars(uint256 _num) public payable {
        num = _num;
        sender = msg.sender;
        value = msg.value;
    }
}
```

### 总结

一个合同有成员变量，每个成员变量根据它的类型占有固定长度的空间并按照规则堆叠上去。这个成员变量序列的布置是固定的，每个成员变量的位置是固定的。**合约编译成的机器码是用这个位置来访问它的！在幕后，合约机器代码只认识位置，不认变量名称！**

## delegatecall的作用和机制

- 语法: `<address>.delegatecall(byte calldata)`
- 作用: 调用别的合约代码，访问自己的成员
- 理解:
  - A合约把B合约的函数借用、搬运到自己合约内部来执行，这些代码访问的是A合约的成员。
  - **由于合约是按照成员变量的存储位置来访问成员变量，如果A和B的成员变量存储布局相同(或兼容)，那么这种借用就能正确执行。**

- 上下文变化:
  - 概念上: 合约B完全运行在合约A的上下文中，`this`、`msg.sender`、`msg.value`均无变化，msg就是合约A的调用者产生的message！
  - 底层实现: 产生了一个message拷贝，比如gaslimit可以重新设置
  - 合约调用链中，delegatecall是调用者合约管辖范围的扩大

### 编码实战

```solidity
// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

contract B {
    uint public num;
    address public sender;
    uint public value;

    function setVars(uint _num) public payable {
        num = _num;
        sender = msg.sender;
        value = msg.value;
    }
}

contract A {
    uint public num;
    address public sender;
    uint public value;

    function setVars(address _contract,uint _num) public payable {
       (bool success, bytes memory data) = _contract.delegatecall{gas: 10000000}(abi.encodeWithSignature("setVars(uint256)", _num));
       require(success, "failed");
    }
}
```

- 先部署合约B，再部署合约A。
- 调用合约A中的`setVars`函数，参数`_contract`是B合约的地址
- 可以观察到，A合约中的成员变量发生了变化，而B合约中的成员变量任然为初始值。

## 代理模式

> 代理模式解决合约升级问题

### 基础工作原理

#### 工作机制

1. 调用者调用proxy的setX()
2.  Proxy的setX()不存在，fallback()触发
3.  fallback使用delegatecall调用logic，将setX()调用的calldata传入
4.  logic的setX()被执行，但访问的是Proxy的x成员

![](/image/proxy1.png)

#### 升级

- 代理模式中，proxy负责数据存储，logic负责数据的逻辑处理
- 升级就是proxy将它的logic成员变量切换到一个新的处理逻辑

![](/image/proxy2.png)

#### 编码实战

```solidity
// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;


contract LogicV1 {
    address placeholder; // 起占位符的作用
    uint256 public count;

    function inc() external {
        count += 1;
    }
}

contract LogicV2 {
    address placeholder; // 起占位符的作用
    uint256 public count;

    function inc() external {
        count += 2; // 业务逻辑的内容发生改变
    }
}


interface LogicInterface {
    function inc() external;
}

contract Proxy {
    address public logic;
    uint256 public count; // 对存储空间的解释

    constructor(address _logic) {
        logic = _logic;
    }

    fallback() external {
        (bool ok, bytes memory res) = logic.delegatecall(msg.data);
        require(ok, "delegatecall failed");
    }

    // 合约升级
    function upgradeTo(address newVersion) external {
        logic = newVersion;
    }
}
```

- 首先部署合约`LogicV1`和`LogicV2`
- 再部署合约`Proxy`, 此时合约`Proxy`的`logic`成员变量指向`LogicV1`合约实例的地址
- 此时，`Proxy`的合约实例没有`inc()`函数，需要用`LogicInterface`接口加载`Proxy`合约实例地址，使得`Proxy`合约实例拥有`inc()`函数
- 在`Proxy`合约实例中调用`inc()`函数，观察`Proxy`合约实例中`count`变量的变化。(每点击一次count变量+1)
- 使用`Proxy`合约实例中的`upgradeTo`函数，将合约升级为`LogicV2`
- 再在`Proxy`合约实例中调用`inc()`函数，观察`Proxy`合约实例中`count`变量的变化。(每点击一次count变量+2)
- 使用 `fallback`和`delegatecall`函数，完成了合约的升级

### 非结构化代理模式

- `proxy`提供空白存储，对存储的”解释权“完全由`logicV1`负责，数据与逻辑相分离
- 想办法让`proxy`中的`logic`不参与storage成员的堆叠，也就不必再`logicV1`中出现`placeholder`
- `proxy`成为通用的、存储的代理，只安排代理控制逻辑，跟具体业务无关，而`logicV1`成为纯粹干净的业务逻辑，跟代理机制无关

#### 编码实战

```solidity
```

## 库合约

> 库函数调用是合约内部调用，上下文不变化

### 库合约定义

- 库定义: `library` 关键字
- 不能继承别的合约，只能实现接口
- 不能有构造函数、成员变量、修饰器

### public 和 internal

- public和external库函数是通过delegatecall调用，但是调用的参数约束与内部调用相同
- internal是通过编译时代码内联实现
- 一个库如果含有public函数必然要单独部署，其地址通过编译时嵌入或者部署时作为部署参数传入

#### 编码实战

```solidity
```

