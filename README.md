# Monday-light

## Cheat sheet

### Go Dependencies

```
go mod init monday-light
go mod tidy
go get github.com/gin-gonic/gin
```

### HTMX

```
[
  { "hx-get": "/get-example", "description": "Sends a GET request to fetch data from the server." },
  { "hx-post": "/post-example", "description": "Sends a POST request with data to the server." },
  { "hx-target": "#result", "description": "Specifies the element to update with the server response." },
  { "hx-trigger": "click", "description": "Triggers the request when the element is clicked." },
  { "hx-trigger": "every 5s", "description": "Triggers the request periodically, every 5 seconds." },
  { "hx-trigger": "load", "description": "Triggers the request when the page loads." },
  { "hx-trigger": "change", "description": "Triggers the request when the value of an input changes." },
  { "hx-swap": "innerHTML", "description": "Replaces the inner HTML of the target element." },
  { "hx-swap": "outerHTML", "description": "Replaces the entire target element with the response." },
  { "hx-swap": "beforeend", "description": "Appends the response content to the end of the target element." },
  { "hx-swap": "afterbegin", "description": "Prepends the response content to the start of the target element." },
  { "hx-swap": "beforebegin", "description": "Inserts the response content before the target element." },
  { "hx-swap": "afterend", "description": "Inserts the response content after the target element." },
  { "hx-on": "htmx:responseError", "description": "Executes JavaScript when a server error occurs." },
  { "hx-vals": "{ 'key': 'value' }", "description": "Sends additional parameters with the request." },
  { "hx-headers": "{ 'Authorization': 'Bearer TOKEN' }", "description": "Adds custom headers to the request." },
  { "hx-indicator": "#loading", "description": "Shows a loading indicator during the request." },
  { "hx-push-url": "true", "description": "Pushes the request URL into the browser history." },
  { "hx-select": ".row", "description": "Selects specific elements from the server response." },
  { "hx-select-oob": ".notification", "description": "Processes out-of-band (OOB) elements from the response." },
  { "hx-confirm": "Are you sure?", "description": "Shows a confirmation dialog before sending the request." },
  { "hx-disable": "true", "description": "Disables the element to prevent duplicate requests." },
  { "hx-history": "false", "description": "Disables browser history for this request." },
  { "hx-include": "#form-id", "description": "Includes additional elements with the request." },
  { "hx-preserve": "true", "description": "Preserves the target element instead of replacing it." },
  { "hx-replace-url": "/new-url", "description": "Replaces the current browser URL with the new URL." }
]
```

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

---

## Exemple de configuration complet

```toml
[build]
cmd = "go build -o ./tmp/main ."
bin = "./tmp/main"
full_bin = "APP_ENV=dev ./tmp/main"
include_ext = ["go"]
exclude_dir = ["vendor", "templates"]
exclude_file = []
work_dir = "."
build_delay = 200
color_scheme = "default"
```
```

### PostgreSQL :
Ajoutez le driver PostgreSQL (`github.com/lib/pq`) à votre projet Go en exécutant la commande suivante :

```bash
go get github.com/lib/pq
```

### Chaîne de connexion :
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

### Fonctionnalités principales :
- **`sql.Open`** : Initialise une connexion à la base de données.
- **`db.Ping`** : Vérifie que la connexion est valide.
- **`db.Query`** : Exécute une requête SQL et retourne les résultats.
- **`rows.Next`** et **`rows.Scan`** : Itère sur les résultats et les extrait.