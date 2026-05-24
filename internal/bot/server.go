package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type Offer struct {
	Store         string `json:"store"`
	Title         string `json:"title"`
	Price         string `json:"price"`
	OriginalPrice string `json:"original_price"`
	Discount      string `json:"discount"`
	Link          string `json:"link"`
	GroupJID      string `json:"group_jid"` // opcional: sobrescreve o(s) grupo(s) do .env
}

func (b *Bot) StartServer() {
	godotenv.Load()
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/offer", b.handleOffer)
	fmt.Printf("Servidor HTTP iniciado na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (b *Bot) handleOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	var offer Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if err := b.SendOffer(offer); err != nil {
		http.Error(w, "Erro ao enviar oferta", http.StatusInternalServerError)
		fmt.Printf("Erro ao enviar oferta: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func storeEmoji(store string) string {
	switch strings.ToLower(store) {
	case "shopee":
		return "🧡"
	case "ml", "mercadolivre", "mercado_livre":
		return "💛"
	case "aliexpress":
		return "❤️"
	default:
		return "🛍️"
	}
}

func storeName(store string) string {
	switch strings.ToLower(store) {
	case "shopee":
		return "SHOPEE"
	case "ml", "mercadolivre", "mercado_livre":
		return "MERCADO LIVRE"
	case "aliexpress":
		return "ALIEXPRESS"
	default:
		return strings.ToUpper(store)
	}
}

func (b *Bot) resolveGroups(offer Offer) ([]types.JID, error) {
	godotenv.Load()

	raw := offer.GroupJID
	if raw == "" {
		raw = os.Getenv("GROUP_JID")
	}
	if raw == "" {
		return nil, fmt.Errorf("GROUP_JID nao configurado")
	}

	parts := strings.Split(raw, ",")
	var jids []types.JID
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		jid, err := types.ParseJID(p)
		if err != nil {
			return nil, fmt.Errorf("JID invalido '%s': %v", p, err)
		}
		jids = append(jids, jid)
	}

	if len(jids) == 0 {
		return nil, fmt.Errorf("nenhum GROUP_JID valido encontrado")
	}
	return jids, nil
}

func (b *Bot) SendOffer(offer Offer) error {
	jids, err := b.resolveGroups(offer)
	if err != nil {
		return err
	}

	emoji := storeEmoji(offer.Store)
	name := storeName(offer.Store)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s *OFERTA %s* %s\n\n", emoji, name, emoji))
	sb.WriteString(fmt.Sprintf("📦 *%s*\n\n", offer.Title))

	if offer.OriginalPrice != "" && offer.Discount != "" {
		sb.WriteString(fmt.Sprintf("💰 De: ~%s~ por *%s* (-%s)\n", offer.OriginalPrice, offer.Price, offer.Discount))
	} else {
		sb.WriteString(fmt.Sprintf("💰 *%s*\n", offer.Price))
	}

	if offer.Link != "" {
		sb.WriteString(fmt.Sprintf("\n🔗 %s\n", offer.Link))
	}

	msg := &waE2E.Message{
		Conversation: proto.String(sb.String()),
	}

	for _, jid := range jids {
		if _, err := b.Client.SendMessage(context.Background(), jid, msg); err != nil {
			fmt.Printf("Erro ao enviar para %s: %v\n", jid, err)
		}
	}
	return nil
}
