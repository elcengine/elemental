package elemental

import "go.mongodb.org/mongo-driver/event"

const EventDeploymentDiscovered = "DeploymentDiscovered"

var eventListeners = make(map[string]map[string]*func())

func triggerEventIfRegistered(alias, eventType string) {
	if eventListeners[alias][eventType] != nil {
		(*eventListeners[alias][eventType])()
	}
}

func defaultPoolMonitor(alias string) *event.PoolMonitor {
	poolMonitor := &event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			if len(eventListeners[alias]) > 0 {
				switch evt.Type {
				case event.ConnectionClosed:
					triggerEventIfRegistered(alias, event.ConnectionClosed)
				case event.ConnectionCreated:
					triggerEventIfRegistered(alias, event.ConnectionCreated)
				case event.ConnectionReady:
					triggerEventIfRegistered(alias, event.ConnectionReady)
				case event.ConnectionReturned:
					triggerEventIfRegistered(alias, event.ConnectionReturned)
				case event.GetFailed:
					triggerEventIfRegistered(alias, event.GetFailed)
				case event.GetStarted:
					triggerEventIfRegistered(alias, event.GetStarted)
				case event.GetSucceeded:
					triggerEventIfRegistered(alias, event.GetSucceeded)
				case event.PoolCleared:
					triggerEventIfRegistered(alias, event.PoolCleared)
				case event.PoolClosedEvent:
					triggerEventIfRegistered(alias, event.PoolClosedEvent)
				case event.PoolCreated:
					triggerEventIfRegistered(alias, event.PoolCreated)
				case event.PoolReady:
					triggerEventIfRegistered(alias, event.PoolReady)
				}
			}
		},
	}
	return poolMonitor
}

// Add a listener to a connection event
//
// @param event - The event to listen for. From the default mongo driver event package or from elemental
//
// @param handler - The function to call when the event is triggered
//
// @param alias - The alias of the connection to listen to
func OnConnectionEvent(event string, handler func(), alias ...string) {
	if len(alias) == 0 {
		alias = []string{"default"}
	}
	if eventListeners[alias[0]] == nil {
		eventListeners[alias[0]] = make(map[string]*func())
	}
	eventListeners[alias[0]][event] = &handler
}

// Remove a listener from a connection event
//
// @param event - The event to remove the listener from
//
// @param alias - The alias of the connection to remove the listener from
func RemoveConnectionEvent(event string, alias ...string) {
	if len(alias) == 0 {
		alias = []string{"default"}
	}
	if eventListeners[alias[0]] != nil {
		eventListeners[alias[0]][event] = nil
	}
}
