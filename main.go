package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ApiUrl = "https://vip.bitcoin.co.id/tapi/"
)

type Info struct {
	Success int    `json:"success"`
	Return  Return `json:"return"`
	Error   string `json:"error"`
}

type CurrentPrice struct {
	Price Rate `json:"prices"`
}

type Return struct {
	Balance     Balance `json:"balance"`
	BalanceHold Balance `json:"balance_hold"`
}

type JsonAsset struct {
	PingTime     time.Time `json:"ping_time"`
	Idr          float64   `json:"idr"`
	Btc          float64   `json:"btc"`
	Ltc          float64   `json:"ltc"`
	Doge         float64   `json:"doge"`
	Xrp          float64   `json:"xrp"`
	Drk          float64   `json:"drk"`
	Bts          float64   `json:"btc"`
	Nxt          float64   `json:"nxt"`
	Str          float64   `json:"str"`
	Nem          float64   `json:"nem"`
	Eth          float64   `json:"eth"`
	IdrHold      float64   `json:"idr_hold"`
	BtcHold      float64   `json:"btc_hold"`
	LtcHold      float64   `json:"ltc_hold"`
	DogeHold     float64   `json:"doge_hold"`
	XrpHold      float64   `json:"xrp_hold"`
	DrkHold      float64   `json:"drk_hold"`
	BtsHold      float64   `json:"btc_hold"`
	NxtHold      float64   `json:"nxt_hold"`
	StrHold      float64   `json:"str_hold"`
	NemHold      float64   `json:"nem_hold"`
	EthHold      float64   `json:"eth_hold"`
	PriceBtcIdr  float64   `json:"price_btc_idr"`
	PriceLtcBtc  float64   `json:"price_ltc_btc"`
	PriceDogeBtc float64   `json:"price_doge_btc"`
	PriceXrpBtc  float64   `json:"price_xrp_btc"`
	PriceDrkBtc  float64   `json:"price_drk_btc"`
	PriceBtsBtc  float64   `json:"price_bts_btc"`
	PriceNxtBtc  float64   `json:"price_nxt_btc"`
	PriceStrBtc  float64   `json:"price_str_btc"`
	PriceNemBtc  float64   `json:"price_nem_btc"`
	PriceEthBtc  float64   `json:"price_eth_btc"`
}

type Balance struct {
	Idr  float64 `json:"idr"`
	Btc  float64 `json:"btc,string"`
	Ltc  float64 `json:"ltc,string"`
	Doge float64 `json:"doge,string"`
	Xrp  float64 `json:"xrp,string"`
	Drk  float64 `json:"drk,string"`
	Bts  float64 `json:"bts,string"`
	Nxt  float64 `json:"nxt,string"`
	Str  float64 `json:"str,string"`
	Nem  float64 `json:"nem,string"`
	Eth  float64 `json:"eth,string"`
}

type Rate struct {
	BtcIdr  float64 `json:"btcidr,string"`
	LtcBtc  float64 `json:"ltcbtc,string"`
	DogeBtc float64 `json:"dogebtc,string"`
	XrpBtc  float64 `json:"xrpbtc,string"`
	DrkBtc  float64 `json:"drkbtc,string"`
	BtsBtc  float64 `json:"btsbtc,string"`
	NxtBtc  float64 `json:"nxtbtc,string"`
	StrBtc  float64 `json:"strbtc,string"`
	NemBtc  float64 `json:"nembtc,string"`
	EthBtc  float64 `json:"ethbtc,string"`
}

func (info Info) IsSuccess() bool {
	return info.Success == 1
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleGetInfo)
	mux.HandleFunc("/record-asset", handleRecordAsset)
	mux.HandleFunc("/get-asset", handleGetAsset)
	http.ListenAndServe(":8000", mux)
}

func handleGetInfo(w http.ResponseWriter, r *http.Request) {
	info := getInfo()
	currentPrice := getCurrentPrice()
	balance := info.Return.Balance
	balanceHold := info.Return.BalanceHold
	price := currentPrice.Price
	totalBtc := calculateTotalBtc(balance, balanceHold, price)
	totalIdr := calculateTotalIdr(balance, balanceHold, price)

	if info.IsSuccess() {
		fmt.Fprintf(w, "== BALANCE ==\n")
		fmt.Fprintf(w, "%s %f \n", "IDR", balance.Idr)
		fmt.Fprintf(w, "%s %g \n", "BTC", balance.Btc)
		fmt.Fprintf(w, "%s %g \n", "LTC", balance.Ltc)
		fmt.Fprintf(w, "%s %g \n", "DOGE", balance.Doge)
		fmt.Fprintf(w, "%s %g \n", "XRP", balance.Xrp)
		fmt.Fprintf(w, "%s %g \n", "DRK", balance.Drk)
		fmt.Fprintf(w, "%s %g \n", "BTS", balance.Bts)
		fmt.Fprintf(w, "%s %g \n", "NXT", balance.Nxt)
		fmt.Fprintf(w, "%s %g \n", "STR", balance.Str)
		fmt.Fprintf(w, "%s %g \n", "NEM", balance.Nem)
		fmt.Fprintf(w, "%s %g \n\n\n", "ETH", balance.Eth)

		fmt.Fprintf(w, "== BALANCE HOLD ==\n")
		fmt.Fprintf(w, "%s %f \n", "IDR", balanceHold.Idr)
		fmt.Fprintf(w, "%s %g \n", "BTC", balanceHold.Btc)
		fmt.Fprintf(w, "%s %g \n", "LTC", balanceHold.Ltc)
		fmt.Fprintf(w, "%s %g \n", "DOGE", balanceHold.Doge)
		fmt.Fprintf(w, "%s %g \n", "XRP", balanceHold.Xrp)
		fmt.Fprintf(w, "%s %g \n", "DRK", balanceHold.Drk)
		fmt.Fprintf(w, "%s %g \n", "BTS", balanceHold.Bts)
		fmt.Fprintf(w, "%s %g \n", "NXT", balanceHold.Nxt)
		fmt.Fprintf(w, "%s %g \n", "STR", balanceHold.Str)
		fmt.Fprintf(w, "%s %g \n", "NEM", balanceHold.Nem)
		fmt.Fprintf(w, "%s %g \n\n\n", "ETH", balanceHold.Eth)

		fmt.Fprintf(w, "== PERCENTAGE ==\n")
		fmt.Fprintf(w, "%s %f \n", "IDR", (balance.Idr+balanceHold.Idr)/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n", "BTC", (balance.Btc+balanceHold.Btc)*price.BtcIdr/totalIdr*100)
		fmt.Fprintf(w, "%s %.2f\n", "LTC", (balance.Ltc+balanceHold.Ltc)*price.LtcBtc/1000000*price.BtcIdr/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n", "DOGE", (balance.Doge+balanceHold.Doge)*price.DogeBtc/1000000*price.BtcIdr/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n", "XRP", (balance.Xrp+balanceHold.Xrp)*price.XrpBtc/1000000*price.BtcIdr/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n", "DRK", (balance.Drk+balanceHold.Drk)*price.DrkBtc/1000000*price.BtcIdr/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n", "BTS", (balance.Bts+balanceHold.Bts)*price.BtsBtc/1000000*price.BtcIdr/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n", "NXT", (balance.Nxt+balanceHold.Nxt)*price.NxtBtc/1000000*price.BtcIdr/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n", "STR", (balance.Str+balanceHold.Str)*price.StrBtc/1000000*price.BtcIdr/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n", "NEM", (balance.Nem+balanceHold.Nem)*price.NemBtc/1000000*price.BtcIdr/totalIdr)
		fmt.Fprintf(w, "%s %.2f\n\n\n", "ETH", (balance.Eth+balanceHold.Eth)*price.EthBtc/1000000*price.BtcIdr/totalIdr)
	} else {
		panic("[ERROR] " + info.Error)
	}

	fmt.Fprintf(w, "== TOTAL ==\n")
	fmt.Fprintf(w, "Total Asset in BTC: %f\n", totalBtc)
	fmt.Fprintf(w, "Total Asset in IDR: %f", totalIdr)
}

func handleRecordAsset(w http.ResponseWriter, r *http.Request) {
	info := getInfo()
	currentPrice := getCurrentPrice()
	id := recordAssetToDb(info.Return.Balance, info.Return.BalanceHold, currentPrice.Price)
	fmt.Fprintf(w, "INSERT ID: %d", id)
}

func handleGetAsset(w http.ResponseWriter, r *http.Request) {
	var jsonAssets []JsonAsset
	limit, _ := strconv.ParseInt(r.FormValue("limit"), 10, 32)
	rows := getAssetFromDb(int(limit))
	for rows.Next() {
		var jsonAsset JsonAsset
		var id int
		var pingTime time.Time
		var idr, btc, ltc, doge, xrp, drk, bts, nxt, str, nem, eth, idrHold, btcHold, ltcHold, dogeHold, xrpHold, drkHold, btsHold, nxtHold, strHold, nemHold, ethHold, priceBtcIdr, priceLtcBtc, priceDogeBtc, priceXrpBtc, priceDrkBtc, priceBtsBtc, priceNxtBtc, priceStrBtc, priceNemBtc, priceEthBtc float64

		err := rows.Scan(&id, &pingTime, &idr, &btc, &ltc, &doge, &xrp, &drk, &bts, &nxt, &str, &nem, &eth, &idrHold, &btcHold, &ltcHold, &dogeHold, &xrpHold, &drkHold, &btsHold, &nxtHold, &strHold, &nemHold, &ethHold, &priceBtcIdr, &priceLtcBtc, &priceDogeBtc, &priceXrpBtc, &priceDrkBtc, &priceBtsBtc, &priceNxtBtc, &priceStrBtc, &priceNemBtc, &priceEthBtc)

		if err != nil {
			panic(err)
		}
		jsonAsset = JsonAsset{
			PingTime:     pingTime,
			Idr:          idr,
			Btc:          btc,
			Ltc:          ltc,
			Doge:         doge,
			Xrp:          xrp,
			Drk:          drk,
			Bts:          bts,
			Nxt:          nxt,
			Str:          str,
			Nem:          nem,
			Eth:          eth,
			IdrHold:      idrHold,
			BtcHold:      btcHold,
			LtcHold:      ltcHold,
			DogeHold:     dogeHold,
			XrpHold:      xrpHold,
			DrkHold:      drkHold,
			BtsHold:      btsHold,
			NxtHold:      nxtHold,
			StrHold:      strHold,
			NemHold:      nemHold,
			EthHold:      ethHold,
			PriceBtcIdr:  priceBtcIdr,
			PriceLtcBtc:  priceLtcBtc,
			PriceDogeBtc: priceDogeBtc,
			PriceXrpBtc:  priceXrpBtc,
			PriceDrkBtc:  priceDrkBtc,
			PriceBtsBtc:  priceBtsBtc,
			PriceNxtBtc:  priceNxtBtc,
			PriceStrBtc:  priceStrBtc,
			PriceNemBtc:  priceNemBtc,
			PriceEthBtc:  priceEthBtc}
		jsonAssets = append(jsonAssets, jsonAsset)
	}
	json.NewEncoder(w).Encode(&jsonAssets)
}

func getInfo() Info {
	var info Info
	data := generateData("getInfo")
	body := sendRequest(data, ApiUrl)

	err := json.Unmarshal(body, &info)
	if err != nil {
		panic(err)
	}
	return info
}

func getCurrentPrice() CurrentPrice {
	res, err := http.Get("https://vip.bitcoin.co.id/api/eth_btc/webdata")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var currentPrice CurrentPrice
	errJson := json.Unmarshal(body, &currentPrice)
	if errJson != nil {
		panic(errJson)
	}

	return currentPrice
}

func calculateTotalBtc(balance Balance, balanceHold Balance, price Rate) float64 {
	totalBtc := 0.0
	totalBtc += (balance.Btc + balanceHold.Btc)
	totalBtc += (balance.Ltc + balanceHold.Ltc) * price.LtcBtc / 100000000
	totalBtc += (balance.Doge + balanceHold.Doge) * price.DogeBtc / 100000000
	totalBtc += (balance.Xrp + balanceHold.Xrp) * price.XrpBtc / 100000000
	totalBtc += (balance.Drk + balanceHold.Drk) * price.DrkBtc / 100000000
	totalBtc += (balance.Bts + balanceHold.Bts) * price.BtsBtc / 100000000
	totalBtc += (balance.Nxt + balanceHold.Nxt) * price.NxtBtc / 100000000
	totalBtc += (balance.Str + balanceHold.Str) * price.StrBtc / 100000000
	totalBtc += (balance.Nem + balanceHold.Nem) * price.NemBtc / 100000000
	totalBtc += (balance.Eth + balanceHold.Eth) * price.EthBtc / 100000000
	return totalBtc
}

func calculateTotalIdr(balance Balance, balanceHold Balance, price Rate) float64 {
	totalBtc := calculateTotalBtc(balance, balanceHold, price)
	total := (totalBtc * price.BtcIdr) + (balance.Idr + balanceHold.Idr)
	return total
}

func generateData(method string) string {
	nonce := int(time.Now().Unix())
	data := "method=getInfo&nonce=" + strconv.Itoa(nonce)
	return data
}

func sendRequest(data string, url string) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	req.Header.Add("Sign", getSign(data, os.Getenv("VIP_SECRET")))
	req.Header.Add("Key", os.Getenv("VIP_KEY"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func getSign(data string, secret string) string {
	sign := hmac.New(sha512.New, []byte(secret))
	sign.Write([]byte(data))

	return hex.EncodeToString(sign.Sum(nil))
}

func getDsn() string {
	return os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(localhost:3306)/" + os.Getenv("DB_NAME") + "?parseTime=true"
}

func recordAssetToDb(balance Balance, balanceHold Balance, price Rate) int {
	db, err := sql.Open("mysql", getDsn())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	errPing := db.Ping()
	if errPing != nil {
		panic(errPing)
	}

	stmt, err := db.Prepare("INSERT INTO assets (ping_time, idr, btc, ltc, doge, xrp, drk, bts, nxt, str, nem, eth, idr_hold, btc_hold, ltc_hold, doge_hold, xrp_hold, drk_hold, bts_hold, nxt_hold, str_hold, nem_hold, eth_hold, price_btc_idr, price_ltc_btc, price_doge_btc, price_xrp_btc, price_drk_btc, price_bts_btc, price_nxt_btc, price_str_btc, price_nem_btc, price_eth_btc) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(time.Now().Format("2006-01-02 03:04:05"), balance.Idr, balance.Btc, balance.Ltc, balance.Doge, balance.Xrp, balance.Drk, balance.Bts, balance.Nxt, balance.Str, balance.Nem, balance.Eth, balanceHold.Idr, balanceHold.Btc, balanceHold.Ltc, balanceHold.Doge, balanceHold.Xrp, balanceHold.Drk, balanceHold.Bts, balanceHold.Nxt, balanceHold.Str, balanceHold.Nem, balanceHold.Eth, price.BtcIdr, price.LtcBtc, price.DogeBtc, price.XrpBtc, price.DrkBtc, price.BtsBtc, price.NxtBtc, price.StrBtc, price.NemBtc, price.EthBtc)
	if err != nil {
		panic(err)
	}
	lastId, _ := res.LastInsertId()
	return int(lastId)
}

func getAssetFromDb(limit int) *sql.Rows {
	db, err := sql.Open("mysql", getDsn())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	errPing := db.Ping()
	if errPing != nil {
		panic(errPing)
	}

	sql := "SELECT * FROM assets ORDER BY id DESC LIMIT " + strconv.Itoa(limit)
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	return rows
}
