# データ準備
```
CREATE TABLE users (
    id SERIAL PRIMARY KEY,             -- ユーザーID（自動採番）
    username VARCHAR(50) NOT NULL,     -- ユーザー名
    email VARCHAR(100) NOT NULL UNIQUE,-- メールアドレス（一意）
    password_hash VARCHAR(255) NOT NULL,-- ハッシュ化されたパスワード
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 登録日時
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 更新日時
);
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `email` (`email`)
)
INSERT INTO users (username, email, password_hash)
VALUES 
  ('alice', 'alice@example.com', 'hashed_password_1'),
  ('bob', 'bob@example.com', 'hashed_password_2'),
  ('charlie', 'charlie@example.com', 'hashed_password_3');
```

# sqlboiler導入・自動生成
```
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
go install github.com/volatiletech/sqlboiler/v4@latest
go get github.com/volatiletech/sqlboiler/v4@latest
go get github.com/volatiletech/null/v8@latest
sqlboiler psql --config sqlboiler_pg.toml --add-global-variants
sqlboiler mysql --config sqlboiler_mysql.toml --add-global-variants
```

# ビルドタグによる接続DBの切替
go run -tags mysql .
go run -tags pg .
