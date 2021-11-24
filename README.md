# spotty

:warning: **Mac only.** :warning:

Easy access from your terminal to your favorite Spotify songs.

No API token required, just the Spotify client.

## Description

Spotify exposes a fairly straightforward AppleScript API,
though it's cumbersome to use the `osascript` CLI.
This simplifies that interaction.

`spotty`'s primary goal is to allow you to `save` the current song to a DB,
so you can `play` it later.

If you forget which songs you `save`d, you can always `list` them.

## Usage

First, clone this repository and `go install`!

Say I'm listening to **Grateful Dead - Touch of Grey**...

From the terminal, I'd type `spotty save` to save it.

Later, I could run `spotty list` to list it,
along with my other saved songs.

I would choose to play the song with either:
1. `spotty play touch`
1. `spotty play touch of grey`
1. `spotty touch%grey`

It's powered by a SQL database, so I could use whatever is compatible with `LIKE`.

I could skip to the next track in context with `spotty next`.

Or toggle play/pause with `spotty pp`.

To get the current track, `spotty get`.
