# Appels d’images payants avec x402

AI ROUTER renvoie d’abord HTTP 402 avec des conditions de paiement lisibles par
machine. Le premier appel ne débite rien : le SDK x402 vérifie les conditions,
signe une autorisation EVM exacte, puis rejoue la même requête. Le client n’a
besoin d’aucun secret marchand ni secret PayAI.

## Contrat de production vérifié

Utilisez Base mainnet, USDC et l’un des quatre modèles configurés : `gpt-image-2`,
`gpt-image-2.5`, `gpt-image-2.5-sunburst`, `gpt-image-2.5-flare`. `n` est requis
(1 à 4) et le prix est `4500 × n` unités atomiques, soit `$0.0045` par image.

Generations reçoit du JSON ; edits reçoit un multipart avec le fichier `image`.
Les exemples Go, Python et Node.js lisent `.env` ou TOML et commencent en mode probe.

Conservez la clé privée dans un `.env` local ignoré par Git et utilisez un wallet
dédié avec un faible solde.

[Page française](https://ai-router.dev/fr) · [Page anglaise](https://ai-router.dev)
