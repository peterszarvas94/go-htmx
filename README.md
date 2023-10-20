# HTTG

- Htmx for interactivity
- Turso db
- Tailwind css for styling
- Go webserwer with html templates

## Commands

- `make dev` - start dev server on port 3000
- `make tw-watch` - generate tailwind styles for dev

### Prerequisites

To run make commands, you need to have:

- gin
- tailwindcss cli

## Todo

- [x] db connection
- [x] env
- [x] add jwt
- [x] rewrite login / signup with form (no boost)
- [x] clean up errors
- [x] signed in/out UI diff
- [x] only logged in can create todo
- [x] input sanitization and validation
- [x] logger
    - [ ] implement on every func
- [x] rework 404 handling
    - [ ] implement on every GET
- [ ] better route handling
- [ ] client side templating
    - [x] render <template> elements
    - [x] create and use variables
    - [x] allow to run js code
    - [ ] find a way to intersect network requests
    - [ ] what else?
