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

4. gorm 与 sqlx 相关的 db 示例都只能初始化一次，然后传参，如果每次请求创建一次 db 则会报 `to many connections` 错误。

   1. 疑问：将 sql db 示例化成单例，这也有隐患，在并发时（insert + update）会出现并发问题，我认为是要每次请求创建一个 db 实例，**然后及时 close 掉。**

5. 设置 mysql 连接池，首先可以先查看 mysql 默认的最大连接池和全局最大连接池：

   ```cmd
   show variables like '%max_connections%';	// 显示最大连接属
   show global status like 'Max_used_connections'; // 服务器响应的最大连接数
   show status like 'Threads%'; // 实时查看占用的连接数和线程数
   ```

   然后在初始化 db 时设置相关的信息：

   ```go
   sqlDB, err := sql.Open("mysql", dsn)	// 初始化 db
   sqlDB.SetMaxIdleConns(10)	// 设置最大空闲连接数
   sqlDB.SetMaxOpenConns(100)	// 设置最大连接数
   sqlDB.SetConnMaxLifetime(time.Millisecond * 200)	// 每个连接的最大生存周期，一般情况是不用设置的
   ```

6. 设置每个连接的生存周期，在高并发测试下（500 个用户，循环 200 次，10 万个请求），会报 `Only one usage of each socket address (protocol/network address/port) is normally permitted.`，这个错误是因为 mysql 连接耗尽了 socket 数导致的错误，在 mysql 服务器下输入以下命令查看 socket 连接数发现在没有设置 `SetConnMaxLifetime`，socket 连接数就是 `SetMaxOpenConns` 设置的最大连接数。就不会存在这个问题，具体原因还没有查明！
