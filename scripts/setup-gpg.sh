#!/bin/bash
# GPG 金鑰生成和導出腳本
# 用於 Terraform Provider 發布

set -e

echo "=== Terraform Provider GPG 金鑰設置 ==="
echo ""

# 檢查是否已安裝 GPG
if ! command -v gpg &> /dev/null; then
    echo "錯誤: 未找到 GPG。請先安裝 GPG。"
    echo "macOS: brew install gnupg"
    echo "Linux: sudo apt-get install gnupg 或 sudo yum install gnupg"
    exit 1
fi

echo "步驟 1: 生成 GPG 金鑰"
echo "----------------------------------------"
echo "請按照提示輸入以下信息："
echo "  - 金鑰類型: 選擇 (1) RSA and RSA (default)"
echo "  - 金鑰長度: 輸入 4096"
echo "  - 過期時間: 建議輸入 2y (2年) 或 3y (3年)"
echo "  - 姓名: 輸入您的姓名"
echo "  - 電子郵件: 輸入您的電子郵件地址"
echo "  - 密碼短語: 設定一個安全的密碼短語（請記住這個密碼！）"
echo ""
read -p "準備好後按 Enter 繼續生成 GPG 金鑰..."

gpg --full-generate-key

echo ""
echo "步驟 2: 查看生成的金鑰"
echo "----------------------------------------"
gpg --list-secret-keys --keyid-format=long

echo ""
echo "步驟 3: 導出金鑰"
echo "----------------------------------------"
read -p "請輸入您的 GPG 金鑰 ID (長格式，例如: 3AA5C34371567BD2): " KEY_ID

if [ -z "$KEY_ID" ]; then
    echo "錯誤: 未輸入金鑰 ID"
    exit 1
fi

# 創建輸出目錄
OUTPUT_DIR="gpg-keys"
mkdir -p "$OUTPUT_DIR"

# 導出公鑰
echo "導出公鑰到 $OUTPUT_DIR/public-key.gpg..."
gpg --armor --export "$KEY_ID" > "$OUTPUT_DIR/public-key.gpg"

# 導出私鑰
echo "導出私鑰到 $OUTPUT_DIR/private-key.gpg..."
gpg --armor --export-secret-keys "$KEY_ID" > "$OUTPUT_DIR/private-key.gpg"

echo ""
echo "✅ 金鑰導出完成！"
echo ""
echo "下一步操作："
echo "1. 查看公鑰內容（用於 Terraform Registry）:"
echo "   cat $OUTPUT_DIR/public-key.gpg"
echo ""
echo "2. 查看私鑰內容（用於 GitHub Secrets）:"
echo "   cat $OUTPUT_DIR/private-key.gpg"
echo ""
echo "3. 在 GitHub 倉庫設置中添加以下 Secrets："
echo "   - GPG_PRIVATE_KEY: 複製 $OUTPUT_DIR/private-key.gpg 的完整內容"
echo "   - PASSPHRASE: 輸入您剛才設定的密碼短語"
echo ""
echo "4. 在 Terraform Registry 註冊時上傳公鑰："
echo "   - 訪問 https://registry.terraform.io/publish/provider"
echo "   - 上傳 $OUTPUT_DIR/public-key.gpg 的內容"
echo ""
echo "⚠️  安全提示："
echo "   - 妥善保管 private-key.gpg 文件"
echo "   - 不要將私鑰提交到 Git 倉庫"
echo "   - 完成設置後可以刪除本地私鑰文件（但請確保已備份）"
