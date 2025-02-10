package workers

import (
	"context"
	"fmt"
	"log"
	"ordermanagement/config"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

// ProcessOrders picks pending orders and completes them asynchronously
func ProcessOrders(db *gorm.DB) {
	for {
		var orders []struct {
			ID uuid.UUID
		}

		// Fetch pending orders from the orders table
		db.Raw("SELECT id FROM orders WHERE status = 'pending'").Scan(&orders)

		// Process each order
		for _, order := range orders {
			// Get the current time for logging
			currentTime := time.Now()

			// Mark the order as processing
			db.Exec("UPDATE orders SET status = 'processing' WHERE id = ?", order.ID)

			// Log the start of the processing in MongoDB with timestamp
			_, err := config.MongoDB.Collection("order_history").UpdateOne(
				context.TODO(),
				bson.M{"order_id": order.ID.String()},
				bson.M{
					"$set": bson.M{
						"status": "processing",
					},
					"$push": bson.M{
						"log": bson.M{
							"status":    "processing",
							"timestamp": currentTime,
						},
					},
				},
			)
			if err != nil {
				log.Println("Failed to update MongoDB:", err)
			}

			// Simulate order processing for 2 seconds
			time.Sleep(100 * time.Second)

			// Mark the order as completed after processing
			db.Exec("UPDATE orders SET status = 'completed' WHERE id = ?", order.ID)

			// Log the completion of the order in MongoDB with timestamp
			_, err = config.MongoDB.Collection("order_history").UpdateOne(
				context.TODO(),
				bson.M{"order_id": order.ID.String()},
				bson.M{
					"$set": bson.M{
						"status": "completed",
					},
					"$push": bson.M{
						"log": bson.M{
							"status":    "completed",
							"timestamp": time.Now(),
						},
					},
				},
			)
			if err != nil {
				log.Println("Failed to update MongoDB:", err)
			}

			// Log that the order has been processed
			fmt.Println("Processed Order:", order.ID)
		}

		// Optional: Sleep for a short duration before checking for new orders
		time.Sleep(10 * time.Second)
	}
}
