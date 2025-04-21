package repository

import (
	"context"
	"fmt"
	"log"
	"telegram-bot/internal/db"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"
	"time"
)

func (client *Client) CheckOrders(update models.Update, callbackId int) {
	dbConn, err := db.ConnectToSQL()
	if err != nil {
		log.Printf("Error connecting to db:%v\n", err)
	}
	repo := NewUserRepository(dbConn)

	user, err := repo.findUser(callbackId)
	if err != nil {
		log.Printf("Error finding user:%v\n", err)
	}
	order, err := helpers.GetRedisOrder(user.NotifierID)
	if err != nil {
		log.Printf("Erorr sending to notifier: %v\n", err)
	}
	textItem := "🛍️ Order details:\n"
	for _, item := range order.Items {
		textItem += fmt.Sprintf("• %s - %d \n", item.Description, item.Quantity)
	}
	orderNum := fmt.Sprintf("📦 Your order number: %s\n", order.OrderNumber)
	orderDate := fmt.Sprintf("📅 Order completion date: %s\n", order.CreatedAt.Format("02 Jan 2006 15:04"))
	orderStatus := fmt.Sprintf("🚚 Order status: %s\n", order.Status)
	if len(order.Items) == 0 {
		client.SendMessage(callbackId, "📦 You don't have any items in your order")
		return
	}
	client.SendMessage(callbackId, textItem+orderNum+orderDate+orderStatus+"💬 We will send you a notification as soon as the order is delivered!")
}

func (client *Client) CheckStatus(userId int, id string, ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			order, err := helpers.GetRedisOrder(id)
			if err != nil {
				log.Println("error getting updated order:", err)
				continue
			}

			if order.Status != order.LastStatus {
				text := helpers.CreateStatusMessage(order)
				if text != "" {
					client.SendMessage(userId, text)
					order.LastStatus = order.Status
					go helpers.UpdateOrder(id, order)
				} else {
					return fmt.Errorf("error during sending a message")
				}
			}

		case <-ctx.Done():
			log.Println("Context canceled, stopping CheckStatus goroutine")
			return ctx.Err()
		}
	}
}

func (client *Client) RunNotifications(id int, session *models.UserSession) {
	ctx, cancel := context.WithCancel(context.Background())
	session.Cancel = cancel
	dbConn, err := db.ConnectToSQL()
	if err != nil {
		log.Fatal(err)
	}
	repo := NewUserRepository(dbConn)
	user, err := repo.findUser(id)
	if err != nil {
		log.Fatalln("Gorutine error: ", err)
	}
	mode, _ := NotifyMode(id)
	if mode == "on" {
		go client.CheckStatus(id, user.NotifierID, ctx)
	} else {
		session.Cancel()
		return
	}
}
