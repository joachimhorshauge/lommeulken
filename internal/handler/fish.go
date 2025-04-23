package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type fishSpecies struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (h *Handler) HandleListSpecies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	species := []fishSpecies{
		{Value: "Aborre", Label: "Aborre"},
		{Value: "Brasen", Label: "Brasen"},
		{Value: "Bækørred", Label: "Bækørred"},
		{Value: "Fladfisk", Label: "Fladfisk"},
		{Value: "Gedde", Label: "Gedde"},
		{Value: "Græskarpe", Label: "Græskarpe"},
		{Value: "Havbars", Label: "Havbars"},
		{Value: "Havørred Kysten", Label: "Havørred Kysten"},
		{Value: "Havørred Åen", Label: "Havørred Åen"},
		{Value: "Helt", Label: "Helt"},
		{Value: "Hornfisk", Label: "Hornfisk"},
		{Value: "Laks", Label: "Laks"},
		{Value: "Makrel", Label: "Makrel"},
		{Value: "Multe", Label: "Multe"},
		{Value: "Pighvar/Slethvar", Label: "Pighvar/Slethvar"},
		{Value: "Put & Take ørred", Label: "Put & Take ørred"},
		{Value: "Rimte", Label: "Rimte"},
		{Value: "Sandart", Label: "Sandart"},
		{Value: "Skalle", Label: "Skalle"},
		{Value: "Skælkarpe", Label: "Skælkarpe"},
		{Value: "Spejlkarpe", Label: "Spejlkarpe"},
		{Value: "Suder", Label: "Suder"},
		{Value: "Søørred", Label: "Søørred"},
		{Value: "Torsk", Label: "Torsk"},
		{Value: "Uden for kategori", Label: "Uden for kategori"},
		{Value: "Ulk", Label: "Ulk"},
	}

	speciesJson, err := json.Marshal(species)
	if err != nil {
		slog.Error("Failed to marshal list of users", "Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(speciesJson)
}
