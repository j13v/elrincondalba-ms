package pubsub

import (
	"net/http"

	"github.com/functionalfoundry/graphqlws"
	"github.com/gorilla/websocket"
	"github.com/jal88/elrincondalba-ms/logger"
	log "github.com/sirupsen/logrus"
)

// NewHandler creates a WebSocket handler for GraphQL WebSocket connections.
// This handler takes a SubscriptionManager and adds/removes subscriptions
// as they are started/stopped by the client.
func NewHandlerFunc(config graphqlws.HandlerConfig) func(http.ResponseWriter, *http.Request) bool {
	// Create a WebSocket upgrader that requires clients to implement
	// the "graphql-ws" protocol
	var upgrader = websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: []string{"graphql-ws"},
	}

	logger := logger.NewLogger("handler")
	subscriptionManager := config.SubscriptionManager

	// Create a map (used like a set) to manage client connections
	var connections = make(map[graphqlws.Connection]bool)

	return func(w http.ResponseWriter, r *http.Request) bool {
		// Establish a WebSocket connection
		var ws, err = upgrader.Upgrade(w, r, nil)

		// Bail out if the WebSocket connection could not be established
		if err != nil {
			logger.Warn("Failed to establish WebSocket connection", err)
			return false
		}

		// Close the connection early if it doesn't implement the graphql-ws protocol
		if ws.Subprotocol() != "graphql-ws" {
			logger.Warn("Connection does not implement the GraphQL WS protocol")
			ws.Close()
			return false
		}

		// Establish a GraphQL WebSocket connection
		conn := graphqlws.NewConnection(ws, graphqlws.ConnectionConfig{
			Authenticate: config.Authenticate,
			EventHandlers: graphqlws.ConnectionEventHandlers{
				Close: func(conn graphqlws.Connection) {
					logger.WithFields(log.Fields{
						"conn": conn.ID(),
						"user": conn.User(),
					}).Debug("Closing connection")

					subscriptionManager.RemoveSubscriptions(conn)

					delete(connections, conn)
				},
				StartOperation: func(
					conn graphqlws.Connection,
					opID string,
					data *graphqlws.StartMessagePayload,
				) []error {
					logger.WithFields(log.Fields{
						"conn": conn.ID(),
						"op":   opID,
						"user": conn.User(),
					}).Debug("Start operation")

					return subscriptionManager.AddSubscription(conn, &graphqlws.Subscription{
						ID:            opID,
						Query:         data.Query,
						Variables:     data.Variables,
						OperationName: data.OperationName,
						Connection:    conn,
						SendData: func(data *graphqlws.DataMessagePayload) {
							conn.SendData(opID, data)
						},
					})
				},
				StopOperation: func(conn graphqlws.Connection, opID string) {
					subscriptionManager.RemoveSubscription(conn, &graphqlws.Subscription{
						ID: opID,
					})
				},
			},
		})
		connections[conn] = true
		return true
	}
}
