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

}function showCategoriesModal() {
    // Implémentation du modal des catégories
}


function showCategoriesModal() {
    const categoriesContent = `
        <select id="category-select">
            ${Object.keys(words).map(category => `<option value="${category}">${category}</option>`).join('')}
        </select>
        <button id="apply-category">Appliquer</button>
    `;
    const modal = createModal("Choisir une catégorie", categoriesContent);
    modal.style.display = "block";

    document.getElementById("apply-category").onclick = function() {
        const selectedCategory = document.getElementById("category-select").value;
        initGame();
        modal.style.display = "none";
    }
}
package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Choisir une catégorie</title>
</head>
<body>
    <h2>{{.Title}}</h2>
    <form action="/" method="POST">
        <label for="category-select">Sélectionnez une catégorie :</label>
        <select name="category" id="category-select">
            {{range .Categories}}
            <option value="{{.}}">{{.}}</option>
            {{end}}
        </select>
        <button type="submit">Appliquer</button>
    </form>
</body>
</html>
`))

// Liste de catégories
var categories = []string{"Sport", "Technologie", "Art", "Science"}

// Struct pour les données de la page
type PageData struct {
	Title      string
	Categories []string
}

func main() {
	http.HandleFunc("/", categoryHandler)
	log.Println("Serveur démarré sur : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		selectedCategory := r.FormValue("category")
		// Appeler la fonction pour démarrer le jeu avec la catégorie sélectionnée
		initGame(selectedCategory)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Rendu de la page pour afficher le formulaire de sélection
	data := PageData{
		Title:      "Choisir une catégorie",
		Categories: categories,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Fonction pour initialiser le jeu avec la catégorie choisie
func initGame(category string) {
	log.Printf("Jeu initialisé avec la catégorie : %s", category)
	// Ajouter la logique de démarrage du jeu ici
}
*/

package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Choisir une catégorie</title>
</head>
<body>
    <h2>{{.Title}}</h2>
    <form action="/" method="POST">
        <label for="category-select">Sélectionnez une catégorie :</label>
        <select name="category" id="category-select">
            {{range .Categories}}
            <option value="{{.}}">{{.}}</option>
            {{end}}
        </select>
        <button type="submit">Appliquer</button>
    </form>
</body>
</html>
`))

// Liste de catégories
var categories = []string{"Sport", "Technologie", "Art", "Science"}

// Struct pour les données de la page
type PageData struct {
	Title      string
	Categories []string
}

func main() {
	http.HandleFunc("/", categoryHandler)
	log.Println("Serveur démarré sur : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		selectedCategory := r.FormValue("category")
		// Appeler la fonction pour démarrer le jeu avec la catégorie sélectionnée
		initGame(selectedCategory)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Rendu de la page pour afficher le formulaire de sélection
	data := PageData{
		Title:      "Choisir une catégorie",
		Categories: categories,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Fonction pour initialiser le jeu avec la catégorie choisie
func initGame(category string) {
	log.Printf("Jeu initialisé avec la catégorie : %s", category)
	// Ajouter la logique de démarrage du jeu ici
}
