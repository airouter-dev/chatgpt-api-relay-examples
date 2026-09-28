# API de geração e edição de imagens com IA: exemplos em seis linguagens

Este guia mostra como criar imagens a partir de texto e depois transformar o resultado usando uma imagem como entrada. O projeto contém clientes de referência em **Python, Node.js, C++, Go, Rust e PHP** que chamam as rotas compatíveis com a Images API da AI-ROUTER. O fluxo é concreto: gerar uma primeira versão, gravá-la em disco e enviá-la a `images/edits` com uma instrução de edição.

> A AI-ROUTER é um serviço independente com uma interface compatível com OpenAI; não é um produto da OpenAI nem afiliada a ela. O nome de um modelo listado neste repositório não garante disponibilidade, preço, quota, formato ou tamanho para toda conta. Confira as opções liberadas na sua conta antes de executar chamadas pagas.

**Acesse a versão em português do site:** [AI-ROUTER em português](https://ai-router.dev/pt). O projeto é um exemplo técnico educacional; ele não promete resultados específicos de geração, disponibilidade de modelos nem posicionamento em mecanismos de busca.

## O que você vai construir

O demo realiza duas operações independentes e encadeadas:

1. Envia um prompt de texto para `POST /v1/images/generations` e recebe a primeira imagem.
2. Materializa a resposta em arquivo local. Os exemplos lidam com resposta `b64_json` e, conforme a implementação e o gateway, com uma resposta `url`.
3. Envia o arquivo salvo para `POST /v1/images/edits` como `multipart/form-data`, junto com um novo prompt.
4. Grava a imagem editada separadamente. Assim, a imagem original permanece disponível para comparação, revisão humana ou uma nova tentativa de edição.

Isso é mais útil do que um exemplo de uma só chamada quando se está avaliando uma API de geração de imagens: permite testar texto-para-imagem e imagem-para-imagem na mesma aplicação e torna explícito o arquivo que conecta as duas etapas. A operação de edição não é simulada; ela usa a imagem criada pela primeira chamada.

## Por que a arquitetura importa

Os exemplos separam responsabilidades para que seja possível reutilizar a integração fora da linha de comando:

| Camada | Responsabilidade |
| --- | --- |
| Configuração | Lê a chave, a URL-base, o modelo, o tamanho e os tempos limite do ambiente ou dos argumentos. |
| Cliente HTTP | Monta JSON ou multipart, autentica, interpreta erros e converte `b64_json`/URL em bytes. |
| Fluxo de trabalho | Coordena geração → gravação → edição → gravação; não implementa transporte HTTP. |
| CLI | Faz a interface de terminal e mostra os caminhos dos arquivos resultantes. |

Essa divisão ajuda a transportar o cliente para um serviço web, worker, fila, aplicação desktop ou processo em lote sem copiar código de rede para a interface. Consulte também [`ARCHITECTURE.md`](../ARCHITECTURE.md), [`API_CONTRACT.md`](../API_CONTRACT.md) e [`SECURITY.md`](../SECURITY.md).

## Escolha a linguagem

Os seis diretórios implementam o mesmo propósito com ferramentas idiomáticas. A escolha depende do ambiente da sua aplicação, não de uma promessa de qualidade diferente do modelo.

| Linguagem | Transporte e ponto de entrada | Bom ponto de partida para |
| --- | --- | --- |
| Python | `urllib`; `python main.py` | scripts, notebooks e automações sem dependência adicional |
| Node.js | `fetch`, `FormData`, `Blob`; `npm run demo --` | serviços JavaScript, tarefas em segundo plano e funções serverless |
| C++ | libcurl; binário `imagegen` | aplicações nativas e ambientes com controle explícito de dependências |
| Go | `net/http`; `go run ./cmd/imagegen` | workers e serviços pequenos com biblioteca padrão |
| Rust | `reqwest`, `serde`, `clap`; `cargo run -- workflow` | clientes tipados e serviços Rust assíncronos |
| PHP | cURL e JSON; `php bin/image-demo.php workflow` | aplicações PHP e tarefas de terminal |

Os exemplos podem ser executados separadamente. Para comparar de forma razoável, use o mesmo modelo, tamanho, prompts e conta; as diferenças de resultado vêm do serviço/modelo e da configuração, não da linguagem do cliente.

## Antes de começar: chave, URL e modelos

Crie uma chave de API na [página da AI-ROUTER em português](https://ai-router.dev/pt) e mantenha-a no ambiente local. O endpoint padrão usado pelos exemplos é `https://api.ai-router.dev/v1`. Não publique a chave no GitHub, em screenshots, em logs de CI nem em aplicações de navegador. Uma aplicação web deve chamar sua própria camada de servidor, não expor o segredo ao usuário.

Exemplo em Bash:

```bash
export AI_ROUTER_API_KEY="sua-chave"
export AI_ROUTER_BASE_URL="https://api.ai-router.dev/v1"
export AI_ROUTER_MODEL="gpt-image-2"
```

Algumas implementações aceitam nomes alternativos de variáveis, como `IMAGE_API_KEY`, `IMAGE_MODEL` ou `OPENAI_API_KEY`; consulte o README de cada linguagem. Os arquivos `.env.example` são referências, mas os CLIs não carregam automaticamente `.env` em todos os runtimes. Configure as variáveis no shell, no gerenciador de segredos ou no ambiente de execução.

Modelos de exemplo incluem `gpt-image-2`, `gpt-image-2.5`, `gpt-image-2.5-flare` e `gpt-image-2.5-burst`. Eles são passados como identificadores configuráveis. **Não presuma que todos estão habilitados para todas as contas.** Para consultar o catálogo visível à sua chave:

```bash
curl https://api.ai-router.dev/v1/models \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY"
```

O endpoint de modelos e a compatibilidade de imagem dependem da configuração atual da conta. Se um modelo não aparecer ou retornar erro de permissão, use um modelo liberado para sua chave. Tamanho e qualidade também devem ser escolhidos de acordo com a documentação e os limites atuais do modelo.

## Executar os seis exemplos

Os comandos abaixo ilustram o fluxo completo. Execute cada bloco a partir do diretório indicado. Em Windows PowerShell, defina variáveis com `$env:AI_ROUTER_API_KEY = "sua-chave"`; no C++, instale também CMake, compilador C++17 e os arquivos de desenvolvimento do libcurl.

### Python

Requer Python 3.10 ou superior e não instala dependências de execução adicionais:

```bash
cd marketing/image-generation-multilang/python
export IMAGE_API_KEY="$AI_ROUTER_API_KEY"
export AI_ROUTER_BASE_URL="https://api.ai-router.dev/v1"
python main.py \
  --prompt "Uma estufa de vidro numa cobertura chuvosa, fotografia editorial" \
  --edit-prompt "Aqueça a luz para o pôr do sol sem alterar a arquitetura" \
  --output-dir ./outputs
python -m unittest discover -s tests -v
```

O fluxo salva `outputs/generated.png` e `outputs/edited.png`. Também é possível usar `IMAGE_MODEL`, `IMAGE_SIZE`, `IMAGE_QUALITY` e as opções correspondentes da CLI.

### Node.js

Requer Node.js 18.17 ou superior; o cliente usa APIs nativas do Node e não precisa de uma biblioteca HTTP adicional:

```bash
cd marketing/image-generation-multilang/nodejs
export AI_ROUTER_API_KEY="sua-chave"
npm test
npm run demo -- \
  --prompt "Uma estufa de vidro numa cobertura chuvosa" \
  --edit-prompt "Adicione luz dourada de fim de tarde" \
  --model gpt-image-2
```

O resultado fica em `outputs/generated.png` e `outputs/edited.png`. `--size`, `--quality`, `--base-url`, `--timeout` e `--output-dir` também podem ser configurados.

### C++

Requer C++17, CMake 3.20+ e libcurl. Com a dependência instalada:

```bash
cd marketing/image-generation-multilang/cpp
export AI_ROUTER_API_KEY="sua-chave"
cmake -S . -B build -DBUILD_TESTING=ON
cmake --build build --config Release
ctest --test-dir build --output-on-failure
./build/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "Foto de produto de uma nave de vidro sobre uma mesa coberta de musgo" \
  --edit-prompt "Adicione uma luz suave de nascer do sol, preservando a silhueta"
```

No Windows, execute o binário gerado na configuração do CMake, normalmente em `build/Release/imagegen.exe`; o caminho pode variar conforme o gerador. O CLI aceita `--output` e `--edited-output` para escolher destinos diferentes.

### Go

Requer uma versão de Go compatível com `go.mod`:

```bash
cd marketing/image-generation-multilang/go
export AI_ROUTER_API_KEY="sua-chave"
go test ./...
go vet ./...
go run ./cmd/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "Uma estufa de vidro em Marte ao nascer do sol" \
  --edit-prompt "Inclua pequenas luzes quentes sem mudar a composição"
```

Por padrão, o CLI grava `generated.png` e `edited.png`. `--output`, `--edited-output`, `--size`, `--timeout` e `--response-format` estão disponíveis.

### Rust

Requer Rust estável 1.80 ou superior. A subcomando `workflow` executa as duas chamadas em sequência:

```bash
cd marketing/image-generation-multilang/rust
export AI_ROUTER_API_KEY="sua-chave"
cargo test
cargo run --release -- workflow \
  --prompt "Uma mochila vermelha em uma fotografia de produto limpa" \
  --edit-prompt "Mude a mochila para azul-cobalto e preserve o enquadramento" \
  --generated output/generated.png \
  --edited output/edited.png
```

Também existem os subcomandos `generate` e `edit`, úteis para inspecionar cada fase. As configurações `AI_ROUTER_BASE_URL`, `AI_ROUTER_MODEL` e `AI_ROUTER_SIZE` podem ser lidas do ambiente.

### PHP

Requer PHP 8.1+, extensões `curl` e `json`; não há dependência de runtime adicional para o demo:

```bash
cd marketing/image-generation-multilang/php
export AI_ROUTER_API_KEY="sua-chave"
php bin/image-demo.php workflow \
  --prompt "Uma mochila vermelha em uma fotografia de produto limpa" \
  --edit-prompt "Mude a mochila para azul-cobalto e preserve o enquadramento" \
  --generated output/generated.png \
  --edited output/edited.png
```

Os comandos `generate` e `edit` também podem ser chamados separadamente. O cliente usa `AI_ROUTER_BASE_URL`, `AI_ROUTER_MODEL` e `AI_ROUTER_SIZE` como padrões de ambiente.

## Contrato das chamadas de imagem

A geração usa JSON e autenticação Bearer. Um corpo típico, com campos que os exemplos entendem, é:

```http
POST https://api.ai-router.dev/v1/images/generations
Authorization: Bearer $AI_ROUTER_API_KEY
Content-Type: application/json
```

```json
{
  "model": "gpt-image-2",
  "prompt": "Uma pequena estufa de vidro numa cobertura chuvosa",
  "size": "1024x1024",
  "n": 1,
  "response_format": "b64_json"
}
```

O resultado bem-sucedido tem uma lista `data`; o primeiro elemento pode conter `b64_json` ou, conforme a configuração do serviço, `url`. Os clientes convertem a resposta em bytes e a gravam em um caminho escolhido.

Para editar, o cliente envia `multipart/form-data` a `POST /v1/images/edits`. Os campos incluem `model`, `prompt`, `size` e a imagem anterior no campo `image`; o formato de resposta usado pelos demos é `b64_json`. A camada HTTP constrói o boundary multipart — não defina manualmente um `Content-Type` incompleto. Detalhes comuns estão em [`API_CONTRACT.md`](../API_CONTRACT.md).

## Escrever prompts que sejam mais fáceis de avaliar

Uma boa comparação precisa de instruções específicas e repetíveis. No primeiro prompt, defina o assunto, o cenário, o enquadramento, a iluminação e a finalidade visual, por exemplo: “foto editorial de uma pequena estufa de vidro numa cobertura urbana depois da chuva, enquadramento amplo, luz fria da manhã”. No prompt de edição, diga o que deve mudar e o que precisa permanecer (“adicione reflexos quentes nas janelas; preserve o formato da estufa e a posição da câmera”).

Faça uma alteração por etapa, mantenha a primeira imagem e compare as duas saídas antes de fazer outra edição. Evite colocar dados pessoais, informações confidenciais ou instruções contraditórias no prompt. A saída depende do modelo, da conta e dos parâmetros aceitos; instruções claras ajudam a testar, mas não garantem um resultado exato.

## Resultados, latência e custos

Uma execução completa faz **duas requisições** potencialmente faturáveis: uma para gerar e outra para editar. Confirme preços, cotas, tamanhos e política de dados na sua conta antes de automatizar execuções. O demo não implementa repetição automática: uma nova tentativa pode criar outra geração e eventual custo. A etapa de geração é salva antes da edição; se a edição falhar, o arquivo inicial não deve ser considerado perdido.

As imagens são gravadas localmente em pastas `output`, `outputs` ou no caminho passado ao CLI, conforme a linguagem. Revise o `.gitignore` antes de compartilhar arquivos e não publique resultados privados ou com direitos de terceiros sem autorização. O nome do arquivo com extensão `.png` é o padrão dos exemplos; valide o conteúdo/MIME do gateway antes de tratar o formato como garantido em um produto.

## Solução de problemas

| Sintoma | O que verificar |
| --- | --- |
| `401` ou `403` | A variável está definida no mesmo terminal? A chave continua válida e habilitada para esse modelo? Não inclua espaços extras nem use aspas dentro do valor. |
| Modelo não encontrado ou sem acesso | Consulte `/v1/models` usando a mesma chave e selecione um identificador disponível para a conta. |
| Tamanho ou qualidade recusados | Volte a `1024x1024` e `auto` quando aplicável; confirme as opções suportadas pela conta/modelo. |
| Resposta sem imagem | Inspecione o status, a mensagem de erro limitada e o ID de requisição, se o serviço o fornecer. Verifique se a resposta contém `data[0].b64_json` ou uma URL que a implementação aceite. |
| Erro apenas na edição | Confira se o primeiro arquivo existe e é legível, se o campo enviado é `image` e se o gateway aceita o tipo/tamanho do arquivo. A imagem gerada é preservada para nova tentativa. |
| Timeout | A geração pode demorar; ajuste o timeout suportado pelo CLI e verifique conectividade. Evite repetir cegamente uma solicitação que pode ter sido concluída no servidor. |
| CMake não encontra curl ou o compilador | Instale o compilador C++17, CMake e os arquivos de desenvolvimento do libcurl; no Windows, confira o toolchain/vcpkg selecionado. |

Os testes locais de protocolo usam fixtures ou servidores de teste e não exigem uma chave real; eles não comprovam que uma conta tenha acesso a um modelo. Execute os testes específicos de cada diretório antes de integrar o cliente.

## Segurança e privacidade

- Mantenha chaves fora do código, do Git e do frontend. Use variáveis de ambiente em desenvolvimento e um cofre de segredos em produção.
- Use HTTPS. Não encaminhe o cabeçalho Bearer para uma URL de imagem de outro host; os clientes evitam propagar credenciais nos downloads de URL.
- Trate prompts e imagens como dados potencialmente sensíveis. Minimize logs, aplique retenção limitada e explique ao usuário para onde os arquivos são enviados.
- Valide MIME, dimensões e tamanho de upload no seu servidor. O exemplo educacional não é um serviço público de upload nem substitui autenticação, moderação, cotas ou proteção contra abuso.
- Confirme direitos da imagem de entrada, políticas do provedor e regras locais. Revise resultados gerados antes de publicá-los.

## Evoluir do demo para produção

Reutilize a camada cliente e substitua a CLI por um handler autenticado. Armazene imagens em object storage com política de ciclo de vida, aplique limites por usuário, timeout e cancelamento, registre IDs de requisição sem segredos, adicione observabilidade e retries limitados somente quando forem apropriados. Uma fila pode desacoplar geração lenta do pedido web; uma camada de armazenamento também evita depender do disco efêmero de containers. Teste erros e limites com respostas simuladas antes de habilitar gastos reais.

## Explore e contribua

Se você procura uma API compatível com OpenAI para integrar geração e edição de imagens em diferentes linguagens, explore a [AI-ROUTER em português](https://ai-router.dev/pt) para obter informações da plataforma, acesso à conta e disponibilidade atual. O endpoint de API usado neste exemplo é `https://api.ai-router.dev/v1`.

Este guia e o código são materiais educacionais, não uma avaliação comparativa independente nem uma garantia de SEO, qualidade, custo ou disponibilidade. Ao referenciar o exemplo, mantenha a atribuição e o link úteis e contextuais; links editoriais honestos podem ajudar leitores a descobrir o projeto, mas nenhum texto consegue prometer melhor ranking. Para adicionar uma linguagem, siga [`CONTRIBUTING.md`](../../CONTRIBUTING.md), preserve o contrato comum e inclua testes offline.
