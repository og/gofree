# scan

> scan 用于将 go sql 返回的 rows (map[string]interface{}) 解析到结构体 

表结构

```sql
CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL DEFAULT '',
  `is_super` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `book` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `price` decimal(10,2) NOT NULL,
  `user_id` int(11) unsigned NOT NULL,
  `last_read_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
```


````go
type User struct {
	ID int `db:"id"`
	Name string `db:"name"`
	IsSuper bool `db:"is_super"`
}
type Book struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Price float64 `db:"price"`
	UserID int `db:"user_id"`
	LastReadTime time.Time `db:"last_read_time"`
}
````

User 与 Book 是一对多的关系

```go
type UserAndBook struct {
	User User `db:"user"`
	BookList []Book `db:"book"`
}
```

数据内容
```sql
INSERT INTO `user` (`id`, `name`, `is_super`)
VALUES
	(1, 'nimo', 0),
	(2, 'nico', 1);


INSERT INTO `book` (`id`, `name`, `price`, `user_id`, `last_read_time`)
VALUES
	(1, 'go action', 12.00, 1, '2019-11-11 00:00:00'),
	(2, 'js action', 11.00, 1, '2019-11-11 00:00:00'),
	(3, 'life', 1.00, 0, '2019-10-25 07:40:53');
```

查询语句
```sql
SELECT
	  user.id             AS "user.id"
	, user.name           AS "user.name"
	, user.is_super       AS "user.is_super"
	
	, book.id             AS "book.id"
	, book.name           AS "book.name"
	, book.price          AS "book.price"
	, book.user_id        AS "book.user_id"
	, book.last_read_time AS "book.last_read_time"
	
FROM 
	user JOIN book ON user.id = book.user_id
WHERE 
	user.id = 1
```
查询结果

```text
user.id	user.name	user.is_super	book.id	book.name	book.price	book.user	book.last_read_time
1	    nimo	        0	            1	go action	  12.00	        1	    2019-11-11 00:00:00
1	    nimo        	0	            2	js action     11.00     	1	    2019-11-11 00:00:00
```

返回结果是多列的，一般我们期望转换为如下数据

```go
userAndBookList := []UserAndBook{
    User: User {
        ID: 1,
        Name: "nimo",
        IsSuper: false,
    },
    BookList: []Book{
        Book: {
            ID 1,
            Name: "go action",
            Price: 12,
            LastReadTime: time.Time{2019-11-11 00:00:00},
        },
        Book: {
            ID 2,
            Name: "js action",
            Price: 11,
            LastReadTime: time.Time{2019-11-11 00:00:00},
        }
    }
} 
```
 通过 sqlx `rows.MapScan()` 可以将sql 查询的行转换为 map[string]interface{}
 即
 
```go
map[string]interface{}{
	{
        "user.id": 1,
        "user.name": "nimo",
        "user.is_super": 0,
        "book.id": 1,
        "book.name": "go action",
        "book.price": 12,
        "book.user": 1,
        "book.last_read_time": time.Time{2019-11-11 00:00:00},
    },
    {
        "user.id": 1,
        "user.name": "nimo",
        "user.is_super": 0,
        "book.id": 2,
        "book.name": "js action",
        "book.price": 11,
        "book.user": 1,
        "book.last_read_time": time.Time{2019-11-11 00:00:00},
    }
}
```

想要将 `rows.MapScan()` 返回的 map 转换为 []UserAndBook 需要根据 UserAndBook 分析到如下信息


```go
type UserAndBook struct {
	User User `db:"user"`
	BookList []Book `db:"book"`
}


type MapRelationData struct {
	FieldIndex int
	TableName string
	DBTag map[int]string `note:"key is fieldName , value is struct tag 'db'"`
}
type MapRelation struct {
	Single []MapRelationData
	Many []MapRelationData
}
mapRelation :=  MapRelation{
	Single: []MapRelationData {
        {
        	FieldIndex: 0,
        	TableName: "user",
        	DBTag map[int]string{
        		0: "id",
        		1: "name",
        		2: "is_super"
        	},
        }		
	},
	Many: []MapRelationData {
		{
			FieldIndex: 1,
			TableName: "book",
			DBTag map[int]string{
				0: "id",
				1: "name",
				2: "price",
				3: "user_id",
				4: "last_read_time",
			},
		},
	},
}
```

然后通过反射根据 `MapRelation` 填充 `[]UserAndBook{}` 。


