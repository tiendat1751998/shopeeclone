# High-Reliability Payment State Machine

The platform integrates Stripe (global credit cards), MoMo (Vietnam wallet), and VNPay (Vietnam local banks via dynamic QR).

## Payment State Transitions
- **CREATED**: Ledger transaction entry generated; dynamic VietQR code displayed.
- **PENDING**: Webhook validation awaiting provider confirmation.
- **COMPLETED**: Signature validated, order marked completed, balances updated.
- **FAILED**: Stock released immediately back to the Inventory pool.

## Webhook Signature Verification Algorithm (MoMo Gateway example)
```java
// Java code representing HMAC verification signature
public class MoMoSecurity {
    public static boolean verifySignature(String rawData, String incomingSignature, String secretKey) {
        try {
            Mac sha256_HMAC = Mac.getInstance("HmacSHA256");
            SecretKeySpec secret_key = new SecretKeySpec(secretKey.getBytes("UTF-8"), "HmacSHA256");
            sha256_HMAC.init(secret_key);
            byte[] hashBytes = sha256_HMAC.doFinal(rawData.getBytes("UTF-8"));
            
            StringBuilder hexString = new StringBuilder();
            for (byte b : hashBytes) {
                String hex = Integer.toHexString(0xff & b);
                if (hex.length() == 1) hexString.append('0');
                hexString.append(hex);
            }
            return hexString.toString().equals(incomingSignature);
        } catch (Exception e) {
            return false;
        }
    }
}
```
