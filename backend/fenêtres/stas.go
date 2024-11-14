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

function showStatsModal() {
    // Implémentation du modal des statistiques
    updateStatsDisplay();
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

function showStatsModal() {
    updateStatsDisplay();
    const statsContent = document.getElementById("stats-container").innerHTML;
    const modal = createModal("Statistiques", statsContent);
    modal.style.display = "block";
}
document.addEventListener("DOMContentLoaded", function() {

    // Ajout de nouveaux écouteurs d'événements pour les boutons de modal
    document.getElementById("stats-link").addEventListener("click", showStatsModal);

    initGame();
});
*/
package main

import (
    "fmt"
    "html/template"
    "net/http"
)

type Stats struct {
    GamesPlayed    int
    GamesWon       int
    BestStreak     int
    CurrentStreak  int
    TotalEssaie    int
    CategoryStats  map[string]CategoryStat
}

type CategoryStat struct {
    Won   int
    Played int
}

var currentUser string
var users map[string]User

type User struct {
    Stats Stats
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/stats", statsHandler)
    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("index.html"))
    tmpl.Execute(w, nil)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
    if currentUser == "" {
        http.Error(w, "No current user", http.StatusBadRequest)
        return
    }

    stats := users[currentUser].Stats
    statsHTML := fmt.Sprintf(`
        <h3>Statistiques de %s</h3>
        <p>Parties jouées : %d</p>
        <p>Parties gagnées : %d</p>
        <p>Pourcentage de victoires : %.2f%%</p>
        <p>Meilleure série : %d</p>
        <p>Série actuelle : %d</p>
        <p>Erreurs moyennes par partie : %.2f</p>
        <h4>Statistiques par catégorie :</h4>
    `, currentUser, stats.GamesPlayed, stats.GamesWon, float64(stats.GamesWon)/float64(stats.GamesPlayed)*100, stats.BestStreak, stats.CurrentStreak, float64(stats.TotalEssaie)/float64(stats.GamesPlayed))

    for category, catStats := range stats.CategoryStats {
        statsHTML += fmt.Sprintf(`
            <p>%s: %d/%d (%.2f%%)</p>
        `, category, catStats.Won, catStats.Played, float64(catStats.Won)/float64(catStats.Played)*100)
    }

    w.Write([]byte(statsHTML))
}