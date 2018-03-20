Lupa
====

Lupa to Makefile-like narzędzie do automatyzowania buildu/deploymentu.

Składnia Lupafile
-------------------
```
<target>: <dependencies>
    <recipe>
...
```

`<target>`: może zawierać tylko alfanumeryki, underscore i kropki. Jest jednocześnie nazwą pliku potencjalnie generowanego przez target.

`<dependencies>`: oddzielone spacjami targety i/lub pliki. Lupa odróżnia targety lupowe od plików tym, że targety nie zawierają innych znaków niż alfanumeryki, underscore i kropki. Lupa wspiera rekurencyjne selektory do plików (np `./**/*.go` oznacza wszystkie pliki z rozszerzeniem `.go` w tym folderze i we wszystkich podfolderach).

`<recipe>`: zwykły skrypt w bashu. **Wszystkie linijki muszą zaczynać się od tabulatora!**

Można wrzucać komentarze zaczynające się od `#`.

Lupa sprawdza czy plik będący nazwą targetu jest nowszy od dependencji. Jeśli tak, nic nie robi. Wpp wykonuje nowe dependencje, oraz target. Lupa wykonuje wszystko co się da równolegle, żeby to zablokować używamy przełącznika `-s`.

Lupafile domyślnie wykonuje target `all`, ale można też dać target jako argument.

Przykład Lupafile:

```
frontendbuild: ./frontend/src/**/*.*
    rm -rf frontendbuild
    cd frontend
    npm run build --prod --ffast-math -O3
    mv dist ../frontendbuild

# tworzy plik wykonywalny o nazwie "backend"
backend: ./**/*.go
    go build

all: backend frontendbuild
```