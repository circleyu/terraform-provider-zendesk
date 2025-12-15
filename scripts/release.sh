#!/bin/bash
# Terraform Provider 發布腳本

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== Terraform Provider 發布腳本 ==="
echo ""

# 檢查是否在 git 倉庫中
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}錯誤: 當前目錄不是 Git 倉庫${NC}"
    exit 1
fi

# 檢查是否有未提交的更改
if ! git diff-index --quiet HEAD --; then
    echo -e "${YELLOW}警告: 檢測到未提交的更改${NC}"
    read -p "是否要先提交這些更改？(y/n): " commit_changes
    if [ "$commit_changes" = "y" ]; then
        git add .
        read -p "請輸入提交訊息: " commit_message
        git commit -m "$commit_message"
    else
        echo -e "${RED}請先提交或暫存所有更改後再繼續${NC}"
        exit 1
    fi
fi

# 獲取版本號
read -p "請輸入版本號 (例如: 1.0.0，將自動添加 v 前綴): " version
if [ -z "$version" ]; then
    echo -e "${RED}錯誤: 未輸入版本號${NC}"
    exit 1
fi

# 驗證版本號格式（簡單檢查）
if ! [[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}錯誤: 版本號格式不正確，應為 X.Y.Z (例如: 1.0.0)${NC}"
    exit 1
fi

tag="v$version"

# 檢查標籤是否已存在
if git rev-parse "$tag" >/dev/null 2>&1; then
    echo -e "${RED}錯誤: 標籤 $tag 已存在${NC}"
    exit 1
fi

# 確認
echo ""
echo "將要執行以下操作："
echo "  - 創建標籤: $tag"
echo "  - 推送標籤到遠程倉庫"
echo ""
read -p "確認繼續？(y/n): " confirm

if [ "$confirm" != "y" ]; then
    echo "已取消"
    exit 0
fi

# 創建標籤
echo ""
echo "創建標籤 $tag..."
git tag -a "$tag" -m "Release $tag"

# 獲取遠程倉庫名稱
remote=$(git remote | head -n 1)
if [ -z "$remote" ]; then
    echo -e "${YELLOW}警告: 未找到遠程倉庫，請手動推送標籤${NC}"
    echo "使用命令: git push origin $tag"
    exit 0
fi

# 推送標籤
echo "推送標籤到 $remote..."
read -p "是否現在推送標籤？(y/n): " push_confirm

if [ "$push_confirm" = "y" ]; then
    git push "$remote" "$tag"
    echo ""
    echo -e "${GREEN}✅ 標籤已推送！${NC}"
    echo ""
    echo "GitHub Actions 將自動觸發發布流程。"
    echo "請前往 GitHub 倉庫的 Actions 頁面查看發布進度。"
    echo ""
    echo "發布完成後，可以在以下位置查看："
    echo "  - GitHub Releases: https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/releases"
    echo "  - Terraform Registry: https://registry.terraform.io/providers/circleyu/zendesk"
else
    echo ""
    echo "標籤已創建但未推送。稍後可以使用以下命令推送："
    echo "  git push $remote $tag"
fi
