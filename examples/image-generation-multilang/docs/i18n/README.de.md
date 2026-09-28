# KI-Bildgenerierung mit Python, Node.js, C++, Go, Rust und PHP

> **Kurz erklärt:** Dieses Repository zeigt den Anschluss an eine mit dem OpenAI-Images-API-Format kompatible Bildgenerierungs-API. Jedes der sechs Sprachbeispiele erzeugt zuerst ein Bild aus Text, speichert es lokal und sendet genau diese Datei anschließend mit einer zweiten Änderungsanweisung an `POST /images/edits`.

Der Leitfaden richtet sich an Entwicklerinnen und Entwickler, die **Text-zu-Bild-Generierung**, **Bildbearbeitung per API** oder die Integration einer Bildgenerierungs-API praktisch nachvollziehen möchten. Die Beispiele verwenden [AI-ROUTER](https://ai-router.dev/de), einen unabhängigen Dienst mit einer OpenAI-kompatiblen API-Schnittstelle. AI-ROUTER ist kein Produkt von OpenAI und dieses Repository ist kein offizielles OpenAI-SDK. Die API-Basisadresse lässt sich für kompatible Gateways anpassen.

> Modellnamen, verfügbare Operationen, Formate, Bildgrößen, Kontingente und Preise richten sich nach Ihrem Konto und der aktuellen Konfiguration des Anbieters. Ein hier aufgeführtes Modell ist nicht automatisch für jedes Konto freigeschaltet. Prüfen Sie den Modellkatalog und die geltenden Konditionen, bevor Sie kostenpflichtige Anfragen senden.

## Was dieses Beispielprojekt vermittelt

Ein kreativer Mini-Workflow ist für sechs Sprachen und Laufzeitumgebungen umgesetzt. Ein erster Prompt beschreibt das gewünschte Bild, zum Beispiel ein Produktfoto, eine redaktionelle Illustration oder einen Konzeptentwurf. Das Programm übermittelt den Text an die API, liest das erste Bild aus der Antwort und speichert es lokal. Danach wird genau diese Datei als Eingabe für den Edit-Endpunkt verwendet; ein zweiter Prompt beschreibt die gewünschte Änderung.

Die Zwischenspeicherung ist wichtig: Das bearbeitete Bild entsteht aus einer echten Bilddatei und nicht aus einem unabhängigen zweiten Textprompt. Der Entwurf kann geprüft, als Variante behalten, manuell freigegeben oder nach einem Fehler des Bearbeitungsschritts erneut verwendet werden. Das Muster eignet sich, um einen Generator für Produktbilder oder Kampagnenmotive zu prototypisieren, Farbgebung und Licht zu variieren, die HTTP-Clients von Python, Node.js, C++, Go, Rust und PHP zu vergleichen oder einen Hintergrundjob und Webdienst vorzubereiten.

Ein vollständiger Durchlauf löst mindestens eine Generierungs- und eine Bearbeitungsoperation aus. Je nach Konto können beide kostenpflichtig sein. Automatische Wiederholungen sind bewusst nicht eingebaut: Prüfen Sie nach einem Netzwerkfehler zuerst, ob der Anbieter die Anfrage bereits verarbeitet hat.

### Prompts so formulieren, dass Iterationen nachvollziehbar bleiben

Beschreiben Sie bei der ersten Generierung zunächst das Hauptmotiv und ergänzen Sie dann Bildaufbau, gewünschte Stimmung oder Stil sowie Details, die für den geplanten Einsatz wichtig sind. Bei einem Produktmotiv können das Blickwinkel, Hintergrund und Licht sein; bei einer redaktionellen Illustration etwa Farbpalette, Ausschnitt und Atmosphäre. So entsteht ein besser überprüfbarer Ausgangspunkt als durch eine unverbundene Liste von Schlagwörtern.

Beim Edit-Prompt sollten Sie benennen, was geändert werden soll, in welche Richtung die Änderung geht und welche Eigenschaften erhalten bleiben sollen. Beispiel: „Den Hintergrund durch einen Abendhimmel ersetzen, das Hauptmotiv, seine Position und den Bildausschnitt beibehalten.“ Ein Prompt bleibt eine Anweisung an ein generatives Modell und garantiert keine pixelgenaue Bearbeitung oder perfekte Erhaltung jedes Details.

## Ablauf in zwei Schritten

~~~text
Textprompt
    │
    ▼
POST /v1/images/generations  (JSON)
    │ Antwort: data[0].b64_json oder data[0].url
    ▼
Bilddatei lokal speichern
    │
    ▼
POST /v1/images/edits        (multipart/form-data, Feld image)
    │ Antwort: data[0].b64_json oder data[0].url
    ▼
Bearbeitetes Bild lokal speichern
~~~

Das generierte Original wird vor der Edit-Anfrage gespeichert. Scheitert der zweite Schritt, bleibt der Entwurf zur Prüfung oder Wiederverwendung erhalten.

## Gemeinsamer Images-API-Vertrag

Die Standardbasisadresse lautet `https://api.ai-router.dev/v1`. Der Client ergänzt den Pfad des jeweiligen Endpunkts. Bei beiden Operationen wird der API-Schlüssel als Bearer-Token im `Authorization`-Header übertragen.

### Text zu Bild: `POST /images/generations`

Die Generierungsanfrage verwendet `application/json`. Beispiel für den Request-Body:

~~~http
POST https://api.ai-router.dev/v1/images/generations
Authorization: Bearer $AI_ROUTER_API_KEY
Content-Type: application/json
~~~

~~~json
{
  "model": "gpt-image-2",
  "prompt": "Ein kleines Gewächshaus aus Glas auf einem Dach nach dem Regen, weiches Morgenlicht",
  "size": "1024x1024",
  "response_format": "b64_json"
}
~~~

Je nach Implementierung kommen weitere Parameter hinzu. Python und Node.js senden außerdem `quality`; einzelne Clients übermitteln `n`. Die Clients fragen in der Regel `b64_json` an und können auch ein Bild unter `url` verarbeiten, wenn der Dienst eine URL zurückgibt. Optionale Felder und Limits sind nicht universell. Maßgeblich bleibt die aktuelle Dokumentation des gewählten Modells und Kontos.

### Bild zu Bild: `POST /images/edits`

Die Bildbearbeitung verwendet `multipart/form-data`. Die im ersten Schritt gespeicherte Datei wird im Feld `image` hochgeladen.

| Multipart-Feld | Bedeutung |
| --- | --- |
| `model` | Kennung eines Modells, das für das Konto freigeschaltet ist und die gewünschte Operation unterstützt. |
| `prompt` | Änderungsanweisung. Beschreiben Sie, was angepasst und gegebenenfalls beibehalten werden soll. |
| `image` | Quelldatei; im Workflow das lokal gespeicherte Ergebnis der Generierung. |
| `size` | Gewünschte Ausgabegröße, sofern das Modell und Konto sie unterstützen. |
| `quality` | Wird in den Python- und Node.js-Beispielen mitgesendet; die Unterstützung hängt vom Dienst ab. |
| `response_format` | Die Beispiele bevorzugen `b64_json`; abhängig von der Dienstkonfiguration kann eine URL zurückkommen. |

Die API-Antwort enthält ein JSON-Array `data`. Die Clients verarbeiten dessen erstes Element (`data[0]`), decodieren `b64_json` oder laden `url` separat herunter. Bild-URLs können ablaufen und sind als externe Eingabe zu behandeln. Der vollständige gemeinsame Vertrag steht in [API_CONTRACT.md](../API_CONTRACT.md).

„OpenAI-kompatibel“ bezeichnet hier ein ähnliches Anfrageformat und ähnliche API-Pfade. Es garantiert nicht, dass alle Gateways dieselben Parameter oder Metadaten akzeptieren; AI-ROUTER ist kein OpenAI-Angebot.

## Modellwahl, Bildgröße und Kontozugriff

Der Standardmodellname der Beispiele ist `gpt-image-2`. Sie können auch `gpt-image-2.5`, `gpt-image-2.5-flare` oder `gpt-image-2.5-burst` ausprobieren, **sofern** das Modell im Katalog Ihres Kontos auftaucht und beide benötigten Operationen unterstützt. Der Modellname im CLI-Aufruf schaltet keinen Zugriff frei und garantiert weder Verfügbarkeit noch Ergebnisqualität oder Preis.

Als Standardgröße wird meist `1024x1024` verwendet. Andere Größen können über `--size` angegeben werden, müssen aber vom gewählten Modell unterstützt werden. Python und Node.js bieten zusätzlich `--quality`. Wählen Sie nur einen Wert, den der Dienst akzeptiert, und testen Sie möglichst jeweils nur eine Änderung.

Die Dateien `.env.example` sind Vorlagen; nicht jedes CLI lädt eine Datei `.env` automatisch. Setzen Sie die Variablen im selben Terminal, aus dem das Programm gestartet wird, oder nutzen Sie einen vertrauenswürdigen dotenv-Loader. Kommandozeilenoptionen überschreiben normalerweise die Umgebungswerte.

Die folgenden Befehle verwenden Bash-Syntax (`export`). In Windows PowerShell setzen Sie den Schlüssel stattdessen mit `$env:AI_ROUTER_API_KEY = "Ihr-api-key"` und starten danach den jeweiligen Befehl im selben Terminal.

| Sprache | API-Key-Variablen | URL, Modell, Größe und weitere Optionen |
| --- | --- | --- |
| Python | `IMAGE_API_KEY`, dann `OPENAI_API_KEY`, dann `AI_ROUTER_API_KEY`; außerdem `--api-key` | `IMAGE_API_BASE_URL` / `AI_ROUTER_BASE_URL`; `IMAGE_MODEL` / `AI_ROUTER_MODEL`; `IMAGE_SIZE` / `AI_ROUTER_SIZE`; `IMAGE_QUALITY`; `IMAGE_TIMEOUT_SECONDS` (Sekunden). |
| Node.js | `IMAGE_API_KEY`, `OPENAI_API_KEY` oder `AI_ROUTER_API_KEY`; außerdem `--api-key` | `IMAGE_API_BASE_URL` / `AI_ROUTER_BASE_URL`; `IMAGE_MODEL`, `IMAGE_SIZE`, `IMAGE_QUALITY`; `IMAGE_TIMEOUT_MS` (Millisekunden). |
| C++ | `OPENAI_API_KEY` oder `AI_ROUTER_API_KEY`; außerdem `--api-key` | `OPENAI_BASE_URL` / `AI_ROUTER_BASE_URL`; `IMAGE_MODEL` / `AI_ROUTER_MODEL`; `IMAGE_SIZE` / `AI_ROUTER_SIZE`; `IMAGE_RESPONSE_FORMAT`; `IMAGE_API_TIMEOUT` oder `AI_ROUTER_IMAGE_TIMEOUT` (Sekunden). |
| Go | `OPENAI_API_KEY` oder `AI_ROUTER_API_KEY`; außerdem `--api-key` | `OPENAI_BASE_URL` / `AI_ROUTER_BASE_URL`; `IMAGE_MODEL` / `AI_ROUTER_MODEL`; `IMAGE_SIZE` / `AI_ROUTER_SIZE`; `IMAGE_RESPONSE_FORMAT`; `IMAGE_API_TIMEOUT` (Sekunden). |
| Rust | `AI_ROUTER_API_KEY` oder `OPENAI_API_KEY` | `AI_ROUTER_BASE_URL`, `AI_ROUTER_TIMEOUT_SECONDS`; `AI_ROUTER_MODEL` und `AI_ROUTER_SIZE` setzen Standardwerte der CLI-Optionen. |
| PHP | `AI_ROUTER_API_KEY` oder `OPENAI_API_KEY` | `AI_ROUTER_BASE_URL`, `AI_ROUTER_TIMEOUT_SECONDS`, `AI_ROUTER_MODEL`, `AI_ROUTER_SIZE`; der Schlüssel wird nicht als CLI-Option übergeben. |

Die API-Basisadresse enthält üblicherweise `/v1`. Python und Node.js normalisieren die Basis-URL automatisch. Bei den übrigen Clients sollten Sie die Adresse einschließlich `/v1` angeben, beispielsweise `https://api.ai-router.dev/v1`.

## Voraussetzungen und ausführbare Beispiele

Die folgenden Pfade beginnen im Repository-Stamm. Jeder Befehl führt Generierung → Speichern → Bearbeiten → Speichern aus. Ersetzen Sie `Ihr-api-key` im Terminal, nicht in einer versionierten Datei.

### Python — Version 3.10 oder neuer

Der Python-Client verwendet ausschließlich die Standardbibliothek und benötigt keine zusätzlichen Laufzeitpakete.

~~~bash
cd marketing/image-generation-multilang/python
export AI_ROUTER_API_KEY="Ihr-api-key"
python main.py \
  --prompt "Eine redaktionelle Aufnahme eines Gewächshauses aus Glas im Regen" \
  --edit-prompt "Das Licht wie bei Sonnenuntergang wärmer gestalten und die Form beibehalten" \
  --model gpt-image-2 \
  --output-dir ./outputs
python -m unittest discover -s tests -v
~~~

Die Bilddateien liegen unter `outputs/`. Weitere Optionen sind `--base-url`, `--size`, `--quality` und `--timeout` (Sekunden). CLI-Optionen überschreiben Umgebungswerte.

### Node.js — Version 18.17 oder neuer

Der Client nutzt `fetch`, `FormData` und `Blob` der Laufzeitumgebung. Ein zusätzliches npm-Laufzeitpaket ist nicht nötig.

~~~bash
cd marketing/image-generation-multilang/nodejs
npm test
export AI_ROUTER_API_KEY="Ihr-api-key"
npm run demo -- \
  --prompt "Eine Produktillustration: roter Rucksack auf einem hellen Tisch" \
  --edit-prompt "Den Rucksack kobaltblau färben und den Bildausschnitt beibehalten" \
  --model gpt-image-2 \
  --output-dir ./outputs
~~~

Standardausgabe: `outputs/generated.png` und `outputs/edited.png`. Das Demo kennt Optionen wie `--base-url`, `--size`, `--quality` und `--timeout` (Millisekunden).

### C++ — C++17, CMake 3.20+ und libcurl

Erforderlich sind ein C++17-Compiler, CMake und libcurl-Entwicklungsdateien. Unter Windows mit vcpkg lässt sich die Beispielabhängigkeit mit `vcpkg install curl:x64-windows` installieren.

~~~bash
cd marketing/image-generation-multilang/cpp
cmake -S . -B build -DBUILD_TESTING=ON
cmake --build build --config Release
ctest --test-dir build --output-on-failure
export AI_ROUTER_API_KEY="Ihr-api-key"
./build/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "Ein filmisches Produktfoto einer Lampe aus Glas" \
  --edit-prompt "Warmes Licht ergänzen und die Proportionen der Lampe beibehalten" \
  --output generated.png --edited-output edited.png
~~~

Mit Visual Studio liegt das Programm eines Multi-Config-Builds häufig unter `build/Release/imagegen.exe`. Wenn CMake die vcpkg-Version von libcurl nicht findet, ergänzen Sie bei der Konfiguration `-DCMAKE_TOOLCHAIN_FILE=.../vcpkg.cmake`. Die CTest-Prüfungen sind offline.

### Go — Version 1.22 oder neuer

Die Go-Implementierung verwendet Standardbibliotheken und benötigt keine zusätzlichen Laufzeitmodule.

~~~bash
cd marketing/image-generation-multilang/go
go test ./...
go vet ./...
export AI_ROUTER_API_KEY="Ihr-api-key"
go run ./cmd/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "Ein modernes kleines Gewächshaus in einem Stadtgarten" \
  --edit-prompt "Rankpflanzen an den Seiten ergänzen und die Fassade beibehalten" \
  --output generated.png --edited-output edited.png
~~~

Die Dateien werden standardmäßig im aktuellen Arbeitsverzeichnis angelegt. Tests nutzen einen simulierten HTTP-Server und benötigen weder Schlüssel noch echte Generierung. Mit `go build ./cmd/imagegen` lässt sich die ausführbare Datei separat bauen.

### Rust — stabile Toolchain, Rust 1.80+

Cargo lädt beim ersten Build die in `Cargo.toml` definierten Crates. Der Client nutzt `reqwest`, `serde`, `clap` und Tokio. Der Unterbefehl `workflow` führt beide Schritte aus; `generate` und `edit` können separat aufgerufen werden.

~~~bash
cd marketing/image-generation-multilang/rust
export AI_ROUTER_API_KEY="Ihr-api-key"
cargo test
cargo run --release -- workflow \
  --prompt "Ein roter Rucksack als Studioaufnahme" \
  --edit-prompt "Den Rucksack kobaltblau färben und den Bildausschnitt beibehalten" \
  --model gpt-image-2 --size 1024x1024 \
  --generated output/generated.png \
  --edited output/edited.png
~~~

Fehlende Ausgabeverzeichnisse werden erstellt. Für eine andere API-Adresse setzen Sie `AI_ROUTER_BASE_URL` oder übergeben die globale Option `--base-url` vor dem Unterbefehl. `AI_ROUTER_MODEL`, `AI_ROUTER_SIZE` und `AI_ROUTER_TIMEOUT_SECONDS` liefern Umgebungs-Standardwerte.

### PHP — Version 8.1+, cURL- und JSON-Erweiterungen

Das Demo enthält einen kleinen eigenen Autoloader und benötigt zur Laufzeit kein Composer-Paket. Prüfen Sie, ob die PHP-CLI-Umgebung die cURL- und JSON-Erweiterung geladen hat.

~~~bash
cd marketing/image-generation-multilang/php
php -m
export AI_ROUTER_API_KEY="Ihr-api-key"
php bin/image-demo.php workflow \
  --prompt "Ein transparentes Gewächshaus im Frühlingsregen" \
  --edit-prompt "Einige Hängelampen ergänzen und den Bildausschnitt erhalten" \
  --model gpt-image-2 --size 1024x1024 \
  --generated output/php-generated.png \
  --edited output/php-edited.png
~~~

Beide Bilder werden unter `output/` gespeichert. `generate` und `edit` sind auch einzeln verfügbar. Setzen Sie den Schlüssel in der Umgebung; das PHP-CLI akzeptiert keine Option `--api-key`. Eine eigene PHP-Testsuite gibt es derzeit nicht; mindestens die CLI-Syntax lässt sich mit `php -l bin/image-demo.php` prüfen.

## Architektur: klare Zuständigkeiten

Die Implementierungen trennen Konfiguration, HTTP-Transport, Ablaufsteuerung und Kommandozeile. So kann der Client in einem anderen Anwendungsteil wiederverwendet werden, ohne CLI-Argumente oder Dateipfade in die API-Schicht zu kopieren:

1. **Konfiguration** liest und prüft API-URL, Schlüssel, Modell, Bildgröße, Timeout und Laufzeitoptionen.
2. **API-Client** übernimmt Bearer-Authentifizierung, JSON- und Multipart-Kodierung, HTTP-Fehler, Antwortdecodierung und optionalen Bilddownload.
3. **Workflow** orchestriert Generierung → Speichern → Bearbeiten → Speichern.
4. **CLI** parst Eingaben und zeigt die Ausgabepfade an; die Netzwerklösung bleibt außerhalb des CLI-Einstiegs.

| Sprache | Zentrale Module | Architekturmerkmal |
| --- | --- | --- |
| Python | `python/image_demo/config.py`, `client.py`, `workflow.py`, `cli.py` | `urllib` aus der Standardbibliothek; Konfigurationstests ohne echten API-Aufruf. |
| Node.js | `nodejs/src/config.mjs`, `client.mjs`, `workflow.mjs`, `cli.mjs` | Native Web-APIs, Tests mit `node:test`. |
| C++ | `cpp/include/image_client.hpp`, `cpp/src/image_client.cpp`, `cpp/src/main.cpp` | Wiederverwendbarer libcurl-Client, CMake und Offline-CTest. |
| Go | `go/internal/config`, `go/internal/client`, `go/cmd/imagegen` | Vom CLI getrennter Client, HTTP-Tests mit `httptest.Server`. |
| Rust | `rust/src/config.rs`, `client.rs`, `workflow.rs`, `artifacts.rs`, `main.rs` | Eigene Typen, `clap`-Befehle, asynchrones `reqwest`. |
| PHP | `php/src/Config.php`, `ImageClient.php`, `Workflow.php`, `ArtifactStore.php`, `php/bin/image-demo.php` | cURL-Transport, Bilddaten und Speicherung getrennt vom CLI. |

Für einen Webdienst oder Worker können Client und Workflow übernommen und die lokalen Dateien durch passenden Objektspeicher ersetzt werden. Weitere Informationen stehen in der [Architekturbeschreibung](../ARCHITECTURE.md) und in den Regeln zum Ergänzen einer Sprache unter [CONTRIBUTING.md](../../CONTRIBUTING.md).

## Lokal testen, ohne die API aufzurufen

Die folgenden Prüfungen benötigen keinen API-Schlüssel und starten keine kostenpflichtige Bildgenerierung. Führen Sie sie jeweils im Sprachverzeichnis aus.

| Sprache | Prüfungen |
| --- | --- |
| Python | `python -m unittest discover -s tests -v` und `python -m compileall -q image_demo main.py` |
| Node.js | `npm test` |
| C++ | `cmake -S . -B build -DBUILD_TESTING=ON`, Build und `ctest --test-dir build --output-on-failure` |
| Go | `go test ./...`, `go vet ./...`, `go build ./cmd/imagegen` |
| Rust | `cargo test` oder `cargo check` |
| PHP | Noch keine eigene Testsuite; `php -l bin/image-demo.php` prüft die Syntax des CLI-Einstiegs. |

Ein erfolgreicher Build bestätigt weder die Kontoberechtigung noch den Modellzugriff, die Kompatibilität einer anderen API oder die Kosten einer echten Anfrage.

## Fehlerbehebung bei der Bildgenerierung

| Symptom | Prüfen |
| --- | --- |
| `401`, `403` oder fehlender API-Key | Setzen Sie die Variable in derselben Shell wie den Prozess. Prüfen Sie Gültigkeit und Kontoberechtigung. Die Variablennamen unterscheiden sich zwischen den Beispielen. |
| Modell nicht gefunden oder kein Zugriff | Prüfen Sie Modellkennung im aktuellen Kontokatalog und Zugriff auf Generierung **und** Bearbeitung. |
| Größe oder Qualität abgelehnt | Verwenden Sie nur Werte, die das gewählte Modell unterstützt. `1024x1024` und `auto` sind Beispiele, keine universellen Zusagen. |
| Timeout, TLS- oder Proxyfehler | Prüfen Sie URL, Zertifikat, Netzwerkproxy und Einheit des Zeitlimits: Node.js verwendet Millisekunden, die übrigen Beispiele Sekunden. Kostenpflichtige Requests nicht blind wiederholen. |
| Keine Bilddaten oder ungültiges Base64 | Prüfen Sie, ob `data[0].b64_json` oder `data[0].url` zurückkommt und ob das Gateway das angefragte Format unterstützt. Entfernen Sie private Daten aus Diagnosen vor dem Teilen. |
| Generierung erfolgreich, Edit fehlgeschlagen | Die Quelldatei muss lesbar sein; der Prompt darf nicht leer sein und das Modell muss `images/edits` unterstützen. Das generierte Original bleibt erhalten. |
| Bild-URL lässt sich nicht laden | Die URL kann abgelaufen oder nicht erreichbar sein. Laden Sie in einem öffentlichen Dienst nicht ungeprüft beliebige externe URLs. |
| C++-/PHP-Problem | CMake muss Header und Bibliotheken von libcurl finden. Die aktive PHP-CLI muss cURL und JSON geladen haben. |

Beginnen Sie mit einem kurzen Prompt, bestätigtem Modellzugriff und einer dokumentierten Größe. Notieren Sie HTTP-Status und gegebenenfalls Request-ID, aber teilen Sie keine Schlüssel, privaten Prompts oder Kundenbilder öffentlich.

## Sicherheit, Datenschutz und Produktiveinsatz

Diese Beispiele sind keine fertige öffentliche Upload-Plattform. Vor einem produktiven Einsatz sollten Sie mindestens folgende Punkte ergänzen:

- API-Schlüssel ausschließlich serverseitig in geschützten Umgebungsvariablen oder einem Secret Manager verwalten; niemals im Browser, Repository, Ticket oder Screenshot veröffentlichen;
- HTTPS nutzen und prüfen, ob die API-Basisadresse auf den vorgesehenen Dienst zeigt;
- Promptlänge, tatsächlichen MIME-Typ, Bildabmessungen und Uploadgröße prüfen sowie Benutzerkontingente, Rate Limits und Inhaltskontrollen einführen;
- lokale Beispieldateien durch dauerhaften Speicher mit Zugriffsschutz, Verschlüsselung, Aufbewahrungsfrist und Löschkonzept ersetzen;
- Prompts, Bilder und temporäre Bild-URLs als potenziell vertraulich behandeln; keine Schlüssel und nur notwendige, redigierte Diagnosedaten protokollieren;
- beim URL-Download erlaubte Hosts, Weiterleitungen, Protokoll, Laufzeit und Dateigröße begrenzen, um SSRF und Ressourcenerschöpfung zu vermeiden;
- Abbruch, Metriken und begrenzte Wiederholungen erst einbauen, wenn Idempotenz und Abrechnung geklärt sind: eine Wiederholung kann eine weitere kostenpflichtige Generierung auslösen;
- aktuelle Datenverarbeitungsbedingungen des verwendeten Dienstes prüfen und nur Bilder übermitteln, zu deren Verarbeitung Sie berechtigt sind.

Aus der API-Form allein lässt sich nicht ableiten, wie ein Dienst Daten verarbeitet oder speichert. Maßgeblich sind die aktuellen Angaben des jeweiligen Anbieters. Weitere projektspezifische Hinweise finden Sie unter [SECURITY.md](../SECURITY.md).

## Dokumentation und nächste Schritte

- [Englische Projektübersicht](../../README.md) · [API-Vertrag](../API_CONTRACT.md) · [Architektur](../ARCHITECTURE.md) · [Sicherheit](../SECURITY.md)
- [Prompt-Beispiele](../../examples/prompts.md) für vergleichbare Versuche in den sechs Sprachimplementierungen.
- Weitere Sprachfassungen: [Chinesisch](README.zh-CN.md), [Russisch](README.ru.md) und [Japanisch](README.ja.md).

Informationen zur Plattform, zum API-Zugang und zu den für Ihr Konto sichtbaren Modellen finden Sie auf der [AI-ROUTER-Website auf Deutsch](https://ai-router.dev/de). Verfügbarkeit und Konditionen richten sich nach Ihrem Konto; dieser Leitfaden garantiert weder Modellzugriff noch ein bestimmtes Ergebnis oder Suchmaschinenranking.
