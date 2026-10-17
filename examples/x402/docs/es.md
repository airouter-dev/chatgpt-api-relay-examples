# Llamadas de imágenes pagadas con x402

AI ROUTER responde primero con HTTP 402 y condiciones de pago legibles por
máquina. La primera solicitud no cobra: el SDK x402 valida las condiciones,
firma una autorización EVM exacta y repite la misma solicitud. El cliente no
necesita el secreto del comerciante ni el secreto de PayAI.

## Contrato de producción verificado

Usa Base mainnet, USDC y uno de los cuatro modelos configurados: `gpt-image-2`,
`gpt-image-2.5`, `gpt-image-2.5-sunburst` y `gpt-image-2.5-flare`. `n` es obligatorio
(1–4); el precio es `4500 × n` unidades atómicas, es decir, `$0.0045` por imagen.

Generations usa JSON y edits usa multipart con el archivo `image`. Los ejemplos
de Go, Python y Node.js leen `.env` o TOML y empiezan en modo probe.

Guarda la clave privada en un `.env` local ignorado por Git y usa una wallet
dedicada con saldo limitado.

[Página en español](https://ai-router.dev/es) · [Página en inglés](https://ai-router.dev)
