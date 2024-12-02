package main

// Importation des packages Go nécessaires pour diverses fonctionnalités
import (
	"bufio" // Pour lire des fichiers ligne par ligne
	"fmt"
	"html/template" // Pour rendre des modèles HTML
	"log"           // Pour enregistrer des erreurs et des informations
	"math/rand"     // Pour générer des nombres aléatoires
	"net/http"      // Pour créer un serveur web et gérer les requêtes HTTP
	"os"            // Pour les opérations sur les fichiers
	"strings"       // Pour la manipulation de chaînes de caractères
	"time"          // Pour initialiser le générateur de nombres aléatoires
	"unicode"
)

// Structure HangmanGame définit l'état et les propriétés d'un jeu du Pendu
type HangmanGame struct {
	Mots            string   // Le mot à deviner
	Motsmasque      string   // Version masquée du mot avec des underscores
	RemainingTries  int      // Nombre de tentatives restantes
	MaxTries        int      // Nombre maximum de tentatives autorisées
	GuessedLetters  []string // Lettres qui ont été devinées
	WrongLetters    []string // Lettres qui ont été devinées incorrectement
	GameStatus      string   // Statut actuel du jeu (en_cours, gagne, perdu)
	RevealedLetter  string   // Une lettre révélée au début du jeu
	GameInitialized bool     // Indicateur pour vérifier si le jeu est initialisé
	Difficulty      string   // Niveau de difficulté du jeu
	State           string
}

// Variables globales pour le jeu
var game HangmanGame
var MotsLists map[string][]string // Map pour stocker les mots par niveau

// Fonction pour charger des mots à partir d'un fichier
func loadMotFromFile(filename string) []string {
	// Ouverture du fichier
	file, err := os.Open(filename)
	if err != nil {
		// Journalisation de l'erreur si le fichier ne peut pas être ouvert
		log.Printf("Erreur lors de l'ouverture du fichier %s : %v", filename, err)
		return []string{}
	}
	defer file.Close() // Fermeture du fichier à la fin de la fonction

	var Mot []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Nettoyage et ajout des mots non vides à la liste
		Mots := strings.TrimSpace(scanner.Text())
		if Mots != "" {
			Mot = append(Mot, Mots)
		}
	}

	// Gestion des erreurs potentielles lors de la lecture
	if err := scanner.Err(); err != nil {
		log.Printf("Erreur lors de la lecture du fichier %s : %v", filename, err)
	}

	return Mot
}

// Fonction pour charger tous les fichiers de mots par difficulté
func loadAllWordLists() {
	// Initialisation de la map pour stocker les listes de mots
	MotsLists = make(map[string][]string)
	// Chargement des mots pour chaque niveau de difficulté
	MotsLists["facile"] = loadMotFromFile("words_easy.txt")
	MotsLists["moyen"] = loadMotFromFile("words_medium.txt")
	MotsLists["difficile"] = loadMotFromFile("words_hard.txt")

	// Listes de mots par défaut si aucun mot n'est chargé
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

// Fonction pour initialiser le jeu avec un niveau de difficulté
func initGame(difficulty string) {
	// Initialisation du jeu uniquement si nécessaire
	if !game.GameInitialized || game.Difficulty != difficulty {
		// Initialisation du générateur de nombres aléatoires
		rand.Seed(time.Now().UnixNano())

		// Sélection de la liste de mots selon la difficulté
		Mots := MotsLists[difficulty]
		if len(Mots) == 0 {
			// Utilisation de la liste de mots de niveau moyen comme secours
			Mots = MotsLists["moyen"]
		}

		// Sélection aléatoire d'un mot
		game.Mots = Mots[rand.Intn(len(Mots))]

		// Configuration des paramètres du jeu
		game.MaxTries = 6
		game.RemainingTries = game.MaxTries
		game.GuessedLetters = []string{}
		game.WrongLetters = []string{}

		// Masquage du mot
		game.Motsmasque = maskWord(game.Mots)

		// Définition du statut du jeu
		game.GameStatus = "en_cours"

		// Révélation d'une lettre aléatoire
		game.RevealedLetter = revealRandomLetter()

		// Marquage du jeu comme initialisé
		game.GameInitialized = true

		// Enregistrement du niveau de difficulté
		game.Difficulty = difficulty
	}
}

// Fonction pour masquer un mot en remplaçant les lettres par des underscores
func maskWord(Mots string) string {
	// Conversion du mot en tableau de caractères
	masked := make([]rune, len(Mots))
	for i := range masked {
		// Conserver les espaces et underscores existants
		if Mots[i] == ' ' || Mots[i] == '_' {
			masked[i] = rune(Mots[i])
		} else {
			// Remplacer les autres caractères par des underscores
			masked[i] = '_'
		}
	}
	return string(masked)
}

// Fonction pour révéler une lettre aléatoire du mot
func revealRandomLetter() string {
	// Collecte des index des lettres (excluant espaces et underscores)
	letterIndexes := []int{}
	for i, char := range game.Mots {
		if char != ' ' && char != '_' {
			letterIndexes = append(letterIndexes, i)
		}
	}

	// Si des lettres sont disponibles
	if len(letterIndexes) > 0 {
		// Sélection d'un index aléatoire
		randomIndex := letterIndexes[rand.Intn(len(letterIndexes))]
		// Conversion de la lettre en minuscules
		revealedLetter := strings.ToLower(string(game.Mots[randomIndex]))

		// Mise à jour du mot masqué
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
func valide(input string) bool {
	for _, char := range input {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

// Fonction pour gérer une tentative de devinette
func handleGuess(letter string) {
	// Conversion de la lettre en minuscules
	letter = strings.ToLower(letter)

	if !valide(letter) {
		//game.GameStatus = "test"
		return // Ignorer les caractères spéciaux sans faire perdre de vie
	}
	// Vérification si la lettre a déjà été devinée
	if contains(game.GuessedLetters, letter) {
		fmt.Printf("ttt")
	}

	// Ajout de la lettre à la liste des lettres devinées
	game.GuessedLetters = append(game.GuessedLetters, letter)

	// Vérification si la lettre est dans le mot
	if strings.Contains(strings.ToLower(game.Mots), letter) {
		// Mise à jour du mot masqué
		updateMotsmasque(letter)
		// Vérification si le mot est completement deviné
		if game.Motsmasque == game.Mots {
			game.GameStatus = "gagne"
			game.RemainingTries = game.MaxTries
		}
	} else {
		// Réduction du nombre de tentatives
		game.RemainingTries--
		// Ajout de la lettre aux lettres incorrectes
		game.WrongLetters = append(game.WrongLetters, letter)
		// Vérification de la défaite
		if game.RemainingTries == 0 {
			game.GameStatus = "perdu"
		}
	}
}

// Fonction pour mettre à jour le mot masqué avec une lettre devinée
func updateMotsmasque(letter string) {
	// Conversion du mot masqué en tableau de caractères
	masked := []rune(game.Motsmasque)
	for i, char := range game.Mots {
		// Remplacement des underscores par les lettres correctes
		if strings.ToLower(string(char)) == letter {
			masked[i] = char
		}
	}
	// Mise à jour du mot masqué
	game.Motsmasque = string(masked)
}

// Fonction utilitaire pour vérifier si un élément est dans un tableau
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Gestionnaire de la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Gestion de la requête GET
	if r.Method == "GET" {

		hangEv()
		fmt.Println("state")
		fmt.Println(game.State)
		fmt.Println(game.RemainingTries)
		// Récupération du niveau de difficulté
		difficulty := r.URL.Query().Get("difficulty")

		// Si aucune difficulté n'est spécifiée
		if difficulty == "" {
			// Si un jeu est déjà en cours, conserver son niveau de difficulté
			if game.GameInitialized {
				difficulty = game.Difficulty
			} else {
				// Sinon, utiliser le niveau moyen par défaut
				difficulty = "moyen"
			}
		}
		log.Printf("Game Status: %s", game.GameStatus)
		log.Printf("Difficulty: %s", difficulty)
		// Vérification pour réinitialiser le jeu
		if !game.GameInitialized ||
			game.GameStatus == "gagne" ||
			game.GameStatus == "perdu" ||
			game.Difficulty != difficulty {
			initGame(difficulty)
		}
		// In your handleGuess function
		if game.Motsmasque == game.Mots {
			game.GameStatus = "gagne" // Ensure it's lowercase
			game.RemainingTries = game.MaxTries
		}

		// And when the player loses
		if game.RemainingTries == 0 {
			game.GameStatus = "perdu" // Ensure it's lowercase
		}
	}

	// Chargement et exécution du template HTML
	tmpl := template.Must(template.ParseFiles("index.html"))
	err := tmpl.Execute(w, game)
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Gestionnaire pour les tentatives de devinettes
func guessHandler(w http.ResponseWriter, r *http.Request) {
	// Gestion de la requête POST
	if r.Method == "POST" {
		hangEv()
		r.ParseForm()
		// Récupération de la lettre devinée
		letter := r.Form.Get("letter")
		// Traitement de la devinette
		handleGuess(letter)
		// Redirection vers la page principale
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Fonction principale pour démarrer le serveur
func main() {
	// Chargement des listes de mots
	loadAllWordLists()

	// Configuration des fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fi := http.FileServer(http.Dir("image")) // Serveur pour les images
	http.Handle("/image/",http.StripPrefix("/image/", fi)) // Serveur pour les images


	// Configuration des routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/guess", guessHandler)

	// Démarrage du serveur
	log.Println("Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Erreur serveur:", err)
	}
}

func hangEv() {
	switch game.RemainingTries {
	case 6:
		game.State = "./image/hangEv/vie.png" // VIE
	case 5:
		game.State = "./image/hangEv/vie-1.png"
	case 4:
		game.State = "./image/hangEv/vie-2.png"
	case 3:
		game.State = "./image/hangEv/vie-3.png"
	case 2:
		game.State = "./image/hangEv/vie-4.png"
	case 1:
		game.State = "./image/hangEv/vie-5.png"
	case 0:
		game.State = "./image/hangEv/perdu.png"
	}
}
