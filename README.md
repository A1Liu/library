# Web Server
A simple web server I'm using to learn the Go programming language.

## Code Responsibilities
- `/database` code must ensure that data taken in and out of the database is in
  valid; this means that it must be well defined, but not necessarily correct.
- `/web` code must ensure that data given to `/database` code is correct; that means
  performing security checks, preventing nulls, etc.
- `/models` code must ensure that type construction can be done correctly.


## Annotations

- `@Responsibility` - Is this code responsible for the functionality it's doing?
  is it responsible for more?
- `@TODO` - Needs to be done
- `@Performance` - There's some benchmarking that can be done here
