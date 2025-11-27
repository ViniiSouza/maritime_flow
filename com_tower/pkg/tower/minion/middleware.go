package minion

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func AuditRequests() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path != "/slots" || ctx.Request.Method != http.MethodPost {
			ctx.Next()
			return
		}

		intercepter := &types.BodyIntercepter{Body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = intercepter

		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
			return
		}

		bodyCopy := bytes.NewBuffer(bodyBytes)                      // save body
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body so handler can read it again

		ctx.Next()

		var slotReq types.SlotRequest
		defer ctx.Request.Body.Close()
		if err := json.NewDecoder(bodyCopy).Decode(&slotReq); err != nil {
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
			Timestamp:     int(time.Now().Unix()),
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
