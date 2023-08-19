# Ugle
![ugle](https://github.com/JustAHippo/Ugle/assets/82006314/e90c94a8-54ef-45ac-8fc7-7620a0226e47)

A search engine for the ucanet!

The [ucanet](https://ucanet.net), or the "U Create a Net" is an alternate DNS made to recreate the old internet from the ground up.

Before you contribute to Ugle!, check out [the ucanet GitHub page](https://github.com/ucanet)!
## What's done
- [x] Basic old Google UI
- [x] Domain name search
- [x] Web scraping/crawling
- [x] Turn registry.txt into a database
- [x] Cache site descriptions in database
- [x] Make API with cached indexing
## What needs work
- [ ] Use new API in frontend
- [ ] Paginate results
## Building and Usage
MongoDB is now required along with Go 1.20.
Edit [db.go](https://github.com/JustAHippo/Ugle/blob/main/db/db.go) to reflect your Mongo URI string

Then run:
```console
user@computer:~$ go build Ugle
```
And you're good to go!
## Contributing
If you're interested in helping out with the project of Ugle! or the [ucanet](https://ucanet.net) as a whole, getting in involved with the community, or registering an ucanet domain, join the [Discord](https://discord.gg/3mjrESssB3) and try it out!

## Credits
Thank you to the ucanet contributors for creating such a unique and interesting product and I'm very happy to be a part of the ucanet!
