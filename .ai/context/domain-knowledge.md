# Vietnam E-Commerce Market Domain Knowledge

## 1. Local Delivery Providers & Coordinate Calculations
- Standard local logistics services integrated:
  - **GHN (Giao Hàng Nhanh)**
  - **GHTK (Giao Hàng Tiết Kiệm)**
  - **Viettel Post**
- Shipping pricing requires passing parcel dimensions, packaging weights, origin warehouse coordinates, and destination coordinates to carrier endpoints.

## 2. VietQR EMVCo Standards (Dynamic instant Bank Transfer)
The payment gateway generates VietQR codes dynamically according to the **EMVCo specification**.
- To dynamic trigger Napas 24/7 bank transfer payments, generate a dynamic QR string formatted as:
  `00020101021238570010A000000727012800069704070114970407138472900208QRIBFTTA5303704540850000.005802VN5916SHOPEE_CLONE_CORP6007HANOI62220818ORDER_REF_123456304`
- Once scanned by any local banking app, the amount and description are pre-filled securely to ensure 0 manual error.
