# Bootdev Pokedex

Uma Pokedex de linha de comando (CLI) feita em Go que permite explorar áreas do mapa, capturar Pokémon e ver seus status usando a [PokeAPI](https://pokeapi.co/).

Versão em inglês: [README.md](README.md)

## Uso

Inicie o REPL executando:

```bash
go run .
```

## Comandos

| Comando | Descrição |
|---------|-----------|
| `map` | Exibe as próximas 20 áreas de localização. Chamadas sucessivas avançam a paginação. |
| `mapb` | Exibe as 20 áreas de localização anteriores, voltando sua posição. |
| `explore <area>` | Lista todos os Pokémon encontrados em uma área de localização específica. |
| `catch <pokemon>` | Tenta capturar um Pokémon. Quanto maior a experiência base, menor a taxa de captura. |
| `inspect <pokemon>` | Mostra os status e tipos de um Pokémon capturado. |
| `pokedex` | Lista todos os Pokémon que você capturou. |
| `help` | Mostra os comandos disponíveis. |
| `exit` | Sai da CLI da Pokedex. |

## Decisões de Design

Este projeto segue o guia da Pokedex da Boot.dev, mas fiz algumas mudanças intencionais.

### 1) Estado global de comandos em vez de uma struct `Config`

Na versão do curso, uma struct `Config` é passada para cada callback de comando. Nesta implementação, o estado compartilhado é mantido no nível do pacote (`apiClient`, `nextLocationAreasURL`, `prevLocationAreasURL` e `pokedex`).

Escolhi isso porque:

- Remove parâmetros que muitos comandos não usam.
- Deixa as assinaturas dos comandos mais simples (`func(args []string) error`).
- Mantém o acesso ao estado comum direto e legível.

### 2) Modelos focados no domínio em vez de formatos brutos de resposta da API

As respostas da PokeAPI são aninhadas e verbosas. Em vez de trabalhar com esses formatos brutos em toda a aplicação, este projeto mapeia as respostas em modelos menores e focados na tarefa.

Exemplos:

- Definimos structs de resposta apenas com os campos que a CLI usa; campos JSON extras são ignorados durante o unmarshaling.
- `GetLocationAreaPokemon` retorna `[]string` (nomes dos Pokémon) em vez de uma resposta completa e aninhada.
- `Pokemon` usa unmarshaling JSON customizado para achatar status e tipos em campos amigáveis para o REPL.

Isso torna o código do REPL mais fácil de ler e a lógica de interação com a API mais fácil de manter.
