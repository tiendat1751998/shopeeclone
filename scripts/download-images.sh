#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

PUBLIC_IMAGES="$PROJECT_DIR/apps/web/public/images"

mkdir -p "$PUBLIC_IMAGES/products" "$PUBLIC_IMAGES/categories" "$PUBLIC_IMAGES/shops" "$PUBLIC_IMAGES/uploads" "$PUBLIC_IMAGES/avatars"

download() {
  local url="$1"
  local path="$2"
  if [ ! -f "$path" ]; then
    echo "Downloading $path ..."
    curl -fsSL -o "$path" "$url" || echo "WARNING: failed to download $url"
  else
    echo "Skipping $path (already exists)"
  fi
}

for i in $(seq 1 10); do
  download "https://picsum.photos/seed/product${i}/400/400" "$PUBLIC_IMAGES/products/product-${i}.jpg"
done

for angle in 2 3 4; do
  download "https://picsum.photos/seed/product-${angle}-angle/400/400" "$PUBLIC_IMAGES/products/product-2-${angle}.jpg"
done

for i in $(seq 1 10); do
  download "https://picsum.photos/seed/category${i}/200/200" "$PUBLIC_IMAGES/categories/category-${i}.jpg"
done

download "https://picsum.photos/seed/shop/200/200" "$PUBLIC_IMAGES/shops/default-shop-avatar.png"
download "https://picsum.photos/seed/user/200/200" "$PUBLIC_IMAGES/avatars/default-avatar.png"

chmod +x "$0"
echo "Done!"
