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
INSERT INTO users (username, email, password_hash)
VALUES 
  ('alice', 'alice@example.com', 'hashed_password_1'),
  ('bob', 'bob@example.com', 'hashed_password_2'),
  ('charlie', 'charlie@example.com', 'hashed_password_3');
```

# sqlboiler導入・自動生成
```
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
go install github.com/volatiletech/sqlboiler/v4@latest
go get github.com/volatiletech/sqlboiler/v4@latest
go get github.com/volatiletech/null/v8@latest
sqlboiler psql --add-global-variants
```

