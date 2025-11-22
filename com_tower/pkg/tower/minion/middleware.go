package minion

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"bytes"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func AuditRequests() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		intercepter := &types.BodyIntercepter{Body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = intercepter

		ctx.Next()

		if ctx.Request.URL.Path == "/slots" && ctx.Request.Method == http.MethodPost {
			var slotReq types.SlotRequest

			defer ctx.Request.Body.Close()
			if err := json.NewDecoder(ctx.Request.Body).Decode(&slotReq); err != nil {
				log.Printf("[minion][audit][middleware] failed to decode request body: %v", err)
				return
			}

			result := types.DeniedResultType
			if ctx.Writer.Status() == http.StatusOK {
				var slotResp types.SlotResponse
				if err := json.Unmarshal(intercepter.Body.Bytes(), &slotResp); err != nil {
					log.Printf("[minion][audit][middleware] failed to unmarshal response body: %v", err)
					return
				}

				result = types.GetResultTypeBySlotState(slotResp.State)
			}

			auditReq := types.AuditRequest{
				VehicleType:   slotReq.VehicleType,
				VehicleUUID:   slotReq.VehicleUUID,
				StructureType: slotReq.StructureType,
				StructureUUID: slotReq.StructureUUID,
				Timestamp:     time.Now(),
				Result:        result,
				SlotNumber:    slotReq.SlotNumber,
			}

			body, err := json.Marshal(auditReq)
			if err != nil {
				log.Printf("[minion][audit][middleware] failed to marshal audit request message body: %v", err)
				return
			}

			payload := amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			}

			if err := config.Configuration.GetRabbitMQChannel().PublishWithContext(ctx, "requests", "requests", false, false, payload); err != nil {
				log.Printf("[minion][audit][middleware] failed to send audit request message to the broker: %v", err)
				return
			}
		}
	}
}
