# Monday-light

projet perso pour apprendre go template avec htmx

method / getter / setter / interface
rajouter des err

hx-on / push url
quand je creer une categorie, penser a refresh aussi le dropdown de creation de tache

## tech debt

coleur du profil limiter a 10 user
if/else dans base.html
param ne sont pas assez modulaire ( il faudrait que chaque champ soit un element d'une liste de un template, ca enleverai le html brut de renderUpdatedField)

## To do

### 1
when add remove tag (refresh dropdown for tags tache)


faire les tasks
design vertical pour les tables
plus de champ a ajouter une tache
creation de la tache ainsi que visu de la tache dans un pop up
dropdown menu for categorie in a task
bar de progression dans le temps reel (front) calculer a partir de taken from
avoir un dashboard user avec leur tache prise

it lack data validation, retour visuelle, confirmation

### 2
searchbar ou trie par bouton genre par tags
chat sur les taches
fichiers dans les taches
suppression de categorie
dashboard avancement du projet ainsi que global (tous les projets)
faire une vu par categorie
pouvoir creer des sous taches
faire une vue par status perso des taches

### 3

proteger route individuel comme modif de compte
changement url et retour en arriere et refresh
gros boulot de front

### 4

connexion via 42 api
bot discord pour prevenir des notifs
swagger
traefik
mvc model
validation data !!!!!
ginkgo gomega
bonne balise html de navigation


do i check if the account exist when checking jwt ????
mdp hacher avant envoie de front a back ?
BUG : si token avec id inconnue, db error

## to do pour page param

plus jolie pop pour le changement de username
icon modify qui se decale apres un update reussit
afficher le succes du changement de mdp

## spa et srr

pour lajout ou la suppresion delement en direct
https://youtu.be/To-Mlm2AwD4?si=Wq2osxv2laQAy1mF
https://youtu.be/To-Mlm2AwD4?si=pJdAwja2CFuuMOF2


# readme a ajouter
connection jwt
mdp hasher
spa ssr

## Cheat sheet

### Go Dependencies

```
go mod init monday-light
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go mod tidy
```

### HTMX

| Attribute         | Example Value             | Description                                                                 |
|-------------------|---------------------------|-----------------------------------------------------------------------------|
| `hx-get`         | `/get-example`            | Sends a GET request to fetch data from the server.                         |
| `hx-post`        | `/post-example`           | Sends a POST request with data to the server.                              |
| `hx-target`      | `#result`                 | Specifies the element to update with the server response.                  |
| `hx-trigger`     | `click`                   | Triggers the request when the element is clicked.                          |
| `hx-trigger`     | `every 5s`                | Triggers the request periodically, every 5 seconds.                        |
| `hx-trigger`     | `load`                    | Triggers the request when the page loads.                                  |
| `hx-trigger`     | `change`                  | Triggers the request when the value of an input changes.                   |
| `hx-swap`        | `innerHTML`               | Replaces the inner HTML of the target element.                             |
| `hx-swap`        | `outerHTML`               | Replaces the entire target element with the response.                      |
| `hx-swap`        | `beforeend`               | Appends the response content to the end of the target element.             |
| `hx-swap`        | `afterbegin`              | Prepends the response content to the start of the target element.          |
| `hx-swap`        | `beforebegin`             | Inserts the response content before the target element.                    |
| `hx-swap`        | `afterend`                | Inserts the response content after the target element.                     |
| `hx-on`          | `htmx:responseError`      | Executes JavaScript when a server error occurs.                            |
| `hx-vals`        | `{ 'key': 'value' }`      | Sends additional parameters with the request.                              |
| `hx-headers`     | `{ 'Authorization': '...' }` | Adds custom headers to the request.                                       |
| `hx-indicator`   | `#loading`                | Shows a loading indicator during the request.                              |
| `hx-push-url`    | `true`                    | Pushes the request URL into the browser history.                           |
| `hx-select`      | `.row`                    | Selects specific elements from the server response.                        |
| `hx-select-oob`  | `.notification`           | Processes out-of-band (OOB) elements from the response.                    |
| `hx-confirm`     | `Are you sure?`           | Shows a confirmation dialog before sending the request.                    |
| `hx-disable`     | `true`                    | Disables the element to prevent duplicate requests.                        |
| `hx-history`     | `false`                   | Disables browser history for this request.                                 |
| `hx-include`     | `#form-id`                | Includes additional elements with the request.                             |
| `hx-preserve`    | `true`                    | Preserves the target element instead of replacing it.                      |
| `hx-replace-url` | `/new-url`                | Replaces the current browser URL with the new URL.                         |

### GIN + Template

```
// Route pour générer les lignes dynamiques
r.GET("/generate-rows", func(c *gin.Context) {
    // Générer des lignes aléatoires
    rows := make([]map[string]interface{}, 5)
    for i := 0; i < 5; i++ {
        rows[i] = map[string]interface{}{
            "ID":    i + 1,
            "Value": rand.Intn(100), // Génère une valeur aléatoire
        }
    }
    c.HTML(http.StatusOK, "rows.html", gin.H{
        "Rows": rows,
    })
})


{{range .Rows}}
<tr>
    <td>{{.ID}}</td>
    <td>{{.Value}}</td>
</tr>
{{end}}
```

#### Fonction Gin

| Fonction                  | Description                                    |
|---------------------------|------------------------------------------------|
| `gin.Default()`           | Initialise un routeur avec middlewares.       |
| `r.LoadHTMLGlob()`        | Charge des templates HTML depuis un dossier.  |
| `r.Static()`              | Sert des fichiers statiques.                  |
| `r.GET()`                 | Crée une route avec la méthode GET.           |
| `c.HTML()`                | Rendu HTML avec des données dynamiques.       |

#### Syntaxe Template

| Syntaxe                   | Description                                    |
|---------------------------|------------------------------------------------|
| `{{ .Key }}`              | Affiche une variable envoyée depuis Gin.       |
| `{{range .List}} ... {{end}}` | Parcourt une liste et affiche son contenu.    |

### Air

```
# Configuration du fichier `.air.toml`

Le fichier `.air.toml` est utilisé pour configurer Air, un outil de rechargement à chaud pour les projets Go. Voici une explication détaillée de tous les arguments disponibles :

## Section `[build]`

### `cmd`
- **Description** : La commande pour compiler votre application.
- **Exemple** : `"go build -o ./tmp/main ."`
- **Détails** : Compile l'application en utilisant `go build` et génère un binaire dans le répertoire spécifié.

---

### `bin`
- **Description** : Chemin vers le fichier binaire généré par `cmd`.
- **Exemple** : `"./tmp/main"`
- **Détails** : Air exécutera ce fichier après la compilation.

---

### `full_bin`
- **Description** : Commande complète pour exécuter le fichier binaire.
- **Exemple** : `"APP_ENV=dev ./tmp/main"`
- **Détails** : Ajoute des variables d'environnement ou des arguments à la commande d'exécution.

---

### `include_ext`
- **Description** : Extensions de fichiers à surveiller pour des changements.
- **Exemple** : `["go"]`
- **Détails** : Spécifiez les extensions de fichiers à surveiller, comme `.go`, `.html`, ou `.css`.

---

### `exclude_dir`
- **Description** : Dossiers à exclure de la surveillance.
- **Exemple** : `["vendor", "templates"]`
- **Détails** : Empêche Air de surveiller les changements dans les dossiers spécifiés.

---

### `exclude_file`
- **Description** : Fichiers spécifiques à exclure de la surveillance.
- **Exemple** : `["config.yaml"]`
- **Détails** : Liste des fichiers qui ne déclencheront pas de rechargement.

---

### `work_dir`
- **Description** : Répertoire de travail pour la construction et l'exécution.
- **Exemple** : `"."`
- **Détails** : Par défaut, le répertoire courant.

---

### `build_delay`
- **Description** : Délai (en millisecondes) avant la reconstruction après un changement.
- **Exemple** : `200`
- **Détails** : Utile pour éviter des reconstructions trop fréquentes lors de changements rapides.

---

### `color_scheme`
- **Description** : Schéma de couleurs pour les messages de la console.
- **Exemple** : `"default"`
- **Options** :
  - `"default"` : Schéma par défaut.
  - `"monochrome"` : Désactive les couleurs.
- **Détails** : Ajuste les couleurs de la sortie en fonction de vos préférences.
```

### PostgreSQL :
Ajoutez le driver PostgreSQL (`github.com/lib/pq`) à votre projet Go en exécutant la commande suivante :

```bash
go get github.com/lib/pq
```

#### Chaîne de connexion :
Voici les paramètres principaux pour configurer la connexion à PostgreSQL :
- **host** : Adresse de votre serveur PostgreSQL.
- **port** : Port utilisé (par défaut `5432`).
- **user** : Nom d'utilisateur.
- **password** : Mot de passe.
- **dbname** : Nom de la base de données.
- **sslmode=disable** : Désactive SSL (utile pour le développement local).

Voici un exemple de chaîne de connexion :
```text
host=localhost port=5432 user=yourusername password=yourpassword dbname=yourdatabase sslmode=disable
```

#### Fonctionnalités principales :
- **`sql.Open`** : Initialise une connexion à la base de données.
- **`db.Ping`** : Vérifie que la connexion est valide.
- **`db.Query`** : Exécute une requête SQL et retourne les résultats.
- **`rows.Next`** et **`rows.Scan`** : Itère sur les résultats et les extrait.
