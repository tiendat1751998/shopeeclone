-- Auto-create all databases if they don't exist
-- This runs on MySQL first startup via /docker-entrypoint-initdb.d/

CREATE DATABASE IF NOT EXISTS shopee_platform;
CREATE DATABASE IF NOT EXISTS shopee_auth;
CREATE DATABASE IF NOT EXISTS shopee_cart;
CREATE DATABASE IF NOT EXISTS shopee_checkout;
CREATE DATABASE IF NOT EXISTS shopee_inventory;
CREATE DATABASE IF NOT EXISTS shopee_order;
CREATE DATABASE IF NOT EXISTS shopee_payment;
CREATE DATABASE IF NOT EXISTS shopee_product;
CREATE DATABASE IF NOT EXISTS shopee_promotion;
CREATE DATABASE IF NOT EXISTS shopee_shipment;
CREATE DATABASE IF NOT EXISTS shopee_catalog;
CREATE DATABASE IF NOT EXISTS shopee_oms;
CREATE DATABASE IF NOT EXISTS shopee_logistics;
CREATE DATABASE IF NOT EXISTS shopee_notification;
CREATE DATABASE IF NOT EXISTS shopee_search;
CREATE DATABASE IF NOT EXISTS shopee_recommendation;
CREATE DATABASE IF NOT EXISTS shopee_analytics;
CREATE DATABASE IF NOT EXISTS shopee_fraud;
CREATE DATABASE IF NOT EXISTS shopee_billing;
CREATE DATABASE IF NOT EXISTS shopee_advertising;
CREATE DATABASE IF NOT EXISTS shopee_live_commerce;
CREATE DATABASE IF NOT EXISTS shopee_user_behavior;
CREATE DATABASE IF NOT EXISTS shopee_developer;
CREATE DATABASE IF NOT EXISTS shopee_sre;
CREATE DATABASE IF NOT EXISTS shopee_service_mesh;
CREATE DATABASE IF NOT EXISTS shopee_global_infra;
CREATE DATABASE IF NOT EXISTS shopee_rec_vector;
CREATE DATABASE IF NOT EXISTS shopee_payment_ledger;
CREATE DATABASE IF NOT EXISTS shopee_oms_fulfillment;
CREATE DATABASE IF NOT EXISTS shopee_live_scale;
CREATE DATABASE IF NOT EXISTS shopee_notification_campaign;

-- Grant permissions to shopee user for all databases
GRANT ALL PRIVILEGES ON shopee_platform.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_auth.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_cart.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_checkout.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_inventory.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_order.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_payment.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_product.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_promotion.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_shipment.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_catalog.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_oms.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_logistics.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_notification.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_search.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_recommendation.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_analytics.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_fraud.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_billing.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_advertising.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_live_commerce.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_user_behavior.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_developer.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_sre.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_service_mesh.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_global_infra.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_rec_vector.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_payment_ledger.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_oms_fulfillment.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_live_scale.* TO 'shopee'@'%';
GRANT ALL PRIVILEGES ON shopee_notification_campaign.* TO 'shopee'@'%';

FLUSH PRIVILEGES;
