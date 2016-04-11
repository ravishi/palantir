palantir
========

It means *seeing stone* in *sindarin*. It's one of those orbs people look into to see the future, the past and all kinds of crazy stuff.

This one is a *media server*. It's purpose is to index your media library (movies, TV shows and music), let you manage it and do all kinds of crazy stuff with it. Stuff like automatically downloading subtitles from the internet, playing from any device, who knows. It's a small experiment.


Architecture
------------

There's going to be a *REST-like* server written in go and a client side application written in React. It'll probably use a file system folder to store it's data. I plan to use quite a bit of sqlite.


Milestones
----------

- [ ] Setup development environment
- [ ] Index media
- [ ] Fetch media metadata
- [ ] Download subtitles :D
- [ ] Index more media, fetch more metadata, make better use of it
- [ ] Improve subtitles system. It should be easy to plug different *subtitle fetcher* in the system. I really hope to find a way to make it easy to 1) automatically fetch the best subtitles from various sources; 2) if needed, select select a specific subtitle; 3) if needed, search for subtitles in all providers; 4) if needed, just send it's own subtitle; 5) adjust subtitle offsets on the fly and save it for later use. It's gonnab be fucking awesome :')
- [ ] Maybe play it?
- [ ] Make the media source system pluggable so we can use media from all kinds of sources. Maybe it could stream torrents directly?
    
Final words
-----------

That's it, I'm already crying over here.
