package supabase

import (
	"github.com/nedpals/supabase-go"
	"os"
)

var Client *supabase.Client

func init() {
	sbHost := os.Getenv("SUPABASE_HOST")
	sbSecret := os.Getenv("SUPABASE_SECRET")
	Client = supabase.CreateClient(sbHost, sbSecret)
}
