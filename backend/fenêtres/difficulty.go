/*
package main

function createModal(title, content) {
    const modal = document.createElement("div");
    modal.className = "modal";
    modal.innerHTML = `
        <div class="modal-content">
            <span class="close">&times;</span>
            <h2>${title}</h2>
            ${content}
        </div>
    `;
    document.body.appendChild(modal);

    const closeBtn = modal.querySelector(".close");
    closeBtn.onclick = function() {
        modal.style.display = "none";
    }

    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    }
    return modal;
}
function showDifficultyModal() {
    // Implémentation du modal de difficulté
}


function updateStatsDisplay() {
    if (!currentUser) return;
    
    const stats = users[currentUser].stats;
    const statsContainer = document.getElementById("stats-container");
    statsContainer.innerHTML = `
        <h3>Statistiques de ${currentUser}</h3>
        <p>Parties jouées : ${stats.gamesPlayed}</p>
        <p>Parties gagnées : ${stats.gamesWon}</p>
        <p>Pourcentage de victoires : ${((stats.gamesWon / stats.gamesPlayed) * 100).toFixed(2)}%</p>
        <p>Meilleure série : ${stats.bestStreak}</p>
        <p>Série actuelle : ${stats.currentStreak}</p>
        <p>Erreurs moyennes par partie : ${(stats.totalEssaie / stats.gamesPlayed).toFixed(2)}</p>
        <h4>Statistiques par catégorie :</h4>
    `;
    for (const [category, catStats] of Object.entries(stats.categoryStats)) {
        statsContainer.innerHTML += `
            <p>${category}: ${catStats.won}/${catStats.played} (${((catStats.won / catStats.played) * 100).toFixed(2)}%)</p>
        `;
    }
}
function showDifficultyModal() {
    const difficultyContent = `
        <button id="easy-mode">Facile</button>
        <button id="medium-mode">Moyen</button>
        <button id="hard-mode">Difficile</button>
    `;
    const modal = createModal("Choisir la difficulté", difficultyContent);
    modal.style.display = "block";

    document.getElementById("easy-mode").onclick = function() {
        maxEssaie = 8;
        modal.style.display = "none";
        resetGame();
    }
    document.getElementById("medium-mode").onclick = function() {
        maxEssaie = 6;
        modal.style.display = "none";
        resetGame();
    }
    document.getElementById("hard-mode").onclick = function() {
        maxEssaie = 4;
        modal.style.display = "none";
        resetGame();
    }
}

document.addEventListener("DOMContentLoaded", function() {

    document.getElementById("difficulty-link").addEventListener("click", showDifficultyModal);

    initGame();
});
*/
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Structures pour les statistiques de l'utilisateur et les paramètres de difficulté
type Stats struct {
	GamesPlayed   int
	GamesWon      int
	BestStreak    int
	CurrentStreak int
	TotalEssaie   int
	CategoryStats map[string]CategoryStats
}

type CategoryStats struct {
	Played int
	Won    int
}

type User struct {
	Name  string
	Stats Stats
}

var users = map[string]User{
	"Joueur1": {
		Name: "Joueur1",
		Stats: Stats{
			GamesPlayed: 10,
			GamesWon:    5,
			BestStreak:  3,
			CurrentStreak: 2,
			TotalEssaie: 50,
			CategoryStats: map[string]CategoryStats{
				"Sport":     {Played: 5, Won: 3},
				"Technologie": {Played: 3, Won: 2},
			},
		},
	},
}

var currentUser = "Joueur1"
var maxEssaie int

var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Statistiques et Difficulté</title>
</head>
<body>
    <h2>Bienvenue, {{.User.Name}}</h2>

    <h3>Statistiques de {{.User.Name}}</h3>
    <p>Parties jouées : {{.User.Stats.GamesPlayed}}</p>
    <p>Parties gagnées : {{.User.Stats.GamesWon}}</p>
    <p>Pourcentage de victoires : {{printf "%.2f" .WinPercentage}}%</p>
    <p>Meilleure série : {{.User.Stats.BestStreak}}</p>
    <p>Série actuelle : {{.User.Stats.CurrentStreak}}</p>
    <p>Erreurs moyennes par partie : {{printf "%.2f" .AvgErrorsPerGame}}</p>

    <h4>Statistiques par catégorie :</h4>
    {{range $category, $stats := .User.Stats.CategoryStats}}
        <p>{{$category}}: {{$stats.Won}}/{{$stats.Played}} ({{printf "%.2f" (calcPercentage $stats.Won $stats.Played)}}%)</p>
    {{end}}

    <h3>Choisir une difficulté</h3>
    <form action="/set-difficulty" method="POST">
        <button name="difficulty" value="easy">Facile</button>
        <button name="difficulty" value="medium">Moyen</button>
        <button name="difficulty" value="hard">Difficile</button>
    </form>
</body>
</html>
`))

// Struct pour transmettre les données au template
type PageData struct {
	User            User
	WinPercentage   float64
	AvgErrorsPerGame float64
}

// Fonction principale
func main() {
	http.HandleFunc("/", statsHandler)
	http.HandleFunc("/set-difficulty", difficultyHandler)
	log.Println("Serveur démarré sur : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler pour afficher les statistiques et les options de difficulté
func statsHandler(w http.ResponseWriter, r *http.Request) {
	user := users[currentUser]
	data := PageData{
		User:            user,
		WinPercentage:   float64(user.Stats.GamesWon) / float64(user.Stats.GamesPlayed) * 100,
		AvgErrorsPerGame: float64(user.Stats.TotalEssaie) / float64(user.Stats.GamesPlayed),
	}

	// Fonction pour calculer le pourcentage de victoire par catégorie
	funcMap := template.FuncMap{
		"calcPercentage": func(won, played int) float64 {
			if played == 0 {
				return 0
			}
			return float64(won) / float64(played) * 100
		},
	}
	tmpl = tmpl.Funcs(funcMap)
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler pour définir la difficulté
func difficultyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		difficulty := r.FormValue("difficulty")
		switch difficulty {
		case "easy":
			maxEssaie = 8
		case "medium":
			maxEssaie = 6
		case "hard":
			maxEssaie = 4
		}
		fmt.Printf("Difficulté sélectionnée: %s, max essaie: %d\n", difficulty, maxEssaie)
		resetGame()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Fonction pour réinitialiser le jeu
func resetGame() {
	log.Println("Jeu réinitialisé avec la difficulté:", maxEssaie)
	// Logique de réinitialisation du jeu ici
}
