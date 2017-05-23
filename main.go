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
	"log"
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
	http.ListenAndServe(":8000", mux)
}

func handleGetInfo(w http.ResponseWriter, r *http.Request) {
	info := getInfo()
	if info.IsSuccess() {
		fmt.Fprintf(w, "== BALANCE ==\n")
		fmt.Fprintf(w, "%s %f \n", "IDR", info.Return.Balance.Idr)
		fmt.Fprintf(w, "%s %g \n", "BTC", info.Return.Balance.Btc)
		fmt.Fprintf(w, "%s %g \n", "LTC", info.Return.Balance.Ltc)
		fmt.Fprintf(w, "%s %g \n", "DOGE", info.Return.Balance.Doge)
		fmt.Fprintf(w, "%s %g \n", "XRP", info.Return.Balance.Xrp)
		fmt.Fprintf(w, "%s %g \n", "DRK", info.Return.Balance.Drk)
		fmt.Fprintf(w, "%s %g \n", "BTS", info.Return.Balance.Bts)
		fmt.Fprintf(w, "%s %g \n", "NXT", info.Return.Balance.Nxt)
		fmt.Fprintf(w, "%s %g \n", "STR", info.Return.Balance.Str)
		fmt.Fprintf(w, "%s %g \n", "NEM", info.Return.Balance.Nem)
		fmt.Fprintf(w, "%s %g \n\n\n", "ETH", info.Return.Balance.Eth)

		fmt.Fprintf(w, "== BALANCE HOLD ==\n")
		fmt.Fprintf(w, "%s %f \n", "IDR", info.Return.BalanceHold.Idr)
		fmt.Fprintf(w, "%s %g \n", "BTC", info.Return.BalanceHold.Btc)
		fmt.Fprintf(w, "%s %g \n", "LTC", info.Return.BalanceHold.Ltc)
		fmt.Fprintf(w, "%s %g \n", "DOGE", info.Return.BalanceHold.Doge)
		fmt.Fprintf(w, "%s %g \n", "XRP", info.Return.BalanceHold.Xrp)
		fmt.Fprintf(w, "%s %g \n", "DRK", info.Return.BalanceHold.Drk)
		fmt.Fprintf(w, "%s %g \n", "BTS", info.Return.BalanceHold.Bts)
		fmt.Fprintf(w, "%s %g \n", "NXT", info.Return.BalanceHold.Nxt)
		fmt.Fprintf(w, "%s %g \n", "STR", info.Return.BalanceHold.Str)
		fmt.Fprintf(w, "%s %g \n", "NEM", info.Return.BalanceHold.Nem)
		fmt.Fprintf(w, "%s %g \n\n\n", "ETH", info.Return.BalanceHold.Eth)
	} else {
		fmt.Fprintln(w, "[ERROR] "+info.Error)
		log.Fatal("[ERROR] " + info.Error)
	}

	currentPrice := getCurrentPrice()
	totalBtc := calculateTotalBtc(info.Return.Balance, info.Return.BalanceHold, currentPrice.Price)
	totalIdr := calculateTotalIdr(info.Return.Balance, info.Return.BalanceHold, currentPrice.Price)

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

func getInfo() Info {
	var info Info
	data := generateData("getInfo")
	body := sendRequest(data, ApiUrl)

	err := json.Unmarshal(body, &info)
	if err != nil {
		log.Fatal(err)
	}
	return info
}

func getCurrentPrice() CurrentPrice {
	res, err := http.Get("https://vip.bitcoin.co.id/api/eth_btc/webdata")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var currentPrice CurrentPrice
	errJson := json.Unmarshal(body, &currentPrice)
	if errJson != nil {
		log.Fatal(errJson)
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
		log.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func getSign(data string, secret string) string {
	sign := hmac.New(sha512.New, []byte(secret))
	sign.Write([]byte(data))

	return hex.EncodeToString(sign.Sum(nil))
}

func recordAssetToDb(balance Balance, balanceHold Balance, price Rate) int {
	dbDsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(localhost:3306)/" + os.Getenv("DB_NAME")
	db, err := sql.Open("mysql", dbDsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	errPing := db.Ping()
	if errPing != nil {
		log.Fatal(errPing)
	}

	stmt, err := db.Prepare("INSERT INTO assets (ping_time, idr, btc, ltc, doge, xrp, drk, bts, nxt, str, nem, eth, idr_hold, btc_hold, ltc_hold, doge_hold, xrp_hold, drk_hold, bts_hold, nxt_hold, str_hold, nem_hold, eth_hold, price_btc_idr, price_ltc_btc, price_doge_btc, price_xrp_btc, price_drk_btc, price_bts_btc, price_nxt_btc, price_str_btc, price_nem_btc, price_eth_btc) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(time.Now().Format("2006-01-02 03:04:05"), balance.Idr, balance.Btc, balance.Ltc, balance.Doge, balance.Xrp, balance.Drk, balance.Bts, balance.Nxt, balance.Str, balance.Nem, balance.Eth, balanceHold.Idr, balanceHold.Btc, balanceHold.Ltc, balanceHold.Doge, balanceHold.Xrp, balanceHold.Drk, balanceHold.Bts, balanceHold.Nxt, balanceHold.Str, balanceHold.Nem, balanceHold.Eth, price.BtcIdr, price.LtcBtc, price.DogeBtc, price.XrpBtc, price.DrkBtc, price.BtsBtc, price.NxtBtc, price.StrBtc, price.NemBtc, price.EthBtc)
	if err != nil {
		log.Fatal(err)
	}
	lastId, _ := res.LastInsertId()
	return int(lastId)
}
