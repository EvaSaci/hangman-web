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
	Word            string
	MaskedWord      string
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
var wordLists map[string][]string // Map pour stocker les mots par niveau

func loadWordsFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture du fichier %s : %v", filename, err)
		return []string{}
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Erreur lors de la lecture du fichier %s : %v", filename, err)
	}

	return words
}
func loadAllWordLists() {
	wordLists = make(map[string][]string)
	wordLists["facile"] = loadWordsFromFile("words_easy.txt")
	wordLists["moyen"] = loadWordsFromFile("words_medium.txt")
	wordLists["difficile"] = loadWordsFromFile("words_hard.txt")

	// Si aucun mot n'est chargé pour un niveau, utiliser une liste par défaut
	if len(wordLists["facile"]) == 0 {
		wordLists["facile"] = []string{"chat", "chien", "oiseau", "lapin"}
	}
	if len(wordLists["moyen"]) == 0 {
		wordLists["moyen"] = []string{"informatique", "programme", "serveur", "reseau"}
	}
	if len(wordLists["difficile"]) == 0 {
		wordLists["difficile"] = []string{"developpement", "infrastructure", "algorithme"}
	}
}


func initGame(difficulty string) {
	if !game.GameInitialized || game.Difficulty != difficulty {
		rand.Seed(time.Now().UnixNano())
		words := wordLists[difficulty]
		if len(words) == 0 {
			words = wordLists["moyen"] // Fallback au niveau moyen si le niveau demandé est vide
		}
		game.Word = words[rand.Intn(len(words))]
		game.MaxTries = 6
		game.RemainingTries = game.MaxTries
		game.GuessedLetters = []string{}
		game.WrongLetters = []string{}
		game.MaskedWord = maskWord(game.Word)
		game.GameStatus = "en_cours"
		game.RevealedLetter = revealRandomLetter()
		game.GameInitialized = true
		game.Difficulty = difficulty
	}
}

func maskWord(word string) string {
	masked := make([]rune, len(word))
	for i := range masked {
		if word[i] == ' ' || word[i] == '_' {
			masked[i] = rune(word[i])
		} else {
			masked[i] = '_'
		}
	}
	return string(masked)
}

func revealRandomLetter() string {
	letterIndexes := []int{}
	for i, char := range game.Word {
		if char != ' ' && char != '_' {
			letterIndexes = append(letterIndexes, i)
		}
	}

	if len(letterIndexes) > 0 {
		randomIndex := letterIndexes[rand.Intn(len(letterIndexes))]
		revealedLetter := strings.ToLower(string(game.Word[randomIndex]))

		maskedRunes := []rune(game.MaskedWord)
		for i, char := range game.Word {
			if strings.ToLower(string(char)) == revealedLetter {
				maskedRunes[i] = char
			}
		}
		game.MaskedWord = string(maskedRunes)

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

	if strings.Contains(strings.ToLower(game.Word), letter) {
		updateMaskedWord(letter)
		if game.MaskedWord == game.Word {
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

func updateMaskedWord(letter string) {
	masked := []rune(game.MaskedWord)
	for i, char := range game.Word {
		if strings.ToLower(string(char)) == letter {
			masked[i] = char
		}
	}
	game.MaskedWord = string(masked)
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
