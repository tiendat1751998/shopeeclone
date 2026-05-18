# E-Commerce Glossary

Definitions of specialized terms used across the codebase:

- **SPU (Standard Product Unit)**: The parent representation of the item. (e.g., "Sony WH-1000XM5"). Represents catalog details.
- **SKU (Stock Keeping Unit)**: The specific sellable physical item variation. (e.g., "Sony WH-1000XM5 - Silver - Warranty 12M"). Contains price and stock.
- **GMV (Gross Merchandise Value)**: Total money processed on the platform before subtraction of cancellations and refunds.
- **Ledger Entries**: Audit lines tracking financial balance changes between buyer wallets, platform escrow, and merchant balances.
- **Outbox Table**: Local transactional database table tracking domain events to guarantee eventual consistency.
