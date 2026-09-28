# Génération d’images par IA en Python, Node.js, C++, Go, Rust et PHP

> **En bref :** ce dépôt montre comment utiliser une API de génération d’images compatible avec le format OpenAI Images API. Dans chacun des six langages, le programme crée une image à partir d’un texte, enregistre le résultat localement, puis envoie ce fichier à `POST /images/edits` avec une seconde instruction pour le transformer.

Ce guide s’adresse aux développeurs qui recherchent un exemple concret de **text-to-image**, d’**image-to-image** ou d’intégration d’une API de génération d’images. Les exemples utilisent [AI-ROUTER](https://ai-router.dev/fr), un service indépendant proposant une interface compatible avec le format d’API OpenAI. AI-ROUTER n’est pas un produit d’OpenAI et ce dépôt n’est pas un SDK officiel. L’URL de base peut être configurée pour une autre passerelle compatible.

> Les modèles, les opérations disponibles, les tailles d’image, les formats, les quotas et les tarifs dépendent de votre compte et de la configuration actuelle du service amont. La présence d’un nom de modèle dans un exemple ne garantit pas qu’il soit activé pour votre compte. Vérifiez le catalogue et les conditions en vigueur avant un appel payant.

## Ce que ce projet vous aide à comprendre

Une même boucle créative est implémentée dans six environnements. Un premier prompt décrit l’image souhaitée — photo produit, illustration éditoriale ou visuel de concept, par exemple. Le programme transmet le texte à l’API, récupère le résultat et l’écrit sur disque. Il réutilise ensuite exactement ce fichier comme entrée de l’endpoint d’édition, accompagné d’un second prompt qui décrit le changement attendu.

Cette étape intermédiaire rend le workflow vérifiable : l’image modifiée provient d’un fichier réel, pas d’un second prompt indépendant. Vous pouvez examiner le brouillon, le conserver comme variante, le faire valider ou relancer seulement l’édition. C’est un patron utile pour prototyper un générateur de visuels, tester une modification de palette ou d’éclairage, comparer plusieurs SDK/clients HTTP ou préparer une intégration dans un worker, une file de tâches ou un service web.

Une exécution complète effectue au minimum deux opérations qui peuvent être facturées selon votre compte. Le code ne relance pas automatiquement une requête : après un incident réseau, vérifiez si l’amont l’a déjà traitée avant de réessayer.

### Écrire des consignes qui facilitent l’itération

Pour la génération, décrivez d’abord le sujet principal, puis sa composition, l’ambiance ou le style recherché et les détails qui comptent pour votre usage. Une consigne de produit peut préciser l’angle de prise de vue, le fond et l’éclairage ; un visuel éditorial peut préciser la palette, le cadrage et l’atmosphère. Ces éléments donnent un point de départ plus facile à comparer qu’une suite de mots-clés sans relation.

Pour l’édition, indiquez l’élément à changer, la direction du changement et ce qui doit rester stable. Par exemple : « Remplacer l’arrière-plan par un ciel de fin de journée, conserver le sujet principal, sa position et le cadrage. » Les consignes restent des instructions à un modèle génératif : elles ne garantissent pas une édition pixel par pixel ni la conservation parfaite de chaque détail.

## Parcours en deux étapes

~~~text
Prompt texte
    │
    ▼
POST /v1/images/generations  (JSON)
    │ réponse : data[0].b64_json ou data[0].url
    ▼
Image enregistrée localement
    │
    ▼
POST /v1/images/edits        (multipart/form-data, champ image)
    │ réponse : data[0].b64_json ou data[0].url
    ▼
Image modifiée enregistrée localement
~~~

Le résultat généré est sauvegardé avant l’appel d’édition. S’il y a une erreur à la seconde étape, le premier fichier reste disponible pour inspection ou réutilisation.

## Contrat HTTP Images API

L’URL de base par défaut est `https://api.ai-router.dev/v1`. Chaque client y ajoute le chemin de l’opération. Les requêtes transmettent la clé dans l’en-tête `Authorization: Bearer …`.

### Texte vers image : `POST /images/generations`

La génération utilise `application/json`. Exemple de corps :

~~~http
POST https://api.ai-router.dev/v1/images/generations
Authorization: Bearer $AI_ROUTER_API_KEY
Content-Type: application/json
~~~

~~~json
{
  "model": "gpt-image-2",
  "prompt": "Une petite serre en verre sur un toit après la pluie, lumière douce du matin",
  "size": "1024x1024",
  "response_format": "b64_json"
}
~~~

Les paramètres effectivement envoyés diffèrent légèrement selon l’implémentation. Python et Node.js transmettent aussi `quality` ; certains clients envoient `n`. Les clients demandent généralement `b64_json` et savent traiter une réponse `url` lorsque le service en retourne une. Les champs facultatifs et les limites ne sont pas universels : référez-vous à la documentation du modèle et du compte utilisés.

### Image vers image : `POST /images/edits`

L’édition utilise `multipart/form-data`. L’image créée à l’étape précédente est jointe dans le champ fichier `image`.

| Champ | Rôle |
| --- | --- |
| `model` | Identifiant d’un modèle accessible au compte et compatible avec l’opération. |
| `prompt` | Instruction de transformation. Précisez ce qui doit changer et, si besoin, ce qui doit rester intact. |
| `image` | Image source, dans ce workflow le fichier généré puis enregistré localement. |
| `size` | Taille de sortie demandée ; elle doit être prise en charge par le modèle et le compte. |
| `quality` | Transmis par les exemples Python et Node.js ; prise en charge dépendante du service. |
| `response_format` | Les clients privilégient `b64_json` ; une URL peut être retournée selon la configuration. |

La réponse est un objet JSON contenant `data`. Les clients traitent le premier résultat (`data[0]`) en décodant `b64_json` ou en téléchargeant `url` séparément. Une URL peut expirer et doit être considérée comme une entrée externe. Consultez le [contrat API commun](../API_CONTRACT.md) pour les détails de compatibilité.

« Compatible OpenAI » décrit ici une forme de requête et des chemins d’API similaires ; cela ne signifie pas que toutes les passerelles acceptent exactement les mêmes options, ni qu’AI-ROUTER est OpenAI.

## Modèles, tailles et accès au compte

Le modèle par défaut des exemples est `gpt-image-2`. Vous pouvez essayer `gpt-image-2.5`, `gpt-image-2.5-flare` ou `gpt-image-2.5-burst` si le modèle apparaît dans le catalogue de votre compte et prend en charge les deux opérations. Le nom passé à la CLI ne donne pas accès au modèle et ne garantit ni disponibilité, ni qualité, ni tarif.

La taille par défaut est généralement `1024x1024`. D’autres tailles peuvent être passées via `--size` par les clients concernés, mais doivent être acceptées par le modèle sélectionné. Python et Node.js exposent aussi `--quality` : choisissez une valeur reconnue par le service. Commencez par les paramètres documentés pour votre compte et modifiez un réglage à la fois.

Les fichiers `.env.example` sont des modèles ; toutes les CLIs ne chargent pas automatiquement un fichier `.env`. Exportez les variables dans le même terminal que celui qui lance le programme, ou utilisez un chargeur dotenv de confiance. Une option CLI remplace généralement la valeur d’environnement.

Les commandes ci-dessous utilisent la syntaxe Bash (`export`). Dans Windows PowerShell, définissez la clé avec `$env:AI_ROUTER_API_KEY = "votre-cle"`, puis lancez la commande du langage depuis ce même terminal.

| Langage | Variables pour la clé API | URL, modèle, taille et options |
| --- | --- | --- |
| Python | `IMAGE_API_KEY`, puis `OPENAI_API_KEY`, puis `AI_ROUTER_API_KEY` ; option `--api-key` | `IMAGE_API_BASE_URL` / `AI_ROUTER_BASE_URL` ; `IMAGE_MODEL` / `AI_ROUTER_MODEL` ; `IMAGE_SIZE` / `AI_ROUTER_SIZE` ; `IMAGE_QUALITY` ; `IMAGE_TIMEOUT_SECONDS` (secondes). |
| Node.js | `IMAGE_API_KEY`, `OPENAI_API_KEY` ou `AI_ROUTER_API_KEY` ; option `--api-key` | `IMAGE_API_BASE_URL` / `AI_ROUTER_BASE_URL` ; variables `IMAGE_MODEL`, `IMAGE_SIZE`, `IMAGE_QUALITY` ; `IMAGE_TIMEOUT_MS` (millisecondes). |
| C++ | `OPENAI_API_KEY` ou `AI_ROUTER_API_KEY` ; option `--api-key` | `OPENAI_BASE_URL` / `AI_ROUTER_BASE_URL` ; `IMAGE_MODEL` / `AI_ROUTER_MODEL` ; `IMAGE_SIZE` / `AI_ROUTER_SIZE` ; `IMAGE_RESPONSE_FORMAT` ; `IMAGE_API_TIMEOUT` ou `AI_ROUTER_IMAGE_TIMEOUT` (secondes). |
| Go | `OPENAI_API_KEY` ou `AI_ROUTER_API_KEY` ; option `--api-key` | `OPENAI_BASE_URL` / `AI_ROUTER_BASE_URL` ; `IMAGE_MODEL` / `AI_ROUTER_MODEL` ; `IMAGE_SIZE` / `AI_ROUTER_SIZE` ; `IMAGE_RESPONSE_FORMAT` ; `IMAGE_API_TIMEOUT` (secondes). |
| Rust | `AI_ROUTER_API_KEY` ou `OPENAI_API_KEY` | `AI_ROUTER_BASE_URL`, `AI_ROUTER_TIMEOUT_SECONDS` ; `AI_ROUTER_MODEL` et `AI_ROUTER_SIZE` fournissent les valeurs par défaut des options CLI. |
| PHP | `AI_ROUTER_API_KEY` ou `OPENAI_API_KEY` | `AI_ROUTER_BASE_URL`, `AI_ROUTER_TIMEOUT_SECONDS`, `AI_ROUTER_MODEL`, `AI_ROUTER_SIZE` ; la clé ne se transmet pas par option CLI. |

La racine d’API contient généralement `/v1`. Python et Node.js normalisent automatiquement l’URL de base ; pour les autres clients, indiquez l’URL conforme à leur valeur par défaut, par exemple `https://api.ai-router.dev/v1`.

## Prérequis et commandes exécutables

Les chemins partent de la racine du dépôt. Chaque exemple ci-dessous exécute génération → sauvegarde → édition → sauvegarde. Remplacez `votre-cle` dans le terminal ; ne mettez jamais la clé dans le code suivi par Git.

### Python — 3.10 ou plus récent

Python utilise la bibliothèque standard, sans dépendance d’exécution supplémentaire.

~~~bash
cd marketing/image-generation-multilang/python
export AI_ROUTER_API_KEY="votre-cle"
python main.py \
  --prompt "Une photo éditoriale d'une serre en verre sous la pluie" \
  --edit-prompt "Réchauffer la lumière comme au coucher du soleil sans modifier la forme de la serre" \
  --model gpt-image-2 \
  --output-dir ./outputs
python -m unittest discover -s tests -v
~~~

Les images sont écrites dans `outputs/`. La CLI accepte notamment `--base-url`, `--size`, `--quality` et `--timeout` (secondes). Les options écrasent l’environnement.

### Node.js — 18.17 ou plus récent

Le transport utilise `fetch`, `FormData` et `Blob` intégrés à Node.js ; pas de dépendance npm d’exécution.

~~~bash
cd marketing/image-generation-multilang/nodejs
npm test
export AI_ROUTER_API_KEY="votre-cle"
npm run demo -- \
  --prompt "Une illustration produit d'un sac rouge sur une table claire" \
  --edit-prompt "Remplacer le rouge par du bleu cobalt et conserver le cadrage" \
  --model gpt-image-2 \
  --output-dir ./outputs
~~~

Les fichiers par défaut sont `outputs/generated.png` et `outputs/edited.png`. Les options incluent `--base-url`, `--size`, `--quality` et `--timeout`, exprimé en millisecondes.

### C++ — C++17, CMake 3.20+ et libcurl

Installez un compilateur C++17, CMake et les fichiers de développement libcurl. Avec vcpkg sous Windows, le paquet d’exemple est `vcpkg install curl:x64-windows`.

~~~bash
cd marketing/image-generation-multilang/cpp
cmake -S . -B build -DBUILD_TESTING=ON
cmake --build build --config Release
ctest --test-dir build --output-on-failure
export AI_ROUTER_API_KEY="votre-cle"
./build/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "Une photo produit cinématographique d'une lampe en verre" \
  --edit-prompt "Ajouter une lumière chaude et conserver les proportions de la lampe" \
  --output generated.png --edited-output edited.png
~~~

Avec Visual Studio, le binaire multi-configuration est généralement sous `build/Release/imagegen.exe`. Si CMake ne trouve pas libcurl de vcpkg, ajoutez l’option `-DCMAKE_TOOLCHAIN_FILE=.../vcpkg.cmake` à la configuration. Les tests CTest sont hors ligne.

### Go — 1.22 ou plus récent

Cette version s’appuie sur la bibliothèque standard et ne demande pas de paquet externe pour le client.

~~~bash
cd marketing/image-generation-multilang/go
go test ./...
go vet ./...
export AI_ROUTER_API_KEY="votre-cle"
go run ./cmd/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "Une serre contemporaine dans un jardin urbain" \
  --edit-prompt "Ajouter des plantes grimpantes sur les côtés sans changer la façade" \
  --output generated.png --edited-output edited.png
~~~

Les fichiers sont enregistrés dans le répertoire courant. Les tests utilisent un serveur HTTP simulé et ne nécessitent pas de clé ; `go build ./cmd/imagegen` construit le binaire.

### Rust — stable, Rust 1.80+

Cargo télécharge les crates définies dans `Cargo.toml` au premier build. Le sous-ordre `workflow` réalise les deux opérations ; `generate` et `edit` peuvent aussi être utilisés séparément.

~~~bash
cd marketing/image-generation-multilang/rust
export AI_ROUTER_API_KEY="votre-cle"
cargo test
cargo run --release -- workflow \
  --prompt "Un sac à dos rouge photographié en studio" \
  --edit-prompt "Le passer au bleu cobalt en gardant le cadrage" \
  --model gpt-image-2 --size 1024x1024 \
  --generated output/generated.png \
  --edited output/edited.png
~~~

Le programme crée les répertoires de sortie manquants. Pour modifier l’adresse, utilisez `AI_ROUTER_BASE_URL` ou l’option globale `--base-url` avant le sous-ordre. `AI_ROUTER_MODEL`, `AI_ROUTER_SIZE` et `AI_ROUTER_TIMEOUT_SECONDS` fournissent des valeurs d’environnement.

### PHP — 8.1+, extensions cURL et JSON

Le démonstrateur contient son propre autoloader ; Composer n’est pas requis pour l’exécuter. Vérifiez que les extensions sont activées dans le PHP CLI.

~~~bash
cd marketing/image-generation-multilang/php
php -m
export AI_ROUTER_API_KEY="votre-cle"
php bin/image-demo.php workflow \
  --prompt "Une serre transparente sous une pluie de printemps" \
  --edit-prompt "Ajouter des lampes suspendues et conserver le cadrage" \
  --model gpt-image-2 --size 1024x1024 \
  --generated output/php-generated.png \
  --edited output/php-edited.png
~~~

Les fichiers sont créés sous `output/`. Les commandes `generate` et `edit` fonctionnent aussi séparément. La clé est lue depuis l’environnement et n’est pas une option CLI. Il n’existe pas encore de suite de tests PHP dédiée ; vérifiez au minimum la syntaxe avec `php -l bin/image-demo.php`.

## Architecture : responsabilités séparées

Les six implémentations séparent les responsabilités pour que le transport HTTP reste réutilisable hors de la CLI :

1. **Configuration** lit et valide la clé, l’URL, le modèle, la taille, le timeout et les options d’exécution.
2. **Client API** gère l’authentification Bearer, le JSON et le multipart, les erreurs HTTP, le décodage de la réponse et le téléchargement éventuel.
3. **Workflow** orchestre génération → sauvegarde → édition → sauvegarde.
4. **CLI** parse les arguments et affiche les chemins produits ; elle ne duplique pas les détails de transport.

| Langage | Modules principaux | Choix d’implémentation |
| --- | --- | --- |
| Python | `python/image_demo/config.py`, `client.py`, `workflow.py`, `cli.py` | `urllib` standard, tests de configuration sans appel réel. |
| Node.js | `nodejs/src/config.mjs`, `client.mjs`, `workflow.mjs`, `cli.mjs` | API Web natives et tests `node:test`. |
| C++ | `cpp/include/image_client.hpp`, `cpp/src/image_client.cpp`, `cpp/src/main.cpp` | Client libcurl réutilisable, CMake et CTest hors ligne. |
| Go | `go/internal/config`, `go/internal/client`, `go/cmd/imagegen` | Client indépendant de la CLI, tests de transport avec `httptest.Server`. |
| Rust | `rust/src/config.rs`, `client.rs`, `workflow.rs`, `artifacts.rs`, `main.rs` | Types dédiés, commandes `clap`, client asynchrone `reqwest`. |
| PHP | `php/src/Config.php`, `ImageClient.php`, `Workflow.php`, `ArtifactStore.php`, `php/bin/image-demo.php` | cURL, données d’image et stockage séparés de la CLI. |

Pour l’intégration dans un serveur ou un worker, réutilisez le client et le workflow mais remplacez le stockage local par un stockage objet adapté. Les détails figurent dans le [guide d’architecture](../ARCHITECTURE.md) et les règles d’ajout de langage dans [CONTRIBUTING.md](../../CONTRIBUTING.md).

## Tests locaux sans appel payant

Exécutez chaque commande depuis le dossier correspondant. Ces vérifications ne nécessitent pas de clé API :

| Langage | Vérifications |
| --- | --- |
| Python | `python -m unittest discover -s tests -v` ; `python -m compileall -q image_demo main.py` |
| Node.js | `npm test` |
| C++ | `cmake -S . -B build -DBUILD_TESTING=ON`, build, puis `ctest --test-dir build --output-on-failure` |
| Go | `go test ./...`, `go vet ./...`, `go build ./cmd/imagegen` |
| Rust | `cargo test` ou `cargo check` |
| PHP | Pas de suite dédiée ; `php -l bin/image-demo.php` vérifie la syntaxe du point d’entrée. |

Une compilation réussie ne valide ni les permissions du compte, ni l’accès au modèle, ni le tarif d’une requête réelle.

## Dépannage

| Symptôme | Pistes de vérification |
| --- | --- |
| `401`, `403` ou clé manquante | Définissez la variable dans le même shell que le processus. Vérifiez la validité de la clé et les droits du compte. Les alias varient selon le langage. |
| Modèle inconnu ou interdit | Vérifiez l’identifiant dans le catalogue actuel et l’accès aux deux opérations, génération et édition. |
| Taille ou qualité refusée | Utilisez uniquement les valeurs prises en charge par le modèle. `1024x1024` et `auto` sont des exemples, pas des garanties globales. |
| Timeout ou problème TLS/proxy | Vérifiez URL, certificat, proxy réseau et unité du délai : Node.js utilise les millisecondes, les autres exemples les secondes. Ne relancez pas aveuglément un appel qui peut être facturé. |
| Pas d’image ou base64 invalide | Contrôlez la réponse `data[0].b64_json` ou `data[0].url` et la compatibilité du gateway. Supprimez les données privées avant de partager les diagnostics. |
| Génération réussie, édition échouée | Vérifiez que le fichier source est lisible, le prompt non vide et le modèle compatible avec `images/edits`. Conservez le brouillon : le workflow ne le supprime pas. |
| Téléchargement d’URL échoué | L’URL peut expirer ou être inaccessible. En service public, ne téléchargez pas une URL arbitraire sans contrôle. |
| Échec C++ ou PHP | C++ doit trouver les en-têtes et bibliothèques libcurl via CMake. Le binaire PHP CLI doit charger les extensions cURL et JSON. |

Pour isoler un problème, commencez par un prompt court, un modèle confirmé et une taille documentée. Notez le statut HTTP et l’identifiant de requête si disponible, sans publier de clé, prompt privé ou image client.

## Sécurité, confidentialité et mise en production

Ces exemples ne sont pas un service d’upload public prêt à déployer. Avant la production :

- conservez la clé côté serveur, dans des variables protégées ou un gestionnaire de secrets ; ne l’intégrez ni au navigateur, ni au dépôt, ni à une capture ;
- utilisez HTTPS et vérifiez que l’URL de base désigne le service attendu ;
- validez la longueur des prompts, le vrai type MIME, les dimensions et le poids des fichiers ; ajoutez quotas, limitation de débit et contrôles de contenu ;
- remplacez les fichiers locaux par un stockage durable avec contrôle d’accès, chiffrement, durée de conservation et suppression adaptés ;
- traitez prompts, images et URL temporaires comme des données potentiellement sensibles ; ne les journalisez pas sans nécessité et n’enregistrez jamais la clé ;
- lors d’un téléchargement par URL, appliquez une liste de domaines autorisés et des limites de redirection, protocole, durée et taille afin de réduire les risques SSRF et d’épuisement de ressources ;
- ajoutez annulation, métriques et reprises limitées seulement après avoir compris l’idempotence et la facturation : une seconde génération peut être une seconde opération facturée ;
- vérifiez les règles actuelles de traitement des données du fournisseur et ne transmettez que les images que vous avez le droit de traiter.

La forme de l’API ne permet pas de déduire la conservation des données. Consultez le service utilisé et son compte. Voir aussi [SECURITY.md](../SECURITY.md).

## Documentation et ressources

- [README principal en anglais](../../README.md) · [contrat API](../API_CONTRACT.md) · [architecture](../ARCHITECTURE.md) · [sécurité](../SECURITY.md)
- [Prompts d’exemple](../../examples/prompts.md) pour comparer les six implémentations avec des entrées similaires.
- Autres guides : [chinois](README.zh-CN.md), [russe](README.ru.md), [japonais](README.ja.md).

Pour découvrir la plateforme, l’accès API et les modèles visibles depuis votre compte, consultez le [site AI-ROUTER en français](https://ai-router.dev/fr). La disponibilité dépend de votre compte ; ce guide ne garantit ni accès à un modèle, ni résultat particulier, ni classement dans les moteurs de recherche.
