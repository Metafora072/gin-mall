## .ini 文件

### INI 文件详解

#### **1. INI 文件基本概念**

INI（Initialization）文件是一种简单的**键值对**配置文件格式，广泛用于存储应用程序的配置信息。其特点如下：

| **特点**       | **说明**                                                   |
| -------------- | ---------------------------------------------------------- |
| **层级结构**   | 使用 `[section]` 划分配置块，每个块内包含多个键值对。      |
| **键值对格式** | `key = value`，等号（`=`）或冒号（`:`）均可作为分隔符。    |
| **注释支持**   | 分号 `;` 或井号 `#` 开头的行视为注释。                     |
| **跨平台兼容** | 无复杂语法，适合轻量级配置，广泛支持 Windows/Linux/macOS。 |

#### **2. INI 文件示例**

```ini
; 数据库配置
[database]
host = 127.0.0.1
port = 3306
username = root
password = secret
enable_ssl = true

; 服务器配置
[server]
listen = :8080
log_level = debug
```

------

### Go 语言读取 INI 文件的实现方法

#### **1. 使用第三方库 `go-ini`**

Go 标准库未内置 INI 解析器，推荐使用社区维护的 **[go-ini/ini](https://github.com/go-ini/ini)** 库。

##### **安装库**

```bash
go get github.com/go-ini/ini
```

#### **2. 基础用法示例**

```go
package main

import (
	"fmt"
	"log"
	"gopkg.in/ini.v1"
)

func main() {
	// 加载 INI 文件
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 读取 [database] 节的配置
	dbSection := cfg.Section("database")
	dbHost := dbSection.Key("host").String()
	dbPort, _ := dbSection.Key("port").Int()
	enableSSL, _ := dbSection.Key("enable_ssl").Bool()

	// 读取 [server] 节的配置
	serverSection := cfg.Section("server")
	listenAddr := serverSection.Key("listen").String()
	logLevel := serverSection.Key("log_level").String()

	// 打印配置
	fmt.Printf("数据库地址: %s:%d (SSL: %t)\n", dbHost, dbPort, enableSSL)
	fmt.Printf("服务器监听: %s, 日志级别: %s\n", listenAddr, logLevel)
}
```

------

#### **3. 进阶用法**

##### **(1) 直接映射到结构体**

```go
type DatabaseConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	EnableSSL bool  `ini:"enable_ssl"`
}

type ServerConfig struct {
	Listen    string `ini:"listen"`
	LogLevel  string `ini:"log_level"`
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	// 映射配置到结构体
	dbCfg := new(DatabaseConfig)
	if err := cfg.Section("database").MapTo(dbCfg); err != nil {
		log.Fatal(err)
	}

	serverCfg := new(ServerConfig)
	if err := cfg.Section("server").MapTo(serverCfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("数据库用户: %s\n", dbCfg.Username)
}
```

##### **(2) 设置默认值**

```go
key := cfg.Section("database").Key("max_connections")
maxConn := key.MustInt(100) // 若配置项不存在，默认值 100
```

##### **(3) 写入 INI 文件**

```go
cfg := ini.Empty()

// 创建 [app] 节并写入配置
appSection := cfg.Section("app")
appSection.Key("version").SetValue("1.0.0")
appSection.Key("debug").SetValue("false")

// 保存到文件
err := cfg.SaveTo("app.ini")
if err != nil {
	log.Fatal(err)
}
```

------

#### **4. 常见问题处理**

| **问题**         | **解决方案**                                                 |
| ---------------- | ------------------------------------------------------------ |
| 键名大小写敏感   | `go-ini` 默认不区分大小写，可通过 `cfg = ini.LoadOptions{Insensitive: false}` 设置敏感。 |
| 节不存在时的处理 | 使用 `cfg.Section("section").Key("key").MustXxx()` 避免 panic。 |
| 解析布尔值       | 支持 `true`/`false`、`1`/`0`、`on`/`off` 等多种格式。        |

------

### **5. 对比其他格式**

| **格式** | **优点**                 | **缺点**               | **适用场景**            |
| -------- | ------------------------ | ---------------------- | ----------------------- |
| **INI**  | 简单易读，轻量级         | 不支持复杂数据类型     | 小型项目、分层配置      |
| **JSON** | 结构化数据，类型明确     | 冗余字符多，可读性稍差 | 前后端交互、复杂配置    |
| **YAML** | 可读性极佳，支持复杂结构 | 缩进敏感，易出错       | DevOps、Kubernetes 配置 |

------

### **总结**

- **INI 文件适用场景**：轻量级配置（如数据库连接、服务器参数）。
- **Go 实现方案**：推荐 `go-ini` 库，支持读取、写入和结构体映射。
- **最佳实践**：将配置映射到结构体，提升代码可维护性。

## 数据库主从复制

### 主从复制

主库（Master）负责处理写操作（INSERT/UPDATE/DELETE），从库（Slave）通过复制主库的二进制日志（Binlog）同步数据，只处理读操作（SELECT）。

![](assets/mermaid-diagram-1741351570687.png)

| **优势**     | **说明**                                                     |
| ------------ | ------------------------------------------------------------ |
| **读写分离** | 写操作集中在主库，读操作分散到多个从库，提升并发性能。       |
| **高可用性** | 主库故障时，可从从库快速切换为新主库（需配合哨兵或集群管理工具）。 |
| **数据备份** | 从库天然是主库的实时备份，避免单点故障导致数据丢失。         |
| **负载均衡** | 通过多个从库分担读请求压力，适合读多写少的场景（如电商、社交应用）。 |

#### **(1) 性能优化**

- **场景**：当应用 80% 以上的请求是读操作时（如商品详情页、新闻列表），主从分离可显著降低主库压力。
- **效果**：通过横向扩展从库数量，线性提升读吞吐量。

#### **(2) 高并发支持**

- **问题**：单库连接数有限（如 MySQL 默认最大连接数 151），高并发时易达到瓶颈。
- **解决**：主库专注写操作，从库处理读操作，突破连接数限制。

#### **(3) 数据安全性**

- **容灾恢复**：从库作为实时备份，主库故障时可快速恢复数据。
- **读写隔离**：避免复杂查询（如报表分析）影响核心业务写入性能。

### 主从配置最佳实践

| **实践要点**     | **说明**                                                     |
| ---------------- | ------------------------------------------------------------ |
| **从库数量**     | 建议至少 2 个从库，避免单点故障。                            |
| **读写分离策略** | 写操作强制走主库，读操作根据业务容忍度选择从库。             |
| **健康检查**     | 定期监控从库复制状态（如 `SHOW SLAVE STATUS`），自动剔除异常节点。 |
| **故障切换**     | 使用 Keepalived 或 Consul 实现主库故障自动切换。             |

## AES 对称加密(加密金额)

**AES（Advanced Encryption Standard）** 是一种广泛使用的对称加密算法，特点如下：

| **特性**     | **说明**                                                     |
| ------------ | ------------------------------------------------------------ |
| **对称性**   | 加密与解密使用 **同一密钥**，需确保密钥传输安全。            |
| **分组加密** | 数据按 **128位（16字节）** 分组处理，密钥长度支持 **128/192/256位**。 |
| **工作模式** | 如 ECB、CBC、CTR 等，不同模式安全性差异大（ECB 不安全，推荐 CBC 或 GCM）。 |
| **填充规则** | 数据长度不足分组大小时需填充，常用 PKCS#7（如缺3字节则填充 `0x03` 三次）。 |

## bcrypt 加密

**bcrypt** 是一种专为密码存储设计的加密算法，它通过 **慢哈希（Slow Hashing）** 和 **盐值（Salt）** 机制，有效抵御暴力破解和彩虹表攻击。

### **bcrypt 的核心特性**

| 特性                            | 说明                                                         |
| ------------------------------- | ------------------------------------------------------------ |
| **自适应成本（Adaptive Cost）** | 可通过调整 `cost` 参数控制计算复杂度，对抗硬件性能提升带来的暴力破解风险。 |
| **内置盐值（Built-in Salt）**   | 自动生成随机盐值并与哈希结果合并存储，无需单独管理盐值。     |
| **基于 Blowfish 加密**          | 核心算法基于 Blowfish 的密钥扩展机制，通过多次迭代增强安全性。 |
| **输出格式固定**                | 生成的哈希字符串包含算法标识、cost 值、盐值和哈希结果，便于存储和验证。 |

### **bcrypt 的加密过程**

1. **输入参数**
   - **明文密码**：用户输入的原始密码（如 `"mypassword123"`）。
   - **Cost 因子**：整数（通常为 10~15），决定哈希的计算复杂度（迭代次数 = 2^cost）。
2. **生成随机盐值**
    bcrypt 自动生成一个 **128 位（16 字节）的随机盐值**，确保相同密码每次加密结果不同。
3. **密钥扩展（Key Schedule）**
    基于 Blowfish 算法，使用盐值和 cost 因子进行多轮迭代，生成一个 **扩展密钥**。
4. **加密流程**
   - 使用扩展密钥对固定字符串 `"OrpheanBeholderScryDoubt"` 进行 **64 次 Blowfish 加密**。
   - 最终生成 **192 位（24 字节）的哈希值**。

5. **结果格式**
    最终输出的哈希字符串格式如下（示例）：

```tex
$2a$10$N9qo8uLOickgx2ZMRZoMye3Z7g7z4Z0eP9aD9Jbyd4uReVfM/4mOy
├──┬──┼──┬──┼───────────────────────────┬───────────────────────────┤
 算法标识 cost      盐值（22字符）                 哈希结果（31字符）
```

- **算法标识**：`$2a$` 表示使用 bcrypt 算法。
- **Cost 值**：`10` 表示迭代次数为 2^10 = 1024 次。
- **盐值**：`N9qo8uLOickgx2ZMRZoMye`（Base64 编码，实际为 16 字节）。
- **哈希结果**：`3Z7g7z4Z0eP9aD9Jbyd4uReVfM/4mOy`（Base64 编码，实际为 24 字节）。

### **代码示例（Go 语言实现）**

#### 1. **密码加密**

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    // cost 设为 12（通常推荐值）
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}
```

#### 2. **密码验证**

```go
func CheckPassword(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(hashedPassword),
        []byte(password),
    )
    return err == nil
}
```

### **与其他哈希算法的对比**

| 算法        | 特点                                                       | 适用场景             |
| ----------- | ---------------------------------------------------------- | -------------------- |
| **bcrypt**  | 慢哈希、内置盐值、抗 GPU 加速                              | 密码存储             |
| **SHA-256** | 快速哈希、无盐值，需自行实现迭代和盐值                     | 数据完整性校验       |
| **Argon2**  | 内存硬哈希、抗 GPU/ASIC 加速，需配置并行度、内存大小等参数 | 高安全性密码存储     |
| **PBKDF2**  | 基于 HMAC 的迭代哈希，配置简单但抗硬件加速能力较弱         | 兼容性要求高的旧系统 |

### 彩虹表攻击（**Rainbow Table Attack**）

#### **定义**

彩虹表攻击是一种利用**预计算哈希值数据库**快速破解密码的逆向攻击方法。攻击者通过预先计算大量常见密码的哈希值，并存储在“彩虹表”中，当获取数据库中的密码哈希后，直接查表匹配明文密码。

#### **攻击流程**

1. **预生成哈希表**：
    提前计算常见密码（如 `123456`、`password`）及其组合的哈希值，形成映射关系表。

   ```
   明文密码   →   哈希值
   --------------------------
   123456    →   e10adc3949ba59abbe56...
   password  →   5f4dcc3b5aa765d61d83...
   qwerty    →   d8578edf8458ce06fbc...
   ```

2. **窃取数据库哈希**：
    攻击者获取数据库中的密码哈希值（如 `e10adc3949ba59abbe56...`）。

3. **查表匹配**：
    在彩虹表中查找哈希值对应的明文密码（如 `123456`），实现快速破解。

#### **防御手段**

- 盐值（Salt）：

在哈希前为每个密码添加随机字符串（盐值），使相同密码的哈希结果不同，彩虹表失效。

```
盐值 + 明文密码 → 哈希值
--------------------------
x1&9 + 123456 → 89a7de3f8d12...
kp5@ + 123456 → 6b2c4f1e9d0a...  (相同密码，不同哈希)
```

### 慢哈希

#### **定义**

慢哈希是一种**故意降低哈希计算速度**的技术，通过增加计算时间和资源消耗，大幅提高暴力破解的难度。

#### **快哈希 vs 慢哈希**

| 特性               | 快哈希（如 MD5、SHA-256）    | 慢哈希（如 bcrypt、PBKDF2） |
| ------------------ | ---------------------------- | --------------------------- |
| **设计目标**       | 快速计算数据完整性校验       | 密码存储，增加暴力破解成本  |
| **计算速度**       | 微秒级（百万次/秒）          | 毫秒级（每秒数次）          |
| **适用场景**       | 文件校验、数字签名           | 用户密码存储                |
| **抗暴力破解能力** | 弱（易被 GPU/ASIC 加速破解） | 强（硬件加速难以优化）      |

#### **慢哈希的实现原理**

1. **迭代计算（Key Stretching）**：
    多次重复哈希过程（如迭代 10,000 次），延长单次计算时间。

   ```go
   // 伪代码示例：迭代哈希
   hash = password
   for i := 0; i < 10000; i++ {
       hash = sha256(hash + salt)
   }
   ```

2. **内存密集型算法（如 Argon2）**：
    要求大量内存访问，限制 GPU/ASIC 的并行计算优势。

#### **实际影响**

假设攻击者尝试破解一个密码：

- **快哈希**（1微秒/次）：每秒可尝试 **100万次**。
- **慢哈希**（100毫秒/次）：每秒仅尝试 **10次**。
  即使密码简单（如 6 位数字），破解时间从 **10秒** 延长到 **27小时**。

## JSON Web Token (JWT) 签发与解析

### **JWT 的基本结构**

JWT 由三部分组成，格式为：
 `Header.Payload.Signature`
 用 Base64URL 编码后拼接成字符串，例如：
 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImV4cCI6MTY4MDAwMDAwMH0.7d5b7d5b7d5b7d5b7d5b7d5b7d5b7d5b`

| 部分          | 内容                                                         | 作用                   |
| ------------- | ------------------------------------------------------------ | ---------------------- |
| **Header**    | 声明算法（如 HMAC SHA256）和 token 类型（JWT）               | 指示如何验证签名       |
| **Payload**   | 包含用户定义的声明（如用户ID、过期时间）和标准声明（如 `exp`, `iss`） | 携带业务相关数据       |
| **Signature** | 对前两部分的签名，使用密钥和指定算法生成                     | 防篡改，确保数据完整性 |

### **JWT 签发（Signing）过程**

1. **定义 Header 和 Payload**

```go
// Header
{
  "alg": "HS256",  // 使用 HMAC-SHA256 算法
  "typ": "JWT"
}

// Payload
{
  "user_id": 123,
  "exp": 1680000000  // 过期时间（Unix 时间戳）
}

```

2. **Base64URL 编码**
    将 Header 和 Payload 分别进行 Base64URL 编码，得到：

- `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`
- `eyJ1c2VyX2lkIjoxMjMsImV4cCI6MTY4MDAwMDAwMH0`

3. **生成签名**
    使用密钥（如 `secret_key`）和算法对 `Base64URL(Header) + "." + Base64URL(Payload)` 签名：

```go
signature = HMAC-SHA256(
  key = "secret_key",
  data = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImV4cCI6MTY4MDAwMDAwMH0"
)
```

4. **拼接完整 Token**

将三部分用 `.` 连接：

```text
JWT = Base64URL(Header) + "." + Base64URL(Payload) + "." + Base64URL(Signature)
```

### **JWT 解析（Parsing）过程**

1. **拆分 Token**
    按 `.` 分割 Token，得到三部分：
   - `Header` → `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`
   - `Payload` → `eyJ1c2VyX2lkIjoxMjMsImV4cCI6MTY4MDAwMDAwMH0`
   - `Signature` → `7d5b7d5b7d5b7d5b7d5b7d5b7d5b7d5b`
2. **验证签名**
    使用相同密钥和算法重新计算签名，与 Token 中的签名比对，确保未被篡改。
3. **解码 Payload**
    将 Base64URL 解码后的 Payload 转换为 JSON 对象，提取声明数据。
4. **验证声明**
    检查标准声明（如 `exp` 是否过期，`iss` 是否合法）和业务声明（如 `user_id` 是否存在）。

#### **Golang 中的 JWT 实现示例**

1. **安装依赖**

```bash
go get github.com/golang-jwt/jwt/v5
```

2. **签发 Token**

```go
package main

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// 定义自定义声明（Payload）
type CustomClaims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims  // 嵌入标准声明（exp, iss 等）
}

func GenerateToken(userID int, secretKey string) (string, error) {
    // 1. 定义过期时间（1小时后）
    expirationTime := time.Now().Add(1 * time.Hour)

    // 2. 创建声明（Claims）
    claims := CustomClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            Issuer:    "my_app",
        },
    }

    // 3. 生成 Token 对象
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // 4. 使用密钥签名生成完整 Token
    signedToken, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }

    return signedToken, nil
}

func main() {
    secretKey := "your_secure_secret_key"  // 生产环境应从安全配置读取
    token, err := GenerateToken(123, secretKey)
    if err != nil {
        panic(err)
    }
    fmt.Println("Generated Token:", token)
}

```

3. **解析并验证 token**

```go
func ParseToken(tokenString, secretKey string) (*CustomClaims, error) {
    // 1. 解析 Token，验证签名和格式
    token, err := jwt.ParseWithClaims(
        tokenString,
        &CustomClaims{},
        func(token *jwt.Token) (interface{}, error) {
            // 验证签名算法
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(secretKey), nil
        },
    )

    if err != nil {
        return nil, err
    }

    // 2. 提取自定义声明
    if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
        return claims, nil
    } else {
        return nil, fmt.Errorf("invalid token")
    }
}

func main() {
    token := "your.jwt.token.here"
    secretKey := "your_secure_secret_key"

    claims, err := ParseToken(token, secretKey)
    if err != nil {
        fmt.Println("Token validation failed:", err)
        return
    }

    fmt.Printf("Valid claims: UserID=%d, ExpiresAt=%v\n", 
        claims.UserID, 
        claims.ExpiresAt.Time.Format(time.RFC3339),
    )
}

```

### **关键安全注意事项**

| 注意事项           | 说明                                                         |
| ------------------ | ------------------------------------------------------------ |
| **密钥管理**       | 使用高强度随机密钥（如 32 字节以上），避免硬编码，推荐从环境变量读取 |
| **签名算法**       | 优先选择 `HS256`（HMAC）或 `RS256`（RSA），禁用 `none` 算法  |
| **Token 过期时间** | 设置合理的短有效期（如 1 小时），结合 Refresh Token 机制更新 |
| **敏感数据存储**   | JWT Payload 默认不加密，切勿存放密码、密钥等敏感信息         |
| **HTTPS 传输**     | 始终通过 HTTPS 传输 Token，防止中间人窃听                    |

## VO view objective

**VO（View Object，视图对象）** 是一种设计模式中的概念，主要用于 **解耦后端数据模型与前端视图展示**，确保传递给前端的数据结构是经过定制、过滤或增强的，而非直接暴露数据库模型。它的核心目标是 **控制数据暴露的粒度和格式**，提高安全性和灵活性。

### **VO 的核心作用**

| 作用             | 说明                                   | 示例场景                                |
| ---------------- | -------------------------------------- | --------------------------------------- |
| **数据裁剪**     | 隐藏敏感字段（如密码、内部状态码）     | 移除 `User` 模型中的 `Password` 字段    |
| **数据格式化**   | 转换日期、拼接URL、枚举值转文本等      | 将 `CreatedAt` 的时间戳转为前端易读格式 |
| **数据聚合**     | 合并多个模型的数据，或添加额外计算字段 | 组合用户信息 + 统计信息生成用户详情VO   |
| **适配前端需求** | 根据视图需求定制字段名或结构           | 将字段名从 `user_name` 改为 `username`  |

### **对比常见数据对象**

| 类型          | 定位                                  | 生命周期             | 典型场景                  |
| ------------- | ------------------------------------- | -------------------- | ------------------------- |
| **VO**        | 为前端视图定制数据                    | 后端到前端的传输过程 | 接口响应数据              |
| **DTO**       | 跨层数据传输（如Service到Controller） | 内部系统层级间传递   | 微服务间通信              |
| **DAO/Model** | 直接映射数据库结构                    | 数据库操作到业务逻辑 | ORM 模型、数据库表对应类  |
| **POJO**      | 纯数据容器，无业务逻辑                | 通用数据存储         | 简单的Java/Kotlin数据对象 |

### **代码示例解析**

用户提供的代码中，`UserVO` 是典型的视图对象实现：

```go
type UserVO struct {
    ID        uint   `json:"id"`
    UserName  string `json:"user_name"`
    // 省略敏感字段（如 Password）
    Avatar    string `json:"avatar"`    // 拼接完整 URL
    CreatedAt int64  `json:"created_at"` // 时间戳格式化
}
func BuildUser(user *model.User) *UserVO {
    return &UserVO{
        // ...
        Avatar: conf.Host + conf.HttpPort + conf.AvatarPath + user.Avatar,
        CreatedAt: user.CreatedAt.Unix(),
    }
}
```

#### **关键实现点**

1. **字段过滤**
    不暴露 `model.User` 中的敏感字段（如密码、内部状态码）。
2. **数据转换**
   - **Avatar**：将数据库存储的相对路径（如 `user1/avatar.jpg`）拼接为完整 URL（如 `http://localhost:8080/static/imgs/avatar/user1/avatar.jpg`）。
   - **CreatedAt**：将 `time.Time` 类型转为 Unix 时间戳，方便前端处理。
3. **结构适配**
    字段命名（如 `user_name`）和类型（如 `int64`）按前端需求定义，而非严格遵循数据库模型。

通过 VO 模式，代码的可维护性和安全性得到显著提升，是分层架构中不可或缺的一环。

## 中间件

