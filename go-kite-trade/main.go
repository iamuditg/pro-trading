package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
	"os"
	"strconv"
	"time"
)

var table *tablewriter.Table

var (
	apiKey    string = getEnv("KITE_API_KEY", "my_api_key")
	apiSecret string = getEnv("KITE_API_SECRET", "my_api_secret")
	instToken uint32 = getEnvUint32("KITE_INSTRUMENT_TOKEN", 256265)
)
var (
	ticker *kiteticker.Ticker
)

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	fmt.Println("Connected")
	fmt.Println("Subscribing to", instToken)
	err := ticker.Subscribe([]uint32{instToken})
	if err != nil {
		fmt.Println("err: ", err)
	}
	// Set subscription mode for given list of tokens
	// Default mode is Quote
	err = ticker.SetMode(kiteticker.ModeFull, []uint32{instToken})
	if err != nil {
		fmt.Println("err: ", err)
	}
}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {
	tickData := TickData{
		Timestamp:          time.Now(),
		LastPrice:          tick.LastPrice,
		LastTradedQuantity: tick.LastTradedQuantity,
	}

	// Clear the console (optional) before updating the table
	fmt.Print("\033[H\033[2J")

	// If the table is not initialized, initialize it
	if table == nil {
		initTable()
	}

	// Clear the table content
	table.ClearRows()

	// Append the updated tick data to the table
	appendTickData([]TickData{tickData})

	// Render the updated table
	table.Render()
}

// Function to append tick data to the table
func appendTickData(ticks []TickData) {
	for _, tick := range ticks {
		table.Append([]string{
			tick.Timestamp.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%.2f", tick.LastPrice),
			fmt.Sprintf("%d", tick.LastTradedQuantity),
		})
	}
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	fmt.Printf("Order: %s", order.OrderID)
}

func main() {
	// Create a new Kite connect instance
	kc := kiteconnect.New(apiKey)

	// Login URL from which request token can be obtained
	fmt.Println(kc.GetLoginURL())

	// Obtained request token after Kite Connect login flow
	// simulated here by scanning from stdin
	var requestToken string
	fmt.Println("Enter request token:")
	fmt.Scanf("%s\n", &requestToken)

	// Get user details and access token
	data, err := kc.GenerateSession(requestToken, apiSecret)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	//kc.SetAccessToken(data.AccessToken)
	//
	//// fetch the instruments list
	//instrument, err := kc.GetInstruments()
	//if err != nil {
	//	fmt.Println("Error fetching instruments", err)
	//	return
	//}
	//
	//// Print the instruments list
	//for _, instru := range instrument {
	//	fmt.Printf("Instrument Token: %d, Trading Symbol: %s\n", instru.InstrumentToken, instru.Name)
	//}

	// Create new Kite ticker instance
	ticker = kiteticker.New(apiKey, data.AccessToken)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// Start the connection
	ticker.Serve()
}

// getEnv returns the value of the environment variable provided.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvUint32 returns the value of the environment variable provided converted as Uint32.
func getEnvUint32(key string, fallback int) uint32 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return uint32(fallback)
		}
		return uint32(i)
	}
	return uint32(fallback)
}

// Initialize the table with headers
func initTable() {
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Timestamp", "Last Price", "Last Traded Quantity"})

	// Define table formatting options (you can customize these as needed)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetHeaderColor(
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiYellowColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiYellowColor),
		tablewriter.Color(tablewriter.Bold, tablewriter.FgHiYellowColor),
	)
	table.SetColumnColor(
		tablewriter.Color(tablewriter.FgHiWhiteColor),
		tablewriter.Color(tablewriter.FgHiGreenColor),
		tablewriter.Color(tablewriter.FgHiGreenColor),
	)
}
