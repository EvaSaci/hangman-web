// script en cours...
// script by Damien, Eva, Guilhem et Lois
// dead line : 12/11/2024
// avancement --> 75% /\ plus que à mettre en golang
// ici bas : script uniquement test et non final

// Définition des variables globales
let chosenWord = ""; // mot sélectionné pour la partie
let displayWord = ""; // mot affiché avec des lettres masquées
let essaie = 0; // nb d'essaie / fautes
let maxEssaie; // nb max d'essaie ou de faute
let isGameOver = false; // perdu ?
let usedLetters = new Set(); // Ajout d'un Set pour suivre les lettres utilisées
let wrongLetters = new Set(); // le set des lettres incorect

// Catégories des mots
const words = {
    superheros: ["superman", "batman", "spiderman", "ironman", "hulk", "thor", "wonder woman", "flash", "aquaman", "black panther", "deadpool", "captain america", "docteur strange", "green lantern", "wolverine", "cyclope", "magneto"],
    instruments: ["guitare", "piano", "violon", "batterie", "flute", "saxophone", "trompette", "clarinette", "harmonica", "ukulélé", "harpe", "banjo", "accordéon", "xylophone", "tuba", "violoncelle", "basson"],
    metiers: ["medecin", "avocat", "enseignant", "plombier", "policier", "ingenieur", "journaliste", "architecte", "pharmacien", "chirurgien", "psychologue", "chef cuisinier", "dentiste", "comptable", "vétérinaire", "électricien", "photographe"],
    autre: ["sérénité", "explosion", "gourmandise", "mystère", "illusion", "réflexion", "invisible", "paradoxe", "vitesse", "profondeur", "évasion", "mouvement", "équilibre", "renaissance", "discrétion", "éclair", "résonance", "crépuscule", "labyrinthe", "cascade", "liberté", "étincelle", "fragment", "mirage", "harmonie", "persévérance", "voyage", "abîme", "spirale", "souvenir", "synergie", "vertige", "élévation", "transformation", "murmure", "labyrinthe", "désir", "énergie", "sublime", "chimère", "lueur", "fantaisie"]
};

// Fonction d'initialisation du jeu modifiée
function initGame() {
    resetGame();
    const category = document.getElementById("category-select").value;
    const wordList = words[category];
    chosenWord = wordList[Math.floor(Math.random() * wordList.length)];
    displayWord = "_".repeat(chosenWord.length);
    essaie = 0;
    isGameOver = false;
    usedLetters.clear();
    wrongLetters.clear(); // Réinitialiser les lettres incorrectes
    updateDisplay();
}

// Fonction de mise à jour de l'affichage modifiée
function updateDisplay() {
    document.getElementById("word").textContent = displayWord.split("").join(" ");
    document.getElementById("hangman").getElementsByTagName("img")[0].src = `hangman${essaie}.png`;
    
    // Affiche les lettres incorrectes
    const wrongLettersContainer = document.getElementById("wrong-letters"); 
    if (wrongLettersContainer) {
        wrongLettersContainer.textContent = `Lettres incorrectes : ${Array.from(wrongLetters).join(', ')}`;
    }
    
    document.getElementById("message").textContent = isGameOver ? 
        (displayWord === chosenWord ? "Félicitations, vous avez gagné !" : `Vous avez perdu ! Le mot était : ${chosenWord}`) 
        : "";
}

// Fonction de devinette de lettre modifiée
function guessLetter(letter) {
    if (isGameOver || usedLetters.has(letter)) return;
    
    usedLetters.add(letter);
    let foundLetter = false;
    let tempDisplayWord = displayWord.split('');
    
    for (let i = 0; i < chosenWord.length; i++) {
        if (chosenWord[i].toLowerCase() === letter.toLowerCase()) {
            tempDisplayWord[i] = chosenWord[i];
            foundLetter = true;
        }
    }
    
    if (!foundLetter) {
        essaie++;
        wrongLetters.add(letter); // Ajouter la lettre aux lettres incorrectes
    }
    
    displayWord = tempDisplayWord.join('');
    updateDisplay();
    checkGameState();
}

// Fonction de vérification de l'état du jeu modifiée
function checkGameState() {
    if (displayWord === chosenWord) {
        document.getElementById("message").textContent = "Félicitations, vous avez gagné !";
        isGameOver = true;
        updateStats(true);
    } else if (essaie >= maxEssaie) {
        document.getElementById("message").textContent = `Vous avez perdu ! Le mot était : ${chosenWord}`;
        isGameOver = true;
        updateStats(false);
    }
}

// Modification de l'event listener pour les touches
document.addEventListener("keydown", function(event) {
    const letter = event.key.toLowerCase();
    if (/^[a-zàáâãäçèéêëìíîïñòóôõöùúûüýÿ]$/.test(letter)) {
        guessLetter(letter);
    }
});

function updateStats(won) {
    if (!currentUser) return;
    
    const userStats = users[currentUser].stats;
    userStats.gamesPlayed++;
    userStats.totalEssaie += essaie;
    
    if (won) {
        userStats.gamesWon++;
        userStats.currentStreak++;
        userStats.bestStreak = Math.max(userStats.bestStreak, userStats.currentStreak);
    } else {
        userStats.currentStreak = 0;
    }

    const category = document.getElementById("category-select").value;
    if (!userStats.categoryStats[category]) {
        userStats.categoryStats[category] = { played: 0, won: 0 };
    }
    userStats.categoryStats[category].played++;
    if (won) userStats.categoryStats[category].won++;

    updateStatsDisplay();
}

function animateHangman(essaie) {
    const parts = document.querySelectorAll('.hangman-part');
    for (let i = 0; i < essaie; i++) {
      if (parts[i]) {
        parts[i].style.animation = 'draw 0.5s linear forwards';
      }
    }
    
  }
  //let essaie = 0;

  function handleGuess(letter) {
    if (!wordToGuess.includes(letter)) {
      essaie++;
      animateHangman(essaie);
      
      if (essaie >= 10) {
        endGame('lose');
      }
    } else {
      // Logique pour une lettre correcte
    }
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
function animateHangman(essaie) {
    const parts = document.querySelectorAll('.hangman-part');
    for (let i = 0; i < essaie; i++) {
      if (parts[i]) {
        parts[i].style.animation = 'draw 0.5s linear forwards';
      }
    }
  }

// Event listeners et initialisation
document.addEventListener("DOMContentLoaded", function() {
    const MainLink = document.getElementById("main-link");
    const playLink = document.getElementById("play-link");
    const categoriesLink = document.getElementById("categories-link");
    const difficultyLink = document.getElementById("difficulty-link");
    const statsLink = document.getElementById("stats-link");
    const startGameBtn = document.getElementById("start-game");
    const resetBtn = document.getElementById("reset-btn");

    MainLink.addEventListener("click", showHeroSection)

    playLink.addEventListener("click", showGameSection);
    categoriesLink.addEventListener("click", showCategoriesModal);
    difficultyLink.addEventListener("click", showDifficultyModal);
    statsLink.addEventListener("click", showStatsModal);
    startGameBtn.addEventListener("click", startGame);
    resetBtn.addEventListener("click", resetGame);

    document.addEventListener("keydown", function(event) {
        const letter = event.key.toLowerCase();
        if (/^[a-z]$/.test(letter)) {
            guessLetter(letter);
        }
    });

    initGame();
});

function showHeroSection() {
    document.querySelector(".hero").style.display = "flex";
    document.getElementById("game-container").style.display = "none";
}

function showGameSection() {
    document.querySelector(".hero").style.display = "none";
    document.getElementById("game-container").style.display = "block";
}

function showCategoriesModal() {
    // Implémentation du modal des catégories
}

function showDifficultyModal() {
    // Implémentation du modal de difficulté
}

function showStatsModal() {
    // Implémentation du modal des statistiques
    updateStatsDisplay();
}

function startGame() {
    showGameSection();
    resetGame();
    updateDisplay();
}

function resetGame() {
    initGame();
}

// Implémentation des modals (suite)
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

function showStatsModal() {
    updateStatsDisplay();
    const statsContent = document.getElementById("stats-container").innerHTML;
    const modal = createModal("Statistiques", statsContent);
    modal.style.display = "block";
}

// Ajout de styles pour les modals
const style = document.createElement('style');
style.textContent = src=style.css
document.head.appendChild(style);

// Initialisation du jeu et des événements
document.addEventListener("DOMContentLoaded", function() {

    // Ajout de nouveaux écouteurs d'événements pour les boutons de modal
    document.getElementById("categories-link").addEventListener("click", showCategoriesModal);
    document.getElementById("difficulty-link").addEventListener("click", showDifficultyModal);
    document.getElementById("stats-link").addEventListener("click", showStatsModal);

    initGame();
});
