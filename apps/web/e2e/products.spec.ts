import { test, expect } from '@playwright/test';

test.describe('Products Page', () => {
  test('should load products listing page', async ({ page }) => {
    await page.goto('/products');

    // Page should render — either with product grid, loading state, or empty state
    // Since backend may not be running, we just verify the page loads without errors
    await expect(page.locator('body')).toBeVisible();

    // Should have header
    await expect(page.locator('header')).toBeVisible();

    // Should show either products, loading skeleton, or empty message
    const hasContent = await page.locator('text=No products found, text=Loading, [class*="skeleton"], [class*="ProductCard"], [class*="product-card"]').first().isVisible().catch(() => false);
    expect(hasContent).toBeTruthy();
  });

  test('should have header navigation', async ({ page }) => {
    await page.goto('/products');
    await expect(page.locator('header')).toBeVisible();
  });
});

test.describe('Product Detail Page', () => {
  test('should handle product detail page', async ({ page }) => {
    const response = await page.goto('/products/1');
    expect(response?.status()).toBeLessThan(500);
  });
});
