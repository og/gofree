# 维护手册

向 gofree 贡献代码需要解读此文档

首先一定要使用并阅读 gofree 所有面向用户的文档并且深度使用过 gofree，了解 gofree 的接口和设计理念。

## QueryBuilder

gofree 在代码实现层的核心是 QueryBuilder，一切都是围绕 QueryBuilder 去实现的。 以下将 QueryBuilder 简称QB

QB 结构体完全是基于SQL语法进行设计的，了解SQL就能了解 QB 各个字段大概的用途。

一切基于QB设计与编码的目的是让使用者尽可能的使用 QB，使用 QB 能让最终写出的 go ORM 代码像SQL一样，
并且 QB 的类型要做到松紧结合，尽可能的避免 `interface{}` 只在必要的时候开放 `interface{}`。

目前 QB 中只有 `type Data map[Column]interface{}` 使用了 `interface{}` 。
`Data` 供 update 与 insert 使用。


### qb.SQL()

