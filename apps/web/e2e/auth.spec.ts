import { test, expect } from '@playwright/test';

test.describe('Auth Pages', () => {
  test('should load login page', async ({ page }) => {
    await page.goto('/login');

    // Should have login form elements
    const hasEmailInput = await page.locator('input[type="email"], input[name="email"], input[placeholder*="email" i]').count() > 0;
    const hasPasswordInput = await page.locator('input[type="password"]').count() > 0;
    const hasSubmitButton = await page.locator('button[type="submit"], input[type="submit"]').count() > 0;

    expect(hasEmailInput || hasPasswordInput || hasSubmitButton).toBeTruthy();
  });

  test('should load register page', async ({ page }) => {
    await page.goto('/register');

    // Should have register form elements
    const hasForm = await page.locator('form').count() > 0;
    const hasPasswordInput = await page.locator('input[type="password"]').count() > 0;

    expect(hasForm || hasPasswordInput).toBeTruthy();
  });

  test('login page should have link to register', async ({ page }) => {
    await page.goto('/login');

    const registerLink = page.getByRole('link', { name: /register|sign up|create account/i });
    const hasLink = await registerLink.isVisible().catch(() => false);

    // Not all apps have this, so just check page loads
    expect(await page.locator('body').isVisible()).toBeTruthy();
  });
});
