# Remote Home Assistant 1.0.2

[English](#english) · [Polski](#polski) · [Deutsch](#deutsch) · [Nederlands](#nederlands) · [Français](#français) · [Svenska](#svenska)

## English

A secure way back to your Home Assistant. Connect through an encrypted tunnel, choose who gets through, and keep control of your home.

### Register and install

Create your account at [remote-ha.com/register](https://remote-ha.com/en/register), verify your email, and open the [dashboard](https://remote-ha.com/en/dashboard). Under **Add a home**, choose a name and an available address, then select **Add installation**.

Requires **Home Assistant OS**, an administrator account and an **aarch64 (ARM64)** or **amd64 (x86-64)** device. Home Assistant Container does not include the app store. Before Home Assistant 2026.2, **Apps** was called **Add-ons**.

**Install with a button:** [Add the Remote HA repository to Home Assistant](https://my.home-assistant.io/redirect/supervisor_addon_repository/?repository=https%3A%2F%2Fgithub.com%2Fremote-ha%2Fremote-ha-addon). Choose your Home Assistant instance, confirm adding the repository, then find **Remote Home Assistant** under **Settings → Apps → Install app** and select **Install**.

**Add the repository manually:**

1. Open **Settings → Apps → Install app**.
2. Open the **⋮** menu at the top right and choose **Repositories**.
3. Paste `https://github.com/remote-ha/remote-ha-addon` and select **Add**.
4. Close the dialog, find **Remote Home Assistant** in **Remote Home Assistant Apps**, and select **Install**. Refresh the store or select **Check for updates** if it does not appear.

The app store downloads the app from the [public GitHub repository](https://github.com/remote-ha/remote-ha-addon); no ZIP download or manual file copying is needed.

**HACS:** HACS does not install Home Assistant apps (formerly add-ons). This repository is an app, not a HACS integration. Use the app-store methods above. See the [HACS FAQ](https://www.hacs.xyz/docs/faq/addons/) and the [full installation guide](https://remote-ha.com/en/docs).

After signing in, follow the setup steps in the dashboard. Select a home to see its overview, access firewall, activity and connection settings. Your subscription is under **Account & billing**. The language and appearance controls are in the top bar; on small screens, appearance is also in the menu.

### Pair your home

Install the app before generating a pairing token. In the dashboard, select **Pairing token** on your installation. In Home Assistant, open **Settings → Apps → Remote Home Assistant → Configuration**, paste the token into **Token** (`pairing_token` in the YAML editor) and **Save**. On **Info**, select **Start**, or **Restart** if already running. Enable **Start on boot**. Wait for **Connected** in the dashboard, open the installation address and use your usual Home Assistant login.

Manage who can connect and which paths are allowed in the Remote HA dashboard. No router settings or certificates are required.

Your connection is remembered after a restart or update. To pair again, generate a fresh token, replace the saved token and restart the app. If a token has expired or was already used, generate another one. Keep tokens private.

If the connection is unavailable, check your internet connection and the installation’s status in the dashboard. Before updating, back up the app and keep its data. Your Home Assistant login remains separate from your Remote HA account.

Configuration labels follow your Home Assistant language. This documentation is available below in all six languages. Remote HA and its connection provider process traffic to deliver the service; see [Security](https://remote-ha.com/en/security) and [Privacy](https://remote-ha.com/en/privacy).

Need help? [contact@remote-ha.com](mailto:contact@remote-ha.com). Never include your token in a support request.

## Polski

Bezpieczny dostęp do Twojego Home Assistanta z dowolnego miejsca. Łącz się przez szyfrowany tunel, decyduj, kogo wpuszczasz, i zachowaj kontrolę nad swoim domem.

### Rejestracja i instalacja

Utwórz konto na [remote-ha.com/register](https://remote-ha.com/pl/register), potwierdź adres e-mail i otwórz [panel](https://remote-ha.com/pl/dashboard). W sekcji **Dodaj dom** wpisz nazwę, wybierz dostępny adres i kliknij **Dodaj instalację**.

Wymagane są **Home Assistant OS**, konto administratora i urządzenie **aarch64 (ARM64)** lub **amd64 (x86-64)**. Home Assistant Container nie zawiera sklepu aplikacji. Przed wersją Home Assistant 2026.2 menu **Aplikacje** nosiło nazwę **Dodatki**.

**Instalacja przyciskiem:** [Dodaj repozytorium Remote HA do Home Assistant](https://my.home-assistant.io/redirect/supervisor_addon_repository/?repository=https%3A%2F%2Fgithub.com%2Fremote-ha%2Fremote-ha-addon). Wybierz swoją instancję Home Assistant i potwierdź dodanie repozytorium. Następnie w **Ustawienia → Aplikacje → Zainstaluj aplikację** znajdź **Remote Home Assistant** i wybierz **Zainstaluj**.

**Ręczne dodanie repozytorium:**

1. Otwórz **Ustawienia → Aplikacje → Zainstaluj aplikację**.
2. W menu **⋮** w prawym górnym rogu wybierz **Repozytoria**.
3. Wklej `https://github.com/remote-ha/remote-ha-addon` i wybierz **Dodaj**.
4. Zamknij okno, znajdź **Remote Home Assistant** w sekcji **Remote Home Assistant Apps** i kliknij **Zainstaluj**. Jeśli aplikacji nie widać, odśwież sklep lub wybierz **Sprawdź aktualizacje**.

Sklep pobiera aplikację z [publicznego repozytorium GitHub](https://github.com/remote-ha/remote-ha-addon). Nie trzeba pobierać pliku ZIP ani ręcznie kopiować plików.

**HACS:** HACS nie instaluje aplikacji Home Assistant (dawniej dodatków). To repozytorium zawiera aplikację, a nie integrację HACS. Użyj jednej z metod instalacji przez sklep opisanych powyżej. Zobacz [wyjaśnienie HACS](https://www.hacs.xyz/docs/faq/addons/) oraz [pełną instrukcję instalacji](https://remote-ha.com/pl/docs).

Po zalogowaniu przejdź przez kroki konfiguracji w panelu. Wybierz dom, aby zobaczyć przegląd, zaporę dostępu, aktywność i ustawienia połączenia. Abonament znajdziesz w sekcji **Konto i abonament**. Język i wygląd zmienisz w górnym pasku; na małych ekranach wygląd jest też w menu.

### Parowanie z kontem

Najpierw zainstaluj aplikację, a następnie w panelu przy swojej instalacji kliknij **Token parowania**. W Home Assistant otwórz **Ustawienia → Aplikacje → Remote Home Assistant → Konfiguracja**, wklej token w pole **Token** (`pairing_token` w edytorze YAML) i zapisz. Na karcie informacji wybierz **Uruchom** lub **Uruchom ponownie**, jeśli aplikacja już działa. Włącz uruchamianie przy starcie systemu. Poczekaj na status **Połączono** w panelu, otwórz adres instalacji i zaloguj się swoim kontem Home Assistant.

W panelu Remote HA określisz, kto może się łączyć i które ścieżki są dozwolone. Nie musisz konfigurować routera ani certyfikatów.

Połączenie jest zapamiętywane po restarcie i aktualizacji. Aby sparować ponownie, wygeneruj nowy token, zastąp zapisany token i zrestartuj aplikację. Jeśli token wygasł lub został już użyty, wygeneruj kolejny. Zachowaj tokeny w tajemnicy.

Jeśli połączenie jest niedostępne, sprawdź dostęp do internetu oraz stan instalacji w panelu. Przed aktualizacją wykonaj kopię zapasową aplikacji i zachowaj jej dane. Konto Home Assistant pozostaje niezależne od konta Remote HA.

Język etykiet konfiguracji odpowiada językowi Home Assistant. Dokumentacja zawiera wszystkie sześć wersji językowych. Remote HA i dostawca połączenia przetwarzają ruch w celu świadczenia usługi; zobacz [bezpieczeństwo](https://remote-ha.com/pl/security) i [prywatność](https://remote-ha.com/pl/privacy).

Potrzebujesz pomocy? [contact@remote-ha.com](mailto:contact@remote-ha.com). Nie dołączaj tokenu do zgłoszenia.

## Deutsch

Ein sicherer Weg zurück zu deinem Home Assistant. Verbinde dich über einen verschlüsselten Tunnel, entscheide, wer Zugang erhält, und behalte die Kontrolle über dein Zuhause.

### Konto erstellen und installieren

Erstelle dein Konto auf [remote-ha.com/register](https://remote-ha.com/de/register), bestätige deine E-Mail-Adresse und öffne das [Dashboard](https://remote-ha.com/de/dashboard). Gib unter **Zuhause hinzufügen** einen Namen ein, wähle eine verfügbare Adresse und klicke auf **Installation hinzufügen**.

Du benötigst **Home Assistant OS**, ein Administratorkonto und ein Gerät mit **aarch64 (ARM64)** oder **amd64 (x86-64)**. Home Assistant Container enthält keinen App-Store. Vor Home Assistant 2026.2 hieß **Apps** noch **Add-ons**.

**Installation per Schaltfläche:** [Remote-HA-Repository zu Home Assistant hinzufügen](https://my.home-assistant.io/redirect/supervisor_addon_repository/?repository=https%3A%2F%2Fgithub.com%2Fremote-ha%2Fremote-ha-addon). Wähle deine Home-Assistant-Instanz, bestätige das Repository und suche unter **Einstellungen → Apps → App installieren** nach **Remote Home Assistant**. Wähle **Installieren**.

**Repository manuell hinzufügen:**

1. Öffne **Einstellungen → Apps → App installieren**.
2. Wähle im Menü **⋮** oben rechts **Repositories**.
3. Füge `https://github.com/remote-ha/remote-ha-addon` ein und wähle **Hinzufügen**.
4. Schließe den Dialog, suche **Remote Home Assistant** unter **Remote Home Assistant Apps** und wähle **Installieren**. Aktualisiere die Seite oder suche im Store nach Updates, falls die App fehlt.

Der App-Store lädt die App aus dem [öffentlichen GitHub-Repository](https://github.com/remote-ha/remote-ha-addon) herunter. Du musst keine ZIP-Datei herunterladen oder Dateien manuell kopieren.

**HACS:** HACS installiert keine Home-Assistant-Apps (früher Add-ons). Dieses Repository enthält eine App und keine HACS-Integration. Verwende die oben beschriebenen Methoden im App-Store. Siehe die [HACS-FAQ](https://www.hacs.xyz/docs/faq/addons/) und die [vollständige Installationsanleitung](https://remote-ha.com/de/docs).

Folge nach dem Anmelden den Einrichtungsschritten im Dashboard. Wähle ein Zuhause aus, um Übersicht, Zugriffsfirewall, Aktivitäten und Verbindungseinstellungen zu öffnen. Dein Abonnement findest du unter **Konto & Abonnement**. Sprache und Darstellung lassen sich in der oberen Leiste ändern; auf kleinen Bildschirmen ist die Darstellung auch im Menü verfügbar.

### Dein Zuhause koppeln

Installiere die App, bevor du einen Token erzeugst. Wähle im Dashboard bei deiner Installation **Kopplungstoken**. Öffne in Home Assistant **Einstellungen → Apps → Remote Home Assistant → Konfiguration**, füge den Token unter **Token** ein (`pairing_token` im YAML-Editor) und speichere. Wähle auf der Infoseite **Starten** oder **Neu starten**, falls die App bereits läuft. Aktiviere den Start beim Systemstart. Warte im Dashboard auf **Verbunden**, öffne die Adresse deiner Installation und melde dich mit deinem Home-Assistant-Konto an.

Im Remote-HA-Dashboard legst du fest, wer sich verbinden darf und welche Pfade erlaubt sind. Du musst weder deinen Router konfigurieren noch Zertifikate einrichten.

Die Verbindung bleibt nach einem Neustart oder Update gespeichert. Für eine erneute Kopplung erstellst du einen neuen Token, ersetzt den gespeicherten Token und startest die App neu. Ist ein Token abgelaufen oder bereits verwendet worden, erstelle einen weiteren. Halte deine Tokens geheim.

Falls die Verbindung nicht verfügbar ist, prüfe deine Internetverbindung und den Status der Installation im Dashboard. Sichere die App vor einem Update und behalte ihre Daten. Dein Home-Assistant-Konto bleibt von deinem Remote-HA-Konto getrennt.

Die Sprache der Konfigurationsfelder richtet sich nach Home Assistant. Diese Dokumentation enthält alle sechs Sprachversionen. Remote HA und sein Verbindungsanbieter verarbeiten den Datenverkehr, um den Dienst bereitzustellen; siehe [Sicherheit](https://remote-ha.com/de/security) und [Datenschutz](https://remote-ha.com/de/privacy).

Brauchst du Hilfe? [contact@remote-ha.com](mailto:contact@remote-ha.com). Gib deinen Token niemals in einer Supportanfrage an.

## Nederlands

Een veilige weg terug naar je Home Assistant. Maak verbinding via een versleutelde tunnel, bepaal wie toegang krijgt en houd de controle over je huis.

### Registreren en installeren

Maak je account aan op [remote-ha.com/register](https://remote-ha.com/nl/register), bevestig je e-mailadres en open het [dashboard](https://remote-ha.com/nl/dashboard). Vul bij **Huis toevoegen** een naam in, kies een beschikbaar adres en klik op **Installatie toevoegen**.

Je hebt **Home Assistant OS**, een beheerdersaccount en een apparaat met **aarch64 (ARM64)** of **amd64 (x86-64)** nodig. Home Assistant Container bevat geen appwinkel. Vóór Home Assistant 2026.2 heette **Apps** nog **Add-ons**.

**Installeren met een knop:** [Voeg de Remote HA-repository toe aan Home Assistant](https://my.home-assistant.io/redirect/supervisor_addon_repository/?repository=https%3A%2F%2Fgithub.com%2Fremote-ha%2Fremote-ha-addon). Kies je Home Assistant-instantie, bevestig de repository en zoek **Remote Home Assistant** via **Instellingen → Apps → App installeren**. Kies **Installeren**.

**De repository handmatig toevoegen:**

1. Open **Instellingen → Apps → App installeren**.
2. Kies rechtsboven in het menu **⋮** voor **Repositories**.
3. Plak `https://github.com/remote-ha/remote-ha-addon` en kies **Toevoegen**.
4. Sluit het venster, zoek **Remote Home Assistant** onder **Remote Home Assistant Apps** en kies **Installeren**. Vernieuw de winkel of controleer op updates als de app ontbreekt.

De appwinkel downloadt de app uit de [openbare GitHub-repository](https://github.com/remote-ha/remote-ha-addon). Je hoeft geen ZIP-bestand te downloaden of bestanden handmatig te kopiëren.

**HACS:** HACS installeert geen Home Assistant-apps (voorheen add-ons). Deze repository bevat een app, geen HACS-integratie. Gebruik de bovenstaande methoden via de appwinkel. Zie de [HACS-FAQ](https://www.hacs.xyz/docs/faq/addons/) en de [volledige installatiehandleiding](https://remote-ha.com/nl/docs).

Volg na het aanmelden de stappen in het dashboard. Selecteer een huis om het overzicht, de toegangsfirewall, activiteit en verbindingsinstellingen te bekijken. Je abonnement staat onder **Account en abonnement**. Taal en weergave wijzig je in de bovenste balk; op kleine schermen staat weergave ook in het menu.

### Je huis koppelen

Installeer de app voordat je een token aanmaakt. Kies in het dashboard **Koppeltoken** bij je installatie. Open in Home Assistant **Instellingen → Apps → Remote Home Assistant → Configuratie**, plak het token in **Token** (`pairing_token` in de YAML-editor) en sla op. Kies op de informatiepagina **Starten** of **Herstarten** als de app al actief is. Schakel starten bij het opstarten in. Wacht op **Verbonden** in het dashboard, open het installatieadres en meld je aan met je Home Assistant-account.

In het Remote HA-dashboard bepaal je wie verbinding mag maken en welke paden zijn toegestaan. Je hoeft geen routerinstellingen te wijzigen of certificaten in te stellen.

Je verbinding blijft bewaard na een herstart of update. Wil je opnieuw koppelen, maak dan een nieuw token aan, vervang het opgeslagen token en herstart de app. Is een token verlopen of al gebruikt, maak dan een ander token aan. Houd je tokens geheim.

Als de verbinding niet beschikbaar is, controleer dan je internetverbinding en de status van de installatie in het dashboard. Maak vóór een update een back-up van de app en bewaar de gegevens. Je Home Assistant-account blijft gescheiden van je Remote HA-account.

De taal van de configuratievelden volgt de taal van Home Assistant. Deze documentatie bevat alle zes taalversies. Remote HA en zijn verbindingsprovider verwerken verkeer om de dienst te leveren; zie [Beveiliging](https://remote-ha.com/nl/security) en [Privacy](https://remote-ha.com/nl/privacy).

Hulp nodig? [contact@remote-ha.com](mailto:contact@remote-ha.com). Vermeld je token nooit in een supportverzoek.

## Français

Retrouvez votre Home Assistant en toute sécurité. Connectez-vous via un tunnel chiffré, choisissez qui peut y accéder et gardez le contrôle de votre maison.

### Créer un compte et installer

Créez votre compte sur [remote-ha.com/register](https://remote-ha.com/fr/register), confirmez votre adresse e-mail et ouvrez le [tableau de bord](https://remote-ha.com/fr/dashboard). Dans **Ajouter une maison**, saisissez un nom, choisissez une adresse disponible et cliquez sur **Ajouter une installation**.

Vous devez disposer de **Home Assistant OS**, d’un compte administrateur et d’un appareil **aarch64 (ARM64)** ou **amd64 (x86-64)**. Home Assistant Container ne comprend pas de boutique d’applications. Avant Home Assistant 2026.2, les **Apps** s’appelaient **Add-ons**.

**Installation avec un bouton :** [Ajouter le dépôt Remote HA à Home Assistant](https://my.home-assistant.io/redirect/supervisor_addon_repository/?repository=https%3A%2F%2Fgithub.com%2Fremote-ha%2Fremote-ha-addon). Choisissez votre instance Home Assistant, confirmez l’ajout du dépôt, puis recherchez **Remote Home Assistant** dans **Paramètres → Applications → Installer une application** et sélectionnez **Installer**.

**Ajouter le dépôt manuellement :**

1. Ouvrez **Paramètres → Applications → Installer une application**.
2. Dans le menu **⋮** en haut à droite, sélectionnez **Dépôts**.
3. Collez `https://github.com/remote-ha/remote-ha-addon` et sélectionnez **Ajouter**.
4. Fermez la fenêtre, recherchez **Remote Home Assistant** dans **Remote Home Assistant Apps** et sélectionnez **Installer**. Actualisez la boutique ou recherchez les mises à jour si l’application n’apparaît pas.

La boutique télécharge l’application depuis le [dépôt GitHub public](https://github.com/remote-ha/remote-ha-addon). Aucun téléchargement ZIP ni copie manuelle de fichiers n’est nécessaire.

**HACS :** HACS n’installe pas les applications Home Assistant (anciennement add-ons). Ce dépôt contient une application et non une intégration HACS. Utilisez les méthodes décrites ci-dessus dans la boutique. Consultez la [FAQ HACS](https://www.hacs.xyz/docs/faq/addons/) et le [guide d’installation complet](https://remote-ha.com/fr/docs).

Après la connexion, suivez les étapes de configuration dans le tableau de bord. Sélectionnez une maison pour consulter sa vue d’ensemble, son pare-feu d’accès, son activité et ses paramètres de connexion. Votre abonnement se trouve dans **Compte et abonnement**. La langue et l’apparence se changent dans la barre supérieure ; sur petit écran, l’apparence est aussi accessible dans le menu.

### Appairer votre maison

Installez l’application avant de générer un jeton. Dans le tableau de bord, sélectionnez **Jeton d’appairage** pour votre installation. Dans Home Assistant, ouvrez **Paramètres → Applications → Remote Home Assistant → Configuration**, collez le jeton dans **Jeton** (`pairing_token` dans l’éditeur YAML) et enregistrez. Sur la page d’informations, sélectionnez **Démarrer**, ou **Redémarrer** si l’application fonctionne déjà. Activez le démarrage automatique. Attendez le statut **Connecté** dans le tableau de bord, ouvrez l’adresse de l’installation et connectez-vous avec votre compte Home Assistant.

Dans le tableau de bord Remote HA, choisissez qui peut se connecter et quels chemins sont autorisés. Vous n’avez pas besoin de configurer votre routeur ni de créer des certificats.

Votre connexion est conservée après un redémarrage ou une mise à jour. Pour effectuer un nouvel appairage, générez un nouveau jeton, remplacez le jeton enregistré et redémarrez l’application. Si un jeton a expiré ou a déjà été utilisé, générez-en un autre. Gardez vos jetons secrets.

Si la connexion est indisponible, vérifiez votre connexion Internet et l’état de l’installation dans le tableau de bord. Avant une mise à jour, sauvegardez l’application et conservez ses données. Votre compte Home Assistant reste distinct de votre compte Remote HA.

La langue des champs de configuration suit celle de Home Assistant. Cette documentation contient les six versions linguistiques. Remote HA et son fournisseur de connexion traitent le trafic pour fournir le service ; consultez les pages [Sécurité](https://remote-ha.com/fr/security) et [Confidentialité](https://remote-ha.com/fr/privacy).

Besoin d’aide ? [contact@remote-ha.com](mailto:contact@remote-ha.com). Ne communiquez jamais votre jeton dans une demande d’assistance.


## Svenska

Säker åtkomst till din Home Assistant, var du än är. Anslut via en krypterad tunnel, välj vem som får åtkomst och behåll kontrollen över ditt hem.

### Skapa ett konto och installera

Skapa ditt konto på [remote-ha.com/register](https://remote-ha.com/sv/register), verifiera din e-postadress och öppna [kontrollpanelen](https://remote-ha.com/sv/dashboard). Under **Lägg till ett hem** anger du ett namn, väljer en ledig adress och klickar på **Lägg till installation**.

Du behöver **Home Assistant OS**, ett administratörskonto och en enhet med **aarch64 (ARM64)** eller **amd64 (x86-64)**. Home Assistant Container innehåller ingen appbutik. Före Home Assistant 2026.2 hette **Appar** i stället **Tillägg**.

**Installera med en knapp:** [Lägg till Remote HA:s arkiv i Home Assistant](https://my.home-assistant.io/redirect/supervisor_addon_repository/?repository=https%3A%2F%2Fgithub.com%2Fremote-ha%2Fremote-ha-addon). Välj din Home Assistant-instans, bekräfta att arkivet ska läggas till och leta sedan upp **Remote Home Assistant** under **Inställningar → Appar → Installera app**. Välj **Installera**.

**Lägg till arkivet manuellt:**

1. Öppna **Inställningar → Appar → Installera app**.
2. Öppna menyn **⋮** uppe till höger och välj **Arkiv**.
3. Klistra in `https://github.com/remote-ha/remote-ha-addon` och välj **Lägg till**.
4. Stäng dialogrutan, leta upp **Remote Home Assistant** under **Remote Home Assistant Apps** och välj **Installera**. Uppdatera butiken eller välj **Sök efter uppdateringar** om appen inte visas.

Appbutiken hämtar appen från det [offentliga GitHub-arkivet](https://github.com/remote-ha/remote-ha-addon). Du behöver inte hämta någon ZIP-fil eller kopiera filer manuellt.

**HACS:** HACS installerar inte Home Assistant-appar (tidigare tillägg). Det här arkivet innehåller en app, inte en HACS-integration. Använd någon av metoderna för appbutiken ovan. Läs [HACS vanliga frågor](https://www.hacs.xyz/docs/faq/addons/) och den [fullständiga installationsguiden](https://remote-ha.com/sv/docs).

Följ installationsstegen i kontrollpanelen när du har loggat in. Välj ett hem för att se översikt, åtkomstbrandvägg, aktivitet och anslutningsinställningar. Din prenumeration finns under **Konto och prenumeration**. Språk och utseende ändras i det övre fältet; på små skärmar finns utseende även i menyn.

### Parkoppla ditt hem

Installera appen innan du skapar en parkopplingstoken. Välj **Parkopplingstoken** vid din installation i kontrollpanelen. I Home Assistant öppnar du **Inställningar → Appar → Remote Home Assistant → Konfiguration**, klistrar in token i **Token** (`pairing_token` i YAML-redigeraren) och väljer **Spara**. På informationssidan väljer du **Starta**, eller **Starta om** om appen redan körs. Aktivera **Starta vid uppstart**. Vänta tills **Ansluten** visas i kontrollpanelen, öppna installationens adress och logga in med ditt vanliga Home Assistant-konto.

Du hanterar vem som får ansluta och vilka sökvägar som tillåts i Remote HA:s kontrollpanel. Du behöver inte konfigurera routern eller hantera certifikat.

Din anslutning sparas efter en omstart eller uppdatering. Om du vill parkoppla på nytt skapar du en ny token, ersätter den sparade token och startar om appen. Skapa en ny token om den har gått ut eller redan har använts. Håll dina token hemliga.

Om anslutningen inte fungerar kontrollerar du internetanslutningen och installationens status i kontrollpanelen. Säkerhetskopiera appen före en uppdatering och behåll dess data. Din inloggning till Home Assistant är separat från ditt Remote HA-konto.

Konfigurationsfältens språk följer språket i Home Assistant. Den här dokumentationen finns på sex språk. Remote HA och dess anslutningsleverantör behandlar trafik för att tillhandahålla tjänsten. Läs mer om [säkerhet](https://remote-ha.com/sv/security) och [integritet](https://remote-ha.com/sv/privacy).

Behöver du hjälp? [contact@remote-ha.com](mailto:contact@remote-ha.com). Skicka aldrig din token i en supportförfrågan.
