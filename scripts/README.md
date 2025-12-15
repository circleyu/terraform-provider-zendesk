# Terraform Provider 發布腳本

本目錄包含用於發布 Terraform Provider 到 Terraform Registry 的輔助腳本。

## 文件說明

- `setup-gpg.sh` - GPG 金鑰生成和導出腳本
- `release.sh` - 版本發布腳本

## 使用步驟

### 1. 生成 GPG 金鑰

運行以下命令生成 GPG 簽名金鑰：

```bash
./scripts/setup-gpg.sh
```

腳本會引導您完成：
- 生成 4096 位 RSA GPG 金鑰對
- 導出公鑰和私鑰到 `gpg-keys/` 目錄

**重要提示：**
- 請記住您設定的密碼短語（passphrase），稍後需要添加到 GitHub Secrets
- 妥善保管私鑰文件，不要提交到 Git 倉庫

### 2. 配置 GitHub Secrets

1. 前往您的 GitHub 倉庫：`Settings` → `Secrets and variables` → `Actions`
2. 添加以下 Secrets：

   **GPG_PRIVATE_KEY**
   - 複製 `gpg-keys/private-key.gpg` 文件的完整內容
   - 包括 `-----BEGIN PGP PRIVATE KEY BLOCK-----` 和 `-----END PGP PRIVATE KEY BLOCK-----`

   **PASSPHRASE**
   - 輸入您在步驟 1 中設定的 GPG 密碼短語

### 3. 在 Terraform Registry 註冊 Provider

1. 訪問 [Terraform Registry](https://registry.terraform.io/)
2. 使用 GitHub 帳號登入
3. 前往 [Publish Provider](https://registry.terraform.io/publish/provider)
4. 選擇您的 GitHub 倉庫
5. 上傳 GPG 公鑰：
   - 複製 `gpg-keys/public-key.gpg` 文件的內容
   - 貼上到 Terraform Registry 的 GPG 公鑰欄位

### 4. 發布版本

在完成所有配置後，使用發布腳本創建新版本：

```bash
./scripts/release.sh
```

腳本會：
- 檢查未提交的更改
- 提示您輸入版本號（例如：1.0.0）
- 創建並推送版本標籤（例如：v1.0.0）
- 觸發 GitHub Actions 自動發布流程

## 版本號格式

版本號必須遵循語義化版本（SemVer）：
- 主版本號.次版本號.修訂號（例如：1.0.0）
- 腳本會自動添加 `v` 前綴（例如：v1.0.0）

## 自動發布流程

當您推送版本標籤後，GitHub Actions 會自動：

1. 觸發 `.github/workflows/release.yml`
2. GoReleaser 會：
   - 編譯多平台二進制文件（Linux, Windows, macOS, FreeBSD）
   - 生成 SHA256 校驗和文件
   - 使用 GPG 簽名校驗和文件
   - 創建 GitHub Release
   - 包含 `terraform-registry-manifest.json`

## 驗證發布

發布完成後，請檢查：

1. **GitHub Release**
   - 前往倉庫的 `Releases` 頁面
   - 確認包含所有必要文件

2. **Terraform Registry**
   - 訪問 https://registry.terraform.io/providers/circleyu/zendesk
   - 新版本應該在幾分鐘內自動出現

## 故障排除

### GPG 金鑰問題

如果遇到 GPG 相關問題：
- 確保已安裝 GPG：`brew install gnupg` (macOS) 或 `sudo apt-get install gnupg` (Linux)
- 檢查金鑰是否存在：`gpg --list-secret-keys --keyid-format=long`

### GitHub Actions 失敗

如果發布流程失敗：
- 檢查 GitHub Secrets 是否正確配置
- 確認 GPG 私鑰和密碼短語正確
- 查看 GitHub Actions 日誌獲取詳細錯誤信息

### Terraform Registry 未顯示版本

- 確認 GitHub Release 已成功創建
- 檢查 Release 中是否包含所有必要文件
- 等待幾分鐘讓 Registry 同步
- 確認 GPG 公鑰已在 Registry 中正確配置
