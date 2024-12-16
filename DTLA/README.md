# Datorsistēmu un tīklu loģiskā aizsardzība

Mājaslapa ar servera aizsardzības ieteikumiem un kā konfigurēt tīkla iestatījumus, lietotāju atļaujas, utt.

Projekta web serveris ir rakstīts iekš [Go](https://go.dev/) programmēšanas valodas.
Es izmantoju Go jo ar to ir bijusi vairāk pieredze un Go "[standarta bibliotēkā](https://en.wikipedia.org/wiki/Standard_library)" ir iekļautas daudz noderīgas bibliotēkas, priekš datubāzes piekļuves, hash algoritmiem, "tīkla ligzdām", WebSockets, veidnes, utt.

Ar web serveri saistītās lietas ir iekš `cmd` un `internal` folderos.
Statiskie faili kā html, css, javascript, bildes, veidnes, utt ir iekš `public`.
Veidnes ir fīča no Go standarta bibliotēkas ([html/template](https://pkg.go.dev/html/template)).
Ar veidnēm var ģenerēt HTML uz servera.
Iekš `public/tmpl` ir veidnes, kuras ir izmantotas priekš vairākām lapām, piemēram navigācijas joslas, footer sekcijas, utt. un pārējie HTML faili iekš `public` tās izmanto.

Sadaļā "Ieteikumi" var skatīt ieteikumus, kas ir glabāti Sqlite datubāzē (`db` fails) tabulā `posts`.
Ja ir ielogojies, tad sadaļā "Ieteikumi" varēs rediģēt, dzēst un veidot jaunus rakstus.
Kad ir ielogojies ir izveidota sessija, kas ilgst 5 minūtes, sessijas identifikātors un sākums ir arī glabāts datubāzē tabulā `users`.

## Kā palaist

Web servera komandai var mainīt konfigurāciju, visas opcijas var apskatīties ar `--help`.
Ja neko nemaina tad web serveris būs palaists izmantojot HTTPS protokolu uz adreses 127.0.0.1 un portu 30000.

