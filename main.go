package main

import (
	"bufio"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type HangmanGame struct {
	Mots            string
	Motsmasque      string
	RemainingTries  int
	MaxTries        int
	GuessedLetters  []string
	WrongLetters    []string
	GameStatus      string
	RevealedLetter  string
	GameInitialized bool
	Difficulty      string // Nouveau champ pour la difficulté
}

var game HangmanGame
var MotsLists map[string][]string // Map pour stocker les mots par niveau

func loadMotFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture du fichier %s : %v", filename, err)
		return []string{}
	}
	defer file.Close()

	var Mot []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Mots := strings.TrimSpace(scanner.Text())
		if Mots != "" {
			Mot = append(Mot, Mots)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Erreur lors de la lecture du fichier %s : %v", filename, err)
	}

	return Mot
}
func loadAllWordLists() {
	MotsLists = make(map[string][]string)
	MotsLists["facile"] = loadMotFromFile("words_easy.txt")
	MotsLists["moyen"] = loadMotFromFile("words_medium.txt")
	MotsLists["difficile"] = loadMotFromFile("words_hard.txt")

	// Si aucun mot n'est chargé pour un niveau, utiliser une liste par défaut
	if len(MotsLists["facile"]) == 0 {
		MotsLists["facile"] = []string{"chat", "chien", "oiseau", "lapin"}
	}
	if len(MotsLists["moyen"]) == 0 {
		MotsLists["moyen"] = []string{"informatique", "programme", "serveur", "reseau"}
	}
	if len(MotsLists["difficile"]) == 0 {
		MotsLists["difficile"] = []string{"developpement", "infrastructure", "algorithme"}
	}
}


func initGame(difficulty string) {
	if !game.GameInitialized || game.Difficulty != difficulty {
		rand.Seed(time.Now().UnixNano())
		Mots := MotsLists[difficulty]
		if len(Mots) == 0 {
			Mots = MotsLists["moyen"] // Fallback au niveau moyen si le niveau demandé est vide
		}
		game.Mots = Mots[rand.Intn(len(Mots))]
		game.MaxTries = 6
		game.RemainingTries = game.MaxTries
		game.GuessedLetters = []string{}
		game.WrongLetters = []string{}
		game.Motsmasque = maskWord(game.Mots)
		game.GameStatus = "en_cours"
		game.RevealedLetter = revealRandomLetter()
		game.GameInitialized = true
		game.Difficulty = difficulty
	}
}

func maskWord(Mots string) string {
	masked := make([]rune, len(Mots))
	for i := range masked {
		if Mots[i] == ' ' || Mots[i] == '_' {
			masked[i] = rune(Mots[i])
		} else {
			masked[i] = '_'
		}
	}
	return string(masked)
}

func revealRandomLetter() string {
	letterIndexes := []int{}
	for i, char := range game.Mots {
		if char != ' ' && char != '_' {
			letterIndexes = append(letterIndexes, i)
		}
	}

	if len(letterIndexes) > 0 {
		randomIndex := letterIndexes[rand.Intn(len(letterIndexes))]
		revealedLetter := strings.ToLower(string(game.Mots[randomIndex]))

		maskedRunes := []rune(game.Motsmasque)
		for i, char := range game.Mots {
			if strings.ToLower(string(char)) == revealedLetter {
				maskedRunes[i] = char
			}
		}
		game.Motsmasque = string(maskedRunes)

		return revealedLetter
	}

	return ""
}

func handleGuess(letter string) {
	letter = strings.ToLower(letter)

	if contains(game.GuessedLetters, letter) {
		return
	}

	game.GuessedLetters = append(game.GuessedLetters, letter)

	if strings.Contains(strings.ToLower(game.Mots), letter) {
		updateMotsmasque(letter)
		if game.Motsmasque == game.Mots {
			game.GameStatus = "gagne"
			game.RemainingTries = game.MaxTries
		}
	} else {
		game.RemainingTries--
		game.WrongLetters = append(game.WrongLetters, letter)
		if game.RemainingTries == 0 {
			game.GameStatus = "perdu"
		}
	}
}

func updateMotsmasque(letter string) {
	masked := []rune(game.Motsmasque)
	for i, char := range game.Mots {
		if strings.ToLower(string(char)) == letter {
			masked[i] = char
		}
	}
	game.Motsmasque = string(masked)
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		difficulty := r.URL.Query().Get("difficulty")
		if difficulty == "" {
			difficulty = "moyen" // Niveau par défaut
		}
		
		if game.GameStatus == "gagne" || game.GameStatus == "perdu" || game.Difficulty != difficulty {
			game.GameInitialized = false
		}
		
		initGame(difficulty)
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	err := tmpl.Execute(w, game)
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		letter := r.Form.Get("letter")
		handleGuess(letter)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func main() {
	loadAllWordLists()
	// Chargement du fichier de mots

	// Configuration des fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/guess", guessHandler)

	// Démarrage du serveur avec logs
	log.Println("Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Erreur serveur:", err)
	}
}
