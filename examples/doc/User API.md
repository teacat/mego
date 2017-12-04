# 會員系統

與會員、使用者有關的 API。

## 建立會員 [createUser]

建立一個新的使用者並在之後可透過註冊時所使用的帳號與密碼進行登入。

+ 請求

    + 欄位
        + username (string!) - 帳號。
        + password (string!) - 密碼。
        + email (string!) - 電子郵件地址。
        + gender (int!)

            使用者的性別採用編號以節省資料空間並在未來有更大的彈性進行變動。0 為男性、1 為女性。

    + 內容

            {
                "username": "YamiOdymel",
                "password": "yami123456",
                "email": "yamiodymel@yami.io",
                "gender": 1
            }

+ 回應

    + StatusOK

    + StatusExists

        + 欄位
            + exists ([]string!) - 已被使用的欄位名稱。

        + 內容

                {
                    "exists": ["username"]
                }

    + StatusInvalid

        + 欄位
            + fields ([]string!) - 格式不正確的欄位名稱。
            + messages ([]string!) - 提示訊息。

        + 內容

                {
                    "fields": ["gender", "email"],
                    "messages": [
                        "`gender` 僅能為數字 0 或 1。",
                        "`email` 長度僅能在 8 到 128 之間。"
                    ]
                }


