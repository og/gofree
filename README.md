# gofree

**目前为内部版本,随时会修改接口,非  nimo 团队外的人不要使用**

> Go style ORM

1. 消除链式调用,让 SQL 一目了然
2. QueryBuilder 支持 DBA 审查 SQL
3. 接口友好,基于正确的场景设计接口
4. 接口设计松紧灵活,减少 `interface{}` 保持类型安全的同时提高易用性



## 消除链式调用

```go
// 某 ORM 的链式风格
db.Where("name = ?" , "nimo").
  Where("age = ?" ,18).
  Order("age DESC").First(&user)
```
链式风格粗看方便快捷，实际有以下缺点：

1. `"name = ?"` 编写繁琐，一旦漏掉 ` = ?`或写错容易出错 
2. 只查看代码难以分析最终执行的SQL是什么
3. 复杂查询会使得 SQL代码 分散

比如复杂情况下

```go
query := db.Where("name = ?" , "nimo").
  Where("age = ?" ,18)
if request.Gender != "" {
    query = query.Where("gender = ?", )
}
query.Order("age DESC").First(&user)
```        
 
 gofree 给出的解决方案是以结构体 QueryBuilder 作为查询条件消除链式调用。
 
 ```go
var foundUser bool
user := User{}
db.OneQB(&user, &foundUser, f.QB{
    Where:
    f.And("name", query.Name).
        And("age", query.Age).
        And("gender", f.EqualIgnoreEmpty(query.Gender)),
})
```

`f.EqualIgnoreEmpty(query string)` 的功能是如果 `query` 不是空字符串，最终执行时的 SQL 会包含`gender = ?`。

`foundUser` 是用于判断是否查询到数据。

`f.QB` 的结构体能一目了然的通过代码知道执行的SQL是什么样的，有利于SQL排查。