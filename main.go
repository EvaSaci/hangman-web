package main

// Importation des packages Go nécessaires pour diverses fonctionnalités
import (
<<<<<<< HEAD
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/common-nighthawk/go-figure"
)

// Déclaration de difficulte
type diff int

// Déclaration des constantes pour la difficulté
const (
	facile    diff = iota // mode facile
	moyen     diff = iota // mode moyen
	difficile diff = iota // mode difficile
	yann      diff = iota // mode yann : mode caché
	raciste   diff = iota // mode raciste : mode caché
	h              = iota // mode HELP
=======
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
>>>>>>> 0c146cc933a6d960a9877620bd98dc89ada9e875
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
<<<<<<< HEAD
	// Initialisation du générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())

	// Gestion du flag pour la difficulté
	difficultyFlag := flag.String("diff", "moyen", "Choisir la difficulté (facile, moyen, difficile)")
	flag.Parse()

	var fichier *os.File
	//	var f *os.File
	var err error

	// Gestion de la difficulté
	switch *difficultyFlag {
	case "facile": // mode facile
		fichier, err = os.Open("mots/motsFacile.txt") // ouvre le fichier des mots facile
		fmt.Println("mode de jeux : facile")
	case "moyen": // mode moyen
		fichier, err = os.Open("mots/motsNormal.txt") // ouvre le fichier des mots Moyen
		fmt.Println("mode de jeux : Moyen")
	case "difficile": // mode difficile
		fichier, err = os.Open("mots/motsHard.txt") // ouvre le fichier des mots Difficile
		fmt.Println("mode de jeux : Difficile")
	case "yann": // mode Yann
		fmt.Println("Mode de jeux : yann")
		fichier, err = os.Open("mots/motsYann.txt") // ouvre le fichier des mots de Yann
	case "raciste": // mode raciste
		fmt.Println("Mode de jeux : raciste")
		fichier, err = os.Open("mots/motsRaciste.txt") // ouvre le fichier des mots Raciste
	case "h": // mode help
		fmt.Println("mode HELP :")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("pour changer de mode : go run . -diff (facile; moyen; difficile)") // TEXTE du mode help
		fmt.Println("vous pouvez annulez le jeu a tous moment avec 'stop'")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("evidemment plusieurs easterEggs sont caché a vous de les trouvez !")
	default:
		fmt.Println("Difficulté inconnue. Utilisation de la difficulté moyenne par défaut.") // difficulté par default = moyen
		fichier, err = os.Open("mots/motsNormal.txt")                                        // ouvre le fichier des mots moyen
	}

	// Vérification d'erreurs lors de l'ouverture du fichier
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err) // messages d'erreur si le fichier est introuvable
		return
	}
	defer fichier.Close()

	// Lire les mots dans un tableau
	var mots []string
	scanner := bufio.NewScanner(fichier)
	for scanner.Scan() {
		mots = append(mots, scanner.Text())
	}

	// Vérifier s'il y a des mots dans la liste
	if len(mots) == 0 {
		fmt.Println("Le fichier ne contient aucun mot")
		return
	}

	// Choisir un mot aléatoirement
	mot := mots[rand.Intn(len(mots))]

	// Définir les vies en fonction de la difficulté
	var pv int
	switch *difficultyFlag {
	case "facile": // pour le mode facile
		pv = 10
	case "moyen": // pour le mode Moyen
		pv = 10
	case "difficile": // pour le mode difficile
		pv = 10
	case "yann": // pour le mode yann
		pv = 10
	case "raciste": // pour le mode raciste
		pv = 10
	default:
		pv = 10 // par défaut moyen
	}

	var test string               // Stocke la lettre ou mot entrée par le joueur
	estla := make(map[rune]bool)  // Lettres correctes
	estpas := make(map[rune]bool) // Lettres incorrectes
	VoF := true
	title := figure.NewFigure("Hangman !", "", true) // texte du début
	fmt.Println("----------------------------------------------------------------")
	title.Print()
	fmt.Println("----------------------------------------------------------------")
	//fmt.Println("		Bienvenue dans le jeu Hangman !")
	//fmt.Println("----------------------------------------------------------------")
	fmt.Println("")
	fmt.Println("plusieurs mode de jeux sont accessible")
	fmt.Println("--> go run . -diff h")
	fmt.Println("")
	fmt.Println("Prêt ?")
	fmt.Println("Tu peux annuler en écrivant 'stop'")
	fmt.Println("(Si tu veux pas y jouer t'es gay)")
	fmt.Println("----------------------------------------------------------------")

	for {
		cmt := 0
		cmt2 := 10 - pv + 1
		lignedebut := (10-pv)*8 + 1
		lignefin := 9 + lignedebut

		fmt.Print("Entrez une lettre ou un mot : ")
		fmt.Scan(&test)

		// Vérification si tous les caractères sont valides (lettres et '-')
		if !valide(test) {
			fmt.Println("Mauvais caractères, veuillez réessayer") // Message affiché si caractère non valide
			continue                                              // continue signifie que le jeux ne s'arrete pas
		}

		if test == "gay" { // conditions secrete qui arrête le jeux
			fmt.Println("la beuteu à Yann")
			break // break signifie la fin du programme
		}
		if test == "non" { // conditions secrete qui arrête le jeux
			fmt.Println("t'es gay")
			break
		}
		if test == "stop" { // conditions de stop du jeux
			fmt.Println("merci d'avoir joué !")
			fmt.Println("à bientôt !")
			break
		}

		// Si le joueur entre un mot entier
		if len(test) != 1 && test != mot {
			pv--
			pv-- // un mot faux = 2 pv en moins
			fmt.Printf("Mot incorrect ! Il vous reste %d chances\n", pv)
		} else if len(test) != 1 && test == mot {
			fmt.Printf("Bravo mec ! Tu as trouvé le mot : %s\n", mot)
			break // arrêt du jeux
		}

		// Si le joueur entre une seule lettre
		if len(test) == 1 {
			lettre := rune(test[0]) // Convertir la chaîne en rune
			if strings.ContainsRune(mot, lettre) {
				estla[lettre] = true
				VoF = false
			} else {
				if !estpas[lettre] {
					estpas[lettre] = true
					pv--
					fmt.Printf("La lettre %s n'est pas dans le mot\n", test)
					fmt.Printf("Il vous reste %d chances\n", pv)

					f, err := os.Open("O.txt") // ouvre le fichié du pendue
					if err != nil {
						//log.Fatal(err)
					}
					defer f.Close()
					scanner := bufio.NewScanner(f) // scanner pour lire le fichier ligne par ligne

					for scanner.Scan() { // Boucle pour parcourir chaque ligne du fichier
						cmt++

						if cmt >= lignedebut && cmt <= lignefin { // Si la ligne actuelle est dans l'intervalle spécifié
							fmt.Println(scanner.Text())
						}
						if cmt%(8*cmt2) == 0 {
							if VoF { // Si VoF est vrai, incrémente `cmt2`
								cmt2++
							}
							break
						}
					}
					if err := scanner.Err(); err != nil {
						//log.Fatal(err)
					}
				}
			}
		}
		// Afficher l'état actuel du mot
		fmt.Print("Mot à trouver : ") // Boucle à travers chaque caractère du mot à deviner
		for _, char := range mot {    // Si le caractère est présent dans le tableau `estla` (lettres devinées correctement)
			if estla[char] { // Affiche le caractère deviné
				fmt.Printf("%c ", char)
			} else {
				fmt.Print("_ ") // Sinon, affiche un underscore pour masquer les lettres non devinées
			}
		}
		fmt.Println()
		// Afficher les lettres incorrectes
		fmt.Print("Lettres incorrectes : ")
		for lettre := range estpas { // Affiche chaque lettre incorrecte devinée
			fmt.Printf("%c ", lettre)
		}
		fmt.Println()
		// Vérifier la victoire
		if checkWin(mot, estla) {
			// affichage spéciaux pour chaque style de victoire en fonction de la vie du joueur
			for pv >= 7 {
				fmt.Println("double Monstre pour avoir échoué moins de 3 fois !")
				break // arrêt du jeux
			}
			for pv == 6 || pv == 5 {
				fmt.Println("en vrais GG mais prochaine fois fait mieux stp !")
				break
			}
			for pv == 4 || pv == 3 {
				fmt.Println("Aie, Aie, Aie, c'était chaud la non ? Bravo quand même")
				break
			}
			for pv == 2 || pv == 1 {
				fmt.Println("ouais bon la soit t'es nul soit le mot été dur... GG quand même")
				break
			}
			title := figure.NewFigure("Victoire !", "", true) // texte de la victoire
			title.Print()                                     // affiché le texte de la victoire
			break
		}
		// Si le joueur perd
		if pv <= 0 {
			fmt.Println("Vous avez perdu")         // message de fin...
			fmt.Println("Le mot à trouver :", mot) // affiche le mot qui été à deviner
			break                                  // arrêt du jeux
		}
	}
}

// Fonction pour vérifier la victoire
func checkWin(mot string, estla map[rune]bool) bool {
	for _, char := range mot { // vérifie le mot
		if !estla[char] {
			return false
		}
	}
	return true
}

// coucou
// Fonction pour valider l'entrée (seulement lettres de a à z et le tiret '-')
func valide(input string) bool {
	for _, char := range input { // Vérifie si le caractère n'est ni une lettre ni un tiret '-'
		if !unicode.IsLetter(char) && char != '-' { // Si un caractère invalide est trouvé
			return false
		}
	}
	return true
=======
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
>>>>>>> 0c146cc933a6d960a9877620bd98dc89ada9e875
}
