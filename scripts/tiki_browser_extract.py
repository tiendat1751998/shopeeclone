#!/usr/bin/env python3
"""
Tiki.vn Browser-Aided Crawler
Outputs JavaScript code to be evaluated via browser automation.
The JS extracts product data from Tiki category pages.
Output: JSON file with product data ready for tiki_processor.py insert
"""

import json

# JavaScript to inject into browser for extracting product data from Tiki listing pages
TIKI_EXTRACT_JS = """
(function() {
    var products = [];
    
    // Method 1: Extract from product card links (full cards with images, ratings, etc.)
    var cardLinks = document.querySelectorAll('a[class*="sc-4eae099-0"]');
    if (cardLinks.length > 0) {
        cardLinks.forEach(function(link) {
            try {
                var href = link.getAttribute('href') || '';
                var productIdMatch = href.match(/-p(\\d+)\\.html/);
                if (!productIdMatch) return;
                var productId = productIdMatch[1];
                
                // Image
                var imgEl = link.querySelector('img[class*="sc-4eae099-1"]');
                var imageUrl = imgEl ? (imgEl.getAttribute('src') || '') : '';
                
                // Name
                var nameEl = link.querySelector('div[class*="sc-4eae099-7"]');
                var name = nameEl ? nameEl.textContent.trim() : '';
                if (!name) return;
                
                // Seller
                var sellerEl = link.querySelector('div[class*="sc-4eae099-6"]');
                var sellerName = sellerEl ? sellerEl.textContent.trim() : null;
                var sellerAvatarEl = link.querySelector('img[class*="sc-4eae099-5"]');
                var sellerAvatarUrl = sellerAvatarEl ? sellerAvatarEl.getAttribute('src') : null;
                var isTikiTrading = sellerName === 'Tiki Trading';
                
                // Sponsored
                var sponsoredEl = link.querySelector('div[class*="sc-4eae099-2"]');
                var isSponsored = sponsoredEl && sponsoredEl.textContent.includes('Tài trợ');
                
                // Rating & sold
                var ratingSoldEl = link.querySelector('div[class*="sc-f5f45ccf-0"]');
                var ratingAverage = null;
                var soldCount = null;
                var soldText = null;
                if (ratingSoldEl) {
                    var starsEl = ratingSoldEl.querySelector('p[class*="sc-eed848ce-0"]');
                    if (starsEl) {
                        var starSvgs = starsEl.querySelectorAll('svg.star-icon');
                        ratingAverage = starSvgs.length;
                    }
                    var divs = ratingSoldEl.querySelectorAll('div');
                    for (var i = 0; i < divs.length; i++) {
                        var txt = divs[i].textContent.trim();
                        if (txt.indexOf('Đã bán') >= 0) {
                            var sm = txt.match(/Đã bán\\s+(.+)/);
                            if (sm) {
                                soldText = sm[1].trim();
                                var ks = soldText.match(/([\\d.]+)k/i);
                                if (ks) {
                                    soldCount = parseInt(parseFloat(ks[1]) * 1000);
                                } else {
                                    soldCount = parseInt(soldText.replace(/[.,]/g, ''));
                                }
                            }
                            break;
                        }
                    }
                }
                
                // Price
                var priceEl = link.querySelector('div[class*="sc-4eae099-9"]');
                var priceText = priceEl ? priceEl.textContent.trim() : '';
                var price = null;
                var priceMatch = priceText.match(/([\\d.]+)/);
                if (priceMatch) {
                    price = parseInt(priceMatch[1].replace(/\\./g, ''));
                }
                
                // Discount
                var discountEl = link.querySelector('div[class*="sc-4eae099-10"]');
                var discountText = discountEl ? discountEl.textContent.trim() : '';
                var discountMatch = discountText.match(/-(\\d+)%/);
                var discountPercent = discountMatch ? parseInt(discountMatch[1]) : null;
                
                // Calculate original price
                var originalPrice = null;
                if (discountPercent && price && discountPercent > 0) {
                    originalPrice = Math.round(price / (1 - discountPercent / 100));
                }
                
                products.push({
                    tiki_product_id: productId,
                    name: name,
                    url: href.charAt(0) === '/' ? 'https://tiki.vn' + href : href,
                    image_url: imageUrl,
                    thumbnail_url: imageUrl,
                    price: price || 0,
                    original_price: originalPrice,
                    discount_percent: discountPercent,
                    rating_average: ratingAverage ? ratingAverage.toFixed(2) : null,
                    sold_count: isNaN(soldCount) ? null : soldCount,
                    quantity_sold_text: soldText,
                    seller_name: sellerName,
                    seller_avatar_url: sellerAvatarUrl,
                    is_tiki_trading: isTikiTrading,
                    is_official: isTikiTrading,
                    is_sponsored: isSponsored,
                    status: 'active'
                });
            } catch(e) {}
        });
    }
    
    // Method 2: Extract from simple product links (fallback for pages without full cards)
    if (products.length === 0) {
        var allLinks = document.querySelectorAll('a[href*="-p"]');
        allLinks.forEach(function(link) {
            var href = link.getAttribute('href') || '';
            if (!href.match(/-p\\d+\\.html/)) return;
            var matched = href.match(/-p(\\d+)\\.html/);
            var productId = matched[1];
            var text = link.textContent || '';
            var priceMatch = text.match(/([\\d.]+)\\s*[₫đ]/);
            var price = priceMatch ? parseInt(priceMatch[1].replace(/\\./g, '')) : 0;
            var title = '';
            if (priceMatch) {
                title = text.substring(0, text.indexOf(priceMatch[0])).trim();
            } else {
                title = text.substring(0, 200).trim();
            }
            // Clean up title
            title = title.replace(/Tài trợ|Tiki Trading|Xem thêm/g, '').trim();
            if (title.length < 5) return;
            var isSponsored = text.includes('Tài trợ');
            var isTikiTrading = text.includes('Tiki Trading');
            products.push({
                tiki_product_id: productId,
                name: title.substring(0, 500),
                url: href.charAt(0) === '/' ? 'https://tiki.vn' + href : href,
                image_url: '',
                thumbnail_url: '',
                price: price,
                original_price: null,
                discount_percent: null,
                rating_average: null,
                sold_count: null,
                seller_name: isTikiTrading ? 'Tiki Trading' : null,
                is_tiki_trading: isTikiTrading,
                is_official: isTikiTrading,
                is_sponsored: isSponsored,
                status: 'active'
            });
        });
    }
    
    // Deduplicate by tiki_product_id
    var seen = {};
    var unique = [];
    products.forEach(function(p) {
        if (!seen[p.tiki_product_id]) {
            seen[p.tiki_product_id] = true;
            unique.push(p);
        }
    });
    
    return JSON.stringify({products: unique, count: unique.length, url: window.location.href, scraped_at: new Date().toISOString()});
})()
"""

# Tiki categories to crawl
TIKI_CATEGORIES = [
    {"slug": "dien-thoai-may-tinh-bang", "id": "1789", "name": "Điện Thoại - Máy Tính Bảng"},
    {"slug": "laptop-may-vi-tinh-linh-kien", "id": "1846", "name": "Laptop - Máy Vi Tính - Linh Kiện"},
    {"slug": "dien-gia-dung", "id": "1882", "name": "Điện Gia Dụng"},
    {"slug": "nha-cua-doi-song", "id": "1883", "name": "Nhà Cửa - Đời Sống"},
    {"slug": "lam-dep-suc-khoe", "id": "1520", "name": "Làm Đẹp - Sức Khỏe"},
    {"slug": "me-be", "id": "2549", "name": "Mẹ & Bé"},
    {"slug": "do-choi-me-be", "id": "2549", "name": "Đồ Chơi - Mẹ & Bé"},
    {"slug": "thoi-trang-nu", "id": "915", "name": "Thời Trang Nữ"},
    {"slug": "thoi-trang-nam", "id": "931", "name": "Thời Trang Nam"},
    {"slug": "giay-dep-nam", "id": "1686", "name": "Giày - Dép Nam"},
    {"slug": "giay-dep-nu", "id": "1703", "name": "Giày - Dép Nữ"},
    {"slug": "dien-tu-dien-lanh", "id": "4221", "name": "Điện Tử - Điện Lạnh"},
    {"slug": "nha-sach-tiki", "id": "8322", "name": "Nhà Sách Tiki"},
    {"slug": "the-thao-da-ngoai", "id": "1975", "name": "Thể Thao - Dã Ngoại"},
    {"slug": "o-to-xe-may-xe-dap", "id": "8594", "name": "Ô Tô - Xe Máy - Xe Đạp"},
    {"slug": "bach-hoa-online", "id": "4384", "name": "Bách Hóa Online"},
    {"slug": "dong-ho-va-trang-suc", "id": "27497", "name": "Đồng Hồ và Trang Sức"},
    {"slug": "balo-va-vali", "id": "6000", "name": "Balo và Vali"},
    {"slug": "phu-kien-thoi-trang", "id": "27498", "name": "Phụ Kiện Thời Trang"},
    {"slug": "tui-thoi-trang-nu", "id": "914", "name": "Túi Thời Trang Nữ"},
]


def get_tiki_url(category_slug, tiki_cat_id, page=1):
    """Generate Tiki category URL."""
    return f"https://tiki.vn/{category_slug}/c{tiki_cat_id}?page={page}"


def save_json(data, filepath):
    """Save data to JSON file."""
    with open(filepath, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    print(f"Saved {len(data.get('products', []))} products to {filepath}")


if __name__ == "__main__":
    print("=== Tiki.vn Crawler - JavaScript Extraction Code ===\n")
    print("This script provides the JavaScript code to extract Tiki product data.")
    print("Use the browser automation tool to:\n")
    print("1. Navigate to a Tiki category page")
    print("2. Evaluate the JS code below to extract products")
    print("3. Save the JSON result")
    print("4. Run tiki_processor.py to insert into database\n")
    print("JS Extraction Code:")
    print("-" * 60)
    print(TIKI_EXTRACT_JS)
    print("-" * 60)
    print(f"\nCategories configured: {len(TIKI_CATEGORIES)}")
    for cat in TIKI_CATEGORIES:
        url = get_tiki_url(cat["slug"], cat["id"])
        print(f"  - {cat['name']}: {url}")
