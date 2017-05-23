CREATE DATABASE vip_asset;
USE vip_asset;

CREATE TABLE assets (
    id INT(6) AUTO_INCREMENT PRIMARY KEY,
    ping_time DATETIME,
    idr  DECIMAL(10, 2),
	btc  DECIMAL(12,8),
	ltc  DECIMAL(10,2),
	doge DECIMAL(10,2),
	xrp  DECIMAL(10,2),
	drk  DECIMAL(10,2),
	bts  DECIMAL(10,2),
	nxt  DECIMAL(10,2),
	str  DECIMAL(10,2),
	nem  DECIMAL(10,2),
	eth  DECIMAL(10,2),
    idr_hold  DECIMAL(10, 2),
	btc_hold  DECIMAL(12,8),
	ltc_hold  DECIMAL(10,2),
	doge_hold DECIMAL(10,2),
	xrp_hold  DECIMAL(10,2),
	drk_hold  DECIMAL(10,2),
	bts_hold  DECIMAL(10,2),
	nxt_hold  DECIMAL(10,2),
	str_hold  DECIMAL(10,2),
	nem_hold  DECIMAL(10,2),
	eth_hold  DECIMAL(10,2),
    price_btc_idr  DECIMAL(10,2),
	price_ltc_btc  DECIMAL(10,2), 
	price_doge_btc DECIMAL(10,2),
	price_xrp_btc  DECIMAL(10,2),
	price_drk_btc  DECIMAL(10,2),
	price_bts_btc  DECIMAL(10,2),
	price_nxt_btc  DECIMAL(10,2),
	price_str_btc  DECIMAL(10,2),
	price_nem_btc  DECIMAL(10,2),
	price_eth_btc  DECIMAL(10,2)
)