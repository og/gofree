# gofree


> Go style ORM

1. 消除链式调用,让 SQL 一目了然
2. QueryBuilder 支持 DBA 审查 SQL
3. 不只是ORM,更提供一套数据层编码指导方案



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

## 审查SQL

不少团队需要DBA审核或自审SQL以避免低性能SQL。gofree 支持SQL审查。

在项目文件例如：`app/model/sql/sql.go` 文件中存放所有 SQL 字符串常量

```go
const SQLSelectUserByNameAge = "SELECT `id`, `name`, `is_super`, `created_at`, `updated_at`, `deleted_at` FROM `user` WHERE `age` = ? AND `name` = ? AND `deleted_at` IS NULL"
const SQLSelectUserByNameAgeGender = "SELECT `id`, `name`, `is_super`, `created_at`, `updated_at`, `deleted_at` FROM `user` WHERE `age` = ? AND `gender` = ? AND `name` = ? AND `deleted_at` IS NULL"
```

使用 gofree 时通过 `f.QB{Check: []string{}}` 或 `db.OneID().Check(sqls ...[]string)` 配置SQL审查

```go
db.OneQB(&user, &foundUser, f.QB{
    Where:
    f.And("name", query.Name).
        And("age", query.Age).
    Check: []string{SQLSelectUserByNameAge},
})
```

使用 `[]string` 而不是 `string` 是因为有些查询在 gofree 内部可能由多个sql组成。

在运行时，gofree 会检查执行时的 SQL 与 Check 配置的 SQL 是否一致，不一致将会安全的提醒SQL与预期不一致。

