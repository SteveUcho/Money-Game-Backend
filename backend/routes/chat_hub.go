package routes

type GetHubRequest struct {
	LobbyID string
	Res     chan *Hub
}

// ChatHub maintains the set of chat hubs
type ChatHub struct {
	// Registered clients.
	hubs map[string]*Hub

	// Register requests from the clients.
	register chan string

	// Unregister requests from clients.
	unregister chan string

	// request to get a hub
	getHub chan GetHubRequest
}

func NewChatHub() *ChatHub {
	hub := &ChatHub{
		hubs:       make(map[string]*Hub),
		register:   make(chan string),
		unregister: make(chan string),
		getHub:     make(chan GetHubRequest),
	}
	go hub.Run()
	return hub
}

func (h *ChatHub) Run() {
	for {
		select {
		case lobbyID := <-h.register:
			if _, ok := h.hubs[lobbyID]; !ok {
				h.hubs[lobbyID] = NewHub()
			}
		case lobbyID := <-h.unregister:
			delete(h.hubs, lobbyID)
		case req := <-h.getHub:
			if hub, ok := h.hubs[req.LobbyID]; ok {
				req.Res <- hub
			}
		}
	}
}
