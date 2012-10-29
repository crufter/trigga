trigga
======

Pub/sub messaging server. Code is very minimalistic, small, it has no dependencies apart from the Go stdlib for the time being.
Currently supports only 3 commands: publish, subscribe, unsubscribe.
Performance is approx 1500 messages per second with 8 subscribers on an 2x2500MHz AMD.
Messages are fire and forget.

I plan to use it as a communication channel of websockets based applications, triggering noncrucial events and sending very short messages, enchancing the user experience of
the web app/website.

No distributed bussiness yet, although it will be kinda trivial to add. Maybe if I find the time...

Drivers
---
- Go: https://github.com/opesun/gotrigga