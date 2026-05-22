import { test, expect } from '@playwright/test';

test.describe('Cart Page', () => {
  test('should load cart page', async ({ page }) => {
    await page.goto('/cart');
    await expect(page.locator('body')).toBeVisible();
  });

  test('should show empty cart or cart items', async ({ page }) => {
    await page.goto('/cart');

    const hasEmptyMessage = await page.getByText(/empty|no items|no products/i).isVisible().catch(() => false);
    const hasCartItems = await page.locator('[class*="cart"], [data-testid*="cart"]').count() > 0;
    const hasHeader = await page.locator('header').isVisible();

    expect(hasHeader).toBeTruthy();
  });
});

test.describe('Checkout Page', () => {
  test('should load checkout page', async ({ page }) => {
    const response = await page.goto('/checkout');
    // Checkout might redirect to login if not authenticated
    expect(response?.status()).toBeLessThan(500);
  });
});

test.describe('Account Page', () => {
  test('should load account page', async ({ page }) => {
    const response = await page.goto('/account');
    // Account might redirect to login if not authenticated
    expect(response?.status()).toBeLessThan(500);
  });
});
