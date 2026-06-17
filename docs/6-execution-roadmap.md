# Extension roadmap

This project might be improved in several ways. None of the extensions listed below are required for the core
functionality to work. It's a list of possible improvements that could be implemented in the future.

## Telegram WebApp

Additional info (price, route stations, ticket office working hours, etc.) might be provided. But it is impossible
to present all this in one small message. It might be done
using [the WebApps](https://core.telegram.org/bots/features#web-apps) feature
and [GitHub Pages](https://github.com/revenkroz/telegram-web-app-bot-example) as a host.

### What it requires

Telegram's WebApp feature opens a web page
inside Telegram's in-app browser, launched from a button in a normal message. A static frontend hosted on GitHub Pages
(zero additional infrastructure cost) can render this. A new backend endpoint is required for processing application
requests.

### What changes

Adds a frontend codebase and corresponding backend API. The frontend might be a single-page app in the first iteration
and can be extended later to a multipage app with complex UI.

### What it costs

- Increases backend API surface area
- Developing, maintaining, and testing frontend codebase with unknown complexity

## Other platforms - Viber and WhatsApp

Telegram bots API is powerful, but Telegram itself is mostly used by foreigners. In Montenegro
Viber is a number one messenger app, and WhatsApp is the best option for any other tourists or foreigners, who don't use
Telegram. It is a good idea to port this solution to other platforms to reach the audience.

ZPCG has a [Viber chat](https://www.zcg-prevoz.me/Informisanje-o-redu-voznje-vozova-putem-aplikacije-Viber.html)
to get info about the timetable. But it replies to messages only on working days from 07:00 to 15:00, probably
because there is a person behind the scenes who is answering all the questions. I think it is a good idea to provide the
same service but working fast and 24/7.

### What it requires

Requires porting the existing chat logic to Viber and WhatsApp. Nothing about the algorithm, the domain model, or the
timetable pipeline is Telegram-specific - all of that is already platform-agnostic.

### What changes

Each platform needs its own bot registration and infrastructure (or at least its own API implementation).
Some Telegram-specific features like slash commands and inline buttons are not available on other platforms. This means
rewriting the chat logic to be platform-agnostic.

### What it costs:

- More budget - Viber API for bots is not free.
- More infrastructure - a new bot registration, Cloud Run services, and CI/CD pipelines.
- More maintenance burden - each version of the bot needs to be tested on every platform.

## Keyboard (rejected)

"I prefer clicking buttons over a raw keyboard input" - this is a common feedback I get from users. Implementing
a keyboard is a great way to improve UX. But how can the implementation still be stateless and fast enough then?

### What it requires

Telegram's [Force Reply](https://core.telegram.org/bots/api#forcereply) mechanism solves this
without server-side state. When the bot replies to a message with `force_reply`, Telegram's client tracks which message
the user is responding to and includes that context in the next webhook payload. The server can reconstruct "this is
step two of a station-selection flow" entirely from the incoming request. Nothing has to be stored between turns.

### What changes

The message-handling logic gains another code path - answering a reply-to-message.
Still no persistence layer, no session table, no TTL cleanup job. The stateless constraint is preserved exactly as is.
The state simply lives in user's Telegram chat.

### What it costs

- A new webhook handler that needs to branch on whether a message is a reply.
- Implementing processing logic with steps like:
    - Get a user's input
  - If it doesn't have two stations names, reply with a keyboard for the second station selection
    - Receive the second station name
    - If everything is correct, send a message with the timetable
- Maintaining and testing the keyboard
- Monitoring not a single message with a reply, but a whole chain from the first message to the last reply.

This feature is not critical for the core functionality, might not be as useful as the WebApp, and will significantly
increase the complexity of the logic and the maintenance burden.

Not a good candidate. Telegram WebApp might be a better option to solve this problem, since it can host a station
picker in a more user-friendly way.
