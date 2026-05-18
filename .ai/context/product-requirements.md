# Product Requirements Document (PRD)

## 1. High-Scale Home Feed & Catalog Experience
- **Home Feed Components**:
  - Banner slider displaying active marketing campaigns.
  - Interactive grid displaying multi-tier Category tree.
  - **Flash Sale Carousel**: Timer countdown countdown, progress bar showing percentage of stock sold out.
  - **Infinite Scroll Product Grid**: Optimized page loading using client-side pre-fetching.
- **Product details screen**:
  - Real-time variation matrix selectors (e.g. size/color buttons deactivate dynamically when variation stock equals 0).
  - Shipping cost preview based on coordinate distances.

## 2. Shopping Cart & Multi-Shop Checkout
- **Cart Grouping**: Cart list must group products based on `SellerID` (shop).
- **Checkout Form**:
  - Voucher Selection modal showing eligible vouchers.
  - Payment Gateways selection (VNPay, MoMo, Cod).
  - Detailed pricing receipt displaying base totals, coupon deductions, shipping, and final payments.
