# Golang 问题

1. 在执行 sql 语句并映射到实体类时，必须要将 sql 中的字段和实体中的属性一一对应。否则会报 `missing destination type` 错误。所以要注意慎用 `select * from table`。

2. 如果实体对应的字段不是必填的，注意将对应的类型设置 `sql.NullString` 或 `sql.NullXXX`

3. 如果实体与数据库中的字段不匹配， 可以显示利用 Go 标签特性来指定映射关系：

   ```go
   type User struct {
   	ID           uint
   	Name         string
   	Email        *string
   	Age          uint8
   	Birthday     *time.Time
   	MemberNumber sql.NullString `db:"member_number"`
   	ActivedAt    sql.NullTime   `db:"actived_at"`
   	CreatedAt    time.Time      `db:"created_at"`
   	UpdatedAt    sql.NullString `db:"updated_at"`
   	DeletedAt    sql.NullString `db:"deleted_at"`
   }
   ```

   

