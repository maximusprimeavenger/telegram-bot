package helpers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"telegram-bot/internal/models"

	"github.com/bytedance/sonic"
)

func CreateStatusMessage(order *models.Order) string {
	switch order.Status {
	case "confirmed":
		return fmt.Sprint("🛒 Your order has been successfully placed!\n" +
			"We’ve received your request and will start processing it shortly. Stay tuned!")
	case "processing":
		return fmt.Sprint("🔧 Your order is being processed.\n" +
			"Our team is carefully preparing your items. Thank you for your patience!")
	case "shipped":
		return fmt.Sprint("🚚 Your order is on its way!\n" +
			"We’ve handed it over to the delivery service. You’ll get it soon!")
	case "delivered":
		return fmt.Sprint("📬 Your order has been delivered!\n" +
			"We hope everything is perfect. Enjoy your purchase! 💝")
	case "returned":
		return fmt.Sprint("↩️ Your order has been marked as returned.\n" +
			"We’re processing your return. Let us know if you need any help.")
	default:
		return ""
	}
}

func GetRedisOrder(id string) (*models.Order, error) {
	url := fmt.Sprintf("http://cart-service:8083/order/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var order *models.Order
	var enc1 = sonic.ConfigDefault.NewDecoder(resp.Body)
	err = enc1.Decode(&order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func UpdateOrder(id string, order *models.Order) {
	orderJSON, err := sonic.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("http://cart-service:8083/order/update/%s", id)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(orderJSON)))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Fatalf("status code isn't a 200 or 201!, error: %v", err)
	}
}
