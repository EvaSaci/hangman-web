***Hangman Web***

*Fonctionnalités*

	•	Jeu classique du pendu : Devinez des mots aléatoires en fonction de différents niveaux de difficulté.
	•	Interface utilisateur : Utilisation de HTML et CSS pour une présentation simple et claire.
	•	Backend performant : Développé en Go pour une exécution rapide.
	•	Multiniveau : Inclut des mots de difficultés facile, moyenne et difficile.

*Structure du projet*

	•	main.go : Contient le backend en Go, qui gère la logique du jeu et le serveur HTTP.
	•	index.html : Interface utilisateur principale pour jouer au jeu.
	•	static/ : Contient les styles CSS et les images pour afficher les états du jeu (pendu, vies restantes, etc.).
	•	Données : Fichiers texte (words_easy.txt, words_medium.txt, words_hard.txt) contenant les listes de mots pour chaque niveau de difficulté.

*Prérequis*

	•	Go installé sur votre système.
	•	Navigateur web moderne.


*Installation et exécution*

1. Clonez le dépôt

Clonez ce dépôt sur votre machine locale avec la commande suivante :

	•	git clone https://github.com/EvaSaci/hangman-web.git
	•	cd hangman-web

2. Installez les dépendances

Assurez-vous que Go est installé sur votre système. Si ce n’est pas déjà fait, téléchargez-le et installez-le depuis golang.org.

3. Lancez le serveur

Exécutez le fichier principal pour démarrer le serveur HTTP :

	•	go run main.go

4. Accédez au jeu

Ouvrez votre navigateur et rendez-vous sur :

	•	http://localhost:8080

Vous êtes prêt à jouer au pendu !